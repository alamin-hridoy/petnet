package handler

import (
	"html/template"
	"net/http"
	"strings"

	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"

	fpb "brank.as/petnet/gunk/dsa/v2/file"
)

type (
	FinancialInfoForm struct {
		CSRFField              template.HTML
		CSRFFieldValue         string
		FinancialStatements    []string
		BankStatements         []string
		FinancialStatementURLs string
		BankStatementURLs      string
		FinancialStatementID   string
		BankStatementID        string
		Errors                 map[string]error
		PresetPermission       map[string]map[string]bool
		ServiceRequest         bool
		SaveDraft              string
		SaveContinue           string
	}
)

func (s *Server) getFinancialInfo(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()

	sess, err := s.sess.Get(r, sessionCookieName)
	if err != nil {
		log.WithError(err).Error("fetching session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	uid := sess.Values[sessionUserID]
	if uid == nil {
		log.Error("missing user id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	oid := sess.Values[sessionOrgID].(string)
	if oid == "" {
		log.Error("missing org id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	template := s.templates.Lookup("onboarding-financial.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	res, err := s.pf.ListFiles(ctx, &fpb.ListFilesRequest{
		OrgID: oid,
		Types: []fpb.UploadType{fpb.UploadType_FinancialStatement, fpb.UploadType_BankStatement},
	})
	if err != nil {
		logging.WithError(err, log).Error("listing files")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	form := fiProtoToForm(res.GetFileUploads())
	form.CSRFField = csrf.TemplateField(r)
	form.CSRFFieldValue = csrf.Token(r)
	etd := s.getEnforceTemplateData(ctx)
	form.PresetPermission = etd.PresetPermission
	form.ServiceRequest = etd.ServiceRequests
	if err := template.Execute(w, form); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) postFinancialInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)

	sess, err := s.sess.Get(r, sessionCookieName)
	if err != nil {
		log.WithError(err).Error("fetching session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	uid, ok := sess.Values[sessionUserID].(string)
	if !ok {
		log.Error("missing user id in session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	oid, ok := sess.Values[sessionOrgID].(string)
	if !ok {
		log.Error("missing org id in session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		logging.WithError(err, log).Error("parsing form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	var form FinancialInfoForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	formErr := validation.Errors{}
	fs, bs := "FinancialStatements", "BankStatements"
	fileTypes := map[string]fpb.UploadType{
		"FinancialStatements": fpb.UploadType_FinancialStatement, "BankStatements": fpb.UploadType_BankStatement,
	}

	mu := map[string][]string{fs: {}, bs: {}}
	for k := range mu {
		err := s.validateMultiFileType(r, k, validedDocFileType)
		if err != nil {
			formErr[k] = err
		}
	}

	form.Errors = formErr
	form.CSRFField = csrf.TemplateField(r)
	form.CSRFFieldValue = csrf.Token(r)
	etd := s.getEnforceTemplateData(ctx)
	form.PresetPermission = etd.PresetPermission
	form.ServiceRequest = etd.ServiceRequests
	template := s.templates.Lookup("onboarding-financial.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if len(formErr) > 0 {
		if err := template.Execute(w, form); err != nil {
			log.Infof("error with template execution: %+v", err)
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}

	mful := map[fpb.UploadType][]string{}
	fns := map[string]map[string]string{}
	for k := range mu {
		fns[k] = map[string]string{}
		_, _, fuerr := r.FormFile(k)
		us, fn, err := s.storeMultiToGCS(r, k, oid)
		if err != nil && fuerr == nil {
			logging.WithError(err, log).Error("storing " + k)
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		mu[k] = us
		if len(us) == 0 {
			continue
		}
		mful[fileTypes[k]] = us
		fns[fileTypes[k].String()] = fn
	}
	fileUploads := []*fpb.FileUpload{}
	for u, v := range mful {
		fileUploads = append(fileUploads, &fpb.FileUpload{
			UserID:    uid,
			OrgID:     oid,
			FileNames: v,
			Type:      u,
			FileName:  fns[u.String()],
		})
	}
	// todo(robin): if we decide to add another bucket make sure the bucket name is saved.
	if len(fileUploads) == 0 {
		http.Redirect(w, r, accountInfoPath, http.StatusSeeOther)
		return
	}
	if form.SaveDraft == "SaveDraft" {
		if _, err := s.pf.UpsertFiles(ctx, &fpb.UpsertFilesRequest{
			FileUploads: fileUploads,
		}); err != nil {
			logging.WithError(err, log).Error("creating files")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, loginPath, http.StatusSeeOther)
		return
	}
	if form.SaveContinue == "SaveContinue" {
		if _, err := s.pf.UpsertFiles(ctx, &fpb.UpsertFilesRequest{
			FileUploads: fileUploads,
		}); err != nil {
			logging.WithError(err, log).Error("creating files")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, accountInfoPath, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, accountInfoPath, http.StatusSeeOther)
	return
}

func fiProtoToForm(fs []*fpb.FileUpload) FinancialInfoForm {
	fm := FinancialInfoForm{}
	for _, f := range fs {
		switch f.Type {
		case fpb.UploadType_FinancialStatement:
			fm.FinancialStatementURLs = strings.Join(f.GetFileNames(), ",")
			fm.FinancialStatementID = f.GetID()
		case fpb.UploadType_BankStatement:
			fm.BankStatementURLs = strings.Join(f.GetFileNames(), ",")
			fm.BankStatementID = f.GetID()
		}
	}
	return fm
}

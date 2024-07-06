package handler

import (
	"html/template"
	"net/http"

	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
	"google.golang.org/protobuf/types/known/timestamppb"

	fpb "brank.as/petnet/gunk/dsa/v2/file"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
)

type (
	DRPInfoForm struct {
		CSRFField        template.HTML
		Questionnaires   string
		Errors           map[string]error
		PresetPermission map[string]map[string]bool
		ServiceRequest   bool
	}
)

func (s *Server) getDRPInfo(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())

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

	template := s.templates.Lookup("onboarding-drp.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	etd := s.getEnforceTemplateData(r.Context())
	f := DRPInfoForm{
		Errors:           map[string]error{},
		PresetPermission: etd.PresetPermission,
		ServiceRequest:   etd.ServiceRequests,
	}
	f.CSRFField = csrf.TemplateField(r)
	if err := template.Execute(w, f); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) postDRPInfo(w http.ResponseWriter, r *http.Request) {
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

	var form DRPInfoForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	formErr := validation.Errors{}
	if err := validation.ValidateStruct(&form); err != nil {
		logging.WithError(err, log).Error("invalid request")
	}

	questionnaires := "Questionnaires"
	mu := map[string][]string{questionnaires: {}}
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
		fns[k] = fn
	}

	template := s.templates.Lookup("onboarding-drp.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	etd := s.getEnforceTemplateData(ctx)
	form.Errors = formErr
	form.CSRFField = csrf.TemplateField(r)
	form.PresetPermission = etd.PresetPermission
	form.ServiceRequest = etd.ServiceRequests
	if len(formErr) > 0 {
		if err := template.Execute(w, form); err != nil {
			log.Infof("error with template execution: %+v", err)
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}

	if _, err := s.pf.UpsertProfile(ctx, &ppb.UpsertProfileRequest{
		Profile: &ppb.OrgProfile{
			UserID:      uid,
			OrgID:       oid,
			DateApplied: timestamppb.Now(),
			Status:      ppb.Status_Pending, // todo(robin): should be set when application validated successfully and sent to api, adjust later
		},
	}); err != nil {
		logging.WithError(err, log).Error("creating profile")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	// todo(robin): if we decide to add another bucket make sure the bucket name is saved
	if _, err := s.pf.UpsertFiles(ctx, &fpb.UpsertFilesRequest{
		FileUploads: []*fpb.FileUpload{
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[questionnaires],
				Type:      fpb.UploadType_Questionnaire,
				FileName:  fns[questionnaires],
			},
		},
	}); err != nil {
		logging.WithError(err, log).Error("creating files")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, registerSuccessPath, http.StatusSeeOther)
}

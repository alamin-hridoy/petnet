package handler

import (
	"errors"
	"html/template"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	fpb "brank.as/petnet/gunk/dsa/v2/file"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gorilla/csrf"
)

type (
	BusinessInfoForm struct {
		CSRFField               template.HTML
		CSRFFieldValue          string
		OrgID                   string
		UserID                  string
		CompanyName             string
		StoreName               string
		PhoneNumber             string
		FaxNumber               string
		Website                 string
		CompanyEmail            string
		ContactPerson           string
		Position                string
		Address1                string
		City                    string
		State                   string
		PostalCode              string
		IDPhotoURLs             string
		PictureURLs             string
		NBIClearanceURLs        string
		CourtClearanceURL       string
		IncorporationPapersURLs string
		MayorsPermitURL         string
		IDPhotoID               string
		PictureID               string
		NBIClearanceID          string
		CourtClearanceID        string
		IncorporationPapersID   string
		MayorsPermitID          string
		NdaURLs                 string
		NdaID                   string
		TransactionTypes        []string
		Errors                  map[string]error
		PresetPermission        map[string]map[string]bool
		ServiceRequest          bool
		SaveDraft               string
		SaveContinue            string
	}
)

var (
	// deprecated ...
	validedImgFileType = []string{
		"image/jpeg", "image/png", "image/gif",
	}

	// deprecated ...
	validedDocFileType = []string{
		"image/jpeg", "image/png", "image/gif", "application/pdf",
	}

	validDocFilesMap = map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"image/gif":       true,
		"application/pdf": true,

		// TODO(vitthal): Fix issue for uploading doc/docx files to GCS gets converted to zip.
		"application/msword": true, // MS-word files (extension .doc)
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // Document extension .docx
	}
)

func (s *Server) getBusinessInfo(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	sess, err := s.sess.Get(r, sessionCookieName)
	if err != nil {
		log.WithError(err).Error("fetching session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	queryParams := r.URL.Query()
	oid, ok := sess.Values[sessionOrgID].(string)
	if !ok {
		oid, err = url.PathUnescape(queryParams.Get("org_id"))
		if err != nil {
			log.Error("missing org id")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}

	uid, ok := sess.Values[sessionUserID].(string)
	if !ok {
		uid, err = url.PathUnescape(queryParams.Get("user_id"))
		if err != nil {
			log.Error("missing user id")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}

	sess.Values[sessionUserID] = uid
	sess.Values[sessionOrgID] = oid
	if err := s.sess.Save(r, w, sess); err != nil {
		log.WithError(err).Error("saving session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	template := s.templates.Lookup("onboarding-business.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	pf, err := s.pf.GetProfile(r.Context(), &ppb.GetProfileRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}

	res, err := s.pf.ListFiles(ctx, &fpb.ListFilesRequest{
		OrgID: oid,
		Types: []fpb.UploadType{
			fpb.UploadType_IDPhoto,
			fpb.UploadType_Picture,
			fpb.UploadType_NBIClearance,
			fpb.UploadType_CourtClearance,
			fpb.UploadType_IncorporationPapers,
			fpb.UploadType_MayorsPermit,
			fpb.UploadType_NDA,
		},
	})
	if err != nil {
		logging.WithError(err, log).Error("listing files")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	transactionTypes := []string{}
	if pf != nil && pf.Profile != nil && pf.Profile.GetTransactionTypes() != "" {
		transactionTypes = strings.Split(pf.Profile.GetTransactionTypes(), ",")
	}

	etd := s.getEnforceTemplateData(ctx)
	form := biProtoToForm(pf.GetProfile().GetBusinessInfo(), res.GetFileUploads())
	form.CSRFField = csrf.TemplateField(r)
	form.CSRFFieldValue = csrf.Token(r)
	form.UserID = uid
	form.OrgID = oid
	form.TransactionTypes = transactionTypes
	form.PresetPermission = etd.PresetPermission
	form.ServiceRequest = etd.ServiceRequests
	if err := template.Execute(w, form); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) postBusinessInfo(w http.ResponseWriter, r *http.Request) {
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

	var form BusinessInfoForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	formErr := validation.Errors{}
	ts := strings.TrimSpace
	if err := validation.ValidateStruct(&form,
		validation.Field(&form.CompanyName, validation.Required),
		validation.Field(&form.PhoneNumber, is.Digit, validation.Required),
		validation.Field(&form.FaxNumber, is.Digit),
		validation.Field(&form.Website, is.URL),
		validation.Field(&form.CompanyEmail, is.Email, validation.Required),
		validation.Field(&form.ContactPerson, validation.Match(regexp.MustCompile("^[a-zA-Z_ ]*$")), validation.Required),
		validation.Field(&form.Position, validation.Match(regexp.MustCompile("^[a-zA-Z_ ]*$")), validation.Required),
		validation.Field(&form.Address1, validation.Required),
		validation.Field(&form.City, validation.Required),
		validation.Field(&form.State, validation.Required),
		validation.Field(&form.TransactionTypes, validation.Required),
		validation.Field(&form.PostalCode, validation.Required),
	); err != nil {
		if err, ok := (err).(validation.Errors); ok {
			if err["CompanyName"] != nil {
				formErr["CompanyName"] = errors.New("Company Name can not be empty")
			}
			if err["PhoneNumber"] != nil {
				formErr["PhoneNumber"] = errors.New("Phone number needs to be digits only")
			}
			if err["FaxNumber"] != nil {
				formErr["FaxNumber"] = errors.New("Fax number should be digit")
			}
			if err["Website"] != nil {
				formErr["Website"] = errors.New("Website need to be a valid URL")
			}
			if err["CompanyEmail"] != nil {
				formErr["CompanyEmail"] = errors.New("Should be a valid Email")
			}
			if err["ContactPerson"] != nil {
				formErr["ContactPerson"] = errors.New("Only allow alphabet and space")
			}
			if err["Position"] != nil {
				formErr["Position"] = errors.New("Only allow alphabet and space")
			}
			if err["Address1"] != nil {
				formErr["Address1"] = errors.New("Building/Stree is required")
			}
			if err["City"] != nil {
				formErr["City"] = errors.New("City is required")
			}
			if err["State"] != nil {
				formErr["State"] = errors.New("Province is required")
			}
			if err["PostalCode"] != nil {
				formErr["PostalCode"] = errors.New("Zip/Postalcode is required")
			}
			if err["TransactionTypes"] != nil {
				formErr["TransactionTypes"] = errors.New("Please select at least 1 Transaction Type to avail")
			}
		}
		logging.WithError(err, log).Error("invalid request")
	}
	cc, mp := "cc", "mp"
	si := map[string]string{cc: "", mp: ""}
	validID, picture, nbi, ip, nda := "validid", "picture", "nbi", "ip", "nda"

	mu := map[string][]string{validID: {}, picture: {}, nbi: {}, ip: {}, nda: {}}
	fileTypes := map[string]fpb.UploadType{
		"cc": fpb.UploadType_CourtClearance, "mp": fpb.UploadType_MayorsPermit,
		"validid": fpb.UploadType_IDPhoto, "picture": fpb.UploadType_Picture, "nbi": fpb.UploadType_NBIClearance, "ip": fpb.UploadType_IncorporationPapers, "nda": fpb.UploadType_NDA,
	}
	requiredDocErr := errors.New("This is a required document, please upload a file")
	for fileField := range si {
		f, _, err := r.FormFile(fileField)
		uploaded := r.FormValue("uploaded-" + fileField)
		// error if not selected and not uploaded
		if err != nil && uploaded == "" {
			// file is required
			formErr[fileField] = requiredDocErr
			continue
		}

		if err != nil {
			// no need to validate if file not present in request
			continue
		}

		defer f.Close()

		err = validateFileType(f, validDocFilesMap)
		if err != nil {
			formErr[fileField] = err
		}
	}

	for fileField := range mu {
		uploaded := r.FormValue("uploaded-" + fileField)
		fileHeaders, ok := r.MultipartForm.File[fileField]
		// error if not selected and not uploaded
		if !ok && uploaded == "" {
			// file is required
			formErr[fileField] = requiredDocErr
			continue
		}

		if !ok {
			// no need to validate if file not present in request
			continue
		}

		err := validateFileHeadersForDuplicateAndFileType(fileHeaders, validDocFilesMap)
		if err != nil {
			formErr[fileField] = err
		}
	}

	if form.StoreName == "" {
		form.StoreName = ts(form.CompanyName)
	}

	template := s.templates.Lookup("onboarding-business.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	transactionTypes := strings.TrimSpace(strings.Join(form.TransactionTypes, ","))
	if transactionTypes == "" {
		formErr["TransactionTypes"] = errors.New("Please select at least 1 Transaction Type to avail")
	}
	form.Errors = formErr
	form.CSRFField = csrf.TemplateField(r)
	form.CSRFFieldValue = csrf.Token(r)
	etd := s.getEnforceTemplateData(ctx)
	form.PresetPermission = etd.PresetPermission
	form.ServiceRequest = etd.ServiceRequests
	form.UserID = mw.GetUserID(ctx)
	if len(formErr) > 0 {
		// TODO(vitthal): Might need to render uploaded files after errors?
		if err := template.Execute(w, form); err != nil {
			log.Infof("error with template execution: %+v", err)
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}
	fs := map[string]map[string]string{}
	sful := map[fpb.UploadType]string{}
	for k := range si {
		fs[k] = map[string]string{}
		_, _, fuerr := r.FormFile(k)
		u, fsn, err := s.storeSingleToGCS(r, k, oid)
		if err != nil && fuerr == nil {
			logging.WithError(err, log).Error("storing " + k)
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		if u == "" {
			continue
		}
		sful[fileTypes[k]] = u
		fs[fileTypes[k].String()] = fsn
	}
	fns := map[string]map[string]string{}
	mful := map[fpb.UploadType][]string{}
	for k := range mu {
		fns[k] = map[string]string{}
		_, _, fuerr := r.FormFile(k)
		us, fn, err := s.storeMultiToGCS(r, k, oid)
		if err != nil && fuerr == nil {
			logging.WithError(err, log).Error("storing " + k)
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		if len(us) == 0 {
			continue
		}
		mful[fileTypes[k]] = us
		fns[fileTypes[k].String()] = fn
	}

	if form.SaveDraft == "SaveDraft" {
		if _, err := s.pf.UpsertProfile(ctx, &ppb.UpsertProfileRequest{
			Profile: &ppb.OrgProfile{
				UserID:           uid,
				OrgID:            oid,
				TransactionTypes: transactionTypes,
				BusinessInfo: &ppb.BusinessInfo{
					CompanyName:   ts(form.CompanyName),
					StoreName:     ts(form.StoreName),
					PhoneNumber:   ts(form.PhoneNumber),
					FaxNumber:     ts(form.FaxNumber),
					Website:       ts(form.Website),
					CompanyEmail:  ts(form.CompanyEmail),
					ContactPerson: ts(form.ContactPerson),
					Position:      ts(form.Position),
					Address: &ppb.Address{
						Address1:   ts(form.Address1),
						City:       ts(form.City),
						State:      ts(form.State),
						PostalCode: ts(form.PostalCode),
					},
				},
			},
		}); err != nil {
			logging.WithError(err, log).Error("creating profile")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}

	if form.SaveDraft == "SaveDraft" {
		if _, err := s.pf.UpsertProfile(ctx, &ppb.UpsertProfileRequest{
			Profile: &ppb.OrgProfile{
				UserID:           uid,
				OrgID:            oid,
				Status:           ppb.Status_Incomplete,
				TransactionTypes: transactionTypes,
				BusinessInfo: &ppb.BusinessInfo{
					CompanyName:   ts(form.CompanyName),
					StoreName:     ts(form.StoreName),
					PhoneNumber:   ts(form.PhoneNumber),
					FaxNumber:     ts(form.FaxNumber),
					Website:       ts(form.Website),
					CompanyEmail:  ts(form.CompanyEmail),
					ContactPerson: ts(form.ContactPerson),
					Position:      ts(form.Position),
					Address: &ppb.Address{
						Address1:   ts(form.Address1),
						City:       ts(form.City),
						State:      ts(form.State),
						PostalCode: ts(form.PostalCode),
					},
				},
			},
		}); err != nil {
			logging.WithError(err, log).Error("creating profile")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, loginPath, http.StatusSeeOther)
	}

	if form.SaveContinue == "SaveContinue" {
		if _, err := s.pf.UpsertProfile(ctx, &ppb.UpsertProfileRequest{
			Profile: &ppb.OrgProfile{
				UserID:           uid,
				OrgID:            oid,
				TransactionTypes: transactionTypes,
				BusinessInfo: &ppb.BusinessInfo{
					CompanyName:   ts(form.CompanyName),
					StoreName:     ts(form.StoreName),
					PhoneNumber:   ts(form.PhoneNumber),
					FaxNumber:     ts(form.FaxNumber),
					Website:       ts(form.Website),
					CompanyEmail:  ts(form.CompanyEmail),
					ContactPerson: ts(form.ContactPerson),
					Position:      ts(form.Position),
					Address: &ppb.Address{
						Address1:   ts(form.Address1),
						City:       ts(form.City),
						State:      ts(form.State),
						PostalCode: ts(form.PostalCode),
					},
				},
			},
		}); err != nil {
			logging.WithError(err, log).Error("creating profile")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}

	// todo(robin): if we decide to add another bucket make sure the bucket name is saved
	// Todo:(vitthal): Need to append already uploaded files or will it be overwritten?

	fileUploads := []*fpb.FileUpload{}
	for u, v := range sful {
		fileUploads = append(fileUploads, &fpb.FileUpload{
			UserID:    uid,
			OrgID:     oid,
			FileNames: []string{v},
			Type:      u,
			FileName:  fs[u.String()],
		})
	}
	for u, v := range mful {
		fileUploads = append(fileUploads, &fpb.FileUpload{
			UserID:    uid,
			OrgID:     oid,
			FileNames: v,
			Type:      u,
			FileName:  fns[u.String()],
		})
	}

	fileReq := fpb.UpsertFilesRequest{
		FileUploads: fileUploads,
	}

	if _, err := s.pf.UpsertFiles(ctx, &fileReq); err != nil {
		logging.WithError(err, log).Error("creating files")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, financialInfoPath, http.StatusSeeOther)
}

func biProtoToForm(bip *ppb.BusinessInfo, fs []*fpb.FileUpload) BusinessInfoForm {
	bi := BusinessInfoForm{}
	bi.CompanyName = bip.GetCompanyName()
	bi.StoreName = bip.GetStoreName()
	bi.PhoneNumber = bip.GetPhoneNumber()
	bi.FaxNumber = bip.GetFaxNumber()
	bi.Website = bip.GetWebsite()
	bi.CompanyEmail = bip.GetCompanyEmail()
	bi.ContactPerson = bip.GetContactPerson()
	bi.Position = bip.GetPosition()
	bi.Address1 = bip.GetAddress().GetAddress1()
	bi.City = bip.GetAddress().GetCity()
	bi.State = bip.GetAddress().GetState()
	bi.PostalCode = bip.GetAddress().GetPostalCode()
	for _, f := range fs {
		switch f.Type {
		case fpb.UploadType_IDPhoto:
			bi.IDPhotoURLs = strings.Join(f.GetFileNames(), ",")
			bi.IDPhotoID = f.GetID()
		case fpb.UploadType_Picture:
			bi.PictureURLs = strings.Join(f.GetFileNames(), ",")
			bi.PictureID = f.GetID()
		case fpb.UploadType_NBIClearance:
			bi.NBIClearanceURLs = strings.Join(f.GetFileNames(), ",")
			bi.NBIClearanceID = f.GetID()
		case fpb.UploadType_CourtClearance:
			bi.CourtClearanceURL = strings.Join(f.GetFileNames(), ",")
			bi.CourtClearanceID = f.GetID()
		case fpb.UploadType_IncorporationPapers:
			bi.IncorporationPapersURLs = strings.Join(f.GetFileNames(), ",")
			bi.IncorporationPapersID = f.GetID()
		case fpb.UploadType_MayorsPermit:
			bi.MayorsPermitURL = strings.Join(f.GetFileNames(), ",")
			bi.MayorsPermitID = f.GetID()
		case fpb.UploadType_NDA:
			bi.NdaURLs = strings.Join(f.GetFileNames(), ",")
			bi.NdaID = f.GetID()
		}
	}
	return bi
}

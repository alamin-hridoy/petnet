package handler

import (
	"errors"
	"html/template"
	"net/http"
	"strings"

	ppf "brank.as/petnet/gunk/dsa/v1/user"
	fpb "brank.as/petnet/gunk/dsa/v2/file"
	sVcpb "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
)

type DsaUploadReqFormMicroInsurance struct {
	CSRFField        template.HTML
	Errors           map[string]error
	PresetPermission map[string]map[string]bool
	ServiceRequest   bool
	MI_mbfForms      string
	MI_mbfID         string
	MI_ndaForms      string
	MI_ndaID         string
	MI_secForms      string
	MI_secID         string
	MI_gisForms      string
	MI_gisID         string
	MI_afsForms      string
	MI_afsID         string
	MI_birForms      string
	MI_birID         string
	MI_scbForms      string
	MI_scbID         string
	MI_viaForms      string
	MI_viaID         string
	MI_moaForms      string
	MI_moaID         string
	SaveDraft        string
	SaveContinue     string
	CSRFFieldValue   string
	UserInfo         *User
	CompanyName      string
}

var validedReqFileTypeMicroIn = map[string]bool{
	"image/jpeg":         true,
	"image/png":          true,
	"image/gif":          true,
	"application/pdf":    true,
	"application/msword": true, // MS-word files (extension .doc)
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // Document extension .docx
}

func (s *Server) getDsaReqDocsMicroInsurance(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	uid := mw.GetUserID(ctx)
	if uid == "" {
		log.Error("missing user id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	oid := mw.GetOrgID(ctx)
	if oid == "" {
		log.Error("missing org id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	template := s.templates.Lookup("dsa-req-docs-micro-insurance.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	etd := s.getEnforceTemplateData(ctx)
	res, _ := s.pf.ListFiles(ctx, &fpb.ListFilesRequest{
		OrgID: oid,
		Types: []fpb.UploadType{fpb.UploadType_MI_mbf, fpb.UploadType_MI_nda, fpb.UploadType_MI_sec, fpb.UploadType_MI_gis, fpb.UploadType_MI_afs, fpb.UploadType_MI_bir, fpb.UploadType_MI_scb, fpb.UploadType_MI_via, fpb.UploadType_MI_moa},
	})

	data := ReqUploadProtoToFormMicroIns(res)
	data.CSRFField = csrf.TemplateField(r)
	data.CSRFFieldValue = csrf.Token(r)
	data.PresetPermission = etd.PresetPermission
	data.ServiceRequest = etd.ServiceRequests
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	data.UserInfo = &usrInfo.UserInfo
	data.CompanyName = usrInfo.CompanyName
	data.UserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, data); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	return
}

func (s *Server) postDsaReqDocsMicroInsurance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	uid := mw.GetUserID(ctx)

	if uid == "" {
		log.Error("missing user id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	oid := mw.GetOrgID(ctx)
	if oid == "" {
		log.Error("missing org id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		logging.WithError(err, log).Error("parsing form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if _, err := s.pf.AddServiceRequest(ctx, &sVcpb.AddServiceRequestRequest{
		OrgID:       oid,
		Type:        sVcpb.ServiceType_MICROINSURANCE,
		Partners:    []string{"RuralNet"},
		AllPartners: true,
	}); err != nil {
		logging.WithError(err, log).Error("Add Service Request Error")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	var form DsaUploadReqForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	formErr := validation.Errors{}
	MI_mbf, MI_nda, MI_sec, MI_gis, MI_afs, MI_bir, MI_scb, MI_via, MI_moa := "MI_mbf", "MI_nda", "MI_sec", "MI_gis", "MI_afs", "MI_bir", "MI_scb", "MI_via", "MI_moa"
	mu := map[string][]string{MI_mbf: {}, MI_nda: {}, MI_sec: {}, MI_gis: {}, MI_afs: {}, MI_bir: {}, MI_scb: {}, MI_via: {}, MI_moa: {}}
	for k := range mu {
		f, _, err := r.FormFile(k)
		if err != nil {
			continue
		}
		err = validateFileType(f, validedReqFileTypeMicroIn)
		if err != nil {
			formErr[k] = err
		}
	}

	requiredDocErr := errors.New("This is a required document, please upload a file")
	res, _ := s.pf.ListFiles(ctx, &fpb.ListFilesRequest{
		OrgID: oid,
		Types: []fpb.UploadType{fpb.UploadType_MI_mbf, fpb.UploadType_MI_nda, fpb.UploadType_MI_sec, fpb.UploadType_MI_gis, fpb.UploadType_MI_afs, fpb.UploadType_MI_bir, fpb.UploadType_MI_scb, fpb.UploadType_MI_via, fpb.UploadType_MI_moa},
	})

	ful := map[string][]string{}
	fulV := map[string][]string{}
	if len(res.GetFileUploads()) > 0 {
		for _, v := range res.GetFileUploads() {
			gt := v.GetType().String()
			if len(v.GetFileNames()) > 0 {
				fulV[gt] = v.GetFileNames()
				continue
			}
			ful[gt] = []string{}
		}
	}
	if len(ful) == 0 {
		ful = mu
	}
	if form.SaveContinue == "SaveContinue" {
		for fileField := range ful {
			uploaded := r.FormValue("uploaded-" + fileField)
			_, ok := r.MultipartForm.File[fileField]
			if !ok && uploaded == "" {
				formErr[fileField] = requiredDocErr
				continue
			}
			if !ok {
				continue
			}
			f, _, err := r.FormFile(fileField)
			if err != nil {
				continue
			}
			err = validateFileType(f, validedReqFileTypeMicroIn)
			if err != nil {
				formErr[fileField] = err
			}
		}
	}
	fns := map[string]map[string]string{}
	for k := range mu {
		fns[k] = map[string]string{}
		if formErr[k] != nil {
			continue
		}
		_, _, fuerr := r.FormFile(k)
		us, fn, err := s.storeMultiToGCS(r, k, oid)
		if err != nil && fuerr == nil {
			logging.WithError(err, log).Error("storing " + k)
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		if len(fulV[k]) > 0 {
			us = append(us, fulV[k]...)
		}
		mu[k] = us
		fns[k] = fn
	}
	// todo(robin): if we decide to add another bucket make sure the bucket name is saved.
	ufr, err := s.pf.UpsertFiles(ctx, &fpb.UpsertFilesRequest{
		FileUploads: []*fpb.FileUpload{
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[MI_mbf],
				Type:      fpb.UploadType_MI_mbf,
				FileName:  fns[MI_mbf],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[MI_nda],
				Type:      fpb.UploadType_MI_nda,
				FileName:  fns[MI_nda],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[MI_sec],
				Type:      fpb.UploadType_MI_sec,
				FileName:  fns[MI_sec],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[MI_gis],
				Type:      fpb.UploadType_MI_gis,
				FileName:  fns[MI_gis],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[MI_afs],
				Type:      fpb.UploadType_MI_afs,
				FileName:  fns[MI_afs],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[MI_bir],
				Type:      fpb.UploadType_MI_bir,
				FileName:  fns[MI_bir],
			},

			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[MI_scb],
				Type:      fpb.UploadType_MI_scb,
				FileName:  fns[MI_scb],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[MI_via],
				Type:      fpb.UploadType_MI_via,
				FileName:  fns[MI_via],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[MI_moa],
				Type:      fpb.UploadType_MI_moa,
				FileName:  fns[MI_moa],
			},
		},
	})
	if err != nil {
		logging.WithError(err, log).Error("creating files")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	for _, vv := range ufr.GetFileUploads() {
		gt := vv.GetType().String()
		if _, err := s.pf.AddUploadSvcRequest(ctx, &sVcpb.AddUploadSvcRequestRequest{
			OrgID:    oid,
			Partner:  "ALL",
			SvcName:  sVcpb.ServiceType_MICROINSURANCE.String(),
			FileType: gt,
			FileID:   vv.GetID(),
			CreateBy: uid,
		}); err != nil {
			if err == storage.Conflict {
				if _, err := s.pf.UpdateUploadSvcRequest(ctx, &sVcpb.UpdateUploadSvcRequestRequest{
					OrgID:    oid,
					Partner:  "ALL",
					SvcName:  sVcpb.ServiceType_MICROINSURANCE.String(),
					FileType: gt,
					FileID:   vv.GetID(),
					CreateBy: uid,
					VerifyBy: uid,
				}); err != nil {
					log.Infof("error with Update Upload Svc Request: %+v", err)
					http.Redirect(w, r, errorPath, http.StatusSeeOther)
					return
				}
			}
		}
	}
	if len(formErr) > 0 {
		res, _ := s.pf.ListFiles(ctx, &fpb.ListFilesRequest{
			OrgID: oid,
			Types: []fpb.UploadType{fpb.UploadType_MI_mbf, fpb.UploadType_MI_nda, fpb.UploadType_MI_sec, fpb.UploadType_MI_gis, fpb.UploadType_MI_afs, fpb.UploadType_MI_bir, fpb.UploadType_MI_scb, fpb.UploadType_MI_via, fpb.UploadType_MI_moa},
		})

		form := ReqUploadProtoToFormMicroIns(res)
		form.Errors = formErr
		form.CSRFField = csrf.TemplateField(r)
		form.CSRFFieldValue = csrf.Token(r)
		etd := s.getEnforceTemplateData(ctx)
		form.PresetPermission = etd.PresetPermission
		form.ServiceRequest = etd.ServiceRequests
		usrInfo := s.GetUserInfoFromCookie(w, r, false)

		form.UserInfo = &usrInfo.UserInfo
		gp, err := s.pf.GetUserProfile(ctx, &ppf.GetUserProfileRequest{
			UserID: uid,
		})
		if err != nil {
			log.Error("failed to get profile")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		form.UserInfo.ProfileImage = gp.GetProfile().ProfilePicture
		template := s.templates.Lookup("dsa-req-docs-micro-insurance.html")
		if template == nil {
			log.Error("unable to load template")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		if err := template.Execute(w, form); err != nil {
			log.Infof("error with template execution: %+v", err)
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}
	if form.SaveDraft == "SaveDraft" {
		if _, err := s.pf.SetStatusUploadSvcRequest(ctx, &sVcpb.SetStatusUploadSvcRequestRequest{
			OrgID:    oid,
			Partners: []string{"RuralNet"},
			SvcName:  sVcpb.ServiceType_MICROINSURANCE.String(),
			Status:   sVcpb.ServiceRequestStatus_REQDOCDRAFT,
		}); err != nil {
			log.Infof("error with set status Svc Request: %+v", err)
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, dsaServicesPath, http.StatusSeeOther)
		return
	}

	if !(len(formErr) > 0) && form.SaveContinue == "SaveContinue" {
		if _, err := s.pf.ApplyServiceRequest(ctx, &sVcpb.ApplyServiceRequestRequest{
			OrgID: oid,
			Type:  sVcpb.ServiceType_MICROINSURANCE,
		}); err != nil {
			if err != storage.NotFound {
				logging.WithError(err, log).Error("Apply Service Request Failed")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, dsaReqDocsMicroInsurance, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, dsaServicesPath, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, dsaReqDocsMicroInsurance, http.StatusSeeOther)
	return
}

func ReqUploadProtoToFormMicroIns(fs *fpb.ListFilesResponse) DsaUploadReqFormMicroInsurance {
	fm := DsaUploadReqFormMicroInsurance{}
	if fs != nil {
		for _, f := range fs.GetFileUploads() {
			switch f.Type {
			case fpb.UploadType_MI_mbf:
				fm.MI_mbfForms = strings.Join(f.GetFileNames(), ",")
				fm.MI_mbfID = f.GetID()
			case fpb.UploadType_MI_nda:
				fm.MI_ndaForms = strings.Join(f.GetFileNames(), ",")
				fm.MI_ndaID = f.GetID()
			case fpb.UploadType_MI_sec:
				fm.MI_secForms = strings.Join(f.GetFileNames(), ",")
				fm.MI_secID = f.GetID()
			case fpb.UploadType_MI_gis:
				fm.MI_gisForms = strings.Join(f.GetFileNames(), ",")
				fm.MI_gisID = f.GetID()
			case fpb.UploadType_MI_afs:
				fm.MI_afsForms = strings.Join(f.GetFileNames(), ",")
				fm.MI_afsID = f.GetID()
			case fpb.UploadType_MI_bir:
				fm.MI_birForms = strings.Join(f.GetFileNames(), ",")
				fm.MI_birID = f.GetID()
			case fpb.UploadType_MI_scb:
				fm.MI_scbForms = strings.Join(f.GetFileNames(), ",")
				fm.MI_scbID = f.GetID()
			case fpb.UploadType_MI_via:
				fm.MI_viaForms = strings.Join(f.GetFileNames(), ",")
				fm.MI_viaID = f.GetID()
			case fpb.UploadType_MI_moa:
				fm.MI_moaForms = strings.Join(f.GetFileNames(), ",")
				fm.MI_moaID = f.GetID()

			}
		}
	}
	return fm
}

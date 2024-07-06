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

type DsaUploadReqCicoForm struct {
	CSRFField        template.HTML
	Errors           map[string]error
	PresetPermission map[string]map[string]bool
	ServiceRequest   bool
	CICO_ndaForms    string
	CICO_ndaID       string
	CICO_sisForms    string
	CICO_sisID       string
	CICO_psfForms    string
	CICO_psfID       string
	CICO_psppForms   string
	CICO_psppID      string
	CICO_secForms    string
	CICO_secID       string
	CICO_gisForms    string
	CICO_gisID       string
	CICO_afsForms    string
	CICO_afsID       string
	CICO_birForms    string
	CICO_birID       string
	CICO_bspForms    string
	CICO_bspID       string
	CICO_amlForms    string
	CICO_amlID       string
	CICO_sccasForms  string
	CICO_sccasID     string
	CICO_vgidForms   string
	CICO_vgidID      string
	CICO_cpForms     string
	CICO_cpID        string
	CICO_moaForms    string
	CICO_moaID       string
	CICO_amlaForms   string
	CICO_amlaID      string
	CICO_mtppForms   string
	CICO_mtppID      string
	CICO_isForms     string
	CICO_isID        string
	CICO_eddForms    string
	CICO_eddID       string
	SaveDraft        string
	SaveContinue     string
	CSRFFieldValue   string
	UserInfo         *User
	CompanyName      string
}

var validedReqFileTypeCico = map[string]bool{
	"image/jpeg":         true,
	"image/png":          true,
	"image/gif":          true,
	"application/pdf":    true,
	"application/msword": true, // MS-word files (extension .doc)
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // Document extension .docx
}

func (s *Server) getDsaReqDocsCico(w http.ResponseWriter, r *http.Request) {
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
	template := s.templates.Lookup("dsa-req-docs-cico.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	etd := s.getEnforceTemplateData(ctx)
	res, _ := s.pf.ListFiles(ctx, &fpb.ListFilesRequest{
		OrgID: oid,
		Types: []fpb.UploadType{fpb.UploadType_CICO_nda, fpb.UploadType_CICO_sis, fpb.UploadType_CICO_psf, fpb.UploadType_CICO_pspp, fpb.UploadType_CICO_sec, fpb.UploadType_CICO_gis, fpb.UploadType_CICO_afs, fpb.UploadType_CICO_bir, fpb.UploadType_CICO_bsp, fpb.UploadType_CICO_aml, fpb.UploadType_CICO_sccas, fpb.UploadType_CICO_vgid, fpb.UploadType_CICO_cp, fpb.UploadType_CICO_moa, fpb.UploadType_CICO_amla, fpb.UploadType_CICO_mtpp, fpb.UploadType_CICO_is, fpb.UploadType_CICO_edd},
	})
	var availablePartners []string
	isPartnerDraft := false
	ress, _ := s.pf.ListServiceRequest(ctx, &sVcpb.ListServiceRequestRequest{
		OrgIDs: []string{oid},
		Types:  []sVcpb.ServiceType{sVcpb.ServiceType_CASHINCASHOUT},
	})
	if ress != nil {
		for _, vv := range ress.GetServiceRequst() {
			availablePartners = append(availablePartners, vv.Partner)
			if vv.Status == sVcpb.ServiceRequestStatus_PARTNERDRAFT {
				isPartnerDraft = true
			}
		}
	}
	if isPartnerDraft || len(availablePartners) == 0 {
		http.Redirect(w, r, dsaPrtCiCoSelPath, http.StatusSeeOther)
	}
	data := CicoReqUploadProtoToForm(res)
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

func (s *Server) postDsaReqDocsCico(w http.ResponseWriter, r *http.Request) {
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
	var form DsaUploadReqCicoForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	formErr := validation.Errors{}
	CICO_nda, CICO_sis, CICO_psf, CICO_pspp, CICO_sec, CICO_gis, CICO_afs, CICO_bir, CICO_bsp, CICO_aml, CICO_sccas, CICO_vgid, CICO_cp, CICO_moa, CICO_amla, CICO_mtpp, CICO_is, CICO_edd := "CICO_nda", "CICO_sis", "CICO_psf", "CICO_pspp", "CICO_sec", "CICO_gis", "CICO_afs", "CICO_bir", "CICO_bsp", "CICO_aml", "CICO_sccas", "CICO_vgid", "CICO_cp", "CICO_moa", "CICO_amla", "CICO_mtpp", "CICO_is", "CICO_edd"
	mu := map[string][]string{CICO_nda: {}, CICO_sis: {}, CICO_psf: {}, CICO_pspp: {}, CICO_sec: {}, CICO_gis: {}, CICO_afs: {}, CICO_bir: {}, CICO_bsp: {}, CICO_aml: {}, CICO_sccas: {}, CICO_vgid: {}, CICO_cp: {}, CICO_moa: {}, CICO_amla: {}, CICO_mtpp: {}, CICO_is: {}, CICO_edd: {}}
	for k := range mu {
		f, _, err := r.FormFile(k)
		if err != nil {
			continue
		}
		err = validateFileType(f, validedReqFileTypeCico)
		if err != nil {
			formErr[k] = err
		}
	}
	var availablePartners []string
	isPartnerDraft := false
	ress, _ := s.pf.ListServiceRequest(ctx, &sVcpb.ListServiceRequestRequest{
		OrgIDs: []string{oid},
		Types:  []sVcpb.ServiceType{sVcpb.ServiceType_CASHINCASHOUT},
	})
	if ress != nil {
		for _, vv := range ress.GetServiceRequst() {
			availablePartners = append(availablePartners, vv.Partner)
			if vv.Status == sVcpb.ServiceRequestStatus_PARTNERDRAFT {
				isPartnerDraft = true
			}
		}
	}
	if isPartnerDraft || len(availablePartners) == 0 {
		http.Redirect(w, r, dsaPrtCiCoSelPath, http.StatusSeeOther)
	}
	requiredDocErr := errors.New("This is a required document, please upload a file")
	res, _ := s.pf.ListFiles(ctx, &fpb.ListFilesRequest{
		OrgID: oid,
		Types: []fpb.UploadType{fpb.UploadType_CICO_nda, fpb.UploadType_CICO_sis, fpb.UploadType_CICO_psf, fpb.UploadType_CICO_pspp, fpb.UploadType_CICO_sec, fpb.UploadType_CICO_gis, fpb.UploadType_CICO_afs, fpb.UploadType_CICO_bir, fpb.UploadType_CICO_bsp, fpb.UploadType_CICO_aml, fpb.UploadType_CICO_sccas, fpb.UploadType_CICO_vgid, fpb.UploadType_CICO_cp, fpb.UploadType_CICO_moa, fpb.UploadType_CICO_amla, fpb.UploadType_CICO_mtpp, fpb.UploadType_CICO_is, fpb.UploadType_CICO_edd},
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
			if gt == "CICO_bsp" || gt == "CICO_aml" || gt == "CICO_edd" {
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
			if fileField == "CICO_bsp" || fileField == "CICO_aml" || fileField == "CICO_edd" {
				continue
			}
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
			err = validateFileType(f, validedReqFileTypeCico)
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
				FileNames: mu[CICO_nda],
				Type:      fpb.UploadType_CICO_nda,
				FileName:  fns[CICO_nda],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[CICO_sis],
				Type:      fpb.UploadType_CICO_sis,
				FileName:  fns[CICO_sis],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[CICO_psf],
				Type:      fpb.UploadType_CICO_psf,
				FileName:  fns[CICO_psf],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[CICO_pspp],
				Type:      fpb.UploadType_CICO_pspp,
				FileName:  fns[CICO_pspp],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[CICO_sec],
				Type:      fpb.UploadType_CICO_sec,
				FileName:  fns[CICO_sec],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[CICO_gis],
				Type:      fpb.UploadType_CICO_gis,
				FileName:  fns[CICO_gis],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[CICO_afs],
				Type:      fpb.UploadType_CICO_afs,
				FileName:  fns[CICO_afs],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[CICO_bir],
				Type:      fpb.UploadType_CICO_bir,
				FileName:  fns[CICO_bir],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[CICO_bsp],
				Type:      fpb.UploadType_CICO_bsp,
				FileName:  fns[CICO_bsp],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[CICO_aml],
				Type:      fpb.UploadType_CICO_aml,
				FileName:  fns[CICO_aml],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[CICO_sccas],
				Type:      fpb.UploadType_CICO_sccas,
				FileName:  fns[CICO_sccas],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[CICO_vgid],
				Type:      fpb.UploadType_CICO_vgid,
				FileName:  fns[CICO_vgid],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[CICO_cp],
				Type:      fpb.UploadType_CICO_cp,
				FileName:  fns[CICO_cp],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[CICO_moa],
				Type:      fpb.UploadType_CICO_moa,
				FileName:  fns[CICO_moa],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[CICO_amla],
				Type:      fpb.UploadType_CICO_amla,
				FileName:  fns[CICO_amla],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[CICO_mtpp],
				Type:      fpb.UploadType_CICO_mtpp,
				FileName:  fns[CICO_mtpp],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[CICO_is],
				Type:      fpb.UploadType_CICO_is,
				FileName:  fns[CICO_is],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[CICO_edd],
				Type:      fpb.UploadType_CICO_edd,
				FileName:  fns[CICO_edd],
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
			SvcName:  sVcpb.ServiceType_CASHINCASHOUT.String(),
			FileType: gt,
			FileID:   vv.GetID(),
			CreateBy: uid,
		}); err != nil {
			if err == storage.Conflict {
				if _, err := s.pf.UpdateUploadSvcRequest(ctx, &sVcpb.UpdateUploadSvcRequestRequest{
					OrgID:    oid,
					Partner:  "ALL",
					SvcName:  sVcpb.ServiceType_CASHINCASHOUT.String(),
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
			Types: []fpb.UploadType{fpb.UploadType_CICO_nda, fpb.UploadType_CICO_sis, fpb.UploadType_CICO_psf, fpb.UploadType_CICO_pspp, fpb.UploadType_CICO_sec, fpb.UploadType_CICO_gis, fpb.UploadType_CICO_afs, fpb.UploadType_CICO_bir, fpb.UploadType_CICO_bsp, fpb.UploadType_CICO_aml, fpb.UploadType_CICO_sccas, fpb.UploadType_CICO_vgid, fpb.UploadType_CICO_cp, fpb.UploadType_CICO_moa, fpb.UploadType_CICO_amla, fpb.UploadType_CICO_mtpp, fpb.UploadType_CICO_is, fpb.UploadType_CICO_edd},
		})
		form := CicoReqUploadProtoToForm(res)
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
		template := s.templates.Lookup("dsa-req-docs.html")
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
			Partners: availablePartners,
			SvcName:  sVcpb.ServiceType_CASHINCASHOUT.String(),
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
			Type:  sVcpb.ServiceType_CASHINCASHOUT,
		}); err != nil {
			if err != storage.NotFound {
				logging.WithError(err, log).Error("Apply Service Request Failed")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, dsaReqDocsCicoPath, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, dsaServicesPath, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, dsaReqDocsCicoPath, http.StatusSeeOther)
	return
}

func CicoReqUploadProtoToForm(fs *fpb.ListFilesResponse) DsaUploadReqCicoForm {
	fm := DsaUploadReqCicoForm{}
	if fs != nil {
		for _, f := range fs.GetFileUploads() {
			switch f.Type {
			case fpb.UploadType_CICO_nda:
				fm.CICO_ndaForms = strings.Join(f.GetFileNames(), ",")
				fm.CICO_ndaID = f.GetID()
			case fpb.UploadType_CICO_sis:
				fm.CICO_sisForms = strings.Join(f.GetFileNames(), ",")
				fm.CICO_sisID = f.GetID()
			case fpb.UploadType_CICO_psf:
				fm.CICO_psfForms = strings.Join(f.GetFileNames(), ",")
				fm.CICO_psfID = f.GetID()
			case fpb.UploadType_CICO_pspp:
				fm.CICO_psppForms = strings.Join(f.GetFileNames(), ",")
				fm.CICO_psppID = f.GetID()
			case fpb.UploadType_CICO_sec:
				fm.CICO_secForms = strings.Join(f.GetFileNames(), ",")
				fm.CICO_secID = f.GetID()
			case fpb.UploadType_CICO_gis:
				fm.CICO_gisForms = strings.Join(f.GetFileNames(), ",")
				fm.CICO_gisID = f.GetID()
			case fpb.UploadType_CICO_afs:
				fm.CICO_afsForms = strings.Join(f.GetFileNames(), ",")
				fm.CICO_afsID = f.GetID()
			case fpb.UploadType_CICO_bir:
				fm.CICO_birForms = strings.Join(f.GetFileNames(), ",")
				fm.CICO_birID = f.GetID()
			case fpb.UploadType_CICO_bsp:
				fm.CICO_bspForms = strings.Join(f.GetFileNames(), ",")
				fm.CICO_bspID = f.GetID()
			case fpb.UploadType_CICO_aml:
				fm.CICO_amlForms = strings.Join(f.GetFileNames(), ",")
				fm.CICO_amlID = f.GetID()
			case fpb.UploadType_CICO_sccas:
				fm.CICO_sccasForms = strings.Join(f.GetFileNames(), ",")
				fm.CICO_sccasID = f.GetID()
			case fpb.UploadType_CICO_vgid:
				fm.CICO_vgidForms = strings.Join(f.GetFileNames(), ",")
				fm.CICO_vgidID = f.GetID()
			case fpb.UploadType_CICO_cp:
				fm.CICO_cpForms = strings.Join(f.GetFileNames(), ",")
				fm.CICO_cpID = f.GetID()
			case fpb.UploadType_CICO_moa:
				fm.CICO_moaForms = strings.Join(f.GetFileNames(), ",")
				fm.CICO_moaID = f.GetID()
			case fpb.UploadType_CICO_amla:
				fm.CICO_amlaForms = strings.Join(f.GetFileNames(), ",")
				fm.CICO_amlaID = f.GetID()
			case fpb.UploadType_CICO_mtpp:
				fm.CICO_mtppForms = strings.Join(f.GetFileNames(), ",")
				fm.CICO_mtppID = f.GetID()
			case fpb.UploadType_CICO_is:
				fm.CICO_isForms = strings.Join(f.GetFileNames(), ",")
				fm.CICO_isID = f.GetID()
			case fpb.UploadType_CICO_edd:
				fm.CICO_eddForms = strings.Join(f.GetFileNames(), ",")
				fm.CICO_eddID = f.GetID()
			}
		}
	}
	return fm
}

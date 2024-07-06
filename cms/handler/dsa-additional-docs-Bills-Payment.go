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

type DsaUploadBillsPaymentAddiForm struct {
	CSRFField        template.HTML
	CSRFFieldValue   string
	Errors           map[string]error
	PresetPermission map[string]map[string]bool
	ServiceRequest   bool
	BP_ndaID         string
	BP_ndaForms      string
	BP_sisID         string
	BP_sisForms      string
	BP_psfID         string
	BP_psfForms      string
	BP_pspID         string
	BP_pspForms      string
	BP_secID         string
	BP_secForms      string
	BP_gisID         string
	BP_gisForms      string
	BP_lpafsID       string
	BP_lpafsForms    string
	BP_birID         string
	BP_birForms      string
	BP_bspID         string
	BP_bspForms      string
	BP_amlID         string
	BP_amlForms      string
	BP_sccasID       string
	BP_sccasForms    string
	BP_vgID          string
	BP_vgForms       string
	BP_cpID          string
	BP_cpForms       string
	BP_moaID         string
	BP_moaForms      string
	BP_amlaID        string
	BP_amlaForms     string
	BP_mttpID        string
	BP_mttpForms     string
	BP_sciID         string
	BP_sciForms      string
	BP_eddID         string
	BP_eddForms      string
	SaveContinue     string
	UserInfo         *User
	SaveDraft        string
	CompanyName      string
}

var validedBillsPaymentAddiFileType = map[string]bool{
	"image/jpeg":         true,
	"image/png":          true,
	"image/gif":          true,
	"application/pdf":    true,
	"application/msword": true, // MS-word files (extension .doc)
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // Document extension .docx
}

func (s *Server) getdsaAdiBillsPaymentSelPath(w http.ResponseWriter, r *http.Request) {
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

	template := s.templates.Lookup("dsa-additional-bills-payment.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	etd := s.getEnforceTemplateData(ctx)
	res, _ := s.pf.ListFiles(ctx, &fpb.ListFilesRequest{
		OrgID: oid,
		Types: []fpb.UploadType{fpb.UploadType_BP_nda, fpb.UploadType_BP_sis, fpb.UploadType_BP_psf, fpb.UploadType_BP_psp, fpb.UploadType_BP_sec, fpb.UploadType_BP_gis, fpb.UploadType_BP_lpafs, fpb.UploadType_BP_bir, fpb.UploadType_BP_bsp, fpb.UploadType_BP_aml, fpb.UploadType_BP_sccas, fpb.UploadType_BP_vg, fpb.UploadType_BP_cp, fpb.UploadType_BP_moa, fpb.UploadType_BP_amla, fpb.UploadType_BP_mttp, fpb.UploadType_BP_sci, fpb.UploadType_BP_edd},
	})
	var availablePartners []string
	isPartnerDraft := false
	ress, _ := s.pf.ListServiceRequest(ctx, &sVcpb.ListServiceRequestRequest{
		OrgIDs: []string{oid},
		Types:  []sVcpb.ServiceType{sVcpb.ServiceType_BILLSPAYMENT},
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
		http.Redirect(w, r, dsaPrtBillsPaymentSelPath, http.StatusSeeOther)
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	data := AddiUploadBillsPaymentProtoToForm(res)
	data.CSRFField = csrf.TemplateField(r)
	data.CSRFFieldValue = csrf.Token(r)
	data.PresetPermission = etd.PresetPermission
	data.ServiceRequest = etd.ServiceRequests
	data.UserInfo = &usrInfo.UserInfo
	data.CompanyName = usrInfo.CompanyName
	data.UserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, data); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) postDsaReqDocsBillsPayment(w http.ResponseWriter, r *http.Request) {
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
	var form DsaUploadBillsPaymentAddiForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	formErr := validation.Errors{}
	BP_nda, BP_sis, BP_psf, BP_psp, BP_sec, BP_gis, BP_lpafs, BP_bir, BP_bsp, BP_aml, BP_sccas, BP_vg, BP_cp, BP_moa, BP_amla, BP_mttp, BP_sci, BP_edd := "BP_nda", "BP_sis", "BP_psf", "BP_psp", "BP_sec", "BP_gis", "BP_lpafs", "BP_bir", "BP_bsp", "BP_aml", "BP_sccas", "BP_vg", "BP_cp", "BP_moa", "BP_amla", "BP_mttp", "BP_sci", "BP_edd"
	mu := map[string][]string{BP_nda: {}, BP_sis: {}, BP_psf: {}, BP_psp: {}, BP_sec: {}, BP_gis: {}, BP_lpafs: {}, BP_bir: {}, BP_bsp: {}, BP_aml: {}, BP_sccas: {}, BP_vg: {}, BP_cp: {}, BP_moa: {}, BP_amla: {}, BP_mttp: {}, BP_sci: {}, BP_edd: {}}
	for k := range mu {
		f, _, err := r.FormFile(k)
		if err != nil {
			continue
		}
		err = validateFileType(f, validedBillsPaymentAddiFileType)
		if err != nil {
			formErr[k] = err
		}
	}
	var availablePartners []string
	isPartnerDraft := false
	ress, _ := s.pf.ListServiceRequest(ctx, &sVcpb.ListServiceRequestRequest{
		OrgIDs: []string{oid},
		Types:  []sVcpb.ServiceType{sVcpb.ServiceType_BILLSPAYMENT},
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
		http.Redirect(w, r, dsaPrtBillsPaymentSelPath, http.StatusSeeOther)
	}
	requiredDocErr := errors.New("This is a required document, please upload a file")
	res, _ := s.pf.ListFiles(ctx, &fpb.ListFilesRequest{
		OrgID: oid,
		Types: []fpb.UploadType{fpb.UploadType_BP_nda, fpb.UploadType_BP_sis, fpb.UploadType_BP_psf, fpb.UploadType_BP_psp, fpb.UploadType_BP_sec, fpb.UploadType_BP_gis, fpb.UploadType_BP_lpafs, fpb.UploadType_BP_bir, fpb.UploadType_BP_bsp, fpb.UploadType_BP_aml, fpb.UploadType_BP_sccas, fpb.UploadType_BP_vg, fpb.UploadType_BP_cp, fpb.UploadType_BP_moa, fpb.UploadType_BP_amla, fpb.UploadType_BP_mttp, fpb.UploadType_BP_sci, fpb.UploadType_BP_edd},
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
			if gt == "BP_bsp" || gt == "BP_aml" || gt == "BP_edd" {
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
			if fileField == "BP_bsp" || fileField == "BP_aml" || fileField == "BP_edd" {
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
			err = validateFileType(f, validedBillsPaymentAddiFileType)
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
				FileNames: mu[BP_nda],
				Type:      fpb.UploadType_BP_nda,
				FileName:  fns[BP_nda],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[BP_sis],
				Type:      fpb.UploadType_BP_sis,
				FileName:  fns[BP_sis],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[BP_psf],
				Type:      fpb.UploadType_BP_psf,
				FileName:  fns[BP_psf],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[BP_psp],
				Type:      fpb.UploadType_BP_psp,
				FileName:  fns[BP_psp],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[BP_sec],
				Type:      fpb.UploadType_BP_sec,
				FileName:  fns[BP_sec],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[BP_gis],
				Type:      fpb.UploadType_BP_gis,
				FileName:  fns[BP_gis],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[BP_lpafs],
				Type:      fpb.UploadType_BP_lpafs,
				FileName:  fns[BP_lpafs],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[BP_bir],
				Type:      fpb.UploadType_BP_bir,
				FileName:  fns[BP_bir],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[BP_bsp],
				Type:      fpb.UploadType_BP_bsp,
				FileName:  fns[BP_bsp],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[BP_aml],
				Type:      fpb.UploadType_BP_aml,
				FileName:  fns[BP_aml],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[BP_sccas],
				Type:      fpb.UploadType_BP_sccas,
				FileName:  fns[BP_sccas],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[BP_vg],
				Type:      fpb.UploadType_BP_vg,
				FileName:  fns[BP_vg],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[BP_cp],
				Type:      fpb.UploadType_BP_cp,
				FileName:  fns[BP_cp],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[BP_moa],
				Type:      fpb.UploadType_BP_moa,
				FileName:  fns[BP_moa],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[BP_amla],
				Type:      fpb.UploadType_BP_amla,
				FileName:  fns[BP_amla],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[BP_mttp],
				Type:      fpb.UploadType_BP_mttp,
				FileName:  fns[BP_mttp],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[BP_sci],
				Type:      fpb.UploadType_BP_sci,
				FileName:  fns[BP_sci],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[BP_edd],
				Type:      fpb.UploadType_BP_edd,
				FileName:  fns[BP_edd],
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
			SvcName:  sVcpb.ServiceType_BILLSPAYMENT.String(),
			FileType: gt,
			FileID:   vv.GetID(),
			CreateBy: uid,
		}); err != nil {
			if err == storage.Conflict {
				if _, err := s.pf.UpdateUploadSvcRequest(ctx, &sVcpb.UpdateUploadSvcRequestRequest{
					OrgID:    oid,
					Partner:  "ALL",
					SvcName:  sVcpb.ServiceType_BILLSPAYMENT.String(),
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
			Types: []fpb.UploadType{fpb.UploadType_BP_nda, fpb.UploadType_BP_sis, fpb.UploadType_BP_psf, fpb.UploadType_BP_psp, fpb.UploadType_BP_sec, fpb.UploadType_BP_gis, fpb.UploadType_BP_lpafs, fpb.UploadType_BP_bir, fpb.UploadType_BP_bsp, fpb.UploadType_BP_aml, fpb.UploadType_BP_sccas, fpb.UploadType_BP_vg, fpb.UploadType_BP_cp, fpb.UploadType_BP_moa, fpb.UploadType_BP_amla, fpb.UploadType_BP_mttp, fpb.UploadType_BP_sci, fpb.UploadType_BP_edd},
		})
		form := AddiUploadBillsPaymentProtoToForm(res)
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
		template := s.templates.Lookup("dsa-additional-bills-payment.html")
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
			SvcName:  sVcpb.ServiceType_BILLSPAYMENT.String(),
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
			Type:  sVcpb.ServiceType_BILLSPAYMENT,
		}); err != nil {
			if err != storage.NotFound {
				logging.WithError(err, log).Error("Apply Service Request Failed")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, dsaAdiBillsPaymentSelPath, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, dsaServicesPath, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, dsaAdiBillsPaymentSelPath, http.StatusSeeOther)
	return
}

func AddiUploadBillsPaymentProtoToForm(fs *fpb.ListFilesResponse) DsaUploadBillsPaymentAddiForm {
	fm := DsaUploadBillsPaymentAddiForm{}
	if fs != nil {
		for _, f := range fs.GetFileUploads() {
			switch f.Type {
			case fpb.UploadType_BP_nda:
				fm.BP_ndaForms = strings.Join(f.GetFileNames(), ",")
				fm.BP_ndaID = f.GetID()
			case fpb.UploadType_BP_sis:
				fm.BP_sisForms = strings.Join(f.GetFileNames(), ",")
				fm.BP_sisID = f.GetID()
			case fpb.UploadType_BP_psf:
				fm.BP_psfForms = strings.Join(f.GetFileNames(), ",")
				fm.BP_psfID = f.GetID()
			case fpb.UploadType_BP_psp:
				fm.BP_pspForms = strings.Join(f.GetFileNames(), ",")
				fm.BP_pspID = f.GetID()
			case fpb.UploadType_BP_sec:
				fm.BP_secForms = strings.Join(f.GetFileNames(), ",")
				fm.BP_secID = f.GetID()
			case fpb.UploadType_BP_gis:
				fm.BP_gisForms = strings.Join(f.GetFileNames(), ",")
				fm.BP_gisID = f.GetID()
			case fpb.UploadType_BP_lpafs:
				fm.BP_lpafsForms = strings.Join(f.GetFileNames(), ",")
				fm.BP_lpafsID = f.GetID()
			case fpb.UploadType_BP_bir:
				fm.BP_birForms = strings.Join(f.GetFileNames(), ",")
				fm.BP_birID = f.GetID()
			case fpb.UploadType_BP_bsp:
				fm.BP_bspForms = strings.Join(f.GetFileNames(), ",")
				fm.BP_bspID = f.GetID()
			case fpb.UploadType_BP_aml:
				fm.BP_amlForms = strings.Join(f.GetFileNames(), ",")
				fm.BP_amlID = f.GetID()
			case fpb.UploadType_BP_sccas:
				fm.BP_sccasForms = strings.Join(f.GetFileNames(), ",")
				fm.BP_sccasID = f.GetID()
			case fpb.UploadType_BP_vg:
				fm.BP_vgForms = strings.Join(f.GetFileNames(), ",")
				fm.BP_vgID = f.GetID()
			case fpb.UploadType_BP_cp:
				fm.BP_cpForms = strings.Join(f.GetFileNames(), ",")
				fm.BP_cpID = f.GetID()
			case fpb.UploadType_BP_moa:
				fm.BP_moaForms = strings.Join(f.GetFileNames(), ",")
				fm.BP_moaID = f.GetID()
			case fpb.UploadType_BP_amla:
				fm.BP_amlaForms = strings.Join(f.GetFileNames(), ",")
				fm.BP_amlaID = f.GetID()
			case fpb.UploadType_BP_mttp:
				fm.BP_mttpForms = strings.Join(f.GetFileNames(), ",")
				fm.BP_mttpID = f.GetID()
			case fpb.UploadType_BP_sci:
				fm.BP_sciForms = strings.Join(f.GetFileNames(), ",")
				fm.BP_sciID = f.GetID()
			case fpb.UploadType_BP_edd:
				fm.BP_eddForms = strings.Join(f.GetFileNames(), ",")
				fm.BP_eddID = f.GetID()
			}
		}
	}
	return fm
}

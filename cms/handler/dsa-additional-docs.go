package handler

import (
	"errors"
	"html/template"
	"net/http"
	"strings"

	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/storage"
	cmsmw "brank.as/petnet/cms/mw"
	ppf "brank.as/petnet/gunk/dsa/v1/user"
	fpb "brank.as/petnet/gunk/dsa/v2/file"
	sVcpb "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
)

type DsaUploadAddiForm struct {
	CSRFField        template.HTML
	CSRFFieldValue   string
	Errors           map[string]error
	PresetPermission map[string]map[string]bool
	ServiceRequest   bool
	WU_srID          string
	WU_srForms       string
	WU_lgisID        string
	WU_lgisForms     string
	WU_dtisspaID     string
	WU_dtisspaForms  string
	WU_birfID        string
	WU_birfForms     string
	WU_bsprID        string
	WU_bsprForms     string
	WU_iqaID         string
	WU_iqaForms      string
	WU_bqaID         string
	WU_bqaForms      string
	AYA_dfedrID      string
	AYA_dfedrForms   string
	AYA_ddfID        string
	AYA_ddfForms     string
	AYA_ialawsID     string
	AYA_ialawsForms  string
	AYA_aialawsID    string
	AYA_aialawsForms string
	AYA_mtlID        string
	AYA_mtlForms     string
	AYA_brdcsaID     string
	AYA_brdcsaForms  string
	AYA_gisID        string
	AYA_gisForms     string
	AYA_ccpwiID      string
	AYA_ccpwiForms   string
	AYA_fasID        string
	AYA_fasForms     string
	AYA_amID         string
	AYA_amForms      string
	AYA_lafID        string
	AYA_lafForms     string
	AYA_birrID       string
	AYA_birrForms    string
	AYA_odID         string
	AYA_odForms      string
	HasAya           bool
	HasWu            bool
	HasAddi          bool
	UserInfo         *User
	SaveDraft        string
	CompanyName      string
}

var validedAddiFileType = map[string]bool{
	"image/jpeg":         true,
	"image/png":          true,
	"image/gif":          true,
	"application/pdf":    true,
	"application/msword": true, // MS-word files (extension .doc)
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // Document extension .docx
}

func (s *Server) getDsaAddiDocs(w http.ResponseWriter, r *http.Request) {
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
	var availablePartners []string
	isPartnerDraft := false
	isReqDraft := false
	ress, _ := s.pf.ListServiceRequest(ctx, &sVcpb.ListServiceRequestRequest{
		OrgIDs: []string{oid},
		Types:  []sVcpb.ServiceType{sVcpb.ServiceType_REMITTANCE},
	})
	if ress != nil {
		for _, vv := range ress.GetServiceRequst() {
			availablePartners = append(availablePartners, vv.Partner)
			if vv.Status == sVcpb.ServiceRequestStatus_PARTNERDRAFT {
				isPartnerDraft = true
			}
			if vv.Status == sVcpb.ServiceRequestStatus_REQDOCDRAFT {
				isReqDraft = true
			}
		}
	}
	if isPartnerDraft || len(availablePartners) == 0 {
		http.Redirect(w, r, dsaPrtSelPath, http.StatusSeeOther)
	}
	if isReqDraft {
		http.Redirect(w, r, dsaReqDocsPath, http.StatusSeeOther)
	}
	ha, _ := cmsmw.InArray(static.AYACode, availablePartners)
	wa, _ := cmsmw.InArray(static.WUCode, availablePartners)
	if !(ha || wa) {
		http.Redirect(w, r, dsaServicesPath, http.StatusSeeOther)
		return
	}
	template := s.templates.Lookup("dsa-additional-docs.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	etd := s.getEnforceTemplateData(ctx)
	res, _ := s.pf.ListFiles(ctx, &fpb.ListFilesRequest{
		OrgID: oid,
		Types: []fpb.UploadType{fpb.UploadType_WU_sr, fpb.UploadType_WU_lgis, fpb.UploadType_WU_dtisspa, fpb.UploadType_WU_birf, fpb.UploadType_WU_bspr, fpb.UploadType_WU_iqa, fpb.UploadType_WU_bqa, fpb.UploadType_AYA_dfedr, fpb.UploadType_AYA_ddf, fpb.UploadType_AYA_ialaws, fpb.UploadType_AYA_aialaws, fpb.UploadType_AYA_mtl, fpb.UploadType_AYA_brdcsa, fpb.UploadType_AYA_gis, fpb.UploadType_AYA_ccpwi, fpb.UploadType_AYA_fas, fpb.UploadType_AYA_am, fpb.UploadType_AYA_laf, fpb.UploadType_AYA_birr, fpb.UploadType_AYA_od},
	})
	data := AddiUploadProtoToForm(res)
	data.HasAya = ha
	data.HasWu = wa
	data.HasAddi = (ha || wa)
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
}

func (s *Server) postDsaAddiDocs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	uid := mw.GetUserID(ctx)
	requiredDocErr := errors.New("This is a required document, please upload a file")
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
	var availablePartners []string
	isPartnerDraft := false
	isReqDraft := false
	ress, _ := s.pf.ListServiceRequest(ctx, &sVcpb.ListServiceRequestRequest{
		OrgIDs: []string{oid},
		Types:  []sVcpb.ServiceType{sVcpb.ServiceType_REMITTANCE},
	})
	if ress != nil {
		for _, vv := range ress.GetServiceRequst() {
			availablePartners = append(availablePartners, vv.Partner)
			if vv.Status == sVcpb.ServiceRequestStatus_PARTNERDRAFT {
				isPartnerDraft = true
			}
			if vv.Status == sVcpb.ServiceRequestStatus_REQDOCDRAFT {
				isReqDraft = true
			}
		}
	}
	if isPartnerDraft || len(availablePartners) == 0 {
		http.Redirect(w, r, dsaPrtSelPath, http.StatusSeeOther)
		return
	}
	if isReqDraft {
		http.Redirect(w, r, dsaReqDocsPath, http.StatusSeeOther)
		return
	}
	ha, _ := cmsmw.InArray(static.AYACode, availablePartners)
	wa, _ := cmsmw.InArray(static.WUCode, availablePartners)
	var form DsaUploadAddiForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	formErr := validation.Errors{}
	WU_sr, WU_lgis, WU_dtisspa, WU_birf, WU_bspr, WU_iqa, WU_bqa, AYA_dfedr, AYA_ddf, AYA_ialaws, AYA_aialaws, AYA_mtl, AYA_brdcsa, AYA_gis, AYA_ccpwi, AYA_fas, AYA_am, AYA_laf, AYA_birr, AYA_od := "WU_sr", "WU_lgis", "WU_dtisspa", "WU_birf", "WU_bspr", "WU_iqa", "WU_bqa", "AYA_dfedr", "AYA_ddf", "AYA_ialaws", "AYA_aialaws", "AYA_mtl", "AYA_brdcsa", "AYA_gis", "AYA_ccpwi", "AYA_fas", "AYA_am", "AYA_laf", "AYA_birr", "AYA_od"
	mu := map[string][]string{WU_sr: {}, WU_lgis: {}, WU_dtisspa: {}, WU_birf: {}, WU_bspr: {}, WU_iqa: {}, WU_bqa: {}, AYA_dfedr: {}, AYA_ddf: {}, AYA_ialaws: {}, AYA_aialaws: {}, AYA_mtl: {}, AYA_brdcsa: {}, AYA_gis: {}, AYA_ccpwi: {}, AYA_fas: {}, AYA_am: {}, AYA_laf: {}, AYA_birr: {}, AYA_od: {}}
	ut := []fpb.UploadType{fpb.UploadType_WU_sr, fpb.UploadType_WU_lgis, fpb.UploadType_WU_dtisspa, fpb.UploadType_WU_birf, fpb.UploadType_WU_bspr, fpb.UploadType_WU_iqa, fpb.UploadType_WU_bqa, fpb.UploadType_AYA_dfedr, fpb.UploadType_AYA_ddf, fpb.UploadType_AYA_ialaws, fpb.UploadType_AYA_aialaws, fpb.UploadType_AYA_mtl, fpb.UploadType_AYA_brdcsa, fpb.UploadType_AYA_gis, fpb.UploadType_AYA_ccpwi, fpb.UploadType_AYA_fas, fpb.UploadType_AYA_am, fpb.UploadType_AYA_laf, fpb.UploadType_AYA_birr, fpb.UploadType_AYA_od}
	if ha && !wa {
		mu = map[string][]string{AYA_dfedr: {}, AYA_ddf: {}, AYA_ialaws: {}, AYA_aialaws: {}, AYA_mtl: {}, AYA_brdcsa: {}, AYA_gis: {}, AYA_ccpwi: {}, AYA_fas: {}, AYA_am: {}, AYA_laf: {}, AYA_birr: {}, AYA_od: {}}
		ut = []fpb.UploadType{fpb.UploadType_AYA_dfedr, fpb.UploadType_AYA_ddf, fpb.UploadType_AYA_ialaws, fpb.UploadType_AYA_aialaws, fpb.UploadType_AYA_mtl, fpb.UploadType_AYA_brdcsa, fpb.UploadType_AYA_gis, fpb.UploadType_AYA_ccpwi, fpb.UploadType_AYA_fas, fpb.UploadType_AYA_am, fpb.UploadType_AYA_laf, fpb.UploadType_AYA_birr, fpb.UploadType_AYA_od}
	} else if !ha && wa {
		mu = map[string][]string{WU_sr: {}, WU_lgis: {}, WU_dtisspa: {}, WU_birf: {}, WU_bspr: {}, WU_iqa: {}, WU_bqa: {}}
		ut = []fpb.UploadType{fpb.UploadType_WU_sr, fpb.UploadType_WU_lgis, fpb.UploadType_WU_dtisspa, fpb.UploadType_WU_birf, fpb.UploadType_WU_bspr, fpb.UploadType_WU_iqa, fpb.UploadType_WU_bqa}
	}
	for k := range mu {
		f, _, err := r.FormFile(k)
		if err != nil {
			continue
		}
		err = validateFileType(f, validedAddiFileType)
		if err != nil {
			formErr[k] = err
		}
	}
	if !(ha || wa) {
		http.Redirect(w, r, dsaServicesPath, http.StatusSeeOther)
		return
	}
	res, _ := s.pf.ListFiles(ctx, &fpb.ListFilesRequest{
		OrgID: oid,
		Types: ut,
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
			if gt == "WU_lgis" || gt == "WU_dtisspa" {
				continue
			}
			ful[gt] = []string{}
		}
	}
	if len(ful) == 0 {
		ful = mu
	}
	for fileField := range ful {
		if fileField == "WU_lgis" || fileField == "WU_dtisspa" {
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
		err = validateFileType(f, validedAddiFileType)
		if err != nil {
			formErr[fileField] = err
		}
	}
	fns := map[string]map[string]string{}
	for k := range mu {
		fns[k] = map[string]string{}
		if len(fulV[k]) > 0 {
			mu[k] = fulV[k]
			continue
		}
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
		mu[k] = us
		fns[k] = fn
	}
	// todo(robin): if we decide to add another bucket make sure the bucket name is saved.
	ufr, err := s.pf.UpsertFiles(ctx, &fpb.UpsertFilesRequest{
		FileUploads: []*fpb.FileUpload{
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[WU_sr],
				Type:      fpb.UploadType_WU_sr,
				FileName:  fns[WU_sr],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[WU_lgis],
				Type:      fpb.UploadType_WU_lgis,
				FileName:  fns[WU_lgis],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[WU_dtisspa],
				Type:      fpb.UploadType_WU_dtisspa,
				FileName:  fns[WU_dtisspa],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[WU_birf],
				Type:      fpb.UploadType_WU_birf,
				FileName:  fns[WU_birf],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[WU_bspr],
				Type:      fpb.UploadType_WU_bspr,
				FileName:  fns[WU_bspr],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[WU_iqa],
				Type:      fpb.UploadType_WU_iqa,
				FileName:  fns[WU_iqa],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[WU_bqa],
				Type:      fpb.UploadType_WU_bqa,
				FileName:  fns[WU_bqa],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[AYA_dfedr],
				Type:      fpb.UploadType_AYA_dfedr,
				FileName:  fns[AYA_dfedr],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[AYA_ddf],
				Type:      fpb.UploadType_AYA_ddf,
				FileName:  fns[AYA_ddf],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[AYA_ialaws],
				Type:      fpb.UploadType_AYA_ialaws,
				FileName:  fns[AYA_ialaws],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[AYA_aialaws],
				Type:      fpb.UploadType_AYA_aialaws,
				FileName:  fns[AYA_aialaws],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[AYA_mtl],
				Type:      fpb.UploadType_AYA_mtl,
				FileName:  fns[AYA_mtl],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[AYA_brdcsa],
				Type:      fpb.UploadType_AYA_brdcsa,
				FileName:  fns[AYA_brdcsa],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[AYA_gis],
				Type:      fpb.UploadType_AYA_gis,
				FileName:  fns[AYA_gis],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[AYA_ccpwi],
				Type:      fpb.UploadType_AYA_ccpwi,
				FileName:  fns[AYA_ccpwi],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[AYA_fas],
				Type:      fpb.UploadType_AYA_fas,
				FileName:  fns[AYA_fas],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[AYA_am],
				Type:      fpb.UploadType_AYA_am,
				FileName:  fns[AYA_am],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[AYA_laf],
				Type:      fpb.UploadType_AYA_laf,
				FileName:  fns[AYA_laf],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[AYA_birr],
				Type:      fpb.UploadType_AYA_birr,
				FileName:  fns[AYA_birr],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[AYA_od],
				Type:      fpb.UploadType_AYA_od,
				FileName:  fns[AYA_od],
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
		ptnt := "ALL"
		if strings.HasPrefix(gt, "WU_") {
			ptnt = "WU"
		} else if strings.HasPrefix(gt, "AYA_") {
			ptnt = "AYA"
		}
		if _, err := s.pf.AddUploadSvcRequest(ctx, &sVcpb.AddUploadSvcRequestRequest{
			OrgID:    oid,
			Partner:  ptnt,
			SvcName:  sVcpb.ServiceType_REMITTANCE.String(),
			FileType: gt,
			FileID:   vv.GetID(),
			CreateBy: uid,
		}); err != nil {
			if err == storage.Conflict {
				if _, err := s.pf.UpdateUploadSvcRequest(ctx, &sVcpb.UpdateUploadSvcRequestRequest{
					OrgID:    oid,
					Partner:  ptnt,
					SvcName:  sVcpb.ServiceType_REMITTANCE.String(),
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
	if len(formErr) > 0 && form.SaveDraft != "SaveDraft" {
		res, _ := s.pf.ListFiles(ctx, &fpb.ListFilesRequest{
			OrgID: oid,
			Types: ut,
		})
		form := AddiUploadProtoToForm(res)
		form.HasAya = ha
		form.HasWu = wa
		form.HasAddi = (ha || wa)
		form.Errors = formErr
		form.CSRFField = csrf.TemplateField(r)
		form.CSRFFieldValue = csrf.Token(r)
		etd := s.getEnforceTemplateData(ctx)
		form.PresetPermission = etd.PresetPermission
		form.ServiceRequest = etd.ServiceRequests
		usrInfo := s.GetUserInfoFromCookie(w, r, false)

		form.UserInfo = &usrInfo.UserInfo
		uid := mw.GetUserID(ctx)
		gp, err := s.pf.GetUserProfile(ctx, &ppf.GetUserProfileRequest{
			UserID: uid,
		})
		if err != nil {
			log.Error("failed to get profile")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		form.UserInfo.ProfileImage = gp.GetProfile().ProfilePicture
		template := s.templates.Lookup("dsa-additional-docs.html")
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
		return
	}

	if form.SaveDraft != "SaveDraft" {
		if _, err := s.pf.ApplyServiceRequest(ctx, &sVcpb.ApplyServiceRequestRequest{
			OrgID: oid,
			Type:  sVcpb.ServiceType_REMITTANCE,
		}); err != nil {
			logging.WithError(err, log).Error("Apply Service Request")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}

	http.Redirect(w, r, dsaServicesPath, http.StatusSeeOther)
}

func AddiUploadProtoToForm(fs *fpb.ListFilesResponse) DsaUploadAddiForm {
	fm := DsaUploadAddiForm{}
	if fs != nil {
		for _, f := range fs.GetFileUploads() {
			switch f.Type {
			case fpb.UploadType_WU_sr:
				fm.WU_srForms = strings.Join(f.GetFileNames(), ",")
				fm.WU_srID = f.GetID()
			case fpb.UploadType_WU_lgis:
				fm.WU_lgisForms = strings.Join(f.GetFileNames(), ",")
				fm.WU_lgisID = f.GetID()
			case fpb.UploadType_WU_dtisspa:
				fm.WU_dtisspaForms = strings.Join(f.GetFileNames(), ",")
				fm.WU_dtisspaID = f.GetID()
			case fpb.UploadType_WU_birf:
				fm.WU_birfForms = strings.Join(f.GetFileNames(), ",")
				fm.WU_birfID = f.GetID()
			case fpb.UploadType_WU_bspr:
				fm.WU_bsprForms = strings.Join(f.GetFileNames(), ",")
				fm.WU_bsprID = f.GetID()
			case fpb.UploadType_WU_iqa:
				fm.WU_iqaForms = strings.Join(f.GetFileNames(), ",")
				fm.WU_iqaID = f.GetID()
			case fpb.UploadType_WU_bqa:
				fm.WU_bqaForms = strings.Join(f.GetFileNames(), ",")
				fm.WU_bqaID = f.GetID()
			case fpb.UploadType_AYA_dfedr:
				fm.AYA_dfedrForms = strings.Join(f.GetFileNames(), ",")
				fm.AYA_dfedrID = f.GetID()
			case fpb.UploadType_AYA_ddf:
				fm.AYA_ddfForms = strings.Join(f.GetFileNames(), ",")
				fm.AYA_ddfID = f.GetID()
			case fpb.UploadType_AYA_ialaws:
				fm.AYA_ialawsForms = strings.Join(f.GetFileNames(), ",")
				fm.AYA_ialawsID = f.GetID()
			case fpb.UploadType_AYA_aialaws:
				fm.AYA_aialawsForms = strings.Join(f.GetFileNames(), ",")
				fm.AYA_aialawsID = f.GetID()
			case fpb.UploadType_AYA_mtl:
				fm.AYA_mtlForms = strings.Join(f.GetFileNames(), ",")
				fm.AYA_mtlID = f.GetID()
			case fpb.UploadType_AYA_brdcsa:
				fm.AYA_brdcsaForms = strings.Join(f.GetFileNames(), ",")
				fm.AYA_brdcsaID = f.GetID()
			case fpb.UploadType_AYA_gis:
				fm.AYA_gisForms = strings.Join(f.GetFileNames(), ",")
				fm.AYA_gisID = f.GetID()
			case fpb.UploadType_AYA_ccpwi:
				fm.AYA_ccpwiForms = strings.Join(f.GetFileNames(), ",")
				fm.AYA_ccpwiID = f.GetID()
			case fpb.UploadType_AYA_fas:
				fm.AYA_fasForms = strings.Join(f.GetFileNames(), ",")
				fm.AYA_fasID = f.GetID()
			case fpb.UploadType_AYA_am:
				fm.AYA_amForms = strings.Join(f.GetFileNames(), ",")
				fm.AYA_amID = f.GetID()
			case fpb.UploadType_AYA_laf:
				fm.AYA_lafForms = strings.Join(f.GetFileNames(), ",")
				fm.AYA_lafID = f.GetID()
			case fpb.UploadType_AYA_birr:
				fm.AYA_birrForms = strings.Join(f.GetFileNames(), ",")
				fm.AYA_birrID = f.GetID()
			case fpb.UploadType_AYA_od:
				fm.AYA_odForms = strings.Join(f.GetFileNames(), ",")
				fm.AYA_odID = f.GetID()
			}
		}
	}
	return fm
}

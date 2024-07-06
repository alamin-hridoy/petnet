package handler

import (
	"errors"
	"html/template"
	"net/http"
	"strings"

	"brank.as/petnet/api/core/static"
	cmsmw "brank.as/petnet/cms/mw"
	ppf "brank.as/petnet/gunk/dsa/v1/user"
	fpb "brank.as/petnet/gunk/dsa/v2/file"
	sVcpb "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
)

type DsaUploadReqForm struct {
	CSRFField        template.HTML
	Errors           map[string]error
	PresetPermission map[string]map[string]bool
	ServiceRequest   bool
	SecForms         string
	SecID            string
	GisForms         string
	GisID            string
	BrsForms         string
	BrsID            string
	CpForms          string
	CpID             string
	AmlForms         string
	AmlID            string
	NndaForms        string
	NndaID           string
	PsfForms         string
	PsfID            string
	PspForms         string
	PspID            string
	KddqForms        string
	KddqID           string
	SisForms         string
	SisID            string
	ScbrForms        string
	ScbrID           string
	ViaForms         string
	ViaID            string
	BsprForms        string
	BsprID           string
	Sec              string
	Gis              string
	Brs              string
	Scbr             string
	Via              string
	Cp               string
	Aml              string
	Nnda             string
	Psf              string
	Psp              string
	Kddq             string
	Sis              string
	SaveDraft        string
	SaveContinue     string
	CSRFFieldValue   string
	HasAya           bool
	HasWu            bool
	HasAddi          bool
	UserInfo         *User
	CompanyName      string
}

var validedReqFileType = map[string]bool{
	"image/jpeg":         true,
	"image/png":          true,
	"image/gif":          true,
	"application/pdf":    true,
	"application/msword": true, // MS-word files (extension .doc)
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // Document extension .docx
}

func (s *Server) getDsaReqDocs(w http.ResponseWriter, r *http.Request) {
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
	template := s.templates.Lookup("dsa-req-docs.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	etd := s.getEnforceTemplateData(ctx)
	res, _ := s.pf.ListFiles(ctx, &fpb.ListFilesRequest{
		OrgID: oid,
		Types: []fpb.UploadType{fpb.UploadType_Sec, fpb.UploadType_Gis, fpb.UploadType_Brs, fpb.UploadType_Scbr, fpb.UploadType_Via, fpb.UploadType_Bspr, fpb.UploadType_Cp, fpb.UploadType_Aml, fpb.UploadType_Nnda, fpb.UploadType_Psf, fpb.UploadType_Psp, fpb.UploadType_Kddq, fpb.UploadType_Sis},
	})
	var availablePartners []string
	isPartnerDraft := false
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
		}
	}
	if isPartnerDraft || len(availablePartners) == 0 {
		http.Redirect(w, r, dsaPrtSelPath, http.StatusSeeOther)
	}
	ha, _ := cmsmw.InArray(static.AYACode, availablePartners)
	wa, _ := cmsmw.InArray(static.WUCode, availablePartners)
	data := ReqUploadProtoToForm(res)
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
	return
}

func (s *Server) postDsaReqDocs(w http.ResponseWriter, r *http.Request) {
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
	var form DsaUploadReqForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	formErr := validation.Errors{}
	sec, gis, brs, scbr, via, bspr, cp, aml, nnda, psf, psp, kddq, sis := "Sec", "Gis", "Brs", "Scbr", "Via", "Bspr", "Cp", "Aml", "Nnda", "Psf", "Psp", "Kddq", "Sis"
	mu := map[string][]string{sec: {}, gis: {}, brs: {}, scbr: {}, via: {}, bspr: {}, cp: {}, aml: {}, nnda: {}, psf: {}, psp: {}, kddq: {}, sis: {}}
	for k := range mu {
		f, _, err := r.FormFile(k)
		if err != nil {
			continue
		}
		err = validateFileType(f, validedReqFileType)
		if err != nil {
			formErr[k] = err
		}
	}
	var availablePartners []string
	var exceptAcceptedPartners []string
	isPartnerDraft := false
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
			if vv.Status != sVcpb.ServiceRequestStatus_ACCEPTED && vv.Status != sVcpb.ServiceRequestStatus_PENDING {
				exceptAcceptedPartners = append(exceptAcceptedPartners, vv.Partner)
			}
		}
	}
	if isPartnerDraft || len(availablePartners) == 0 {
		http.Redirect(w, r, dsaPrtSelPath, http.StatusSeeOther)
	}
	ha, _ := cmsmw.InArray(static.AYACode, availablePartners)
	wa, _ := cmsmw.InArray(static.WUCode, availablePartners)
	requiredDocErr := errors.New("This is a required document, please upload a file")
	res, _ := s.pf.ListFiles(ctx, &fpb.ListFilesRequest{
		OrgID: oid,
		Types: []fpb.UploadType{fpb.UploadType_Sec, fpb.UploadType_Gis, fpb.UploadType_Brs, fpb.UploadType_Scbr, fpb.UploadType_Via, fpb.UploadType_Bspr, fpb.UploadType_Cp, fpb.UploadType_Aml, fpb.UploadType_Nnda, fpb.UploadType_Psf, fpb.UploadType_Psp, fpb.UploadType_Kddq, fpb.UploadType_Sis},
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
			if gt == "Aml" || gt == "Bspr" {
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
			if fileField == "Aml" || fileField == "Bspr" {
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
			err = validateFileType(f, validedReqFileType)
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
				FileNames: mu[sec],
				Type:      fpb.UploadType_Sec,
				FileName:  fns[sec],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[gis],
				Type:      fpb.UploadType_Gis,
				FileName:  fns[gis],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[brs],
				Type:      fpb.UploadType_Brs,
				FileName:  fns[brs],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[scbr],
				Type:      fpb.UploadType_Scbr,
				FileName:  fns[scbr],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[via],
				Type:      fpb.UploadType_Via,
				FileName:  fns[via],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[bspr],
				Type:      fpb.UploadType_Bspr,
				FileName:  fns[bspr],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[cp],
				Type:      fpb.UploadType_Cp,
				FileName:  fns[cp],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[aml],
				Type:      fpb.UploadType_Aml,
				FileName:  fns[aml],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[nnda],
				Type:      fpb.UploadType_Nnda,
				FileName:  fns[nnda],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[psf],
				Type:      fpb.UploadType_Psf,
				FileName:  fns[psf],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[psp],
				Type:      fpb.UploadType_Psp,
				FileName:  fns[psp],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[kddq],
				Type:      fpb.UploadType_Kddq,
				FileName:  fns[kddq],
			},
			{
				UserID:    uid,
				OrgID:     oid,
				FileNames: mu[sis],
				Type:      fpb.UploadType_Sis,
				FileName:  fns[sis],
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
			SvcName:  sVcpb.ServiceType_REMITTANCE.String(),
			FileType: gt,
			FileID:   vv.GetID(),
			CreateBy: uid,
		}); err != nil {
			if err == storage.Conflict {
				if _, err := s.pf.UpdateUploadSvcRequest(ctx, &sVcpb.UpdateUploadSvcRequestRequest{
					OrgID:    oid,
					Partner:  "ALL",
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
	if len(formErr) > 0 {
		res, _ := s.pf.ListFiles(ctx, &fpb.ListFilesRequest{
			OrgID: oid,
			Types: []fpb.UploadType{fpb.UploadType_Sec, fpb.UploadType_Gis, fpb.UploadType_Brs, fpb.UploadType_Scbr, fpb.UploadType_Via, fpb.UploadType_Bspr, fpb.UploadType_Cp, fpb.UploadType_Aml, fpb.UploadType_Nnda, fpb.UploadType_Psf, fpb.UploadType_Psp, fpb.UploadType_Kddq, fpb.UploadType_Sis},
		})
		form := ReqUploadProtoToForm(res)
		form.Errors = formErr
		form.CSRFField = csrf.TemplateField(r)
		form.CSRFFieldValue = csrf.Token(r)
		form.HasAya = ha
		form.HasWu = wa
		form.HasAddi = (ha || wa)
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
			Partners: exceptAcceptedPartners,
			SvcName:  sVcpb.ServiceType_REMITTANCE.String(),
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
		if ha || wa {
			s.pf.SetStatusUploadSvcRequest(ctx, &sVcpb.SetStatusUploadSvcRequestRequest{
				OrgID:    oid,
				Partners: exceptAcceptedPartners,
				SvcName:  sVcpb.ServiceType_REMITTANCE.String(),
				Status:   sVcpb.ServiceRequestStatus_ADDIDOCDRAFT,
			})
			http.Redirect(w, r, dsaAddiDocsPath, http.StatusSeeOther)
			return
		}
		if _, err := s.pf.ApplyServiceRequest(ctx, &sVcpb.ApplyServiceRequestRequest{
			OrgID: oid,
			Type:  sVcpb.ServiceType_REMITTANCE,
		}); err != nil {
			if err != storage.NotFound {
				logging.WithError(err, log).Error("Apply Service Request Failed")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, dsaReqDocsPath, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, dsaServicesPath, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, dsaReqDocsPath, http.StatusSeeOther)
	return
}

func ReqUploadProtoToForm(fs *fpb.ListFilesResponse) DsaUploadReqForm {
	fm := DsaUploadReqForm{}
	if fs != nil {
		for _, f := range fs.GetFileUploads() {
			switch f.Type {
			case fpb.UploadType_Sec:
				fm.SecForms = strings.Join(f.GetFileNames(), ",")
				fm.SecID = f.GetID()
			case fpb.UploadType_Gis:
				fm.GisForms = strings.Join(f.GetFileNames(), ",")
				fm.GisID = f.GetID()
			case fpb.UploadType_Brs:
				fm.BrsForms = strings.Join(f.GetFileNames(), ",")
				fm.BrsID = f.GetID()
			case fpb.UploadType_Scbr:
				fm.ScbrForms = strings.Join(f.GetFileNames(), ",")
				fm.ScbrID = f.GetID()
			case fpb.UploadType_Via:
				fm.ViaForms = strings.Join(f.GetFileNames(), ",")
				fm.ViaID = f.GetID()
			case fpb.UploadType_Bspr:
				fm.BsprForms = strings.Join(f.GetFileNames(), ",")
				fm.BsprID = f.GetID()
			case fpb.UploadType_Cp:
				fm.CpForms = strings.Join(f.GetFileNames(), ",")
				fm.CpID = f.GetID()
			case fpb.UploadType_Aml:
				fm.AmlForms = strings.Join(f.GetFileNames(), ",")
				fm.AmlID = f.GetID()
			case fpb.UploadType_Nnda:
				fm.NndaForms = strings.Join(f.GetFileNames(), ",")
				fm.NndaID = f.GetID()
			case fpb.UploadType_Psf:
				fm.PsfForms = strings.Join(f.GetFileNames(), ",")
				fm.PsfID = f.GetID()
			case fpb.UploadType_Psp:
				fm.PspForms = strings.Join(f.GetFileNames(), ",")
				fm.PspID = f.GetID()
			case fpb.UploadType_Kddq:
				fm.KddqForms = strings.Join(f.GetFileNames(), ",")
				fm.KddqID = f.GetID()
			case fpb.UploadType_Sis:
				fm.SisForms = strings.Join(f.GetFileNames(), ",")
				fm.SisID = f.GetID()
			}
		}
	}
	return fm
}

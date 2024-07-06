package handler

import (
	"html/template"
	"net/http"
	"net/url"
	"time"

	"brank.as/petnet/api/core/static"
	cmsmw "brank.as/petnet/cms/mw"
	fpb "brank.as/petnet/gunk/dsa/v2/file"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	svr "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	rbupb "brank.as/rbac/gunk/v1/user"
	"github.com/gorilla/csrf"
	"github.com/kenshaw/goji"
)

type (
	ServicesDetails struct {
		CSRFField         template.HTML
		CSRFFieldValue    string
		ID                string
		OrgID             string
		Status            string
		FileStatus        map[string]string
		ServiceRequests   []*ServiceRequests
		CompanyName       string
		Files             ReqFiles
		BaseURL           string
		ActivePartnerName string
		UserInfo          *User
		PresetPermission  map[string]map[string]bool
		ServiceRequest    bool
		ReqFilesDetails   map[string]string
		HasAya            bool
		HasWu             bool
		HasAddi           bool
	}

	ServiceRequests struct {
		OrgID       string
		CompanyName string
		Partner     string
		PartnerName string
		Type        svr.ServiceType
		Status      svr.ServiceRequestStatus
		Enabled     bool
		Remarks     string
		Applied     time.Time
		Created     time.Time
		Updated     time.Time
		UpdatedBy   string
		ID          string
	}
)

func (s *Server) getServicesDetails(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	oid := goji.Param(r, "id")
	if oid == "" {
		log.Error("missing org id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	queryParams := r.URL.Query()
	partner, err := url.PathUnescape(queryParams.Get("partner"))
	if err != nil {
		partner = "ALL"
	}
	if partner == "" {
		partner = "ALL"
	}
	template := s.templates.Lookup("services-details.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	pf, err := s.pf.GetProfile(r.Context(), &ppb.GetProfileRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}
	lsr, err := s.pf.ListServiceRequest(r.Context(), &svr.ListServiceRequestRequest{
		OrgIDs:   []string{oid},
		Types:    []svr.ServiceType{svr.ServiceType_REMITTANCE},
		Statuses: []svr.ServiceRequestStatus{svr.ServiceRequestStatus_PENDING, svr.ServiceRequestStatus_ACCEPTED, svr.ServiceRequestStatus_REJECTED},
	})
	if err != nil {
		http.Redirect(w, r, serviceReqPath, http.StatusSeeOther)
		logging.WithError(err, log).Info("List Service Request")
	}
	PartnerListMaps := s.getPartnerListMaps(ctx, svr.ServiceType_REMITTANCE.String())
	var availablePartners []string
	if lsr != nil {
		for _, vv := range lsr.GetServiceRequst() {
			availablePartners = append(availablePartners, vv.Partner)
		}
	}
	ha, _ := cmsmw.InArray(static.AYACode, availablePartners)
	wa, _ := cmsmw.InArray(static.WUCode, availablePartners)
	fs := map[string]string{}
	ulsr, err := s.pf.ListUploadSvcRequest(ctx, &svr.ListUploadSvcRequestRequest{
		OrgID:    oid,
		SvcNames: []string{svr.ServiceType_REMITTANCE.String()},
	})
	if err != nil {
		http.Redirect(w, r, serviceReqPath, http.StatusSeeOther)
		logging.WithError(err, log).Info("List Upload Service Request")
	}
	if len(ulsr.GetResults()) > 0 {
		for _, v := range ulsr.GetResults() {
			fs[v.FileType] = v.Status
		}
	}
	lf, _ := s.pf.ListFiles(ctx, &fpb.ListFilesRequest{
		OrgID: oid,
		Types: []fpb.UploadType{fpb.UploadType_Mbf, fpb.UploadType_Sec, fpb.UploadType_Gis, fpb.UploadType_Afs, fpb.UploadType_Brs, fpb.UploadType_Bmp, fpb.UploadType_Scbr, fpb.UploadType_Via, fpb.UploadType_WU_cd, fpb.UploadType_WU_lbp, fpb.UploadType_WU_sr, fpb.UploadType_WU_lgis, fpb.UploadType_WU_dtisspa, fpb.UploadType_WU_birf, fpb.UploadType_WU_bspr, fpb.UploadType_WU_iqa, fpb.UploadType_WU_bqa, fpb.UploadType_AYA_dfedr, fpb.UploadType_AYA_ddf, fpb.UploadType_AYA_ialaws, fpb.UploadType_AYA_aialaws, fpb.UploadType_AYA_mtl, fpb.UploadType_AYA_cbpr, fpb.UploadType_AYA_brdcsa, fpb.UploadType_AYA_gis, fpb.UploadType_AYA_ccpwi, fpb.UploadType_AYA_fas, fpb.UploadType_AYA_am, fpb.UploadType_AYA_laf, fpb.UploadType_AYA_birr, fpb.UploadType_AYA_od, fpb.UploadType_Bspr, fpb.UploadType_Cp, fpb.UploadType_Aml, fpb.UploadType_Nnda, fpb.UploadType_Psf, fpb.UploadType_Psp, fpb.UploadType_Kddq, fpb.UploadType_Sis},
	})
	_, err = s.rbac.GetUser(ctx, &rbupb.GetUserRequest{ID: pf.Profile.UserID})
	if err != nil {
		logging.WithError(err, log).Info("getting user")
		return
	}
	cn := ""
	sr := []*ServiceRequests{}
	if len(lsr.GetServiceRequst()) > 0 {
		for _, v := range lsr.GetServiceRequst() {
			cn = v.GetCompanyName()
			sr = append(sr, &ServiceRequests{
				OrgID:       v.GetOrgID(),
				CompanyName: v.GetCompanyName(),
				Partner:     v.GetPartner(),
				PartnerName: getPartnerFullName(PartnerListMaps, v.GetPartner()),
				Type:        v.GetType(),
				Status:      v.GetStatus(),
				Enabled:     v.GetEnabled(),
				Remarks:     v.GetRemarks(),
				Applied:     v.GetApplied().AsTime(),
				Created:     v.GetCreated().AsTime(),
				Updated:     v.GetUpdated().AsTime(),
				UpdatedBy:   v.GetUpdatedBy(),
				ID:          v.GetID(),
			})
		}
	}
	etd := s.getEnforceTemplateData(ctx)
	details := serviceDetails(lf.GetFileUploads())
	details.OrgID = oid
	if len(lsr.GetServiceRequst()) > 0 {
		details.Status = lsr.GetServiceRequst()[0].Status.String()
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	details.UserInfo = &usrInfo.UserInfo
	details.CSRFFieldValue = csrf.Token(r)
	details.ServiceRequests = sr
	details.ActivePartnerName = partner
	details.CompanyName = cn
	details.FileStatus = fs
	details.HasAya = ha
	details.HasWu = wa
	details.HasAddi = (ha || wa)
	details.BaseURL = s.urls.Base + "/u/files/"
	details.CSRFField = csrf.TemplateField(r)
	details.PresetPermission = etd.PresetPermission
	details.ServiceRequest = etd.ServiceRequests
	details.UserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, details); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) getChangeStatusSvcReq(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	oid := goji.Param(r, "id")
	if oid == "" {
		log.Error("missing org id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	status := goji.Param(r, "status")
	if status == "" {
		log.Error("missing status query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	partner := goji.Param(r, "partner")
	if partner == "" {
		log.Error("missing partner query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	file := goji.Param(r, "file")
	if file == "" {
		log.Error("missing file type query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	uid := mw.GetUserID(ctx)
	if uid == "" {
		log.Error("missing user id")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if status == Accept {
		if _, err := s.pf.AcceptUploadSvcRequest(ctx, &svr.AcceptUploadSvcRequestRequest{
			OrgID:    oid,
			Partner:  partner,
			SvcName:  svr.ServiceType_REMITTANCE.String(),
			FileType: file,
			VerifyBy: uid,
		}); err != nil {
			log.Error("Accept failed")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	} else if status == "reject" {
		if _, err := s.pf.RejectUploadSvcRequest(ctx, &svr.RejectUploadSvcRequestRequest{
			OrgID:    oid,
			Partner:  partner,
			SvcName:  svr.ServiceType_REMITTANCE.String(),
			FileType: file,
			VerifyBy: uid,
		}); err != nil {
			log.Error("Reject failed")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}

	partnerU := ""
	if partner != "" {
		partnerU = "?partner=" + partner
	}
	http.Redirect(w, r, "/dashboard/services-details/"+oid+partnerU, http.StatusSeeOther)
}

func (s *Server) getChangeAllStatusSvcReq(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	oid := goji.Param(r, "id")
	if oid == "" {
		log.Error("missing org id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	status := goji.Param(r, "status")
	if status == "" {
		log.Error("missing status query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	partner := goji.Param(r, "partner")
	if partner == "" {
		log.Error("missing partner query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	uid := mw.GetUserID(ctx)
	if uid == "" {
		log.Error("missing user id")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if status == Accept {
		if _, err := s.pf.AcceptServiceRequest(ctx, &svr.ServiceStatusRequestRequest{
			OrgID:     oid,
			UpdatedBy: uid,
			Partner:   partner,
			SvcName:   svr.ServiceType_REMITTANCE.String(),
		}); err != nil {
			log.WithError(err).Error("Accept failed")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}

	if status == Reject {
		if _, err := s.pf.RejectServiceRequest(ctx, &svr.ServiceStatusRequestRequest{
			OrgID:     oid,
			UpdatedBy: uid,
			Partner:   partner,
			SvcName:   svr.ServiceType_REMITTANCE.String(),
		}); err != nil {
			log.WithError(err).Error("Reject failed")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}
	http.Redirect(w, r, "/dashboard/services-details/"+oid, http.StatusSeeOther)
}

func (s *Server) getChangeCiCoAllStatusSvcReq(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	oid := goji.Param(r, "id")
	if oid == "" {
		log.Error("missing org id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	status := goji.Param(r, "status")
	if status == "" {
		log.Error("missing status query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	partner := goji.Param(r, "partner")
	if partner == "" {
		log.Error("missing partner query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	uid := mw.GetUserID(ctx)
	if uid == "" {
		log.Error("missing user id")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if status == Accept {
		if _, err := s.pf.AcceptServiceRequest(ctx, &svr.ServiceStatusRequestRequest{
			OrgID:     oid,
			UpdatedBy: uid,
			Partner:   partner,
			SvcName:   svr.ServiceType_CASHINCASHOUT.String(),
		}); err != nil {
			log.Error("Accept failed")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	} else if status == Reject {
		if _, err := s.pf.RejectServiceRequest(ctx, &svr.ServiceStatusRequestRequest{
			OrgID:     oid,
			UpdatedBy: uid,
			Partner:   partner,
			SvcName:   svr.ServiceType_CASHINCASHOUT.String(),
		}); err != nil {
			log.Error("Reject failed")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}
	http.Redirect(w, r, "/dashboard/services-details-cico/"+oid, http.StatusSeeOther)
}

func (s *Server) getChangeMIAllStatusSvcReq(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	oid := goji.Param(r, "id")
	if oid == "" {
		log.Error("missing org id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	status := goji.Param(r, "status")
	if status == "" {
		log.Error("missing status query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	partner := goji.Param(r, "partner")
	if partner == "" {
		log.Error("missing partner query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	uid := mw.GetUserID(ctx)
	if uid == "" {
		log.Error("missing user id")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if status == Accept {
		if _, err := s.pf.AcceptServiceRequest(ctx, &svr.ServiceStatusRequestRequest{
			OrgID:     oid,
			UpdatedBy: uid,
			Partner:   partner,
			SvcName:   svr.ServiceType_MICROINSURANCE.String(),
		}); err != nil {
			log.Error("Accept failed")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	} else if status == Reject {
		if _, err := s.pf.RejectServiceRequest(ctx, &svr.ServiceStatusRequestRequest{
			OrgID:     oid,
			UpdatedBy: uid,
			Partner:   partner,
			SvcName:   svr.ServiceType_MICROINSURANCE.String(),
		}); err != nil {
			log.Error("Reject failed")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}
	http.Redirect(w, r, "/dashboard/services-details-mi/"+oid, http.StatusSeeOther)
}

func serviceDetails(fs []*fpb.FileUpload) ServicesDetails {
	fm := ServicesDetails{}
	for _, f := range fs {
		switch f.Type {
		case fpb.UploadType_Mbf:
			fm.Files.Mbf = f.GetFileNames()
			fm.Files.MbfName = f.GetFileName()
		case fpb.UploadType_Sec:
			fm.Files.Sec = f.GetFileNames()
			fm.Files.SecName = f.GetFileName()
		case fpb.UploadType_Gis:
			fm.Files.Gis = f.GetFileNames()
			fm.Files.GisName = f.GetFileName()
		case fpb.UploadType_Afs:
			fm.Files.Afs = f.GetFileNames()
			fm.Files.AfsName = f.GetFileName()
		case fpb.UploadType_Brs:
			fm.Files.Brs = f.GetFileNames()
			fm.Files.BrsName = f.GetFileName()
		case fpb.UploadType_Bmp:
			fm.Files.Bmp = f.GetFileNames()
			fm.Files.BmpName = f.GetFileName()
		case fpb.UploadType_Scbr:
			fm.Files.Scbr = f.GetFileNames()
			fm.Files.ScbrName = f.GetFileName()
		case fpb.UploadType_Via:
			fm.Files.Via = f.GetFileNames()
			fm.Files.ViaName = f.GetFileName()
		case fpb.UploadType_WU_cd:
			fm.Files.WU_cd = f.GetFileNames()
			fm.Files.WU_cdName = f.GetFileName()
		case fpb.UploadType_WU_lbp:
			fm.Files.WU_lbp = f.GetFileNames()
			fm.Files.WU_lbpName = f.GetFileName()
		case fpb.UploadType_WU_sr:
			fm.Files.WU_sr = f.GetFileNames()
			fm.Files.WU_srName = f.GetFileName()
		case fpb.UploadType_WU_lgis:
			fm.Files.WU_lgis = f.GetFileNames()
			fm.Files.WU_lgisName = f.GetFileName()
		case fpb.UploadType_WU_dtisspa:
			fm.Files.WU_dtisspa = f.GetFileNames()
			fm.Files.WU_dtisspaName = f.GetFileName()
		case fpb.UploadType_WU_birf:
			fm.Files.WU_birf = f.GetFileNames()
			fm.Files.WU_birfName = f.GetFileName()
		case fpb.UploadType_WU_bspr:
			fm.Files.WU_bspr = f.GetFileNames()
			fm.Files.WU_bsprName = f.GetFileName()
		case fpb.UploadType_WU_iqa:
			fm.Files.WU_iqa = f.GetFileNames()
			fm.Files.WU_iqaName = f.GetFileName()
		case fpb.UploadType_WU_bqa:
			fm.Files.WU_bqa = f.GetFileNames()
			fm.Files.WU_bqaName = f.GetFileName()
		case fpb.UploadType_AYA_dfedr:
			fm.Files.AYA_dfedr = f.GetFileNames()
			fm.Files.AYA_dfedrName = f.GetFileName()
		case fpb.UploadType_AYA_ddf:
			fm.Files.AYA_ddf = f.GetFileNames()
			fm.Files.AYA_ddfName = f.GetFileName()
		case fpb.UploadType_AYA_ialaws:
			fm.Files.AYA_ialaws = f.GetFileNames()
			fm.Files.AYA_ialawsName = f.GetFileName()
		case fpb.UploadType_AYA_aialaws:
			fm.Files.AYA_aialaws = f.GetFileNames()
			fm.Files.AYA_aialawsName = f.GetFileName()
		case fpb.UploadType_AYA_mtl:
			fm.Files.AYA_mtl = f.GetFileNames()
			fm.Files.AYA_mtlName = f.GetFileName()
		case fpb.UploadType_AYA_cbpr:
			fm.Files.AYA_cbpr = f.GetFileNames()
			fm.Files.AYA_cbprName = f.GetFileName()
		case fpb.UploadType_AYA_brdcsa:
			fm.Files.AYA_brdcsa = f.GetFileNames()
			fm.Files.AYA_brdcsaName = f.GetFileName()
		case fpb.UploadType_AYA_gis:
			fm.Files.AYA_gis = f.GetFileNames()
			fm.Files.AYA_gisName = f.GetFileName()
		case fpb.UploadType_AYA_ccpwi:
			fm.Files.AYA_ccpwi = f.GetFileNames()
			fm.Files.AYA_ccpwiName = f.GetFileName()
		case fpb.UploadType_AYA_fas:
			fm.Files.AYA_fas = f.GetFileNames()
			fm.Files.AYA_fasName = f.GetFileName()
		case fpb.UploadType_AYA_am:
			fm.Files.AYA_am = f.GetFileNames()
			fm.Files.AYA_amName = f.GetFileName()
		case fpb.UploadType_AYA_laf:
			fm.Files.AYA_laf = f.GetFileNames()
			fm.Files.AYA_lafName = f.GetFileName()
		case fpb.UploadType_AYA_birr:
			fm.Files.AYA_birr = f.GetFileNames()
			fm.Files.AYA_birrName = f.GetFileName()
		case fpb.UploadType_AYA_od:
			fm.Files.AYA_od = f.GetFileNames()
			fm.Files.AYA_odName = f.GetFileName()
		case fpb.UploadType_Bspr:
			fm.Files.Bspr = f.GetFileNames()
			fm.Files.BsprName = f.GetFileName()
		case fpb.UploadType_Cp:
			fm.Files.Cp = f.GetFileNames()
			fm.Files.CpName = f.GetFileName()
		case fpb.UploadType_Aml:
			fm.Files.Aml = f.GetFileNames()
			fm.Files.AmlName = f.GetFileName()
		case fpb.UploadType_Nnda:
			fm.Files.Nnda = f.GetFileNames()
			fm.Files.NndaName = f.GetFileName()
		case fpb.UploadType_Psf:
			fm.Files.Psf = f.GetFileNames()
			fm.Files.PsfName = f.GetFileName()
		case fpb.UploadType_Psp:
			fm.Files.Psp = f.GetFileNames()
			fm.Files.PspName = f.GetFileName()
		case fpb.UploadType_Kddq:
			fm.Files.Kddq = f.GetFileNames()
			fm.Files.KddqName = f.GetFileName()
		case fpb.UploadType_Sis:
			fm.Files.Sis = f.GetFileNames()
			fm.Files.SisName = f.GetFileName()
		}
	}
	return fm
}

func (s *Server) getChangeBillsPaymentAllStatusSvcReq(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	oid := goji.Param(r, "id")
	if oid == "" {
		log.Error("missing org id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	status := goji.Param(r, "status")
	if status == "" {
		log.Error("missing status query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	partner := goji.Param(r, "partner")
	if partner == "" {
		log.Error("missing partner query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	uid := mw.GetUserID(ctx)
	if uid == "" {
		log.Error("missing user id")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if status == Accept {
		if _, err := s.pf.AcceptServiceRequest(ctx, &svr.ServiceStatusRequestRequest{
			OrgID:     oid,
			UpdatedBy: uid,
			Partner:   partner,
			SvcName:   svr.ServiceType_BILLSPAYMENT.String(),
		}); err != nil {
			log.Error("Accept failed")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	} else if status == Reject {
		if _, err := s.pf.RejectServiceRequest(ctx, &svr.ServiceStatusRequestRequest{
			OrgID:     oid,
			UpdatedBy: uid,
			Partner:   partner,
			SvcName:   svr.ServiceType_BILLSPAYMENT.String(),
		}); err != nil {
			log.Error("Reject failed")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}
	http.Redirect(w, r, "/dashboard/services-details-bills-payment/"+oid, http.StatusSeeOther)
}

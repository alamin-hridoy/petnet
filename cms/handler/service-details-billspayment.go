package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"time"

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
	ServicesDetailsBillsPayment struct {
		CSRFField                   template.HTML
		CSRFFieldValue              string
		ID                          string
		OrgID                       string
		Status                      string
		FileStatus                  map[string]string
		ServiceRequestsBillsPayment []*ServiceRequestsBillsPayment
		CompanyName                 string
		Files                       ReqFiles
		BaseURL                     string
		ActivePartnerName           string
		UserInfo                    *User
		PresetPermission            map[string]map[string]bool
		ServiceRequest              bool
		ReqFilesDetails             map[string]string
	}

	ServiceRequestsBillsPayment struct {
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

const (
	Accept = "accept"
	Reject = "reject"
)

func (s *Server) getServicesDetailsBillsPayment(w http.ResponseWriter, r *http.Request) {
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
	template := s.templates.Lookup("service-details-bills-payment.html")
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
		Types:    []svr.ServiceType{svr.ServiceType_BILLSPAYMENT},
		Statuses: []svr.ServiceRequestStatus{svr.ServiceRequestStatus_PENDING, svr.ServiceRequestStatus_ACCEPTED, svr.ServiceRequestStatus_REJECTED},
	})
	if err != nil {
		http.Redirect(w, r, serviceReqPath, http.StatusSeeOther)
		logging.WithError(err, log).Info("List Service Request")
	}

	BillsPaymentPartnerListMaps := s.getPartnerListMaps(ctx, svr.ServiceType_BILLSPAYMENT.String())
	fs := map[string]string{}
	ulsr, err := s.pf.ListUploadSvcRequest(ctx, &svr.ListUploadSvcRequestRequest{
		OrgID:    oid,
		SvcNames: []string{svr.ServiceType_BILLSPAYMENT.String()},
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
		Types: []fpb.UploadType{fpb.UploadType_BP_nda, fpb.UploadType_BP_sis, fpb.UploadType_BP_psf, fpb.UploadType_BP_psp, fpb.UploadType_BP_sec, fpb.UploadType_BP_aib, fpb.UploadType_BP_gis, fpb.UploadType_BP_lpafs, fpb.UploadType_BP_bir, fpb.UploadType_BP_bp, fpb.UploadType_BP_bsp, fpb.UploadType_BP_aml, fpb.UploadType_BP_sccas, fpb.UploadType_BP_vg, fpb.UploadType_BP_cp, fpb.UploadType_BP_moa, fpb.UploadType_BP_amla, fpb.UploadType_BP_mttp, fpb.UploadType_BP_sci, fpb.UploadType_BP_edd},
	})
	_, err = s.rbac.GetUser(ctx, &rbupb.GetUserRequest{ID: pf.Profile.UserID})
	if err != nil {
		logging.WithError(err, log).Info("getting user")
		return
	}
	cn := ""
	sr := []*ServiceRequestsBillsPayment{}
	if len(lsr.GetServiceRequst()) > 0 {
		for _, v := range lsr.GetServiceRequst() {
			cn = v.GetCompanyName()
			sr = append(sr, &ServiceRequestsBillsPayment{
				OrgID:       v.GetOrgID(),
				CompanyName: v.GetCompanyName(),
				Partner:     v.GetPartner(),
				PartnerName: getPartnerFullName(BillsPaymentPartnerListMaps, v.GetPartner()),
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
	details := serviceDetailsBillsPayment(lf.GetFileUploads())
	details.OrgID = oid
	if len(lsr.GetServiceRequst()) > 0 {
		details.Status = lsr.GetServiceRequst()[0].Status.String()
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	details.UserInfo = &usrInfo.UserInfo
	details.CSRFFieldValue = csrf.Token(r)
	details.ServiceRequestsBillsPayment = sr
	details.ActivePartnerName = partner
	details.CompanyName = cn
	details.FileStatus = fs
	details.BaseURL = s.urls.Base + "/u/files/"
	details.CSRFField = csrf.TemplateField(r)
	details.PresetPermission = etd.PresetPermission
	details.ServiceRequest = etd.ServiceRequests
	details.UserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, details); err != nil {
		log.WithError(err).Error("error with template execution")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) getChangeAllStatusSvcReqBillsPayment(w http.ResponseWriter, r *http.Request) {
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
			SvcName:  svr.ServiceType_BILLSPAYMENT.String(),
			FileType: file,
			VerifyBy: uid,
		}); err != nil {
			log.WithError(err).Error("Accept failed")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}

	if status == Reject {
		if _, err := s.pf.RejectUploadSvcRequest(ctx, &svr.RejectUploadSvcRequestRequest{
			OrgID:    oid,
			Partner:  partner,
			SvcName:  svr.ServiceType_BILLSPAYMENT.String(),
			FileType: file,
			VerifyBy: uid,
		}); err != nil {
			log.WithError(err).Error("Reject failed")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}
	http.Redirect(w, r, "/dashboard/services-details-bills-payment/"+oid, http.StatusSeeOther)
}

func (s *Server) getAjaxChangeAllStatusSvcReqBillsPayment(w http.ResponseWriter, r *http.Request) {
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
			SvcName:  svr.ServiceType_BILLSPAYMENT.String(),
			FileType: file,
			VerifyBy: uid,
		}); err != nil {
			log.WithError(err).Error("Accept failed")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}
	if status == Reject {
		if _, err := s.pf.RejectUploadSvcRequest(ctx, &svr.RejectUploadSvcRequestRequest{
			OrgID:    oid,
			Partner:  partner,
			SvcName:  svr.ServiceType_BILLSPAYMENT.String(),
			FileType: file,
			VerifyBy: uid,
		}); err != nil {
			log.WithError(err).Error("Reject failed")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}
	jsn, _ := json.Marshal([]string{})
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsn)
}

func serviceDetailsBillsPayment(fs []*fpb.FileUpload) ServicesDetailsBillsPayment {
	fm := ServicesDetailsBillsPayment{}
	for _, f := range fs {
		switch f.Type {
		case fpb.UploadType_BP_nda:
			fm.Files.BP_nda = f.GetFileNames()
			fm.Files.BP_ndaName = f.GetFileName()
		case fpb.UploadType_BP_sis:
			fm.Files.BP_sis = f.GetFileNames()
			fm.Files.BP_sisName = f.GetFileName()
		case fpb.UploadType_BP_psf:
			fm.Files.BP_psf = f.GetFileNames()
			fm.Files.BP_psfName = f.GetFileName()
		case fpb.UploadType_BP_psp:
			fm.Files.BP_psp = f.GetFileNames()
			fm.Files.BP_pspName = f.GetFileName()
		case fpb.UploadType_BP_sec:
			fm.Files.BP_sec = f.GetFileNames()
			fm.Files.BP_secName = f.GetFileName()
		case fpb.UploadType_BP_aib:
			fm.Files.BP_aib = f.GetFileNames()
			fm.Files.BP_aibName = f.GetFileName()
		case fpb.UploadType_BP_gis:
			fm.Files.BP_gis = f.GetFileNames()
			fm.Files.BP_gisName = f.GetFileName()
		case fpb.UploadType_BP_lpafs:
			fm.Files.BP_lpafs = f.GetFileNames()
			fm.Files.BP_lpafsName = f.GetFileName()
		case fpb.UploadType_BP_bir:
			fm.Files.BP_bir = f.GetFileNames()
			fm.Files.BP_birName = f.GetFileName()
		case fpb.UploadType_BP_bp:
			fm.Files.BP_bp = f.GetFileNames()
			fm.Files.BP_bpName = f.GetFileName()
		case fpb.UploadType_BP_bsp:
			fm.Files.BP_bsp = f.GetFileNames()
			fm.Files.BP_bspName = f.GetFileName()
		case fpb.UploadType_BP_aml:
			fm.Files.BP_aml = f.GetFileNames()
			fm.Files.BP_amlName = f.GetFileName()
		case fpb.UploadType_BP_sccas:
			fm.Files.BP_sccas = f.GetFileNames()
			fm.Files.BP_sccasName = f.GetFileName()
		case fpb.UploadType_BP_vg:
			fm.Files.BP_vg = f.GetFileNames()
			fm.Files.BP_vgName = f.GetFileName()
		case fpb.UploadType_BP_cp:
			fm.Files.BP_cp = f.GetFileNames()
			fm.Files.BP_cpName = f.GetFileName()
		case fpb.UploadType_BP_moa:
			fm.Files.BP_moa = f.GetFileNames()
			fm.Files.BP_moaName = f.GetFileName()
		case fpb.UploadType_BP_amla:
			fm.Files.BP_amla = f.GetFileNames()
			fm.Files.BP_amlaName = f.GetFileName()
		case fpb.UploadType_BP_mttp:
			fm.Files.BP_mttp = f.GetFileNames()
			fm.Files.BP_mttpName = f.GetFileName()
		case fpb.UploadType_BP_sci:
			fm.Files.BP_sci = f.GetFileNames()
			fm.Files.BP_sciName = f.GetFileName()
		case fpb.UploadType_BP_edd:
			fm.Files.BP_edd = f.GetFileNames()
			fm.Files.BP_eddName = f.GetFileName()
		}
	}
	return fm
}

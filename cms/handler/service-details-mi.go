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
	ServicesDetailsMI struct {
		CSRFField         template.HTML
		CSRFFieldValue    string
		ID                string
		OrgID             string
		Status            string
		FileStatus        map[string]string
		ServiceRequestsMI []*ServiceRequestsMI
		CompanyName       string
		Files             ReqFiles
		BaseURL           string
		ActivePartnerName string
		UserInfo          *User
		PresetPermission  map[string]map[string]bool
		ServiceRequest    bool
		ReqFilesDetails   map[string]string
	}

	ServiceRequestsMI struct {
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

var MIPartnerListMaps = map[string]string{"RuralNet": "RuralNet"}

func (s *Server) getServicesDetailsMI(w http.ResponseWriter, r *http.Request) {
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
	if err != nil || partner == "" {
		partner = "ALL"
	}
	template := s.templates.Lookup("service-details-mi.html")
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
		Types:    []svr.ServiceType{svr.ServiceType_MICROINSURANCE},
		Statuses: []svr.ServiceRequestStatus{svr.ServiceRequestStatus_PENDING, svr.ServiceRequestStatus_ACCEPTED, svr.ServiceRequestStatus_REJECTED},
	})
	if err != nil {
		http.Redirect(w, r, serviceReqPath, http.StatusSeeOther)
		logging.WithError(err, log).Info("List Service Request")
	}
	fs := map[string]string{}
	ulsr, err := s.pf.ListUploadSvcRequest(ctx, &svr.ListUploadSvcRequestRequest{
		OrgID:    oid,
		SvcNames: []string{svr.ServiceType_MICROINSURANCE.String()},
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
		Types: []fpb.UploadType{fpb.UploadType_MI_mbf, fpb.UploadType_MI_nda, fpb.UploadType_MI_sec, fpb.UploadType_MI_gis, fpb.UploadType_MI_afs, fpb.UploadType_MI_bir, fpb.UploadType_MI_scb, fpb.UploadType_MI_via, fpb.UploadType_MI_moa},
	})
	_, err = s.rbac.GetUser(ctx, &rbupb.GetUserRequest{ID: pf.Profile.UserID})
	if err != nil {
		logging.WithError(err, log).Info("getting user")
		return
	}
	cn := ""
	sr := []*ServiceRequestsMI{}
	if len(lsr.GetServiceRequst()) > 0 {
		for _, v := range lsr.GetServiceRequst() {
			cn = v.GetCompanyName()
			sr = append(sr, &ServiceRequestsMI{
				OrgID:       v.GetOrgID(),
				CompanyName: v.GetCompanyName(),
				Partner:     v.GetPartner(),
				PartnerName: getPartnerFullName(MIPartnerListMaps, v.GetPartner()),
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
	details := serviceDetailsMI(lf.GetFileUploads())
	details.OrgID = oid
	if len(lsr.GetServiceRequst()) > 0 {
		details.Status = lsr.GetServiceRequst()[0].Status.String()
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	details.UserInfo = &usrInfo.UserInfo
	details.CSRFFieldValue = csrf.Token(r)
	details.ServiceRequestsMI = sr
	details.ActivePartnerName = partner
	details.CompanyName = cn
	details.FileStatus = fs
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

func (s *Server) getChangeAllStatusSvcReqMI(w http.ResponseWriter, r *http.Request) {
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
			SvcName:  svr.ServiceType_MICROINSURANCE.String(),
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
			SvcName:  svr.ServiceType_MICROINSURANCE.String(),
			FileType: file,
			VerifyBy: uid,
		}); err != nil {
			log.WithError(err).Error("Reject failed")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}
	http.Redirect(w, r, "/dashboard/services-details-mi/"+oid, http.StatusSeeOther)
}

func (s *Server) getAjaxChangeAllStatusSvcReqMI(w http.ResponseWriter, r *http.Request) {
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
			SvcName:  svr.ServiceType_MICROINSURANCE.String(),
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
			SvcName:  svr.ServiceType_MICROINSURANCE.String(),
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

func serviceDetailsMI(fs []*fpb.FileUpload) ServicesDetailsMI {
	fm := ServicesDetailsMI{}
	for _, f := range fs {
		switch f.Type {
		case fpb.UploadType_MI_mbf:
			fm.Files.MI_mbf = f.GetFileNames()
			fm.Files.MI_mbfName = f.GetFileName()
		case fpb.UploadType_MI_nda:
			fm.Files.MI_nda = f.GetFileNames()
			fm.Files.MI_ndaName = f.GetFileName()
		case fpb.UploadType_MI_sec:
			fm.Files.MI_sec = f.GetFileNames()
			fm.Files.MI_secName = f.GetFileName()
		case fpb.UploadType_MI_gis:
			fm.Files.MI_gis = f.GetFileNames()
			fm.Files.MI_gisName = f.GetFileName()
		case fpb.UploadType_MI_afs:
			fm.Files.MI_afs = f.GetFileNames()
			fm.Files.MI_afsName = f.GetFileName()
		case fpb.UploadType_MI_bir:
			fm.Files.MI_bir = f.GetFileNames()
			fm.Files.MI_birName = f.GetFileName()
		case fpb.UploadType_MI_scb:
			fm.Files.MI_scb = f.GetFileNames()
			fm.Files.MI_scbName = f.GetFileName()
		case fpb.UploadType_MI_via:
			fm.Files.MI_via = f.GetFileNames()
			fm.Files.MI_viaName = f.GetFileName()
		case fpb.UploadType_MI_moa:
			fm.Files.MI_moa = f.GetFileNames()
			fm.Files.MI_moaName = f.GetFileName()
		}
	}
	return fm
}

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
	ServicesDetailsCico struct {
		CSRFField           template.HTML
		CSRFFieldValue      string
		ID                  string
		OrgID               string
		Status              string
		FileStatus          map[string]string
		ServiceRequestsCico []*ServiceRequestsCico
		CompanyName         string
		Files               ReqFiles
		BaseURL             string
		ActivePartnerName   string
		UserInfo            *User
		PresetPermission    map[string]map[string]bool
		ServiceRequest      bool
		ReqFilesDetails     map[string]string
	}

	ServiceRequestsCico struct {
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

func (s *Server) getServicesDetailsCico(w http.ResponseWriter, r *http.Request) {
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
	template := s.templates.Lookup("service-details-cico.html")
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
		Types:    []svr.ServiceType{svr.ServiceType_CASHINCASHOUT},
		Statuses: []svr.ServiceRequestStatus{svr.ServiceRequestStatus_PENDING, svr.ServiceRequestStatus_ACCEPTED, svr.ServiceRequestStatus_REJECTED},
	})
	if err != nil {
		http.Redirect(w, r, serviceReqPath, http.StatusSeeOther)
		logging.WithError(err, log).Info("List Service Request")
	}

	CiCoPartnerListMaps := s.getPartnerListMaps(ctx, svr.ServiceType_CASHINCASHOUT.String())
	fs := map[string]string{}
	ulsr, err := s.pf.ListUploadSvcRequest(ctx, &svr.ListUploadSvcRequestRequest{
		OrgID:    oid,
		SvcNames: []string{svr.ServiceType_CASHINCASHOUT.String()},
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
		Types: []fpb.UploadType{fpb.UploadType_CICO_nda, fpb.UploadType_CICO_sis, fpb.UploadType_CICO_psf, fpb.UploadType_CICO_pspp, fpb.UploadType_CICO_sec, fpb.UploadType_CICO_gis, fpb.UploadType_CICO_afs, fpb.UploadType_CICO_bir, fpb.UploadType_CICO_bsp, fpb.UploadType_CICO_aml, fpb.UploadType_CICO_sccas, fpb.UploadType_CICO_vgid, fpb.UploadType_CICO_cp, fpb.UploadType_CICO_moa, fpb.UploadType_CICO_amla, fpb.UploadType_CICO_mtpp, fpb.UploadType_CICO_is, fpb.UploadType_CICO_edd},
	})
	_, err = s.rbac.GetUser(ctx, &rbupb.GetUserRequest{ID: pf.Profile.UserID})
	if err != nil {
		logging.WithError(err, log).Info("getting user")
		return
	}
	cn := ""
	sr := []*ServiceRequestsCico{}
	if len(lsr.GetServiceRequst()) > 0 {
		for _, v := range lsr.GetServiceRequst() {
			cn = v.GetCompanyName()
			sr = append(sr, &ServiceRequestsCico{
				OrgID:       v.GetOrgID(),
				CompanyName: v.GetCompanyName(),
				Partner:     v.GetPartner(),
				PartnerName: getPartnerFullName(CiCoPartnerListMaps, v.GetPartner()),
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
	details := serviceDetailsCico(lf.GetFileUploads())
	details.OrgID = oid
	if len(lsr.GetServiceRequst()) > 0 {
		details.Status = lsr.GetServiceRequst()[0].Status.String()
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	details.UserInfo = &usrInfo.UserInfo
	details.CSRFFieldValue = csrf.Token(r)
	details.ServiceRequestsCico = sr
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

func (s *Server) getChangeAllStatusSvcReqCico(w http.ResponseWriter, r *http.Request) {
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
			SvcName:  svr.ServiceType_CASHINCASHOUT.String(),
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
			SvcName:  svr.ServiceType_CASHINCASHOUT.String(),
			FileType: file,
			VerifyBy: uid,
		}); err != nil {
			log.WithError(err).Error("Reject failed")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}
	http.Redirect(w, r, "/dashboard/services-details-cico/"+oid, http.StatusSeeOther)
}

func (s *Server) getAjaxChangeAllStatusSvcReqCico(w http.ResponseWriter, r *http.Request) {
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
			SvcName:  svr.ServiceType_CASHINCASHOUT.String(),
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
			SvcName:  svr.ServiceType_CASHINCASHOUT.String(),
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

func serviceDetailsCico(fs []*fpb.FileUpload) ServicesDetailsCico {
	fm := ServicesDetailsCico{}
	for _, f := range fs {
		switch f.Type {
		case fpb.UploadType_CICO_nda:
			fm.Files.CICO_nda = f.GetFileNames()
			fm.Files.CICO_ndaName = f.GetFileName()
		case fpb.UploadType_CICO_sis:
			fm.Files.CICO_sis = f.GetFileNames()
			fm.Files.CICO_sisName = f.GetFileName()
		case fpb.UploadType_CICO_psf:
			fm.Files.CICO_psf = f.GetFileNames()
			fm.Files.CICO_psfName = f.GetFileName()
		case fpb.UploadType_CICO_pspp:
			fm.Files.CICO_pspp = f.GetFileNames()
			fm.Files.CICO_psppName = f.GetFileName()
		case fpb.UploadType_CICO_sec:
			fm.Files.CICO_sec = f.GetFileNames()
			fm.Files.CICO_secName = f.GetFileName()
		case fpb.UploadType_CICO_gis:
			fm.Files.CICO_gis = f.GetFileNames()
			fm.Files.CICO_gisName = f.GetFileName()
		case fpb.UploadType_CICO_afs:
			fm.Files.CICO_afs = f.GetFileNames()
			fm.Files.CICO_afsName = f.GetFileName()
		case fpb.UploadType_CICO_bir:
			fm.Files.CICO_bir = f.GetFileNames()
			fm.Files.CICO_birName = f.GetFileName()
		case fpb.UploadType_CICO_bsp:
			fm.Files.CICO_bsp = f.GetFileNames()
			fm.Files.CICO_bspName = f.GetFileName()
		case fpb.UploadType_CICO_aml:
			fm.Files.CICO_aml = f.GetFileNames()
			fm.Files.CICO_amlName = f.GetFileName()
		case fpb.UploadType_CICO_sccas:
			fm.Files.CICO_sccas = f.GetFileNames()
			fm.Files.CICO_sccasName = f.GetFileName()
		case fpb.UploadType_CICO_vgid:
			fm.Files.CICO_vgid = f.GetFileNames()
			fm.Files.CICO_vgidName = f.GetFileName()
		case fpb.UploadType_CICO_cp:
			fm.Files.CICO_cp = f.GetFileNames()
			fm.Files.CICO_cpName = f.GetFileName()
		case fpb.UploadType_CICO_moa:
			fm.Files.CICO_moa = f.GetFileNames()
			fm.Files.CICO_moaName = f.GetFileName()
		case fpb.UploadType_CICO_amla:
			fm.Files.CICO_amla = f.GetFileNames()
			fm.Files.CICO_amlaName = f.GetFileName()
		case fpb.UploadType_CICO_mtpp:
			fm.Files.CICO_mtpp = f.GetFileNames()
			fm.Files.CICO_mtppName = f.GetFileName()
		case fpb.UploadType_CICO_is:
			fm.Files.CICO_is = f.GetFileNames()
			fm.Files.CICO_isName = f.GetFileName()
		case fpb.UploadType_CICO_edd:
			fm.Files.CICO_edd = f.GetFileNames()
			fm.Files.CICO_eddName = f.GetFileName()
		}
	}
	return fm
}

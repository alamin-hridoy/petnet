package handler

import (
	"errors"
	"html/template"
	"net/http"

	cmsmw "brank.as/petnet/cms/mw"
	ppf "brank.as/petnet/gunk/dsa/v1/user"
	spb "brank.as/petnet/gunk/dsa/v2/partner"
	spbl "brank.as/petnet/gunk/dsa/v2/partnerlist"
	sVcpb "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
)

type (
	DsaSelectPartnerCiCoForm struct {
		SaveDraft             string
		SaveContinue          string
		CSRFField             template.HTML
		PartnerName           []string
		SelectedPartners      []string
		Errors                map[string]error
		PresetPermission      map[string]map[string]bool
		PartnerListApplicants []PartnerListApplicant
		ServiceRequest        bool
		UserInfo              *User
		CompanyName           string
	}
)

func (s *Server) getDsaSelectCiCoPartners(w http.ResponseWriter, r *http.Request) {
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
	template := s.templates.Lookup("dsa-services-select-partners-cashin-cashout.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	CiCoPartnerListMaps := s.getPartnerListMaps(ctx, sVcpb.ServiceType_CASHINCASHOUT.String())
	ptnr := ""
	newReq := &spbl.GetPartnerListRequest{
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_CASHINCASHOUT.String(),
	}
	if ptnr != "" {
		newReq.Stype = ptnr
	}
	svc, err := s.pf.GetPartnerList(ctx, newReq)
	var partnerListApplicantList []PartnerListApplicant
	ap := PartnerListApplicant{
		Stype: AllPartners,
		Name:  AllPartners,
	}
	partnerListApplicantList = append(partnerListApplicantList, ap)
	if err != nil {
		log.Error("failed to Get Partner List")
	} else {
		for _, sv := range svc.GetPartnerList() {
			partnerListApplicantList = append(partnerListApplicantList, PartnerListApplicant{
				Stype: sv.GetStype(),
				Name:  sv.GetName(),
			})
		}
	}
	var selectedPartnerList []string
	var availablePartners []string
	var acceptedPartners []string

	ress, _ := s.pf.ListServiceRequest(ctx, &sVcpb.ListServiceRequestRequest{
		OrgIDs: []string{oid},
		Types:  []sVcpb.ServiceType{sVcpb.ServiceType_CASHINCASHOUT},
	})
	if ress != nil {
		for _, vv := range ress.GetServiceRequst() {
			if vv.Status != sVcpb.ServiceRequestStatus_ACCEPTED && vv.Status != sVcpb.ServiceRequestStatus_PENDING && vv.Status != sVcpb.ServiceRequestStatus_REJECTED {
				selectedPartnerList = append(selectedPartnerList, getPartnerFullName(CiCoPartnerListMaps, vv.Partner))
				availablePartners = append(availablePartners, vv.Partner)
			}
			if vv.Status == sVcpb.ServiceRequestStatus_ACCEPTED || vv.Status == sVcpb.ServiceRequestStatus_PENDING {
				acceptedPartners = append(acceptedPartners, vv.Partner)
			}
		}
	}
	if partnerListApplicantList != nil && acceptedPartners != nil {
		partnerListApplicantList = removeAcceptedPartnerList(partnerListApplicantList, acceptedPartners)
	}

	hv, _ := cmsmw.InArray(ap, partnerListApplicantList)
	cpl := len(partnerListApplicantList)
	if hv {
		cpl = cpl - 1
	}
	if len(selectedPartnerList) >= cpl {
		selectedPartnerList = []string{AllPartners}
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	etd := s.getEnforceTemplateData(ctx)
	data := DsaSelectPartnerCiCoForm{
		CSRFField:             csrf.TemplateField(r),
		Errors:                map[string]error{},
		PartnerName:           []string{},
		ServiceRequest:        etd.ServiceRequests,
		PartnerListApplicants: partnerListApplicantList,
		SelectedPartners:      selectedPartnerList,
		UserInfo:              &usrInfo.UserInfo,
		CompanyName:           usrInfo.CompanyName,
	}
	data.UserInfo.ProfileImage = usrInfo.ProfileImage

	if err := template.Execute(w, data); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) postDsaSelectCicoPartners(w http.ResponseWriter, r *http.Request) {
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

	var form DsaSelectPartnerCiCoForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	formErr := validation.Errors{}
	if err := validation.ValidateStruct(&form,
		validation.Field(&form.PartnerName, validation.Required),
	); err != nil {
		if err, ok := (err).(validation.Errors); ok {
			if err["PartnerName"] != nil {
				formErr["PartnerName"] = errors.New("Select at least one partner")
			}
		}
		logging.WithError(err, log).Error("invalid request")
	}
	CiCoPartnerListMaps := s.getPartnerListMaps(ctx, sVcpb.ServiceType_CASHINCASHOUT.String())
	var partnerListApplicantList []PartnerListApplicant
	ap := PartnerListApplicant{
		Stype: AllPartners,
		Name:  AllPartners,
	}
	partnerListApplicantList = append(partnerListApplicantList, ap)
	ptnrLists := []string{}
	ptnr := ""
	newReq := &spbl.GetPartnerListRequest{
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_CASHINCASHOUT.String(),
	}
	if ptnr != "" {
		newReq.Stype = ptnr
	}
	svc, err := s.pf.GetPartnerList(r.Context(), newReq)
	if err != nil {
		log.Error("failed to Get CICO Partner List")
	} else {
		for _, sv := range svc.GetPartnerList() {
			ptnrLists = append(ptnrLists, sv.GetStype())
			partnerListApplicantList = append(partnerListApplicantList, PartnerListApplicant{
				Stype: sv.GetStype(),
				Name:  sv.GetName(),
			})
		}
	}
	hv, _ := cmsmw.InArray(ap, partnerListApplicantList)
	cpl := len(partnerListApplicantList)
	if hv {
		cpl = cpl - 1
	}
	var exceptAcceptedPartners []string
	var selectedPartnerList []string
	var availablePartners []string
	var acceptedPartners []string
	ress, _ := s.pf.ListServiceRequest(ctx, &sVcpb.ListServiceRequestRequest{
		OrgIDs: []string{oid},
		Types:  []sVcpb.ServiceType{sVcpb.ServiceType_CASHINCASHOUT},
	})

	if ress != nil {
		for _, vv := range ress.GetServiceRequst() {
			if vv.Status == sVcpb.ServiceRequestStatus_ACCEPTED || vv.Status == sVcpb.ServiceRequestStatus_PENDING {
				acceptedPartners = append(acceptedPartners, vv.Partner)
			}
			selectedPartnerList = append(selectedPartnerList, getPartnerFullName(CiCoPartnerListMaps, vv.Partner))
			availablePartners = append(availablePartners, vv.Partner)
			if vv.Status != sVcpb.ServiceRequestStatus_ACCEPTED && vv.Status != sVcpb.ServiceRequestStatus_PENDING {
				exceptAcceptedPartners = append(exceptAcceptedPartners, vv.Partner)
			}
		}
	}

	if partnerListApplicantList != nil && acceptedPartners != nil {
		partnerListApplicantList = removeAcceptedPartnerList(partnerListApplicantList, acceptedPartners)
	}

	if len(selectedPartnerList) >= cpl {
		selectedPartnerList = []string{AllPartners}
	}
	if len(formErr) > 0 {
		template := s.templates.Lookup("dsa-services-select-partners-cashin-cashout.html")
		if template == nil {
			log.Error("unable to load template")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		usrInfo := s.GetUserInfoFromCookie(w, r, false)

		etd := s.getEnforceTemplateData(r.Context())
		data := DsaSelectPartnerCiCoForm{
			CSRFField:             csrf.TemplateField(r),
			PartnerName:           []string{},
			SelectedPartners:      selectedPartnerList,
			Errors:                map[string]error{},
			PresetPermission:      etd.PresetPermission,
			ServiceRequest:        etd.ServiceRequests,
			PartnerListApplicants: partnerListApplicantList,
			UserInfo:              &usrInfo.UserInfo,
		}

		uid := mw.GetUserID(ctx)
		gp, err := s.pf.GetUserProfile(ctx, &ppf.GetUserProfileRequest{
			UserID: uid,
		})
		if err != nil {
			log.Error("failed to get profile")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		data.UserInfo.ProfileImage = gp.GetProfile().ProfilePicture
		if err := template.Execute(w, data); err != nil {
			log.Infof("error with template execution: %+v", err)
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}
	ptnrAppList := []string{}
	if len(partnerListApplicantList) > 0 {
		for _, v := range partnerListApplicantList {
			ptnrAppList = append(ptnrAppList, v.Stype)
		}
	}

	var allPtnr bool
	hAll, _ := cmsmw.InArray(AllPartners, form.PartnerName)
	frmPartner := getPartnerShortName(CiCoPartnerListMaps, form.PartnerName)
	if hAll || len(form.PartnerName) >= cpl {
		allPtnr = false
	}

	notAcceptedPtnr := frmPartner
	if len(acceptedPartners) > 0 {
		notAcceptedPtnr = []string{}
		for _, fp := range frmPartner {
			hv, _ := cmsmw.InArray(fp, ptnrAppList)
			if hv {
				notAcceptedPtnr = append(notAcceptedPtnr, fp)
			}
		}
	}

	if _, err := s.pf.AddServiceRequest(ctx, &sVcpb.AddServiceRequestRequest{
		OrgID:       oid,
		Type:        sVcpb.ServiceType_CASHINCASHOUT,
		Partners:    notAcceptedPtnr,
		AllPartners: allPtnr,
	}); err != nil {
		logging.WithError(err, log).Error("Add Service Request")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if form.SaveDraft == "SaveDraft" {
		s.pf.SetStatusUploadSvcRequest(ctx, &sVcpb.SetStatusUploadSvcRequestRequest{
			OrgID:    oid,
			Partners: notAcceptedPtnr,
			SvcName:  sVcpb.ServiceType_CASHINCASHOUT.String(),
			Status:   sVcpb.ServiceRequestStatus_PARTNERDRAFT,
		})
		http.Redirect(w, r, dsaServicesPath, http.StatusSeeOther)
		return
	}
	if form.SaveContinue == "SaveContinue" {
		s.pf.SetStatusUploadSvcRequest(ctx, &sVcpb.SetStatusUploadSvcRequestRequest{
			OrgID:    oid,
			Partners: notAcceptedPtnr,
			SvcName:  sVcpb.ServiceType_CASHINCASHOUT.String(),
			Status:   sVcpb.ServiceRequestStatus_REQDOCDRAFT,
		})
		http.Redirect(w, r, dsaReqDocsCicoPath, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, dsaPrtCiCoSelPath, http.StatusSeeOther)
	return
}

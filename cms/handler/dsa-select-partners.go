package handler

import (
	"context"
	"errors"
	"html/template"
	"net/http"

	"brank.as/petnet/api/core/static"
	cmsmw "brank.as/petnet/cms/mw"
	ppu "brank.as/petnet/gunk/dsa/v1/user"
	spb "brank.as/petnet/gunk/dsa/v2/partner"
	spbl "brank.as/petnet/gunk/dsa/v2/partnerlist"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	sVcpb "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
)

const AllPartners = "All partners"

type (
	DsaSelectPartnerForm struct {
		SaveDraft             string
		SaveContinue          string
		CSRFField             template.HTML
		PartnerName           []string
		SelectedPartners      []string
		Errors                map[string]error
		PresetPermission      map[string]map[string]bool
		ServiceRequest        bool
		PartnerListApplicants []PartnerListApplicant
		HasAya                bool
		HasWu                 bool
		HasAddi               bool
		UserInfo              *User
		CompanyName           string
	}
)

func (s *Server) getDsaSelectPartners(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	uid := mw.GetUserID(ctx)
	if uid == "" {
		log.Error("missing user id")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	oid := mw.GetOrgID(ctx)
	if oid == "" {
		log.Error("missing org id")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	template := s.templates.Lookup("dsa-services-select-partners.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	svc, err := s.pf.GetPartnerList(ctx, &spbl.GetPartnerListRequest{
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTANCE.String(),
	})
	if err != nil {
		log.WithError(err).Error("failed to get PartnerList")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	var partnerListApplicantList []PartnerListApplicant
	ap := PartnerListApplicant{
		Stype: AllPartners,
		Name:  AllPartners,
	}
	partnerListApplicantList = append(partnerListApplicantList, ap)

	for _, sv := range svc.GetPartnerList() {
		partnerListApplicantList = append(partnerListApplicantList, PartnerListApplicant{
			Stype: sv.GetStype(),
			Name:  sv.GetName(),
		})
	}

	var selectedPartnerList []string
	var availablePartners []string
	var acceptedPartners []string
	ress, err := s.pf.ListServiceRequest(ctx, &sVcpb.ListServiceRequestRequest{
		OrgIDs: []string{oid},
		Types:  []sVcpb.ServiceType{sVcpb.ServiceType_REMITTANCE},
	})
	if err != nil {
		log.WithError(err).Error("failed to List Service Request")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	ptnrListMap := s.getPartnerListMaps(ctx, sVcpb.ServiceType_REMITTANCE.String())

	if ress != nil {
		for _, vv := range ress.GetServiceRequst() {
			if vv.Status != sVcpb.ServiceRequestStatus_ACCEPTED && vv.Status != sVcpb.ServiceRequestStatus_PENDING && vv.Status != sVcpb.ServiceRequestStatus_REJECTED {
				selectedPartnerList = append(selectedPartnerList, getPartnerFullName(ptnrListMap, vv.Partner))
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

	ha, _ := cmsmw.InArray(static.AYACode, availablePartners)
	wa, _ := cmsmw.InArray(static.WUCode, availablePartners)
	etd := s.getEnforceTemplateData(ctx)
	hv, _ := cmsmw.InArray(ap, partnerListApplicantList)
	cpl := len(partnerListApplicantList)
	if hv {
		cpl = cpl - 1
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	data := DsaSelectPartnerForm{
		CSRFField:             csrf.TemplateField(r),
		PartnerName:           []string{},
		SelectedPartners:      selectedPartnerList,
		Errors:                map[string]error{},
		ServiceRequest:        etd.ServiceRequests,
		PartnerListApplicants: partnerListApplicantList,
		HasAya:                ha,
		HasWu:                 wa,
		UserInfo:              &usrInfo.UserInfo,
		HasAddi:               (ha || wa),
		CompanyName:           usrInfo.CompanyName,
	}
	data.UserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, data); err != nil {
		log.WithError(err).Error("error with template execution")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) postDsaSelectPartners(w http.ResponseWriter, r *http.Request) {
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

	var form DsaSelectPartnerForm
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
	svc, err := s.pf.GetPartnerList(ctx, &spbl.GetPartnerListRequest{
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTANCE.String(),
	})
	if err != nil {
		log.WithError(err).Error("failed to get PartnerList")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	var partnerListApplicantList []PartnerListApplicant
	ap := PartnerListApplicant{
		Stype: AllPartners,
		Name:  AllPartners,
	}
	partnerListApplicantList = append(partnerListApplicantList, ap)
	for _, sv := range svc.GetPartnerList() {
		partnerListApplicantList = append(partnerListApplicantList, PartnerListApplicant{
			Stype: sv.GetStype(),
			Name:  sv.GetName(),
		})
	}

	var selectedPartnerList []string
	var availablePartners []string
	var acceptedPartners []string
	var exceptAcceptedPartners []string
	ress, err := s.pf.ListServiceRequest(ctx, &sVcpb.ListServiceRequestRequest{
		OrgIDs: []string{oid},
		Types:  []sVcpb.ServiceType{sVcpb.ServiceType_REMITTANCE},
	})
	if err != nil {
		log.WithError(err).Error("failed to List Service Request")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	ptnrListMap := s.getPartnerListMaps(ctx, sVcpb.ServiceType_REMITTANCE.String())

	if ress != nil {
		for _, vv := range ress.GetServiceRequst() {
			if vv.Status == sVcpb.ServiceRequestStatus_ACCEPTED || vv.Status == sVcpb.ServiceRequestStatus_PENDING {
				acceptedPartners = append(acceptedPartners, vv.Partner)
			}
			selectedPartnerList = append(selectedPartnerList, getPartnerFullName(ptnrListMap, vv.Partner))
			availablePartners = append(availablePartners, vv.Partner)
			if vv.Status != sVcpb.ServiceRequestStatus_ACCEPTED && vv.Status != sVcpb.ServiceRequestStatus_PENDING {
				exceptAcceptedPartners = append(exceptAcceptedPartners, vv.Partner)
			}
		}
	}
	if partnerListApplicantList != nil && acceptedPartners != nil {
		partnerListApplicantList = removeAcceptedPartnerList(partnerListApplicantList, acceptedPartners)
	}
	if len(formErr) > 0 {
		template := s.templates.Lookup("dsa-services-select-partners.html")
		if template == nil {
			log.Error("unable to load template")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		etd := s.getEnforceTemplateData(r.Context())
		usrInfo := s.GetUserInfoFromCookie(w, r, false)

		ha, _ := cmsmw.InArray(static.AYACode, availablePartners)
		wa, _ := cmsmw.InArray(static.WUCode, availablePartners)
		data := DsaSelectPartnerForm{
			CSRFField:             csrf.TemplateField(r),
			PartnerName:           []string{},
			SelectedPartners:      selectedPartnerList,
			Errors:                map[string]error{},
			PresetPermission:      etd.PresetPermission,
			ServiceRequest:        etd.ServiceRequests,
			PartnerListApplicants: partnerListApplicantList,
			HasAya:                ha,
			HasWu:                 wa,
			HasAddi:               (ha || wa),
			UserInfo:              &usrInfo.UserInfo,
		}

		uid := mw.GetUserID(ctx)
		gp, err := s.pf.GetUserProfile(ctx, &ppu.GetUserProfileRequest{
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
	cpl := len(partnerListApplicantList)
	if hAll {
		cpl = cpl - 1
	}
	frmPartner := getPartnerShortName(ptnrListMap, form.PartnerName)
	if hAll || len(form.PartnerName) >= cpl {
		allPtnr = false
	}
	notAcceptedPtnr := frmPartner
	if len(exceptAcceptedPartners) > 0 {
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
		Type:        sVcpb.ServiceType_REMITTANCE,
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
			SvcName:  sVcpb.ServiceType_REMITTANCE.String(),
			Status:   sVcpb.ServiceRequestStatus_PARTNERDRAFT,
		})
		http.Redirect(w, r, dsaServicesPath, http.StatusSeeOther)
		return
	}
	if form.SaveContinue == "SaveContinue" {
		s.pf.SetStatusUploadSvcRequest(ctx, &sVcpb.SetStatusUploadSvcRequestRequest{
			OrgID:    oid,
			Partners: notAcceptedPtnr,
			SvcName:  sVcpb.ServiceType_REMITTANCE.String(),
			Status:   sVcpb.ServiceRequestStatus_REQDOCDRAFT,
		})
		http.Redirect(w, r, dsaReqDocsPath, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, dsaPrtSelPath, http.StatusSeeOther)
}

func getTransactionTypes(bip *ppb.OrgProfile) (string, error) {
	return bip.GetTransactionTypes(), nil
}

func getPartnerShortName(ptnrListMaps map[string]string, ptnr []string) (res []string) {
	hAll, _ := cmsmw.InArray(AllPartners, ptnr)
	if hAll {
		for _, v := range ptnrListMaps {
			res = append(res, v)
		}
		return
	}
	for _, v := range ptnr {
		res = append(res, ptnrListMaps[v])
	}
	return
}

func getPartnerFullName(ptnrListMaps map[string]string, ptnr string) (res string) {
	if ptnr == AllPartners {
		return AllPartners
	}
	for k, v := range ptnrListMaps {
		if v == ptnr {
			res = k
		}
	}
	return
}

func removeAcceptedPartnerList(p []PartnerListApplicant, SType []string) []PartnerListApplicant {
	for _, v := range SType {
		for i, pv := range p {
			if v == pv.Stype {
				p = append(p[:i], p[i+1:]...)
			}
		}
	}
	return p
}

func (s *Server) getPartnerListMaps(ctx context.Context, serviceType string) map[string]string {
	ptnrs := map[string]string{}
	res, err := s.pf.GetPartnerList(ctx, &spbl.GetPartnerListRequest{
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: serviceType,
	})
	if err != nil || res == nil {
		return ptnrs
	}

	for _, v := range res.GetPartnerList() {
		ptnrs[v.Name] = v.Stype
	}
	return ptnrs
}

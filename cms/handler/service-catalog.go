package handler

import (
	"errors"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	ppf "brank.as/petnet/gunk/dsa/v1/user"
	fpb "brank.as/petnet/gunk/dsa/v2/file"
	spb "brank.as/petnet/gunk/dsa/v2/partner"
	spbl "brank.as/petnet/gunk/dsa/v2/partnerlist"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	psvc "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	rbupb "brank.as/rbac/gunk/v1/user"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gorilla/csrf"
	"github.com/kenshaw/goji"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	svcCatTempDate struct {
		BusinessInfo               *ppb.BusinessInfo
		WUInfo                     *spb.WesternUnionPartner
		TFInfo                     *spb.TransfastPartner
		IRInfo                     *spb.IRemitPartner
		RIAInfo                    *spb.RiaPartner
		MBInfo                     *spb.MetroBankPartner
		RMInfo                     *spb.RemitlyPartner
		BPIInfo                    *spb.BPIPartner
		USSCInfo                   *spb.USSCPartner
		JPRInfo                    *spb.JapanRemitPartner
		ICInfo                     *spb.InstantCashPartner
		UNTInfo                    *spb.UnitellerPartner
		CEBInfo                    *spb.CebuanaPartner
		WISEInfo                   *spb.TransferWisePartner
		CEBIInfo                   *spb.CebuanaIntlPartner
		AYAInfo                    *spb.AyannahPartner
		IEInfo                     *spb.IntelExpressPartner
		CSRFField                  template.HTML
		WUID                       string
		PartnerName                string
		Param1                     string
		Param2                     string
		OrgID                      string
		Status                     string
		MessageDetails             string
		StartDate                  string
		EndDate                    string
		UpdatedBy                  string
		AllDocsSubmitted           bool
		Errors                     map[string]string
		PresetPermission           map[string]map[string]bool
		ServiceRequest             bool
		PartnerListApplicants      []PartnerListApplicant
		User                       *User
		UserInfo                   *User
		DsaCode                    string
		TerminalIdOtc              string
		TerminalIdDigital          string
		TransactionTypesForDSACode TransactionTypesForDSACode
	}
	PartnerDelete struct {
		ID    string
		OrgID string
	}
)

func (s *Server) getPartnerCatalog(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()

	template := s.templates.Lookup("service-catalog.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	oid := goji.Param(r, "id")
	if oid == "" {
		log.Error("missing org id in url param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	queryParams := r.URL.Query()
	msgType, _ := url.PathUnescape(queryParams.Get("msgType"))
	messageDetails := ""
	if msgType == "success" {
		messageDetails = "Successfully Deleted."
	}

	pf, err := s.pf.GetProfile(ctx, &ppb.GetProfileRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}
	businessInfo := pf.GetProfile().GetBusinessInfo()
	status := pf.GetProfile().GetStatus().String()
	if pf.GetProfile().GetStatus() == ppb.Status_UnknownStatus {
		status = "Incomplete"
	}
	if pf.GetProfile().GetStatus() == ppb.Status_PendingDocuments {
		status = "Pending Documents"
	}

	svc, err := s.pf.GetPartners(ctx, &spb.GetPartnersRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting services")
	}
	if svc == nil {
		svc = &spb.GetPartnersResponse{}
	}
	fs, err := s.pf.ListFiles(r.Context(), &fpb.ListFilesRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting files")
	}
	etd := s.getEnforceTemplateData(ctx)
	UpdatedBy := ""
	if svc != nil {
		cu, err := s.rbac.GetUser(ctx, &rbupb.GetUserRequest{ID: svc.GetPartners().GetUpdatedBy()})
		if err == nil && cu != nil {
			UpdatedBy = cu.User.FirstName + " " + cu.User.LastName
		}
	}
	newReq := &spbl.GetPartnerListRequest{
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: psvc.ServiceType_REMITTANCE.String(),
	}
	svcl, err := s.pf.GetPartnerList(r.Context(), newReq)
	var partnerListApplicantList []PartnerListApplicant
	if err != nil {
		log.Error("failed to Get Partner List")
	} else {
		for _, sv := range svcl.GetPartnerList() {
			partnerListApplicantList = append(partnerListApplicantList, PartnerListApplicant{
				Stype: sv.GetStype(),
				Name:  sv.GetName(),
			})
		}
	}
	transactionTypesForDSACode := checkTrnTypeForDSACode(pf.Profile.GetTransactionTypes())
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	data := svcCatTempDate{
		BusinessInfo:               businessInfo,
		WUInfo:                     svc.GetPartners().GetWesternUnionPartner(),
		TFInfo:                     svc.GetPartners().GetTransfastPartner(),
		IRInfo:                     svc.GetPartners().GetIRemitPartner(),
		MBInfo:                     svc.GetPartners().GetMetroBankPartner(),
		BPIInfo:                    svc.GetPartners().GetBPIPartner(),
		RIAInfo:                    svc.GetPartners().GetRiaPartner(),
		RMInfo:                     svc.GetPartners().GetRemitlyPartner(),
		USSCInfo:                   svc.GetPartners().GetUSSCPartner(),
		JPRInfo:                    svc.GetPartners().GetJapanRemitPartner(),
		ICInfo:                     svc.GetPartners().GetInstantCashPartner(),
		UNTInfo:                    svc.GetPartners().GetUnitellerPartner(),
		CEBInfo:                    svc.GetPartners().GetCebuanaPartner(),
		WISEInfo:                   svc.GetPartners().GetTransferWisePartner(),
		CEBIInfo:                   svc.GetPartners().GetCebuanaIntlPartner(),
		AYAInfo:                    svc.GetPartners().GetAyannahPartner(),
		IEInfo:                     svc.GetPartners().GetIntelExpressPartner(),
		CSRFField:                  csrf.TemplateField(r),
		OrgID:                      oid,
		Status:                     status,
		MessageDetails:             messageDetails,
		UpdatedBy:                  UpdatedBy,
		AllDocsSubmitted:           allDocsSubmitted(fs.GetFileUploads()),
		PresetPermission:           etd.PresetPermission,
		ServiceRequest:             etd.ServiceRequests,
		PartnerListApplicants:      partnerListApplicantList,
		User:                       &usrInfo.UserInfo,
		UserInfo:                   &usrInfo.UserInfo,
		DsaCode:                    pf.Profile.DsaCode,
		TerminalIdOtc:              pf.Profile.TerminalIdOtc,
		TerminalIdDigital:          pf.Profile.TerminalIdDigital,
		TransactionTypesForDSACode: transactionTypesForDSACode,
	}

	data.User.ProfileImage = usrInfo.ProfileImage

	if err := template.Execute(w, data); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) deletePartnerCatalog(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	var form PartnerDelete
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	formErr := validation.Errors{}
	if err := validation.ValidateStruct(&form,
		validation.Field(&form.ID, validation.Required),
		validation.Field(&form.OrgID, validation.Required),
	); err != nil {
		if err, ok := (err).(validation.Errors); ok {
			if err["ID"] != nil {
				formErr["ID"] = errors.New("Partner ID is required")
			}
			if err["OrgID"] != nil {
				formErr["OrgID"] = errors.New("Org ID is required")
			}
		}
		logging.WithError(err, log).Error("invalid request")
	}

	if len(formErr) > 0 {
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	_, err := s.pf.DeletePartner(ctx, &spb.DeletePartnerRequest{
		ID: form.ID,
	})
	if err != nil {
		logging.WithError(err, log).Error("deleting services")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, strings.Replace(serviceCatalogPath, ":id", form.OrgID, 1)+"?msgType=success", http.StatusSeeOther)
}

func (s *Server) postPartnerCatalog(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	oid := goji.Param(r, "id")
	if oid == "" {
		log.Error("missing org id in url param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	var form svcCatTempDate
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	formErr := map[string]string{}
	if form.PartnerName == "WU" {
		if err := validation.ValidateStruct(&form,
			validation.Field(&form.Param1, is.Digit, validation.Required),
			validation.Field(&form.Param2, is.Digit, validation.Required),
		); err != nil {
			if err, ok := (err).(validation.Errors); ok {
				if err["Param1"] != nil {
					formErr["Param1"] = "Terminal ID should be digit"
				}
				if err["Param2"] != nil {
					formErr["Param2"] = "COY ID should be digit"
				}
			}
			logging.WithError(err, log).Error("invalid request")
		}
	}

	form.Errors = formErr
	form.CSRFField = csrf.TemplateField(r)
	if len(formErr) > 0 {
		template := s.templates.Lookup("service-catalog.html")
		if template == nil {
			log.Error("unable to load template")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}

		oid := goji.Param(r, "id")
		if oid == "" {
			log.Error("missing org id in url param")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}

		pf, err := s.pf.GetProfile(ctx, &ppb.GetProfileRequest{OrgID: oid})
		if err != nil {
			logging.WithError(err, log).Info("getting profile")
		}
		businessInfo := pf.GetProfile().GetBusinessInfo()
		status := pf.GetProfile().GetStatus().String()
		if pf.GetProfile().GetStatus() == ppb.Status_UnknownStatus {
			status = "Incomplete"
		}
		if pf.GetProfile().GetStatus() == ppb.Status_PendingDocuments {
			status = "Pending Documents"
		}

		svc, err := s.pf.GetPartners(ctx, &spb.GetPartnersRequest{OrgID: oid})
		if err != nil {
			logging.WithError(err, log).Info("getting services")
		}
		wu := svc.GetPartners().GetWesternUnionPartner()
		wuInfo := &spb.WesternUnionPartner{
			Coy:        wu.GetCoy(),
			TerminalID: wu.GetTerminalID(),
			Created:    wu.GetCreated(),
			Updated:    wu.GetUpdated(),
		}

		fs, err := s.pf.ListFiles(r.Context(), &fpb.ListFilesRequest{OrgID: oid})
		if err != nil {
			logging.WithError(err, log).Info("getting files")
		}
		etd := s.getEnforceTemplateData(ctx)
		UpdatedBy := ""
		cu, err := s.rbac.GetUser(ctx, &rbupb.GetUserRequest{ID: svc.Partners.UpdatedBy})
		if err == nil {
			UpdatedBy = cu.User.FirstName + " " + cu.User.LastName
		}
		newReq := &spbl.GetPartnerListRequest{
			Status:      spb.PartnerStatusType_ENABLED.String(),
			ServiceName: psvc.ServiceType_REMITTANCE.String(),
		}
		svcl, err := s.pf.GetPartnerList(r.Context(), newReq)
		var partnerListApplicantList []PartnerListApplicant
		if err != nil {
			log.Error("failed to Get Partner List")
		} else {
			for _, sv := range svcl.GetPartnerList() {
				partnerListApplicantList = append(partnerListApplicantList, PartnerListApplicant{
					Stype: sv.GetStype(),
					Name:  sv.GetName(),
				})
			}
		}
		usrInfo := s.GetUserInfoFromCookie(w, r, false)

		data := svcCatTempDate{
			BusinessInfo:          businessInfo,
			WUInfo:                wuInfo,
			CSRFField:             csrf.TemplateField(r),
			WUID:                  form.WUID,
			Param1:                form.Param1,
			Param2:                form.Param2,
			OrgID:                 oid,
			Status:                status,
			StartDate:             form.StartDate,
			EndDate:               form.EndDate,
			UpdatedBy:             UpdatedBy,
			AllDocsSubmitted:      allDocsSubmitted(fs.GetFileUploads()),
			Errors:                form.Errors,
			PresetPermission:      etd.PresetPermission,
			ServiceRequest:        etd.ServiceRequests,
			PartnerListApplicants: partnerListApplicantList,
			UserInfo:              &usrInfo.UserInfo,
			User:                  &usrInfo.UserInfo,
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
		return
	}
	layout := "2006-01-02"
	t1, err := time.Parse(layout, form.StartDate)
	if err != nil {
		logging.WithError(err, log).Error("decoding form")
		return
	}
	t2, err := time.Parse(layout, form.EndDate)
	if err != nil {
		logging.WithError(err, log).Error("decoding form")
		return
	}

	convertedStartDate := timestamppb.New(t1)
	convertedEndDate := timestamppb.New(t2)
	tc := convertedStartDate.AsTime().Before(convertedEndDate.AsTime())
	if !tc {
		logging.WithError(err, log).Error("End date must be greater than start date")
		http.Redirect(w, r, strings.Replace(serviceCatalogPath, ":id", oid, 1), http.StatusSeeOther)
		return
	}
	sn := form.PartnerName
	var isCreated bool
	serviceID := "0"
	res, err := s.pf.GetPartner(ctx, &spb.GetPartnersRequest{OrgID: oid, Type: sn})
	if err != nil {
		logging.WithError(err, log).Error("getting services")
		isCreated = true
	}
	if res.GetPartner().GetID() == "" {
		isCreated = true
	}
	if res.GetPartner().GetID() != "" {
		serviceID = res.Partner.ID
	}
	srv := &spb.Partners{
		OrgID:     oid,
		UpdatedBy: mw.GetUserID(ctx),
	}
	switch sn {
	case "WU":
		srv.WesternUnionPartner = &spb.WesternUnionPartner{
			ID:         serviceID,
			Coy:        form.Param2,
			TerminalID: form.Param1,
			Status:     spb.PartnerStatusType_ENABLED,
			Created:    timestamppb.Now(),
			Updated:    timestamppb.Now(),
			StartDate:  convertedStartDate,
			EndDate:    convertedEndDate,
		}
	case "TF":
		srv.TransfastPartner = &spb.TransfastPartner{
			ID:        serviceID,
			Param1:    form.Param1,
			Param2:    form.Param2,
			Status:    spb.PartnerStatusType_ENABLED,
			Created:   timestamppb.Now(),
			Updated:   timestamppb.Now(),
			StartDate: convertedStartDate,
			EndDate:   convertedEndDate,
		}
	case "JPR":
		srv.JapanRemitPartner = &spb.JapanRemitPartner{
			ID:        serviceID,
			Param1:    form.Param1,
			Param2:    form.Param2,
			Status:    spb.PartnerStatusType_ENABLED,
			Created:   timestamppb.Now(),
			Updated:   timestamppb.Now(),
			StartDate: convertedStartDate,
			EndDate:   convertedEndDate,
		}
	case "IR":
		srv.IRemitPartner = &spb.IRemitPartner{
			ID:        serviceID,
			Param1:    form.Param1,
			Param2:    form.Param2,
			Status:    spb.PartnerStatusType_ENABLED,
			Created:   timestamppb.Now(),
			Updated:   timestamppb.Now(),
			StartDate: convertedStartDate,
			EndDate:   convertedEndDate,
		}
	case "RIA":
		srv.RiaPartner = &spb.RiaPartner{
			ID:        serviceID,
			Param1:    form.Param1,
			Param2:    form.Param2,
			Status:    spb.PartnerStatusType_ENABLED,
			Created:   timestamppb.Now(),
			Updated:   timestamppb.Now(),
			StartDate: convertedStartDate,
			EndDate:   convertedEndDate,
		}
	case "MB":
		srv.MetroBankPartner = &spb.MetroBankPartner{
			ID:        serviceID,
			Param1:    form.Param1,
			Param2:    form.Param2,
			Status:    spb.PartnerStatusType_ENABLED,
			Created:   timestamppb.Now(),
			Updated:   timestamppb.Now(),
			StartDate: convertedStartDate,
			EndDate:   convertedEndDate,
		}
	case "RM":
		srv.RemitlyPartner = &spb.RemitlyPartner{
			ID:        serviceID,
			Param1:    form.Param1,
			Param2:    form.Param2,
			Status:    spb.PartnerStatusType_ENABLED,
			Created:   timestamppb.Now(),
			Updated:   timestamppb.Now(),
			StartDate: convertedStartDate,
			EndDate:   convertedEndDate,
		}
	case "BPI":
		srv.BPIPartner = &spb.BPIPartner{
			ID:        serviceID,
			Param1:    form.Param1,
			Param2:    form.Param2,
			Status:    spb.PartnerStatusType_ENABLED,
			Created:   timestamppb.Now(),
			Updated:   timestamppb.Now(),
			StartDate: convertedStartDate,
			EndDate:   convertedEndDate,
		}
	case "USSC":
		srv.USSCPartner = &spb.USSCPartner{
			ID:        serviceID,
			Param1:    form.Param1,
			Param2:    form.Param2,
			Status:    spb.PartnerStatusType_ENABLED,
			Created:   timestamppb.Now(),
			Updated:   timestamppb.Now(),
			StartDate: convertedStartDate,
			EndDate:   convertedEndDate,
		}
	case "IC":
		srv.InstantCashPartner = &spb.InstantCashPartner{
			ID:        serviceID,
			Param1:    form.Param1,
			Param2:    form.Param2,
			Status:    spb.PartnerStatusType_ENABLED,
			Created:   timestamppb.Now(),
			Updated:   timestamppb.Now(),
			StartDate: convertedStartDate,
			EndDate:   convertedEndDate,
		}
	case "UNT":
		srv.UnitellerPartner = &spb.UnitellerPartner{
			ID:        serviceID,
			Param1:    form.Param1,
			Param2:    form.Param2,
			Status:    spb.PartnerStatusType_ENABLED,
			Created:   timestamppb.Now(),
			Updated:   timestamppb.Now(),
			StartDate: convertedStartDate,
			EndDate:   convertedEndDate,
		}
	case "CEB":
		srv.CebuanaPartner = &spb.CebuanaPartner{
			ID:        serviceID,
			Param1:    form.Param1,
			Param2:    form.Param2,
			Status:    spb.PartnerStatusType_ENABLED,
			Created:   timestamppb.Now(),
			Updated:   timestamppb.Now(),
			StartDate: convertedStartDate,
			EndDate:   convertedEndDate,
		}
	case "WISE":
		srv.TransferWisePartner = &spb.TransferWisePartner{
			ID:        serviceID,
			Param1:    form.Param1,
			Param2:    form.Param2,
			Status:    spb.PartnerStatusType_ENABLED,
			Created:   timestamppb.Now(),
			Updated:   timestamppb.Now(),
			StartDate: convertedStartDate,
			EndDate:   convertedEndDate,
		}
	case "CEBI":
		srv.CebuanaIntlPartner = &spb.CebuanaIntlPartner{
			ID:        serviceID,
			Param1:    form.Param1,
			Param2:    form.Param2,
			Status:    spb.PartnerStatusType_ENABLED,
			Created:   timestamppb.Now(),
			Updated:   timestamppb.Now(),
			StartDate: convertedStartDate,
			EndDate:   convertedEndDate,
		}
	case "AYA":
		srv.AyannahPartner = &spb.AyannahPartner{
			ID:        serviceID,
			Param1:    form.Param1,
			Param2:    form.Param2,
			Status:    spb.PartnerStatusType_ENABLED,
			Created:   timestamppb.Now(),
			Updated:   timestamppb.Now(),
			StartDate: convertedStartDate,
			EndDate:   convertedEndDate,
		}
	case "IE":
		srv.IntelExpressPartner = &spb.IntelExpressPartner{
			ID:        serviceID,
			Param1:    form.Param1,
			Param2:    form.Param2,
			Status:    spb.PartnerStatusType_ENABLED,
			Created:   timestamppb.Now(),
			Updated:   timestamppb.Now(),
			StartDate: convertedStartDate,
			EndDate:   convertedEndDate,
		}
	}
	if isCreated {
		if _, err := s.pf.CreatePartners(ctx, &spb.CreatePartnersRequest{
			Partners: srv,
		}); err != nil {
			logging.WithError(err, log).Error("updating services")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	} else {
		if _, err := s.pf.UpdatePartners(ctx, &spb.UpdatePartnersRequest{
			Partners: srv,
		}); err != nil {
			logging.WithError(err, log).Error("updating services")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}
	http.Redirect(w, r, strings.Replace(serviceCatalogPath, ":id", oid, 1), http.StatusSeeOther)
}

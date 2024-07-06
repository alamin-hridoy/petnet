package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	fpb "brank.as/petnet/gunk/dsa/v2/file"
	ptnrcom "brank.as/petnet/gunk/dsa/v2/partnercommission"
	spbl "brank.as/petnet/gunk/dsa/v2/partnerlist"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	revcom "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	svc "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	"github.com/gorilla/csrf"
	"github.com/kenshaw/goji"
)

type RevenueSharingTemplateData struct {
	UserInfo                   *User
	User                       *User
	PresetPermission           map[string]map[string]bool
	ServiceRequest             bool
	PartnerList                RevenueSharingList
	CSRFField                  template.HTML
	OrgID                      string
	Status                     string
	AllDocsSubmitted           bool
	CSRFFieldValue             string
	BusinessInfo               *ppb.BusinessInfo
	AccountInfo                *ppb.AccountInfo
	OrgTransactionTypes        []string
	DsaCode                    string
	TerminalIdOtc              string
	TerminalIdDigital          string
	TransactionTypesForDSACode TransactionTypesForDSACode
	SvcReqStatus               string
}

type PartnerRevenueSharingTier struct {
	ID               string
	RevenueSharingID string
	MinValue         string
	MaxValue         string
	Amount           string
}

type PartnerRevenueSharing struct {
	ID              string
	Partner         string
	BoundType       string
	RemitType       string
	TransactionType string
	TierType        string
	Amount          string
	TierList        []PartnerRevenueSharingTier
}

type RevenueSharingList struct {
	ID               string
	Stype            string
	Name             string
	Status           string
	Reqstatus        string
	TransactionTypes []string
	InBound          map[string]PartnerRevenueSharing
	OutBound         map[string]PartnerRevenueSharing
}

func (s *Server) FormatPartnerRevenueSharingTier(ctx context.Context, id string) []PartnerRevenueSharingTier {
	res := []PartnerRevenueSharingTier{}
	revComList, err := s.pf.GetRevenueSharingTierList(ctx, &revcom.GetRevenueSharingTierListRequest{
		RevenueSharingID: id,
	})
	if err != nil {
		return res
	}
	for _, v := range revComList.GetResults() {
		res = append(res, PartnerRevenueSharingTier{
			ID:               v.GetID(),
			RevenueSharingID: v.GetRevenueSharingID(),
			MinValue:         v.GetMinValue(),
			MaxValue:         v.GetMaxValue(),
			Amount:           v.GetAmount(),
		})
	}
	return res
}

func (s *Server) FormatPartnerRevenueSharing(ctx context.Context, revMgts []*revcom.RevenueSharing, tranType revcom.TransactionType) PartnerRevenueSharing {
	list := PartnerRevenueSharing{}
	for _, v := range revMgts {
		if v.GetTransactionType() == tranType {
			list = PartnerRevenueSharing{
				ID:              v.ID,
				Partner:         v.GetPartner(),
				BoundType:       v.GetBoundType().String(),
				RemitType:       v.GetRemitType().String(),
				TransactionType: v.GetTransactionType().String(),
				TierType:        v.GetTierType().String(),
				Amount:          v.GetAmount(),
				TierList:        s.FormatPartnerRevenueSharingTier(ctx, v.GetID()),
			}
		}
	}
	return list
}

func (s *Server) postRevenueSharingMgt(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	oid := goji.Param(r, "id")
	if oid == "" {
		log.Error("org id shouldn't empty")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	uid := mw.GetUserID(ctx)
	if uid == "" {
		log.Error("user id shouldn't empty")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var partnerList []*spbl.PartnerList
	partnerList = []*spbl.PartnerList{
		{
			Stype:       "all",
			Name:        "Remittance Service",
			Status:      "ENABLED",
			ServiceName: ptnrcom.RemitType_REMITTANCE.String(),
			UpdatedBy:   mw.GetUserID(ctx),
		},
	}
	remitTypes := []map[string]string{{ptnrcom.RemitType_REMITTANCE.String(): "rmt"}}
	boundTypes := []map[string]string{{ptnrcom.BoundType_OUTBOUND.String(): "otb"}, {ptnrcom.BoundType_INBOUND.String(): "inb"}}
	if len(remitTypes) == 0 || len(boundTypes) == 0 {
		log.Error("remitTypes, boundTypes shouldn't empty")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	for _, rtVMaps := range remitTypes {
		if err := s.saveRevenueSharingRemitTypesValue(r, saveRemitTypesValueReq{
			PartnerList: partnerList,
			RtVMaps:     rtVMaps,
			BoundTypes:  boundTypes,
			Uid:         uid,
			OrgID:       oid,
		}); err != nil {
			log.Error(err)
			continue
		}
	}

	// sync drp values to perahub in background
	// TODO(vitthal): Make it in background. Getting authentication error for background context
	s.remcoCommSvc.SyncDSACommissionConfigForRemittance(ctx, oid)

	http.Redirect(w, r, "/dashboard/revenue-sharing-mgt/"+oid, http.StatusSeeOther)
}

func (s *Server) getRevenueSharingMgt(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	template := s.templates.Lookup("revenue-sharing-mgt.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	oid := goji.Param(r, "id")
	if oid == "" {
		log.Error("org id shouldn't empty")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	uid := mw.GetUserID(ctx)
	pf, err := s.pf.GetProfile(r.Context(), &ppb.GetProfileRequest{OrgID: string(oid)})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}
	transactionType := pf.GetProfile().TransactionTypes
	transactionTypes := []string{}
	if transactionType != "" {
		transactionTypes = strings.Split(transactionType, ",")
		if len(transactionTypes) == 0 {
			transactionTypes = []string{revcom.TransactionType_DIGITAL.String(), revcom.TransactionType_OTC.String()}
		}
	}
	loadProfile := pf.GetProfile()
	status := loadProfile.GetStatus().String()
	if loadProfile.GetStatus() == ppb.Status_UnknownStatus {
		status = "Incomplete"
	}
	if loadProfile.GetStatus() == ppb.Status_PendingDocuments {
		status = "Pending Documents"
	}
	fs, err := s.pf.ListFiles(r.Context(), &fpb.ListFilesRequest{OrgID: string(oid)})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}
	svcs, err := s.pf.ListServiceRequest(ctx, &svc.ListServiceRequestRequest{
		OrgIDs:   []string{oid},
		Types:    []svc.ServiceType{svc.ServiceType_REMITTANCE},
		Statuses: []svc.ServiceRequestStatus{svc.ServiceRequestStatus_ACCEPTED},
	})
	if err != nil {
		log.Error("failed to get list service request")
	}
	var sts string
	if len(svcs.GetServiceRequst()) > 0 {
		for _, lv := range svcs.GetServiceRequst() {
			if lv.GetStatus().String() == svc.ServiceRequestStatus_ACCEPTED.String() {
				sts = svc.ServiceRequestStatus_ACCEPTED.String()
				break
			}
		}
	}

	ptnrOrgLists := []string{}
	if svcs != nil {
		for _, v := range svcs.ServiceRequst {
			ptnrOrgLists = append(ptnrOrgLists, v.Partner)
		}
	}

	gptt, _ := s.pf.GetPartnerTransactionType(ctx, &revcom.GetPartnerTransactionTypeRequest{
		Partners: strings.Join(ptnrOrgLists, ","),
	})
	ptnrTransactions := make(map[string][]string)
	if gptt != nil {
		for _, v := range gptt.PartnerDetails {
			ptnrTransactions[v.Stype] = v.TransactionTypes
		}
	}
	boundTypes := []revcom.BoundType{revcom.BoundType_INBOUND, revcom.BoundType_OUTBOUND}

	partnerList := RevenueSharingList{
		Stype:            "all",
		Name:             "Remittance Service",
		Status:           "ENABLED",
		TransactionTypes: []string{"DIGITAL", "OTC"},
	}
	partnerList.InBound = make(map[string]PartnerRevenueSharing)
	partnerList.OutBound = make(map[string]PartnerRevenueSharing)
	for _, bT := range boundTypes {
		ptnrList, err := s.pf.GetRevenueSharingList(ctx, &revcom.GetRevenueSharingListRequest{
			OrgID:     oid,
			UserID:    uid,
			BoundType: bT,
			RemitType: revcom.RemitType_REMITTANCE,
			Partner:   "all",
		})
		if err != nil {
			log.Error(err)
			continue
		}
		switch bT {
		case revcom.BoundType_INBOUND:
			partnerList.InBound[revcom.TransactionType_DIGITAL.String()] = s.FormatPartnerRevenueSharing(ctx, ptnrList.GetResults(), revcom.TransactionType_DIGITAL)
			partnerList.InBound[revcom.TransactionType_OTC.String()] = s.FormatPartnerRevenueSharing(ctx, ptnrList.GetResults(), revcom.TransactionType_OTC)
		case revcom.BoundType_OUTBOUND:
			partnerList.OutBound[revcom.TransactionType_DIGITAL.String()] = s.FormatPartnerRevenueSharing(ctx, ptnrList.GetResults(), revcom.TransactionType_DIGITAL)
			partnerList.OutBound[revcom.TransactionType_OTC.String()] = s.FormatPartnerRevenueSharing(ctx, ptnrList.GetResults(), revcom.TransactionType_OTC)
		}
	}

	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	transactionTypesForDSACode := checkTrnTypeForDSACode(loadProfile.GetTransactionTypes())
	etd := s.getEnforceTemplateData(r.Context())
	templateData := RevenueSharingTemplateData{
		OrgID:                      oid,
		CSRFField:                  csrf.TemplateField(r),
		AllDocsSubmitted:           allDocsSubmitted(fs.GetFileUploads()),
		UserInfo:                   &usrInfo.UserInfo,
		User:                       &usrInfo.UserInfo,
		PresetPermission:           etd.PresetPermission,
		ServiceRequest:             etd.ServiceRequests,
		PartnerList:                partnerList,
		CSRFFieldValue:             csrf.Token(r),
		Status:                     status,
		BusinessInfo:               loadProfile.GetBusinessInfo(),
		AccountInfo:                loadProfile.GetAccountInfo(),
		OrgTransactionTypes:        transactionTypes,
		DsaCode:                    loadProfile.DsaCode,
		TerminalIdOtc:              loadProfile.TerminalIdOtc,
		TerminalIdDigital:          loadProfile.TerminalIdDigital,
		TransactionTypesForDSACode: transactionTypesForDSACode,
		SvcReqStatus:               sts,
	}
	templateData.UserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, templateData); err != nil {
		logging.WithError(err, log).Error("error with template execution")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) removeRevenueSharingMgt(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	tierID := goji.Param(r, "id")
	if tierID != "" {
		_, err := s.pf.DeleteRevenueSharingTierById(ctx, &revcom.DeleteRevenueSharingTierByIdRequest{
			ID: tierID,
		})
		if err != nil {
			log.Error("failed to remove revenue sharing tier")
		} else {
			// sync drp values to perahub in background
			oid := goji.Param(r, "oid")
			if oid != "" {
				s.remcoCommSvc.SyncDSACommissionConfigForRemittance(ctx, oid)
			}
		}
	}
	jsn, _ := json.Marshal([]string{})
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsn)
}

func (s *Server) saveRevenueSharingBoundTypesValue(r *http.Request, config saveBoundTypesValueReq) error {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	if len(config.BtVMaps) == 0 {
		log.Error("boundTypes shouldn't empty")
		return errors.New("boundTypes shouldn't empty")
	}
	for btK, btV := range config.BtVMaps {
		for _, ptV := range config.PartnerList {
			if ptV.Stype == "" || config.RtV == "" || btV == "" {
				log.Error("partner type, remit type, bound type shouldn't empty")
				return errors.New("partner type, remit type, bound type shouldn't empty")
			}
			pref := fmt.Sprintf("%s%s%s", config.RtV, ptV.Stype, btV)
			comConfig := CommissionMgtConfig{
				Pref:     pref,
				PtVStype: ptV.Stype,
				BtK:      btK,
				RtK:      config.RtK,
				Uid:      config.Uid,
				OrgID:    config.OrgID,
			}
			trDgt := getFormValue(r.Form, pref, "TransactionDigital")
			trOtc := getFormValue(r.Form, pref, "TransactionOTC")
			dcomConfig := comConfig
			ocomConfig := comConfig
			dcomConfig.TrDgt = trDgt
			ocomConfig.TrOtc = trOtc
			if err := s.saveDigitalFixedRevenueSharing(r, dcomConfig); err != nil {
				log.Error(err)
			}
			if err := s.saveOTCFixedRevenueSharing(r, ocomConfig); err != nil {
				log.Error(err)
			}
		}
	}
	return nil
}

func (s *Server) saveRevenueSharingRemitTypesValue(r *http.Request, config saveRemitTypesValueReq) error {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	if len(config.RtVMaps) == 0 {
		log.Error("remitTypes shouldn't empty")
		return errors.New("remitTypes shouldn't empty")
	}
	for rtK, rtV := range config.RtVMaps {
		if rtK == "" || rtV == "" {
			log.Error("remit type shouldn't empty")
			continue
		}
		for _, btVMaps := range config.BoundTypes {
			if err := s.saveRevenueSharingBoundTypesValue(r, saveBoundTypesValueReq{
				PartnerList: config.PartnerList,
				BtVMaps:     btVMaps,
				Uid:         config.Uid,
				RtK:         rtK,
				RtV:         rtV,
				OrgID:       config.OrgID,
			}); err != nil {
				log.Error("remitTypes shouldn't empty")
				continue
			}
		}
	}
	return nil
}

func (s *Server) saveOTCFixedRevenueSharing(r *http.Request, config CommissionMgtConfig) error {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	oid := config.OrgID
	if oid == "" {
		return errors.New("org id shouldn't empty")
	}
	if config.TrOtc == "" || config.TrOtc != "OTC" {
		if _, err := s.pf.DeleteRevenueSharing(ctx, &revcom.DeleteRevenueSharingRequest{
			TransactionType: revcom.TransactionType(revcom.TransactionType_OTC),
			OrgID:           oid,
			UserID:          config.Uid,
			RemitType:       revcom.RemitType(revcom.RemitType_value[config.RtK]),
			Partner:         config.PtVStype,
			BoundType:       revcom.BoundType(revcom.BoundType_value[config.BtK]),
		}); err != nil {
			return err
		}
		return nil
	}
	trTypeOtc := getFormValue(r.Form, config.Pref, "TransactionTypeOTC")
	transactionOTCID := getFormValue(r.Form, config.Pref, "TransactionOTCID")
	switch trTypeOtc {
	case "PERCENTAGE":
		perAmtOtc := getFormValue(r.Form, config.Pref, "FixedPercentageOTC")
		fmt.Printf("otc per req: %+v", &revcom.UpsertRevenueSharingRequest{
			ID:              transactionOTCID,
			OrgID:           config.OrgID,
			UserID:          config.Uid,
			Partner:         config.PtVStype,
			BoundType:       revcom.BoundType(revcom.BoundType_value[config.BtK]),
			RemitType:       revcom.RemitType(revcom.RemitType_value[config.RtK]),
			TransactionType: revcom.TransactionType(revcom.TransactionType_OTC),
			TierType:        revcom.TierType(revcom.TierType_value[revcom.TierType_PERCENTAGE.String()]),
			Amount:          perAmtOtc,
			CreatedBy:       config.Uid,
		})
		comres, err := s.pf.UpsertRevenueSharing(ctx, &revcom.UpsertRevenueSharingRequest{
			ID:              transactionOTCID,
			OrgID:           config.OrgID,
			UserID:          config.Uid,
			Partner:         config.PtVStype,
			BoundType:       revcom.BoundType(revcom.BoundType_value[config.BtK]),
			RemitType:       revcom.RemitType(revcom.RemitType_value[config.RtK]),
			TransactionType: revcom.TransactionType(revcom.TransactionType_OTC),
			TierType:        revcom.TierType(revcom.TierType_value[revcom.TierType_PERCENTAGE.String()]),
			Amount:          perAmtOtc,
			CreatedBy:       config.Uid,
		})
		if err != nil {
			log.Error("failed to create revenue sharing mgt")
			return errors.New("failed to create revenue sharing mgt")
		}
		if _, err := s.pf.DeleteRevenueSharingTier(ctx, &revcom.DeleteRevenueSharingTierRequest{
			RevenueSharingID: comres.ID,
		}); err != nil {
			log.Error("failed to delete revenue sharing mgt tier")
			return errors.New("failed to delete revenue sharing mgt tier")
		}
	case "TIERPERCENTAGE":
		lnTirPrAmnOtc := getFormValue(r.Form, config.Pref, "LenTieredPercentageOTC")
		lnTirPrAmnOtcInt, _ := strconv.Atoi(lnTirPrAmnOtc)
		tierPerRes, err := s.pf.UpsertRevenueSharing(ctx, &revcom.UpsertRevenueSharingRequest{
			ID:              transactionOTCID,
			OrgID:           config.OrgID,
			UserID:          config.Uid,
			Partner:         config.PtVStype,
			BoundType:       revcom.BoundType(revcom.BoundType_value[config.BtK]),
			RemitType:       revcom.RemitType(revcom.RemitType_value[config.RtK]),
			TransactionType: revcom.TransactionType(revcom.TransactionType_OTC),
			TierType:        revcom.TierType(revcom.TierType_value[revcom.TierType_TIERPERCENTAGE.String()]),
			Amount:          "",
			CreatedBy:       config.Uid,
		})
		if err != nil {
			log.Error("failed to create revenue sharing mgt")
			return errors.New("failed to create revenue sharing mgt")
		}
		if tierPerRes.ID != "" && lnTirPrAmnOtcInt > 0 {
			for i := 0; i < lnTirPrAmnOtcInt; i++ {
				mnTirPrAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("MinTieredPercentageOTC%d", i))
				mxTirPrAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("MaxTieredPercentageOTC%d", i))
				feTirPrAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("FeeTieredPercentageOTC%d", i))
				idTirPrAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("IDTieredPercentageOTC%d", i))
				_, err := s.pf.UpsertRevenueSharingTier(ctx, &revcom.UpsertRevenueSharingTierRequest{
					ID:               idTirPrAmtOtc,
					RevenueSharingID: tierPerRes.ID,
					MinValue:         mnTirPrAmtOtc,
					MaxValue:         mxTirPrAmtOtc,
					Amount:           feTirPrAmtOtc,
				})
				if err != nil {
					log.Error("failed to create revenue sharing mgt tier")
					return errors.New("failed to create revenue sharing mgt tier")
				}
			}
		}
	}
	return nil
}

func (s *Server) saveDigitalFixedRevenueSharing(r *http.Request, config CommissionMgtConfig) error {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	oid := config.OrgID
	if oid == "" {
		return errors.New("org id shouldn't empty")
	}
	if config.TrDgt == "" || config.TrDgt != "DIGITAL" {
		if _, err := s.pf.DeleteRevenueSharing(ctx, &revcom.DeleteRevenueSharingRequest{
			TransactionType: revcom.TransactionType(revcom.TransactionType_DIGITAL),
			OrgID:           oid,
			UserID:          config.Uid,
			RemitType:       revcom.RemitType(revcom.RemitType_value[config.RtK]),
			Partner:         config.PtVStype,
			BoundType:       revcom.BoundType(revcom.BoundType_value[config.BtK]),
		}); err != nil {
			return err
		}
		return nil
	}
	trTypeDgt := getFormValue(r.Form, config.Pref, "TransactionTypeDigital")
	transactionDIGITALID := getFormValue(r.Form, config.Pref, "TransactionDIGITALID")
	switch trTypeDgt {
	case "PERCENTAGE":
		perAmtDgt := getFormValue(r.Form, config.Pref, "FixedPercentageDigital")
		fmt.Printf("dgt per req: %+v", &revcom.UpsertRevenueSharingRequest{
			ID:              transactionDIGITALID,
			OrgID:           oid,
			UserID:          config.Uid,
			Partner:         config.PtVStype,
			BoundType:       revcom.BoundType(revcom.BoundType_value[config.BtK]),
			RemitType:       revcom.RemitType(revcom.RemitType_value[config.RtK]),
			TransactionType: revcom.TransactionType(revcom.TransactionType_DIGITAL),
			TierType:        revcom.TierType(revcom.TierType_value[revcom.TierType_PERCENTAGE.String()]),
			Amount:          perAmtDgt,
			CreatedBy:       config.Uid,
		})
		comIdres, err := s.pf.UpsertRevenueSharing(ctx, &revcom.UpsertRevenueSharingRequest{
			ID:              transactionDIGITALID,
			OrgID:           oid,
			UserID:          config.Uid,
			Partner:         config.PtVStype,
			BoundType:       revcom.BoundType(revcom.BoundType_value[config.BtK]),
			RemitType:       revcom.RemitType(revcom.RemitType_value[config.RtK]),
			TransactionType: revcom.TransactionType(revcom.TransactionType_DIGITAL),
			TierType:        revcom.TierType(revcom.TierType_value[revcom.TierType_PERCENTAGE.String()]),
			Amount:          perAmtDgt,
			CreatedBy:       config.Uid,
		})
		if err != nil {
			log.Error("failed to update revenue sharing mgt")
			return errors.New("failed to update revenue sharing mgt")
		}
		if _, err := s.pf.DeleteRevenueSharingTier(ctx, &revcom.DeleteRevenueSharingTierRequest{
			RevenueSharingID: comIdres.ID,
		}); err != nil {
			log.Error("failed to delete revenue sharing mgt tier")
			return errors.New("failed to delete revenue sharing mgt tier")
		}
	case "TIERPERCENTAGE":
		lnTirPrAmnDgt := getFormValue(r.Form, config.Pref, "LenTieredPercentageDigital")
		lnTirPrAmnDgtInt, _ := strconv.Atoi(lnTirPrAmnDgt)
		tierPerRes, err := s.pf.UpsertRevenueSharing(ctx, &revcom.UpsertRevenueSharingRequest{
			ID:              transactionDIGITALID,
			OrgID:           oid,
			UserID:          config.Uid,
			Partner:         config.PtVStype,
			BoundType:       revcom.BoundType(revcom.BoundType_value[config.BtK]),
			RemitType:       revcom.RemitType(revcom.RemitType_value[config.RtK]),
			TransactionType: revcom.TransactionType(revcom.TransactionType_DIGITAL),
			TierType:        revcom.TierType(revcom.TierType_value[revcom.TierType_TIERPERCENTAGE.String()]),
			Amount:          "",
			CreatedBy:       config.Uid,
		})
		if err != nil {
			log.Error("failed to update revenue sharing mgt")
			return errors.New("failed to update revenue sharing mgt")
		}
		if tierPerRes.ID != "" && lnTirPrAmnDgtInt > 0 {
			for i := 0; i < lnTirPrAmnDgtInt; i++ {
				mnTirPrAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("MinTieredPercentageDigital%d", i))
				mxTirPrAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("MaxTieredPercentageDigital%d", i))
				feTirPrAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("FeeTieredPercentageDigital%d", i))
				idTirPrAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("IDTieredPercentageDigital%d", i))
				_, err := s.pf.UpsertRevenueSharingTier(ctx, &revcom.UpsertRevenueSharingTierRequest{
					ID:               idTirPrAmtDgt,
					RevenueSharingID: tierPerRes.ID,
					MinValue:         mnTirPrAmtDgt,
					MaxValue:         mxTirPrAmtDgt,
					Amount:           feTirPrAmtDgt,
				})
				if err != nil {
					log.Error("failed to update revenue sharing mgt tier")
					return errors.New("failed to update revenue sharing mgt tier")
				}
			}
		}
	}
	return nil
}

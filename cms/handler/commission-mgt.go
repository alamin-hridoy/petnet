package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/gorilla/csrf"
	"github.com/kenshaw/goji"
	"google.golang.org/protobuf/types/known/timestamppb"

	spb "brank.as/petnet/gunk/dsa/v2/partner"
	ptnrcom "brank.as/petnet/gunk/dsa/v2/partnercommission"
	spbl "brank.as/petnet/gunk/dsa/v2/partnerlist"
	svc "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
)

type CommissionFeeTemplateData struct {
	UserInfo         *User
	PresetPermission map[string]map[string]bool
	ServiceRequest   bool
	PartnerList      []PartnerList
	CSRFField        template.HTML
	CSRFFieldValue   string
	ErrorMsg         string
}

type PartnerCommissionTier struct {
	ID                  string
	PartnerCommissionID string
	MinValue            string
	MaxValue            string
	Amount              string
}

type PartnerCommission struct {
	ID              string
	Partner         string
	BoundType       string
	RemitType       string
	TransactionType string
	TierType        string
	Amount          string
	TierList        []PartnerCommissionTier
	StartDate       time.Time
	EndDate         time.Time
}

type PartnerList struct {
	ID       string
	Stype    string
	Name     string
	Status   string
	InBound  map[string]PartnerCommission
	OutBound map[string]PartnerCommission
}

type CommissionMgtConfig struct {
	Pref                 string
	TrDgt                string
	TrOtc                string
	TrTypeDgt            string
	TransactionDIGITALID string
	ComId                string
	PtVStype             string
	BtK                  string
	RtK                  string
	Uid                  string
	OrgID                string
}

type saveRemitTypesValueReq struct {
	PartnerList []*spbl.PartnerList
	RtVMaps     map[string]string
	BoundTypes  []map[string]string
	Uid         string
	OrgID       string
}

type saveBoundTypesValueReq struct {
	PartnerList []*spbl.PartnerList
	BtVMaps     map[string]string
	Uid         string
	RtK         string
	RtV         string
	OrgID       string
}

type noCancel struct {
	ctx context.Context
}

func (c noCancel) Deadline() (time.Time, bool)       { return time.Time{}, false }
func (c noCancel) Done() <-chan struct{}             { return nil }
func (c noCancel) Err() error                        { return nil }
func (c noCancel) Value(key interface{}) interface{} { return c.ctx.Value(key) }

// WithoutCancel returns a context that is never canceled.
func ctxWithoutCancel(ctx context.Context) context.Context {
	return noCancel{ctx: ctx}
}

// doPostCommissionMgt is post action for commission management form
func (s *Server) doPostCommissionMgt(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	uid := mw.GetUserID(ctx)
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	svc, err := s.pf.GetPartnerList(ctx, &spbl.GetPartnerListRequest{
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: svc.ServiceType_REMITTANCE.String(),
	})
	if err != nil {
		log.Error("failed to get profile")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	var partnerList []*spbl.PartnerList
	ptnr := svc.GetPartnerList()
	if len(ptnr) > 0 {
		partnerList = ptnr
	}

	remitTypes := []map[string]string{{ptnrcom.RemitType_REMITTANCE.String(): "rmt"}}
	boundTypes := []map[string]string{{ptnrcom.BoundType_OUTBOUND.String(): "otb"}, {ptnrcom.BoundType_INBOUND.String(): "inb"}}
	if len(remitTypes) == 0 || len(partnerList) == 0 || len(boundTypes) == 0 {
		log.Error("remitTypes, partnerList, boundTypes shouldn't empty")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	for _, rtVMaps := range remitTypes {
		err = s.validateRemitTypesValue(r, saveRemitTypesValueReq{
			PartnerList: partnerList,
			RtVMaps:     rtVMaps,
			BoundTypes:  boundTypes,
			Uid:         uid,
		})
		if err != nil {
			http.Redirect(w, r, commissionFeePath+"?errmsg="+err.Error(), http.StatusSeeOther)
			return
		}
	}

	success := false
	for _, rtVMaps := range remitTypes {
		err = s.saveRemitTypesValue(r, saveRemitTypesValueReq{
			PartnerList: partnerList,
			RtVMaps:     rtVMaps,
			BoundTypes:  boundTypes,
			Uid:         uid,
		})

		if err != nil {
			log.Error(err)
			continue
		}

		success = true
	}

	if success {
		// sync drp values to perahub in background
		// TODO(vitthal): Make it in background. Getting authentication error for background context
		s.remcoCommSvc.SyncRemcoCommissionConfigForRemittance(ctx)
	}

	http.Redirect(w, r, commissionFeePath, http.StatusSeeOther)
}

// doGetCommissionMgt is get action for commission management form to display
func (s *Server) doGetCommissionMgt(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	template := s.templates.Lookup("commission-mgt.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	svc, err := s.pf.GetPartnerList(ctx, &spbl.GetPartnerListRequest{
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: svc.ServiceType_REMITTANCE.String(),
	})
	if err != nil {
		logging.WithError(err, log).Error("failed to get partner list")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	boundTypes := []ptnrcom.BoundType{ptnrcom.BoundType_INBOUND, ptnrcom.BoundType_OUTBOUND}
	var partnerList []PartnerList
	ptnr := svc.GetPartnerList()
	if len(ptnr) > 0 {
		for _, v := range ptnr {
			ptn := PartnerList{
				ID:     v.GetID(),
				Stype:  v.GetStype(),
				Name:   v.GetName(),
				Status: v.GetStatus(),
			}
			ptn.InBound = make(map[string]PartnerCommission)
			ptn.OutBound = make(map[string]PartnerCommission)
			for _, bT := range boundTypes {
				ptnrList, err := s.pf.GetPartnerCommissionsList(ctx, &ptnrcom.GetPartnerCommissionsListRequest{
					BoundType: bT,
					RemitType: ptnrcom.RemitType_REMITTANCE,
					Partner:   v.GetStype(),
				})
				if err != nil {
					continue
				}
				switch bT {
				case ptnrcom.BoundType_INBOUND:
					ptn.InBound[ptnrcom.TransactionType_DIGITAL.String()] = s.formatPartnerCommission(ctx, ptnrList.GetResults(), ptnrcom.TransactionType_DIGITAL)
					ptn.InBound[ptnrcom.TransactionType_OTC.String()] = s.formatPartnerCommission(ctx, ptnrList.GetResults(), ptnrcom.TransactionType_OTC)
				case ptnrcom.BoundType_OUTBOUND:
					ptn.OutBound[ptnrcom.TransactionType_DIGITAL.String()] = s.formatPartnerCommission(ctx, ptnrList.GetResults(), ptnrcom.TransactionType_DIGITAL)
					ptn.OutBound[ptnrcom.TransactionType_OTC.String()] = s.formatPartnerCommission(ctx, ptnrList.GetResults(), ptnrcom.TransactionType_OTC)
				}
			}
			partnerList = append(partnerList, ptn)
		}
	}
	queryParams := r.URL.Query()
	errmsg, err := url.PathUnescape(queryParams.Get("errmsg"))
	if err != nil {
		log.Error("unable to decode url type param")
	}

	usrInfo := s.GetUserInfoFromCookie(w, r, false)
	etd := s.getEnforceTemplateData(r.Context())
	templateData := CommissionFeeTemplateData{
		CSRFField:        csrf.TemplateField(r),
		UserInfo:         &usrInfo.UserInfo,
		PresetPermission: etd.PresetPermission,
		ServiceRequest:   etd.ServiceRequests,
		PartnerList:      partnerList,
		CSRFFieldValue:   csrf.Token(r),
		ErrorMsg:         errmsg,
	}
	templateData.UserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, templateData); err != nil {
		logging.WithError(err, log).Error("error with template execution")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

// doDeleteCommissionMgt is "delete" action for partner commission tier
func (s *Server) doDeleteCommissionMgt(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	tierID := goji.Param(r, "id")
	type errStruct struct {
		Message string
	}

	if tierID == "" {
		jsn, _ := json.Marshal(errStruct{Message: "tire id required"})

		w.WriteHeader(422)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsn)
		return
	}

	_, err := s.pf.DeletePartnerCommissionTierById(r.Context(), &ptnrcom.DeletePartnerCommissionTierByIdRequest{
		ID: tierID,
	})
	if err != nil {
		log.WithError(err).Error("failed to remove commission tier")
		jsn, _ := json.Marshal(errStruct{Message: "failed to remove commission tier"})
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsn)
		return
	}

	// sync drp values to perahub in background
	// TODO(vitthal): Make it in background. Getting authentication error for background context
	s.remcoCommSvc.SyncRemcoCommissionConfigForRemittance(r.Context())

	jsn, _ := json.Marshal([]string{})
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsn)
}

func (s *Server) getPartnerCommissionTiers(ctx context.Context, id string) []PartnerCommissionTier {
	ptnrComList, err := s.pf.GetPartnerCommissionsTierList(ctx, &ptnrcom.GetPartnerCommissionsTierListRequest{
		PartnerCommissionID: id,
	})
	if err != nil || ptnrComList == nil {
		return []PartnerCommissionTier{}
	}

	res := make([]PartnerCommissionTier, 0, len(ptnrComList.GetResults()))
	for _, v := range ptnrComList.GetResults() {
		res = append(res, PartnerCommissionTier{
			ID:                  v.GetID(),
			PartnerCommissionID: v.GetPartnerCommissionID(),
			MinValue:            v.GetMinValue(),
			MaxValue:            v.GetMaxValue(),
			Amount:              v.GetAmount(),
		})
	}
	return res
}

func (s *Server) formatPartnerCommission(ctx context.Context, ptnrCommissions []*ptnrcom.PartnerCommission, tranType ptnrcom.TransactionType) PartnerCommission {
	list := PartnerCommission{}
	for _, v := range ptnrCommissions {
		if v.GetTransactionType() == tranType {
			list = PartnerCommission{
				ID:              v.ID,
				Partner:         v.GetPartner(),
				BoundType:       v.GetBoundType().String(),
				RemitType:       v.GetRemitType().String(),
				TransactionType: v.GetTransactionType().String(),
				TierType:        v.GetTierType().String(),
				Amount:          v.GetAmount(),
				TierList:        s.getPartnerCommissionTiers(ctx, v.GetID()),
				StartDate:       v.GetStartDate().AsTime(),
				EndDate:         v.GetEndDate().AsTime(),
			}
		}
	}
	return list
}

func (s *Server) saveBoundTypesValue(r *http.Request, config saveBoundTypesValueReq) error {
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
			trDgt := getFormValue(r.Form, pref, "TransactionDigital")
			trOtc := getFormValue(r.Form, pref, "TransactionOTC")
			if err := s.saveDigitalFixedCommission(r, CommissionMgtConfig{
				Pref:     pref,
				TrDgt:    trDgt,
				PtVStype: ptV.Stype,
				BtK:      btK,
				RtK:      config.RtK,
				Uid:      config.Uid,
			}); err != nil {
				log.Error(err)
			}
			if err := s.saveOTCFixedCommission(r, CommissionMgtConfig{
				Pref:     pref,
				TrOtc:    trOtc,
				PtVStype: ptV.Stype,
				BtK:      btK,
				RtK:      config.RtK,
				Uid:      config.Uid,
			}); err != nil {
				log.Error(err)
			}
		}
	}
	return nil
}

func (s *Server) saveRemitTypesValue(r *http.Request, config saveRemitTypesValueReq) error {
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
			if err := s.saveBoundTypesValue(r, saveBoundTypesValueReq{
				PartnerList: config.PartnerList,
				BtVMaps:     btVMaps,
				Uid:         config.Uid,
				RtK:         rtK,
				RtV:         rtV,
			}); err != nil {
				log.Error("remitTypes shouldn't empty")
				return errors.New("remitTypes shouldn't empty")
			}
		}
	}
	return nil
}

func (s *Server) saveOTCFixedCommission(r *http.Request, config CommissionMgtConfig) error {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	if config.TrOtc == "" || config.TrOtc != "OTC" {
		_, err := s.pf.DeletePartnerCommission(ctx, &ptnrcom.DeletePartnerCommissionRequest{
			Partner:         config.PtVStype,
			BoundType:       ptnrcom.BoundType(ptnrcom.BoundType_value[config.BtK]),
			RemitType:       ptnrcom.RemitType(ptnrcom.RemitType_value[config.RtK]),
			TransactionType: ptnrcom.TransactionType(ptnrcom.TransactionType_OTC),
		})
		if err != nil {
			return err
		}
		return errors.New("transaction type should be OTC")
	}
	trTypeOtc := getFormValue(r.Form, config.Pref, "TransactionTypeOTC")
	transactionOTCID := getFormValue(r.Form, config.Pref, "TransactionOTCID")
	comId := transactionOTCID
	switch trTypeOtc {
	case "FIXED":
		fixAmtOtc := getFormValue(r.Form, config.Pref, "FixedAmountOTC")
		sDtOtcFx := getFormValue(r.Form, config.Pref, "StartDateOTCFixed")
		eDtOtcFx := getFormValue(r.Form, config.Pref, "EndDateOTCFixed")
		if transactionOTCID == "" {
			comres, err := s.pf.CreatePartnerCommission(ctx, &ptnrcom.CreatePartnerCommissionRequest{
				Partner:         config.PtVStype,
				BoundType:       ptnrcom.BoundType(ptnrcom.BoundType_value[config.BtK]),
				RemitType:       ptnrcom.RemitType(ptnrcom.RemitType_value[config.RtK]),
				TransactionType: ptnrcom.TransactionType(ptnrcom.TransactionType_value[config.TrOtc]),
				TierType:        ptnrcom.TierType(ptnrcom.TierType_value[trTypeOtc]),
				Amount:          fixAmtOtc,
				StartDate:       strToTimeConvert(sDtOtcFx),
				EndDate:         strToTimeConvert(eDtOtcFx),
				CreatedBy:       config.Uid,
			})
			if err != nil {
				log.Error("failed to create partner commission")
				return errors.New("failed to create partner commission")
			}
			comId = comres.GetID()
		} else {
			_, err := s.pf.UpdatePartnerCommission(ctx, &ptnrcom.UpdatePartnerCommissionRequest{
				ID:              transactionOTCID,
				Partner:         config.PtVStype,
				BoundType:       ptnrcom.BoundType(ptnrcom.BoundType_value[config.BtK]),
				RemitType:       ptnrcom.RemitType(ptnrcom.RemitType_value[config.RtK]),
				TransactionType: ptnrcom.TransactionType(ptnrcom.TransactionType_value[config.TrOtc]),
				TierType:        ptnrcom.TierType(ptnrcom.TierType_value[trTypeOtc]),
				Amount:          fixAmtOtc,
				StartDate:       strToTimeConvert(sDtOtcFx),
				EndDate:         strToTimeConvert(eDtOtcFx),
				CreatedBy:       config.Uid,
			})
			if err != nil {
				log.Error("failed to update partner commission")
				return errors.New("failed to update partner commission")
			}
		}
		if _, err := s.pf.DeletePartnerCommissionTier(ctx, &ptnrcom.DeletePartnerCommissionTierRequest{
			PartnerCommissionID: comId,
		}); err != nil {
			log.Error("failed to delete partner commission tier")
			return errors.New("failed to delete partner commission tier")
		}
	case "PERCENTAGE":
		perAmtOtc := getFormValue(r.Form, config.Pref, "FixedPercentageOTC")
		sDtOtcPer := getFormValue(r.Form, config.Pref, "StartDateOTCPercentage")
		eDtOtcPer := getFormValue(r.Form, config.Pref, "EndDateOTCPercentage")
		if transactionOTCID == "" {
			comres, err := s.pf.CreatePartnerCommission(ctx, &ptnrcom.CreatePartnerCommissionRequest{
				Partner:         config.PtVStype,
				BoundType:       ptnrcom.BoundType(ptnrcom.BoundType_value[config.BtK]),
				RemitType:       ptnrcom.RemitType(ptnrcom.RemitType_value[config.RtK]),
				TransactionType: ptnrcom.TransactionType(ptnrcom.TransactionType_value[config.TrOtc]),
				TierType:        ptnrcom.TierType(ptnrcom.TierType_value[trTypeOtc]),
				Amount:          perAmtOtc,
				StartDate:       strToTimeConvert(sDtOtcPer),
				EndDate:         strToTimeConvert(eDtOtcPer),
				CreatedBy:       config.Uid,
			})
			if err != nil {
				log.Error("failed to create partner commission")
				return errors.New("failed to create partner commission")
			}
			comId = comres.GetID()
		} else {
			_, err := s.pf.UpdatePartnerCommission(ctx, &ptnrcom.UpdatePartnerCommissionRequest{
				ID:              transactionOTCID,
				Partner:         config.PtVStype,
				BoundType:       ptnrcom.BoundType(ptnrcom.BoundType_value[config.BtK]),
				RemitType:       ptnrcom.RemitType(ptnrcom.RemitType_value[config.RtK]),
				TransactionType: ptnrcom.TransactionType(ptnrcom.TransactionType_value[config.TrOtc]),
				TierType:        ptnrcom.TierType(ptnrcom.TierType_value[trTypeOtc]),
				Amount:          perAmtOtc,
				StartDate:       strToTimeConvert(sDtOtcPer),
				EndDate:         strToTimeConvert(eDtOtcPer),
				CreatedBy:       config.Uid,
			})
			if err != nil {
				log.Error("failed to update partner commission")
				return errors.New("failed to update partner commission")
			}
		}
		if _, err := s.pf.DeletePartnerCommissionTier(ctx, &ptnrcom.DeletePartnerCommissionTierRequest{
			PartnerCommissionID: comId,
		}); err != nil {
			log.Error("failed to delete partner commission tier")
			return errors.New("failed to delete partner commission tier")
		}
	case "TIERAMOUNT":
		lnTirAmtOtc := getFormValue(r.Form, config.Pref, "LenTieredAmountOTC")
		lnTirAmtOtcInt, _ := strconv.Atoi(lnTirAmtOtc)
		tierAmtResID := ""
		if transactionOTCID == "" {
			tierAmtRes, err := s.pf.CreatePartnerCommission(ctx, &ptnrcom.CreatePartnerCommissionRequest{
				Partner:         config.PtVStype,
				BoundType:       ptnrcom.BoundType(ptnrcom.BoundType_value[config.BtK]),
				RemitType:       ptnrcom.RemitType(ptnrcom.RemitType_value[config.RtK]),
				TransactionType: ptnrcom.TransactionType(ptnrcom.TransactionType_value[config.TrOtc]),
				TierType:        ptnrcom.TierType(ptnrcom.TierType_value[trTypeOtc]),
				Amount:          "",
				CreatedBy:       config.Uid,
			})
			if err != nil {
				log.Error("failed to create partner commission")
				return errors.New("failed to create partner commission")
			}
			tierAmtResID = tierAmtRes.ID
		} else {
			tierAmtRes, err := s.pf.UpdatePartnerCommission(ctx, &ptnrcom.UpdatePartnerCommissionRequest{
				ID:              transactionOTCID,
				Partner:         config.PtVStype,
				BoundType:       ptnrcom.BoundType(ptnrcom.BoundType_value[config.BtK]),
				RemitType:       ptnrcom.RemitType(ptnrcom.RemitType_value[config.RtK]),
				TransactionType: ptnrcom.TransactionType(ptnrcom.TransactionType_value[config.TrOtc]),
				TierType:        ptnrcom.TierType(ptnrcom.TierType_value[trTypeOtc]),
				Amount:          "",
				CreatedBy:       config.Uid,
			})
			if err != nil {
				log.Error("failed to update partner commission")
				return errors.New("failed to update partner commission")
			}
			tierAmtResID = tierAmtRes.ID
		}
		if tierAmtResID != "" && lnTirAmtOtcInt > 0 {
			cpcttotc := []*ptnrcom.PartnerCommissionTier{}
			upcttotc := []*ptnrcom.PartnerCommissionTier{}
			for i := 0; i < lnTirAmtOtcInt; i++ {
				mnTirAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("MinTieredAmountOTC%d", i))
				mxTirAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("MaxTieredAmountOTC%d", i))
				feTirAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("FeeTieredAmountOTC%d", i))
				idTirAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("IDTieredAmountOTC%d", i))
				if idTirAmtOtc == "" {
					cpcttotc = append(cpcttotc, &ptnrcom.PartnerCommissionTier{
						PartnerCommissionID: tierAmtResID,
						MinValue:            mnTirAmtOtc,
						MaxValue:            mxTirAmtOtc,
						Amount:              feTirAmtOtc,
					})
				} else {
					upcttotc = append(upcttotc, &ptnrcom.PartnerCommissionTier{
						ID:                  idTirAmtOtc,
						PartnerCommissionID: tierAmtResID,
						MinValue:            mnTirAmtOtc,
						MaxValue:            mxTirAmtOtc,
						Amount:              feTirAmtOtc,
					})
				}
			}
			if len(cpcttotc) > 0 {
				_, err := s.pf.CreatePartnerCommissionTier(ctx, &ptnrcom.CreatePartnerCommissionTierRequest{
					CommissionTier: cpcttotc,
				})
				if err != nil {
					log.Error("failed to create partner commission tier")
					return errors.New("failed to create partner commission tier")
				}
			}
			if len(upcttotc) > 0 {
				_, err := s.pf.UpdatePartnerCommissionTier(ctx, &ptnrcom.UpdatePartnerCommissionTierRequest{
					CommissionTier: upcttotc,
				})
				if err != nil {
					log.Error("failed to update partner commission tier")
					return errors.New("failed to update partner commission tier")
				}
			}
		}
	case "TIERPERCENTAGE":
		lnTirPrAmnOtc := getFormValue(r.Form, config.Pref, "LenTieredPercentageOTC")
		lnTirPrAmnOtcInt, _ := strconv.Atoi(lnTirPrAmnOtc)
		tierPerResID := ""
		if transactionOTCID == "" {
			tierPerRes, err := s.pf.CreatePartnerCommission(ctx, &ptnrcom.CreatePartnerCommissionRequest{
				Partner:         config.PtVStype,
				BoundType:       ptnrcom.BoundType(ptnrcom.BoundType_value[config.BtK]),
				RemitType:       ptnrcom.RemitType(ptnrcom.RemitType_value[config.RtK]),
				TransactionType: ptnrcom.TransactionType(ptnrcom.TransactionType_value[config.TrOtc]),
				TierType:        ptnrcom.TierType(ptnrcom.TierType_value[trTypeOtc]),
				Amount:          "",
				CreatedBy:       config.Uid,
			})
			if err != nil {
				log.Error("failed to create partner commission")
				return errors.New("failed to create partner commission")
			}
			tierPerResID = tierPerRes.ID
		} else {
			tierPerRes, err := s.pf.UpdatePartnerCommission(ctx, &ptnrcom.UpdatePartnerCommissionRequest{
				ID:              transactionOTCID,
				Partner:         config.PtVStype,
				BoundType:       ptnrcom.BoundType(ptnrcom.BoundType_value[config.BtK]),
				RemitType:       ptnrcom.RemitType(ptnrcom.RemitType_value[config.RtK]),
				TransactionType: ptnrcom.TransactionType(ptnrcom.TransactionType_value[config.TrOtc]),
				TierType:        ptnrcom.TierType(ptnrcom.TierType_value[trTypeOtc]),
				Amount:          "",
				CreatedBy:       config.Uid,
			})
			if err != nil {
				log.Error("failed to update partner commission")
				return errors.New("failed to update partner commission")
			}
			tierPerResID = tierPerRes.ID
		}
		if tierPerResID != "" && lnTirPrAmnOtcInt > 0 {
			cpctotc := []*ptnrcom.PartnerCommissionTier{}
			upctotc := []*ptnrcom.PartnerCommissionTier{}
			for i := 0; i < lnTirPrAmnOtcInt; i++ {
				mnTirPrAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("MinTieredPercentageOTC%d", i))
				mxTirPrAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("MaxTieredPercentageOTC%d", i))
				feTirPrAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("FeeTieredPercentageOTC%d", i))
				idTirPrAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("IDTieredPercentageOTC%d", i))
				if idTirPrAmtOtc == "" {
					cpctotc = append(cpctotc, &ptnrcom.PartnerCommissionTier{
						PartnerCommissionID: tierPerResID,
						MinValue:            mnTirPrAmtOtc,
						MaxValue:            mxTirPrAmtOtc,
						Amount:              feTirPrAmtOtc,
					})
				} else {
					upctotc = append(upctotc, &ptnrcom.PartnerCommissionTier{
						ID:                  idTirPrAmtOtc,
						PartnerCommissionID: tierPerResID,
						MinValue:            mnTirPrAmtOtc,
						MaxValue:            mxTirPrAmtOtc,
						Amount:              feTirPrAmtOtc,
					})
				}
			}
			if len(cpctotc) > 0 {
				_, err := s.pf.CreatePartnerCommissionTier(ctx, &ptnrcom.CreatePartnerCommissionTierRequest{
					CommissionTier: cpctotc,
				})
				if err != nil {
					log.Error("failed to create partner commission tier")
					return errors.New("failed to create partner commission tier")
				}
			}
			if len(upctotc) > 0 {
				_, err := s.pf.UpdatePartnerCommissionTier(ctx, &ptnrcom.UpdatePartnerCommissionTierRequest{
					CommissionTier: upctotc,
				})
				if err != nil {
					log.Error("failed to update partner commission tier")
					return errors.New("failed to update partner commission tier")
				}
			}
		}
	}
	return nil
}

func (s *Server) saveDigitalFixedCommission(r *http.Request, config CommissionMgtConfig) error {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	if config.TrDgt == "" || config.TrDgt != "DIGITAL" {
		_, err := s.pf.DeletePartnerCommission(ctx, &ptnrcom.DeletePartnerCommissionRequest{
			Partner:         config.PtVStype,
			BoundType:       ptnrcom.BoundType(ptnrcom.BoundType_value[config.BtK]),
			RemitType:       ptnrcom.RemitType(ptnrcom.RemitType_value[config.RtK]),
			TransactionType: ptnrcom.TransactionType(ptnrcom.TransactionType_DIGITAL),
		})
		if err != nil {
			return err
		}
		return errors.New("transaction type should be Digital")
	}
	trTypeDgt := getFormValue(r.Form, config.Pref, "TransactionTypeDigital")
	transactionDIGITALID := getFormValue(r.Form, config.Pref, "TransactionDIGITALID")
	comId := transactionDIGITALID
	switch trTypeDgt {
	case "FIXED":
		fixAmtDgt := getFormValue(r.Form, config.Pref, "FixedAmountDigital")
		sDtDgtFx := getFormValue(r.Form, config.Pref, "StartDateDigitalFixed")
		eDtDgtFx := getFormValue(r.Form, config.Pref, "EndDateDigitalFixed")
		if transactionDIGITALID == "" {
			comres, err := s.pf.CreatePartnerCommission(ctx, &ptnrcom.CreatePartnerCommissionRequest{
				Partner:         config.PtVStype,
				BoundType:       ptnrcom.BoundType(ptnrcom.BoundType_value[config.BtK]),
				RemitType:       ptnrcom.RemitType(ptnrcom.RemitType_value[config.RtK]),
				TransactionType: ptnrcom.TransactionType(ptnrcom.TransactionType_value[config.TrDgt]),
				TierType:        ptnrcom.TierType(ptnrcom.TierType_value[trTypeDgt]),
				Amount:          fixAmtDgt,
				StartDate:       strToTimeConvert(sDtDgtFx),
				EndDate:         strToTimeConvert(eDtDgtFx),
				CreatedBy:       config.Uid,
			})
			if err != nil {
				log.Error("failed to create partner commission")
				return errors.New("failed to create partner commission")
			}
			comId = comres.GetID()
		} else {
			_, err := s.pf.UpdatePartnerCommission(ctx, &ptnrcom.UpdatePartnerCommissionRequest{
				ID:              transactionDIGITALID,
				Partner:         config.PtVStype,
				BoundType:       ptnrcom.BoundType(ptnrcom.BoundType_value[config.BtK]),
				RemitType:       ptnrcom.RemitType(ptnrcom.RemitType_value[config.RtK]),
				TransactionType: ptnrcom.TransactionType(ptnrcom.TransactionType_value[config.TrDgt]),
				TierType:        ptnrcom.TierType(ptnrcom.TierType_value[trTypeDgt]),
				Amount:          fixAmtDgt,
				StartDate:       strToTimeConvert(sDtDgtFx),
				EndDate:         strToTimeConvert(eDtDgtFx),
				CreatedBy:       config.Uid,
			})
			if err != nil {
				log.Error("failed to update partner commission")
				return errors.New("failed to update partner commission")
			}
		}
		if _, err := s.pf.DeletePartnerCommissionTier(ctx, &ptnrcom.DeletePartnerCommissionTierRequest{
			PartnerCommissionID: comId,
		}); err != nil {
			log.Error("failed to delete partner commission tier")
			return errors.New("failed to delete partner commission tier")
		}
	case "PERCENTAGE":
		perAmtDgt := getFormValue(r.Form, config.Pref, "FixedPercentageDigital")
		sDtDgtPer := getFormValue(r.Form, config.Pref, "StartDateDigitalPercentage")
		eDtDgtPer := getFormValue(r.Form, config.Pref, "EndDateDigitalPercentage")
		if transactionDIGITALID == "" {
			comIdres, err := s.pf.CreatePartnerCommission(ctx, &ptnrcom.CreatePartnerCommissionRequest{
				Partner:         config.PtVStype,
				BoundType:       ptnrcom.BoundType(ptnrcom.BoundType_value[config.BtK]),
				RemitType:       ptnrcom.RemitType(ptnrcom.RemitType_value[config.RtK]),
				TransactionType: ptnrcom.TransactionType(ptnrcom.TransactionType_value[config.TrDgt]),
				TierType:        ptnrcom.TierType(ptnrcom.TierType_value[trTypeDgt]),
				Amount:          perAmtDgt,
				StartDate:       strToTimeConvert(sDtDgtPer),
				EndDate:         strToTimeConvert(eDtDgtPer),
				CreatedBy:       config.Uid,
			})
			if err != nil {
				log.Error("failed to create partner commission")
				return errors.New("failed to create partner commission")
			}
			comId = comIdres.GetID()
		} else {
			_, err := s.pf.UpdatePartnerCommission(ctx, &ptnrcom.UpdatePartnerCommissionRequest{
				ID:              transactionDIGITALID,
				Partner:         config.PtVStype,
				BoundType:       ptnrcom.BoundType(ptnrcom.BoundType_value[config.BtK]),
				RemitType:       ptnrcom.RemitType(ptnrcom.RemitType_value[config.RtK]),
				TransactionType: ptnrcom.TransactionType(ptnrcom.TransactionType_value[config.TrDgt]),
				TierType:        ptnrcom.TierType(ptnrcom.TierType_value[trTypeDgt]),
				Amount:          perAmtDgt,
				StartDate:       strToTimeConvert(sDtDgtPer),
				EndDate:         strToTimeConvert(eDtDgtPer),
				CreatedBy:       config.Uid,
			})
			if err != nil {
				log.Error("failed to update partner commission")
				return errors.New("failed to update partner commission")
			}
		}
		if _, err := s.pf.DeletePartnerCommissionTier(ctx, &ptnrcom.DeletePartnerCommissionTierRequest{
			PartnerCommissionID: comId,
		}); err != nil {
			log.Error("failed to delete partner commission tier")
			return errors.New("failed to delete partner commission tier")
		}
	case "TIERAMOUNT":
		lnTirAmtDgt := getFormValue(r.Form, config.Pref, "LenTieredAmountDigital")
		lnTirAmtDgtInt, _ := strconv.Atoi(lnTirAmtDgt)
		tierAmtResId := ""
		if transactionDIGITALID == "" {
			tierAmtRes, err := s.pf.CreatePartnerCommission(ctx, &ptnrcom.CreatePartnerCommissionRequest{
				Partner:         config.PtVStype,
				BoundType:       ptnrcom.BoundType(ptnrcom.BoundType_value[config.BtK]),
				RemitType:       ptnrcom.RemitType(ptnrcom.RemitType_value[config.RtK]),
				TransactionType: ptnrcom.TransactionType(ptnrcom.TransactionType_value[config.TrDgt]),
				TierType:        ptnrcom.TierType(ptnrcom.TierType_value[trTypeDgt]),
				Amount:          "",
				CreatedBy:       config.Uid,
			})
			if err != nil {
				log.Error("failed to create partner commission")
				return errors.New("failed to create partner commission")
			}
			tierAmtResId = tierAmtRes.ID
		} else {
			tierAmtRes, err := s.pf.UpdatePartnerCommission(ctx, &ptnrcom.UpdatePartnerCommissionRequest{
				ID:              transactionDIGITALID,
				Partner:         config.PtVStype,
				BoundType:       ptnrcom.BoundType(ptnrcom.BoundType_value[config.BtK]),
				RemitType:       ptnrcom.RemitType(ptnrcom.RemitType_value[config.RtK]),
				TransactionType: ptnrcom.TransactionType(ptnrcom.TransactionType_value[config.TrDgt]),
				TierType:        ptnrcom.TierType(ptnrcom.TierType_value[trTypeDgt]),
				Amount:          "",
				CreatedBy:       config.Uid,
			})
			if err != nil {
				log.Error("failed to update partner commission")
				return errors.New("failed to update partner commission")
			}
			tierAmtResId = tierAmtRes.ID
		}
		if tierAmtResId != "" && lnTirAmtDgtInt > 0 {
			cpct := []*ptnrcom.PartnerCommissionTier{}
			upct := []*ptnrcom.PartnerCommissionTier{}
			for i := 0; i < lnTirAmtDgtInt; i++ {
				mnTirAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("MinTieredAmountDigital%d", i))
				mxTirAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("MaxTieredAmountDigital%d", i))
				feTirAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("FeeTieredAmountDigital%d", i))
				idTirAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("IDTieredAmountDigital%d", i))
				if idTirAmtDgt == "" {
					cpct = append(cpct, &ptnrcom.PartnerCommissionTier{
						PartnerCommissionID: tierAmtResId,
						MinValue:            mnTirAmtDgt,
						MaxValue:            mxTirAmtDgt,
						Amount:              feTirAmtDgt,
					})
				} else {
					upct = append(upct, &ptnrcom.PartnerCommissionTier{
						ID:                  idTirAmtDgt,
						PartnerCommissionID: tierAmtResId,
						MinValue:            mnTirAmtDgt,
						MaxValue:            mxTirAmtDgt,
						Amount:              feTirAmtDgt,
					})
				}
			}
			if len(cpct) > 0 {
				_, err := s.pf.CreatePartnerCommissionTier(ctx, &ptnrcom.CreatePartnerCommissionTierRequest{
					CommissionTier: cpct,
				})
				if err != nil {
					log.Error("failed to create partner commission tier")
					return errors.New("failed to create partner commission tier")
				}
			}
			if len(upct) > 0 {
				_, err := s.pf.UpdatePartnerCommissionTier(ctx, &ptnrcom.UpdatePartnerCommissionTierRequest{
					CommissionTier: upct,
				})
				if err != nil {
					log.Error("failed to update partner commission tier")
					return errors.New("failed to update partner commission tier")
				}
			}
		}
	case "TIERPERCENTAGE":
		lnTirPrAmnDgt := getFormValue(r.Form, config.Pref, "LenTieredPercentageDigital")
		lnTirPrAmnDgtInt, _ := strconv.Atoi(lnTirPrAmnDgt)
		tierPerResID := ""
		if transactionDIGITALID == "" {
			tierPerRes, err := s.pf.CreatePartnerCommission(ctx, &ptnrcom.CreatePartnerCommissionRequest{
				Partner:         config.PtVStype,
				BoundType:       ptnrcom.BoundType(ptnrcom.BoundType_value[config.BtK]),
				RemitType:       ptnrcom.RemitType(ptnrcom.RemitType_value[config.RtK]),
				TransactionType: ptnrcom.TransactionType(ptnrcom.TransactionType_value[config.TrDgt]),
				TierType:        ptnrcom.TierType(ptnrcom.TierType_value[trTypeDgt]),
				Amount:          "",
				CreatedBy:       config.Uid,
			})
			if err != nil {
				log.Error("failed to create partner commission")
				return errors.New("failed to create partner commission")
			}
			tierPerResID = tierPerRes.ID
		} else {
			tierPerRes, err := s.pf.UpdatePartnerCommission(ctx, &ptnrcom.UpdatePartnerCommissionRequest{
				ID:              transactionDIGITALID,
				Partner:         config.PtVStype,
				BoundType:       ptnrcom.BoundType(ptnrcom.BoundType_value[config.BtK]),
				RemitType:       ptnrcom.RemitType(ptnrcom.RemitType_value[config.RtK]),
				TransactionType: ptnrcom.TransactionType(ptnrcom.TransactionType_value[config.TrDgt]),
				TierType:        ptnrcom.TierType(ptnrcom.TierType_value[trTypeDgt]),
				Amount:          "",
				CreatedBy:       config.Uid,
			})
			if err != nil {
				log.Error("failed to update partner commission")
				return errors.New("failed to update partner commission")
			}
			tierPerResID = tierPerRes.ID
		}
		if tierPerResID != "" && lnTirPrAmnDgtInt > 0 {
			cpctdgt := []*ptnrcom.PartnerCommissionTier{}
			upctdgt := []*ptnrcom.PartnerCommissionTier{}
			for i := 0; i < lnTirPrAmnDgtInt; i++ {
				mnTirPrAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("MinTieredPercentageDigital%d", i))
				mxTirPrAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("MaxTieredPercentageDigital%d", i))
				feTirPrAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("FeeTieredPercentageDigital%d", i))
				idTirPrAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("IDTieredPercentageDigital%d", i))
				if idTirPrAmtDgt == "" {
					cpctdgt = append(cpctdgt, &ptnrcom.PartnerCommissionTier{
						PartnerCommissionID: tierPerResID,
						MinValue:            mnTirPrAmtDgt,
						MaxValue:            mxTirPrAmtDgt,
						Amount:              feTirPrAmtDgt,
					})
				} else {
					upctdgt = append(upctdgt, &ptnrcom.PartnerCommissionTier{
						ID:                  idTirPrAmtDgt,
						PartnerCommissionID: tierPerResID,
						MinValue:            mnTirPrAmtDgt,
						MaxValue:            mxTirPrAmtDgt,
						Amount:              feTirPrAmtDgt,
					})
				}
			}

			if len(cpctdgt) > 0 {
				_, err := s.pf.CreatePartnerCommissionTier(ctx, &ptnrcom.CreatePartnerCommissionTierRequest{
					CommissionTier: cpctdgt,
				})
				if err != nil {
					log.Error("failed to create partner commission tier")
					return errors.New("failed to create partner commission tier")
				}
			}
			if len(upctdgt) > 0 {
				_, err := s.pf.UpdatePartnerCommissionTier(ctx, &ptnrcom.UpdatePartnerCommissionTierRequest{
					CommissionTier: upctdgt,
				})
				if err != nil {
					log.Error("failed to update partner commission tier")
					return errors.New("failed to update partner commission tier")
				}
			}
		}
	}
	return nil
}

func extractDynamicArray(form url.Values, key string) (result map[string]string, err error) {
	result = make(map[string]string)
	reg, err := regexp.Compile(`^([a-z0-9]|[A-z0-9]+)\[([a-z0-9]|[A-z0-9]+)\]$`)
	if err != nil {
		return
	}
	var matches [][]string
	for k, v := range form {
		matches = reg.FindAllStringSubmatch(k, -1)
		if len(matches) != 1 {
			continue
		}
		if key != "" && matches[0][1] != key {
			continue
		}
		if len(matches[0]) != 3 {
			continue
		}
		val := ""
		if len(v) > 0 {
			val = v[0]
		}
		result[matches[0][2]] = val
	}
	return
}

func getFormValue(f url.Values, pref string, key string) string {
	vMap, err := extractDynamicArray(f, key)
	if err != nil {
		return ""
	}
	return vMap[pref]
}

func strToTimeConvert(str string) *timestamppb.Timestamp {
	layout := "2006-01-02"
	t1, err := time.Parse(layout, str)
	if err != nil {
		return timestamppb.Now()
	}
	return timestamppb.New(t1)
}

func (s *Server) validateDigitalFixedCommission(r *http.Request, config CommissionMgtConfig) error {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	if config.TrDgt == "" || config.TrDgt != "DIGITAL" {
		return nil
	}
	trTypeDgt := getFormValue(r.Form, config.Pref, "TransactionTypeDigital")
	switch trTypeDgt {
	case "FIXED":
		fixAmtDgt := getFormValue(r.Form, config.Pref, "FixedAmountDigital")
		fixAmtDgtF, err := strconv.ParseFloat(fixAmtDgt, 64)
		if err != nil {
			logging.WithError(err, log).Error("Value must be greater than 0")
			return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
		}
		if fixAmtDgtF < 0 {
			logging.WithError(err, log).Error("Value must be greater than 0 for partner " + config.PtVStype)
			return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
		}

	case "PERCENTAGE":
		perAmtDgt := getFormValue(r.Form, config.Pref, "FixedPercentageDigital")
		perAmtDgtF, err := strconv.ParseFloat(perAmtDgt, 64)
		if err != nil {
			logging.WithError(err, log).Error("The Value should be greater than 0 and less than 100")
			return errors.New("The Value should be greater than 0 and less than 100 for partner " + config.PtVStype)
		}
		if perAmtDgtF < 0 {
			logging.WithError(err, log).Error("The Value should be greater than 0 and less than 100")
			return errors.New("The Value should be greater than 0 and less than 100 for partner " + config.PtVStype)
		}

	case "TIERAMOUNT":
		lnTirAmtDgt := getFormValue(r.Form, config.Pref, "LenTieredAmountDigital")
		lnTirAmtDgtInt, _ := strconv.Atoi(lnTirAmtDgt)
		if lnTirAmtDgtInt == 0 {
			return nil
		}
		for i := 0; i < lnTirAmtDgtInt; i++ {
			mnTirAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("MinTieredAmountDigital%d", i))
			mxTirAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("MaxTieredAmountDigital%d", i))
			feTirAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("FeeTieredAmountDigital%d", i))
			if mnTirAmtDgt == "" && mxTirAmtDgt == "" && feTirAmtDgt == "" {
				continue
			}
			mnTirAmtDgtF, err := strconv.ParseFloat(mnTirAmtDgt, 64)
			if err != nil {
				logging.WithError(err, log).Error("Value must be greater than 0")
				return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
			}
			mxTirAmtDgtF, err := strconv.ParseFloat(mxTirAmtDgt, 64)
			if err != nil {
				logging.WithError(err, log).Error("Value must be greater than 0")
				return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
			}
			feTirAmtDgtF, err := strconv.ParseFloat(feTirAmtDgt, 64)
			if err != nil {
				logging.WithError(err, log).Error("Value must be greater than 0")
				return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
			}
			if feTirAmtDgtF < 0 || feTirAmtDgtF == 0 {
				logging.WithError(err, log).Error("Value must be greater than 0")
				return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
			}
			if mnTirAmtDgtF >= mxTirAmtDgtF {
				logging.WithError(err, log).Error("Minimum amount should not be greater than the maximum amount")
				return errors.New("Minimum amount should not be greater than the maximum amount for partner " + config.PtVStype)
			}
			if i != 0 {
				prvMxTirAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("MaxTieredAmountDigital%d", i-1))
				prvMxTirAmtDgtF, err := strconv.ParseFloat(prvMxTirAmtDgt, 64)
				if err != nil {
					logging.WithError(err, log).Error("Value must be greater than 0")
					return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
				}
				if prvMxTirAmtDgtF > mnTirAmtDgtF {
					logging.WithError(err, log).Error("Amount should not be within range of the previous input.")
					return errors.New("Minimum amount should not be within range of the previous input for partner " + config.PtVStype)
				}
			}
			if i < (lnTirAmtDgtInt - 1) {
				nextmnTirAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("MinTieredAmountDigital%d", i+1))
				nextmnTirAmtDgtF, err := strconv.ParseFloat(nextmnTirAmtDgt, 64)
				if err != nil {
					logging.WithError(err, log).Error("Value must be greater than 0")
					return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
				}
				if nextmnTirAmtDgtF < mxTirAmtDgtF {
					logging.WithError(err, log).Error("Amount should not be within range of the previous input.")
					return errors.New("Maximum amount should not be within range of the previous input for partner " + config.PtVStype)
				}
			}
		}
	case "TIERPERCENTAGE":
		lnTirPrAmnDgt := getFormValue(r.Form, config.Pref, "LenTieredPercentageDigital")
		lnTirPrAmnDgtInt, _ := strconv.Atoi(lnTirPrAmnDgt)
		if lnTirPrAmnDgtInt == 0 {
			return nil
		}
		for i := 0; i < lnTirPrAmnDgtInt; i++ {
			mnTirPrAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("MinTieredPercentageDigital%d", i))
			mxTirPrAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("MaxTieredPercentageDigital%d", i))
			feTirPrAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("FeeTieredPercentageDigital%d", i))
			mnTirPrAmtDgtF, err := strconv.ParseFloat(mnTirPrAmtDgt, 64)
			if err != nil {
				logging.WithError(err, log).Error("Value must be greater than 0")
				return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
			}
			mxTirPrAmtDgtF, err := strconv.ParseFloat(mxTirPrAmtDgt, 64)
			if err != nil {
				logging.WithError(err, log).Error("Value must be greater than 0")
				return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
			}
			feTirPrAmtDgtF, err := strconv.ParseFloat(feTirPrAmtDgt, 64)
			if err != nil {
				logging.WithError(err, log).Error("The Value should be greater than 0 and less than 100")
				return errors.New("The Value should be greater than 0 and less than 100 for partner " + config.PtVStype)
			}
			if feTirPrAmtDgtF < 0 || feTirPrAmtDgtF == 0 || feTirPrAmtDgtF > 100 {
				logging.WithError(err, log).Error("The Value should be greater than 0 and less than 100")
				return errors.New("The Value should be greater than 0 and less than 100 for partner " + config.PtVStype)
			}
			if mnTirPrAmtDgtF >= mxTirPrAmtDgtF {
				logging.WithError(err, log).Error("Minimum amount should not be greater than the maximum amount")
				return errors.New("Minimum amount should not be greater than the maximum amount for partner " + config.PtVStype)
			}
			if i != 0 {
				prvmxTirPrAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("MaxTieredPercentageDigital%d", i-1))
				prvmxTirPrAmtDgtF, err := strconv.ParseFloat(prvmxTirPrAmtDgt, 64)
				if err != nil {
					logging.WithError(err, log).Error("Value must be greater than 0")
					return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
				}
				if prvmxTirPrAmtDgtF > mnTirPrAmtDgtF {
					logging.WithError(err, log).Error("Amount should not be within range of the previous input.")
					return errors.New("Minimum amount should not be within range of the previous input for partner " + config.PtVStype)
				}
			}
			if i < (lnTirPrAmnDgtInt - 1) {
				nextmnTirPrAmtDgt := getFormValue(r.Form, config.Pref, fmt.Sprintf("MinTieredPercentageDigital%d", i+1))
				nextmnTirPrAmtDgtF, err := strconv.ParseFloat(nextmnTirPrAmtDgt, 64)
				if err != nil {
					logging.WithError(err, log).Error("Value must be greater than 0")
					return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
				}
				if nextmnTirPrAmtDgtF < mxTirPrAmtDgtF {
					logging.WithError(err, log).Error("Amount should not be within range of the previous input.")
					return errors.New("Maximum amount should not be within range of the previous input for partner " + config.PtVStype)
				}
			}
		}
	}
	return nil
}

func (s *Server) validateOTCFixedCommission(r *http.Request, config CommissionMgtConfig) error {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	if config.TrOtc == "" || config.TrOtc != "OTC" {
		return nil
	}
	trTypeOtc := getFormValue(r.Form, config.Pref, "TransactionTypeOTC")
	switch trTypeOtc {
	case "FIXED":
		fixAmtOtc := getFormValue(r.Form, config.Pref, "FixedAmountOTC")
		fixAmtOtcF, err := strconv.ParseFloat(fixAmtOtc, 64)
		if err != nil {
			logging.WithError(err, log).Error("Value must be greater than 0")
			return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
		}
		if fixAmtOtcF < 0 {
			logging.WithError(err, log).Error("Value must be greater than 0")
			return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
		}
	case "PERCENTAGE":
		perAmtOtc := getFormValue(r.Form, config.Pref, "FixedPercentageOTC")
		perAmtOtcF, err := strconv.ParseFloat(perAmtOtc, 64)
		if err != nil {
			logging.WithError(err, log).Error("The Value should be greater than 0 and less than 100")
			return errors.New("The Value should be greater than 0 and less than 100 for partner " + config.PtVStype)
		}
		if perAmtOtcF < 0 {
			logging.WithError(err, log).Error("The Value should be greater than 0 and less than 100")
			return errors.New("The Value should be greater than 0 and less than 100 for partner " + config.PtVStype)
		}
	case "TIERAMOUNT":
		lnTirAmtOtc := getFormValue(r.Form, config.Pref, "LenTieredAmountOTC")
		lnTirAmtOtcInt, _ := strconv.Atoi(lnTirAmtOtc)
		if lnTirAmtOtcInt == 0 {
			return nil
		}
		for i := 0; i < lnTirAmtOtcInt; i++ {
			mnTirAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("MinTieredAmountOTC%d", i))
			mxTirAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("MaxTieredAmountOTC%d", i))
			feTirAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("FeeTieredAmountOTC%d", i))
			if mnTirAmtOtc == "" && mxTirAmtOtc == "" && feTirAmtOtc == "" {
				continue
			}
			mnTirAmtOtcF, err := strconv.ParseFloat(mnTirAmtOtc, 64)
			if err != nil {
				logging.WithError(err, log).Error("Value must be greater than 0")
				return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
			}
			mxTirAmtOtcF, err := strconv.ParseFloat(mxTirAmtOtc, 64)
			if err != nil {
				logging.WithError(err, log).Error("Value must be greater than 0")
				return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
			}
			feTirAmtOtcF, err := strconv.ParseFloat(feTirAmtOtc, 64)
			if err != nil {
				logging.WithError(err, log).Error("Value must be greater than 0")
				return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
			}
			if feTirAmtOtcF < 0 || feTirAmtOtcF == 0 {
				logging.WithError(err, log).Error("Value must be greater than 0")
				return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
			}
			if mnTirAmtOtcF >= mxTirAmtOtcF {
				logging.WithError(err, log).Error("Minimum amount should not be greater than the maximum amount")
				return errors.New("Minimum amount should not be greater than the maximum amount for partner " + config.PtVStype)
			}
			if i != 0 {
				prvmxTirAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("MaxTieredAmountOTC%d", i-1))
				prvmxTirAmtOtcF, err := strconv.ParseFloat(prvmxTirAmtOtc, 64)
				if err != nil {
					logging.WithError(err, log).Error("Value must be greater than 0")
					return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
				}
				if prvmxTirAmtOtcF > mnTirAmtOtcF {
					logging.WithError(err, log).Error("Amount should not be within range of the previous input.")
					return errors.New("Minimum amount should not be within range of the previous input for partner " + config.PtVStype)
				}
			}
			if i < (lnTirAmtOtcInt - 1) {
				nextmnTirAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("MinTieredAmountOTC%d", i+1))
				nextmnTirAmtOtcF, err := strconv.ParseFloat(nextmnTirAmtOtc, 64)
				if err != nil {
					logging.WithError(err, log).Error("Value must be greater than 0")
					return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
				}
				if nextmnTirAmtOtcF < mxTirAmtOtcF {
					logging.WithError(err, log).Error("Amount should not be within range of the previous input.")
					return errors.New("Maximum amount should not be within range of the previous input for partner " + config.PtVStype)
				}
			}
		}
	case "TIERPERCENTAGE":
		lnTirPrAmnOtc := getFormValue(r.Form, config.Pref, "LenTieredPercentageOTC")
		lnTirPrAmnOtcInt, _ := strconv.Atoi(lnTirPrAmnOtc)
		if lnTirPrAmnOtcInt == 0 {
			return nil
		}
		for i := 0; i < lnTirPrAmnOtcInt; i++ {
			mnTirPrAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("MinTieredPercentageOTC%d", i))
			mxTirPrAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("MaxTieredPercentageOTC%d", i))
			feTirPrAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("FeeTieredPercentageOTC%d", i))
			if mnTirPrAmtOtc == "" && mxTirPrAmtOtc == "" && feTirPrAmtOtc == "" {
				continue
			}
			mnTirPrAmtOtcF, err := strconv.ParseFloat(mnTirPrAmtOtc, 64)
			if err != nil {
				logging.WithError(err, log).Error("Value must be greater than 0")
				return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
			}
			mxTirPrAmtOtcF, err := strconv.ParseFloat(mxTirPrAmtOtc, 64)
			if err != nil {
				logging.WithError(err, log).Error("Value must be greater than 0")
				return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
			}
			feTirPrAmtOtcF, err := strconv.ParseFloat(feTirPrAmtOtc, 64)
			if err != nil {
				logging.WithError(err, log).Error("The Value should be greater than 0 and less than 100")
				return errors.New("The Value should be greater than 0 and less than 100 for partner " + config.PtVStype)
			}
			if feTirPrAmtOtcF < 0 || feTirPrAmtOtcF == 0 || feTirPrAmtOtcF > 100 {
				logging.WithError(err, log).Error("The Value should be greater than 0 and less than 100")
				return errors.New("The Value should be greater than 0 and less than 100 for partner " + config.PtVStype)
			}
			if mnTirPrAmtOtcF >= mxTirPrAmtOtcF {
				logging.WithError(err, log).Error("Minimum amount should not be greater than the maximum amount")
				return errors.New("Minimum amount should not be greater than the maximum amount for partner " + config.PtVStype)
			}
			if i != 0 {
				prvmxTirPrAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("MaxTieredPercentageOTC%d", i-1))
				prvmxTirPrAmtOtcF, err := strconv.ParseFloat(prvmxTirPrAmtOtc, 64)
				if err != nil {
					logging.WithError(err, log).Error("Value must be greater than 0")
					return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
				}
				if prvmxTirPrAmtOtcF > mnTirPrAmtOtcF {
					logging.WithError(err, log).Error("Amount should not be within range of the previous input.")
					return errors.New("Minimum amount should not be within range of the previous input for partner " + config.PtVStype)
				}
			}
			if i < (lnTirPrAmnOtcInt - 1) {
				nextmnTirPrAmtOtc := getFormValue(r.Form, config.Pref, fmt.Sprintf("MinTieredPercentageOTC%d", i+1))
				nextmnTirPrAmtOtcF, err := strconv.ParseFloat(nextmnTirPrAmtOtc, 64)
				if err != nil {
					logging.WithError(err, log).Error("Value must be greater than 0")
					return errors.New("Value must be greater than 0 for partner " + config.PtVStype)
				}
				if nextmnTirPrAmtOtcF < mxTirPrAmtOtcF {
					logging.WithError(err, log).Error("Amount should not be within range of the previous input.")
					return errors.New("Maximum amount should not be within range of the previous input for partner " + config.PtVStype)
				}
			}
		}
	}
	return nil
}

func (s *Server) validateRemitTypesValue(r *http.Request, config saveRemitTypesValueReq) error {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	if len(config.RtVMaps) == 0 {
		log.Error("remitTypes shouldn't empty")
		return errors.New("remitTypes shouldn't empty")
	}
	for rtK, rtV := range config.RtVMaps {
		if rtK == "" || rtV == "" {
			log.Error("remit type shouldn't empty")
			return errors.New("remit type shouldn't empty")
		}
		for _, btVMaps := range config.BoundTypes {
			if err := s.validateBoundTypesValue(r, saveBoundTypesValueReq{
				PartnerList: config.PartnerList,
				BtVMaps:     btVMaps,
				Uid:         config.Uid,
				RtK:         rtK,
				RtV:         rtV,
			}); err != nil {
				log.Error("remitTypes shouldn't empty")
				return err
			}
		}
	}
	return nil
}

func (s *Server) validateBoundTypesValue(r *http.Request, config saveBoundTypesValueReq) error {
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
			trDgt := getFormValue(r.Form, pref, "TransactionDigital")
			trOtc := getFormValue(r.Form, pref, "TransactionOTC")
			if err := s.validateDigitalFixedCommission(r, CommissionMgtConfig{
				Pref:     pref,
				TrDgt:    trDgt,
				PtVStype: ptV.Stype,
				BtK:      btK,
				RtK:      config.RtK,
				Uid:      config.Uid,
			}); err != nil {
				return err
			}
			if err := s.validateOTCFixedCommission(r, CommissionMgtConfig{
				Pref:     pref,
				TrOtc:    trOtc,
				PtVStype: ptV.Stype,
				BtK:      btK,
				RtK:      config.RtK,
				Uid:      config.Uid,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

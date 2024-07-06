package core

import (
	"context"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"brank.as/petnet/gunk/drp/v1/dsa"
	revcom "brank.as/petnet/gunk/drp/v1/revenue-commission"
	spb "brank.as/petnet/gunk/dsa/v2/partner"
	ptnrcom "brank.as/petnet/gunk/dsa/v2/partnercommission"
	spbl "brank.as/petnet/gunk/dsa/v2/partnerlist"
	pfpb "brank.as/petnet/gunk/dsa/v2/profile"
	revsrng "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	pfsvc "brank.as/petnet/gunk/dsa/v2/service"
)

const userDRP = "DRP"

// TODO(vitthal): Move remco ids to partner_list table to make it dynamic
var remcoIDsByRemcoCode = map[string]uint32{
	"IR":   1,
	"TF":   7,
	"BPI":  2,
	"RIA":  12,
	"USSC": 10,
	"AYA":  22,
	"MB":   8,
	"UNT":  20,
	"JPR":  17,
	"IC":   16,
	"RM":   21,
	"CEB":  9,
	"IE":   24,
	// TODO: Make Cebuana Int code consistant in cms and api(CEBINT)
	"CEBI": 19,

	// TODO(vitthal): Get remco id for WU and WISE
	"WISE": 100000,
	"WU":   100001,
}

type iDRP interface {
	revcom.RevenueCommissionServiceClient
	dsa.DSAServiceClient
}

type iProfile interface {
	pfpb.OrgProfileServiceClient
	spbl.PartnerListServiceClient
	ptnrcom.PartnerCommissionServiceClient
	revsrng.RevenueSharingServiceClient
}

type RemcoCommissionSvc struct {
	env   string
	log   *logrus.Entry
	pf    iProfile
	drpSB iDRP
	drpLV iDRP
}

func NewRemcoCommissionSvc(pf iProfile, drpSB, drpLV iDRP, log *logrus.Entry) *RemcoCommissionSvc {
	return &RemcoCommissionSvc{
		log:   log,
		pf:    pf,
		drpSB: drpSB,
		drpLV: drpLV,
	}
}

func (s *RemcoCommissionSvc) resolveDRP() iDRP {
	if s.env == "production" {
		return s.drpLV
	}

	return s.drpSB
}

func (s *RemcoCommissionSvc) SyncRemcoCommissionConfigForRemittance(ctx context.Context) {
	log := s.log
	svc, err := s.pf.GetPartnerList(ctx, &spbl.GetPartnerListRequest{
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: pfsvc.ServiceType_REMITTANCE.String(),
	})
	if err != nil {
		log.WithError(err).Error("SyncRemcoCommissionForRemittance get partner list error")
		return
	}

	if svc == nil || len(svc.GetPartnerList()) == 0 {
		log.Error("SyncRemcoCommissionForRemittance partner list is empty")
		return
	}

	phAllCommRes, err := s.resolveDRP().ListRemcoCommissionFee(ctx, &emptypb.Empty{})
	if err != nil {
		sErr := status.Convert(err)
		// if not "Not found" error
		if !(sErr.Code() == codes.NotFound || int(sErr.Code()) == 404) {
			log.WithError(err).Error("SyncRemcoCommissionForRemittance list drp remco commission error")
			return
		}
	}

	phAllRemcoCommMap := toCommMapByRemco(ctx, phAllCommRes, log, s.pf)
	// TODO(vitthal): Should we delete extra partners in perahub if not present in DRP?
	for _, partner := range svc.GetPartnerList() {
		if partner == nil {
			continue
		}

		provider := partner.GetIsProvider()

		if remcoID := getRemcoIDByRemcoCode(ctx, partner.GetStype(), s.pf); !provider && remcoID <= 0 {
			log.Error("remcoID is empty for partner " + partner.GetStype())
			continue
		}

		phRemcoCommList, ok := phAllRemcoCommMap[partner.GetStype()]
		if !ok {
			log.Info("SyncRemcoCommissionForRemittance partner commission config is not present in perahub for " + partner.GetStype())
		}

		err = s.syncPartnerCommissionConfig(ctx, log, partner.GetStype(), phRemcoCommList)
		if err != nil {
			log.WithError(err).Error("SyncRemcoCommissionForRemittance please fix error for sync")
			return
		}
	}
}

func (s *RemcoCommissionSvc) syncPartnerCommissionConfig(ctx context.Context,
	log *logrus.Entry,
	partner string,
	phRemcoCommList []revcom.RemcoCommissionFee,
) error {
	commRes, err := s.pf.GetPartnerCommissionsList(ctx, &ptnrcom.GetPartnerCommissionsListRequest{
		RemitType: ptnrcom.RemitType_REMITTANCE,
		Partner:   partner,
	})
	if err != nil {
		log.WithError(err).Error("SyncRemcoCommissionForRemittance get partner commission list error for partner " + partner)
		return err
	}

	if commRes != nil && len(commRes.GetResults()) == 0 {
		log.Error("SyncRemcoCommissionForRemittance partner commission list empty for partner " + partner)
		return nil
	}

	drpCommByBound := toDRPPartnerCommListToMapByBoundType(commRes.GetResults())
	phCommByBound := toPerahubPartnerCommListToMapByBoundType(phRemcoCommList)

	if len(drpCommByBound) == 0 {
		log.Error("SyncRemcoCommissionForRemittance drp commission by bound is empty for partner " + partner)
		return nil
	}

	createList := make([]revcom.RemcoCommissionFee, 0, len(commRes.GetResults()))
	deleteList := make([]revcom.RemcoCommissionFee, 0, len(phRemcoCommList))

	for bound, drpPrtnrCommList := range drpCommByBound {
		if len(drpPrtnrCommList) == 0 {
			continue
		}

		drpList, err := s.drpCommissionListOfBoundTypeToPerahubRemcoCommission(ctx, drpPrtnrCommList, partner)
		if err != nil {
			log.WithError(err).Error("SyncRemcoCommissionForRemittance drp commission to perahub commission convert error for partner " + partner)
			return err
		}

		phList := phCommByBound[bound]

		toCreate, toDelete := getCreateAndDeleteList(drpList, phList)

		if len(toCreate) > 0 {
			createList = append(createList, toCreate...)
		}

		if len(toDelete) > 0 {
			deleteList = append(deleteList, toDelete...)
		}
	}

	// delete first
	for _, rc := range deleteList {
		_, err = s.resolveDRP().DeleteRemcoCommissionFee(ctx, &revcom.DeleteRemcoCommissionFeeRequest{FeeID: rc.FeeID})
		if err != nil {
			log.WithError(err).
				WithField("remcoID", rc.RemcoID).WithField("FeeID", rc.FeeID).
				Error("SyncRemcoCommissionForRemittance delete remco commission perahub error")
			continue
		}
	}

	// create new ones
	for _, rc := range createList {
		crRes, err := s.resolveDRP().CreateRemcoCommissionFee(ctx, toCreateRemcoCommissionFeeRequest(rc))
		if err != nil {
			log.WithError(err).
				WithField("remcoID", rc.RemcoID).WithField("FeeID", rc.FeeID).
				Error("SyncRemcoCommissionForRemittance create remco commission perahub error")
			continue
		}

		log.Info(fmt.Sprintf("SyncRemcoCommissionForRemittance remco commission fee created with perahub id %d", crRes.FeeID))
	}

	return nil
}

func toCreateRemcoCommissionFeeRequest(rc revcom.RemcoCommissionFee) *revcom.CreateRemcoCommissionFeeRequest {
	return &revcom.CreateRemcoCommissionFeeRequest{
		RemcoID:             rc.RemcoID,
		MinAmount:           rc.MinAmount,
		MaxAmount:           rc.MaxAmount,
		ServiceFee:          rc.ServiceFee,
		CommissionType:      rc.CommissionType,
		CommissionAmount:    rc.CommissionAmount,
		CommissionAmountOTC: rc.CommissionAmountOTC,
		TrxType:             rc.TrxType,
		UpdatedBy:           userDRP,
	}
}

func getCreateAndDeleteList(drpList []revcom.RemcoCommissionFee,
	phList []revcom.RemcoCommissionFee,
) (createList []revcom.RemcoCommissionFee, deleteList []revcom.RemcoCommissionFee) {
	if len(phList) == 0 {
		return drpList, nil
	}

	drpMap := map[string]revcom.RemcoCommissionFee{}
	for _, drpCm := range drpList {
		drpMap[generateRemcoCommKey(drpCm)] = drpCm
	}

	deleteList = make([]revcom.RemcoCommissionFee, 0, len(phList))
	phMap := map[string]revcom.RemcoCommissionFee{}
	for _, phCm := range phList {
		key := generateRemcoCommKey(phCm)

		// If perahub commission not present in DRP, add to delete list
		if _, ok := drpMap[key]; !ok {
			deleteList = append(deleteList, phCm)
			continue
		}

		// if duplicate values present then delete those
		if _, ok := phMap[key]; ok {
			deleteList = append(deleteList, phCm)
			continue
		}

		phMap[key] = phCm
	}

	createList = make([]revcom.RemcoCommissionFee, 0, len(drpList))
	for _, drpCm := range drpList {
		key := generateRemcoCommKey(drpCm)

		// If drp commission not present in perahub, add to create list
		if _, ok := phMap[key]; !ok {
			createList = append(createList, drpCm)
			continue
		}
	}

	return createList, deleteList
}

// generateRemcoCommKey generates key from remco commission values
// in the format "RC_{CommissionType}_{TrxType}_{MinAmount}_{MaxAmount}_{CommissionAmount}_{CommissionAmountOTC}"
// Ex. "RC_OTC_RANGE_0_100_45_0"
func generateRemcoCommKey(rc revcom.RemcoCommissionFee) string {
	return fmt.Sprintf("RC_%s_%s_%s_%s_%s_%s",
		rc.CommissionType.String(), rc.TrxType.String(), rc.MinAmount, rc.MaxAmount, rc.CommissionAmount, rc.CommissionAmountOTC)
}

func (s *RemcoCommissionSvc) drpCommissionListOfBoundTypeToPerahubRemcoCommission(ctx context.Context, drpList []ptnrcom.PartnerCommission, partner string) ([]revcom.RemcoCommissionFee, error) {
	rangeRemcoCommList := make([]revcom.RemcoCommissionFee, 0, len(drpList))
	for _, drpCm := range drpList {
		if drpCm.TierType == ptnrcom.TierType_TIERAMOUNT || drpCm.TierType == ptnrcom.TierType_TIERPERCENTAGE {
			tiersRes, err := s.pf.GetPartnerCommissionsTierList(ctx, &ptnrcom.GetPartnerCommissionsTierListRequest{
				PartnerCommissionID: drpCm.GetID(),
			})
			if err != nil {
				return nil, err
			}

			rangeRemcoCommList = append(rangeRemcoCommList, commTiersToPerhubRemcoCommissionFeeList(ctx, tiersRes, drpCm, partner, s.pf)...)

			continue
		}

		rc := toPerhubRemcoCommissionFee(ctx, drpCm, drpCm.Amount, "0", "0", partner, s.pf)
		rangeRemcoCommList = append(rangeRemcoCommList, *rc)
	}

	return rangeRemcoCommList, nil
}

func commTiersToPerhubRemcoCommissionFeeList(ctx context.Context, tiersRes *ptnrcom.GetPartnerCommissionsTierListResponse,
	drpCm ptnrcom.PartnerCommission,
	partner string,
	pf iProfile,
) []revcom.RemcoCommissionFee {
	if tiersRes == nil || len(tiersRes.GetResults()) == 0 {
		rc := toPerhubRemcoCommissionFee(ctx, drpCm, drpCm.Amount, "0", "0", partner, pf)
		return []revcom.RemcoCommissionFee{*rc}
	}

	rangeRemcoCommList := make([]revcom.RemcoCommissionFee, 0, len(tiersRes.GetResults()))
	for _, tier := range tiersRes.GetResults() {
		rc := toPerhubRemcoCommissionFee(ctx, drpCm, tier.Amount, tier.MinValue, tier.MaxValue, partner, pf)
		rangeRemcoCommList = append(rangeRemcoCommList, *rc)
	}

	return rangeRemcoCommList
}

func toPerhubRemcoCommissionFee(ctx context.Context, drpCm ptnrcom.PartnerCommission, amt, minAmt, maxAmt, partner string, pf iProfile) *revcom.RemcoCommissionFee {
	commissionType := revcom.CommissionType_CommissionTypeAbsolute
	switch drpCm.TierType {
	case ptnrcom.TierType_TIERAMOUNT:
		commissionType = revcom.CommissionType_CommissionTypeRange
	case ptnrcom.TierType_TIERPERCENTAGE, ptnrcom.TierType_PERCENTAGE:
		commissionType = revcom.CommissionType_CommissionTypePercent
	}

	transactionType := revcom.TrxType_TrxTypeInbound
	if drpCm.BoundType == ptnrcom.BoundType_OUTBOUND {
		transactionType = revcom.TrxType_TrxTypeOutbound
	}

	if _, err := strconv.ParseFloat(amt, 32); err != nil {
		amt = "0"
	}

	digiAmt, otcAmt := amt, "0"
	// If its otc, digital amount is zero
	if drpCm.TransactionType == ptnrcom.TransactionType_OTC {
		otcAmt = amt
		digiAmt = "0"
	}

	// TODO(vitthal): Do we need to pass logged in user name?
	updatedBy := userDRP

	return &revcom.RemcoCommissionFee{
		RemcoID:             getRemcoIDByRemcoCode(ctx, partner, pf),
		MinAmount:           minAmt,
		MaxAmount:           maxAmt,
		ServiceFee:          "0",
		CommissionAmount:    digiAmt,
		CommissionAmountOTC: otcAmt,
		CommissionType:      commissionType,
		TrxType:             transactionType,
		UpdatedBy:           updatedBy,
	}
}

func toDRPPartnerCommListToMapByBoundType(list []*ptnrcom.PartnerCommission) map[ptnrcom.BoundType][]ptnrcom.PartnerCommission {
	commMap := map[ptnrcom.BoundType][]ptnrcom.PartnerCommission{}
	if len(list) == 0 {
		return commMap
	}

	for _, pc := range list {
		if pc.BoundType == ptnrcom.BoundType_OTHERS || pc.BoundType == ptnrcom.BoundType_EMPTYBOUNDTYPE {
			// TODO(vitthal): Should we delete these from DB as they are invalid bound types?
			continue
		}

		commListByBound, ok := commMap[pc.BoundType]
		if !ok {
			// should be maximum 2, one for inbound, one for outbound
			commListByBound = make([]ptnrcom.PartnerCommission, 0, 2)
		}

		commMap[pc.BoundType] = append(commListByBound, *pc)
	}

	return commMap
}

func toPerahubPartnerCommListToMapByBoundType(list []revcom.RemcoCommissionFee) map[ptnrcom.BoundType][]revcom.RemcoCommissionFee {
	commMap := map[ptnrcom.BoundType][]revcom.RemcoCommissionFee{}
	if len(list) == 0 {
		return commMap
	}

	for _, c := range list {
		bndType := ptnrcom.BoundType_INBOUND
		if c.TrxType == revcom.TrxType_TrxTypeOutbound {
			bndType = ptnrcom.BoundType_OUTBOUND
		}

		commListByBound, ok := commMap[bndType]
		if !ok {
			commListByBound = make([]revcom.RemcoCommissionFee, 0, 2)
		}

		commMap[bndType] = append(commListByBound, c)
	}

	return commMap
}

func toCommMapByRemco(ctx context.Context, res *revcom.ListRemcoCommissionFeeResponse, log *logrus.Entry, pf iProfile) map[string][]revcom.RemcoCommissionFee {
	allRemcoCommMap := map[string][]revcom.RemcoCommissionFee{}
	if res == nil || len(res.GetRemcoCommissionFeeList()) == 0 {
		return allRemcoCommMap
	}

	list := res.GetRemcoCommissionFeeList()
	for _, c := range list {
		if c == nil {
			continue
		}

		remcoCode := getRemcoCodeByRemcoID(ctx, c.RemcoID, pf)
		if remcoCode == "" {
			log.Error(fmt.Sprintf("remco code not found for remco ID %d", c.RemcoID))
			continue
		}

		remcoCommList, ok := allRemcoCommMap[remcoCode]
		if !ok {
			remcoCommList = make([]revcom.RemcoCommissionFee, 0, 5)
		}

		allRemcoCommMap[remcoCode] = append(remcoCommList, *c)
	}

	return allRemcoCommMap
}

func getRemcoIDByRemcoCode(ctx context.Context, partnerCode string, pf iProfile) uint32 {
	gpRes, err := pf.GetPartnerByStype(ctx, &spbl.GetPartnerByStypeRequest{
		Stype: partnerCode,
	})
	if err != nil {
		return 0
	}

	if gpRes == nil || gpRes.PartnerList == nil {
		return 0
	}

	id, err := strconv.Atoi(gpRes.PartnerList.RemcoID)
	if err != nil {
		return 0
	}

	return uint32(id)
}

func getRemcoCodeByRemcoID(ctx context.Context, id uint32, pf iProfile) string {
	gpRes, err := pf.GetPartnerList(ctx, &spbl.GetPartnerListRequest{})
	if err != nil {
		return ""
	}

	if gpRes == nil || gpRes.PartnerList == nil || len(gpRes.PartnerList) == 0 {
		return ""
	}

	for _, v := range gpRes.GetPartnerList() {
		if v.RemcoID == "" {
			continue
		}
		remID, err := strconv.Atoi(v.RemcoID)
		if err != nil {
			continue
		}
		if uint32(remID) == id {
			return v.Stype
		}
	}

	return ""
}

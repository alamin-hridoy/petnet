package core

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"brank.as/petnet/api/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"brank.as/petnet/gunk/drp/v1/dsa"
	revcomm "brank.as/petnet/gunk/drp/v1/revenue-commission"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	revshar "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	"brank.as/petnet/serviceutil/logging"
)

const commissionCurrency = "1"

func (s *RemcoCommissionSvc) SyncDSACommissionConfigForRemittance(ctx context.Context, orgID string) {
	log := s.log
	if strings.TrimSpace(orgID) == "" {
		return
	}

	log.WithField("orgID", orgID).Info("synchronizing dsa commission records")

	orgProfile, err := s.pf.GetProfile(ctx, &ppb.GetProfileRequest{OrgID: orgID})
	if err != nil {
		logging.WithError(err, log).WithField("orgID", orgID).
			Error("SyncDSACommissionConfigForRemittance get org profile error")
		return
	}

	if orgProfile == nil || orgProfile.Profile == nil {
		log.WithField("orgID", orgID).Error("SyncDSACommissionConfigForRemittance empty org profile")
		return
	}

	if orgProfile.GetProfile().GetDsaCode() == "" {
		log.WithField("orgID", orgID).Error("SyncDSACommissionConfigForRemittance empty dsa code")
		return
	}

	err = s.upsertPerahubDSACode(ctx, orgProfile.Profile)
	if err != nil {
		logging.WithError(err, log).WithField("orgID", orgID).
			Error("SyncDSACommissionConfigForRemittance create perahub dsa error")
	}

	// Create Perahub and DRP commission map to find out present and absent records on perahub side
	phCommMap, err := s.getPerahubDSACommMap(ctx, orgProfile.Profile.DsaCode)
	if err != nil {
		return
	}

	// Need drp commission tier map to create tier record in perahub
	// Note: we will not delete any tier record on perahub side because of risk of associating with other dsa commission
	drpCommMap, drpCommTiers, err := s.getDRPDSACommAndTierMap(ctx, orgProfile.GetProfile().GetDsaCode(), orgID)
	if err != nil {
		return
	}

	// Find out create and delete list
	createList, deleteList, err := s.getCreateAndDeletePerahubDSACommissionList(ctx, drpCommMap, phCommMap, drpCommTiers)
	if err != nil {
		return
	}

	s.createPerahubDSACommissionRecords(ctx, createList)
	s.deletePerahubDSACommissionRecords(ctx, deleteList)
}

func (s *RemcoCommissionSvc) upsertPerahubDSACode(ctx context.Context, p *ppb.OrgProfile) error {
	allDSARes, err := s.resolveDRP().ListDSA(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}

	if ok := isDSACodePresent(p.DsaCode, allDSARes); ok {
		return nil
	}

	busInfo := &ppb.BusinessInfo{}
	if p.GetBusinessInfo() != nil {
		busInfo = p.GetBusinessInfo()
	}

	addr := &ppb.Address{}
	if busInfo.GetAddress() != nil {
		addr = busInfo.GetAddress()
	}

	// TODO: Hoew to get these values?
	vatable, tin := "1", "1234567890"

	contact, address1 := "Unknown", "Unknown"
	if p.GetBusinessInfo().GetContactPerson() != "" {
		contact = p.GetBusinessInfo().GetContactPerson()
	}

	if addr.Address1 != "" {
		address1 = addr.Address1
	}

	// Not present then create DSA in perahub
	dsaRes, err := s.resolveDRP().CreateDSA(ctx, &dsa.CreateDSARequest{
		DsaCode:       p.DsaCode,
		DsaName:       busInfo.CompanyName,
		EmailAddress:  busInfo.CompanyEmail,
		Vatable:       vatable,
		Address:       address1,
		Tin:           tin,
		UpdatedBy:     userDRP,
		ContactPerson: contact,
		City:          addr.City,
		Province:      addr.State,
		Zipcode:       addr.PostalCode,
	})
	if err != nil {
		logging.WithError(err, s.log).Error("SyncDSACommissionConfigForRemittance create dsa perahub error")
		return err
	}

	s.log.WithField("dsaCode", dsaRes.DsaCode).
		Info("SyncDSACommissionConfigForRemittance dsa created in perahub")

	return nil
}

func (s *RemcoCommissionSvc) getDRPDSACommAndTierMap(ctx context.Context,
	dsaCode, orgID string,
) (map[string]revcomm.DSACommission, map[string]revcomm.DSACommissionTier, error) {
	log := s.log
	drpCommList, err := s.pf.GetRevenueSharingList(ctx, &revshar.GetRevenueSharingListRequest{
		OrgID:     orgID,
		RemitType: revshar.RemitType_REMITTANCE,
	})
	if err != nil {
		sErr := status.Convert(err)
		if !(sErr.Code() == codes.NotFound || sErr.Code() == http.StatusNotFound) {
			logging.WithError(err, log).WithField("orgID", orgID).
				Error("SyncDSACommissionConfigForRemittance get revenue sharing list error")
			return nil, nil, err
		}

		log.Info("SyncDSACommissionConfigForRemittance no revenue sharing config found in DRP")
	}

	commMap := map[string]revcomm.DSACommission{}
	tierMap := map[string]revcomm.DSACommissionTier{}
	if drpCommList == nil || len(drpCommList.Results) == 0 {
		log.Info("SyncDSACommissionConfigForRemittance drp commission list is empty")
		return commMap, tierMap, nil
	}

	for _, d := range drpCommList.Results {
		if d == nil {
			continue
		}

		tiersRes, err := s.pf.GetRevenueSharingTierList(ctx, &revshar.GetRevenueSharingTierListRequest{
			RevenueSharingID: d.ID,
		})
		if err != nil {
			sErr := status.Convert(err)
			if !(sErr.Code() == codes.NotFound || sErr.Code() == http.StatusNotFound) {
				s.log.WithError(err).Error("SyncDSACommissionConfigForRemittance list dsa commission perahub error")
				return nil, nil, err
			}
		}

		phComm := toPerahubDSACommission(d, dsaCode)
		if tiersRes == nil || len(tiersRes.Results) == 0 {
			commMap[generateDSACommKey(phComm, "0", "0")] = phComm
			continue
		}

		for _, t := range tiersRes.Results {
			if t == nil {
				continue
			}

			phComm.CommissionAmount = t.Amount
			key := generateDSACommKey(phComm, t.MinValue, t.MaxValue)
			commMap[key] = phComm
			tierMap[key] = toPerahubDSACommissionTier(*t)
		}
	}

	return commMap, tierMap, nil
}

func (s *RemcoCommissionSvc) getPerahubDSACommMap(ctx context.Context, dsaCode string) (map[string]revcomm.DSACommission, error) {
	phCommRes, err := s.resolveDRP().ListDSACommission(ctx, &emptypb.Empty{})
	if err != nil {
		sErr := status.Convert(err)
		if !(sErr.Code() == codes.NotFound || sErr.Code() == http.StatusNotFound) {
			s.log.WithError(err).Error("SyncDSACommissionConfigForRemittance list dsa commission perahub error")
			return nil, err
		}

		s.log.Info("SyncDSACommissionConfigForRemittance no revenue sharing config found in Perahub")
	}

	commMap := map[string]revcomm.DSACommission{}
	if phCommRes == nil || len(phCommRes.CommissionList) == 0 {
		return commMap, nil
	}

	for _, d := range phCommRes.CommissionList {
		if d == nil || d.DsaCode != dsaCode {
			continue
		}

		tier := s.geDSACommissionCommTier(ctx, d.TierID)
		commMap[generateDSACommKey(*d, tier.Minimum, tier.Minimum)] = *d
	}

	return commMap, nil
}

func (s *RemcoCommissionSvc) createPerahubDSACommissionRecords(ctx context.Context, createList []revcomm.DSACommission) {
	log := s.log
	for _, cm := range createList {
		req := &revcomm.CreateDSACommissionRequest{
			DsaCode:            cm.DsaCode,
			CommissionType:     cm.CommissionType,
			TierID:             cm.TierID,
			CommissionAmount:   cm.CommissionAmount,
			CommissionCurrency: cm.CommissionCurrency,
			UpdatedBy:          cm.UpdatedBy,
			EffectiveDate:      cm.EffectiveDate,
			TrxType:            cm.TrxType,
			RemitType:          cm.RemitType,
		}

		res, err := s.resolveDRP().CreateDSACommission(ctx, req)
		if err != nil {
			logging.WithError(err, log).WithField("request", req).Error("create perahub dsa commission error")
			continue
		}

		log.WithField("CommissionID", res.CommID).
			Info("perahub dsa commission created")
		log.WithField("response", res).
			Debug("perahub dsa commission created")
	}
}

func (s *RemcoCommissionSvc) deletePerahubDSACommissionRecords(ctx context.Context, deleteList []revcomm.DSACommission) {
	log := s.log
	for _, cm := range deleteList {
		req := &revcomm.DeleteDSACommissionRequest{
			CommID: cm.CommID,
		}

		_, err := s.resolveDRP().DeleteDSACommission(ctx, req)
		if err != nil {
			logging.WithError(err, log).WithField("request", req).Error("delete perahub dsa commission error")
			continue
		}

		log.WithField("CommissionID", cm.CommID).
			Info("perahub dsa commission deleted")
	}
}

func (s *RemcoCommissionSvc) getCreateAndDeletePerahubDSACommissionList(ctx context.Context,
	drpMap map[string]revcomm.DSACommission,
	phMap map[string]revcomm.DSACommission,
	drpCommTiersMap map[string]revcomm.DSACommissionTier,
) (createList []revcomm.DSACommission, deleteList []revcomm.DSACommission, err error) {
	log := s.log

	createList = make([]revcomm.DSACommission, 0, len(drpMap))
	for key, cm := range drpMap {
		if _, ok := phMap[key]; ok {
			// record is present in perahub
			continue
		}

		// If tier does not present, just create commission record
		tier, ok := drpCommTiersMap[key]
		if !ok {
			createList = append(createList, cm)
			continue
		}

		// Create tier
		tierRes, err := s.resolveDRP().CreateDSACommissionTier(ctx, &revcomm.CreateDSACommissionTierRequest{
			TierNo:    tier.TierNo,
			Minimum:   tier.Minimum,
			Maximum:   tier.Maximum,
			UpdatedBy: tier.UpdatedBy,
		})
		if err != nil {
			logging.WithError(err, log).Error("create perahub dsa commissioin tier error")
			return nil, nil, err
		}

		cm.TierID = tierRes.TierID
		createList = append(createList, cm)
	}

	deleteList = make([]revcomm.DSACommission, 0, len(phMap))
	for key, pcm := range phMap {
		if _, ok := drpMap[key]; ok {
			// record is present in perahub
			continue
		}

		deleteList = append(deleteList, pcm)
	}

	return createList, deleteList, nil
}

func (s *RemcoCommissionSvc) geDSACommissionCommTier(ctx context.Context, tierID uint32) revcomm.DSACommissionTier {
	log := s.log
	// need empty tier for generating key
	tier := revcomm.DSACommissionTier{
		TierID:  tierID,
		Minimum: "0",
		Maximum: "0",
	}

	if tierID == 0 {
		return tier
	}

	res, err := s.resolveDRP().GetDSACommissionTierByID(ctx, &revcomm.GetDSACommissionTierByIDRequest{TierID: tierID})
	if err != nil {
		logging.WithError(err, log).Error("SyncDSACommissionConfigForRemittance error getting dsa commission tier from perahub")
		return tier
	}

	if res == nil {
		log.WithField("tierID", tierID).Error("SyncDSACommissionConfigForRemittance empty dsa commission tier")
		return tier
	}

	return *res
}

func isDSACodePresent(dsaCode string, r *dsa.ListDSAResponse) bool {
	if r == nil || len(r.DSAList) == 0 {
		return false
	}

	for _, d := range r.DSAList {
		if d.DsaCode == dsaCode {
			return true
		}
	}

	return false
}

// generateDSACommKey generates key from dsa commission values
func generateDSACommKey(dc revcomm.DSACommission, min, max string) string {
	return fmt.Sprintf("DC_%s_%s_%s_%s_%s_%s",
		dc.CommissionType.String(), dc.TrxType, dc.RemitType, min, max, dc.CommissionAmount)
}

func drpTierTypeToCommissionType(tierType revshar.TierType) revcomm.CommissionType {
	if tierType == revshar.TierType_TIERPERCENTAGE {
		return revcomm.CommissionType_CommissionTypeRange
	}

	return revcomm.CommissionType_CommissionTypePercent
}

func toPerahubDSACommission(rs *revshar.RevenueSharing, dsaCode string) revcomm.DSACommission {
	trxType := util.PerahubTrxTypeDigital
	if rs.TransactionType == revshar.TransactionType_OTC {
		trxType = util.PerahubTrxTypeOTC
	}

	remitType := util.PerahubRemitTypeInbound
	if rs.BoundType == revshar.BoundType_OUTBOUND {
		remitType = util.PerahubRemitTypeOutbound
	}

	return revcomm.DSACommission{
		DsaCode:        dsaCode,
		CommissionType: drpTierTypeToCommissionType(rs.TierType),
		// perahub TierID is not present in drp
		TierID:             0,
		CommissionAmount:   rs.Amount,
		CommissionCurrency: commissionCurrency,
		// TODO: get user name from rs.UserID if needed
		UpdatedBy: userDRP,
		CreatedAt: rs.Created,
		UpdatedAt: rs.Updated,
		TrxType:   trxType,
		RemitType: remitType,
	}
}

func toPerahubDSACommissionTier(rs revshar.RevenueSharingTier) revcomm.DSACommissionTier {
	return revcomm.DSACommissionTier{
		TierNo:    rs.ID,
		Minimum:   rs.MinValue,
		Maximum:   rs.MaxValue,
		UpdatedBy: userDRP,
	}
}

package revenue_commission

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	revcom_int "brank.as/petnet/api/integration/revenue-commission"
	"brank.as/petnet/api/storage"
	revcom "brank.as/petnet/gunk/drp/v1/revenue-commission"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	rsr "brank.as/petnet/gunk/dsa/v2/revenuesharingreport"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SendTransactionReport ...
func (s *Svc) SendTransactionReport(ctx context.Context, req *revcom.SendTransactionReportRequest) (*emptypb.Empty, error) {
	log := logging.FromContext(ctx)
	dsaList, err := s.dsaStore.ListDSA(ctx)
	if err != nil {
		log.WithError(err).Error("failed to get dsa list")
		return nil, err
	}
	for _, dsaDtls := range dsaList {
		trxCount := s.getTransactionReport(ctx, dsaDtls.DsaCode)
		err := s.sendDataToPeraHub(ctx, dsaDtls.DsaCode, trxCount)
		if err != nil {
			log.WithError(err).Error("failed to send data to perahub")
			continue
		}
		err = s.sendDataToDbProfile(ctx, dsaDtls.DsaCode, trxCount)
		if err != nil {
			log.WithError(err).Error("failed to send data to db profile")
			continue
		}
	}
	return new(emptypb.Empty), nil
}

// SendTransactionReport ...
func (s *Svc) SyncTransactionReport(ctx context.Context) error {
	log := logging.FromContext(ctx)
	dsaList, err := s.dsaStore.ListDSA(ctx)
	if err != nil {
		log.WithError(err).Error("failed to get dsa list")
		return err
	}
	for _, dsaDtls := range dsaList {
		trxCount := s.getTransactionReport(ctx, dsaDtls.DsaCode)
		err := s.sendDataToPeraHub(ctx, dsaDtls.DsaCode, trxCount)
		if err != nil {
			log.WithError(err).Error("failed to send data to perahub")
			continue
		}
		err = s.sendDataToDbProfile(ctx, dsaDtls.DsaCode, trxCount)
		if err != nil {
			log.WithError(err).Error("failed to send data to db profile")
			continue
		}
	}
	return nil
}

// getTransactionReport ...
func (s *Svc) getTransactionReport(ctx context.Context, dsaCode string) int32 {
	log := logging.FromContext(ctx).WithField("method", "getTransactionReport")
	frm, untl := convertYearMonth()
	pf, err := s.profileService.GetProfileByDsaCode(ctx, &ppb.GetProfileByDsaCodeRequest{DsaCode: dsaCode})
	if err != nil {
		log.Error("failed to get profile by dsa code", err)
		return 0
	}

	str, err := s.store.GetTransactionReport(ctx, &storage.LRHFilter{
		DsaID:     pf.GetProfile().GetOrgID(), // should be dsa id
		TxnStatus: string(storage.SuccessStatus),
		TxnStep:   string(storage.ConfirmStep),
		From:      frm,
		Until:     untl,
	})
	if err != nil {
		log.Error("failed to get transaction report", err)
		return 0
	}
	return int32(str.RemitTransactionCount)
}

// sendDataToPeraHub ...
func (s *Svc) sendDataToPeraHub(ctx context.Context, dsaCode string, trxCount int32) error {
	_, err := s.commissionFeeStore.CreateTransactionCount(ctx, &revcom_int.SaveTransactionCountRequest{
		DsaCode:           dsaCode,
		YearMonth:         getYearMonth(),
		RemittanceCount:   json.Number(fmt.Sprintf("%d", trxCount)),
		CiCoCount:         json.Number("0"), // TODO: Add CashInCashOut Count
		BillsPaymentCount: json.Number("0"), // TODO: Add BillsPayment Count
		InsuranceCount:    json.Number("0"), // TODO: Add Insurance Count
		UpdatedBy:         "Admin",
		DsaCommission:     json.Number("0"), // TODO: Add Dsa Commission
		DsaCommissionType: "0",              // TODO: Add Commission Type
	})
	if err != nil {
		return err
	}
	return nil
}

// sendDataToDbProfile ...
func (s *Svc) sendDataToDbProfile(ctx context.Context, dsaCode string, trxCount int32) error {
	profile, err := s.profileService.GetProfileByDsaCode(ctx, &ppb.GetProfileByDsaCodeRequest{DsaCode: dsaCode})
	if err != nil {
		return err
	}
	orgId := profile.GetProfile().GetOrgID()

	if orgId == "" {
		return errors.New("dsa code not found" + dsaCode)
	}
	_, err = s.revenueReport.CreateRevenueSharingReport(ctx, &rsr.CreateRevenueSharingReportRequest{
		OrgID:             orgId,
		Created:           timestamppb.Now(),
		RemittanceCount:   trxCount,
		CicoCount:         0,   // TODO: Add CashInCashOut Count
		BillsPaymentCount: 0,   // TODO: Add BillsPayment Count
		InsuranceCount:    0,   // TODO: Add Insurance Count
		DsaCommission:     "0", // TODO: Add Dsa Commission
		CommissionType:    "0", // TODO: Add Commission Type
		DsaCode:           dsaCode,
		YearMonth:         getYearMonth(),
	})
	if err != nil {
		return err
	}
	return nil
}

func getYearMonth() string {
	t := time.Now().AddDate(0, -1, 0)
	year, month, _ := t.Date()
	return fmt.Sprintf("%d%02d", year, month)
}

// beginningOfMonth ...
func beginningOfMonth(date time.Time) time.Time {
	return date.AddDate(0, -1, -date.Day()+1)
}

// endOfMonth ...
func endOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day())
}

// convertYearMonth ...
func convertYearMonth() (time.Time, time.Time) {
	layout := "2006-1-02"
	fromDate := time.Now()
	y := fromDate.Year()
	m := fromDate.Month()
	dt := fmt.Sprintf("%d-%d-%s", y, m, "01")
	tm, _ := time.Parse(layout, dt)
	bt := beginningOfMonth(tm)
	et := endOfMonth(tm)
	return bt, et
}

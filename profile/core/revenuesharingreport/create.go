package revenuesharingreport

import (
	"context"

	rsp "brank.as/petnet/gunk/dsa/v2/revenuesharingreport"
	"brank.as/petnet/profile/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) CreateRevenueSharingReport(ctx context.Context, req *rsp.CreateRevenueSharingReportRequest) (*rsp.CreateRevenueSharingReportResponse, error) {
	res, err := s.st.CreateRevenueSharingReport(ctx, storage.RevenueSharingReport{
		OrgID:             req.GetOrgID(),
		DsaCode:           req.GetDsaCode(),
		YearMonth:         req.GetYearMonth(),
		RemittanceCount:   int(req.GetRemittanceCount()),
		CicoCount:         int(req.GetCicoCount()),
		BillsPaymentCount: int(req.GetBillsPaymentCount()),
		InsuranceCount:    int(req.GetInsuranceCount()),
		DsaCommission:     req.GetDsaCommission(),
		CommissionType:    req.GetCommissionType(),
		Status:            int(req.GetStatus()),
		Created:           req.GetCreated().AsTime(),
	})
	if err != nil {
		return nil, err
	}
	return &rsp.CreateRevenueSharingReportResponse{
		ID:                res.ID,
		OrgID:             res.OrgID,
		DsaCode:           res.DsaCode,
		YearMonth:         res.YearMonth,
		Status:            int32(res.Status),
		Created:           timestamppb.New(res.Created),
		Updated:           timestamppb.New(res.Updated),
		RemittanceCount:   int32(res.RemittanceCount),
		CicoCount:         int32(res.CicoCount),
		BillsPaymentCount: int32(res.BillsPaymentCount),
		InsuranceCount:    int32(res.InsuranceCount),
		DsaCommission:     res.DsaCommission,
		CommissionType:    res.CommissionType,
	}, nil
}

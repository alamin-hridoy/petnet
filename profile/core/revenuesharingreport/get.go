package revenuesharingreport

import (
	"context"

	rsp "brank.as/petnet/gunk/dsa/v2/revenuesharingreport"
	"brank.as/petnet/profile/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) GetRevenueSharingReportList(ctx context.Context, req *rsp.GetRevenueSharingReportListRequest) (*rsp.GetRevenueSharingReportListResponse, error) {
	um, err := s.st.GetRevenueSharingReportList(ctx, storage.RevenueSharingReport{
		OrgID:        req.GetOrgID(),
		SortByColumn: storage.RevenueSharingReportColumn(req.GetSortByColumn().String()),
		SortOrder:    storage.SortOrder(req.GetSortOrder().String()),
		Limit:        int(req.GetLimit()),
		Offset:       int(req.GetOffset()),
	})
	if err != nil {
		return nil, err
	}
	var trns []*rsp.RevenueSharingReport
	for _, v := range um {
		trns = append(trns, &rsp.RevenueSharingReport{
			ID:                v.ID,
			OrgID:             v.OrgID,
			DsaCode:           v.DsaCode,
			YearMonth:         v.YearMonth,
			Status:            int32(v.Status),
			Created:           timestamppb.New(v.Created),
			Updated:           timestamppb.New(v.Updated),
			Count:             int32(v.Count),
			RemittanceCount:   int32(v.RemittanceCount),
			CicoCount:         int32(v.CicoCount),
			BillsPaymentCount: int32(v.BillsPaymentCount),
			InsuranceCount:    int32(v.InsuranceCount),
			DsaCommission:     v.DsaCommission,
			CommissionType:    v.CommissionType,
		})
	}

	return &rsp.GetRevenueSharingReportListResponse{
		Results: trns,
	}, nil
}

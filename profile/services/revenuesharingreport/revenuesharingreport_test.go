package revenuesharingreport

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"brank.as/petnet/profile/core/revenuesharingreport"
	"brank.as/petnet/profile/storage/postgres"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	tpb "brank.as/petnet/gunk/dsa/v2/revenuesharingreport"
)

func TestRevenueSharingReport(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	st, cleanup := postgres.NewTestStorage(os.Getenv("DATABASE_CONNECTION"), filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(cleanup)
	h := New(revenuesharingreport.New(st))
	oid := "20000000-0000-0000-0000-000000000000"
	oid2 := "22000000-0000-0000-0000-000000000000"
	oid3 := "21000000-0000-0000-0000-000000000000"
	want := &tpb.GetRevenueSharingReportListRequest{
		OrgID:  oid,
		Status: 0,
	}
	_, err := h.CreateRevenueSharingReport(ctx, &tpb.CreateRevenueSharingReportRequest{
		OrgID:             want.OrgID,
		DsaCode:           "test",
		YearMonth:         "2022",
		Status:            1,
		RemittanceCount:   1,
		CicoCount:         1,
		BillsPaymentCount: 1,
		InsuranceCount:    1,
		DsaCommission:     "test",
		CommissionType:    "test",
	})
	if err != nil {
		t.Fatal("h.CreateRevenueSharingReport: ", err)
	}
	want2 := []*tpb.RevenueSharingReport{
		{
			OrgID:             oid,
			DsaCode:           "test",
			YearMonth:         "2022",
			Status:            1,
			Created:           timestamppb.Now(),
			Updated:           timestamppb.Now(),
			Count:             1,
			RemittanceCount:   1,
			CicoCount:         1,
			BillsPaymentCount: 1,
			InsuranceCount:    1,
			DsaCommission:     "test",
			CommissionType:    "test",
		},
		{
			OrgID:             oid2,
			DsaCode:           "test2",
			YearMonth:         "2021",
			Status:            1,
			Created:           timestamppb.Now(),
			Updated:           timestamppb.Now(),
			Count:             1,
			RemittanceCount:   1,
			CicoCount:         1,
			BillsPaymentCount: 1,
			InsuranceCount:    1,
			DsaCommission:     "test2",
			CommissionType:    "test2",
		},
		{
			OrgID:             oid3,
			DsaCode:           "test3",
			YearMonth:         "2024",
			Status:            2,
			Created:           timestamppb.Now(),
			Updated:           timestamppb.Now(),
			Count:             2,
			RemittanceCount:   2,
			CicoCount:         2,
			BillsPaymentCount: 2,
			InsuranceCount:    2,
			DsaCommission:     "test3",
			CommissionType:    "test3",
		},
	}
	want3 := &tpb.GetRevenueSharingReportListResponse{
		Results: []*tpb.RevenueSharingReport{want2[0]},
	}
	oo := cmp.Options{cmpopts.IgnoreFields(tpb.RevenueSharingReport{}, "ID", "Created", "Updated", "Count"), cmpopts.IgnoreUnexported(tpb.RevenueSharingReport{}, timestamppb.Timestamp{})}
	_, err = h.CreateRevenueSharingReport(ctx, &tpb.CreateRevenueSharingReportRequest{
		OrgID:             want2[1].OrgID,
		DsaCode:           want2[1].DsaCode,
		YearMonth:         want2[1].YearMonth,
		Status:            1,
		Created:           timestamppb.Now(),
		RemittanceCount:   want2[1].RemittanceCount,
		CicoCount:         want2[1].CicoCount,
		BillsPaymentCount: want2[1].BillsPaymentCount,
		InsuranceCount:    want2[1].InsuranceCount,
		DsaCommission:     want2[1].DsaCommission,
		CommissionType:    want2[1].CommissionType,
	})
	if err != nil {
		t.Fatal("Create Revenue Sharing Report: ", err)
	}
	res2, err := h.GetRevenueSharingReportList(ctx, &tpb.GetRevenueSharingReportListRequest{
		OrgID:             want2[0].OrgID,
		DsaCode:           want2[0].DsaCode,
		YearMonth:         want2[0].YearMonth,
		Status:            want2[0].Status,
		Created:           timestamppb.Now(),
		Updated:           timestamppb.Now(),
		RemittanceCount:   want2[0].RemittanceCount,
		CicoCount:         want2[0].CicoCount,
		BillsPaymentCount: want2[0].BillsPaymentCount,
		InsuranceCount:    want2[0].InsuranceCount,
		DsaCommission:     want2[0].DsaCommission,
		CommissionType:    want2[0].CommissionType,
	})
	if err != nil {
		t.Error("GetRevenueSharingReportList: ", err)
	}
	if !cmp.Equal(want2[0], res2.Results[0], oo) {
		t.Error("(-want +got): ", cmp.Diff(want2[0], res2.Results[0], oo))
	}
	want3.Results = []*tpb.RevenueSharingReport{want2[1]}
	if _, err = h.CreateRevenueSharingReport(ctx, &tpb.CreateRevenueSharingReportRequest{
		OrgID:             want2[2].OrgID,
		DsaCode:           want2[2].DsaCode,
		YearMonth:         want2[2].YearMonth,
		Status:            want2[2].Status,
		Created:           timestamppb.Now(),
		RemittanceCount:   want2[2].RemittanceCount,
		CicoCount:         want2[2].CicoCount,
		BillsPaymentCount: want2[2].BillsPaymentCount,
		InsuranceCount:    want2[2].InsuranceCount,
		DsaCommission:     want2[2].DsaCommission,
		CommissionType:    want2[2].CommissionType,
	}); err != nil {
		t.Fatal("Create Revenue Sharing Report: ", err)
	}
	res2, err = h.GetRevenueSharingReportList(ctx, &tpb.GetRevenueSharingReportListRequest{
		OrgID:             want2[1].OrgID,
		DsaCode:           want2[1].DsaCode,
		YearMonth:         want2[1].YearMonth,
		Status:            want2[1].Status,
		Created:           timestamppb.Now(),
		Updated:           timestamppb.Now(),
		RemittanceCount:   want2[1].RemittanceCount,
		CicoCount:         want2[1].CicoCount,
		BillsPaymentCount: want2[1].BillsPaymentCount,
		InsuranceCount:    want2[1].InsuranceCount,
		DsaCommission:     want2[1].DsaCommission,
		CommissionType:    want2[1].CommissionType,
	})
	if err != nil {
		t.Error("GetRevenueSharingReportList: ", err)
	}
	if !cmp.Equal(want3.Results[0], res2.Results[0], oo) {
		t.Error("(-want +got): ", cmp.Diff(want3.Results[0], res2.Results[0], oo))
	}

	_, err = h.UpdateRevenueSharingReport(ctx, &tpb.UpdateRevenueSharingReportRequest{
		OrgID:             want2[1].OrgID,
		DsaCode:           want2[1].DsaCode,
		YearMonth:         want2[1].YearMonth,
		Status:            want2[1].Status,
		Created:           timestamppb.Now(),
		RemittanceCount:   want2[1].RemittanceCount,
		CicoCount:         want2[1].CicoCount,
		BillsPaymentCount: want2[1].BillsPaymentCount,
		InsuranceCount:    want2[1].InsuranceCount,
		DsaCommission:     want2[1].DsaCommission,
		CommissionType:    want2[1].CommissionType,
	})
	if err != nil {
		t.Error("UpdateRevenueSharingReport: ", err)
	}
}

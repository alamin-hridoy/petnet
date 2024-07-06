package postgres_test

import (
	"context"
	"sort"
	"testing"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sirupsen/logrus"
)

func TestRevenueSharingReport(t *testing.T) {
	ts := newTestStorage(t)

	oid := "10000000-0000-0000-0000-000000000000"
	oid2 := "20000000-0000-0000-0000-000000000000"
	want := []storage.RevenueSharingReport{
		{
			OrgID:             oid,
			DsaCode:           "test",
			YearMonth:         "2022",
			RemittanceCount:   1,
			CicoCount:         1,
			BillsPaymentCount: 1,
			InsuranceCount:    1,
			DsaCommission:     "test",
			CommissionType:    "test",
			Status:            0,
		},
		{
			OrgID:             oid2,
			DsaCode:           "test2",
			YearMonth:         "2021",
			RemittanceCount:   2,
			CicoCount:         2,
			BillsPaymentCount: 2,
			InsuranceCount:    2,
			DsaCommission:     "test2",
			CommissionType:    "test2",
			Status:            1,
		},
	}

	logr := logging.NewLogger(nil)
	logr.SetFormatter(&logrus.JSONFormatter{})
	ctx := logging.WithLogger(context.TODO(), logr)

	_, err := ts.CreateRevenueSharingReport(ctx, want[0])
	if err != nil {
		t.Fatal(err)
	}

	_, err = ts.CreateRevenueSharingReport(ctx, want[1])
	if err != nil {
		t.Fatal(err)
	}

	tOps := []cmp.Option{
		cmpopts.IgnoreFields(storage.RevenueSharingReport{}, "ID", "Created", "Updated", "Total"),
	}

	gotlist, err := ts.GetRevenueSharingReportList(ctx, want[0])
	if err != nil {
		t.Fatal(err)
	}

	upid, err := ts.UpdateRevenueSharingReport(ctx, want[0])
	if err != nil {
		t.Fatal(err)
	}
	if upid.ID != gotlist[0].ID {
		t.Error("id mismatch")
	}

	sort.Slice(want, func(i, j int) bool {
		return want[i].OrgID < want[j].OrgID
	})
	sort.Slice(gotlist, func(i, j int) bool {
		return gotlist[i].OrgID < gotlist[j].OrgID
	})
	for i, pf := range gotlist {
		if !cmp.Equal(want[i], pf, tOps...) {
			t.Error("(-want +got): ", cmp.Diff(want[i], pf, tOps...))
		}
		if pf.OrgID == "" {
			t.Error("org id should not be empty")
		}
		if pf.Created.IsZero() {
			t.Error("created shouldn't be empty")
		}
	}
}

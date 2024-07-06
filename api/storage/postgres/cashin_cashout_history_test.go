package postgres

import (
	"context"
	"sort"
	"testing"
	"time"

	"brank.as/petnet/api/storage"
	"brank.as/petnet/serviceutil/logging"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sirupsen/logrus"
)

func TestCashInCashOutHistory(t *testing.T) {
	ts := newTestStorage(t)
	oid := "10000000-0000-0000-0000-000000000000"
	oid2 := "20000000-0000-0000-0000-000000000000"
	want := []storage.CashInCashOutHistory{
		{
			OrgID:              oid,
			PartnerCode:        "DSA",
			SvcProvider:        "GCASH_CASHIN",
			Provider:           "GCASH",
			TrxType:            "Cash In",
			ReferenceNumber:    "09654767706",
			PetnetTrackingNo:   "3115bc3f587d747cf8f5",
			ProviderTrackingNo: "7000001521345",
			PrincipalAmount:    100,
			TxnStatus:          "SUCCESS",
			Charges:            10,
			TotalAmount:        110,
			TrxDate:            time.Time{},
		},
		{
			OrgID:              oid2,
			PartnerCode:        "DSA2",
			SvcProvider:        "GCASH_CASHIN",
			Provider:           "GCASH2",
			TrxType:            "Cash In2",
			ReferenceNumber:    "09654767708",
			PetnetTrackingNo:   "3115bc3f587d747cf8f9",
			ProviderTrackingNo: "7000001521347",
			TxnStatus:          "SUCCESS",
			PrincipalAmount:    200,
			Charges:            20,
			TotalAmount:        220,
			TrxDate:            time.Time{},
		},
	}
	logr := logging.NewLogger(nil)
	logr.SetFormatter(&logrus.JSONFormatter{})
	ctx := logging.WithLogger(context.TODO(), logr)
	_, err := ts.CreateCICOHistory(ctx, want[0])
	if err != nil {
		t.Fatal(err)
	}
	_, err = ts.CreateCICOHistory(ctx, want[1])
	if err != nil {
		t.Fatal(err)
	}
	tOps := []cmp.Option{
		cmpopts.IgnoreFields(storage.CashInCashOutHistory{}, "ID", "Created", "Updated", "Total", "Details", "ProviderTrackingNo"),
	}
	_, err = ts.GetCICOHistory(ctx, want[0].OrgID)
	if err != nil {
		t.Fatalf("Get CICO History = got error %v, want nil", err)
	}
	gotlist, err := ts.ListCICOHistory(ctx, want[0])
	if err != nil {
		t.Fatal(err)
	}
	upData := storage.CashInCashOutHistory{
		OrgID: want[0].OrgID,
	}
	upid, err := ts.UpdateCICOHistory(ctx, upData)
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

	}
}

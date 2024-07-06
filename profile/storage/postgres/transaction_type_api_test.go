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

func TestTransactionTypeApi(t *testing.T) {
	ts := newTestStorage(t)

	oid := "10000000-0000-0000-0000-000000000000"
	oid2 := "20000000-0000-0000-0000-000000000000"
	uid := "11000000-0000-0000-0000-000000000000"
	uid2 := "22000000-0000-0000-0000-000000000000"
	want := []*storage.ApiKeyTransactionType{
		{
			OrgID:           oid,
			UserID:          uid,
			ClientID:        "12345",
			Environment:     "Sandbox",
			TransactionType: "Digital",
		},
		{
			OrgID:           oid2,
			UserID:          uid2,
			ClientID:        "123456",
			Environment:     "Production",
			TransactionType: "OTC",
		},
	}

	logr := logging.NewLogger(nil)
	logr.SetFormatter(&logrus.JSONFormatter{})
	ctx := logging.WithLogger(context.TODO(), logr)

	pid, err := ts.InsertApiKeyTransactionType(ctx, want[0])
	if err != nil {
		t.Fatal(err)
	}
	if pid == "" {
		t.Error("id should not be empty")
	}

	pid2, err := ts.InsertApiKeyTransactionType(ctx, want[1])
	if err != nil {
		t.Fatal(err)
	}
	if pid2 == "" {
		t.Error("id should not be empty")
	}

	got, err := ts.GetAPITransactionType(ctx, &storage.ApiKeyTransactionType{
		UserID:          uid,
		OrgID:           oid,
		Environment:     "Sandbox",
		TransactionType: "Digital",
	})
	if err != nil {
		t.Fatal(err)
	}

	tOps := []cmp.Option{
		cmpopts.IgnoreFields(storage.ApiKeyTransactionType{}, "ID", "Created"),
	}
	if !cmp.Equal(want[0], got, tOps...) {
		t.Error("(-want +got): ", cmp.Diff(&want[0], got, tOps...))
	}
	if got.OrgID == "" {
		t.Error("org id should not be empty")
	}
	if got.UserID == "" {
		t.Error("user id should not be empty")
	}
	if got.Environment == "" {
		t.Error("environment should not be empty")
	}
	if got.TransactionType == "" {
		t.Error("transaction type should not be empty")
	}
	if got.Created.IsZero() {
		t.Error("created shouldn't be empty")
	}

	gotlist, err := ts.ListUserAPIKeyTransactionType(ctx, oid, uid)
	if err != nil {
		t.Fatal(err)
	}

	sort.Slice(want, func(i, j int) bool {
		return want[i].OrgID < want[j].OrgID
	})
	sort.Slice(gotlist, func(i, j int) bool {
		return gotlist[i].OrgID < gotlist[j].OrgID
	})
	for i, pf := range gotlist {
		if !cmp.Equal(want[i], &pf, tOps...) {
			t.Error("(-want +got): ", cmp.Diff(want[i], pf, tOps...))
		}
		if pf.OrgID == "" {
			t.Error("org id should not be empty")
		}
		if pf.UserID == "" {
			t.Error("user id should not be empty")
		}
		if pf.Created.IsZero() {
			t.Error("created shouldn't be empty")
		}
	}
}

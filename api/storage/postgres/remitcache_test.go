package postgres

import (
	"context"
	"testing"

	"brank.as/petnet/api/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestCreateRemitCache(t *testing.T) {
	ts := newTestStorage(t)

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()
		in := &storage.RemitCache{
			DsaID:             uuid.NewString(),
			UserID:            "test-user-id",
			RemcoID:           "test-remco-id",
			RemcoMemberID:     "test-remco-member-id",
			RemcoControlNo:    "test-remco-reference",
			RemcoAltControlNo: "test-remco-alt-reference",
			Step:              "test-step",
			Remit:             []byte{},
		}
		r, err := ts.CreateRemitCache(context.TODO(), *in)
		if err != nil {
			t.Fatalf("CreateRemitCache() = got error %v, want nil", err)
		}
		if r.TxnID == "" {
			t.Fatal("CreateRemitCache() = returned empty ID")
		}
	})
}

func TestUpdateRemitCache(t *testing.T) {
	ts := newTestStorage(t)

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()
		in := &storage.RemitCache{
			DsaID:             uuid.NewString(),
			UserID:            "test-user-id",
			RemcoID:           "test-remco-id",
			RemcoMemberID:     "test-remco-member-id",
			RemcoControlNo:    "test-remco-reference",
			RemcoAltControlNo: "test-remco-alt-reference",
			Step:              "test-step",
			Remit:             []byte{},
		}
		r, err := ts.CreateRemitCache(context.TODO(), *in)
		if err != nil {
			t.Fatalf("CreateRemitCache() = got error %v, want nil", err)
		}
		if r.TxnID == "" {
			t.Fatal("CreateRemitCache() = returned empty ID")
		}

		r.Step = "test-status-update"
		ru, err := ts.UpdateRemitCache(context.TODO(), *r)
		if err != nil {
			t.Fatalf("UpdateRemitCache() = got error %v, want nil", err)
		}

		if !cmp.Equal(r.Step, ru.Step) {
			t.Fatal(cmp.Diff(r.Step, ru.Step))
		}
	})
}

func TestGetRemitCache(t *testing.T) {
	ts := newTestStorage(t)

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()
		in := &storage.RemitCache{
			DsaID:             uuid.NewString(),
			UserID:            "test-user-id",
			RemcoID:           "test-remco-id",
			RemcoMemberID:     "test-remco-member-id",
			RemcoControlNo:    "test-remco-reference",
			RemcoAltControlNo: "test-remco-alt-reference",
			Step:              "test-step",
			Remit:             []byte{},
		}
		r, err := ts.CreateRemitCache(context.TODO(), *in)
		if err != nil {
			t.Fatalf("CreateRemitCache() = got error %v, want nil", err)
		}
		if r.TxnID == "" {
			t.Fatal("CreateRemitCache() = returned empty ID")
		}

		gr, err := ts.GetRemitCache(context.TODO(), r.TxnID)
		if err != nil {
			t.Fatalf("GetRemitCache() = got error %v, want nil", err)
		}

		opt := cmpopts.IgnoreFields(storage.RemitCache{}, "Remit")

		if !cmp.Equal(*r, *gr, opt) {
			t.Fatal(cmp.Diff(*r, *gr, opt))
		}
	})
}

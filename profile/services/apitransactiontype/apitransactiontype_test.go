package apitransactiontype

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"brank.as/petnet/profile/core/apitransactiontype"
	"brank.as/petnet/profile/storage/postgres"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	tpb "brank.as/petnet/gunk/dsa/v2/transactiontype"
)

func TestApiTransactionType(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	st, cleanup := postgres.NewTestStorage(os.Getenv("DATABASE_CONNECTION"), filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(cleanup)
	h := New(apitransactiontype.New(st))
	transtyp := tpb.TransactionType_DIGITAL
	transtypOTC := tpb.TransactionType_OTC
	oid := "20000000-0000-0000-0000-000000000000"
	uid := uuid.New().String()
	oid2 := "22000000-0000-0000-0000-000000000000"
	uid2 := uuid.New().String()
	want := &tpb.ApiKeyTransactionType{
		UserID:          uid,
		OrgID:           oid,
		ClientID:        "12345",
		Environment:     "Sandbox",
		TransactionType: transtyp,
	}
	rr, err := h.CreateApiKeyTransactionType(ctx, &tpb.CreateApiKeyTransactionTypeRequest{
		UserID:          uid,
		OrgID:           want.OrgID,
		ClientID:        want.ClientID,
		Environment:     want.Environment,
		TransactionType: want.TransactionType,
	})
	if err != nil {
		t.Fatal("h.CreateApiKeyTransactionType: ", err)
	}

	if want.UserID != rr.ApiTransactionTypes.UserID {
		t.Fatal("user id should be equal")
	}

	o := cmp.Options{
		cmpopts.IgnoreUnexported(
			tpb.ApiKeyTransactionType{},
		),
		cmpopts.IgnoreFields(tpb.ApiKeyTransactionType{}, "ID"),
	}
	got, err := h.GetAPITransactionType(ctx, &tpb.GetAPITransactionTypeRequest{
		UserID:          uid,
		OrgID:           oid,
		Environment:     "Sandbox",
		TransactionType: transtyp,
	})
	if err != nil {
		t.Fatal("get api transaction type data: ", err)
	}

	if !cmp.Equal(want, got, o...) {
		t.Error("(-want +got): ", cmp.Diff(want, got, o...))
	}

	want2 := []*tpb.ApiKeyTransactionType{
		{
			UserID:          uid,
			OrgID:           oid,
			ClientID:        "12345",
			Environment:     "Sandbox",
			TransactionType: transtyp,
		},
		{
			UserID:          uid2,
			OrgID:           oid2,
			ClientID:        "123456",
			Environment:     "Production",
			TransactionType: transtypOTC,
		},
	}

	want3 := &tpb.ListUserAPIKeyTransactionTypeResponse{
		ApiTransactionType: []*tpb.ApiKeyTransactionType{want2[0]},
	}
	oo := cmp.Options{
		cmpopts.IgnoreFields(tpb.ApiKeyTransactionType{}, "ID", "Created"),
		cmpopts.IgnoreUnexported(tpb.ApiKeyTransactionType{}),
	}
	res, err := h.CreateApiKeyTransactionType(ctx, &tpb.CreateApiKeyTransactionTypeRequest{
		UserID:          want2[0].UserID,
		OrgID:           want2[0].OrgID,
		ClientID:        want2[0].ClientID,
		Environment:     want2[0].Environment,
		TransactionType: want2[0].TransactionType,
	})
	if err != nil {
		t.Fatal("create branch: ", err)
	}
	res2, err := h.ListUserAPIKeyTransactionType(ctx, &tpb.ListUserAPIKeyTransactionTypeRequest{
		OrgID:  want2[0].OrgID,
		UserID: want2[0].UserID,
	})
	if err != nil {
		t.Error("ListUserAPIKeyTransactionType: ", err)
	}
	if !cmp.Equal(want2[0], res2.ApiTransactionType[0], oo) {
		t.Error("(-want +got): ", cmp.Diff(want2[0], res2.ApiTransactionType[0], oo))
	}
	want2[1].UserID = res.ApiTransactionTypes.UserID
	want3.ApiTransactionType = []*tpb.ApiKeyTransactionType{want2[1]}
	if _, err = h.CreateApiKeyTransactionType(ctx, &tpb.CreateApiKeyTransactionTypeRequest{
		UserID:          want2[1].UserID,
		OrgID:           want2[1].OrgID,
		ClientID:        want2[1].ClientID,
		Environment:     want2[1].Environment,
		TransactionType: want2[1].TransactionType,
	}); err != nil {
		t.Fatal("create api key transaction type: ", err)
	}
	res2, err = h.ListUserAPIKeyTransactionType(ctx, &tpb.ListUserAPIKeyTransactionTypeRequest{
		OrgID:  want2[1].OrgID,
		UserID: want2[1].UserID,
	})
	if err != nil {
		t.Error("ListUserAPIKeyTransactionType: ", err)
	}
	if !cmp.Equal(want3.ApiTransactionType[0], res2.ApiTransactionType[0], oo) {
		t.Error("(-want +got): ", cmp.Diff(want3.ApiTransactionType[0], res2.ApiTransactionType[0], oo))
	}
}

package fees

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus/hooks/test"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/profile/storage/postgres"

	fpb "brank.as/petnet/gunk/dsa/v2/fees"
	fc "brank.as/petnet/profile/core/fees"
	"google.golang.org/protobuf/types/known/timestamppb"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func TestFees(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	oid := uuid.New().String()
	fs := []*fpb.Fee{
		{
			OrgID: oid,
			Type:  1,
			Schedule: &fpb.Schedule{
				Status:    fpb.FeeStatus_Disabled,
				StartDate: tspb.Now(),
				EndDate:   tspb.Now(),
			},
		},
		{
			OrgID: oid,
			Type:  1,
			Schedule: &fpb.Schedule{
				Status:    fpb.FeeStatus_Disabled,
				StartDate: tspb.Now(),
				EndDate:   tspb.Now(),
			},
		},
	}
	test.NewNullLogger()
	st, cleanup := postgres.NewTestStorage(os.Getenv("DATABASE_CONNECTION"), filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(cleanup)

	_, err := st.CreateOrgProfile(ctx, &storage.OrgProfile{
		OrgID:  oid,
		UserID: uuid.NewString(),
	})
	if err != nil {
		t.Fatal(err)
	}

	want := &fpb.ListFeesResponse{
		Fees:  []*fpb.Fee{fs[0]},
		Total: 1,
	}
	s := New(fc.New(st))
	o := cmp.Options{
		cmpopts.IgnoreFields(fpb.Fee{}, "ID", "Rates"),
		cmpopts.IgnoreFields(timestamppb.Timestamp{}, "Nanos", "Seconds"),
		cmpopts.IgnoreUnexported(
			fpb.Fee{}, fpb.ListFeesResponse{}, fpb.Rate{},
			fpb.Amount{}, fpb.Schedule{}, timestamppb.Timestamp{},
		),
	}
	res, err := s.UpsertFee(ctx, &fpb.UpsertFeeRequest{Fee: fs[0]})
	if err != nil {
		t.Fatal("create fee: ", err)
	}
	res2, err := s.ListFees(ctx, &fpb.ListFeesRequest{
		OrgID: oid,
		Type:  fpb.FeeType_TypeFee.String(),
	})
	if err != nil {
		t.Error("ListFees: ", err)
	}
	if !cmp.Equal(want, res2, o) {
		t.Error("ListFees (-want +got): ", cmp.Diff(want, res2, o))
	}

	fs[1].ID = res.ID
	want.Fees = []*fpb.Fee{fs[1]}
	if _, err = s.UpsertFee(ctx, &fpb.UpsertFeeRequest{Fee: fs[1]}); err != nil {
		t.Fatal("update fee: ", err)
	}

	res2, err = s.ListFees(ctx, &fpb.ListFeesRequest{
		OrgID: oid,
		Type:  fpb.FeeType_TypeFee.String(),
	})
	if err != nil {
		t.Error("ListFees: ", err)
	}
	if !cmp.Equal(want, res2, o) {
		t.Error("ListFees (-want +got): ", cmp.Diff(want, res2, o))
	}
}

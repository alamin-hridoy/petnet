package branch

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

	bpb "brank.as/petnet/gunk/dsa/v2/branch"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	bc "brank.as/petnet/profile/core/branch"
)

func TestBranch(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	oid := uuid.New().String()
	bs := []*bpb.Branch{
		{
			OrgID: oid,
			Title: "title",
			Address: &ppb.Address{
				Address1:   "addr1",
				City:       "city",
				State:      "state",
				PostalCode: "12345",
			},
			PhoneNumber:   "12345",
			FaxNumber:     "12345",
			ContactPerson: "contact",
		},
		{
			OrgID: oid,
			Title: "title2",
			Address: &ppb.Address{
				Address1:   "addr12",
				City:       "city2",
				State:      "state2",
				PostalCode: "123452",
			},
			PhoneNumber:   "123452",
			FaxNumber:     "123452",
			ContactPerson: "contact2",
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

	want := &bpb.ListBranchesResponse{
		Branches: []*bpb.Branch{bs[0]},
		Total:    1,
	}
	s := New(bc.New(st))
	o := cmp.Options{
		cmpopts.IgnoreFields(bpb.Branch{}, "ID", "Created", "Updated"),
		cmpopts.IgnoreUnexported(
			bpb.Branch{}, bpb.ListBranchesResponse{}, ppb.Address{},
		),
	}
	res, err := s.UpsertBranch(ctx, &bpb.UpsertBranchRequest{Branch: bs[0]})
	if err != nil {
		t.Fatal("create branch: ", err)
	}
	res2, err := s.ListBranches(ctx, &bpb.ListBranchesRequest{OrgID: oid})
	if err != nil {
		t.Error("ListBranches: ", err)
	}
	if !cmp.Equal(want, res2, o) {
		t.Error("(-want +got): ", cmp.Diff(want, res2, o))
	}

	bs[1].ID = res.ID
	want.Branches = []*bpb.Branch{bs[1]}
	if _, err = s.UpsertBranch(ctx, &bpb.UpsertBranchRequest{Branch: bs[1]}); err != nil {
		t.Fatal("update branch: ", err)
	}
	res2, err = s.ListBranches(ctx, &bpb.ListBranchesRequest{OrgID: oid})
	if err != nil {
		t.Error("ListBranches: ", err)
	}
	if !cmp.Equal(want, res2, o) {
		t.Error("(-want +got): ", cmp.Diff(want, res2, o))
	}
}

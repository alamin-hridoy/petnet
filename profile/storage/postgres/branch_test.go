package postgres_test

import (
	"context"
	"testing"

	"brank.as/petnet/profile/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestCreateBranch(t *testing.T) {
	ts := newTestStorage(t)

	ctx := context.Background()
	oid := uuid.NewString()
	_, err := ts.CreateOrgProfile(ctx, &storage.OrgProfile{
		OrgID:  oid,
		UserID: uuid.NewString(),
	})
	if err != nil {
		t.Fatal(err)
	}

	o := cmpopts.IgnoreFields(storage.Branch{}, "ID", "Created", "Updated")
	fs := []storage.Branch{
		{
			OrgID:        oid,
			OrgProfileID: uuid.NewString(),
			Title:        "title1",
			BranchAddress: storage.BranchAddress{
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
			OrgID:        oid,
			OrgProfileID: uuid.NewString(),
			Title:        "title2",
			BranchAddress: storage.BranchAddress{
				Address1:   "addr2",
				City:       "city2",
				State:      "state2",
				PostalCode: "123456",
			},
			PhoneNumber:   "123456",
			FaxNumber:     "123456",
			ContactPerson: "contact2",
		},
	}
	got, err := ts.CreateBranch(ctx, fs[0])
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(fs[0], *got, o) {
		t.Error("(-want +got): ", cmp.Diff(fs[0], *got, o))
	}
	got, err = ts.UpsertBranch(ctx, fs[1])
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(fs[1], *got, o) {
		t.Error("(-want +got): ", cmp.Diff(fs[1], *got, o))
	}
}

func TestListBranches(t *testing.T) {
	ts := newTestStorage(t)

	ctx := context.Background()
	oid := uuid.NewString()
	title := ""
	_, err := ts.CreateOrgProfile(ctx, &storage.OrgProfile{
		OrgID:  oid,
		UserID: uuid.NewString(),
	})
	if err != nil {
		t.Fatal(err)
	}

	o := cmpopts.IgnoreFields(storage.Branch{}, "ID", "Created", "Updated")
	fs := []storage.Branch{
		{
			OrgID:        oid,
			OrgProfileID: uuid.NewString(),
			Title:        "title1",
			BranchAddress: storage.BranchAddress{
				Address1:   "addr1",
				City:       "city",
				State:      "state",
				PostalCode: "12345",
			},
			PhoneNumber:   "12345",
			FaxNumber:     "12345",
			ContactPerson: "contact",
			Count:         5,
		},
		{
			OrgID:        oid,
			OrgProfileID: uuid.NewString(),
			Title:        "title2",
			BranchAddress: storage.BranchAddress{
				Address1:   "addr2",
				City:       "city2",
				State:      "state2",
				PostalCode: "123452",
			},
			PhoneNumber:   "123452",
			FaxNumber:     "123452",
			ContactPerson: "contact2",
			Count:         5,
		},
		{
			OrgID:        oid,
			OrgProfileID: uuid.NewString(),
			Title:        "title3",
			BranchAddress: storage.BranchAddress{
				Address1:   "addr13",
				City:       "city3",
				State:      "state3",
				PostalCode: "123453",
			},
			PhoneNumber:   "123453",
			FaxNumber:     "123453",
			ContactPerson: "contact3",
			Count:         5,
		},
		{
			OrgID:        oid,
			OrgProfileID: uuid.NewString(),
			Title:        "title4",
			BranchAddress: storage.BranchAddress{
				Address1:   "addr14",
				City:       "city4",
				State:      "state4",
				PostalCode: "123454",
			},
			PhoneNumber:   "123454",
			FaxNumber:     "123454",
			ContactPerson: "contact4",
			Count:         5,
		},
		{
			OrgID:        oid,
			OrgProfileID: uuid.NewString(),
			Title:        "title5",
			BranchAddress: storage.BranchAddress{
				Address1:   "addr15",
				City:       "city5",
				State:      "state5",
				PostalCode: "123455",
			},
			PhoneNumber:   "123455",
			FaxNumber:     "123455",
			ContactPerson: "contact5",
			Count:         5,
		},
	}
	for _, f := range fs {
		if _, err := ts.CreateBranch(ctx, f); err != nil {
			t.Error(err)
		}
	}

	tests := []struct {
		name string
		f    storage.LimitOffsetFilter
		want []storage.Branch
	}{
		{
			name: "No Limit",
			f:    storage.LimitOffsetFilter{},
			want: []storage.Branch{fs[0], fs[1], fs[2], fs[3], fs[4]},
		},
		{
			name: "First Two",
			f: storage.LimitOffsetFilter{
				Limit:  2,
				Offset: 0,
			},
			want: []storage.Branch{fs[0], fs[1]},
		},
		{
			name: "Next Two",
			f: storage.LimitOffsetFilter{
				Limit:  2,
				Offset: 2,
			},
			want: []storage.Branch{fs[2], fs[3]},
		},
		{
			name: "Last One",
			f: storage.LimitOffsetFilter{
				Limit:  2,
				Offset: 4,
			},
			want: []storage.Branch{fs[4]},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := ts.ListBranches(ctx, oid, test.f, title)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.want, got, o) {
				t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
			}
		})
	}
}

package userprofile

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"brank.as/petnet/profile/core/profile"
	"brank.as/petnet/profile/storage/postgres"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	ppb "brank.as/petnet/gunk/dsa/v1/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestListProfile(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	st, cleanup := postgres.NewTestStorage(os.Getenv("DATABASE_CONNECTION"), filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(cleanup)

	h := New(profile.New(st))

	oid := "20000000-0000-0000-0000-000000000000"
	want := []*ppb.Profile{
		{
			UserID: uuid.New().String(),
			OrgID:  oid,
			Email:  "email@example.com",
		},
		{
			UserID: uuid.New().String(),
			OrgID:  oid,
			Email:  "email2@example.com",
		},
	}
	for _, pf := range want {
		_, err := h.CreateUserProfile(ctx, &ppb.CreateUserProfileRequest{
			Profile: pf,
		})
		if err != nil {
			t.Fatal("h.CreateProfile: ", err)
		}
	}

	o := cmp.Options{
		cmpopts.IgnoreFields(ppb.Profile{}, "ID", "Created", "Updated", "Deleted"),
		cmpopts.IgnoreUnexported(
			ppb.Profile{}, timestamppb.Timestamp{},
		),
	}
	got, err := h.ListUserProfiles(ctx, &ppb.ListUserProfilesRequest{
		OrgID: oid,
	})
	if err != nil {
		t.Fatal("h.ListProfiles: ", err)
	}
	if !cmp.Equal(want, got.Profiles, o) {
		t.Error("(-want +got): ", cmp.Diff(want, got.Profiles), o)
	}
}

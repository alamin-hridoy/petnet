package userprofile

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"brank.as/petnet/profile/core/profile"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/profile/storage/postgres"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/types/known/timestamppb"

	ppb "brank.as/petnet/gunk/dsa/v1/user"
)

func TestGetProfile(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	st, cleanup := postgres.NewTestStorage(os.Getenv("DATABASE_CONNECTION"), filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(cleanup)
	h := New(profile.New(st))
	uid := "10000000-0000-0000-0000-000000000000"
	want := &ppb.CreateUserProfileRequest{
		Profile: &ppb.Profile{
			UserID: uid,
			OrgID:  "20000000-0000-0000-0000-000000000000",
			Email:  "email@example.com",
		},
	}
	_, err := h.CreateUserProfile(ctx, want)
	if err != nil {
		t.Fatal("h.CreateProfile: ", err)
	}
	o := []cmp.Option{
		cmpopts.IgnoreFields(ppb.Profile{}, "ID", "Created", "Updated", "Deleted"),
		cmpopts.IgnoreFields(timestamppb.Timestamp{}, "Seconds", "Nanos"),
		cmpopts.IgnoreUnexported(
			ppb.Profile{}, timestamppb.Timestamp{},
		),
	}
	got, err := h.GetUserProfile(ctx, &ppb.GetUserProfileRequest{
		UserID: uid,
	})
	if err != nil && err != storage.NotFound {
		t.Fatal("h.GetProfile: ", err)
	}
	if err != storage.NotFound {
		if !cmp.Equal(want.Profile, got.Profile, o...) {
			t.Error("(-want +got): ", cmp.Diff(want.Profile, got.Profile), o)
		}
		if got.Profile.ID == "" {
			t.Error("profile id not be created")
		}
	}
}

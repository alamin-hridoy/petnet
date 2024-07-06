package userprofile

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	ppb "brank.as/petnet/gunk/dsa/v1/user"
	"brank.as/petnet/profile/core/profile"
	"brank.as/petnet/profile/storage/postgres"

	"github.com/sirupsen/logrus/hooks/test"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func TestUpdateProfile(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	st, cleanup := postgres.NewTestStorage(os.Getenv("DATABASE_CONNECTION"), filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(cleanup)
	h := New(profile.New(st))
	ts := tspb.New(time.Unix(1515151515, 0))
	res, err := h.CreateUserProfile(ctx, &ppb.CreateUserProfileRequest{
		Profile: &ppb.Profile{
			UserID:  "10000000-0000-0000-0000-000000000000",
			OrgID:   "20000000-0000-0000-0000-000000000000",
			Email:   "email@example.com",
			Deleted: ts,
		},
	})
	if err != nil {
		t.Fatal("h.CreateProfile: ", err)
	}
	tests := []struct {
		desc            string
		profileConflict bool
		req             *ppb.UpdateUserProfileRequest
		wantErr         bool
	}{
		{
			desc: "All Success",
			req: &ppb.UpdateUserProfileRequest{
				Profile: &ppb.Profile{
					ID:      res.GetID(),
					Deleted: ts,
				},
			},
		},
		{
			desc:    "Missing ID",
			wantErr: true,
			req: &ppb.UpdateUserProfileRequest{
				Profile: &ppb.Profile{
					Deleted: ts,
				},
			},
		},
	}
	test.NewNullLogger()
	for _, test := range tests {
		tt := test
		t.Run(tt.desc, func(t *testing.T) {
			_, err = h.UpdateUserProfile(ctx, test.req)
			if err == nil && tt.wantErr {
				t.Fatal("h.UpdateProfile: want error got nil")
			}
		})
	}
}

package userprofile

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus/hooks/test"

	"brank.as/petnet/profile/core/profile"
	"brank.as/petnet/profile/storage/postgres"

	ppb "brank.as/petnet/gunk/dsa/v1/user"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateProfile(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ts := tspb.New(time.Unix(1515151515, 0))
	tests := []struct {
		desc            string
		profileConflict bool
		req             *ppb.CreateUserProfileRequest
		wantErr         bool
	}{
		{
			desc: "All Success",
			req: &ppb.CreateUserProfileRequest{
				Profile: &ppb.Profile{
					UserID:  "10000000-0000-0000-0000-000000000000",
					OrgID:   "20000000-0000-0000-0000-000000000000",
					Email:   "email@example.com",
					Deleted: ts,
				},
			},
		},
		{
			desc:    "Missing orgID",
			wantErr: true,
			req: &ppb.CreateUserProfileRequest{
				Profile: &ppb.Profile{
					UserID: "10000000-0000-0000-0000-000000000000",
					Email:  "email2@example.com",
				},
			},
		},
		{
			desc:    "Missing userID",
			wantErr: true,
			req: &ppb.CreateUserProfileRequest{
				Profile: &ppb.Profile{
					OrgID: "20000000-0000-0000-0000-000000000000",
					Email: "email3@example.com",
				},
			},
		},
		{
			desc:    "Missing email",
			wantErr: true,
			req: &ppb.CreateUserProfileRequest{
				Profile: &ppb.Profile{
					UserID: "10000000-0000-0000-0000-000000000000",
					OrgID:  "20000000-0000-0000-0000-000000000000",
				},
			},
		},
	}
	test.NewNullLogger()
	for _, test := range tests {
		tt := test
		t.Run(tt.desc, func(t *testing.T) {
			st, cleanup := postgres.NewTestStorage(os.Getenv("DATABASE_CONNECTION"), filepath.Join("..", "..", "migrations", "sql"))
			t.Cleanup(cleanup)

			h := New(profile.New(st))
			_, err := h.CreateUserProfile(ctx, tt.req)
			if err != nil && !tt.wantErr {
				t.Fatal("h.CreateProfile: ", err)
			}
			if err == nil && tt.wantErr {
				t.Fatal("h.CreateProfile: want error got nil")
			}
		})
	}
}

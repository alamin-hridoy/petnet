package userprofile

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"

	"brank.as/petnet/profile/core/profile"
	"brank.as/petnet/profile/storage/postgres"

	ppb "brank.as/petnet/gunk/dsa/v1/user"
)

func TestUpdateProfileOrgID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	st, cleanup := postgres.NewTestStorage(os.Getenv("DATABASE_CONNECTION"), filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(cleanup)

	h := New(profile.New(st))

	tests := []struct {
		desc            string
		profileConflict bool
		reqq            *ppb.UpdateUserProfileByOrgIDRequest
		wantErr         bool
	}{
		{
			desc: "All Success",
			reqq: &ppb.UpdateUserProfileByOrgIDRequest{
				OldOrgID: "20000000-0000-0000-0000-000000000000",
				NewOrgID: "28000000-0000-0000-0000-000000000000",
				UserID:   "10000000-0000-0000-0000-000000000000",
			},
		},
	}
	test.NewNullLogger()
	for _, test := range tests {
		tt := test
		t.Run(tt.desc, func(t *testing.T) {
			_, err := h.UpdateUserProfileByOrgID(ctx, test.reqq)
			if err == nil && tt.wantErr {
				t.Fatal("h.UpdateUserProfileByOrgID: want error got nil")
			}
		})
	}
}

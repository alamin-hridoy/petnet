package hydra

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUserDetails(t *testing.T) {
	tests := []struct {
		name  string
		extra map[string]string
		orgID string
		user  string
		err   error
	}{
		{
			name: "org/user",
			extra: map[string]string{
				"org_id":   "testorg",
				"username": "testuser",
				"other":    "randomstuff",
			},
			orgID: "testorg",
			user:  "testuser",
		},
		{
			name: "missingExtra",
			err:  status.Error(codes.Unauthenticated, "session details not found"),
		},
	}
	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = metautils.ExtractIncoming(ctx).ToIncoming(ctx)
			if tst.extra != nil {
				ctx = context.WithValue(ctx, &extraData{}, tst.extra)
			}
			u := &UserDetails{}
			got, err := u.Metadata(ctx)
			if !cmp.Equal(tst.err, err, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(tst.err, err, cmpopts.EquateErrors()))
			}
			if got == nil {
				if tst.err != nil {
					return
				}
				t.Fatal("unexpected nil context returned")
			}
			orgID := OrgID(got)
			if !cmp.Equal(tst.orgID, orgID) {
				t.Error(cmp.Diff(tst.orgID, orgID))
			}
			usr := Username(got)
			if !cmp.Equal(tst.user, usr) {
				t.Error(cmp.Diff(tst.user, usr))
			}
		})
	}
}

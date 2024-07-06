package keto

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDeleteRole(t *testing.T) {
	baseURL := os.Getenv("KETO_URL")
	if baseURL == "" {
		t.Skip("missing env 'KETO_URL'")
	}
	t.Parallel()
	s := New(baseURL)

	tests := []struct {
		name    string
		r       Role
		err     error
		wantErr bool
	}{
		{
			name: "GetSuccess",
			r: Role{
				Members: []string{"member-1"},
			},
		},
		{
			name: "DoesNotExist",
			r: Role{
				Members: []string{"member-2"},
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			roleID, err := s.CreateRole(ctx, test.r)
			if !cmp.Equal(test.err, err) {
				t.Error(cmp.Diff(test.err, err))
			}

			id := roleID
			if test.wantErr {
				id = "wrong"
			}

			t.Cleanup(func() { s.DeleteRole(ctx, roleID) })

			if err := s.DeleteRole(ctx, id); !cmp.Equal(test.err, err) {
				t.Error(cmp.Diff(test.err, err))
			}
		})
	}
}

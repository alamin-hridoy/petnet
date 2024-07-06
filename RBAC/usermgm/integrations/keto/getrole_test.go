package keto

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetRole(t *testing.T) {
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

			role, err := s.GetRole(ctx, id)

			if (err != nil) != test.wantErr {
				t.Fatalf("Svc.GetRole() error = %v, wantErr %v", err, test.wantErr)
			}

			test.r.ID = roleID
			if !test.wantErr && !cmp.Equal(role, test.r) {
				t.Errorf("Svc.GetRole() = %v, want %v", role, test.r)
			}
		})
	}
}

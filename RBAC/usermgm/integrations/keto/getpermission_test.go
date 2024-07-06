package keto

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"brank.as/rbac/usermgm/core"
)

func TestSvc_GetPermission(t *testing.T) {
	baseURL := os.Getenv("KETO_URL")
	if baseURL == "" {
		t.Skip("missing env 'KETO_URL'")
	}
	t.Parallel()
	s := New(baseURL)

	tests := []struct {
		name    string
		p       Permission
		wantErr bool
	}{
		{
			name: "DoesNotExist",
			p: Permission{
				Description: "missing",
				Environment: "some-id",
				Allow:       true,
				Actions:     []string{"get"},
				Resources:   []string{"resource-object"},
				Groups:      []string{"test-user-get-1"},
			},
			wantErr: true,
		},
		{
			name: "Success",
			p: Permission{
				Description: "missing",
				Environment: "some-other-id",
				Allow:       true,
				Actions:     []string{"delete"},
				Resources:   []string{"resource-object-get"},
				Groups:      []string{"test-user-get-2"},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			id, err := s.CreatePermission(ctx, test.p)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { s.DeletePermission(ctx, id) })
			del := id
			if test.wantErr {
				del = "wrong"
			}
			val, err := s.ValidateRequest(ctx, core.Validation{
				Environment: test.p.Environment,
				Action:      test.p.Actions[0],
				Resource:    test.p.Resources[0],
				ID:          test.p.Groups[0],
			})
			if err != nil {
				t.Errorf("%#v", err)
			}
			if !val {
				t.Fatal("permission not registered")
			}

			perm, err := s.GetPermission(ctx, del)
			if (err != nil) != test.wantErr {
				t.Fatalf("Svc.GetPermission() error = %v, wantErr %v", err, test.wantErr)
			}

			test.p.ID = del
			if !test.wantErr && !cmp.Equal(perm, test.p) {
				t.Errorf("Svc.GetPermission() = %v, want %v", perm, test.p)
			}
		})
	}
}

package keto

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ory/keto-client-go/client/engines"

	"brank.as/rbac/usermgm/core"
)

func TestSvc_UpdatePermission(t *testing.T) {
	baseURL := os.Getenv("KETO_URL")
	if baseURL == "" {
		t.Skip("missing env 'KETO_URL'")
	}
	t.Parallel()
	s := New(baseURL)

	t.Cleanup(func() {
		l, err := s.cl.Engines.ListOryAccessControlPolicies(
			engines.NewListOryAccessControlPoliciesParams().
				WithFlavor("exact"),
		)
		if err != nil {
			return
		}
		for _, p := range l.Payload {
			s.DeletePermission(context.Background(), p.ID)
		}
	})

	tests := []struct {
		name  string
		p     Permission
		newIn bool
		err   error
	}{
		{
			name: "DoesNotUpdate",
			p: Permission{
				Description: "missing",
				Environment: "some-id",
				Allow:       true,
				Actions:     []string{"get"},
				Resources:   []string{"resource-object"},
				Groups:      []string{"test-user-update-1"},
			},
		},
		{
			name: "Success",
			p: Permission{
				Description: "missing",
				Environment: "some-other-id",
				Allow:       true,
				Actions:     []string{"delete"},
				Resources:   []string{"resource-object-update"},
				Groups:      []string{"test-user-update-2"},
			},
			newIn: true,
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

			groups := test.p.Groups
			if test.newIn {
				groups = []string{test.p.Groups[0] + "updated"}
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

			test.p.ID = id
			test.p.Groups = groups
			if err := s.UpdatePermission(ctx, test.p); !cmp.Equal(test.err, err) {
				t.Error(cmp.Diff(test.err, err))
			}
		})
	}
}

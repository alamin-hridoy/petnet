package keto

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestCreateRole(t *testing.T) {
	baseURL := os.Getenv("KETO_URL")
	if baseURL == "" {
		t.Skip("missing env 'KETO_URL'")
	}
	t.Parallel()
	s := New(baseURL)

	tests := []struct {
		name string
		r    Role
		err  error
	}{
		{
			name: "CreateSuccess",
			r: Role{
				Members: []string{"test-group"},
			},
		},
		{
			name: "EmptyPermission",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			got, err := s.CreateRole(ctx, test.r)
			if !cmp.Equal(test.err, err) {
				t.Error(cmp.Diff(test.err, err))
			}

			t.Cleanup(func() { s.DeletePermission(ctx, got) })
			if _, err := uuid.Parse(got); err != nil && test.err == nil {
				t.Fatalf("got empty ID")
			}
			// TODO(Chad): Get role and validate
		})
	}
}

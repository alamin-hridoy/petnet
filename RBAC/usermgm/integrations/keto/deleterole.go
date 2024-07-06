package keto

import (
	"context"
	"fmt"

	"github.com/ory/keto-client-go/client/engines"
)

func (s *Svc) DeleteRole(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}
	pm := engines.NewDeleteOryAccessControlPolicyRoleParams().
		WithFlavor("exact").
		WithID(id).
		WithContext(ctx)

	if _, err := s.cl.Engines.DeleteOryAccessControlPolicyRole(pm); err != nil {
		return err
	}

	return nil
}

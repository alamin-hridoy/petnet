package keto

import (
	"context"
	"errors"
	"fmt"

	"github.com/ory/keto-client-go/client/engines"
)

func (s *Svc) DeletePermission(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}

	pm := engines.NewDeleteOryAccessControlPolicyParams().
		WithFlavor("exact").WithID(id).WithContext(ctx)
	if _, err := s.cl.Engines.DeleteOryAccessControlPolicy(pm); err != nil {
		e := &engines.GetOryAccessControlPolicyNotFound{}
		if errors.As(err, &e) {
			return nil
		}
		return err
	}
	return nil
}

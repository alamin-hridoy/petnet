package keto

import (
	"context"
	"fmt"

	"github.com/ory/keto-client-go/client/engines"
)

func (s *Svc) GetRolePermissions(ctx context.Context, role string) ([]string, error) {
	if role == "" {
		return nil, fmt.Errorf("id is required")
	}

	pm := engines.NewListOryAccessControlPoliciesParamsWithContext(ctx).
		WithFlavor("exact").WithSubject(&role)
	r, err := s.cl.Engines.ListOryAccessControlPolicies(pm)
	if err != nil {
		return nil, err
	}

	payload := r.GetPayload()
	if payload == nil {
		return nil, fmt.Errorf("received nil payload")
	}

	p := make([]string, len(payload))
	for i, v := range payload {
		p[i] = v.ID
	}

	return p, nil
}

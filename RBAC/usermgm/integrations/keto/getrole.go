package keto

import (
	"context"
	"fmt"

	"github.com/ory/keto-client-go/client/engines"
)

func (s *Svc) GetRole(ctx context.Context, id string) (Role, error) {
	if id == "" {
		return Role{}, fmt.Errorf("id is required")
	}
	pm := engines.NewGetOryAccessControlPolicyRoleParams().
		WithFlavor("exact").
		WithID(id).
		WithContext(ctx)
	r, err := s.cl.Engines.GetOryAccessControlPolicyRole(pm)
	if err != nil {
		return Role{}, err
	}

	payload := r.GetPayload()
	if payload == nil {
		return Role{}, fmt.Errorf("received nil payload")
	}

	return Role{
		ID:      payload.ID,
		Members: payload.Members,
	}, nil
}

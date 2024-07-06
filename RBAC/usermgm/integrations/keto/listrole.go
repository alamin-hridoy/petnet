package keto

import (
	"context"
	"fmt"

	"github.com/ory/keto-client-go/client/engines"
)

func (s *Svc) ListRoles(ctx context.Context, uid string) ([]string, error) {
	if uid == "" {
		return nil, fmt.Errorf("user id required")
	}
	pm := engines.NewListOryAccessControlPolicyRolesParamsWithContext(ctx).
		WithFlavor("exact").
		WithMember(&uid)
	r, err := s.cl.Engines.ListOryAccessControlPolicyRoles(pm)
	if err != nil {
		return nil, err
	}

	payload := r.GetPayload()
	if payload == nil {
		return nil, fmt.Errorf("received nil payload")
	}
	lst := make([]string, 0, len(payload))
	for _, p := range payload {
		lst = append(lst, p.ID)
	}

	return lst, nil
}

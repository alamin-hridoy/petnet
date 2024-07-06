package keto

import (
	"context"
	"fmt"

	"github.com/ory/keto-client-go/client/engines"
	"github.com/ory/keto-client-go/models"
)

func (s *Svc) UpdateRole(ctx context.Context, ro Role) (string, error) {
	if ro.ID == "" {
		return "", fmt.Errorf("role id is required")
	}
	m := make(map[string]struct{}, len(ro.Members))
	for _, g := range ro.Members {
		m[g] = struct{}{}
	}
	// Deduplicate
	ro.Members = ro.Members[:0]
	for mem := range m {
		ro.Members = append(ro.Members, mem)
	}
	pm := engines.NewUpsertOryAccessControlPolicyRoleParams().
		WithFlavor("exact").
		WithBody(&models.OryAccessControlPolicyRole{
			ID:      ro.ID,
			Members: ro.Members,
		}).WithContext(ctx)
	r, err := s.cl.Engines.UpsertOryAccessControlPolicyRole(pm)
	if err != nil {
		return "", err
	}

	payload := r.GetPayload()
	if payload == nil {
		return "", fmt.Errorf("received nil payload")
	}

	return payload.ID, nil
}

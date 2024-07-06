package keto

import (
	"context"
	"fmt"

	"github.com/ory/keto-client-go/client/engines"
	"github.com/ory/keto-client-go/models"
)

type Role struct {
	// ID is the permission identifier.
	ID string
	// Members is the list of accounts assigned the role
	Members []string
}

func (s *Svc) CreateRole(ctx context.Context, ro Role) (string, error) {
	if ro.ID != "" {
		return "", fmt.Errorf("role with id should be updated")
	}
	pm := engines.NewUpsertOryAccessControlPolicyRoleParams().
		WithFlavor("exact").
		WithBody(&models.OryAccessControlPolicyRole{
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

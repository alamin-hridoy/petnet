package keto

import (
	"context"
	"fmt"

	"github.com/ory/keto-client-go/client/engines"
)

func (s *Svc) GetPermission(ctx context.Context, id string) (Permission, error) {
	if id == "" {
		return Permission{}, fmt.Errorf("id is required")
	}

	pm := engines.NewGetOryAccessControlPolicyParams().
		WithFlavor("exact").WithID(id).WithContext(ctx)
	r, err := s.cl.Engines.GetOryAccessControlPolicy(pm)
	if err != nil {
		return Permission{}, err
	}

	payload := r.GetPayload()
	if payload == nil {
		return Permission{}, fmt.Errorf("received nil payload")
	}

	e := make(map[string]string, len(payload.Conditions))
	for k, v := range payload.Conditions {
		if v2, ok := v.(map[string]interface{}); ok {
			for _, v3 := range v2 {
				if v4, ok := v3.(map[string]interface{}); ok {
					for _, value := range v4 {
						e[k] = value.(string)
					}
				}
			}
		}
	}

	return Permission{
		ID:          payload.ID,
		Description: payload.Description,
		Environment: e["env"],
		Actions:     payload.Actions,
		Allow:       (payload.Effect == string(allow)),
		Resources:   payload.Resources,
		Groups:      payload.Subjects,
	}, nil
}

package keto

import (
	"context"
	"fmt"

	"github.com/ory/keto-client-go/client/engines"
	"github.com/ory/keto-client-go/models"
)

type oryEffect string

const (
	allow oryEffect = "allow"
	deny  oryEffect = "deny"
)

type Permission struct {
	// ID is the permission identifier.
	ID string
	// Description gives a user understandable summary of the permission.
	Description string
	// Environment is the target OrgID, prod/sandbox runtime, etc.
	// Which define the operating scope of the permission within the entire data/operation space.
	Environment string
	// Permission should allow (grant) or deny (restrict) the action.
	Allow bool
	// Actions to be taken on the resource(s).
	// Can be data operations like read/write/delete or
	// more basic endpoint action like call/use.
	Actions []string
	// Resources should uniquely identify the data or endpoint that is being acted upon.
	// Can be multiple data objects, especially if an endpoint operates on them simultaneously.
	Resources []string
	// Groups or Users assigned the permission.
	// In most cases, these should correspond to permission Groups.
	// Assignment to user(s) should be limited to Admin or other highly-restricted permissions.
	Groups []string
}

type Condition struct {
	Type    string `json:"type"`
	Options `json:"options"`
}

type Options struct {
	Equals string `json:"equals"`
}

func (s *Svc) CreatePermission(ctx context.Context, p Permission) (string, error) {
	if p.ID != "" {
		return "", fmt.Errorf("permissions with id should be updated")
	}
	c := map[string]interface{}{}
	if p.Environment != "" {
		c["env"] = Condition{
			Type: "StringEqualCondition",
			Options: Options{
				Equals: p.Environment,
			},
		}
	}

	e := allow
	if !p.Allow {
		e = deny
	}
	pm := engines.NewUpsertOryAccessControlPolicyParams().
		WithFlavor("exact").
		WithBody(&models.OryAccessControlPolicy{
			Actions:     p.Actions,
			Conditions:  c,
			Description: p.Description,
			Effect:      string(e),
			Resources:   p.Resources,
			Subjects:    p.Groups,
		}).WithContext(ctx)
	r, err := s.cl.Engines.UpsertOryAccessControlPolicy(pm)
	if err != nil {
		return "", err
	}

	payload := r.GetPayload()
	if payload == nil {
		return "", fmt.Errorf("received nil payload")
	}

	return payload.ID, nil
}

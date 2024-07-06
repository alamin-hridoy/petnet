package keto

import (
	"context"
	"fmt"

	"github.com/ory/keto-client-go/client/engines"
	"github.com/ory/keto-client-go/models"
)

func (s *Svc) UpdatePermission(ctx context.Context, p Permission) error {
	if p.ID == "" {
		return fmt.Errorf("updating permissions require id")
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
	m := make(map[string]struct{}, len(p.Groups))
	for _, g := range p.Groups {
		m[g] = struct{}{}
	}
	// Deduplicate
	p.Groups = p.Groups[:0]
	for pr := range m {
		p.Groups = append(p.Groups, pr)
	}
	pm := engines.NewUpsertOryAccessControlPolicyParams().
		WithFlavor("exact").
		WithBody(&models.OryAccessControlPolicy{
			ID:          p.ID,
			Actions:     p.Actions,
			Conditions:  c,
			Description: p.Description,
			Effect:      string(e),
			Resources:   p.Resources,
			Subjects:    p.Groups,
		}).WithContext(ctx)
	r, err := s.cl.Engines.UpsertOryAccessControlPolicy(pm)
	if err != nil {
		return err
	}

	payload := r.GetPayload()
	if payload == nil {
		return fmt.Errorf("received nil payload")
	}

	return nil
}

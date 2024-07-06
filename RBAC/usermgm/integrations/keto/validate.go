package keto

import (
	"context"
	"errors"

	"github.com/ory/keto-client-go/client/engines"
	"github.com/ory/keto-client-go/models"

	"brank.as/rbac/usermgm/core"
)

func (s *Svc) ValidateRequest(ctx context.Context, v core.Validation) (bool, error) {
	b := &models.OryAccessControlPolicyAllowedInput{
		Action:   v.Action,
		Resource: v.Resource,
		Subject:  v.ID,
	}
	if v.Environment != "" {
		b.Context = map[string]interface{}{
			"env": v.Environment,
		}
	}
	p := engines.NewDoOryAccessControlPoliciesAllowParams().
		WithFlavor("exact").
		WithBody(b).WithContext(ctx)
	r, err := s.cl.Engines.DoOryAccessControlPoliciesAllow(p)
	if err != nil {
		e := &engines.DoOryAccessControlPoliciesAllowForbidden{}
		if errors.As(err, &e) {
			return false, nil
		}
		return false, err
	}
	if r.Error() != "" && (r.GetPayload() == nil || r.GetPayload().Allowed == nil) {
		return false, r
	}
	return *r.GetPayload().Allowed, nil
}

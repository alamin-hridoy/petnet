package challenge

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/usermgm/core"

	"brank.as/rbac/serviceutil/logging"
)

func (s *Svc) Validate(ctx context.Context, v core.Validation) (*core.Identity, error) {
	log := logging.FromContext(ctx).WithField("method", "core.permissions.validate")
	idn, err := s.fetchAccount(ctx, v.ID)
	if err != nil {
		logging.WithError(err, log).Error("account not found")
		return nil, status.Error(codes.PermissionDenied, "not authorized")
	}

	org, err := s.id.GetOrgByID(ctx, idn.OrgID)
	if err != nil {
		logging.WithError(err, log).Error("org not found")
		return nil, status.Error(codes.PermissionDenied, "not authorized")
	}
	if !org.Active {
		log.WithField("org", org.ID).Error("org inactive")
		return nil, status.Error(codes.PermissionDenied, "organization inactive")
	}
	if v.OrgID == "" {
		v.OrgID = org.ID
	}

	v.Resource = fmt.Sprintf("org:%s:%s", v.OrgID, v.Resource)
	ok, err := s.val.ValidateRequest(ctx, v)
	switch {
	case err != nil:
		logging.WithError(err, log).Error("keto validate")
		fallthrough
	case !ok:
		log.WithField("action", v).Error("permission denied")
		return nil, status.Error(codes.PermissionDenied, "not authorized")
	}
	return idn, nil
}

func (s *Svc) fetchAccount(ctx context.Context, id string) (*core.Identity, error) {
	if u, err := s.id.GetUserByID(ctx, id); err == nil {
		if u.Deleted.Valid {
			return nil, status.Error(codes.FailedPrecondition, "account disabled")
		}
		return &core.Identity{
			ID:    id,
			OrgID: u.OrgID,
			Name:  u.Username,
		}, nil
	}
	if s, err := s.id.GetSvcAccountByID(ctx, id); err == nil {
		return &core.Identity{
			ID:    id,
			OrgID: s.OrgID,
			Name:  s.ClientName,
		}, nil
	}
	return nil, status.Error(codes.NotFound, "account invalid")
}

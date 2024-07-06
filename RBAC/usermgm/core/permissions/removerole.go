package permissions

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"
)

func (s *Svc) RemoveRole(ctx context.Context, g core.Grant) (*core.Role, error) {
	log := logging.FromContext(ctx).WithField("method", "core.permissions.removerole")

	meta, err := s.store.GetRole(ctx, g.RoleID)
	if err != nil {
		logging.WithError(err, log).Error("metadata fetch")
		return nil, status.Error(codes.Internal, "role not found")
	}
	if meta.OrgID != hydra.OrgID(ctx) {
		log.WithField("request_org", hydra.OrgID(ctx)).
			WithField("role_org", meta.OrgID).Error("incorrect org")
		return nil, status.Error(codes.NotFound, "role not found")
	}

	r, err := s.keto.GetRole(ctx, g.RoleID)
	if err != nil {
		logging.WithError(err, log).Error("role fetch")
		return nil, status.Error(codes.NotFound, "role not found")
	}

	mm := r.Members
	r.Members = r.Members[:0]
	for _, id := range mm {
		if id != g.GrantID {
			r.Members = append(r.Members, id)
		}
	}

	if _, err := s.keto.UpdateRole(ctx, r); err != nil {
		logging.WithError(err, log).Error("permission update")
		return nil, status.Error(codes.Internal, "processing failed")
	}

	r, err = s.keto.GetRole(ctx, g.RoleID)
	if err != nil {
		logging.WithError(err, log).Error("role refresh")
		return nil, status.Error(codes.Internal, "processing failed")
	}

	p, err := s.keto.GetRolePermissions(ctx, g.RoleID)
	if err != nil {
		logging.WithError(err, log).Error("permission fetch")
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return &core.Role{
		ID:          r.ID,
		OrgID:       meta.OrgID,
		Name:        meta.Name,
		Desc:        meta.Description,
		Permissions: p,
		Members:     r.Members,
	}, nil
}
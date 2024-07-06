package permissions

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"
)

func (s *Svc) RoleRevoke(ctx context.Context, g core.Grant) (*core.Role, error) {
	log := logging.FromContext(ctx).WithField("method", "core.permissions.rolegrant")

	meta, err := s.store.GetRole(ctx, g.RoleID)
	if err != nil {
		logging.WithError(err, log).Error("metadata fetch")
		return nil, status.Error(codes.Internal, "processing failed")
	}
	if meta.OrgID != hydra.OrgID(ctx) {
		log.WithField("request_org", hydra.OrgID(ctx)).
			WithField("role_org", meta.OrgID).Error("incorrect org")
		return nil, status.Error(codes.NotFound, "role not found")
	}

	p, err := s.keto.GetPermission(ctx, g.GrantID)
	if err != nil {
		logging.WithError(err, log).Error("permission fetch")
		return nil, status.Error(codes.NotFound, "role not found")
	}

	gr := p.Groups
	p.Groups = p.Groups[:0]
	for _, id := range gr {
		if id != g.RoleID {
			p.Groups = append(p.Groups, id)
		}
	}

	if err := s.keto.UpdatePermission(ctx, p); err != nil {
		logging.WithError(err, log).Error("permission update")
		return nil, status.Error(codes.Internal, "processing failed")
	}

	r, err := s.keto.GetRole(ctx, g.RoleID)
	if err != nil {
		logging.WithError(err, log).Error("role refresh")
		return nil, status.Error(codes.Internal, "processing failed")
	}

	pr, err := s.keto.GetRolePermissions(ctx, g.RoleID)
	if err != nil {
		logging.WithError(err, log).Error("permission fetch")
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return &core.Role{
		ID:          r.ID,
		OrgID:       meta.OrgID,
		Name:        meta.Name,
		Desc:        meta.Description,
		Permissions: pr,
		Members:     r.Members,
	}, nil
}

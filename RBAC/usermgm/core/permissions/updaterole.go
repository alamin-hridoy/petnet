package permissions

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

func (s *Svc) UpdateRole(ctx context.Context, g core.Role) (core.Role, error) {
	log := logging.FromContext(ctx).WithField("method", "core.permissions.updaterole")

	r, err := s.store.UpdateRole(ctx, storage.Role{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Desc,
		UpdatedUID:  g.UpdatedUID,
	})
	if err != nil {
		logging.WithError(err, log).Error("db update")
		return core.Role{}, status.Error(codes.Internal, "processing failed")
	}

	ketoRole, err := s.keto.GetRole(ctx, r.ID)
	if err != nil {
		logging.WithError(err, log).Error("read keto role")
		return core.Role{}, status.Error(codes.Internal, "processing failed")
	}

	ketoRolePerms, err := s.keto.GetRolePermissions(ctx, r.ID)
	if err != nil {
		logging.WithError(err, log).Error("read keto role permission")
		return core.Role{}, status.Error(codes.Internal, "processing failed")
	}

	return core.Role{
		ID:          r.ID,
		OrgID:       r.OrgID,
		Name:        r.Name,
		Desc:        r.Description,
		CreateUID:   r.CreateUID,
		DeleteUID:   r.DeleteUID.String,
		Permissions: ketoRolePerms,
		Members:     ketoRole.Members,
	}, nil
}

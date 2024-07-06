package permissions

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/integrations/keto"
	"brank.as/rbac/usermgm/storage"
)

func (s *Svc) CreateRole(ctx context.Context, g core.Role) (string, error) {
	log := logging.FromContext(ctx).WithField("method", "core.permissions.createrole")

	id, err := s.keto.CreateRole(ctx, keto.Role{
		Members: g.Members,
	})
	if err != nil {
		logging.WithError(err, log).Error("keto create")
		return "", status.Error(codes.Internal, "processing failed")
	}
	g.ID = id

	if _, err := s.store.CreateRole(ctx, storage.Role{
		ID:          id,
		OrgID:       g.OrgID,
		Name:        g.Name,
		Description: g.Desc,
		CreateUID:   g.CreateUID,
		UpdatedUID:  g.UpdatedUID,
	}); err != nil {
		logging.WithError(err, log).Error("db create")
		return "", status.Error(codes.Internal, "processing failed")
	}

	return id, nil
}

package permissions

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/svcutil/mw/hydra"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

func (s *Svc) RevokeService(ctx context.Context, g core.Grant) error {
	log := logging.FromContext(ctx).WithField("method", "core.permissions.createrole")

	org, err := s.store.GetOrgByID(ctx, g.RoleID)
	if err != nil {
		logging.WithError(err, log).Error("getting org")
		return status.Error(codes.NotFound, "org not found")
	}
	if !org.Active {
		logging.WithError(err, log).Error("org inactive")
		return status.Error(codes.NotFound, "org inactive")
	}
	gr, err := s.store.GetAssignedService(ctx, storage.ServiceAssignment{
		OrgID:       g.RoleID,
		ServiceID:   g.GrantID,
		Environment: g.Environment,
	})
	if err != nil {
		logging.WithError(err, log).Error("getting grant")
		return err // status.Error(codes.NotFound, "service not found")
	}

	sp, err := s.store.ListOrgPermissionGrant(ctx, gr.GrantID)
	if err != nil {
		logging.WithError(err, log).Error("getting service permissions")
		return err // status.Error(codes.NotFound, "service not found")
	}

	uid := hydra.ClientID(ctx)
	if err := s.store.RevokeService(ctx, storage.ServiceAssignment{
		GrantID: g.GrantID,
		RevokeUserID: sql.NullString{
			Valid:  true,
			String: uid,
		},
	}); err != nil {
		logging.WithError(err, log).Error("store service revoke")
		return status.Error(codes.Internal, "processing failed")
	}

	for _, p := range sp {
		if err := s.DeleteOrgPermission(ctx, core.OrgPermission{
			ID:        p.ID,
			DeleteUID: uid,
		}); err != nil {
			logging.WithError(err, log).Error("delete org permission")
		}
	}

	return nil
}

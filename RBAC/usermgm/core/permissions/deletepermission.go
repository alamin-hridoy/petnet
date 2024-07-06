package permissions

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

func (s *Svc) DeleteOrgPermission(ctx context.Context, p core.OrgPermission) error {
	log := logging.FromContext(ctx).WithField("method", "core.permissions.deletepermission")

	if _, err := s.store.GetOrgPermission(ctx, p.ID); err != nil {
		if err == storage.NotFound {
			log.Info("permission deleted")
			return nil
		}
		logging.WithError(err, log).Error("read storage")
		return err // status.Error(codes.Internal, "processing failed")
	}

	if err := s.keto.DeletePermission(ctx, p.ID); err != nil {
		logging.WithError(err, log).Error("delete keto")
		return status.Error(codes.Internal, "processing failed")
	}

	if err := s.store.DeleteOrgPermission(ctx, p.ID); err != nil {
		return err
	}
	return nil
}

func (s *Svc) DeletePermission(ctx context.Context, svcPermID string) error {
	log := logging.FromContext(ctx).WithField("method", "core.permissions.deletepermission")

	if _, err := s.store.GetServicePermission(ctx, svcPermID); err != nil {
		logging.WithError(err, log).Error("read storage")
		return status.Error(codes.Internal, "processing failed")
	}

	if err := s.store.DeleteServicePermission(ctx, svcPermID); err != nil {
		logging.WithError(err, log).Error("delete storage")
		return status.Error(codes.Internal, "processing failed")
	}

	perms, err := s.store.ListPermBySvcPermID(ctx, svcPermID)
	if err != nil {
		logging.WithError(err, log).Error("read storage")
		return status.Error(codes.Internal, "processing failed")
	}

	for _, p := range perms {
		if err := s.keto.DeletePermission(ctx, p.ID); err != nil {
			logging.WithError(err, log).Error("delete keto")
			return status.Error(codes.Internal, "processing failed")
		}

		storePermission := storage.Permission{
			ID: p.ID,
			DeleteUID: sql.NullString{
				Valid:  true,
				String: p.DeleteUID.String,
			},
		}
		if _, err := s.store.DeletePermission(ctx, storePermission); err != nil {
			logging.WithError(err, log).Error("delete storage")
			return status.Error(codes.Internal, "processing failed")
		}
	}

	return nil
}

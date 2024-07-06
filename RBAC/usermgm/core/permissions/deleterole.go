package permissions

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
)

func (s *Svc) DeleteRole(ctx context.Context, id string) error {
	log := logging.FromContext(ctx).WithField("method", "core.permissions.deleterole")

	r, err := s.store.GetRole(ctx, id)
	if err != nil {
		logging.WithError(err, log).Error("read storage")
		return status.Error(codes.Internal, "processing failed")
	} else if r.Delete.Valid {
		log.Info("role deleted")
		return nil
	}

	if err := s.keto.DeleteRole(ctx, id); err != nil {
		logging.WithError(err, log).Error("keto delete")
		return status.Error(codes.Internal, "processing failed")
	}

	if _, err = s.store.DeleteRole(ctx, *r); err != nil {
		logging.WithError(err, log).Error("delete storage")
		return status.Error(codes.Internal, "processing failed")
	}

	return nil
}

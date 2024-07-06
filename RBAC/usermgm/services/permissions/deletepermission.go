package permissions

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"brank.as/rbac/serviceutil/logging"

	tspb "google.golang.org/protobuf/types/known/timestamppb"

	ppb "brank.as/rbac/gunk/v1/permissions"
)

func (s *Svc) DeletePermission(ctx context.Context, req *ppb.DeletePermissionRequest) (*ppb.DeletePermissionResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.permissions.delete")
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ID, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.perm.DeletePermission(ctx, req.ID); err != nil {
		logging.WithError(err, log).Error("delete failed")
		return nil, err

	}

	return &ppb.DeletePermissionResponse{
		Deleted: tspb.Now(),
	}, nil
}

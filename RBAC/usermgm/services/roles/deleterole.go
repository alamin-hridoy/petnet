package role

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/rbac/serviceutil/logging"

	ppb "brank.as/rbac/gunk/v1/permissions"
)

func (s *Svc) DeleteRole(ctx context.Context, req *ppb.DeleteRoleRequest) (*ppb.DeleteRoleResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.role.deleterole")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.ID, validation.Required, is.UUIDv4),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.perm.DeleteRole(ctx, req.GetID()); err != nil {
		logging.WithError(err, log).Error("processing")
		return nil, err
	}

	return &ppb.DeleteRoleResponse{
		Deleted: timestamppb.Now(),
	}, nil
}

package user

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"

	upb "brank.as/rbac/gunk/v1/user"

	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Handler) DisableUser(ctx context.Context, req *upb.DisableUserRequest) (*upb.DisableUserResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.user.DisableUser")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.UserID, validation.Required, is.UUIDv4),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	dis := core.UserActivation{
		ID:              req.GetUserID(),
		CustomEmailData: req.GetCustomEmailData(),
	}
	if err := s.usr.DisableUser(ctx, dis); err != nil {
		logging.WithError(err, log).Error("disable user")
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to disable user")
	}
	return &upb.DisableUserResponse{Updated: tspb.Now()}, nil
}

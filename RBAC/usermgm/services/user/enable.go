package user

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"

	tspb "google.golang.org/protobuf/types/known/timestamppb"

	upb "brank.as/rbac/gunk/v1/user"
)

func (s *Handler) EnableUser(ctx context.Context, req *upb.EnableUserRequest) (*upb.EnableUserResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.user.EnableUser")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.UserID, validation.Required, is.UUIDv4),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	en := core.UserActivation{
		ID:              req.GetUserID(),
		CustomEmailData: req.GetCustomEmailData(),
	}
	if err := s.usr.EnableUser(ctx, en); err != nil {
		logging.WithError(err, log).Error("enable user")
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to enable user")
	}
	return &upb.EnableUserResponse{Updated: tspb.Now()}, nil
}

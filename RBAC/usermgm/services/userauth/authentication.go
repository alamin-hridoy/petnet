package userauth

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	upb "brank.as/rbac/gunk/v1/user"
	"brank.as/rbac/serviceutil/logging"
)

func (s *Svc) AuthenticateUser(ctx context.Context, req *upb.AuthenticateUserRequest) (*upb.AuthenticateUserResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.auth.Authenticate")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.Username, validation.Required),
		validation.Field(&req.Password, validation.Required),
	); err != nil {
		logging.WithError(err, log).Info("invalid request")
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	u, err := s.perm.GetUser(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		logging.WithError(err, log).Error("failed to authenticate user")
		if status.Code(err) == codes.NotFound {
			return nil, status.Error(codes.InvalidArgument, "user not found")
		}
		return nil, status.Error(codes.InvalidArgument, "invalid argument given")
	}

	return &upb.AuthenticateUserResponse{
		UserID: u.ID,
		OrgID:  u.OrgID,
	}, nil
}

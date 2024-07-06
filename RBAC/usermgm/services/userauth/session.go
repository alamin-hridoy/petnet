package userauth

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	apb "brank.as/rbac/gunk/v1/authenticate"
)

func (s *Svc) GetSession(ctx context.Context, req *apb.GetSessionRequest) (*apb.Session, error) {
	log := logging.FromContext(ctx).WithField("method", "service.auth.getsession")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.UserID, validation.Required, is.UUIDv4),
	); err != nil {
		logging.WithError(err, log).Info("invalid request")
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	u, err := s.auth.UserSession(ctx, req.GetUserID())
	if err != nil {
		logging.WithError(err, log).Error("failed to authenticate user")
		if status.Code(err) == codes.NotFound {
			return nil, status.Error(codes.InvalidArgument, "user not found")
		}
		return nil, status.Error(codes.InvalidArgument, "invalid argument given")
	}

	return &apb.Session{
		UserID: u.ID,
		OrgID:  u.OrgID,
	}, nil
}

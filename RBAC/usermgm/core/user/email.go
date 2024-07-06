package user

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/storage"
)

// GetConfirmationCode
func (s *Svc) GetConfirmationCode(ctx context.Context, uid string) (*storage.User, string, error) {
	log := logging.FromContext(ctx).WithField("method", "user.GetConfirmationCode")

	u, err := s.usr.GetUserByID(ctx, uid)
	if err != nil {
		logging.WithError(err, log).Error("user not found")
		return nil, "", err
	}

	code, err := s.usr.GetConfirmationCode(ctx, uid)
	if err != nil {
		logging.WithError(err, log).Error("code not found")
		return nil, "", err
	}
	return u, code, nil
}

// ConfirmEmail
func (s *Svc) ConfirmEmail(ctx context.Context, code string) (*storage.User, error) {
	log := logging.FromContext(ctx).WithField("method", "user.confirmemail")

	u, err := s.usr.VerifyConfirmationCode(ctx, code)
	if err != nil {
		if err == storage.NotFound {
			logging.WithError(err, log).Error("invalid code")
			return nil, status.Error(codes.InvalidArgument, "invalid code")
		}
		logging.WithError(err, log).Error("verification")
		return nil, status.Error(codes.Internal, "failed to verify code")
	}
	return u, nil
}

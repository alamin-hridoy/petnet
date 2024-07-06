package user

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/storage"
)

// ResetPasswordInit
func (s *Svc) ResetPasswordInit(ctx context.Context, email string) error {
	log := logging.FromContext(ctx).WithField("method", "user.resetpasswordinit")

	user, err := s.usr.GetUserByEmail(ctx, email)
	if err != nil && err != storage.NotFound {
		logging.WithError(err, log).Error("user record")
		return status.Error(codes.Internal, "processing failed")
	}
	code, err := s.usr.CreatePasswordReset(ctx, user.ID, s.reset)
	if err != nil {
		log.WithError(err).Error("failed to create reset code")
		return err
	}

	log.Debug("sending password reset email")
	err = s.mail.ForgotPassword(email, user.FirstName, code)
	if err != nil {
		const errMsg = "failed to send confirmation email"
		logging.WithError(err, log).Error(errMsg)
		return status.Error(codes.Internal, errMsg)
	}
	return nil
}

// ResetPassword after confirmation.
func (s *Svc) ResetPassword(ctx context.Context, code, pass string) error {
	log := logging.FromContext(ctx).WithField("method", "user.resetpassword")

	if err := s.usr.PasswordReset(ctx, code, pass); err != nil {
		logging.WithError(err, log).Error("change password with code")
		if err == storage.NotFound {
			return status.Error(codes.InvalidArgument, "invalid code")
		}
		return err
	}
	return nil
}

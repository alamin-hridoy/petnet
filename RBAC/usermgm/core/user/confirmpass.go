package user

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"
)

// ChangePass Confirm MFA and change the password
func (s *Svc) ConfirmPass(ctx context.Context, m core.MFAChallenge) error {
	log := logging.FromContext(ctx).WithField("method", "user.confirmpass")

	if _, err := s.mfa.MFAuth(ctx, m); err != nil {
		logging.WithError(err, log).Error("validate MFA failed")
		return status.Error(codes.Internal, "failed to validate mfa")
	}

	if err := s.usr.ChangePassword(ctx, m.UserID, m.EventID); err != nil {
		logging.WithError(err, log).Error("change password failed")
		return status.Error(codes.Internal, "failed to change user password")
	}
	return nil
}

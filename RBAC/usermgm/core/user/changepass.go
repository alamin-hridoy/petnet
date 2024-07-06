package user

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

// ChangePass Validate the old password and confirm to change the password given user id
func (s *Svc) ChangePass(ctx context.Context, userID, oldPass, newPass string) (*core.MFAChallenge, error) {
	log := logging.FromContext(ctx).WithField("method", "user.changepass")

	if _, _, err := s.usr.ValidateUserPass(ctx, userID, oldPass); err == storage.NotFound {
		return nil, status.Error(codes.PermissionDenied, "old password is invalid")
	}

	ev, err := s.mfa.InitiateMFA(ctx, core.MFAChallenge{
		UserID:    userID,
		EventDesc: "Change User Password",
	})
	if err != nil {
		logging.WithError(err, log).Error("storage create change password")
		return nil, status.Error(codes.Internal, "failed to create change password")
	}

	if _, err := s.usr.CreateChangePassword(ctx, userID, ev.EventID, newPass); err != nil {
		logging.WithError(err, log).Error("storage create change password")
		return nil, status.Error(codes.Internal, "failed to create change password")
	}

	return &core.MFAChallenge{
		EventID:  ev.EventID,
		UserID:   userID,
		SourceID: ev.SourceID,
		Type:     ev.Type,
		Token:    ev.Token,
		Sources:  ev.Sources,
	}, nil
}

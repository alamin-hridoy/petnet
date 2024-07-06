package user

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

func (s *Svc) UpdateUser(ctx context.Context, u core.User) (*core.MFAChallenge, error) {
	log := logging.FromContext(ctx).WithField("method", "core.user.updateuser")
	usr, err := s.usr.GetUserByID(ctx, u.ID)
	if err != nil {
		logging.WithError(err, log).Error("fetch user")
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to read user record")
	}
	if u.Email != "" && u.Email != usr.Email {
		return nil, status.Error(codes.Unimplemented, "update email not implemented")
	}

	if u.PreferredMFA != "" && u.PreferredMFA != usr.PreferredMFA {
		m, err := s.usr.GetActiveMFAByType(ctx, u.ID, u.PreferredMFA)
		if err != nil {
			logging.WithError(err, log).Error("list mfa type from storage")
			if err == storage.NotFound {
				return nil, status.Error(codes.NotFound, "mfa source not found")
			}
			return nil, err
		}
		if len(m) == 0 {
			return nil, status.Error(codes.NotFound, "mfa source not found")
		}
		usr.PreferredMFA = u.PreferredMFA
	}

	switch {
	case u.EnableMFA:
		if usr.PreferredMFA == "" {
			return nil, status.Error(codes.FailedPrecondition, "no mfa sources available")
		}
		usr.MFALogin = true
	case u.DisableMFA:
		usr.MFALogin = false
	}
	if u.FName != "" {
		usr.FirstName = u.FName
	}
	if u.LName != "" {
		usr.LastName = u.LName
	}
	usr, err = s.usr.UpdateUserByID(ctx, *usr)
	if err != nil {
		logging.WithError(err, log).Error("list mfa type from storage")
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "mfa source not found")
		}
		return nil, err
	}
	return nil, nil
}

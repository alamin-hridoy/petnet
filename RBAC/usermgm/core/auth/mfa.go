package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

// NewMFA issues a new MFA event, replacing existing event.
func (s *Svc) NewMFA(ctx context.Context, c core.AuthCredential) (*core.Identity, error) {
	log := logging.FromContext(ctx).WithField("method", "core.user.newmfa")

	ev, err := s.usr.GetMFAEventByID(ctx, c.MFA.EventID)
	if err != nil {
		logging.WithError(err, log).Error("mfa event")
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "login mfa event not found")
		}
		return nil, status.Error(codes.Internal, "mfa authentication failed")
	}
	if ev.Desc != loginEvent || !ev.Active || c.MFA.UserID != ev.UserID {
		return nil, status.Error(codes.InvalidArgument, "invalid mfa retry")
	}

	u, err := s.usr.GetUserByID(ctx, ev.UserID)
	if err != nil {
		logging.WithError(err, log).Error("mfa event")
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "mfa authentication failed")
	}
	if u.Deleted.Valid {
		return nil, status.Error(codes.FailedPrecondition, "user account not active")
	}
	m, err := s.mfa.RestartMFA(ctx, core.MFAChallenge{
		EventID: c.MFA.EventID,
		UserID:  c.MFA.UserID,
	})
	if err != nil {
		return nil, err
	}

	return &core.Identity{
		ID:       u.ID,
		Name:     u.Username,
		OrgID:    u.OrgID,
		EventID:  m.EventID,
		MFA:      m.Type,
		MFATrial: m.Attempt,
		Token:    m.Token,
	}, nil
}

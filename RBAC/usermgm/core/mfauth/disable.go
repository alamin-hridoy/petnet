package mfauth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"

	"brank.as/rbac/serviceutil/logging"
)

func (s *Svc) DisableMFA(ctx context.Context, c core.MFA) (*core.MFA, error) {
	log := logging.FromContext(ctx).WithField("method", "core.mfauth.disable")
	log.WithField("request", c).Trace("received")

	u, err := s.st.GetUserByID(ctx, c.UserID)
	if err != nil {
		logging.WithError(err, log).Error("user from storage")
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, err
	}
	m, err := s.st.GetMFAByID(ctx, c.MFAID)
	if err != nil {
		logging.WithError(err, log).Error("mfa from storage")
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "mfa source not found")
		}
		return nil, err
	}
	if m.UserID != u.ID {
		return nil, status.Error(codes.NotFound, "mfa source not found")
	}
	t, err := s.st.DisableMFA(ctx, m.ID)
	if err != nil {
		logging.WithError(err, log).Error("disable in storage")
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "update failed")
	}
	c.Revoked = t
	return &c, nil
}

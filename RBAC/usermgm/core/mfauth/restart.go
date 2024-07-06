package mfauth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

// RestartMFA cancels the mfa challenge and initiates a new challenge.
func (s *Svc) RestartMFA(ctx context.Context, c core.MFAChallenge) (*core.MFAChallenge, error) {
	log := logging.FromContext(ctx).WithField("method", "core.mfauth.restartmfa")

	ctx, err := s.st.NewTransacton(ctx)
	if err != nil {
		logging.WithError(err, log).Error("storage transaction")
		return nil, status.Error(codes.Internal, "failed to retry MFA")
	}
	defer func() {
		if err := s.st.Rollback(ctx); err != nil {
			logging.WithError(err, log).Error("transaction rollback")
		}
	}()

	ev, err := s.st.DisableMFAEvent(ctx, storage.MFAEvent{
		EventID: c.EventID,
		UserID:  c.UserID,
	})
	if err != nil {
		return nil, err
	}
	c.EventDesc = ev.Desc
	c.Type = ev.MFAType
	c.SourceID = ev.MFAID
	c.Attempt = ev.Attempt + 1

	chg, err := s.InitiateMFA(ctx, c)
	if err != nil {
		return nil, err
	}
	if err := s.st.Commit(ctx); err != nil {
		logging.WithError(err, log).Error("transaction commit")
		return nil, status.Error(codes.Internal, "failed to re-initialize MFA")
	}

	return chg, nil
}

package mfauth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

func (s *Svc) CancelMFA(ctx context.Context, c core.MFAChallenge) error {
	log := logging.FromContext(ctx).WithField("method", "core.mfauth.cancelmfa")

	if _, err := s.st.DisableMFAEvent(ctx, storage.MFAEvent{
		EventID: c.EventID,
		UserID:  c.UserID,
	}); err != nil {
		logging.WithError(err, log).Error("storage cancel")
		if err == storage.NotFound {
			return status.Error(codes.NotFound, "event not found")
		}
		return status.Error(codes.Internal, "failed to cancel event")
	}
	return nil
}

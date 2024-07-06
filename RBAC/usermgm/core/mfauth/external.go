package mfauth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

// ExternalMFA ...
func (s *Svc) ExternalMFA(ctx context.Context, m core.MFAChallenge) (*core.MFAChallenge, error) {
	log := logging.FromContext(ctx).WithField("method", "core.mfauth.externalmfa")
	log.Trace("received")

	ev, err := s.st.GetMFAEventByID(ctx, m.EventID)
	if err != nil {
		if err == storage.NotFound {
			return nil, status.Error(codes.InvalidArgument, "invalid event")
		}
		return nil, err
	}
	if ev.Confirmed.Valid || ev.Expired {
		return nil, status.Error(codes.OutOfRange, "expired/confirmed events cannot be modified")
	}
	switch ev.MFAType {
	case storage.SMS, storage.EMail:
	default:
	}
	ev.Token = m.Token

	e, err := s.st.UpdateMFAEventToken(ctx, *ev)
	if err != nil {
		if err == storage.NotFound {
			return nil, status.Error(codes.InvalidArgument, "invalid event")
		}
		return nil, err
	}
	m.SourceID = e.MFAID
	return &m, nil
}

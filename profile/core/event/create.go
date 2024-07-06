package event

import (
	"context"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) CreateEventData(ctx context.Context, e *storage.EventData) error {
	log := logging.FromContext(ctx)

	if err := s.st.CreateEventData(ctx, e); err != nil {
		if err == storage.Conflict {
			logging.WithError(err, log).Error("event data already exists")
			return status.Error(codes.AlreadyExists, "failed to record event data record")
		}
		logging.WithError(err, log).Error("store event data")
		return status.Error(codes.Internal, "failed to record event data record")
	}
	return nil
}

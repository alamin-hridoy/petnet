package event

import (
	"context"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) GetEventData(ctx context.Context, eid string) (*storage.EventData, error) {
	log := logging.FromContext(ctx)

	e, err := s.st.GetEventData(ctx, eid)
	if err != nil {
		if err == storage.NotFound {
			logging.WithError(err, log).Error("event data doesn't exists")
			return nil, status.Error(codes.NotFound, "failed to get event data record")
		}
		logging.WithError(err, log).Error("get event data")
		return nil, status.Error(codes.Internal, "failed to get event data record")
	}

	// deleting the event data here for now, we might want to make deletion into
	// its own svc endpoint instead to have more control on when to delete
	if err := s.st.DeleteEventData(ctx, eid); err != nil {
		logging.WithError(err, log).Error("delete event data")
		return nil, status.Error(codes.Internal, "failed to get event data record")
	}
	return e, nil
}

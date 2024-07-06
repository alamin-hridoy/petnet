package partner

import (
	"context"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) DeletePartner(ctx context.Context, id string) error {
	log := logging.FromContext(ctx)

	if _, err := s.st.DeletePartner(ctx, id); err != nil {
		if err == storage.NotFound {
			logging.WithError(err, log).Error("partner doesn't exists")
			return status.Error(codes.NotFound, "failed to delete partner record")
		}
		logging.WithError(err, log).Error("delete partner")
		return status.Error(codes.Internal, "failed to delete partner record")
	}
	return nil
}

package partner

import (
	"context"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) EnablePartner(ctx context.Context, oid, svc string) error {
	log := logging.FromContext(ctx)

	err := s.st.EnablePartner(ctx, oid, svc)
	if err != nil {
		logging.WithError(err, log).WithField("org_id", oid).WithField("partner", svc).Error("enabling partner")
		if err == storage.NotFound {
			return status.Error(codes.NotFound, "partner not found")
		}
		return err
	}
	return nil
}

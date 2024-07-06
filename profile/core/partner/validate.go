package partner

import (
	"context"

	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) ValidatePartnerAccess(ctx context.Context, oid, pnr string) error {
	log := logging.FromContext(ctx)

	err := s.st.ValidatePartnerAccess(ctx, oid, pnr)
	if err != nil {
		log.WithField("org_id", oid).WithField("partner", pnr).Error("partner disabled")
		return status.Error(codes.PermissionDenied, "partner disabled for org")
	}
	return nil
}

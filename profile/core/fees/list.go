package fees

import (
	"context"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
)

// ListFees ...
func (s *Svc) ListFees(ctx context.Context, oid string, f storage.LimitOffsetFilter) ([]storage.FeeCommission, error) {
	log := logging.FromContext(ctx)

	fs, err := s.st.ListOrgFees(ctx, oid, f)
	if err != nil {
		logging.WithError(err, log).Error("listing fees")
		return nil, err
	}
	return fs, nil
}

// ListFeesCommissionRate ...
func (s *Svc) ListRates(ctx context.Context, fcid string) ([]storage.Rate, error) {
	log := logging.FromContext(ctx)

	lr, err := s.st.ListFeesCommissionRate(ctx, fcid)
	if err != nil {
		logging.WithError(err, log).Error("listing rates")
		return nil, err
	}
	return lr, nil
}

package fees

import (
	"context"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
)

// UpsertFee ...
func (s *Svc) UpsertFee(ctx context.Context, f storage.FeeCommission) (*storage.FeeCommission, error) {
	log := logging.FromContext(ctx)

	var err error
	fs := &storage.FeeCommission{}
	if f.ID == "" {
		fs, err = s.st.CreateOrgFees(ctx, f)
		if err != nil {
			if err != nil {
				logging.WithError(err, log).Error("creating fee")
				return nil, err
			}
		}
	} else {
		fs, err = s.st.UpsertOrgFees(ctx, f)
		if err != nil {
			logging.WithError(err, log).Error("upserting fee")
			return nil, err
		}
	}
	return fs, nil
}

// UpsertRate ...
func (s *Svc) UpsertRate(ctx context.Context, f storage.Rate) (*storage.Rate, error) {
	log := logging.FromContext(ctx)
	var err error
	fs := &storage.Rate{}
	if f.ID == "" {
		fs, err = s.st.CreateFeeCommissionRate(ctx, f)
		if err != nil {
			if err != nil {
				logging.WithError(err, log).Error("creating rate")
				return nil, err
			}
		}
	} else {
		fs, err = s.st.UpsertRate(ctx, f)
		if err != nil {
			logging.WithError(err, log).Error("upserting rate")
			return nil, err
		}
	}
	return fs, nil
}

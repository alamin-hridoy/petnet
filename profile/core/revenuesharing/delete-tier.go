package revenuesharing

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	"brank.as/petnet/profile/storage"
)

func (s *Svc) DeleteRevenueSharingTier(ctx context.Context, req *rc.DeleteRevenueSharingTierRequest) error {
	if err := s.st.DeleteRevenueSharingTier(ctx, storage.RevenueSharingTier{
		ID:               req.GetID(),
		RevenueSharingID: req.GetRevenueSharingID(),
		MinValue:         req.GetMinValue(),
		MaxValue:         req.GetMaxValue(),
		Amount:           req.GetAmount(),
	}); err != nil {
		return err
	}
	return nil
}

func (s *Svc) DeleteRevenueSharingTierById(ctx context.Context, req *rc.DeleteRevenueSharingTierByIdRequest) error {
	if err := s.st.DeleteRevenueSharingTierById(ctx, storage.RevenueSharingTier{
		ID: req.GetID(),
	}); err != nil {
		return err
	}
	return nil
}

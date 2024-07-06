package revenuesharing

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	"brank.as/petnet/profile/storage"
)

func (s *Svc) UpdateRevenueSharingTier(ctx context.Context, req *rc.UpdateRevenueSharingTierRequest) (*rc.UpdateRevenueSharingTierResponse, error) {
	res, err := s.st.UpdateRevenueSharingTier(ctx, storage.RevenueSharingTier{
		ID:               req.GetID(),
		RevenueSharingID: req.GetRevenueSharingID(),
		MinValue:         req.GetMinValue(),
		MaxValue:         req.GetMaxValue(),
		Amount:           req.GetAmount(),
	})
	if err != nil {
		return nil, err
	}
	return &rc.UpdateRevenueSharingTierResponse{
		ID:               res.ID,
		RevenueSharingID: res.RevenueSharingID,
		MinValue:         res.MinValue,
		MaxValue:         res.MaxValue,
		Amount:           res.Amount,
	}, nil
}

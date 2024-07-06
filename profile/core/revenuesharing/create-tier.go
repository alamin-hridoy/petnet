package revenuesharing

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	"brank.as/petnet/profile/storage"
)

func (s *Svc) CreateRevenueSharingTier(ctx context.Context, req *rc.CreateRevenueSharingTierRequest) (*rc.CreateRevenueSharingTierResponse, error) {
	res, err := s.st.CreateRevenueSharingTier(ctx, storage.RevenueSharingTier{
		RevenueSharingID: req.GetRevenueSharingID(),
		MinValue:         req.GetMinValue(),
		MaxValue:         req.GetMaxValue(),
		Amount:           req.GetAmount(),
	})
	if err != nil {
		return nil, err
	}
	return &rc.CreateRevenueSharingTierResponse{
		ID:               res.ID,
		RevenueSharingID: res.RevenueSharingID,
		MinValue:         res.MinValue,
		MaxValue:         res.MaxValue,
		Amount:           res.Amount,
	}, nil
}

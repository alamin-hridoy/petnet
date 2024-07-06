package revenuesharing

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	"brank.as/petnet/profile/storage"
)

func (s *Svc) GetRevenueSharingTierList(ctx context.Context, req *rc.GetRevenueSharingTierListRequest) (*rc.GetRevenueSharingTierListResponse, error) {
	res, err := s.st.GetRevenueSharingTierList(ctx, storage.RevenueSharingTier{
		ID:               req.GetID(),
		RevenueSharingID: req.GetRevenueSharingID(),
		MinValue:         req.GetMinValue(),
		MaxValue:         req.GetMaxValue(),
		Amount:           req.GetAmount(),
	})
	if err != nil {
		return nil, err
	}
	var rcList []*rc.RevenueSharingTier
	for _, v := range res {
		rcList = append(rcList, &rc.RevenueSharingTier{
			ID:               v.ID,
			RevenueSharingID: v.RevenueSharingID,
			MinValue:         v.MinValue,
			MaxValue:         v.MaxValue,
			Amount:           v.Amount,
		})
	}
	return &rc.GetRevenueSharingTierListResponse{
		Results: rcList,
	}, nil
}

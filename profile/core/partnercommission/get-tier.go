package partnercommission

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/partnercommission"
	"brank.as/petnet/profile/storage"
)

func (s *Svc) GetPartnerCommissionsTierList(ctx context.Context, req *rc.GetPartnerCommissionsTierListRequest) (*rc.GetPartnerCommissionsTierListResponse, error) {
	res, err := s.st.GetPartnerCommissionsTierList(ctx, storage.PartnerCommissionTier{
		ID:                  req.GetID(),
		PartnerCommissionID: req.GetPartnerCommissionID(),
		MinValue:            req.GetMinValue(),
		MaxValue:            req.GetMaxValue(),
		Amount:              req.GetAmount(),
	})
	if err != nil {
		return nil, err
	}
	var rcList []*rc.PartnerCommissionTier
	for _, v := range res {
		rcList = append(rcList, &rc.PartnerCommissionTier{
			ID:                  v.ID,
			PartnerCommissionID: v.PartnerCommissionID,
			MinValue:            v.MinValue,
			MaxValue:            v.MaxValue,
			Amount:              v.Amount,
		})
	}
	return &rc.GetPartnerCommissionsTierListResponse{
		Results: rcList,
	}, nil
}

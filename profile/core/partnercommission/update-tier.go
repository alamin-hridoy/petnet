package partnercommission

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/partnercommission"
	"brank.as/petnet/profile/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) UpdatePartnerCommissionTier(ctx context.Context, req *rc.UpdatePartnerCommissionTierRequest) (*rc.UpdatePartnerCommissionTierResponse, error) {
	if req == nil || req.CommissionTier == nil {
		return nil, status.Error(codes.InvalidArgument, "commission tier is required.")
	}

	tierRes := []*rc.PartnerCommissionTier{}
	for _, r := range req.CommissionTier {
		res, err := s.st.UpdatePartnerCommissionTier(ctx, storage.PartnerCommissionTier{
			ID:                  r.GetID(),
			PartnerCommissionID: r.GetPartnerCommissionID(),
			MinValue:            r.GetMinValue(),
			MaxValue:            r.GetMaxValue(),
			Amount:              r.GetAmount(),
		})
		if err != nil {
			return nil, err
		}
		tierRes = append(tierRes, &rc.PartnerCommissionTier{
			ID:                  res.ID,
			PartnerCommissionID: res.PartnerCommissionID,
			MinValue:            res.MinValue,
			MaxValue:            res.MaxValue,
			Amount:              res.Amount,
		})
	}

	return &rc.UpdatePartnerCommissionTierResponse{
		CommissionTier: tierRes,
	}, nil
}

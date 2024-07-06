package partnercommission

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/partnercommission"
	"brank.as/petnet/profile/storage"
)

func (s *Svc) DeletePartnerCommissionTier(ctx context.Context, req *rc.DeletePartnerCommissionTierRequest) error {
	if err := s.st.DeletePartnerCommissionTier(ctx, storage.PartnerCommissionTier{
		ID:                  req.GetID(),
		PartnerCommissionID: req.GetPartnerCommissionID(),
		MinValue:            req.GetMinValue(),
		MaxValue:            req.GetMaxValue(),
		Amount:              req.GetAmount(),
	}); err != nil {
		return err
	}
	return nil
}

func (s *Svc) DeletePartnerCommissionTierById(ctx context.Context, req *rc.DeletePartnerCommissionTierByIdRequest) error {
	if err := s.st.DeletePartnerCommissionTierById(ctx, storage.PartnerCommissionTier{
		ID: req.GetID(),
	}); err != nil {
		return err
	}
	return nil
}

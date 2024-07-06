package partnercommission

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/partnercommission"
	"brank.as/petnet/profile/storage"
)

func (s *Svc) DeletePartnerCommission(ctx context.Context, req *rc.DeletePartnerCommissionRequest) error {
	if err := s.st.DeletePartnerCommission(ctx, storage.PartnerCommission{
		ID:              req.GetID(),
		Partner:         req.GetPartner(),
		BoundType:       req.GetBoundType().String(),
		RemitType:       req.GetRemitType().String(),
		TransactionType: req.GetTransactionType().String(),
		TierType:        req.GetTierType().String(),
		Amount:          req.GetAmount(),
		StartDate:       req.GetStartDate().AsTime(),
		EndDate:         req.GetEndDate().AsTime(),
		CreatedBy:       req.GetCreatedBy(),
		Created:         req.GetCreated().AsTime(),
	}); err != nil {
		return err
	}
	return nil
}

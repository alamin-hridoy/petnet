package revenuesharing

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	"brank.as/petnet/profile/storage"
)

func (s *Svc) DeleteRevenueSharing(ctx context.Context, req *rc.DeleteRevenueSharingRequest) error {
	if err := s.st.DeleteRevenueSharing(ctx, storage.RevenueSharing{
		ID:              req.GetID(),
		OrgID:           req.GetOrgID(),
		UserID:          req.GetUserID(),
		Partner:         req.GetPartner(),
		BoundType:       req.GetBoundType().String(),
		RemitType:       req.GetRemitType().String(),
		TransactionType: req.GetTransactionType().String(),
		TierType:        req.GetTierType().String(),
		Amount:          req.GetAmount(),
		CreatedBy:       req.GetCreatedBy(),
		Created:         req.GetCreated().AsTime(),
	}); err != nil {
		return err
	}
	return nil
}

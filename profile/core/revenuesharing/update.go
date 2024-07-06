package revenuesharing

import (
	"context"
	"time"

	rc "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/svcutil/mw"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) UpdateRevenueSharing(ctx context.Context, req *rc.UpdateRevenueSharingRequest) (*rc.UpdateRevenueSharingResponse, error) {
	res, err := s.st.UpdateRevenueSharing(ctx, storage.RevenueSharing{
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
		UpdatedBy:       mw.GetUserID(ctx),
		Created:         req.GetCreated().AsTime(),
		Updated:         time.Now(),
	})
	if err != nil {
		return nil, err
	}
	return &rc.UpdateRevenueSharingResponse{
		ID:              res.ID,
		Partner:         res.Partner,
		BoundType:       rc.BoundType(rc.BoundType_value[res.BoundType]),
		RemitType:       rc.RemitType(rc.RemitType_value[res.RemitType]),
		TransactionType: rc.TransactionType(rc.TransactionType_value[res.TransactionType]),
		TierType:        rc.TierType(rc.TierType_value[res.TierType]),
		Amount:          res.Amount,
		CreatedBy:       res.CreatedBy,
		UpdatedBy:       res.UpdatedBy,
		Created:         timestamppb.New(res.Created),
		Updated:         timestamppb.New(res.Updated),
	}, nil
}

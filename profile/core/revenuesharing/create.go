package revenuesharing

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	"brank.as/petnet/profile/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) CreateRevenueSharing(ctx context.Context, req *rc.CreateRevenueSharingRequest) (*rc.CreateRevenueSharingResponse, error) {
	res, err := s.st.CreateRevenueSharing(ctx, storage.RevenueSharing{
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
	})
	if err != nil {
		return nil, err
	}
	return &rc.CreateRevenueSharingResponse{
		ID:              res.ID,
		OrgID:           res.OrgID,
		UserID:          res.UserID,
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

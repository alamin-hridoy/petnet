package revenuesharing

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	"brank.as/petnet/profile/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) GetRevenueSharingList(ctx context.Context, req *rc.GetRevenueSharingListRequest) (*rc.GetRevenueSharingListResponse, error) {
	res, err := s.st.GetRevenueSharingList(ctx, storage.RevenueSharing{
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
	})
	if err != nil {
		return nil, err
	}
	var rcList []*rc.RevenueSharing
	for _, v := range res {
		rcList = append(rcList, &rc.RevenueSharing{
			ID:              v.ID,
			OrgID:           v.OrgID,
			UserID:          v.UserID,
			Partner:         v.Partner,
			BoundType:       rc.BoundType(rc.BoundType_value[v.BoundType]),
			RemitType:       rc.RemitType(rc.RemitType_value[v.RemitType]),
			TransactionType: rc.TransactionType(rc.TransactionType_value[v.TransactionType]),
			TierType:        rc.TierType(rc.TierType_value[v.TierType]),
			Amount:          v.Amount,
			CreatedBy:       v.CreatedBy,
			UpdatedBy:       v.UpdatedBy,
			Created:         timestamppb.New(v.Created),
			Updated:         timestamppb.New(v.Updated),
		})
	}
	return &rc.GetRevenueSharingListResponse{
		Results: rcList,
	}, nil
}

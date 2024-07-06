package partnercommission

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/partnercommission"
	"brank.as/petnet/profile/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) GetPartnerCommissionsList(ctx context.Context, req *rc.GetPartnerCommissionsListRequest) (*rc.GetPartnerCommissionsListResponse, error) {
	res, err := s.st.GetPartnerCommissionsList(ctx, storage.PartnerCommission{
		BoundType: req.GetBoundType().String(),
		RemitType: req.GetRemitType().String(),
		Partner:   req.GetPartner(),
	})
	if err != nil {
		return nil, err
	}
	var rcList []*rc.PartnerCommission
	for _, v := range res {
		rcList = append(rcList, &rc.PartnerCommission{
			ID:              v.ID,
			Partner:         v.Partner,
			BoundType:       rc.BoundType(rc.BoundType_value[v.BoundType]),
			RemitType:       rc.RemitType(rc.RemitType_value[v.RemitType]),
			TransactionType: rc.TransactionType(rc.TransactionType_value[v.TransactionType]),
			TierType:        rc.TierType(rc.TierType_value[v.TierType]),
			Amount:          v.Amount,
			StartDate:       timestamppb.New(v.StartDate),
			EndDate:         timestamppb.New(v.EndDate),
			CreatedBy:       v.CreatedBy,
			UpdatedBy:       v.UpdatedBy,
			Created:         timestamppb.New(v.Created),
			Updated:         timestamppb.New(v.Updated),
		})
	}
	return &rc.GetPartnerCommissionsListResponse{
		Results: rcList,
	}, nil
}

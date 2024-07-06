package partnercommission

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/partnercommission"
	"brank.as/petnet/profile/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) CreatePartnerCommission(ctx context.Context, req *rc.CreatePartnerCommissionRequest) (*rc.CreatePartnerCommissionResponse, error) {
	res, err := s.st.CreatePartnerCommission(ctx, storage.PartnerCommission{
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
	})
	if err != nil {
		return nil, err
	}
	return &rc.CreatePartnerCommissionResponse{
		ID:              res.ID,
		Partner:         res.Partner,
		BoundType:       rc.BoundType(rc.BoundType_value[res.BoundType]),
		RemitType:       rc.RemitType(rc.RemitType_value[res.RemitType]),
		TransactionType: rc.TransactionType(rc.TransactionType_value[res.TransactionType]),
		TierType:        rc.TierType(rc.TierType_value[res.TierType]),
		Amount:          res.Amount,
		StartDate:       timestamppb.New(res.StartDate),
		EndDate:         timestamppb.New(res.EndDate),
		CreatedBy:       res.CreatedBy,
		UpdatedBy:       res.UpdatedBy,
		Created:         timestamppb.New(res.Created),
		Updated:         timestamppb.New(res.Updated),
	}, nil
}

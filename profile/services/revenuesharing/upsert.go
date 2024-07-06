package revenuesharing

import (
	"context"

	"brank.as/petnet/api/storage"
	rc "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) UpsertRevenueSharing(ctx context.Context, req *rc.UpsertRevenueSharingRequest) (*rc.UpsertRevenueSharingResponse, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, required, is.UUID),
		validation.Field(&req.UserID, required, is.UUID),
		validation.Field(&req.Partner, required),
		validation.Field(&req.BoundType, required),
		validation.Field(&req.TransactionType, required),
		validation.Field(&req.RemitType, required),
		validation.Field(&req.TierType, required),
	); err != nil {
		logging.WithError(err, log).Error("upsert revenue sharing validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res := &rc.UpsertRevenueSharingResponse{
		ID:              req.GetID(),
		OrgID:           req.GetOrgID(),
		UserID:          req.GetUserID(),
		Partner:         req.GetPartner(),
		BoundType:       req.GetBoundType(),
		RemitType:       req.GetRemitType(),
		TransactionType: req.GetTransactionType(),
		TierType:        req.GetTierType(),
		Amount:          req.GetAmount(),
		CreatedBy:       req.GetCreatedBy(),
		Created:         req.GetCreated(),
	}
	if req.ID == "" {
		crRes, err := s.core.CreateRevenueSharing(ctx, &rc.CreateRevenueSharingRequest{
			OrgID:           req.GetOrgID(),
			UserID:          req.GetUserID(),
			Partner:         req.GetPartner(),
			BoundType:       req.GetBoundType(),
			RemitType:       req.GetRemitType(),
			TransactionType: req.GetTransactionType(),
			TierType:        req.GetTierType(),
			Amount:          req.GetAmount(),
			CreatedBy:       req.GetCreatedBy(),
			Created:         req.GetCreated(),
		})
		if err != nil {
			if err != storage.Conflict {
				logging.WithError(err, log).Error("creating revenue sharing")
				return nil, err
			}
		}
		res.ID = crRes.ID
	} else {
		_, err := s.core.UpdateRevenueSharing(ctx, &rc.UpdateRevenueSharingRequest{
			ID:              req.GetID(),
			OrgID:           req.GetOrgID(),
			UserID:          req.GetUserID(),
			Partner:         req.GetPartner(),
			BoundType:       req.GetBoundType(),
			RemitType:       req.GetRemitType(),
			TransactionType: req.GetTransactionType(),
			TierType:        req.GetTierType(),
			Amount:          req.GetAmount(),
			CreatedBy:       req.GetCreatedBy(),
			Created:         req.GetCreated(),
		})
		if err != nil {
			logging.WithError(err, log).Error("updating revenue sharing")
			return nil, err
		}
	}
	return res, nil
}

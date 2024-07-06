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

func (s *Svc) UpsertRevenueSharingTier(ctx context.Context, req *rc.UpsertRevenueSharingTierRequest) (*rc.UpsertRevenueSharingTierResponse, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RevenueSharingID, required, is.UUID),
		validation.Field(&req.MinValue, required),
		validation.Field(&req.MaxValue, required),
		validation.Field(&req.Amount, required),
	); err != nil {
		logging.WithError(err, log).Error("upsert revenue sharing tier validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res := &rc.UpsertRevenueSharingTierResponse{
		ID:               req.GetID(),
		RevenueSharingID: req.GetRevenueSharingID(),
		MinValue:         req.GetMinValue(),
		MaxValue:         req.GetMaxValue(),
		Amount:           req.GetAmount(),
	}

	if req.ID == "" {
		id, err := s.NewRevenueSharingTier(ctx, req)
		if err != nil {
			logging.WithError(err, log).Error("creating revenue sharing tier")
			return nil, err
		}
		res.ID = id
		return res, nil
	}

	id, err := s.RevenueSharingTier(ctx, req)
	if err != nil {
		logging.WithError(err, log).Error("updating revenue sharing tier")
		return nil, err
	}
	res.ID = id
	return res, nil
}

func (s *Svc) NewRevenueSharingTier(ctx context.Context, req *rc.UpsertRevenueSharingTierRequest) (string, error) {
	log := logging.FromContext(ctx)
	res, err := s.core.CreateRevenueSharingTier(ctx, &rc.CreateRevenueSharingTierRequest{
		RevenueSharingID: req.GetRevenueSharingID(),
		MinValue:         req.GetMinValue(),
		MaxValue:         req.GetMaxValue(),
		Amount:           req.GetAmount(),
	})
	if err != nil {
		if err == storage.Conflict {
			return s.RevenueSharingTier(ctx, req)
		}
		if err != storage.Conflict {
			logging.WithError(err, log).Error("creating revenue sharing tier")
			return "", err
		}
	}

	return res.ID, nil
}

func (s *Svc) RevenueSharingTier(ctx context.Context, req *rc.UpsertRevenueSharingTierRequest) (string, error) {
	log := logging.FromContext(ctx)
	res, err := s.core.UpdateRevenueSharingTier(ctx, &rc.UpdateRevenueSharingTierRequest{
		ID:               req.GetID(),
		RevenueSharingID: req.GetRevenueSharingID(),
		MinValue:         req.GetMinValue(),
		MaxValue:         req.GetMaxValue(),
		Amount:           req.GetAmount(),
	})
	if err != nil {
		logging.WithError(err, log).Error("updating revenue sharing tier")
		return "", err
	}

	return res.ID, nil
}

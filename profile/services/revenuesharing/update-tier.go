package revenuesharing

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) UpdateRevenueSharingTier(ctx context.Context, req *rc.UpdateRevenueSharingTierRequest) (*rc.UpdateRevenueSharingTierResponse, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ID, required, is.UUID),
		validation.Field(&req.RevenueSharingID, required, is.UUID),
		validation.Field(&req.MinValue, required),
		validation.Field(&req.MaxValue, required),
		validation.Field(&req.Amount, required),
	); err != nil {
		logging.WithError(err, log).Error("update revenue sharing tier validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := s.core.UpdateRevenueSharingTier(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

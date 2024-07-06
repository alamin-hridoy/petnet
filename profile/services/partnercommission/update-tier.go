package partnercommission

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/partnercommission"
	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) UpdatePartnerCommissionTier(ctx context.Context, req *rc.UpdatePartnerCommissionTierRequest) (*rc.UpdatePartnerCommissionTierResponse, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.CommissionTier, required),
	); err != nil {
		logging.WithError(err, log).Error("create partner commission tier validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	for _, ct := range req.CommissionTier {
		if err := validation.ValidateStruct(ct,
			validation.Field(&ct.ID, required, is.UUID),
			validation.Field(&ct.PartnerCommissionID, required, is.UUID),
			validation.Field(&ct.MinValue, required),
			validation.Field(&ct.MaxValue, required),
			validation.Field(&ct.Amount, required),
		); err != nil {
			logging.WithError(err, log).Error("update partner commission tier validation error")
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	res, err := s.core.UpdatePartnerCommissionTier(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

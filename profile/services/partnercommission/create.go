package partnercommission

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/partnercommission"
	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) CreatePartnerCommission(ctx context.Context, req *rc.CreatePartnerCommissionRequest) (*rc.CreatePartnerCommissionResponse, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.BoundType, required),
		validation.Field(&req.RemitType, required),
		validation.Field(&req.Partner, required),
		validation.Field(&req.TransactionType, required),
		validation.Field(&req.TierType, required),
	); err != nil {
		logging.WithError(err, log).Error("create partner commission validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := s.core.CreatePartnerCommission(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

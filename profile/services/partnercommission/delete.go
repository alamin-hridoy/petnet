package partnercommission

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/partnercommission"
	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Svc) DeletePartnerCommission(ctx context.Context, req *rc.DeletePartnerCommissionRequest) (*emptypb.Empty, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Partner, required),
		validation.Field(&req.TransactionType, required),
		validation.Field(&req.RemitType, required),
		validation.Field(&req.BoundType, required),
	); err != nil {
		logging.WithError(err, log).Error("delete partner commission validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.core.DeletePartnerCommission(ctx, req); err != nil {
		return nil, err
	}
	return new(emptypb.Empty), nil
}

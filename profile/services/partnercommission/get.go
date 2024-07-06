package partnercommission

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/partnercommission"
	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) GetPartnerCommissionsList(ctx context.Context, req *rc.GetPartnerCommissionsListRequest) (*rc.GetPartnerCommissionsListResponse, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitType, required),
	); err != nil {
		logging.WithError(err, log).Error("get partner commission validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := s.core.GetPartnerCommissionsList(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

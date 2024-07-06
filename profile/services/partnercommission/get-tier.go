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

func (s *Svc) GetPartnerCommissionsTierList(ctx context.Context, req *rc.GetPartnerCommissionsTierListRequest) (*rc.GetPartnerCommissionsTierListResponse, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.PartnerCommissionID, required, is.UUID),
	); err != nil {
		logging.WithError(err, log).Error("get partner commission tier validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := s.core.GetPartnerCommissionsTierList(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

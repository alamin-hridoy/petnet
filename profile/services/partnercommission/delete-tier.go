package partnercommission

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/partnercommission"
	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Svc) DeletePartnerCommissionTier(ctx context.Context, req *rc.DeletePartnerCommissionTierRequest) (*emptypb.Empty, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.PartnerCommissionID, required, is.UUID),
	); err != nil {
		logging.WithError(err, log).Error("delete partner commission tier validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.core.DeletePartnerCommissionTier(ctx, req); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *Svc) DeletePartnerCommissionTierById(ctx context.Context, req *rc.DeletePartnerCommissionTierByIdRequest) (*emptypb.Empty, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ID, required, is.UUID),
	); err != nil {
		logging.WithError(err, log).Error("delete partner commission tier validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.core.DeletePartnerCommissionTierById(ctx, req); err != nil {
		return nil, err
	}
	return nil, nil
}

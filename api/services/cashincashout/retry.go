package cashincashout

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/util"
	cico "brank.as/petnet/gunk/drp/v1/cashincashout"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) CiCoRetry(ctx context.Context, req *cico.CiCoRetryRequest) (*cico.CiCoRetryResponse, error) {
	if err := s.CiCoRetryValidate(ctx, req); err != nil {
		return nil, util.HandleServiceErr(err)
	}
	res, err := s.core.CiCoRetry(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}
	return res, err
}

func (s *Svc) CiCoRetryValidate(ctx context.Context, req *cico.CiCoRetryRequest) error {
	log := logging.FromContext(ctx)
	if req == nil {
		return status.Error(codes.InvalidArgument, "please provide required fields")
	}
	if err := validation.ValidateStruct(req,
		validation.Field(&req.PartnerCode, validation.Required),
		validation.Field(&req.PetnetTrackingno, validation.Required),
		validation.Field(&req.TrxDate, validation.Required),
		validation.Field(&req.Provider, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}

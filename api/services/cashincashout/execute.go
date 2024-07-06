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

func (s *Svc) CiCoExecute(ctx context.Context, req *cico.CiCoExecuteRequest) (*cico.CiCoExecuteResponse, error) {
	if err := s.CiCoEexcuteValidate(ctx, req); err != nil {
		return nil, util.HandleServiceErr(err)
	}
	res, err := s.core.CiCoExecute(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}
	return res, nil
}

func (s *Svc) CiCoEexcuteValidate(ctx context.Context, req *cico.CiCoExecuteRequest) error {
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

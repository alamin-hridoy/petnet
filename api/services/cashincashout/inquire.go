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

const (
	GCASH      = "GCASH"
	DragonPAY  = "DRAGONPAY"
	PeraHUB    = "PERAHUB"
	DISKARTECH = "DISKARTECH"
	PAYMAYA    = "PAYMAYA"
)

func (s *Svc) CiCoInquire(ctx context.Context, req *cico.CiCoInquireRequest) (*cico.CiCoInquireResponse, error) {
	if err := s.CiCoInquireValidate(ctx, req); err != nil {
		return nil, util.HandleServiceErr(err)
	}
	res, err := s.core.CiCoInquire(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}
	return res, nil
}

func (s *Svc) CiCoInquireValidate(ctx context.Context, req *cico.CiCoInquireRequest) error {
	log := logging.FromContext(ctx)
	if req == nil {
		return status.Error(codes.InvalidArgument, "please provide required fields")
	}
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Provider, validation.Required),
		validation.Field(&req.PartnerCode, validation.When(req.Provider == GCASH || req.Provider == DragonPAY || req.Provider == PeraHUB || req.Provider == DISKARTECH || req.Provider == PAYMAYA, validation.Required)),
		validation.Field(&req.Provider, validation.When(req.Provider == GCASH || req.Provider == DragonPAY || req.Provider == PeraHUB || req.Provider == DISKARTECH || req.Provider == PAYMAYA, validation.Required)),
		validation.Field(&req.TrxType, validation.When(req.Provider == GCASH || req.Provider == DragonPAY || req.Provider == PeraHUB || req.Provider == DISKARTECH || req.Provider == PAYMAYA, validation.Required)),
		validation.Field(&req.ReferenceNumber, validation.When(req.Provider == GCASH || req.Provider == DragonPAY || req.Provider == PeraHUB || req.Provider == DISKARTECH || req.Provider == PAYMAYA, validation.Required)),
		validation.Field(&req.PetnetTrackingno, validation.When(req.Provider == GCASH || req.Provider == DragonPAY || req.Provider == PeraHUB || req.Provider == DISKARTECH || req.Provider == PAYMAYA, validation.Required)),
		validation.Field(&req.ProviderTrackingno, validation.When(req.Provider == PAYMAYA, validation.Required)),
		validation.Field(&req.Message, validation.When(req.Provider != GCASH && req.Provider != DragonPAY && req.Provider != PeraHUB && req.Provider != DISKARTECH && req.Provider != PAYMAYA, validation.Required)),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}

package remittance

import (
	"context"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/util"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) ValidateReceiveMoney(ctx context.Context, req *bpa.ValidateReceiveMoneyRequest) (*bpa.ValidateReceiveMoneyResponse, error) {
	if err := s.ValidateReceiveMoneyValidation(ctx, req); err != nil {
		return nil, err
	}

	res, err := s.remittanceStore.ValidateReceiveMoney(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

func (s *Svc) ValidateReceiveMoneyValidation(ctx context.Context, req *bpa.ValidateReceiveMoneyRequest) error {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Phrn, validation.Required),
		validation.Field(&req.PrincipalAmount, validation.Required),
		validation.Field(&req.IsoOriginatingCountry, validation.Required),
		validation.Field(&req.IsoDestinationCountry, validation.Required),
		validation.Field(&req.SenderLastName, validation.Required),
		validation.Field(&req.SenderFirstName, validation.Required),
		validation.Field(&req.SenderMiddleName, validation.Required),
		validation.Field(&req.ReceiverLastName, validation.Required),
		validation.Field(&req.ReceiverFirstName, validation.Required),
		validation.Field(&req.ReceiverMiddleName, validation.Required),
		validation.Field(&req.PayoutPartnerCode, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return status.Error(http.StatusUnprocessableEntity, err.Error())
	}

	return nil
}

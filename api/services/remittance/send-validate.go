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

func (s *Svc) ValidateSendMoney(ctx context.Context, req *bpa.ValidateSendMoneyRequest) (*bpa.ValidateSendMoneyResponse, error) {
	if err := s.ValidateSendMoneyValidate(ctx, req); err != nil {
		return nil, err
	}

	res, err := s.remittanceStore.ValidateSendMoney(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

func (s *Svc) ValidateSendMoneyValidate(ctx context.Context, req *bpa.ValidateSendMoneyRequest) error {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ConversionRate, validation.Required),
		validation.Field(&req.IsoCurrency, validation.Required),
		validation.Field(&req.IsoDestinationCountry, validation.Required),
		validation.Field(&req.IsoOriginatingCountry, validation.Required),
		validation.Field(&req.PartnerReferenceNumber, validation.Required),
		validation.Field(&req.PrincipalAmount, validation.Required),
		validation.Field(&req.ReceiverFirstName, validation.Required),
		validation.Field(&req.ReceiverLastName, validation.Required),
		validation.Field(&req.ReceiverMiddleName, validation.Required),
		validation.Field(&req.SendPartnerCode, validation.Required),
		validation.Field(&req.SenderBirthCountry, validation.Required),
		validation.Field(&req.SenderBirthDate, validation.Required),
		validation.Field(&req.SenderBirthPlace, validation.Required),
		validation.Field(&req.SenderEmploymentNature, validation.Required),
		validation.Field(&req.SenderFirstName, validation.Required),
		validation.Field(&req.SenderGender, validation.Required),
		validation.Field(&req.SenderLastName, validation.Required),
		validation.Field(&req.SenderMiddleName, validation.Required),
		validation.Field(&req.SenderOccupation, validation.Required),
		validation.Field(&req.SenderPurpose, validation.Required),
		validation.Field(&req.SenderRelationship, validation.Required),
		validation.Field(&req.SenderSourceOfFund, validation.Required),
		validation.Field(&req.ServiceFee, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return status.Error(http.StatusUnprocessableEntity, err.Error())
	}

	return nil
}

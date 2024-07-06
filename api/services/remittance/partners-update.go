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

func (s *Svc) PartnersUpdate(ctx context.Context, req *bpa.PartnersUpdateRequest) (*bpa.PartnersUpdateResponse, error) {
	if err := s.PartnersUpdateValidate(ctx, req); err != nil {
		return nil, err
	}

	res, err := s.remittanceStore.PartnersUpdate(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

func (s *Svc) PartnersUpdateValidate(ctx context.Context, req *bpa.PartnersUpdateRequest) error {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.PartnerCode, validation.Required),
		validation.Field(&req.PartnerName, validation.Required),
		validation.Field(&req.Service, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return status.Error(http.StatusUnprocessableEntity, err.Error())
	}

	return nil
}

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

func (s *Svc) Inquire(ctx context.Context, req *bpa.InquireRequest) (*bpa.InquireResponse, error) {
	if err := s.InquireValidate(ctx, req); err != nil {
		return nil, err
	}

	res, err := s.remittanceStore.Inquire(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

func (s *Svc) InquireValidate(ctx context.Context, req *bpa.InquireRequest) error {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Phrn, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return status.Error(http.StatusUnprocessableEntity, err.Error())
	}

	return nil
}

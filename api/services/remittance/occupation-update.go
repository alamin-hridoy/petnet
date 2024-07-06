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

func (s *Svc) OccupationUpdate(ctx context.Context, req *bpa.OccupationUpdateRequest) (*bpa.OccupationUpdateResponse, error) {
	if err := s.OccupationUpdateValidation(ctx, req); err != nil {
		return nil, err
	}

	res, err := s.remittanceStore.OccupationUpdate(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

func (s *Svc) OccupationUpdateValidation(ctx context.Context, req *bpa.OccupationUpdateRequest) error {
	log := logging.FromContext(ctx)
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Occupation, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return status.Error(http.StatusUnprocessableEntity, err.Error())
	}

	return nil
}

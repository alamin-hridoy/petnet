package remittance

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/petnet/api/integration/perahub"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) PurposeOfRemittanceUpdate(ctx context.Context, req *bpa.PurposeOfRemittanceUpdateRequest) (res *bpa.PurposeOfRemittanceUpdateResponse, err error) {
	log := logging.FromContext(ctx)
	rvsm, err := s.ph.PurposeOfRemittanceUpdate(ctx, perahub.PurposeOfRemittanceUpdateReq{
		PurposeOfRemittance: req.GetPurposeOfRemittance(),
	}, req.ID)
	if err != nil {
		logging.WithError(err, log).Error("PurposeOfRemittanceUpdate error")
		return nil, handlePerahubError(err)
	}

	return &bpa.PurposeOfRemittanceUpdateResponse{
		Code:    int32(rvsm.Code),
		Message: rvsm.Message,
		Result: &bpa.PurposeOfRemittanceUpdateResult{
			ID:                  int32(rvsm.Result.ID),
			PurposeOfRemittance: rvsm.Result.PurposeOfRemittance,
			CreatedAt:           timestamppb.New(rvsm.Result.CreatedAt),
			UpdatedAt:           timestamppb.New(rvsm.Result.UpdatedAt),
		},
	}, nil
}

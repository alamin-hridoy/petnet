package remittance

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/petnet/api/integration/perahub"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) PurposeOfRemittanceCreate(ctx context.Context, req *bpa.PurposeOfRemittanceCreateRequest) (res *bpa.PurposeOfRemittanceCreateResponse, err error) {
	log := logging.FromContext(ctx)
	rvsm, err := s.ph.PurposeOfRemittanceCreate(ctx, perahub.PurposeOfRemittanceCreateReq{
		PurposeOfRemittance: req.GetPurposeOfRemittance(),
	})
	if err != nil {
		logging.WithError(err, log).Error("PurposeOfRemittanceCreate error")
		return nil, handlePerahubError(err)
	}

	return &bpa.PurposeOfRemittanceCreateResponse{
		Code:    int32(rvsm.Code),
		Message: rvsm.Message,
		Result: &bpa.PurposeOfRemittanceCreateResult{
			ID:                  int32(rvsm.Result.ID),
			PurposeOfRemittance: rvsm.Result.PurposeOfRemittance,
			CreatedAt:           timestamppb.New(rvsm.Result.CreatedAt),
			UpdatedAt:           timestamppb.New(rvsm.Result.UpdatedAt),
		},
	}, nil
}

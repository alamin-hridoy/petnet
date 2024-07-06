package remittance

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/petnet/api/integration/perahub"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) OccupationCreate(ctx context.Context, req *bpa.OccupationCreateRequest) (res *bpa.OccupationCreateResponse, err error) {
	log := logging.FromContext(ctx)
	rvsm, err := s.ph.OccupationCreate(ctx, perahub.OccupationCreateReq{
		Occupation: req.GetOccupation(),
	})
	if err != nil {
		logging.WithError(err, log).Error("OccupationCreate error")
		return nil, handlePerahubError(err)
	}

	return &bpa.OccupationCreateResponse{
		Code:    int32(rvsm.Code),
		Message: rvsm.Message,
		Result: &bpa.OccupationCreateResult{
			ID:         int32(rvsm.Result.ID),
			Occupation: rvsm.Result.Occupation,
			CreatedAt:  timestamppb.New(rvsm.Result.CreatedAt),
			UpdatedAt:  timestamppb.New(rvsm.Result.UpdatedAt),
		},
	}, nil
}

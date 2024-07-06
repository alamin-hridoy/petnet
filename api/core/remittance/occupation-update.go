package remittance

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/petnet/api/integration/perahub"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) OccupationUpdate(ctx context.Context, req *bpa.OccupationUpdateRequest) (res *bpa.OccupationUpdateResponse, err error) {
	log := logging.FromContext(ctx)
	rvsm, err := s.ph.OccupationUpdate(ctx, perahub.OccupationUpdateReq{
		Occupation: req.GetOccupation(),
	}, req.ID)
	if err != nil {
		logging.WithError(err, log).Error("OccupationUpdate error")
		return nil, handlePerahubError(err)
	}

	return &bpa.OccupationUpdateResponse{
		Code:    int32(rvsm.Code),
		Message: rvsm.Message,
		Result: &bpa.OccupationUpdateResult{
			ID:         int32(rvsm.Result.ID),
			Occupation: rvsm.Result.Occupation,
			CreatedAt:  timestamppb.New(rvsm.Result.CreatedAt),
			UpdatedAt:  timestamppb.New(rvsm.Result.UpdatedAt),
		},
	}, nil
}

package remittance

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/petnet/api/integration/perahub"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) SourceOfFundUpdate(ctx context.Context, req *bpa.SourceOfFundUpdateRequest) (res *bpa.SourceOfFundUpdateResponse, err error) {
	log := logging.FromContext(ctx)
	rvsm, err := s.ph.SourceOfFundUpdate(ctx, perahub.SourceOfFundUpdateReq{
		SourceOfFund: req.GetSourceOfFund(),
	}, req.ID)
	if err != nil {
		logging.WithError(err, log).Error("SourceOfFundUpdate error")
		return nil, handlePerahubError(err)
	}

	return &bpa.SourceOfFundUpdateResponse{
		Code:    int32(rvsm.Code),
		Message: rvsm.Message,
		Result: &bpa.SourceOfFundUpdateResult{
			ID:           int32(rvsm.Result.ID),
			SourceOfFund: rvsm.Result.SourceOfFund,
			CreatedAt:    timestamppb.New(rvsm.Result.CreatedAt),
			UpdatedAt:    timestamppb.New(rvsm.Result.UpdatedAt),
		},
	}, nil
}

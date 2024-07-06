package remittance

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/petnet/api/integration/perahub"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) SourceOfFundCreate(ctx context.Context, req *bpa.SourceOfFundCreateRequest) (res *bpa.SourceOfFundCreateResponse, err error) {
	log := logging.FromContext(ctx)
	rvsm, err := s.ph.SourceOfFundCreate(ctx, perahub.SourceOfFundCreateReq{
		SourceOfFund: req.GetSourceOfFund(),
	})
	if err != nil {
		logging.WithError(err, log).Error("SourceOfFundCreate error")
		return nil, handlePerahubError(err)
	}

	return &bpa.SourceOfFundCreateResponse{
		Code:    int32(rvsm.Code),
		Message: rvsm.Message,
		Result: &bpa.SourceOfFundCreateResult{
			ID:           int32(rvsm.Result.ID),
			SourceOfFund: rvsm.Result.SourceOfFund,
			CreatedAt:    timestamppb.New(rvsm.Result.CreatedAt),
			UpdatedAt:    timestamppb.New(rvsm.Result.UpdatedAt),
		},
	}, nil
}

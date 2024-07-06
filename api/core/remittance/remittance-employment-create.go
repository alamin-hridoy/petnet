package remittance

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/petnet/api/integration/perahub"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) RemittanceEmploymentCreate(ctx context.Context, req *bpa.RemittanceEmploymentCreateRequest) (res *bpa.RemittanceEmploymentCreateResponse, err error) {
	log := logging.FromContext(ctx)
	rvsm, err := s.ph.RemittanceEmploymentCreate(ctx, perahub.RemittanceEmploymentCreateReq{
		Employment:       req.GetEmployment(),
		EmploymentNature: req.GetEmploymentNature(),
	})
	if err != nil {
		logging.WithError(err, log).Error("RemittanceEmploymentCreate error")
		return nil, handlePerahubError(err)
	}

	return &bpa.RemittanceEmploymentCreateResponse{
		Code:    int32(rvsm.Code),
		Message: rvsm.Message,
		Result: &bpa.RemittanceEmploymentCreateResult{
			ID:               int32(rvsm.Result.ID),
			EmploymentNature: rvsm.Result.EmploymentNature,
			CreatedAt:        timestamppb.New(rvsm.Result.CreatedAt),
			UpdatedAt:        timestamppb.New(rvsm.Result.UpdatedAt),
		},
	}, nil
}

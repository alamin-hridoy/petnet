package remittance

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/petnet/api/integration/perahub"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) RemittanceEmploymentUpdate(ctx context.Context, req *bpa.RemittanceEmploymentUpdateRequest) (res *bpa.RemittanceEmploymentUpdateResponse, err error) {
	log := logging.FromContext(ctx)
	rvsm, err := s.ph.RemittanceEmploymentUpdate(ctx, perahub.RemittanceEmploymentUpdateReq{
		Employment:       req.GetEmployment(),
		EmploymentNature: req.GetEmploymentNature(),
	}, req.ID)
	if err != nil {
		logging.WithError(err, log).Error("RemittanceEmploymentUpdate error")
		return nil, handlePerahubError(err)
	}

	return &bpa.RemittanceEmploymentUpdateResponse{
		Code:    int32(rvsm.Code),
		Message: rvsm.Message,
		Result: &bpa.RemittanceEmploymentUpdateResult{
			ID:               int32(rvsm.Result.ID),
			EmploymentNature: rvsm.Result.EmploymentNature,
			CreatedAt:        timestamppb.New(rvsm.Result.CreatedAt),
			UpdatedAt:        timestamppb.New(rvsm.Result.UpdatedAt),
		},
	}, nil
}

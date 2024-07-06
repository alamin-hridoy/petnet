package remittance

import (
	"context"

	"google.golang.org/grpc/codes"
	timestamps "google.golang.org/protobuf/types/known/timestamppb"

	coreerror "brank.as/petnet/api/core/error"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) EmploymentGrid(ctx context.Context) (res *bpa.EmploymentGridResponse, err error) {
	log := logging.FromContext(ctx)
	um, err := s.ph.RemittanceEmploymentGrid(ctx)
	if err != nil {
		logging.WithError(err, log).Error("EmploymentGrid error")
		return nil, handlePerahubError(err)
	}

	if um == nil || len(um.Result) == 0 {
		return nil, coreerror.NewCoreError(codes.NotFound, coreerror.MsgNotFound)
	}

	results := make([]*bpa.EmploymentGridResult, 0, len(um.Result))
	for _, v := range um.Result {
		results = append(results, &bpa.EmploymentGridResult{
			ID:               int32(v.ID),
			EmploymentNature: v.EmploymentNature,
			CreatedAt:        timestamps.New(v.CreatedAt),
			UpdatedAt:        timestamps.New(v.UpdatedAt),
			DeletedAt:        timestamps.New(v.DeletedAt),
		})
	}

	return &bpa.EmploymentGridResponse{
		Code:    int32(um.Code),
		Message: um.Message,
		Result:  results,
	}, nil
}

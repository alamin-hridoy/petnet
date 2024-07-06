package remittance

import (
	"context"

	"google.golang.org/grpc/codes"
	timestamps "google.golang.org/protobuf/types/known/timestamppb"

	coreerror "brank.as/petnet/api/core/error"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) EmploymentGet(ctx context.Context, req *bpa.EmploymentGetRequest) (res *bpa.EmploymentGetResponse, err error) {
	log := logging.FromContext(ctx)
	um, err := s.ph.RemittanceEmploymentGet(ctx, req.GetID())
	if err != nil {
		logging.WithError(err, log).Error("EmploymentGet error")
		return nil, handlePerahubError(err)
	}

	if um == nil || um.Result == nil {
		return nil, coreerror.NewCoreError(codes.NotFound, coreerror.MsgNotFound)
	}

	return &bpa.EmploymentGetResponse{
		Code:    int32(um.Code),
		Message: um.Message,
		Result: &bpa.EmploymentGetResult{
			ID:               int32(um.Result.ID),
			EmploymentNature: um.Result.EmploymentNature,
			CreatedAt:        timestamps.New(um.Result.CreatedAt),
			UpdatedAt:        timestamps.New(um.Result.UpdatedAt),
			DeletedAt:        timestamps.New(um.Result.DeletedAt),
		},
	}, nil
}

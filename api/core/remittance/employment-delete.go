package remittance

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	coreerror "brank.as/petnet/api/core/error"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) RemittanceEmploymentDelete(ctx context.Context, req *bpa.RemittanceEmploymentDeleteRequest) (res *bpa.RemittanceEmploymentDeleteResponse, err error) {
	log := logging.FromContext(ctx)
	um, err := s.ph.RemittanceEmploymentDelete(ctx, req.GetID())
	if err != nil {
		logging.WithError(err, log).Error("RemittanceEmploymentDelete error")
		return nil, handlePerahubError(err)
	}

	if um == nil || um.Result == nil {
		return nil, coreerror.NewCoreError(codes.NotFound, coreerror.MsgNotFound)
	}

	return &bpa.RemittanceEmploymentDeleteResponse{
		Code:    int32(um.Code),
		Message: um.Message,
		Result: &bpa.RemittanceEmploymentDeleteResult{
			ID:               int32(um.Result.ID),
			EmploymentNature: um.Result.EmploymentNature,
			CreatedAt:        timestamppb.New(um.Result.CreatedAt),
			UpdatedAt:        timestamppb.New(um.Result.UpdatedAt),
			DeletedAt:        timestamppb.New(um.Result.DeletedAt),
		},
	}, nil
}

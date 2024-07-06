package remittance

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	coreerror "brank.as/petnet/api/core/error"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) OccupationGet(ctx context.Context, req *bpa.OccupationGetRequest) (res *bpa.OccupationGetResponse, err error) {
	log := logging.FromContext(ctx)
	um, err := s.ph.OccupationGet(ctx, req.GetID())
	if err != nil {
		logging.WithError(err, log).Error("OccupationGet error")
		return nil, handlePerahubError(err)
	}

	if um == nil || um.Result == nil {
		return nil, coreerror.NewCoreError(codes.NotFound, coreerror.MsgNotFound)
	}

	return &bpa.OccupationGetResponse{
		Code:    int32(um.Code),
		Message: um.Message,
		Result: &bpa.OccupationGetResult{
			ID:         int32(um.Result.ID),
			Occupation: um.Result.Occupation,
			CreatedAt:  timestamppb.New(um.Result.CreatedAt),
			UpdatedAt:  timestamppb.New(um.Result.UpdatedAt),
			DeletedAt:  timestamppb.New(um.Result.DeletedAt),
		},
	}, nil
}

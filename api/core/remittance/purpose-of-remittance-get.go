package remittance

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	coreerror "brank.as/petnet/api/core/error"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) PurposeOfRemittanceGet(ctx context.Context, req *bpa.PurposeOfRemittanceGetRequest) (res *bpa.PurposeOfRemittanceGetResponse, err error) {
	log := logging.FromContext(ctx)
	um, err := s.ph.PurposeOfRemittanceGet(ctx, req.GetID())
	if err != nil {
		logging.WithError(err, log).Error("PurposeOfRemittanceGet error")
		return nil, handlePerahubError(err)
	}

	if um == nil || um.Result == nil {
		return nil, coreerror.NewCoreError(codes.NotFound, coreerror.MsgNotFound)
	}

	return &bpa.PurposeOfRemittanceGetResponse{
		Code:    int32(um.Code),
		Message: um.Message,
		Result: &bpa.PurposeOfRemittanceGetResult{
			ID:                  int32(um.Result.ID),
			PurposeOfRemittance: um.Result.PurposeOfRemittance,
			CreatedAt:           timestamppb.New(um.Result.CreatedAt),
			UpdatedAt:           timestamppb.New(um.Result.UpdatedAt),
			DeletedAt:           timestamppb.New(um.Result.DeletedAt),
		},
	}, nil
}

package remittance

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	coreerror "brank.as/petnet/api/core/error"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) PartnersGet(ctx context.Context, req *bpa.PartnersGetRequest) (res *bpa.PartnersGetResponse, err error) {
	log := logging.FromContext(ctx)
	rpc, err := s.ph.RemittancePartnersGet(ctx, req.GetID())
	if err != nil {
		logging.WithError(err, log).Error("PartnersGet error")
		return nil, handlePerahubError(err)
	}

	if rpc == nil || rpc.Result == nil {
		return nil, coreerror.NewCoreError(codes.NotFound, coreerror.MsgNotFound)
	}

	return &bpa.PartnersGetResponse{
		Code:    int32(rpc.Code),
		Message: rpc.Message,
		Result: &bpa.PartnersGetResult{
			ID:           int32(rpc.Result.ID),
			PartnerCode:  rpc.Result.PartnerCode,
			PartnerName:  rpc.Result.PartnerName,
			ClientSecret: rpc.Result.ClientSecret,
			Status:       int32(rpc.Result.Status),
			CreatedAt:    timestamppb.New(rpc.Result.CreatedAt),
			UpdatedAt:    timestamppb.New(rpc.Result.UpdatedAt),
			DeletedAt:    rpc.Result.DeletedAt,
		},
	}, nil
}

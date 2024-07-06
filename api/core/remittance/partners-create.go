package remittance

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/petnet/api/integration/perahub"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) PartnersCreate(ctx context.Context, req *bpa.PartnersCreateRequest) (res *bpa.PartnersCreateResponse, err error) {
	log := logging.FromContext(ctx)
	rpc, err := s.ph.RemittancePartnersCreate(ctx, perahub.RemittancePartnersCreateReq{
		PartnerCode: req.GetPartnerCode(),
		PartnerName: req.GetPartnerName(),
		Service:     req.GetService(),
	})
	if err != nil {
		logging.WithError(err, log).Error("PartnersCreate error")
		return nil, handlePerahubError(err)
	}

	return &bpa.PartnersCreateResponse{
		Code:    int32(rpc.Code),
		Message: rpc.Message,
		Result: &bpa.PartnersCreateResult{
			ID:           int32(rpc.Result.ID),
			PartnerCode:  rpc.Result.PartnerCode,
			PartnerName:  rpc.Result.PartnerName,
			ClientSecret: rpc.Result.ClientSecret,
			CreatedAt:    timestamppb.New(rpc.Result.CreatedAt),
			UpdatedAt:    timestamppb.New(rpc.Result.UpdatedAt),
		},
	}, nil
}

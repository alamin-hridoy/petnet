package remittance

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/petnet/api/integration/perahub"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) PartnersUpdate(ctx context.Context, req *bpa.PartnersUpdateRequest) (res *bpa.PartnersUpdateResponse, err error) {
	log := logging.FromContext(ctx)
	rpu, err := s.ph.RemittancePartnersUpdate(ctx, perahub.RemittancePartnersUpdateReq{
		ID:          req.GetID(),
		PartnerCode: req.GetPartnerCode(),
		PartnerName: req.GetPartnerName(),
		Service:     req.GetService(),
	}, req.GetID())
	if err != nil {
		logging.WithError(err, log).Error("PartnersUpdate error")
		return nil, handlePerahubError(err)
	}

	return &bpa.PartnersUpdateResponse{
		Code:    int32(rpu.Code),
		Message: rpu.Message,
		Result: &bpa.PartnersUpdateResult{
			ID:           int32(rpu.Result.ID),
			PartnerCode:  rpu.Result.PartnerCode,
			PartnerName:  rpu.Result.PartnerName,
			ClientSecret: rpu.Result.ClientSecret,
			Status:       int32(rpu.Result.Status),
			CreatedAt:    timestamppb.New(rpu.Result.CreatedAt),
			UpdatedAt:    timestamppb.New(rpu.Result.UpdatedAt),
			DeletedAt:    rpu.Result.DeletedAt,
		},
	}, nil
}

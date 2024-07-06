package remittance

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/petnet/api/integration/perahub"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) PartnersDelete(ctx context.Context, req *bpa.PartnersDeleteRequest) (res *bpa.PartnersDeleteResponse, err error) {
	log := logging.FromContext(ctx)
	um, err := s.ph.RemittancePartnersDelete(ctx, perahub.RemittancePartnersDeleteReq{
		ID:          req.GetID(),
		PartnerCode: req.GetPartnerCode(),
		PartnerName: req.GetPartnerName(),
	})
	if err != nil {
		logging.WithError(err, log).Error("PartnersDelete error")
		return nil, handlePerahubError(err)
	}

	return &bpa.PartnersDeleteResponse{
		Code:    int32(um.Code),
		Message: um.Message,
		Result: &bpa.PartnersDeleteResult{
			ID:           int32(um.Result.ID),
			PartnerCode:  um.Result.PartnerCode,
			PartnerName:  um.Result.PartnerName,
			ClientSecret: um.Result.ClientSecret,
			CreatedAt:    timestamppb.New(um.Result.CreatedAt),
			UpdatedAt:    timestamppb.New(um.Result.UpdatedAt),
			DeletedAt:    timestamppb.New(um.Result.DeletedAt),
		},
	}, nil
}

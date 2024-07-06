package remittance

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	coreerror "brank.as/petnet/api/core/error"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) PartnersGrid(ctx context.Context) (res *bpa.PartnersGridResponse, err error) {
	log := logging.FromContext(ctx)
	um, err := s.ph.RemittancePartnersGrid(ctx)
	if err != nil {
		logging.WithError(err, log).Error("PartnersGrid error")
		return nil, handlePerahubError(err)
	}

	if um == nil || len(um.Result) == 0 {
		return nil, coreerror.NewCoreError(codes.NotFound, coreerror.MsgNotFound)
	}

	pg := make([]*bpa.PartnersGridResult, 0, len(um.Result))
	for _, v := range um.Result {
		pg = append(pg, &bpa.PartnersGridResult{
			ID:           int32(v.ID),
			PartnerCode:  v.PartnerCode,
			PartnerName:  v.PartnerName,
			ClientSecret: v.ClientSecret,
			Status:       int32(v.Status),
			CreatedAt:    timestamppb.New(v.CreatedAt),
			UpdatedAt:    timestamppb.New(v.UpdatedAt),
			DeletedAt:    timestamppb.New(v.DeletedAt),
		})
	}

	return &bpa.PartnersGridResponse{
		Code:    int32(um.Code),
		Message: um.Message,
		Result:  pg,
	}, nil
}

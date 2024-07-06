package remittance

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) PurposeOfRemittanceDelete(ctx context.Context, req *bpa.PurposeOfRemittanceDeleteRequest) (res *bpa.PurposeOfRemittanceDeleteResponse, err error) {
	log := logging.FromContext(ctx)
	um, err := s.ph.PurposeOfRemittanceDelete(ctx, req.GetID())
	if err != nil {
		logging.WithError(err, log).Error("PurposeOfRemittanceDelete error")
		return nil, handlePerahubError(err)
	}

	return &bpa.PurposeOfRemittanceDeleteResponse{
		Code:    int32(um.Code),
		Message: um.Message,
		Result: &bpa.PurposeOfRemittanceDeleteResult{
			ID:                  int32(um.Result.ID),
			PurposeOfRemittance: um.Result.PurposeOfRemittance,
			CreatedAt:           timestamppb.New(um.Result.CreatedAt),
			UpdatedAt:           timestamppb.New(um.Result.UpdatedAt),
			DeletedAt:           timestamppb.New(um.Result.DeletedAt),
		},
	}, nil
}

package remittance

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) OccupationDelete(ctx context.Context, req *bpa.OccupationDeleteRequest) (res *bpa.OccupationDeleteResponse, err error) {
	log := logging.FromContext(ctx)
	um, err := s.ph.OccupationDelete(ctx, req.GetID())
	if err != nil {
		logging.WithError(err, log).Error("OccupationDelete error")
		return nil, handlePerahubError(err)
	}

	return &bpa.OccupationDeleteResponse{
		Code:    int32(um.Code),
		Message: um.Message,
		Result: &bpa.OccupationDeleteResult{
			ID:         int32(um.Result.ID),
			Occupation: um.Result.Occupation,
			CreatedAt:  timestamppb.New(um.Result.CreatedAt),
			UpdatedAt:  timestamppb.New(um.Result.UpdatedAt),
			DeletedAt:  timestamppb.New(um.Result.DeletedAt),
		},
	}, nil
}

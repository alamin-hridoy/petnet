package remittance

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) SourceOfFundDelete(ctx context.Context, req *bpa.SourceOfFundDeleteRequest) (res *bpa.SourceOfFundDeleteResponse, err error) {
	log := logging.FromContext(ctx)
	um, err := s.ph.SourceOfFundDelete(ctx, req.GetID())
	if err != nil {
		logging.WithError(err, log).Error("SourceOfFundDelete error")
		return nil, handlePerahubError(err)
	}

	return &bpa.SourceOfFundDeleteResponse{
		Code:    int32(um.Code),
		Message: um.Message,
		Result: &bpa.SourceOfFundDeleteResult{
			ID:           int32(um.Result.ID),
			SourceOfFund: um.Result.SourceOfFund,
			CreatedAt:    timestamppb.New(um.Result.CreatedAt),
			UpdatedAt:    timestamppb.New(um.Result.UpdatedAt),
			DeletedAt:    timestamppb.New(um.Result.DeletedAt),
		},
	}, nil
}

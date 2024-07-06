package remittance

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	coreerror "brank.as/petnet/api/core/error"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) SourceOfFundGrid(ctx context.Context) (res *bpa.SourceOfFundGridResponse, err error) {
	log := logging.FromContext(ctx)
	um, err := s.ph.SourceOfFundGrid(ctx)
	if err != nil {
		logging.WithError(err, log).Error("SourceOfFundGrid error")
		return nil, handlePerahubError(err)
	}

	if um == nil || len(um.Result) == 0 {
		return nil, coreerror.NewCoreError(codes.NotFound, coreerror.MsgNotFound)
	}

	br := make([]*bpa.SourceOfFundGridResult, 0, len(um.Result))
	for _, v := range um.Result {
		br = append(br, &bpa.SourceOfFundGridResult{
			ID:           string(v.ID),
			SourceOfFund: v.SourceOfFund,
			CreatedAt:    timestamppb.New(v.CreatedAt),
			UpdatedAt:    timestamppb.New(v.UpdatedAt),
			DeletedAt:    timestamppb.New(v.DeletedAt),
		})
	}

	return &bpa.SourceOfFundGridResponse{
		Code:    int32(um.Code),
		Message: um.Message,
		Result:  br,
	}, nil
}

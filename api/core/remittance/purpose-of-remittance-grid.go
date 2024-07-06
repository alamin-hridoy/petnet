package remittance

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	coreerror "brank.as/petnet/api/core/error"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) PurposeOfRemittanceGrid(ctx context.Context) (res *bpa.PurposeOfRemittanceGridResponse, err error) {
	log := logging.FromContext(ctx)
	um, err := s.ph.PurposeOfRemittanceGrid(ctx)
	if err != nil {
		logging.WithError(err, log).Error("PurposeOfRemittanceGrid error")
		return nil, handlePerahubError(err)
	}

	if um == nil || len(um.Result) == 0 {
		return nil, coreerror.NewCoreError(codes.NotFound, coreerror.MsgNotFound)
	}

	br := make([]*bpa.PurposeOfRemittanceGridResult, 0, len(um.Result))
	for _, v := range um.Result {
		br = append(br, &bpa.PurposeOfRemittanceGridResult{
			ID:                  string(v.ID),
			PurposeOfRemittance: v.PurposeOfRemittance,
			CreatedAt:           timestamppb.New(v.CreatedAt),
			UpdatedAt:           timestamppb.New(v.UpdatedAt),
			DeletedAt:           timestamppb.New(v.DeletedAt),
		})
	}

	return &bpa.PurposeOfRemittanceGridResponse{
		Code:    int32(um.Code),
		Message: um.Message,
		Result:  br,
	}, nil
}

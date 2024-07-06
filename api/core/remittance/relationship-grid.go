package remittance

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	coreerror "brank.as/petnet/api/core/error"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) RelationshipGrid(ctx context.Context, em *emptypb.Empty) (res *bpa.RelationshipGridResponse, err error) {
	log := logging.FromContext(ctx)
	um, err := s.ph.RemittanceRelationshiptGrid(ctx)
	if err != nil {
		logging.WithError(err, log).Error("RelationshipGrid error")
		return nil, handlePerahubError(err)
	}

	if um == nil || len(um.Result) == 0 {
		return nil, coreerror.NewCoreError(codes.NotFound, coreerror.MsgNotFound)
	}

	results := make([]*bpa.RelationshipGridResult, 0, len(um.Result))
	for _, v := range um.Result {
		results = append(results, &bpa.RelationshipGridResult{
			ID:           int32(v.ID),
			Relationship: v.Relationship,
			CreatedAt:    timestamppb.New(v.CreatedAt),
			UpdatedAt:    timestamppb.New(v.UpdatedAt),
			DeletedAt:    timestamppb.New(v.DeletedAt),
		})
	}

	return &bpa.RelationshipGridResponse{
		Code:    int32(um.Code),
		Message: um.Message,
		Result:  results,
	}, nil
}

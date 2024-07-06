package remittance

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/petnet/api/integration/perahub"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) RelationshipUpdate(ctx context.Context, req *bpa.RelationshipUpdateRequest) (res *bpa.RelationshipUpdateResponse, err error) {
	log := logging.FromContext(ctx)
	rvsm, err := s.ph.RelationshipUpdate(ctx, perahub.RelationshipUpdateReq{
		Relationship: req.GetRelationship(),
	}, req.ID)
	if err != nil {
		logging.WithError(err, log).Error("RelationshipUpdate error")
		return nil, handlePerahubError(err)
	}

	return &bpa.RelationshipUpdateResponse{
		Code:    int32(rvsm.Code),
		Message: rvsm.Message,
		Result: &bpa.RelationshipUpdateResult{
			ID:           int32(rvsm.Result.ID),
			Relationship: rvsm.Result.Relationship,
			CreatedAt:    timestamppb.New(rvsm.Result.CreatedAt),
			UpdatedAt:    timestamppb.New(rvsm.Result.UpdatedAt),
		},
	}, nil
}

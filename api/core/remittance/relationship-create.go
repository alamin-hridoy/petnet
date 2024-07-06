package remittance

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/petnet/api/integration/perahub"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) RelationshipCreate(ctx context.Context, req *bpa.RelationshipCreateRequest) (res *bpa.RelationshipCreateResponse, err error) {
	log := logging.FromContext(ctx)
	rrc, err := s.ph.RemittanceRelationshipCreate(ctx, perahub.RemittanceRelationshipCreateReq{
		Relationship: req.GetRelationship(),
	})
	if err != nil {
		logging.WithError(err, log).Error("RelationshipCreate error")
		return nil, handlePerahubError(err)
	}

	return &bpa.RelationshipCreateResponse{
		Code:    int32(rrc.Code),
		Message: rrc.Message,
		Result: &bpa.RelationshipCreateResult{
			Relationship: rrc.Result.Relationship,
			CreatedAt:    timestamppb.New(rrc.Result.CreatedAt),
			UpdatedAt:    timestamppb.New(rrc.Result.UpdatedAt),
			ID:           int32(rrc.Result.ID),
		},
	}, nil
}

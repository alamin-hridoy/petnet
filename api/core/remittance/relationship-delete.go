package remittance

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) RelationshipDelete(ctx context.Context, req *bpa.RelationshipDeleteRequest) (res *bpa.RelationshipDeleteResponse, err error) {
	log := logging.FromContext(ctx)
	um, err := s.ph.RelationshipDelete(ctx, req.GetID())
	if err != nil {
		logging.WithError(err, log).Error("RelationshipDelete error")
		return nil, handlePerahubError(err)
	}

	return &bpa.RelationshipDeleteResponse{
		Code:    int32(um.Code),
		Message: um.Message,
		Result: &bpa.RelationshipDeleteResult{
			ID:           int32(um.Result.ID),
			Relationship: um.Result.Relationship,
			CreatedAt:    timestamppb.New(um.Result.CreatedAt),
			UpdatedAt:    timestamppb.New(um.Result.UpdatedAt),
			DeletedAt:    timestamppb.New(um.Result.DeletedAt),
		},
	}, nil
}

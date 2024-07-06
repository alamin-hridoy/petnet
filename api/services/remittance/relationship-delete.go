package remittance

import (
	"context"

	"brank.as/petnet/api/util"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
)

func (s *Svc) RelationshipDelete(ctx context.Context, req *bpa.RelationshipDeleteRequest) (*bpa.RelationshipDeleteResponse, error) {
	res, err := s.remittanceStore.RelationshipDelete(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

package remittance

import (
	"context"

	"brank.as/petnet/api/util"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
)

func (s *Svc) RelationshipGet(ctx context.Context, req *bpa.RelationshipGetRequest) (*bpa.RelationshipGetResponse, error) {
	res, err := s.remittanceStore.RelationshipGet(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

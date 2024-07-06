package remittance

import (
	"context"

	"brank.as/petnet/api/util"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
)

func (s *Svc) OccupationDelete(ctx context.Context, req *bpa.OccupationDeleteRequest) (*bpa.OccupationDeleteResponse, error) {
	res, err := s.remittanceStore.OccupationDelete(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

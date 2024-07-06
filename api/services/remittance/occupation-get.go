package remittance

import (
	"context"

	"brank.as/petnet/api/util"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
)

func (s *Svc) OccupationGet(ctx context.Context, req *bpa.OccupationGetRequest) (*bpa.OccupationGetResponse, error) {
	res, err := s.remittanceStore.OccupationGet(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

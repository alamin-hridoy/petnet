package remittance

import (
	"context"

	"brank.as/petnet/api/util"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
)

func (s *Svc) SourceOfFundGet(ctx context.Context, req *bpa.SourceOfFundGetRequest) (*bpa.SourceOfFundGetResponse, error) {
	res, err := s.remittanceStore.SourceOfFundGet(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

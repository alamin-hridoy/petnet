package remittance

import (
	"context"

	"brank.as/petnet/api/util"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
)

func (s *Svc) SourceOfFundDelete(ctx context.Context, req *bpa.SourceOfFundDeleteRequest) (*bpa.SourceOfFundDeleteResponse, error) {
	res, err := s.remittanceStore.SourceOfFundDelete(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

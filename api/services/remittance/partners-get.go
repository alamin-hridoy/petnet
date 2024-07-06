package remittance

import (
	"context"

	"brank.as/petnet/api/util"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
)

func (s *Svc) PartnersGet(ctx context.Context, req *bpa.PartnersGetRequest) (*bpa.PartnersGetResponse, error) {
	res, err := s.remittanceStore.PartnersGet(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

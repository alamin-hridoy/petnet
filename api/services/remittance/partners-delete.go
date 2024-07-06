package remittance

import (
	"context"

	"brank.as/petnet/api/util"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
)

func (s *Svc) PartnersDelete(ctx context.Context, req *bpa.PartnersDeleteRequest) (*bpa.PartnersDeleteResponse, error) {
	res, err := s.remittanceStore.PartnersDelete(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

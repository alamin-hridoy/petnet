package remittance

import (
	"context"

	"brank.as/petnet/api/util"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
)

func (s *Svc) PurposeOfRemittanceGet(ctx context.Context, req *bpa.PurposeOfRemittanceGetRequest) (*bpa.PurposeOfRemittanceGetResponse, error) {
	res, err := s.remittanceStore.PurposeOfRemittanceGet(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

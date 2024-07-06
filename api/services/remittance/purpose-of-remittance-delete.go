package remittance

import (
	"context"

	"brank.as/petnet/api/util"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
)

func (s *Svc) PurposeOfRemittanceDelete(ctx context.Context, req *bpa.PurposeOfRemittanceDeleteRequest) (*bpa.PurposeOfRemittanceDeleteResponse, error) {
	res, err := s.remittanceStore.PurposeOfRemittanceDelete(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

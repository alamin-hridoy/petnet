package remittance

import (
	"context"

	"brank.as/petnet/api/util"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
)

func (s *Svc) RemittanceEmploymentDelete(ctx context.Context, req *bpa.RemittanceEmploymentDeleteRequest) (*bpa.RemittanceEmploymentDeleteResponse, error) {
	res, err := s.remittanceStore.RemittanceEmploymentDelete(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

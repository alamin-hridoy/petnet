package remittance

import (
	"context"

	"brank.as/petnet/api/util"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
)

func (s *Svc) EmploymentGet(ctx context.Context, req *bpa.EmploymentGetRequest) (*bpa.EmploymentGetResponse, error) {
	res, err := s.remittanceStore.EmploymentGet(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

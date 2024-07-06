package remittance

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"brank.as/petnet/api/util"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
)

func (s *Svc) EmploymentGrid(ctx context.Context, em *emptypb.Empty) (*bpa.EmploymentGridResponse, error) {
	res, err := s.remittanceStore.EmploymentGrid(ctx)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

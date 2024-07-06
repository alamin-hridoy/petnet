package remittance

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"brank.as/petnet/api/util"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
)

func (s *Svc) PurposeOfRemittanceGrid(ctx context.Context, em *emptypb.Empty) (*bpa.PurposeOfRemittanceGridResponse, error) {
	res, err := s.remittanceStore.PurposeOfRemittanceGrid(ctx)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

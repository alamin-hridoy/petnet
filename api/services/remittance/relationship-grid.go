package remittance

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"brank.as/petnet/api/util"
	bpa "brank.as/petnet/gunk/drp/v1/remittance"
)

func (s *Svc) RelationshipGrid(ctx context.Context, em *emptypb.Empty) (*bpa.RelationshipGridResponse, error) {
	res, err := s.remittanceStore.RelationshipGrid(ctx, em)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}
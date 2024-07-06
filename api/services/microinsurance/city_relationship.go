package microinsurance

import (
	"context"

	migunk "brank.as/petnet/gunk/drp/v1/microinsurance"
	"google.golang.org/protobuf/types/known/emptypb"

	"brank.as/petnet/api/util"
)

// GetRelationships ...
func (s *Svc) GetRelationships(ctx context.Context, r *emptypb.Empty) (*migunk.GetRelationshipsResult, error) {
	res, err := s.store.GetRelationships(ctx)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

// GetAllCities ...
func (s *Svc) GetAllCities(ctx context.Context, r *emptypb.Empty) (*migunk.CityListResult, error) {
	res, err := s.store.GetAllCities(ctx)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

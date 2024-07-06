package microinsurance

import (
	"context"

	"google.golang.org/grpc/codes"

	coreerror "brank.as/petnet/api/core/error"
	migunk "brank.as/petnet/gunk/drp/v1/microinsurance"
)

// GetAllCities ...
func (s *MICoreSvc) GetAllCities(ctx context.Context) (*migunk.CityListResult, error) {
	res, err := s.cl.GetAllCities(ctx)
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	if res == nil || res.AllCities == nil {
		return nil, coreerror.NewCoreError(codes.NotFound, "no cities found")
	}

	cities := make([]*migunk.City, 0, len(res.AllCities))
	for _, c := range res.AllCities {
		cities = append(cities, &migunk.City{
			CityCode:     c.CityCode,
			CityName:     c.CityName,
			ProvinceCode: c.ProvinceCode,
			ProvinceName: c.ProvinceName,
		})
	}

	return &migunk.CityListResult{
		Cities: cities,
	}, nil
}

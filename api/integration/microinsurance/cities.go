package microinsurance

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/codes"

	coreerror "brank.as/petnet/api/core/error"
)

// City ...
type City struct {
	CityCode     string `json:"cityCode"`
	CityName     string `json:"cityName"`
	ProvinceCode string `json:"provinceCode"`
	ProvinceName string `json:"provinceName"`
}

// CityListResult ...
type CityListResult struct {
	SessionID  string `json:"sessionID"`
	StatusCode string `json:"statusCode"`
	StatusDesc string `json:"statusDesc"`
	AllCities  []City `json:"allcities"`
}

// GetAllCitiesResult ...
type GetAllCitiesResult struct {
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Result  *CityListResult `json:"result"`
}

func (c *Client) GetAllCities(ctx context.Context) (*CityListResult, error) {
	rawRes, err := c.phService.GetMicroInsurance(ctx, c.getUrl("get-all-cities"))
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	var res GetAllCitiesResult
	err = json.Unmarshal(rawRes, &res)
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, "Invalid perahub response")
	}

	return res.Result, nil
}

package perahub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetProvincesCityListRequest struct {
	ID          int    `json:"id"`
	PartnerCode string `json:"partner_code"`
}

type GetProvincesCityListResults struct {
	Province string   `json:"province"`
	CityList []string `json:"cityList"`
}

type GetProvincesCityListResponse struct {
	Code    int                           `json:"code"`
	Message string                        `json:"message"`
	Result  []GetProvincesCityListResults `json:"result"`
}

type GetBrgyListRequest struct {
	City string `json:"city"`
}

type GetBrgyListResults struct {
	Barangay string `json:"barangay"`
	Zipcode  string `json:"zipcode"`
}

type GetBrgyListResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Result  []GetBrgyListResults `json:"result"`
}

type GetUtilityPurposeResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Result  []string `json:"result"`
}

type GetUtilityRelationshipResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Result  []string `json:"result"`
}

type GetUtilityPartnerResult struct {
	PartnerCode string `json:"partner_code"`
	PartnerName string `json:"partner_name"`
	Status      int    `json:"status"`
}

type GetUtilityPartnerResponse struct {
	Code    int                       `json:"code"`
	Message string                    `json:"message"`
	Result  []GetUtilityPartnerResult `json:"result"`
}

type GetUtilityOccupationResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Result  []string `json:"result"`
}

type GetUtilityEmploymentResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Result  []string `json:"result"`
}

type GetUtilitySourceFundResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Result  []string `json:"result"`
}

func (s *Svc) GetProvincesCityList(ctx context.Context, req GetProvincesCityListRequest) (*GetProvincesCityListResponse, error) {
	if req.ID == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.PartnerCode == "" {
		return nil, status.Error(codes.InvalidArgument, "partner_code is required")
	}

	baseUrl := fmt.Sprintf("perahub-remit/address/provinces?id=%d&partner_code=%s", req.ID, req.PartnerCode)
	nonexUrl := s.nonexURL(baseUrl)
	decodedUrl, _ := url.QueryUnescape(nonexUrl)
	res, err := s.getNonex(ctx, decodedUrl)
	if err != nil {
		return nil, err
	}

	rb := &GetProvincesCityListResponse{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Svc) GetBrgyList(ctx context.Context, req GetBrgyListRequest) (*GetBrgyListResponse, error) {
	if req.City == "" {
		return nil, status.Error(codes.InvalidArgument, "city is required")
	}

	baseUrl := fmt.Sprintf("perahub-remit/address/brgy/%s", req.City)
	nonexUrl := s.nonexURL(baseUrl)
	decodedUrl, _ := url.QueryUnescape(nonexUrl)
	res, err := s.getNonex(ctx, decodedUrl)
	if err != nil {
		return nil, err
	}

	rb := &GetBrgyListResponse{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Svc) GetUtilityPurpose(ctx context.Context) (*GetUtilityPurposeResponse, error) {
	res, err := s.getNonex(ctx, s.nonexURL("perahub-remit/utility/purpose"))
	if err != nil {
		return nil, err
	}

	rb := &GetUtilityPurposeResponse{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Svc) GetUtilityRelationship(ctx context.Context) (*GetUtilityRelationshipResponse, error) {
	res, err := s.getNonex(ctx, s.nonexURL("perahub-remit/utility/relationship"))
	if err != nil {
		return nil, err
	}

	rb := &GetUtilityRelationshipResponse{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Svc) GetUtilityPartner(ctx context.Context) (*GetUtilityPartnerResponse, error) {
	res, err := s.getNonex(ctx, s.nonexURL("perahub-remit/utility/partner"))
	if err != nil {
		return nil, err
	}

	rb := &GetUtilityPartnerResponse{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Svc) GetUtilityOccupation(ctx context.Context) (*GetUtilityOccupationResponse, error) {
	res, err := s.getNonex(ctx, s.nonexURL("perahub-remit/utility/occupation"))
	if err != nil {
		return nil, err
	}

	rb := &GetUtilityOccupationResponse{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Svc) GetUtilityEmployment(ctx context.Context) (*GetUtilityEmploymentResponse, error) {
	res, err := s.getNonex(ctx, s.nonexURL("perahub-remit/utility/employment"))
	if err != nil {
		return nil, err
	}

	rb := &GetUtilityEmploymentResponse{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Svc) GetUtilitySourceFund(ctx context.Context) (*GetUtilitySourceFundResponse, error) {
	res, err := s.getNonex(ctx, s.nonexURL("perahub-remit/utility/sourcefund"))
	if err != nil {
		return nil, err
	}

	rb := &GetUtilitySourceFundResponse{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

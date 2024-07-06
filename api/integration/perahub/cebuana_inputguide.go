package perahub

import (
	"context"
	"encoding/json"
	"net/url"
)

type CEBCountryResponse struct {
	CountryID         json.Number `json:"CountryID"`
	CountryName       string      `json:"CountryName"`
	CountryCodeAplha2 string      `json:"CountryCodeAplha2"`
	CountryCodeAlpha3 string      `json:"CountryCodeAlpha3"`
	PhoneCode         string      `json:"PhoneCode"`
}

type CEBCountryListResponse struct {
	Country []CEBCountryResponse `json:"Country"`
}

type CEBCountryBaseResponse struct {
	Code    json.Number            `json:"code"`
	Message string                 `json:"message"`
	Result  CEBCountryListResponse `json:"result"`
	RemcoID json.Number            `json:"remco_id"`
}

type CEBCurrencyResponse struct {
	CurrencyID  json.Number `json:"CurrencyID"`
	Code        string      `json:"Code"`
	Description string      `json:"Description"`
}

type CEBCurrencyListResponse struct {
	Currency CEBCurrencyResponse `json:"Currency"`
}

type CEBCurrencyReq struct {
	AgentCode string `json:"agent_code"`
}

type CEBCurrencyBaseResponse struct {
	Code    json.Number             `json:"code"`
	Message string                  `json:"message"`
	Result  CEBCurrencyListResponse `json:"result"`
	RemcoID json.Number             `json:"remco_id"`
}

type CEBSourceFundsResponse struct {
	SourceOfFundID json.Number `json:"SourceOfFundID"`
	SourceOfFund   string      `json:"SourceOfFund"`
}

type CEBSourceFundsListResponse struct {
	ClientSourceOfFund []CEBSourceFundsResponse `json:"ClientSourceOfFund"`
}

type CEBSourceFundsBaseResponse struct {
	Code    json.Number                `json:"code"`
	Message string                     `json:"message"`
	Result  CEBSourceFundsListResponse `json:"result"`
	RemcoID json.Number                `json:"remco_id"`
}

type CEBIdentificationTypeResponse struct {
	IdentificationTypeID json.Number `json:"IdentificationTypeID"`
	Description          string      `json:"Description"`
	SmsCode              string      `json:"SmsCode"`
}

type (
	CEBIdTypesResponseList []CEBIdentificationTypeResponse
)

func (s *Svc) CEBgetCountries(ctx context.Context) (*CEBCountryBaseResponse, error) {
	res, err := s.getNonex(ctx, s.nonexURL("cebuana/get-country"))
	if err != nil {
		return nil, err
	}

	rb := &CEBCountryBaseResponse{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Svc) CEBgetCurrencies(ctx context.Context, agentCode *CEBCurrencyReq) (*CEBCurrencyBaseResponse, error) {
	nonexUrl := s.nonexURL("cebuana/send-currency-collection?agent_code=" + agentCode.AgentCode)
	decodedUrl, _ := url.QueryUnescape(nonexUrl)
	res, err := s.getNonex(ctx, decodedUrl)
	if err != nil {
		return nil, err
	}

	rb := &CEBCurrencyBaseResponse{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Svc) CEBgetSourceFunds(ctx context.Context) (*CEBSourceFundsBaseResponse, error) {
	res, err := s.getNonex(ctx, s.nonexURL("cebuana/get-source-of-fund"))
	if err != nil {
		return nil, err
	}

	rb := &CEBSourceFundsBaseResponse{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Svc) CEBgetIdTypes(ctx context.Context) (*CEBIdTypesResponseList, error) {
	res, err := s.getNonex(ctx, s.nonexURL("cebuana/identification-types"))
	if err != nil {
		return nil, err
	}

	rb := &CEBIdTypesResponseList{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

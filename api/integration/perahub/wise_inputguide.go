package perahub

import (
	"context"
	"encoding/json"
)

type WISEBaseResponse struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type WISECurrencyResponse struct {
	Currency    string `json:"currency"`
	Description string `json:"description"`
}

type (
	WISEBaseResponseList     []WISEBaseResponse
	WISECurrencyResponseList []WISECurrencyResponse
)

func (s *Svc) WISEgetCountries(ctx context.Context) (*WISEBaseResponseList, error) {
	res, err := s.getNonex(ctx, s.nonexURL("transferwise/utils/country"))
	if err != nil {
		return nil, err
	}

	rb := &WISEBaseResponseList{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Svc) WISEgetStates(ctx context.Context, countryCode string) (*WISEBaseResponseList, error) {
	res, err := s.getNonex(ctx, s.nonexURL("transferwise/utils/state/"+countryCode))
	if err != nil {
		return nil, err
	}

	rb := &WISEBaseResponseList{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Svc) WISEgetCurrencies(ctx context.Context) (*WISECurrencyResponseList, error) {
	res, err := s.getNonex(ctx, s.nonexURL("transferwise/utils/currency"))
	if err != nil {
		return nil, err
	}

	rb := &WISECurrencyResponseList{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

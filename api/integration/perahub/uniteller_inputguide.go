package perahub

import (
	"context"
	"encoding/json"
)

type UNTBaseResponse struct {
	Code string `json:"CODE"`
	Name string `json:"NAME"`
}

type UNTCommonResponse struct {
	Country     string `json:"COUNTRY"`
	Code        string `json:"CODE"`
	Description string `json:"DESCRIPCION"`
}

type UNTStatesResponse struct {
	Country   string `json:"COUNTRY_NAME"`
	StateName string `json:"STATE_NAME"`
	UtlCode   string `json:"UTL_CODE"`
}

type UNTUsStatesResponse struct {
	Country   string `json:"COUNTRY_NAME"`
	StateName string `json:"STATE_NAME"`
	UtlCode   string `json:"UTL_CODE"`
}

type UNTBaseResponseList struct {
	Data []UNTBaseResponse `json:"Data"`
}

type UNTCommonResponseList struct {
	Data []UNTCommonResponse `json:"Data"`
}

type UNTStateResponseList struct {
	Data []UNTStatesResponse `json:"Data"`
}

type UNTUsStateResponseList struct {
	Data []UNTUsStatesResponse `json:"Data"`
}

func (s *Svc) UNTgetCountries(ctx context.Context) (*UNTBaseResponseList, error) {
	res, err := s.getNonex(ctx, s.nonexURL("uniteller/countries"))
	if err != nil {
		return nil, err
	}

	rb := &UNTBaseResponseList{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Svc) UNTgetCurrencies(ctx context.Context) (*UNTBaseResponseList, error) {
	res, err := s.getNonex(ctx, s.nonexURL("uniteller/currencies"))
	if err != nil {
		return nil, err
	}

	rb := &UNTBaseResponseList{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Svc) UNTgetOccupations(ctx context.Context) (*UNTCommonResponseList, error) {
	res, err := s.getNonex(ctx, s.nonexURL("uniteller/occupations"))
	if err != nil {
		return nil, err
	}

	rb := &UNTCommonResponseList{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Svc) UNTgetIds(ctx context.Context) (*UNTCommonResponseList, error) {
	res, err := s.getNonex(ctx, s.nonexURL("uniteller/ids"))
	if err != nil {
		return nil, err
	}

	rb := &UNTCommonResponseList{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Svc) UNTgetStates(ctx context.Context) (*UNTStateResponseList, error) {
	res, err := s.getNonex(ctx, s.nonexURL("uniteller/ph-states"))
	if err != nil {
		return nil, err
	}

	rb := &UNTStateResponseList{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Svc) UNTgetUsStates(ctx context.Context) (*UNTUsStateResponseList, error) {
	res, err := s.getNonex(ctx, s.nonexURL("uniteller/usa-states"))
	if err != nil {
		return nil, err
	}

	rb := &UNTUsStateResponseList{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

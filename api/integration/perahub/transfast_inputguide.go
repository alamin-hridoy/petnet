package perahub

import (
	"context"
	"encoding/json"
)

type TFBaseResponseBody struct {
	Code    json.Number `json:"code"`
	Msg     string      `json:"message"`
	RemcoID json.Number `json:"remco_id"`
}

type TFIDsRespBody struct {
	TFBaseResponseBody
	Result TFIDsResult `json:"result"`
}

type TFIDsResult struct {
	IDs []TFID `json:"ReceiverTypeOfIds"`
}

type TFID struct {
	ID                     json.Number `json:"Id"`
	Name                   string      `json:"Name"`
	RequiredExpirationDate string      `json:"RequiredExpirationDate"`
	RequiredIssueDate      string      `json:"RequiredIssueDate"`
	CountryIsoCode         string      `json:"CountryIsoCode"`
}

type TFRelationsRespBody struct {
	TFBaseResponseBody
	Result TFRelationsResult `json:"result"`
}

type TFRelationsResult struct {
	Relations []TFIDName `json:"RelationshipWithSenders"`
}

type TFIDName struct {
	ID   json.Number `json:"Id"`
	Name string      `json:"Name"`
}

type TFOccupsRespBody struct {
	TFBaseResponseBody
	Result TFOccupsResult `json:"result"`
}

type TFOccupsResult struct {
	Occups []TFIDName `json:"Occupations"`
}

type TFPrpsRespBody struct {
	TFBaseResponseBody
	Result TFPrpsResult `json:"result"`
}

type TFPrpsResult struct {
	Prps []TFPrp `json:"RemittancePurposes"`
}

type TFPrp struct {
	ID             json.Number `json:"Id"`
	Name           string      `json:"Name"`
	CountryIsoCode string      `json:"CountryIsoCode"`
}

func (s *Svc) TFIDs(ctx context.Context) (*TFIDsRespBody, error) {
	res, err := s.getNonex(ctx, s.nonexURL("transfast/receiver-ids"))
	if err != nil {
		return nil, err
	}

	rb := &TFIDsRespBody{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Svc) TFRelations(ctx context.Context) (*TFRelationsRespBody, error) {
	res, err := s.getNonex(ctx, s.nonexURL("transfast/relationships"))
	if err != nil {
		return nil, err
	}

	rb := &TFRelationsRespBody{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Svc) TFOccupations(ctx context.Context) (*TFOccupsRespBody, error) {
	res, err := s.getNonex(ctx, s.nonexURL("transfast/beneficiary-occupations"))
	if err != nil {
		return nil, err
	}

	rb := &TFOccupsRespBody{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

func (s *Svc) TFPrps(ctx context.Context) (*TFPrpsRespBody, error) {
	res, err := s.getNonex(ctx, s.nonexURL("transfast/remittance-purposes"))
	if err != nil {
		return nil, err
	}

	rb := &TFPrpsRespBody{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

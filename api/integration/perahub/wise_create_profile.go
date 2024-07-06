package perahub

import (
	"context"
	"encoding/json"
)

type WISECreateProfileReq struct {
	Email      string              `json:"email"`
	Type       string              `json:"type"`
	Details    WISECreatePFDetails `json:"details"`
	Address    WISEPFAddress       `json:"address"`
	Occupation string              `json:"occupation"`
}

type WISECreateProfileResp struct {
	Msg   string `json:"message"`
	Error string `json:"error"`
}

type WISECreatePFDetails struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	BirthDate   string `json:"dateOfBirth"`
	PhoneNumber string `json:"phoneNumber"`
}

type WISEPFAddress struct {
	Country   string `json:"country"`
	FirstLine string `json:"firstLine"`
	PostCode  string `json:"postCode"`
	City      string `json:"city"`
	State     string `json:"state"`
}

func (s *Svc) WISECreateProfile(ctx context.Context, req WISECreateProfileReq) (*WISECreateProfileResp, error) {
	res, err := s.postNonex(ctx, s.nonexURL("transferwise/profiles"), req)
	if err != nil {
		return nil, err
	}

	rb := &WISECreateProfileResp{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

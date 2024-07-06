package perahub

import (
	"context"
	"encoding/json"
)

type WISEGetProfileReq struct {
	Email string `json:"email"`
}

type WISEGetProfileResp struct {
	Type      string           `json:"type"`
	Details   WISEGetPFDetails `json:"details"`
	Address   WISEPFAddress    `json:"address"`
	ProfileID json.Number      `json:"profileID"`
	Error     string           `json:"error"`
}

type WISEGetPFDetails struct {
	FirstName       string           `json:"firstName"`
	LastName        string           `json:"lastName"`
	BirthDate       string           `json:"dateOfBirth"`
	PhoneNumber     string           `json:"phoneNumber"`
	Avatar          string           `json:"avatar"`
	Occupation      string           `json:"occupation"`
	Occupations     []WISEOccupation `json:"occupations"`
	PrimaryAddress  json.Number      `json:"primaryAddress"`
	FirstNameInKana string           `json:"firstNameInKana"`
	LastNameInKana  string           `json:"lastNameInKana"`
}

type WISEOccupation struct {
	Code   string `json:"code"`
	Format string `json:"format"`
}

func (s *Svc) WISEGetProfile(ctx context.Context, req WISEGetProfileReq) (*WISEGetProfileResp, error) {
	res, err := s.postNonex(ctx, s.nonexURL("transferwise/profiles/personal"), req)
	if err != nil {
		return nil, err
	}

	rb := &WISEGetProfileResp{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

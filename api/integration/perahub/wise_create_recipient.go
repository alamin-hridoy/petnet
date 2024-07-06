package perahub

import (
	"context"
	"encoding/json"
)

type WISECreateRecipientReq struct {
	Email             string            `json:"email"`
	Currency          string            `json:"currency"`
	Type              string            `json:"type"`
	OwnedByCustomer   bool              `json:"ownedByCustomer"`
	AccountHolderName string            `json:"accountHolderName"`
	Requirements      []WISERequirement `json:"requirements"`
}

type WISECreateRecipientResp struct {
	RecipientID       string        `json:"recipientID"`
	Details           WISECRDetails `json:"details"`
	AccountHolderName string        `json:"accountHolderName"`
	Currency          string        `json:"currency"`
	OwnedByCustomer   bool          `json:"ownedByCustomer"`
	Country           string        `json:"country"`
	Msg               string        `json:"message"`
}

type WISECRDetails struct {
	Address       WISECRAddress `json:"address"`
	LegalType     string        `json:"legalType"`
	AccountNumber string        `json:"accountNumber"`
	SortCode      string        `json:"sortCode"`
}

type WISERequirement struct {
	PropName string      `json:"propName"`
	Value    interface{} `json:"value"`
}

type WISECRAddress struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	FirstLine   string `json:"firstLine"`
	PostCode    string `json:"postCode"`
	City        string `json:"city"`
	State       string `json:"state"`
}

func (s *Svc) WISECreateRecipient(ctx context.Context, req WISECreateRecipientReq) (*WISECreateRecipientResp, error) {
	res, err := s.postNonex(ctx, s.nonexURL("transferwise/recipients"), req)
	if err != nil {
		return nil, err
	}

	rb := &WISECreateRecipientResp{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

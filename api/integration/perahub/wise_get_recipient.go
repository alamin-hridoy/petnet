package perahub

import (
	"context"
	"encoding/json"
)

type WISEGetRecipientsReq struct {
	Email    string `json:"email"`
	Currency string `json:"currency"`
}

type WISEGetRecipientsResp struct {
	Recipients []WISERecipient `json:"recipients"`
}

type WISERecipient struct {
	RecipientID        json.Number        `json:"recipientID"`
	Details            WISEGRDetails      `json:"details"`
	AccountSummary     string             `json:"accountSummary"`
	LongAccountSummary string             `json:"longAccountSummary"`
	DisplayFields      []WISEDisplayField `json:"displayFields"`
	FullName           string             `json:"fullName"`
	Currency           string             `json:"currency"`
	Country            string             `json:"country"`
}

type WISEGRDetails struct {
	AccountNumber    string `json:"accountNumber"`
	SortCode         string `json:"sortCode"`
	HashedByLooseAlg string `json:"hashedByLooseHashAlgorithm"`
}

type WISEDisplayField struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func (s *Svc) WISEGetRecipients(ctx context.Context, req WISEGetRecipientsReq) (*WISEGetRecipientsResp, error) {
	res, err := s.postNonex(ctx, s.nonexURL("transferwise/recipients/list"), req)
	if err != nil {
		return nil, err
	}

	rcp := &[]WISERecipient{}
	if err := json.Unmarshal(res, rcp); err != nil {
		return nil, err
	}
	rb := &WISEGetRecipientsResp{Recipients: *rcp}
	return rb, nil
}

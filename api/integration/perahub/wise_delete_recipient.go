package perahub

import (
	"context"
	"encoding/json"
)

type WISEDeleteRecipientReq struct {
	RecipientID string
	Email       string `json:"email"`
}

type WISEDeleteRecipientResp struct {
	Msg string `json:"message"`
}

func (s *Svc) WISEDeleteRecipient(ctx context.Context, req WISEDeleteRecipientReq) (*WISEDeleteRecipientResp, error) {
	res, err := s.postNonex(ctx, s.nonexURL("transferwise/recipients/"+req.RecipientID), req)
	if err != nil {
		return nil, err
	}

	rb := &WISEDeleteRecipientResp{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

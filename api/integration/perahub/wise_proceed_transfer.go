package perahub

import (
	"context"
	"encoding/json"
)

type WISEProceedTransferReq struct {
	Email       string        `json:"email"`
	RecipientID json.Number   `json:"recipientID"`
	Details     WISEPCDetails `json:"details"`
}

type WISEProceedTransferResp struct {
	TransferID     string        `json:"transferID"`
	Details        WISEPCDetails `json:"details"`
	CustomerTxnID  string        `json:"customerTransactionId"`
	RecipientID    json.Number   `json:"recipientID"`
	Status         string        `json:"status"`
	SourceCurrency string        `json:"sourceCurrency"`
	TargetCurrency string        `json:"targetCurrency"`
	SourceAmount   json.Number   `json:"sourceAmount"`
	DateCreated    string        `json:"dateCreated"`
}

type WISEPCDetails struct {
	Reference string `json:"reference"`
}

func (s *Svc) WISEProceedTransfer(ctx context.Context, req WISEProceedTransferReq) (*WISEProceedTransferResp, error) {
	res, err := s.postNonex(ctx, s.nonexURL("transferwise/transfers/proceed"), req)
	if err != nil {
		return nil, err
	}

	rb := &WISEProceedTransferResp{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

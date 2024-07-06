package perahub

import (
	"context"
	"encoding/json"
)

type WISEPrepareTransferReq struct {
	Email               string `json:"email"`
	RecipientID         string `json:"recipientID"`
	AccountHolderName   string `json:"accountHolderName"`
	SourceAccountNumber string `json:"sourceAccountNumber"`
}

type WISEPrepareTransferResp struct {
	Requirements        []WISERequirementsResp `json:"requirements"`
	UpdatedQuoteSummary WISEQuoteInquiryResp   `json:"updatedQuoteSummary"`
}

func (s *Svc) WISEPrepareTransfer(ctx context.Context, req WISEPrepareTransferReq) (*WISEPrepareTransferResp, error) {
	res, err := s.postNonex(ctx, s.nonexURL("transferwise/transfers/prepare"), req)
	if err != nil {
		return nil, err
	}

	rb := &WISEPrepareTransferResp{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

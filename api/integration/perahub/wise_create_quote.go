package perahub

import (
	"context"
	"encoding/json"
)

type WISECreateQuoteReq struct {
	Email          string      `json:"email"`
	SourceCurrency string      `json:"sourceCurrency"`
	TargetCurrency string      `json:"targetCurrency"`
	SourceAmount   json.Number `json:"sourceAmount"`
}

type WISECreateQuoteResp struct {
	Requirements []WISERequirementsResp `json:"requirements"`
	QuoteSummary WISEQuoteInquiryResp   `json:"quoteSummary"`
}

func (s *Svc) WISECreateQuote(ctx context.Context, req WISECreateQuoteReq) (*WISECreateQuoteResp, error) {
	res, err := s.postNonex(ctx, s.nonexURL("transferwise/quotes"), req)
	if err != nil {
		return nil, err
	}

	rb := &WISECreateQuoteResp{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

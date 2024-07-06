package perahub

import (
	"context"
	"encoding/json"
)

type WISEGetQuoteRequirementsReq struct {
	SourceCurrency string      `json:"sourceCurrency"`
	TargetCurrency string      `json:"targetCurrency"`
	SourceAmount   json.Number `json:"sourceAmount"`
}

type WISEGetQuoteRequirementsResp struct {
	Requirements []WISERequirementsResp `json:"requirements"`
	Quote        WISEQuote              `json:"quote"`
}

type WISEQuote struct {
	SourceCurrency string      `json:"sourceCurrency"`
	TargetCurrency string      `json:"targetCurrency"`
	SourceAmount   json.Number `json:"sourceAmount"`
}

func (s *Svc) WISEGetQuoteRequirements(ctx context.Context, req WISEGetQuoteRequirementsReq) (*WISEGetQuoteRequirementsResp, error) {
	res, err := s.postNonex(ctx, s.nonexURL("transferwise/quotes/requirements"), req)
	if err != nil {
		return nil, err
	}

	rb := &WISEGetQuoteRequirementsResp{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

package perahub

import (
	"context"
	"encoding/json"
)

type WISEQuoteInquiryReq struct {
	SourceCurrency string      `json:"sourceCurrency"`
	TargetCurrency string      `json:"targetCurrency"`
	SourceAmount   json.Number `json:"sourceAmount"`
}

type WISEQuoteInquiryResp struct {
	SourceCurrency string           `json:"sourceCurrency"`
	TargetCurrency string           `json:"targetCurrency"`
	SourceAmount   json.Number      `json:"sourceAmount"`
	TargetAmount   json.Number      `json:"targetAmount"`
	FeeBreakdown   WISEFeeBreakdown `json:"feeBreakdown"`
	TotalFee       json.Number      `json:"totalFee"`
	TransferAmount json.Number      `json:"transferAmount"`
	PayOut         string           `json:"payOut"`
	Rate           json.Number      `json:"rate"`
}

type WISEFeeBreakdown struct {
	Transferwise json.Number `json:"transferwise"`
	PayIn        json.Number `json:"payIn"`
	Discount     json.Number `json:"discount"`
	Total        json.Number `json:"total"`
	PriceSetID   json.Number `json:"priceSetId"`
	Partner      json.Number `json:"partner"`
}

func (s *Svc) WISEQuoteInquiry(ctx context.Context, req WISEQuoteInquiryReq) (*WISEQuoteInquiryResp, error) {
	res, err := s.postNonex(ctx, s.nonexURL("transferwise/quotes/inquiry"), req)
	if err != nil {
		return nil, err
	}

	rb := &WISEQuoteInquiryResp{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

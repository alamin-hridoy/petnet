package perahub

import (
	"context"
	"encoding/json"
	"net/url"
)

const (
	CebuanaSFInquire = "Successful"
)

type CebuanaSFInquiryRequest struct {
	PrincipalAmount json.Number `json:"principal_amount"`
	CurrencyID      json.Number `json:"currency_id"`
	AgentCode       string      `json:"agent_code"`
}

type CebuanaSFInquiryRespBody struct {
	Code    json.Number            `json:"code"`
	Message string                 `json:"message"`
	Result  CebuanaSFInquiryResult `json:"result"`
	RemcoID json.Number            `json:"remco_id"`
}

type CebuanaSFInquiryResult struct {
	ResultStatus string      `json:"ResultStatus"`
	MessageID    json.Number `json:"MessageID"`
	LogID        json.Number `json:"LogID"`
	ServiceFee   string      `json:"ServiceFee"`
}

func (s *Svc) CebuanaSFInquiry(ctx context.Context, sr CebuanaSFInquiryRequest) (*CebuanaSFInquiryRespBody, error) {
	nonexUrl := s.nonexURL("cebuana/get-service-fee?principal_amount=" + string(sr.PrincipalAmount) + "&currency_id=" + string(sr.CurrencyID) + "&agent_code=" + sr.AgentCode)
	decodedUrl, _ := url.QueryUnescape(nonexUrl)
	res, err := s.getNonex(ctx, decodedUrl)
	if err != nil {
		return nil, err
	}

	rb := &CebuanaSFInquiryRespBody{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

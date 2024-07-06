package perahub

import (
	"context"
	"encoding/json"
)

type CicoInquireRequest struct {
	PartnerCode        string `json:"partner_code"`
	Provider           string `json:"provider"`
	TrxType            string `json:"trx_type"`
	ReferenceNumber    string `json:"reference_number"`
	PetnetTrackingno   string `json:"petnet_trackingno"`
	ProviderTrackingno string `json:"provider_trackingno"`
	Message            string `json:"message"`
}

type CicoInquireResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Result  *CicoInquireResult `json:"result"`
}

type CicoInquireResult struct {
	StatusMessage      string `json:"status_message"`
	PetnetTrackingno   string `json:"petnet_trackingno"`
	TrxType            string `json:"trx_type"`
	ReferenceNumber    string `json:"reference_number"`
	Amount             string `json:"amount"`
	ProviderTrackingno string `json:"provider_trackingno"`
	Expiry             string `json:"expiry"`
	CustomerName       string `json:"customer_name"`
	CustomerFirstname  string `json:"customer_firstname"`
	CustomerLastname   string `json:"customer_lastname"`
	MerchantID         string `json:"MerchantId"`
	PartnerCode        string `json:"partner_code"`
	AccountNumber      string `json:"account_number"`
	ServiceCharge      int    `json:"service_charge"`
	CreatedAt          string `json:"created_at"`
}

func (s *Svc) CicoInquire(ctx context.Context, sr CicoInquireRequest) (*CicoInquireResponse, error) {
	res, err := s.cicoPost(ctx, s.cicoURL("inquiry"), sr)
	if err != nil {
		return nil, err
	}

	CicoRes := &CicoInquireResponse{}
	if err := json.Unmarshal(res, CicoRes); err != nil {
		return nil, err
	}
	return CicoRes, nil
}

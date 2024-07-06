package perahub

import (
	"context"
	"encoding/json"
)

type CicoExecuteRequest struct {
	PartnerCode      string `json:"partner_code"`
	PetnetTrackingno string `json:"petnet_trackingno"`
	TrxDate          string `json:"trx_date"`
}

type CicoExecuteResult struct {
	PartnerCode        string `json:"partner_code"`
	Provider           string `json:"provider"`
	PetnetTrackingno   string `json:"petnet_trackingno"`
	TrxDate            string `json:"trx_date"`
	TrxType            string `json:"trx_type"`
	ProviderTrackingno string `json:"provider_trackingno"`
	ReferenceNumber    string `json:"reference_number"`
	PrincipalAmount    int    `json:"principal_amount"`
	Charges            int    `json:"charges"`
	TotalAmount        int    `json:"total_amount"`
}

type CicoExecuteResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Result  *CicoExecuteResult `json:"result"`
}

func (s *Svc) CicoExecute(ctx context.Context, sr CicoExecuteRequest) (*CicoExecuteResponse, error) {
	res, err := s.cicoPost(ctx, s.cicoURL("execute"), sr)
	if err != nil {
		return nil, err
	}

	CicoRes := &CicoExecuteResponse{}
	if err := json.Unmarshal(res, CicoRes); err != nil {
		return nil, err
	}
	return CicoRes, nil
}

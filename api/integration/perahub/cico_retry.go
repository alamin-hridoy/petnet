package perahub

import (
	"context"
	"encoding/json"
)

type CicoRetryRequest struct {
	PartnerCode      string `json:"partner_code"`
	PetnetTrackingno string `json:"petnet_trackingno"`
	TrxDate          string `json:"trx_date"`
}
type CicoRetryResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Result  *CicoRetryResult `json:"result"`
}

type CicoRetryResult struct {
	PartnerCode        string     `json:"partner_code"`
	Provider           string     `json:"provider"`
	PetnetTrackingno   string     `json:"petnet_trackingno"`
	TrxDate            string     `json:"trx_date"`
	TrxType            string     `json:"trx_type"`
	ProviderTrackingno string     `json:"provider_trackingno"`
	ReferenceNumber    string     `json:"reference_number"`
	PrincipalAmount    int        `json:"principal_amount"`
	Charges            int        `json:"charges"`
	TotalAmount        int        `json:"total_amount"`
	OTPPayload         OTPPayload `json:"otp_payload"`
}

type OTPPayload struct {
	CommandID int    `json:"commandId"`
	Payload   string `json:"payload"`
}

func (s *Svc) CicoRetry(ctx context.Context, sr CicoRetryRequest) (*CicoRetryResponse, error) {
	res, err := s.cicoPost(ctx, s.cicoURL("retry"), sr)
	if err != nil {
		return nil, err
	}

	CicoRes := &CicoRetryResponse{}
	if err := json.Unmarshal(res, CicoRes); err != nil {
		return nil, err
	}
	return CicoRes, nil
}

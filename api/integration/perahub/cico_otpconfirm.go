package perahub

import (
	"context"
	"encoding/json"
)

type CicoOTPConfirmRequest struct {
	PartnerCode      string     `json:"partner_code"`
	PetnetTrackingno string     `json:"petnet_trackingno"`
	TrxDate          string     `json:"trx_date"`
	OTP              string     `json:"otp"`
	OTPPayload       OTPpayload `json:"otp_payload"`
}

type CicoOTPConfirmResponse struct {
	Code    int                   `json:"code"`
	Message string                `json:"message"`
	Result  *CicoOTPConfirmResult `json:"result"`
}

type CicoOTPConfirmResult struct {
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

type OTPpayload struct {
	CommandID int    `json:"commandId"`
	Payload   string `json:"payload"`
}

func (s *Svc) CicoOTPConfirm(ctx context.Context, sr CicoOTPConfirmRequest) (*CicoOTPConfirmResponse, error) {
	res, err := s.cicoPost(ctx, s.cicoURL("otp"), sr)
	if err != nil {
		return nil, err
	}

	CicoRes := &CicoOTPConfirmResponse{}
	if err := json.Unmarshal(res, CicoRes); err != nil {
		return nil, err
	}
	return CicoRes, nil
}

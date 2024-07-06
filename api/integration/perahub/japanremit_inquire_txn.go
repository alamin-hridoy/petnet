package perahub

import (
	"context"
	"encoding/json"
)

const (
	JapanRemitAwaitPayment = "Available For Pickup"
)

type JPRInquireRequest struct {
	Branch          string      `json:"branch"`
	OutletCode      string      `json:"outlet_code"`
	ReferenceNumber string      `json:"reference_number"`
	ControlNumber   string      `json:"control_number"`
	LocationID      json.Number `json:"location_id"`
	UserID          int         `json:"user_id"`
	LocationName    string      `json:"location_name"`
	DeviceId        string      `json:"device_id"`
	AgentId         string      `json:"agent_id"`
	AgentCode       string      `json:"agent_code"`
	BranchCode      string      `json:"branch_code"`
	LocationCode    string      `json:"location_code"`
	Currency        string      `json:"currency"`
}

type JPRInquireResponseBody struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Result  JPRResult   `json:"result"`
	RemcoID json.Number `json:"remco_id"`
}

type JPRResult struct {
	ControlNumber      string `json:"control_number"`
	ReferenceNumber    string `json:"reference_number"`
	OriginatingCountry string `json:"originating_country"`
	DestinationCountry string `json:"destination_country"`
	SenderName         string `json:"sender_name"`
	ReceiverName       string `json:"receiver_name"`
	PrincipalAmount    string `json:"principal_amount"`
	Currency           string `json:"currency"`
	PayTokenId         string `json:"pay_token_id"`
}

func (s *Svc) JPRInquire(ctx context.Context, sr JPRInquireRequest) (*JPRInquireResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("japanremit/inquire"), sr)
	if err != nil {
		return nil, err
	}

	JPRRes := &JPRInquireResponseBody{}
	if err := json.Unmarshal(res, JPRRes); err != nil {
		return nil, err
	}
	return JPRRes, nil
}

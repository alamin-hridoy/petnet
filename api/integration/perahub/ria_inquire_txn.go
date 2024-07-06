package perahub

import (
	"context"
	"encoding/json"
)

const (
	RiaAwaitPayment = "Order is available for payout."
)

type RiaInquireRequest struct {
	DeviceID        string      `json:"device_id"`
	AgentID         string      `json:"agent_id"`
	AgentCode       string      `json:"agent_code"`
	ReferenceNumber string      `json:"reference_number"`
	ControlNumber   string      `json:"control_number"`
	LocationID      json.Number `json:"location_id"`
	UserID          json.Number `json:"user_id"`
	LocationName    string      `json:"location_name"`
}

type RiaInquireResponseBody struct {
	Code    string      `json:"code"`
	Msg     string      `json:"message"`
	Result  RiaResult   `json:"result"`
	RemcoID json.Number `json:"remco_id"`
}

type RiaResult struct {
	ControlNumber      string      `json:"control_number"`
	ClientReferenceNo  string      `json:"client_reference_no"`
	OriginatingCountry string      `json:"originating_country"`
	DestinationCountry string      `json:"destination_country"`
	SenderName         string      `json:"sender_name"`
	ReceiverName       string      `json:"receiver_name"`
	PrincipalAmount    string      `json:"principal_amount"`
	Currency           string      `json:"currency"`
	IsDomestic         json.Number `json:"is_domestic"`
	OrderNo            string      `json:"order_number"`
}

func (s *Svc) RiaInquire(ctx context.Context, sr RiaInquireRequest) (*RiaInquireResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("ria/inquire"), sr)
	if err != nil {
		return nil, err
	}

	RiaRes := &RiaInquireResponseBody{}
	if err := json.Unmarshal(res, RiaRes); err != nil {
		return nil, err
	}
	return RiaRes, nil
}

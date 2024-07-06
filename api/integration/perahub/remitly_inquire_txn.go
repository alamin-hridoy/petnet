package perahub

import (
	"context"
	"encoding/json"
)

const (
	RMAwaitPayment = "PAYABLE"
)

type RMInquireRequest struct {
	Branch       string      `json:"branch"`
	OutletCode   string      `json:"outlet_code"`
	RefNo        string      `json:"reference_number"`
	ControlNo    string      `json:"control_number"`
	LocationID   json.Number `json:"location_id"`
	UserID       json.Number `json:"user_id"`
	LocationName string      `json:"location_name"`
	DeviceID     string      `json:"device_id"`
	AgentID      string      `json:"agent_id"`
	AgentCode    string      `json:"agent_code"`
	BranchCode   string      `json:"branch_code"`
	LocationCode string      `json:"location_code"`
	CurrencyCode string      `json:"currency"`
}

type RMInquireResponseBody struct {
	Code    json.Number `json:"code"`
	Msg     string      `json:"message"`
	Result  RMInqResult `json:"result"`
	RemcoID json.Number `json:"remco_id"`
}

type RMInqResult struct {
	ControlNo     string `json:"control_number"`
	PnplAmt       string `json:"principal_amount"`
	CurrencyCode  string `json:"currency"`
	RcvName       string `json:"receiver_name"`
	Address       string `json:"address"`
	City          string `json:"city"`
	Country       string `json:"country"`
	SenderName    string `json:"sender_name"`
	ContactNumber string `json:"contact_number"`
}

func (s *Svc) RMInquire(ctx context.Context, sr RMInquireRequest) (*RMInquireResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("remitly/inquire"), sr)
	if err != nil {
		return nil, err
	}

	rb := &RMInquireResponseBody{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

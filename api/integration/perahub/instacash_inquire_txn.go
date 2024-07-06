package perahub

import (
	"context"
	"encoding/json"
)

const (
	InstaCashAwaitPayment = "Client Details"
)

type InstaCashInquireRequest struct {
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
}

type InstaCashInquireResponseBody struct {
	Code    json.Number     `json:"code"`
	Message string          `json:"message"`
	Result  InstaCashResult `json:"result"`
	RemcoID json.Number     `json:"remco_id"`
}

type InstaCashResult struct {
	ControlNumber      string      `json:"control_number"`
	ReferenceNumber    string      `json:"reference_number"`
	OriginatingCountry string      `json:"originating_country"`
	DestinationCountry string      `json:"destination_country"`
	SenderName         string      `json:"sender_name"`
	ReceiverName       string      `json:"receiver_name"`
	PrincipalAmount    json.Number `json:"principal_amount"`
	Currency           string      `json:"currency"`
	Purpose            string      `json:"purpose"`
	Status             string      `json:"status"`
}

func (s *Svc) InstaCashInquire(ctx context.Context, sr InstaCashInquireRequest) (*InstaCashInquireResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("instacash/inquire"), sr)
	if err != nil {
		return nil, err
	}

	ICRes := &InstaCashInquireResponseBody{}
	if err := json.Unmarshal(res, ICRes); err != nil {
		return nil, err
	}
	return ICRes, nil
}

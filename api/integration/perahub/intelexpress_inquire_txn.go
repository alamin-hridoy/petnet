package perahub

import (
	"context"
	"encoding/json"
)

const (
	IEAwaitPayment = "Success"
)

type IEInquireRequest struct {
	Branch          string      `json:"branch"`
	OutletCode      string      `json:"outlet_code"`
	ReferenceNumber string      `json:"reference_number"`
	ControlNumber   string      `json:"control_number"`
	LocationID      json.Number `json:"location_id"`
	UserID          json.Number `json:"user_id"`
	LocationName    string      `json:"location_name"`
	DeviceID        string      `json:"device_id"`
	AgentID         string      `json:"agent_id"`
	AgentCode       string      `json:"agent_code"`
	BranchCode      string      `json:"branch_code"`
	LocationCode    string      `json:"location_code"`
	Currency        string      `json:"currency"`
}

type IEInquireResponse struct {
	Code    json.Number     `json:"code"`
	Message string          `json:"message"`
	Result  IEInquireResult `json:"result"`
	RemcoID json.Number     `json:"remco_id"`
}

type IEInquireResult struct {
	ControlNumber      string `json:"control_number"`
	TrxDate            string `json:"trx_date"`
	PrincipalAmount    string `json:"principal_amount"`
	Currency           string `json:"currency"`
	ReceiverName       string `json:"receiver_name"`
	SenderName         string `json:"sender_name"`
	Address            string `json:"address"`
	Country            string `json:"country"`
	OriginatingCountry string `json:"originating_country"`
	DestinationCountry string `json:"destination_country"`
	ContactNumber      string `json:"contact_number"`
	ReferenceNumber    string `json:"reference_number"`
}

func (s *Svc) IEInquire(ctx context.Context, sr IEInquireRequest) (*IEInquireResponse, error) {
	res, err := s.postNonex(ctx, s.nonexURL("intelexpress/inquire"), sr)
	if err != nil {
		return nil, err
	}

	IERes := &IEInquireResponse{}
	if err := json.Unmarshal(res, IERes); err != nil {
		return nil, err
	}
	return IERes, nil
}

package perahub

import (
	"context"
	"encoding/json"
)

const (
	UnitellerAwaitPayment = "nonex success"
)

type UNTInquireRequest struct {
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

type UNTInquireResponseBody struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Result  UNTResult   `json:"result"`
	RemcoID json.Number `json:"remco_id"`
}

type UNTResult struct {
	ResponseCode       string      `json:"response_code"`
	ControlNumber      string      `json:"control_number"`
	PrincipalAmount    json.Number `json:"principal_amount"`
	Currency           string      `json:"currency"`
	CreationDate       string      `json:"creation_date"`
	ReceiverName       string      `json:"receiver_name"`
	Address            string      `json:"address"`
	City               string      `json:"city"`
	Country            string      `json:"country"`
	SenderName         string      `json:"sender_name"`
	ZipCode            string      `json:"zip_code"`
	OriginatingCountry string      `json:"originating_country"`
	DestinationCountry string      `json:"destination_country"`
	ContactNumber      string      `json:"contact_number"`
	FmtSenderName      string      `json:"formatted_sender_name"`
	FmtReceiverName    string      `json:"formatted_receiver_name"`
}

func (s *Svc) UNTInquire(ctx context.Context, sr UNTInquireRequest) (*UNTInquireResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("uniteller/inquire"), sr)
	if err != nil {
		return nil, err
	}

	UNTRes := &UNTInquireResponseBody{}
	if err := json.Unmarshal(res, UNTRes); err != nil {
		return nil, err
	}
	return UNTRes, nil
}

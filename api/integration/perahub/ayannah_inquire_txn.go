package perahub

import (
	"context"
	"encoding/json"
)

const (
	AYANNAHAwaitPayment = "Transaction Available For Payout."
)

type AYANNAHInquireRequest struct {
	Branch          string      `json:"branch"`
	OutletCode      string      `json:"outlet_code"`
	ReferenceNumber string      `json:"reference_number"`
	ControlNumber   string      `json:"control_number"`
	LocationID      json.Number `json:"location_id"`
	UserID          json.Number `json:"user_id"`
	LocationName    string      `json:"location_name"`
	DeviceID        string      `json:"device_id"`
	AgentCode       string      `json:"agent_code"`
	LocationCode    string      `json:"location_code"`
	Currency        string      `json:"currency"`
}

type AYANNAHInquireResponseBody struct {
	Code    json.Number          `json:"code"`
	Message string               `json:"message"`
	Result  AYANNAHInquireResult `json:"result"`
	RemcoID json.Number          `json:"remco_id"`
}

type AYANNAHInquireResult struct {
	ResponseCode       string `json:"response_code"`
	ResponseMessage    string `json:"response_message"`
	ControlNumber      string `json:"control_number"`
	PrincipalAmount    string `json:"principal_amount"`
	Currency           string `json:"currency"`
	CreationDate       string `json:"creation_date"`
	ReceiverName       string `json:"receiver_name"`
	SenderName         string `json:"sender_name"`
	Address            string `json:"address"`
	City               string `json:"city"`
	Country            string `json:"country"`
	ZipCode            string `json:"zip_code"`
	OriginatingCountry string `json:"originating_country"`
	DestinationCountry string `json:"destination_country"`
	ContactNumber      string `json:"contact_number"`
	ReferenceNumber    string `json:"reference_number"`
}

func (s *Svc) AYANNAHInquire(ctx context.Context, sr AYANNAHInquireRequest) (*AYANNAHInquireResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("ayannah/inquire"), sr)
	if err != nil {
		return nil, err
	}

	cebRes := &AYANNAHInquireResponseBody{}
	if err := json.Unmarshal(res, cebRes); err != nil {
		return nil, err
	}
	return cebRes, nil
}

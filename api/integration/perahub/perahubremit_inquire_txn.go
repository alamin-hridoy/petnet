package perahub

import (
	"context"
	"encoding/json"
)

type PerahubRemitInquireRequest struct {
	ControlNumber string `json:"control_number"`
	LocationID    int    `json:"location_id"`
}

type PerahubRemitInquireResult struct {
	PrincipalAmount    int    `json:"principal_amount"`
	IsoCurrency        string `json:"iso_currency"`
	ConversionRate     int    `json:"conversion_rate"`
	SenderLastName     string `json:"sender_last_name"`
	SenderFirstName    string `json:"sender_first_name"`
	SenderMiddleName   string `json:"sender_middle_name"`
	ReceiverLastName   string `json:"receiver_last_name"`
	ReceiverFirstName  string `json:"receiver_first_name"`
	ReceiverMiddleName string `json:"receiver_middle_name"`
	ControlNumber      string `json:"control_number"`
	OriginatingCountry string `json:"originating_country"`
	DestinationCountry string `json:"destination_country"`
	SenderName         string `json:"sender_name"`
	ReceiverName       string `json:"receiver_name"`
	PartnerCode        string `json:"partner_code"`
}

type PerahubRemitInquireResponse struct {
	Code    int                       `json:"code"`
	Message string                    `json:"message"`
	Result  PerahubRemitInquireResult `json:"result"`
}

func (s *Svc) PerahubRemitInquire(ctx context.Context, req PerahubRemitInquireRequest) (*PerahubRemitInquireResponse, error) {
	res, err := s.postNonex(ctx, s.nonexURL("/perahub-remit/payout/inquire"), req)
	if err != nil {
		return nil, err
	}

	inqRes := &PerahubRemitInquireResponse{}
	if err := json.Unmarshal(res, inqRes); err != nil {
		return nil, err
	}
	return inqRes, nil
}

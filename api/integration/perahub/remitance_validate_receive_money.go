package perahub

import (
	"context"
	"encoding/json"
)

type RemitanceValidateReceiveMoneyReq struct {
	Phrn                  string `json:"phrn"`
	PrincipalAmount       string `json:"principal_amount"`
	IsoOriginatingCountry string `json:"iso_originating_country"`
	IsoDestinationCountry string `json:"iso_destination_country"`
	SenderLastName        string `json:"sender_last_name"`
	SenderFirstName       string `json:"sender_first_name"`
	SenderMiddleName      string `json:"sender_middle_name"`
	ReceiverLastName      string `json:"receiver_last_name"`
	ReceiverFirstName     string `json:"receiver_first_name"`
	ReceiverMiddleName    string `json:"receiver_middle_name"`
	PayoutPartnerCode     string `json:"payout_partner_code"`
}

type RemitanceValidateReceiveMoneyResult struct {
	PayoutValidateReferenceNumber string `json:"payout_validate_reference_number"`
}

type RemitanceValidateReceiveMoneyRes struct {
	Code    int                                 `json:"code"`
	Message string                              `json:"message"`
	Result  RemitanceValidateReceiveMoneyResult `json:"result"`
}

func (s *Svc) RemitanceValidateReceiveMoney(ctx context.Context, req RemitanceValidateReceiveMoneyReq) (*RemitanceValidateReceiveMoneyRes, error) {
	res, err := s.remitancePost(ctx, s.remitanceURL("receive/validate"), req)
	if err != nil {
		return nil, err
	}

	rb := &RemitanceValidateReceiveMoneyRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

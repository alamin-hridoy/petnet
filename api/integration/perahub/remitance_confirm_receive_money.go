package perahub

import (
	"context"
	"encoding/json"
)

type RemitanceConfirmReceiveMoneyReq struct {
	PayoutValidateReferenceNumber string `json:"payout_validate_reference_number"`
}

type RemitanceConfirmReceiveMoneyResult struct {
	Phrn                  string `json:"phrn"`
	PrincipalAmount       int    `json:"principal_amount"`
	IsoOriginatingCountry string `json:"iso_originating_country"`
	IsoDestinationCountry string `json:"iso_destination_country"`
	SenderLastName        string `json:"sender_last_name"`
	SenderFirstName       string `json:"sender_first_name"`
	SenderMiddleName      string `json:"sender_middle_name"`
	ReceiverLastName      string `json:"receiver_last_name"`
	ReceiverFirstName     string `json:"receiver_first_name"`
	ReceiverMiddleName    string `json:"receiver_middle_name"`
}

type RemitanceConfirmReceiveMoneyRes struct {
	Code    int                                `json:"code"`
	Message string                             `json:"message"`
	Result  RemitanceConfirmReceiveMoneyResult `json:"result"`
}

func (s *Svc) RemitanceConfirmReceiveMoney(ctx context.Context, req RemitanceConfirmReceiveMoneyReq) (*RemitanceConfirmReceiveMoneyRes, error) {
	res, err := s.remitancePost(ctx, s.remitanceURL("receive/confirm"), req)
	if err != nil {
		return nil, err
	}

	rb := &RemitanceConfirmReceiveMoneyRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

package perahub

import (
	"context"
	"encoding/json"
)

type RemitanceInquireReq struct {
	Phrn string `json:"phrn"`
}

type RemitanceInquireResult struct {
	Phrn                  string `json:"phrn"`
	PrincipalAmount       int    `json:"principal_amount"`
	IsoCurrency           string `json:"iso_currency"`
	ConversionRate        int    `json:"conversion_rate"`
	IsoOriginatingCountry string `json:"iso_originating_country"`
	IsoDestinationCountry string `json:"iso_destination_country"`
	SenderLastName        string `json:"sender_last_name"`
	SenderFirstName       string `json:"sender_first_name"`
	SenderMiddleName      string `json:"sender_middle_name"`
	ReceiverLastName      string `json:"receiver_last_name"`
	ReceiverFirstName     string `json:"receiver_first_name"`
	ReceiverMiddleName    string `json:"receiver_middle_name"`
}

type RemitanceInquireRes struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Result  RemitanceInquireResult `json:"result"`
}

func (s *Svc) RemitanceInquire(ctx context.Context, req RemitanceInquireReq) (*RemitanceInquireRes, error) {
	res, err := s.remitancePost(ctx, s.remitanceURL("inquire"), req)
	if err != nil {
		return nil, err
	}

	ri := &RemitanceInquireRes{}
	if err := json.Unmarshal(res, ri); err != nil {
		return nil, err
	}
	return ri, nil
}

package perahub

import (
	"context"
	"encoding/json"
)

type RemitanceValidateSendMoneyReq struct {
	PartnerReferenceNumber string `json:"partner_reference_number"`
	PrincipalAmount        string `json:"principal_amount"`
	ServiceFee             string `json:"service_fee"`
	IsoCurrency            string `json:"iso_currency"`
	ConversionRate         string `json:"conversion_rate"`
	IsoOriginatingCountry  string `json:"iso_originating_country"`
	IsoDestinationCountry  string `json:"iso_destination_country"`
	SenderLastName         string `json:"sender_last_name"`
	SenderFirstName        string `json:"sender_first_name"`
	SenderMiddleName       string `json:"sender_middle_name"`
	ReceiverLastName       string `json:"receiver_last_name"`
	ReceiverFirstName      string `json:"receiver_first_name"`
	ReceiverMiddleName     string `json:"receiver_middle_name"`
	SenderBirthDate        string `json:"sender_birth_date"`
	SenderBirthPlace       string `json:"sender_birth_place"`
	SenderBirthCountry     string `json:"sender_birth_country"`
	SenderGender           string `json:"sender_gender"`
	SenderRelationship     string `json:"sender_relationship"`
	SenderPurpose          string `json:"sender_purpose"`
	SenderOfFund           string `json:"sender_of_fund"`
	SenderOccupation       string `json:"sender_occupation"`
	SenderEmploymentNature string `json:"sender_employment_nature"`
	SendPartnerCode        string `json:"send_partner_code"`
	SenderSourceOfFund     string `json:"sender_source_of_fund"`
}

type RemitanceValidateSendMoneyResult struct {
	SendValidateReferenceNumber string `json:"send_validate_reference_number"`
}

type RemitanceValidateSendMoneyRes struct {
	Code    int                              `json:"code"`
	Message string                           `json:"message"`
	Result  RemitanceValidateSendMoneyResult `json:"result"`
}

func (s *Svc) RemitanceValidateSendMoney(ctx context.Context, req RemitanceValidateSendMoneyReq) (*RemitanceValidateSendMoneyRes, error) {
	res, err := s.remitancePost(ctx, s.remitanceURL("send/validate"), req)
	if err != nil {
		return nil, err
	}

	rb := &RemitanceValidateSendMoneyRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

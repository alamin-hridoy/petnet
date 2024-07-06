package perahub

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
)

type PaySIRequest struct {
	MTCN                    string `json:"mtcn"`
	OriginatingCurrencyCode string `json:"originating_currency_code"`
	OriginatingCountryCode  string `json:"originating_country_code"`
	ForeignReferenceNo      string `json:"foreign_reference_no"`
	TerminalID              string `json:"terminal_id"`
	OperatorID              string `json:"operator_id"`
}

type Name struct {
	NameType  string `json:"name_type"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Sender struct {
	Name Name `json:"name"`
}

type Receiver struct {
	Name Name `json:"name"`
}

type IsoCode struct {
	CountryCode  string `json:"country_code"`
	CurrencyCode string `json:"currency_code"`
}

type OriginatingCountryCurrency struct {
	IsoCode IsoCode `json:"iso_code"`
}

type InquiryPaymentDetails struct {
	OriginatingCountryCurrency OriginatingCountryCurrency `json:"originating_country_currency"`
}

type PaymentTransaction struct {
	Sender               Sender                `json:"sender"`
	Receiver             Receiver              `json:"receiver"`
	Financials           Financials            `json:"financials"`
	PaymentDetails       InquiryPaymentDetails `json:"payment_details"`
	FilingDate           string                `json:"filing_date"`
	FilingTime           string                `json:"filing_time"`
	MoneyTransferKey     string                `json:"money_transfer_key"`
	PayStatusDescription string                `json:"pay_status_description"`
}

type PaymentTransactions struct {
	PaymentTransaction PaymentTransaction `json:"payment_transaction"`
}

type ForeignRemoteSystem struct {
	Identifier  string `json:"identifier"`
	ReferenceNo string `json:"reference_no"`
	CounterID   string `json:"counter_id"`
}

type PaySIResponseBody struct {
	PaymentTransactions PaymentTransactions `json:"payment_transactions"`
	ForeignRemoteSystem ForeignRemoteSystem `json:"foreign_remote_system"`
	NumberMatches       int                 `json:"number_matches"`
	CurrentPageNumber   int                 `json:"current_page_number"`
	LastPageNumber      int                 `json:"last_page_number"`
}

type PaySIResponseWU struct {
	Header ResponseHeader    `json:"header"`
	Body   PaySIResponseBody `json:"body"`
}

type PaySIResponse struct {
	WU PaySIResponseWU `json:"uspwuapi"`
}

func (s *Svc) PayStatusInquiry(ctx context.Context, psiReq PaySIRequest) (*PaySIResponse, error) {
	req, err := s.newParahubRequest(ctx, "wupo", "checkstat", psiReq)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := s.cl.Post(s.moduleURL("wupo", "checkstat"), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var psiRes PaySIResponse
	if err := json.Unmarshal(body, &psiRes); err != nil {
		return nil, err
	}

	return &psiRes, nil
}

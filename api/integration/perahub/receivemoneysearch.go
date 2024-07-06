package perahub

import (
	"context"
	"encoding/json"
)

type RMSearchRequest struct {
	FrgnRefNo    string `json:"foreign_reference_no"`
	MTCN         string `json:"mtcn"`
	DestCurrency string `json:"pd_dest_currency_code"`
	TerminalID   string `json:"terminal_id"`
	OperatorID   string `json:"operator_id"`
}

type RMSearchResponseBody struct {
	Txn         RMSPaymentTransaction `json:"payment_transaction"`
	GrossPayout json.Number           `json:"gross_payout"`
	DST         json.Number           `json:"dst"`
	NetPayout   json.Number           `json:"net_payout"`
}

type RMSPaymentTransaction struct {
	Sender                  Contact           `json:"sender"`
	Receiver                Contact           `json:"receiver"`
	Financials              RMSFinancials     `json:"financials"`
	Payment                 RMSPaymentDetails `json:"payment_details"`
	FilingDate              string            `json:"filing_date"`
	FilingTime              string            `json:"filing_time"`
	MoneyTransferKey        string            `json:"money_transfer_key"`
	PayStatus               string            `json:"pay_status_description"`
	Mtcn                    string            `json:"mtcn"`
	NewMtcn                 string            `json:"new_mtcn"`
	Fusion                  RMSFusion         `json:"fusion"`
	WuNetworkAgentIndicator string            `json:"wu_network_agent_indicator"`
}

type Contact struct {
	Name             Name            `json:"name"`
	Address          RMSAddress      `json:"address"`
	Phone            string          `json:"contact_phone"`
	Mobile           RMSMobilePhone  `json:"mobile_phone"`
	RawMobileDetails json.RawMessage `json:"mobile_details"`
	MobileDetails    RMSMobileDetails
}

type RMSAddress struct {
	Street      string         `json:"street"`
	City        string         `json:"city"`
	State       string         `json:"state"`
	PostalCode  string         `json:"state_zip"`
	CountryCode RMSCountryCode `json:"country_code"`
}

type RMSCountryCode struct {
	IsoCode RMSIsoCode `json:"iso_code"`
}

type RMSIsoCode struct {
	Country  string `json:"country_code"`
	Currency string `json:"currency_code"`
}

type RMSPhoneNumber struct {
	CountryCode string `json:"country_code"`
	Number      string `json:"national_number"`
}

type RMSMobilePhone struct {
	RawPhone json.RawMessage `json:"phone_number"`
	Phone    RMSPhoneNumber  `json:"phone_number2"`
}

type RMSMobileDetails struct {
	CountryCode string `json:"city_code"`
	Number      string `json:"number"`
}

type RMSTaxes struct {
	TaxWorksheet string `json:"tax_worksheet"`
}

type RMSFinancials struct {
	Taxes      RMSTaxes    `json:"taxes"`
	GrossTotal json.Number `json:"gross_total_amount"`
	PayAmount  json.Number `json:"pay_amount"`
	Principal  json.Number `json:"principal_amount"`
	Charges    json.Number `json:"charges"`
	Tolls      json.Number `json:"tolls"`
}

type RMSExpectedPayoutLocation struct {
	State string `json:"state_code"`
	City  string `json:"city"`
}

type RMSCountryCurrency struct {
	IsoCode RMSIsoCode `json:"iso_code"`
}

type RMSPaymentDetails struct {
	ExpectedPayoutLocation RMSExpectedPayoutLocation `json:"expected_payout_location"`
	DestCountry            RMSCountryCurrency        `json:"destination_country_currency"`
	OrigCountry            RMSCountryCurrency        `json:"originating_country_currency"`
	OriginatingCity        string                    `json:"originating_city"`
	TransactionType        string                    `json:"transaction_type"`
	ExchangeRate           json.Number               `json:"exchange_rate"`
	SenderDestCountry      RMSCountryCurrency        `json:"original_destination_country_currency"`
}

type RMSFusion struct {
	FusionStatus  string `json:"fusion_status"`
	AccountNumber string `json:"account_number"`
}

func (s *Svc) RMSearch(ctx context.Context, sr RMSearchRequest) (*RMSearchResponseBody, error) {
	const mod, modReq = "wupo", "search"
	req, err := s.newParahubRequest(ctx, mod, modReq, sr)
	if err != nil {
		return nil, err
	}

	res, err := s.post(ctx, s.moduleURL(mod, modReq), *req)
	if err != nil {
		return nil, err
	}

	var smvRes RMSearchResponseBody
	if err := json.Unmarshal(res, &smvRes); err != nil {
		return nil, err
	}

	// ignore error as the sender and receiver phone can change data type to array which we ignore
	json.Unmarshal(smvRes.Txn.Sender.Mobile.RawPhone, &smvRes.Txn.Sender.Mobile.Phone)
	json.Unmarshal(smvRes.Txn.Receiver.Mobile.RawPhone, &smvRes.Txn.Receiver.Mobile.Phone)
	json.Unmarshal(smvRes.Txn.Sender.RawMobileDetails, &smvRes.Txn.Sender.MobileDetails)
	json.Unmarshal(smvRes.Txn.Receiver.RawMobileDetails, &smvRes.Txn.Receiver.MobileDetails)

	return &smvRes, nil
}

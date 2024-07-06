package perahub

import (
	"context"
	"encoding/json"
)

type SMVRequest struct {
	FrgnRefNo string `json:"foreign_reference_no"`

	CustomerCode    string `json:"customer_code"`
	UserCode        string `json:"user_code"`
	TransactionType string `json:"transaction_type"`

	SenderNameType   string `json:"sender_name_type"`
	SenderFirstName  string `json:"sender_first_name"`
	SenderMiddleName string `json:"sender_middle_name"`
	SenderLastName   string `json:"sender_last_name"`

	SenderAddrLine1         string `json:"sender_addr_line1"`
	SenderAddrLine2         string `json:"sender_addr_line2"`
	SenderCity              string `json:"sender_city"`
	SenderState             string `json:"sender_state"`
	SenderPostalCode        string `json:"sender_postal_code"`
	SenderAddrCountryCode   string `json:"sender_addr_country_code"`
	SenderAddrCurrencyCode  string `json:"sender_addr_currency_code"`
	SenderAddrCountryName   string `json:"sender_addr_country_name"`
	SenderContactPhone      string `json:"sender_contact_phone"`
	SenderMobileCountryCode string `json:"sender_mobile_country_code"`
	SenderMobileNo          string `json:"sender_mobile_no"`

	Email          string `json:"email"`
	Gender         string `json:"gender"`
	DateOfBirth    string `json:"date_of_birth"`
	CountryOfBirth string `json:"Country_of_Birth"`
	Nationality    string `json:"nationality"`

	IDType    string `json:"id_type"`
	IDCountry string `json:"id_country_of_issue"`
	IDNumber  string `json:"id_number"`
	IDIssued  string `json:"id_issue_date"`
	IDExpiry  string `json:"id_expiration_date"`

	Occupation              string  `json:"occupation"`
	EmploymentPositionLevel string  `json:"employment_position_level"`
	SourceOfFunds           string  `json:"Source_of_Funds"`
	ReceiverRelationship    string  `json:"Relationship_to_Receiver"`
	TransactionPurpose      string  `json:"Purpose_of_Transaction"`
	SenderReason            string  `json:"reason_for_send"`
	GalacticID              *string `json:"galactic_id"`
	MyWUNumber              string  `json:"mywu_number"`
	MyWUEnrollTag           string  `json:"mywu_enroll_tag"`

	ReceiverNameType          string `json:"receiver_name_type"`
	ReceiverFirstName         string `json:"receiver_first_name"`
	ReceiverMiddleName        string `json:"receiver_middle_name"`
	ReceiverLastName          string `json:"receiver_last_name"`
	ReceiverAddrLine1         string `json:"receiver_addr_line1"`
	ReceiverAddrLine2         string `json:"receiver_addr_line2"`
	ReceiverCity              string `json:"receiver_city"`
	ReceiverState             string `json:"receiver_state"`
	ReceiverPostalCode        string `json:"receiver_postal_code"`
	ReceiverAddrCurrencyCode  string `json:"receiver_addr_currency_code"`
	ReceiverAddrCountryCode   string `json:"receiver_addr_country_code"`
	ReceiverMobileCountryCode string `json:"receiver_mobile_country_code"`
	ReceiverMobileNo          string `json:"receiver_mobile_no"`
	ReceiverContactPhone      string `json:"receiver_contact_phone"`

	PromoCode       string      `json:"promo_code"`
	PrincipalAmount json.Number `json:"principal_amount"`
	ExchangeRate    string      `json:"exchange_rate"`
	FixedAmountFlag string      `json:"fixed_amount_flag"`

	DestCurrency  string `json:"destination_currency_code"`
	DestCountry   string `json:"destination_country_code"`
	DestCity      string `json:"destination_city_name"`
	DestStateCode string `json:"destination_state_code"`

	Message []string `json:"message"`

	BankName      string `json:"bank_name"`
	BankLocation  string `json:"bank_location"`
	BankCode      string `json:"bank_code"`
	AccountSuffix string `json:"account_suffix"`
	AccountNumber string `json:"account_number"`

	IsOnBehalf string `json:"is_on_behalf"`
	AckFlag    string `json:"ack_flag"`

	TerminalID       string `json:"terminal_id"`
	OperatorID       string `json:"operator_id"`
	RemoteTerminalID string `json:"remote_terminal_id"`
	RemoteOperatorID string `json:"remote_operator_id"`
}

type SMVResponseBody struct {
	ServiceCode       ServiceCode       `json:"service_code"`
	Compliance        Compliance        `json:"compliance"`
	PaymentDetails    PaymentDetails    `json:"payment_details"`
	Fin               Financials        `json:"financials"`
	Promotions        Promotions        `json:"promotions"`
	NewDetails        NewDetails        `json:"new_details"`
	PreferredCustomer PreferredCustomer `json:"preferred_customer"`
}

type ServiceCode struct {
	AddlSvcChg string `json:"addl_service_charges"`
}

type Compliance struct {
	ComplianceBuf string `json:"compliance_data_buffer"`
}

type PaymentDetails struct {
	OrigCity  string `json:"originating_city"`
	OrigState string `json:"originating_state"`
}

type Financials struct {
	Taxes         Taxes       `json:"taxes"`
	OrigPnplAmt   json.Number `json:"originators_principal_amount"`
	DestPcplAmt   json.Number `json:"destination_principal_amount"`
	GrossTotal    json.Number `json:"gross_total_amount"`
	Charges       json.Number `json:"charges"`
	AddlCharges   json.Number `json:"plus_charges_amount"`
	MessageCharge json.Number `json:"message_charge"`
}

type Taxes struct {
	// Currency  string      `json:"currency"`
	MuniTax   json.Number `json:"municipal_tax"`
	StateTax  json.Number `json:"state_tax"`
	CountyTax json.Number `json:"county_tax"`
	// Total     json.Number `json:"total"`
}

type Promotions struct {
	PromoDesc       string `json:"promo_code_description"`
	PromoMessage    string `json:"promo_message"`
	SenderPromoCode string `json:"sender_promo_code"`
}

type NewDetails struct {
	MTCN       string `json:"mtcn"`
	NewMTCN    string `json:"new_mtcn"`
	FilingDate string `json:"filing_date"`
	FilingTime string `json:"filing_time"`
}

type PreferredCustomer struct {
	MyWUNumber string `json:"mywu_number"`
}

func (s *Svc) SendMoneyValidate(ctx context.Context, smvr SMVRequest) (*SMVResponseBody, error) {
	// Static elements
	smvr.ReceiverNameType = "D"
	smvr.SenderNameType = "D"
	smvr.MyWUEnrollTag = "none"

	const mod, reqMod, modReq = "wuso", "validate", "SMvalidate"
	req, err := s.newParahubRequest(ctx,
		mod, modReq, smvr,
		WithUserCode(json.Number(smvr.UserCode)),
		WithLocationCode(smvr.UserCode))
	if err != nil {
		return nil, err
	}

	res, err := s.post(ctx, s.moduleURL(mod, reqMod), *req)
	if err != nil {
		return nil, err
	}

	var smvRes SMVResponseBody
	if err := json.Unmarshal(res, &smvRes); err != nil {
		return nil, err
	}

	return &smvRes, nil
}

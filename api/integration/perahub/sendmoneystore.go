package perahub

import (
	"context"
	"encoding/json"
)

type SMStoreRequest struct {
	FrgnRefNo    string `json:"foreign_reference_no"`
	UserCode     string `json:"user_code"`
	CustomerCode string `json:"customer_code"`

	SenderNameType          string `json:"sender_name_type"`
	SenderFirstName         string `json:"sender_first_name"`
	SenderMiddleName        string `json:"sender_middle_name"`
	SenderLastName          string `json:"sender_last_name"`
	SenderAddrCountryCode   string `json:"sender_addr_country_code"`
	SenderAddrCurrencyCode  string `json:"sender_addr_currency_code"`
	SenderContactPhone      string `json:"sender_contact_phone"`
	SenderMobileCountryCode string `json:"sender_mobile_country_code"`
	SenderMobileNo          string `json:"sender_mobile_no"`

	SenderReason string `json:"reason_for_send"`
	Email        string `json:"email"`

	ReceiverNameType          string `json:"receiver_name_type"`
	ReceiverFirstName         string `json:"receiver_first_name"`
	ReceiverMiddleName        string `json:"receiver_middle_name"`
	ReceiverLastName          string `json:"receiver_last_name"`
	ReceiverAddrLine1         string `json:"receiver_addr_line1"`
	ReceiverAddrLine2         string `json:"receiver_addr_line2"`
	ReceiverCity              string `json:"receiver_city"`
	ReceiverState             string `json:"receiver_state"`
	ReceiverPostalCode        string `json:"receiver_postal_code"`
	ReceiverAddrCountryCode   string `json:"receiver_addr_country_code"`
	ReceiverAddrCurrencyCode  string `json:"receiver_addr_currency_code"`
	ReceiverContactPhone      string `json:"receiver_contact_phone"`
	ReceiverMobileCountryCode string `json:"receiver_mobile_country_code"`
	ReceiverMobileNo          string `json:"receiver_mobile_no"`

	DestinationCountryCode  string `json:"destination_country_code"`
	DestinationCurrencyCode string `json:"destination_currency_code"`
	DestinationStateCode    string `json:"destination_state_code"`
	DestinationCityName     string `json:"destination_city_name"`

	TransactionType string      `json:"transaction_type"`
	PrincipalAmount json.Number `json:"principal_amount"`
	FixedAmountFlag string      `json:"fixed_amount_flag"`
	PromoCode       string      `json:"promo_code"`
	Message         []string    `json:"message"`

	AddlServiceChg       string `json:"addl_service_charges"`
	ComplianceDataBuffer string `json:"compliance_data_buffer"`

	OriginatingCity  string `json:"originating_city"`
	OriginatingState string `json:"originating_state"`
	MTCN             string `json:"mtcn"`
	NewMTCN          string `json:"new_mtcn"`
	ExchangeRate     string `json:"exchange_rate"`

	Fin               Financials        `json:"financials"`
	Promo             Promotions        `json:"promotions"`
	ComplianceDetails ComplianceDetails `json:"compliance_details"`

	BankName      string `json:"bank_name"`
	BankLocation  string `json:"bank_location"`
	AccountNumber string `json:"account_number"`
	BankCode      string `json:"bank_code"`
	AccountSuffix string `json:"account_suffix"`

	MyWUNumber string `json:"mywu_number"`
	WUEnroll   string `json:"mywu_enroll_tag"`

	TerminalID       string `json:"terminal_id"`
	OperatorID       string `json:"operator_id"`
	RemoteTerminalID string `json:"remote_terminal_id"`
	RemoteOperatorID string `json:"remote_operator_id"`
}

type SMSTaxes struct {
	MuniTax   int `json:"municipal_tax"`
	StateTax  int `json:"state_tax"`
	CountyTax int `json:"county_tax"`
}

type IDDetails struct {
	IDType    string `json:"id_type"`
	IDCountry string `json:"id_country_of_issue"`
	IDNumber  string `json:"id_number"`
}

type Address struct {
	AddrLine1  string `json:"addr_line1"`
	AddrLine2  string `json:"addr_line2"`
	City       string `json:"city"`
	State      string `json:"state_name"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

type ComplianceDetails struct {
	IDDetails               IDDetails `json:"id_details"`
	IDIssued                string    `json:"id_issue_date"`
	IDExpiry                string    `json:"id_expiration_date"`
	Birthdate               string    `json:"date_of_birth"`
	BirthCountry            string    `json:"Country_of_Birth"`
	Nationality             string    `json:"nationality"`
	SourceFunds             string    `json:"Source_of_Funds"`
	Occupation              string    `json:"occupation"`
	TransactionReason       string    `json:"transaction_reason"`
	CurrentAddress          Address   `json:"Current_address"`
	TxnRelationship         string    `json:"Relationship_to_Receiver_Sender"`
	IActOnMyBehalf          string    `json:"I_act_on_My_Behalf"`
	GalacticID              string    `json:"galactic_id,omitempty"`
	EmploymentPositionLevel string    `json:"employment_position_level"`
}

type ConfirmedDetails struct {
	AdvisoryText     string   `json:"advisory_text"`
	MTCN             string   `json:"mtcn"`
	NewMTCN          string   `json:"new_mtcn"`
	FilingDate       string   `json:"filing_date"`
	FilingTime       string   `json:"filing_time"`
	PinMessage       []string `json:"pin_message"`
	PromoTextMessage []string `json:"promo_text_message"`
	MyWUNumber       string   `json:"mywu_number"`
	NewPointsEarned  int      `json:"new_points_earned"`
	OtherMessage1    []string `json:"other_message_1"`
	OtherMessage2    []string `json:"other_message_2"`
}

func (s *Svc) SendMoneyStore(ctx context.Context, smsr SMStoreRequest) (*ConfirmedDetails, error) {
	smsr.SenderNameType = "D"
	smsr.ReceiverNameType = "D"
	smsr.WUEnroll = "none"

	const mod, reqMod, modReq = "wuso", "store", "SMstore"
	req, err := s.newParahubRequest(ctx,
		mod, modReq, smsr,
		WithUserCode(json.Number(smsr.UserCode)),
		WithLocationCode(smsr.UserCode))
	if err != nil {
		return nil, err
	}

	res, err := s.post(ctx, s.moduleURL(mod, reqMod), *req)
	if err != nil {
		return nil, err
	}

	var smsRes ConfirmedDetails
	if err := json.Unmarshal(res, &struct {
		ConfirmedDetails *ConfirmedDetails `json:"confirmed_details"`
	}{ConfirmedDetails: &smsRes}); err != nil {
		return nil, err
	}

	return &smsRes, nil
}

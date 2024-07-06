package perahub

import (
	"context"
	"encoding/json"
)

type QPVRequest struct {
	FrgnRefNo       string `json:"foreign_reference_no"`
	UserCode        string `json:"user_code"`
	CustomerCode    string `json:"customer_code"`
	TransactionType string `json:"transaction_type"`

	SenderNameType          string `json:"sender_name_type"`
	SenderFirstName         string `json:"sender_first_name"`
	SenderMiddleName        string `json:"sender_middle_name"`
	SenderLastName          string `json:"sender_last_name"`
	SenderAddrLine1         string `json:"sender_addr_line1"`
	SenderAddrLine2         string `json:"sender_addr_line2"`
	SenderCity              string `json:"sender_city"`
	SenderState             string `json:"sender_state"`
	SenderPostalCode        string `json:"sender_postal_code"`
	SenderAddrCountry       string `json:"sender_addr_country_code"`
	SenderAddrCurrency      string `json:"sender_addr_currency_code"`
	SenderAddrCountryName   string `json:"sender_addr_country_name"`
	SenderContactPhone      string `json:"sender_contact_phone"`
	SenderMobileCountryCode string `json:"sender_mobile_country_code"`
	SenderMobileNo          string `json:"sender_mobile_no"`

	Email        string `json:"email"`
	Gender       string `json:"Gender"`
	BirthCountry string `json:"Country_of_Birth"`
	Birthdate    string `json:"date_of_birth"`
	Nationality  string `json:"nationality"`

	IDType    string `json:"id_type"`
	IDCountry string `json:"id_country_of_issue"`
	IDNumber  string `json:"id_number"`
	IDIssued  string `json:"id_issue_date"`
	IDExpiry  string `json:"id_expiration_date"`

	Occupation              string `json:"occupation"`
	EmploymentPositionLevel string `json:"employment_position_level"`
	FundSource              string `json:"Source_of_Funds"`
	ReceiverRelationship    string `json:"Relationship_to_Receiver"`
	TransactionPurpose      string `json:"Purpose_of_Transaction"`

	SendingReason string  `json:"reason_for_send"`
	GalacticID    *string `json:"galactic_id"`
	MyWUNumber    string  `json:"mywu_number"`
	MyWUEnrollTag string  `json:"mywu_enroll_tag"`

	CompanyName        string `json:"company_name"`
	CompanyCode        string `json:"company_code"`
	CompanyAccountCode string `json:"company_account_code"`
	ReferenceNo        string `json:"reference_no"`

	PromoCode       string      `json:"promo_code"`
	PrincipalAmount json.Number `json:"principal_amount"`
	ExchangeRate    json.Number `json:"exchange_rate,omitempty"`
	FixedAmountFlag string      `json:"fixed_amount_flag"`

	DestCountry  string `json:"destination_country_code"`
	DestCurrency string `json:"destination_currency_code"`
	DestState    string `json:"destination_state_code"`
	DestCity     string `json:"destination_city_name"`

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

type QPVResponse struct {
	ServiceCode       QPVServiceCode       `json:"service_code"`
	Compliance        QPVCompliance        `json:"compliance"`
	PaymentDetails    QPVPaymentDetails    `json:"payment_details"`
	Fin               Financials           `json:"financials"`
	Promotions        QPVPromotions        `json:"promotions"`
	NewDetails        QPVNewDetails        `json:"new_details"`
	PreferredCustomer QPVPreferredCustomer `json:"preferred_customer"`
}

type QPVServiceCode struct {
	AddlServiceCharges string              `json:"addl_service_charges"`
	AddlServiceBlock   QPVAddlServiceBlock `json:"addl_service_block"`
}

type QPVAddlServiceBlock struct {
	AddlServiceLength     int    `json:"addl_service_length"`
	AddlServiceDataBuffer string `json:"addl_service_data_buffer"`
}

type QPVCompliance struct {
	ComplianceFlagsBuffer string `json:"compliance_flags_buffer"`
	ComplianceDataBuffer  string `json:"compliance_data_buffer"`
}

type QPVPaymentDetails struct {
	OrigCity  string `json:"originating_city"`
	OrigState string `json:"originating_state"`
}

type QPVTaxes struct {
	MuniTax   json.Number `json:"municipal_tax"`
	StateTax  json.Number `json:"state_tax"`
	CountyTax json.Number `json:"county_tax"`
}

type QPVFinancials struct {
	Taxes         Taxes       `json:"taxes"`
	OrigPrincipal json.Number `json:"originators_principal_amount"`
	DestPrincipal json.Number `json:"destination_principal_amount"`
	GrossTotal    json.Number `json:"gross_total_amount"`
	PlusCharges   json.Number `json:"plus_charges_amount"`
	Charges       json.Number `json:"charges"`
}

type QPVPromotions struct {
	PromoDesc       string      `json:"promo_code_description"`
	Message         string      `json:"promo_message"`
	DiscountAmount  json.Number `json:"promo_discount_amount"`
	PromotionError  string      `json:"promotion_error"`
	SenderPromoCode string      `json:"sender_promo_code"`
}

type QPVNewDetails struct {
	MTCN       string `json:"mtcn"`
	NewMTCN    string `json:"new_mtcn"`
	FilingDate string `json:"filing_date"`
	FilingTime string `json:"filing_time"`
}

type QPVPreferredCustomer struct {
	MyWUNumber string `json:"mywu_number"`
}

func (s *Svc) QuickPayValidate(ctx context.Context, smsr QPVRequest) (*QPVResponse, error) {
	smsr.SenderNameType = "D"

	const mod, reqMod, modReq = "wuqp", "wuqp-validate", "validate"
	req, err := s.newParahubRequest(ctx,
		mod, reqMod, smsr,
		WithUserCode(json.Number(smsr.UserCode)),
		WithLocationCode(smsr.UserCode))
	if err != nil {
		return nil, err
	}

	resp, err := s.post(ctx, s.moduleURL(mod, modReq), *req)
	if err != nil {
		return nil, err
	}

	var res QPVResponse
	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

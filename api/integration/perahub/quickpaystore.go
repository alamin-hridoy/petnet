package perahub

import (
	"context"
	"encoding/json"
)

type QPSRequest struct {
	FrgnRefNo    string `json:"foreign_reference_no"`
	UserCode     string `json:"user_code"`
	CustomerCode string `json:"customer_code"`

	SenderNameType   string `json:"sender_name_type"`
	SenderFirstName  string `json:"sender_first_name"`
	SenderMiddleName string `json:"sender_middle_name"`
	SenderLastName   string `json:"sender_last_name"`

	SenderAddrCountryCode   string `json:"sender_addr_country_code"`
	SenderAddrCurrencyCode  string `json:"sender_addr_currency_code"`
	SenderContactPhone      string `json:"sender_contact_phone"`
	SenderMobileCountryCode string `json:"sender_mobile_country_code"`
	SenderMobileNo          string `json:"sender_mobile_no"`

	SendingReason string `json:"reason_for_send"`
	Email         string `json:"email"`

	CompanyName        string `json:"company_name"`
	CompanyCode        string `json:"company_code"`
	CompanyAccountCode string `json:"company_account_code"`
	ReferenceNo        string `json:"reference_no"`

	DestCountry     string      `json:"destination_country_code"`
	DestCurrency    string      `json:"destination_currency_code"`
	DestState       string      `json:"destination_state_code"`
	DestCity        string      `json:"destination_city_name"`
	TransactionType string      `json:"transaction_type"`
	PrincipalAmount json.Number `json:"principal_amount"`
	FixedAmountFlag string      `json:"fixed_amount_flag"`

	PromoCode string   `json:"promo_code"`
	Message   []string `json:"message"`

	AddlServiceCharges   string `json:"addl_service_charges"`
	ComplianceDataBuffer string `json:"compliance_data_buffer"`

	OrigCity     string      `json:"orig_city"`
	OrigState    string      `json:"orig_state"`
	MTCN         string      `json:"mtcn"`
	NewMTCN      string      `json:"new_mtcn"`
	ExchangeRate json.Number `json:"exchange_rate"`

	Fin        Financials        `json:"financials"`
	Promo      Promotions        `json:"promotions"`
	Compliance ComplianceDetails `json:"compliance_details"`

	BankName      string `json:"bank_name"`
	BankLocation  string `json:"bank_location"`
	AccountNumber string `json:"account_number"`
	BankCode      string `json:"bank_code"`
	AccountSuffix string `json:"account_suffix"`

	MyWUNumber    string `json:"mywu_number"`
	MyWuPoints    string `json:"my_wu_current_points"`
	MyWUEnrollTag string `json:"mywu_enroll_tag"`

	TerminalID        string `json:"terminal_id"`
	OperatorID        string `json:"operator_id"`
	RemoteTerminalID  string `json:"remote_terminal_id"`
	RemoteOperatorID  string `json:"remote_operator_id"`
	ClientReferenceNo string `json:"client_reference_no"`
}

func (s *Svc) QuickPayStore(ctx context.Context, smsr QPSRequest) (*ConfirmedDetails, error) {
	smsr.SenderNameType = "D"

	const mod, reqMod, modReq = "wuqp", "wuqp-store", "store"
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

	var res ConfirmedDetails
	if err := json.Unmarshal(resp, &struct {
		ConfirmedDetails *ConfirmedDetails `json:"confirmed_details"`
	}{ConfirmedDetails: &res}); err != nil {
		return nil, err
	}

	return &res, nil
}

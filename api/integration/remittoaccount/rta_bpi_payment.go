package remittoaccount

import (
	"context"
	"encoding/json"
)

type RTABPIPaymentRequest struct {
	ReferenceNumber             string `json:"reference_number"`
	TrxDate                     string `json:"trx_date"`
	AccountNumber               string `json:"account_number"`
	Currency                    string `json:"currency"`
	ServiceCharge               string `json:"service_charge"`
	Remarks                     string `json:"remarks"`
	Particulars                 string `json:"particulars"`
	MerchantName                string `json:"merchant_name"`
	BankID                      string `json:"bank_id"`
	LocationID                  int    `json:"location_id"`
	UserID                      int    `json:"user_id"`
	CurrencyID                  string `json:"currency_id"`
	CustomerID                  string `json:"customer_id"`
	FormType                    string `json:"form_type"`
	FormNumber                  string `json:"form_number"`
	TrxType                     string `json:"trx_type"`
	RemoteLocationID            int    `json:"remote_location_id"`
	RemoteUserID                int    `json:"remote_user_id"`
	BillerName                  string `json:"biller_name"`
	TrxTime                     string `json:"trx_time"`
	TotalAmount                 string `json:"total_amount"`
	AccountName                 string `json:"account_name"`
	BeneficiaryAddress          string `json:"beneficiary_address"`
	BeneficiaryBirthdate        string `json:"beneficiary_birthdate"`
	BeneficiaryCity             string `json:"beneficiary_city"`
	BeneficiaryCivil            string `json:"beneficiary_civil"`
	BeneficiaryCountry          string `json:"beneficiary_country"`
	BeneficiaryCustomertype     string `json:"beneficiary_customertype"`
	BeneficiaryFirstname        string `json:"beneficiary_firstname"`
	BeneficiaryLastname         string `json:"beneficiary_lastname"`
	BeneficiaryMiddlename       string `json:"beneficiary_middlename"`
	BeneficiaryTin              string `json:"beneficiary_tin"`
	BeneficiarySex              string `json:"beneficiary_sex"`
	BeneficiaryState            string `json:"beneficiary_state"`
	CurrencyCodePrincipalAmount string `json:"currency_code_principal_amount"`
	PrincipalAmount             string `json:"principal_amount"`
	RecordType                  string `json:"record_type"`
	RemitterAddress             string `json:"remitter_address"`
	RemitterBirthdate           string `json:"remitter_birthdate"`
	RemitterCity                string `json:"remitter_city"`
	RemitterCivil               string `json:"remitter_civil"`
	RemitterCountry             string `json:"remitter_country"`
	RemitterCustomerType        string `json:"remitter_customer_type"`
	RemitterFirstname           string `json:"remitter_firstname"`
	RemitterGender              string `json:"remitter_gender"`
	RemitterID                  int    `json:"remitter_id"`
	RemitterLastname            string `json:"remitter_lastname"`
	RemitterMiddlename          string `json:"remitter_middlename"`
	RemitterState               string `json:"remitter_state"`
	SettlementMode              string `json:"settlement_mode"`
}

type RTABPIPaymentResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Result  RTABPIPaymentResult `json:"result"`
	RemcoID int                 `json:"remco_id"`
}

type RTABPIPaymentResult struct{}

func (c *Client) RTABPIPayment(ctx context.Context, req RTABPIPaymentRequest) (*RTABPIPaymentResponse, error) {
	res, err := c.phService.RtaPost(ctx, c.getUrl("bpi/payment"), req)
	if err != nil {
		return nil, err
	}

	bpip := &RTABPIPaymentResponse{}
	if err := json.Unmarshal(res, bpip); err != nil {
		return nil, err
	}
	return bpip, nil
}

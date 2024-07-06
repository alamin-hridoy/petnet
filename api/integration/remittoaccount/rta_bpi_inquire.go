package remittoaccount

import (
	"context"
	"encoding/json"
)

type RTABPIInquireRequest struct {
	ReferenceNumber string `json:"reference_number"`
	LocationID      string `json:"location_id"`
}

type RTABPIInquireResponse struct {
	Code     int                 `json:"code"`
	Message  string              `json:"message"`
	Result   RTABPIInquireResult `json:"result"`
	BankCode string              `json:"bank_code"`
}

type RTABPIInquireResult struct {
	BeneAmount                     string `json:"beneAmount"`
	BeneficiaryAddress             string `json:"beneficiaryAddress"`
	BeneficiaryBankAccountno       string `json:"beneficiaryBankAccountno"`
	BeneficiaryCity                string `json:"beneficiaryCity"`
	BeneficiaryCountry             string `json:"beneficiaryCountry"`
	BeneficiaryFirstName           string `json:"beneficiaryFirstName"`
	BeneficiaryLastName            string `json:"beneficiaryLastName"`
	BeneficiaryMiddleName          string `json:"beneficiaryMiddleName"`
	BeneficiaryStateOrProvince     string `json:"beneficiaryStateOrProvince"`
	BpiBranchCode                  string `json:"bpiBranchCode"`
	CurrencyCodeOfFundingAmount    string `json:"currencyCodeOfFundingAmount"`
	CurrencyCodeOfSettlementAmount string `json:"currencyCodeOfSettlementAmount"`
	TxnDistributionDate            string `json:"txnDistributionDate"`
	FundingAmount                  string `json:"fundingAmount"`
	Reason                         string `json:"reason"`
	RemitterCity                   string `json:"remitterCity"`
	RemitterCountry                string `json:"remitterCountry"`
	RemitterFirstName              string `json:"remitterFirstName"`
	RemitterLastName               string `json:"remitterLastName"`
	RemitterMiddleName             string `json:"remitterMiddleName"`
	RemitterStateOrProvince        string `json:"remitterStateOrProvince"`
	SettlementMode                 string `json:"settlementMode"`
	StatusCode                     string `json:"statusCode"`
	TransactionDate                string `json:"transactionDate"`
	TransactionReferenceNo         string `json:"transactionReferenceNo"`
}

func (c *Client) RTABPIInquire(ctx context.Context, req RTABPIInquireRequest) (*RTABPIInquireResponse, error) {
	res, err := c.phService.RtaPost(ctx, c.getUrl("bpi/inquire"), req)
	if err != nil {
		return nil, err
	}

	bpii := &RTABPIInquireResponse{}
	if err := json.Unmarshal(res, bpii); err != nil {
		return nil, err
	}
	return bpii, nil
}

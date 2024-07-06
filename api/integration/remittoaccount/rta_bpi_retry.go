package remittoaccount

import (
	"context"
	"encoding/json"
)

type RTABPIRetryRequest struct {
	ReferenceNumber string `json:"reference_number"`
	ID              string `json:"id"`
	LocationID      int    `json:"location_id"`
	FormNumber      string `json:"form_number"`
}

type RTABPIRetryResponse struct {
	Code     int               `json:"code"`
	Message  string            `json:"message"`
	Result   RTABPIRetryResult `json:"result"`
	BankCode string            `json:"bank_code"`
}

type RTABPIRetryResult struct {
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
	RemitterMessageToBeneficiary   string `json:"remitterMessageToBeneficiary"`
	RemitterMiddleName             string `json:"remitterMiddleName"`
	RemitterStateOrProvince        string `json:"remitterStateOrProvince"`
	SettlementMode                 string `json:"settlementMode"`
	StatusCode                     string `json:"statusCode"`
	TransactionDate                string `json:"transactionDate"`
	TransactionReferenceNo         string `json:"transactionReferenceNo"`
}

func (c *Client) RTABPIRetry(ctx context.Context, req RTABPIRetryRequest) (*RTABPIRetryResponse, error) {
	res, err := c.phService.RtaPost(ctx, c.getUrl("bpi/retry"), req)
	if err != nil {
		return nil, err
	}

	bpir := &RTABPIRetryResponse{}
	if err := json.Unmarshal(res, bpir); err != nil {
		return nil, err
	}
	return bpir, nil
}

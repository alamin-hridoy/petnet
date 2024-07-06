package remittoaccount

import (
	"context"
	"encoding/json"
)

type RTAUBCashinRequest struct {
	TrxType          int    `json:"trx_type"`
	BillerName       string `json:"biller_name"`
	BankID           int    `json:"bank_id"`
	LocationID       string `json:"location_id"`
	UserID           string `json:"user_id"`
	RemoteLocationID string `json:"remote_location_id"`
	RemoteUserID     string `json:"remote_user_id"`
	CurrencyID       string `json:"currency_id"`
	FormType         string `json:"form_type"`
	FormNumber       string `json:"form_number"`
	CustomerID       int    `json:"customer_id"`
	ReferenceNumber  string `json:"reference_number"`
	TrxDate          string `json:"trx_date"`
	TrxTime          string `json:"trx_time"`
	AccountNumber    string `json:"account_number"`
	Currency         string `json:"currency"`
	PrincipalAmount  string `json:"principal_amount"`
	ServiceCharge    string `json:"service_charge"`
	TotalAmount      string `json:"total_amount"`
	Remarks          string `json:"remarks"`
	Particulars      string `json:"particulars"`
	MerchantName     string `json:"merchant_name"`
	Notification     string `json:"notification"`
	AccountName      string `json:"account_name"`
	Info             []Info `json:"info"`
}

type Info struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type RTAUBCashinResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Result  RTAUBCashinResult `json:"result"`
	RemcoID int               `json:"remco_id"`
}

type RTAUBCashinResult struct {
	Code            string `json:"code"`
	SenderRefID     string `json:"senderRefId"`
	State           string `json:"state"`
	UUID            string `json:"uuid"`
	Description     string `json:"description"`
	Type            string `json:"type"`
	Amount          string `json:"amount"`
	UbpTranID       string `json:"ubpTranId"`
	TranRequestDate string `json:"tranRequestDate"`
	TranFinacleDate string `json:"tranFinacleDate"`
}

func (c *Client) RTAUBCashin(ctx context.Context, req RTAUBCashinRequest) (*RTAUBCashinResponse, error) {
	res, err := c.phService.RtaPost(ctx, c.getUrl("unionbank/cashin"), req)
	if err != nil {
		return nil, err
	}

	ubc := &RTAUBCashinResponse{}
	if err := json.Unmarshal(res, ubc); err != nil {
		return nil, err
	}
	return ubc, nil
}

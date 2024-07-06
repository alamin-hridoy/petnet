package bills_payment

import (
	"context"
	"encoding/json"
)

type BillsPaymentEcpayValidateRequest struct {
	AccountNo  string `json:"account_no"`
	Identifier string `json:"identifier"`
	BillerTag  string `json:"biller_tag"`
}

type BillsPaymentEcpayValidateResponseBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  string `json:"result"`
	RemcoID int    `json:"remco_id"`
}

func (c *Client) BillsPaymentEcpayValidate(ctx context.Context, vr BillsPaymentEcpayValidateRequest) (*BillsPaymentEcpayValidateResponseBody, error) {
	res, err := c.phService.BillsPost(ctx, c.getUrl("ecpay/validate-account"), vr)
	if err != nil {
		return nil, err
	}

	prvRes := &BillsPaymentEcpayValidateResponseBody{}
	if err := json.Unmarshal(res, prvRes); err != nil {
		return nil, err
	}
	return prvRes, nil
}

package bills_payment

import (
	"context"
	"encoding/json"
)

type BillsPaymentEcpayCheckBalanceResponseBody struct {
	Code    int                                 `json:"code"`
	Message string                              `json:"message"`
	Result  BillsPaymentEcpayCheckBalanceResult `json:"result"`
	RemcoID int                                 `json:"remco_id"`
}

type BillsPaymentEcpayCheckBalanceResult struct {
	RemBal string `json:"RemBal"`
}

func (c *Client) BillsPaymentEcpayCheckBalance(ctx context.Context) (*BillsPaymentEcpayCheckBalanceResponseBody, error) {
	res, err := c.phService.BillsGet(ctx, c.getUrl("ecpay/check-balance"))
	if err != nil {
		return nil, err
	}

	prvRes := &BillsPaymentEcpayCheckBalanceResponseBody{}
	if err := json.Unmarshal(res, prvRes); err != nil {
		return nil, err
	}
	return prvRes, nil
}

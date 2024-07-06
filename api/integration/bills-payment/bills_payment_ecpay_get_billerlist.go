package bills_payment

import (
	"context"
	"encoding/json"
)

type BillsPaymentEcpayBillerlistResponseBody struct {
	Code    int                                 `json:"code"`
	Message string                              `json:"message"`
	Result  []BillsPaymentEcpayBillerlistResult `json:"result"`
	RemcoID int                                 `json:"remco_id"`
}

type BillsPaymentEcpayBillerlistResult struct {
	BillerTag         string `json:"BillerTag"`
	Description       string `json:"Description"`
	FirstField        string `json:"FirstField"`
	FirstFieldFormat  string `json:"FirstFieldFormat"`
	FirstFieldWidth   string `json:"FirstFieldWidth"`
	SecondField       string `json:"SecondField"`
	SecondFieldFormat string `json:"SecondFieldFormat"`
	SecondFieldWidth  string `json:"SecondFieldWidth"`
	ServiceCharge     int    `json:"ServiceCharge"`
}

func (c *Client) BillsPaymentEcpayBillerlist(ctx context.Context) (*BillsPaymentEcpayBillerlistResponseBody, error) {
	res, err := c.phService.BillsGet(ctx, c.getUrl("ecpay/biller-list"))
	if err != nil {
		return nil, err
	}

	prvRes := &BillsPaymentEcpayBillerlistResponseBody{}
	if err := json.Unmarshal(res, prvRes); err != nil {
		return nil, err
	}
	return prvRes, nil
}

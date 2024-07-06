package bills_payment

import (
	"context"
	"encoding/json"
)

type BillsPaymentMultiPayBillerProcessRequest struct {
	AccountNumber string `json:"account_number"`
	Amount        int    `json:"amount"`
	ContactNumber string `json:"contact_number"`
}

type BillsPaymentMultiPayBillerProcessResponseBody struct {
	Status int                                   `json:"status"`
	Reason string                                `json:"reason"`
	Data   BillsPaymentMultiPayBillerProcessData `json:"data"`
}

type BillsPaymentMultiPayBillerProcessData struct {
	Refno  string        `json:"refno"`
	Txnid  string        `json:"txnid"`
	Biller string        `json:"biller"`
	Meta   []interface{} `json:"meta"`
}

func (c *Client) BillsPaymentMultiPayBillerProcess(ctx context.Context, vr BillsPaymentMultiPayBillerProcessRequest) (*BillsPaymentMultiPayBillerProcessResponseBody, error) {
	res, err := c.phService.BillsPost(ctx, c.getUrl("multipay/biller-process"), vr)
	if err != nil {
		return nil, err
	}

	prvRes := &BillsPaymentMultiPayBillerProcessResponseBody{}
	if err := json.Unmarshal(res, prvRes); err != nil {
		return nil, err
	}
	return prvRes, nil
}

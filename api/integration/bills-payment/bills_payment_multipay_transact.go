package bills_payment

import (
	"context"
	"encoding/json"
)

type BillsPaymentMultiPayTransactRequest struct {
	Amount string `json:"amount"`
	Txnid  string `json:"txnid"`
}

type BillsPaymentMultiPayTransactResponseBody struct {
	Data BillsPaymentMultiPayTransactData `json:"data"`
}

type BillsPaymentMultiPayTransactData struct {
	URL string `json:"url"`
}

func (c *Client) BillsPaymentMultiPayTransact(ctx context.Context, vr BillsPaymentMultiPayTransactRequest) (*BillsPaymentMultiPayTransactResponseBody, error) {
	res, err := c.phService.BillsPost(ctx, c.getUrl("multipay/generate-transaction"), vr)
	if err != nil {
		return nil, err
	}

	prvRes := &BillsPaymentMultiPayTransactResponseBody{}
	if err := json.Unmarshal(res, prvRes); err != nil {
		return nil, err
	}
	return prvRes, nil
}

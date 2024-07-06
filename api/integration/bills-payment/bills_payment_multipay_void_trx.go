package bills_payment

import (
	"context"
	"encoding/json"
)

type BillsPaymentMultiPayVoidTrxRequest struct {
	ReferenceNo string `json:"reference_no"`
}

type BillsPaymentMultiPayVoidTrxResponseBody struct {
	Data BillsPaymentMultiPayVoidTrxData `json:"data"`
}

type BillsPaymentMultiPayVoidTrxData struct {
	Txnid                string `json:"txnid"`
	Refno                string `json:"refno"`
	Amount               string `json:"amount"`
	Fee                  string `json:"fee"`
	Status               string `json:"status"`
	PaymentChannel       string `json:"payment_channel"`
	IsTransactionExpired bool   `json:"is_transaction_expired"`
	CreatedAt            string `json:"created_at"`
	ExpiresAt            string `json:"expires_at"`
}

func (c *Client) BillsPaymentMultiPayVoidTrx(ctx context.Context, vr BillsPaymentMultiPayVoidTrxRequest) (*BillsPaymentMultiPayVoidTrxResponseBody, error) {
	res, err := c.phService.BillsPost(ctx, c.getUrl("multipay/void-transaction"), vr)
	if err != nil {
		return nil, err
	}

	prvRes := &BillsPaymentMultiPayVoidTrxResponseBody{}
	if err := json.Unmarshal(res, prvRes); err != nil {
		return nil, err
	}
	return prvRes, nil
}

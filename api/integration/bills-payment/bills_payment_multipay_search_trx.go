package bills_payment

import (
	"context"
	"encoding/json"
)

type BillsPaymentMultiPaySearchTrxRequest struct {
	ReferenceNo string `json:"reference_no"`
}

type BillsPaymentMultiPaySearchTrxResponseBody struct {
	Code    int                                 `json:"code"`
	Message string                              `json:"message"`
	Result  BillsPaymentMultiPaySearchTrxResult `json:"result"`
}

type BillsPaymentMultiPaySearchTrxResult struct {
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

func (c *Client) BillsPaymentMultiPaySearchTrx(ctx context.Context, vr BillsPaymentMultiPaySearchTrxRequest) (*BillsPaymentMultiPaySearchTrxResponseBody, error) {
	res, err := c.phService.BillsPost(ctx, c.getUrl("multipay/search-transaction"), vr)
	if err != nil {
		return nil, err
	}

	prvRes := &BillsPaymentMultiPaySearchTrxResponseBody{}
	if err := json.Unmarshal(res, prvRes); err != nil {
		return nil, err
	}
	return prvRes, nil
}

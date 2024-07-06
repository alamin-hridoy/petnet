package bills_payment

import (
	"context"
	"encoding/json"
)

type BillsPaymentEcpayRetryRequest struct {
	ID int `json:"id"`
}

type BillsPaymentEcpayRetryResponseBody struct {
	Code    string                       `json:"code"`
	Message string                       `json:"message"`
	Result  BillsPaymentEcpayRetryResult `json:"result"`
	RemcoID int                          `json:"remco_id"`
}

type BillsPaymentEcpayRetryResult struct {
	Status          string `json:"Status"`
	Message         string `json:"Message"`
	ServiceCharge   int    `json:"ServiceCharge"`
	Timestamp       string `json:"timestamp"`
	ReferenceNumber string `json:"referenceNumber"`
}

func (c *Client) BillsPaymentEcpayRetry(ctx context.Context, vr BillsPaymentEcpayRetryRequest) (*BillsPaymentEcpayRetryResponseBody, error) {
	res, err := c.phService.BillsPost(ctx, c.getUrl("ecpay/retry"), vr)
	if err != nil {
		return nil, err
	}

	prvRes := &BillsPaymentEcpayRetryResponseBody{}
	if err := json.Unmarshal(res, prvRes); err != nil {
		return nil, err
	}
	return prvRes, nil
}

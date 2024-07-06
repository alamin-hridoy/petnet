package bills_payment

import (
	"context"
	"encoding/json"
)

type BillsPaymentMultiPayBillerInquireRequest struct {
	AccountNumber string `json:"account_number"`
	Amount        int    `json:"amount"`
	ContactNumber string `json:"contact_number"`
}

type BillsPaymentMultiPayBillerInquireResponseBody struct {
	Status int                                   `json:"status"`
	Reason string                                `json:"reason"`
	Data   BillsPaymentMultiPayBillerInquireData `json:"data"`
}

type BillsPaymentMultiPayBillerInquireData struct {
	AccountNumber string `json:"account_number"`
	Amount        int    `json:"amount"`
	Biller        string `json:"biller"`
}

func (c *Client) BillsPaymentMultiPayBillerInquire(ctx context.Context, vr BillsPaymentMultiPayBillerInquireRequest) (*BillsPaymentMultiPayBillerInquireResponseBody, error) {
	res, err := c.phService.BillsPost(ctx, c.getUrl("multipay/biller-inquire"), vr)
	if err != nil {
		return nil, err
	}

	prvRes := &BillsPaymentMultiPayBillerInquireResponseBody{}
	if err := json.Unmarshal(res, prvRes); err != nil {
		return nil, err
	}
	return prvRes, nil
}

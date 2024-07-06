package bills_payment

import (
	"context"
	"encoding/json"
)

type BillsPaymentEcpayBillerCategoryResponseBody struct {
	ID           int    `json:"id"`
	BillID       int    `json:"bill_id"`
	CategoryName string `json:"category_name"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func (c *Client) BillsPaymentEcpayBillerCategory(ctx context.Context) ([]*BillsPaymentEcpayBillerCategoryResponseBody, error) {
	res, err := c.phService.BillsGet(ctx, c.getUrl("ecpay/biller-category"))
	if err != nil {
		return nil, err
	}

	prvRes := []*BillsPaymentEcpayBillerCategoryResponseBody{}
	if err := json.Unmarshal(res, &prvRes); err != nil {
		return nil, err
	}
	return prvRes, nil
}

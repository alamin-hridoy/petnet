package bills_payment

import (
	"context"
	"encoding/json"
)

type BillsPaymentMultiPayBillerCategoryResponseBody struct {
	Code    int                                        `json:"code"`
	Message string                                     `json:"message"`
	Result  []BillsPaymentMultiPayBillerCategoryResult `json:"result"`
}

type BillsPaymentMultiPayBillerCategoryResult struct {
	ID           int    `json:"id"`
	BillID       int    `json:"bill_id"`
	CategoryName string `json:"category_name"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func (c *Client) BillsPaymentMultiPayBillerCategory(ctx context.Context) (*BillsPaymentMultiPayBillerCategoryResponseBody, error) {
	res, err := c.phService.BillsGet(ctx, c.getUrl("multipay/biller-category"))
	if err != nil {
		return nil, err
	}

	prvRes := &BillsPaymentMultiPayBillerCategoryResponseBody{}
	if err := json.Unmarshal(res, prvRes); err != nil {
		return nil, err
	}
	return prvRes, nil
}

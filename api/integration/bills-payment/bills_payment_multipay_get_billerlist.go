package bills_payment

import (
	"context"
	"encoding/json"
)

type BillsPaymentMultiPayBillerlistResponseBody struct {
	Code    int                                    `json:"code"`
	Message string                                 `json:"message"`
	Result  []BillsPaymentMultiPayBillerlistResult `json:"result"`
}

type BillsPaymentMultiPayBillerlistResult struct {
	PartnerID     int         `json:"partner_id"`
	BillerTag     string      `json:"BillerTag"`
	Description   string      `json:"Description"`
	Category      int         `json:"Category"`
	FieldList     []FieldList `json:"FieldList"`
	ServiceCharge int         `json:"ServiceCharge"`
}

type FieldList struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Label       string  `json:"label"`
	Order       int     `json:"order"`
	Rules       []Rules `json:"rules"`
	Description string  `json:"description"`
	Placeholder string  `json:"placeholder"`
}

type Rules struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Value   string `json:"value"`
	Format  string `json:"format"`
	Message string `json:"message"`
	Options string `json:"options"`
}

func (c *Client) BillsPaymentMultiPayBillerlist(ctx context.Context) (*BillsPaymentMultiPayBillerlistResponseBody, error) {
	res, err := c.phService.BillsGet(ctx, c.getUrl("multipay/biller-list"))
	if err != nil {
		return nil, err
	}

	prvRes := &BillsPaymentMultiPayBillerlistResponseBody{}
	if err := json.Unmarshal(res, prvRes); err != nil {
		return nil, err
	}
	return prvRes, nil
}

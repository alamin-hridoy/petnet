package bills_payment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type BillsPaymentMultiPayBillerByCategoryResponseBody struct {
	Code    int                                          `json:"code"`
	Message string                                       `json:"message"`
	Result  []BillsPaymentMultiPayBillerByCategoryResult `json:"result"`
}

type BillsPaymentMultiPayBillerByCategoryResult struct {
	PartnerID     int         `json:"partner_id"`
	BillerTag     string      `json:"BillerTag"`
	Description   string      `json:"Description"`
	Category      int         `json:"Category"`
	FieldList     []FieldList `json:"FieldList"`
	ServiceCharge int         `json:"ServiceCharge"`
}

func (c *Client) BillsPaymentMultiPayBillerByCategory(ctx context.Context, categoryID int) (*BillsPaymentMultiPayBillerByCategoryResponseBody, error) {
	billUrl := c.getUrl(fmt.Sprintf("multipay/biller-by-category/%d", categoryID))
	decodedUrl, _ := url.QueryUnescape(billUrl)
	res, err := c.phService.BillsGet(ctx, decodedUrl)
	if err != nil {
		return nil, err
	}

	prvRes := &BillsPaymentMultiPayBillerByCategoryResponseBody{}
	if err := json.Unmarshal(res, prvRes); err != nil {
		return nil, err
	}
	return prvRes, nil
}

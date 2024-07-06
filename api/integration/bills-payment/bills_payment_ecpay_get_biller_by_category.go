package bills_payment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type BillsPaymentEcpayBillerByCategoryResponseBody struct {
	Code    int                                       `json:"code"`
	Message string                                    `json:"message"`
	Result  []BillsPaymentEcpayBillerByCategoryResult `json:"result"`
}

type BillsPaymentEcpayBillerByCategoryResult struct {
	BillerTag         string `json:"BillerTag"`
	Description       string `json:"Description"`
	FirstField        string `json:"FirstField"`
	FirstFieldFormat  string `json:"FirstFieldFormat"`
	FirstFieldWidth   string `json:"FirstFieldWidth"`
	SecondField       string `json:"SecondField"`
	SecondFieldFormat string `json:"SecondFieldFormat"`
	SecondFieldWidth  string `json:"SecondFieldWidth"`
	ServiceCharge     int    `json:"ServiceCharge"`
}

func (c *Client) BillsPaymentEcpayBillerByCategory(ctx context.Context, categoryID int) (*BillsPaymentEcpayBillerByCategoryResponseBody, error) {
	billUrl := c.getUrl(fmt.Sprintf("ecpay/biller-by-category/%d", categoryID))
	decodedUrl, _ := url.QueryUnescape(billUrl)
	res, err := c.phService.BillsGet(ctx, decodedUrl)
	if err != nil {
		return nil, err
	}

	prvRes := &BillsPaymentEcpayBillerByCategoryResponseBody{}
	if err := json.Unmarshal(res, prvRes); err != nil {
		return nil, err
	}
	return prvRes, nil
}

package bills_payment

import (
	"context"
	"encoding/json"
)

type BillsPaymentEcpayTransactRequest struct {
	BillID                int    `json:"bill_id"`
	BillerTag             string `json:"biller_tag"`
	TrxDate               string `json:"trx_date"`
	UserID                string `json:"user_id"`
	RemoteUserID          string `json:"remote_user_id"`
	CustomerID            string `json:"customer_id"`
	LocationID            string `json:"location_id"`
	RemoteLocationID      string `json:"remote_location_id"`
	LocationName          string `json:"location_name"`
	Coy                   string `json:"coy"`
	CurrencyID            string `json:"currency_id"`
	FormType              string `json:"form_type"`
	FormNumber            string `json:"form_number"`
	AccountNumber         string `json:"account_number"`
	Identifier            string `json:"identifier"`
	Amount                int    `json:"amount"`
	ServiceCharge         int    `json:"service_charge"`
	TotalAmount           int    `json:"total_amount"`
	ClientReferenceNumber string `json:"client_reference_number"`
}

type BillsPaymentEcpayTransactResponseBody struct {
	Code    string                          `json:"code"`
	Message string                          `json:"message"`
	Result  BillsPaymentEcpayTransactResult `json:"result"`
	RemcoID int                             `json:"remco_id"`
}

type BillsPaymentEcpayTransactResult struct {
	Status          string `json:"Status"`
	Message         string `json:"Message"`
	ServiceCharge   int    `json:"ServiceCharge"`
	Timestamp       string `json:"timestamp"`
	ReferenceNumber string `json:"referenceNumber"`
}

func (c *Client) BillsPaymentEcpayTransact(ctx context.Context, vr BillsPaymentEcpayTransactRequest) (*BillsPaymentEcpayTransactResponseBody, error) {
	res, err := c.phService.BillsPost(ctx, c.getUrl("ecpay/transact"), vr)
	if err != nil {
		return nil, err
	}

	prvRes := &BillsPaymentEcpayTransactResponseBody{}
	if err := json.Unmarshal(res, prvRes); err != nil {
		return nil, err
	}
	return prvRes, nil
}

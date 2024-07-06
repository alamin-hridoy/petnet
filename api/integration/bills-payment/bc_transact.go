package bills_payment

import (
	"context"
	"encoding/json"

	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
)

type BCTransactRequest struct {
	UserID                  string      `json:"user_id"`
	CustomerID              int         `json:"customer_id"`
	LocationID              string      `json:"location_id"`
	LocationName            string      `json:"location_name"`
	Coy                     string      `json:"coy"`
	CallbackURL             string      `json:"callbackUrl"`
	BillID                  string      `json:"bill_id"`
	BillerTag               string      `json:"biller_tag"`
	BillerName              string      `json:"biller_name"`
	TrxDate                 string      `json:"trx_date"`
	Amount                  string      `json:"amount"`
	ServiceCharge           string      `json:"service_charge"`
	PartnerCharge           string      `json:"partner_charge"`
	TotalAmount             int         `json:"total_amount"`
	Identifier              string      `json:"identifier"`
	AccountNumber           string      `json:"account_number"`
	PaymentMethod           string      `json:"payment_method"`
	ClientReferenceNumber   string      `json:"client_reference_number"`
	ReferenceNumber         string      `json:"reference_number"`
	ValidationNumber        string      `json:"validation_number"`
	ReceiptValidationNumber string      `json:"receipt_validation_number"`
	TpaID                   string      `json:"tpa_id"`
	CurrencyID              string      `json:"currency_id"`
	FormType                string      `json:"form_type"`
	FormNumber              string      `json:"form_number"`
	OtherInfo               BCOtherInfo `json:"otherInfo"`
	Type                    string      `json:"type"`
}

type BCTransactResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Result  BCTransactResult `json:"result"`
	RemcoID int              `json:"remco_id"`
}

type BCTransactResult struct {
	TransactionID   string        `json:"transactionId"`
	ReferenceNumber string        `json:"referenceNumber"`
	ClientReference string        `json:"clientReference"`
	BillerReference string        `json:"billerReference"`
	PaymentMethod   string        `json:"paymentMethod"`
	Amount          string        `json:"amount"`
	OtherCharges    string        `json:"otherCharges"`
	Status          string        `json:"status"`
	Message         string        `json:"message"`
	Details         []*bp.Details `json:"details"`
	CreatedAt       string        `json:"createdAt"`
	Timestamp       string        `json:"timestamp"`
}

func (c *Client) BCTransact(ctx context.Context, req BCTransactRequest) (*BCTransactResponse, error) {
	res, err := c.phService.BillsPost(ctx, c.getUrl("bayad/bayad-center/transact"), req)
	if err != nil {
		return nil, err
	}

	bct := &BCTransactResponse{}
	if err := json.Unmarshal(res, bct); err != nil {
		return nil, err
	}
	return bct, nil
}

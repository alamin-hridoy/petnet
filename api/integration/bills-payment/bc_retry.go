package bills_payment

import (
	"context"
	"encoding/json"
)

type BCRetryRequest struct {
	Coy                     string      `json:"coy"`
	Type                    string      `json:"type"`
	Amount                  string      `json:"amount"`
	TpaID                   string      `json:"tpa_id"`
	BillID                  string      `json:"bill_id"`
	UserID                  string      `json:"user_id"`
	TrxDate                 string      `json:"trx_date"`
	FormType                string      `json:"form_type"`
	OtherInfo               BCOtherInfo `json:"otherInfo"`
	BillerTag               string      `json:"biller_tag"`
	Identifier              string      `json:"identifier"`
	BillerName              string      `json:"biller_name"`
	CurrencyID              string      `json:"currency_id"`
	CustomerID              int         `json:"customer_id"`
	FormNumber              string      `json:"form_number"`
	LocationID              string      `json:"location_id"`
	TotalAmount             int         `json:"total_amount"`
	LocationName            string      `json:"location_name"`
	AccountNumber           string      `json:"account_number"`
	PartnerCharge           string      `json:"partner_charge"`
	PaymentMethod           string      `json:"payment_method"`
	ServiceCharge           string      `json:"service_charge"`
	ReferenceNumber         string      `json:"reference_number"`
	ValidationNumber        string      `json:"validation_number"`
	ClientReferenceNumber   string      `json:"client_reference_number"`
	ReceiptValidationNumber string      `json:"receipt_validation_number"`
	ID                      int         `json:"id"`
}

type BCOtherInfo struct {
	LastName        string `json:"LastName"`
	FirstName       string `json:"FirstName"`
	MiddleName      string `json:"MiddleName"`
	Name            string `json:"Name"`
	PaymentType     string `json:"PaymentType"`
	Course          string `json:"Course"`
	TotalAssessment string `json:"TotalAssessment"`
	SchoolYear      string `json:"SchoolYear"`
	Term            string `json:"Term"`
}

type BCRetryResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Result  BCRetryResult `json:"result"`
	RemcoID int           `json:"remco_id"`
}

type BCRetryResult struct {
	TransactionID   string        `json:"transactionId"`
	ReferenceNumber string        `json:"referenceNumber"`
	ClientReference string        `json:"clientReference"`
	BillerReference string        `json:"billerReference"`
	PaymentMethod   string        `json:"paymentMethod"`
	Amount          string        `json:"amount"`
	OtherCharges    string        `json:"otherCharges"`
	Status          string        `json:"status"`
	Message         string        `json:"message"`
	Details         []interface{} `json:"details"`
	CreatedAt       string        `json:"createdAt"`
	Timestamp       string        `json:"timestamp"`
}

func (c *Client) BCRetry(ctx context.Context, req BCRetryRequest) (*BCRetryResponse, error) {
	res, err := c.phService.BillsPost(ctx, c.getUrl("bayad/bayad-center/retry"), req)
	if err != nil {
		return nil, err
	}

	bcr := &BCRetryResponse{}
	if err := json.Unmarshal(res, bcr); err != nil {
		return nil, err
	}
	return bcr, nil
}

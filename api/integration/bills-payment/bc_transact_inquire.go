package bills_payment

import (
	"context"
	"encoding/json"

	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
)

type BCTransactInquireRequest struct {
	Code            string `json:"code"`
	ClientReference string `json:"clientReference"`
}

type BCTransactInquireResponse struct {
	Code    int                     `json:"code"`
	Message string                  `json:"message"`
	Result  BCTransactInquireResult `json:"result"`
	RemcoID int                     `json:"remco_id"`
}

type BCTransactInquireResult struct {
	TransactionID   string                         `json:"transactionId"`
	ReferenceNumber string                         `json:"referenceNumber"`
	ClientReference string                         `json:"clientReference"`
	BillerReference string                         `json:"billerReference"`
	PaymentMethod   string                         `json:"paymentMethod"`
	Amount          string                         `json:"amount"`
	OtherCharges    string                         `json:"otherCharges"`
	Status          string                         `json:"status"`
	Message         BCTransactInquireResultMessage `json:"message"`
	Details         []*bp.Details                  `json:"details"`
	CreatedAt       string                         `json:"createdAt"`
}

type BCTransactInquireResultMessage struct {
	Header  string `json:"header"`
	Message string `json:"message"`
	Footer  string `json:"footer"`
}

func (c *Client) BCTransactInquire(ctx context.Context, req BCTransactInquireRequest) (*BCTransactInquireResponse, error) {
	res, err := c.phService.BillsPost(ctx, c.getUrl("bayad/bayad-center/transact-inquire"), req)
	if err != nil {
		return nil, err
	}

	bcti := &BCTransactInquireResponse{}
	if err := json.Unmarshal(res, bcti); err != nil {
		return nil, err
	}
	return bcti, nil
}

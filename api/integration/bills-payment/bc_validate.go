package bills_payment

import (
	"context"
	"encoding/json"
)

type BCValidateRequest struct {
	BillPartnerID int         `json:"bill_partner_id"`
	BillerTag     string      `json:"biller_tag"`
	Code          string      `json:"code"`
	AccountNumber string      `json:"accountNumber"`
	AccountNo     string      `json:"account_no"`
	Identifier    string      `json:"identifier"`
	PaymentMethod string      `json:"paymentMethod"`
	OtherCharges  string      `json:"otherCharges"`
	Amount        string      `json:"amount"`
	OtherInfo     BCOtherInfo `json:"otherInfo"`
}

type BCValidateResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Result  BCValidateResult `json:"result"`
	RemcoID int              `json:"remco_id"`
}

type BCValidateResult struct {
	Valid            bool          `json:"valid"`
	Code             int           `json:"code"`
	Account          string        `json:"account"`
	Details          []interface{} `json:"details"`
	ValidationNumber string        `json:"validationNumber"`
}

func (c *Client) BCValidate(ctx context.Context, req BCValidateRequest) (*BCValidateResponse, error) {
	res, err := c.phService.BillsPost(ctx, c.getUrl("bayad/bayad-center/validate"), req)
	if err != nil {
		return nil, err
	}

	bcv := &BCValidateResponse{}
	if err := json.Unmarshal(res, bcv); err != nil {
		return nil, err
	}
	return bcv, nil
}

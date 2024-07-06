package bills_payment

import (
	"context"
	"encoding/json"
)

type BCGettingOtherChargesRequest struct {
	Code       string `json:"code"`
	Amount     string `json:"amount"`
	UserID     string `json:"user_id"`
	LocationID string `json:"location_id"`
}

type BCGettingOtherChargesResponse struct {
	Code    int                         `json:"code"`
	Message string                      `json:"message"`
	Result  BCGettingOtherChargesResult `json:"result"`
	RemcoID int                         `json:"remco_id"`
}

type BCGettingOtherChargesResult struct {
	OtherCharges string `json:"otherCharges"`
}

func (c *Client) BCGettingOtherCharges(ctx context.Context, req BCGettingOtherChargesRequest) (*BCGettingOtherChargesResponse, error) {
	res, err := c.phService.BillsPost(ctx, c.getUrl("bayad/bayad-center/fees"), req)
	if err != nil {
		return nil, err
	}

	bcbi := &BCGettingOtherChargesResponse{}
	if err := json.Unmarshal(res, bcbi); err != nil {
		return nil, err
	}
	return bcbi, nil
}

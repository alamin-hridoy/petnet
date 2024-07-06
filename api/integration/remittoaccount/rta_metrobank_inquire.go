package remittoaccount

import (
	"context"
	"encoding/json"
)

type RTAMetrobankInquireRequest struct {
	ReferenceNumber string `json:"reference_number"`
	LocationID      string `json:"location_id"`
}

type RTAMetrobankInquireResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Result   string `json:"result"`
	BankCode string `json:"bank_code"`
}

func (c *Client) RTAMetrobankInquire(ctx context.Context, req RTAMetrobankInquireRequest) (*RTAMetrobankInquireResponse, error) {
	res, err := c.phService.RtaPost(ctx, c.getUrl("metrobank-rta/inquire"), req)
	if err != nil {
		return nil, err
	}

	mbi := &RTAMetrobankInquireResponse{}
	if err := json.Unmarshal(res, mbi); err != nil {
		return nil, err
	}
	return mbi, nil
}

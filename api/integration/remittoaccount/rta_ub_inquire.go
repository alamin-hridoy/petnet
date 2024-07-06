package remittoaccount

import (
	"context"
	"encoding/json"
	"time"
)

type RTAUBInquireRequest struct {
	ReferenceNumber string `json:"reference_number"`
}

type RTAUBInquireResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Result  RTAUBInquireResult `json:"result"`
	RemcoID string             `json:"remco_id"`
}

type RTAUBInquireResult struct {
	Code            string    `json:"code"`
	SenderRefID     string    `json:"senderRefId"`
	State           string    `json:"state"`
	UUID            string    `json:"uuid"`
	Description     string    `json:"description"`
	Type            string    `json:"type"`
	Amount          string    `json:"amount"`
	UbpTranID       string    `json:"ubpTranId"`
	TranRequestDate string    `json:"tranRequestDate"`
	TranFinacleDate string    `json:"tranFinacleDate"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (c *Client) RTAUBInquire(ctx context.Context, req RTAUBInquireRequest) (*RTAUBInquireResponse, error) {
	res, err := c.phService.RtaPost(ctx, c.getUrl("unionbank/inquire"), req)
	if err != nil {
		return nil, err
	}

	ubi := &RTAUBInquireResponse{}
	if err := json.Unmarshal(res, ubi); err != nil {
		return nil, err
	}
	return ubi, nil
}

package remittoaccount

import (
	"context"
	"encoding/json"
	"time"
)

type RTAUBRetryRequest struct {
	ReferenceNumber string `json:"reference_number"`
	ID              string `json:"id"`
}

type RTAUBRetryResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Result  RTAUBRetryResult `json:"result"`
	RemcoID string           `json:"remco_id"`
}

type RTAUBRetryResult struct {
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

func (c *Client) RTAUBRetry(ctx context.Context, req RTAUBRetryRequest) (*RTAUBRetryResponse, error) {
	res, err := c.phService.RtaPost(ctx, c.getUrl("unionbank/retry"), req)
	if err != nil {
		return nil, err
	}

	ubr := &RTAUBRetryResponse{}
	if err := json.Unmarshal(res, ubr); err != nil {
		return nil, err
	}
	return ubr, nil
}

package remittoaccount

import (
	"context"
	"encoding/json"
)

type RTAMetrobankRetryRequest struct {
	ReferenceNumber string `json:"reference_number"`
	ID              string `json:"id"`
	LocationID      int    `json:"location_id"`
	PrincipalAmount string `json:"principal_amount"`
	FormNumber      string `json:"form_number"`
}

type RTAMetrobankRetryResponse struct {
	Code     int                     `json:"code"`
	Message  string                  `json:"message"`
	Result   RTAMetrobankRetryResult `json:"result"`
	BankCode string                  `json:"bank_code"`
}
type RTAMetrobankRetryResult struct {
	Message string `json:"message"`
}

func (c *Client) RTAMetrobankRetry(ctx context.Context, req RTAMetrobankRetryRequest) (*RTAMetrobankRetryResponse, error) {
	res, err := c.phService.RtaPost(ctx, c.getUrl("metrobank-rta/retry"), req)
	if err != nil {
		return nil, err
	}

	mbr := &RTAMetrobankRetryResponse{}
	if err := json.Unmarshal(res, mbr); err != nil {
		return nil, err
	}
	return mbr, nil
}

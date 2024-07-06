package bills_payment

import (
	"context"
	"encoding/json"
)

type BCBillerListResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Result  []BCBillerListResult `json:"result"`
}

type BCBillerListResult struct {
	Name            string `json:"name"`
	Code            string `json:"code"`
	Description     string `json:"description"`
	Category        string `json:"category"`
	Type            string `json:"type"`
	Logo            string `json:"logo"`
	IsMultipleBills int    `json:"isMultipleBills"`
	IsCde           int    `json:"isCde"`
	IsAsync         int    `json:"isAsync"`
}

func (c *Client) BCBillerList(ctx context.Context) (*BCBillerListResponse, error) {
	res, err := c.phService.BillsGet(ctx, c.getUrl("bayad/bayad-center/billers"))
	if err != nil {
		return nil, err
	}

	bcbl := &BCBillerListResponse{}
	if err := json.Unmarshal(res, bcbl); err != nil {
		return nil, err
	}
	return bcbl, nil
}

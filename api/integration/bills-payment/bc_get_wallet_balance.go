package bills_payment

import (
	"context"
	"encoding/json"
)

type BCGetWalletBalanceResponse struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Result  BCGetWalletBalanceResult `json:"result"`
	RemcoID int                      `json:"remco_id"`
}

type BCGetWalletBalanceResult struct {
	Balance string `json:"balance"`
}

func (c *Client) BCGetWalletBalance(ctx context.Context) (*BCGetWalletBalanceResponse, error) {
	res, err := c.phService.BillsGet(ctx, c.getUrl("bayad/bayad-center/wallets"))
	if err != nil {
		return nil, err
	}

	bcgwb := &BCGetWalletBalanceResponse{}
	if err := json.Unmarshal(res, bcgwb); err != nil {
		return nil, err
	}
	return bcgwb, nil
}

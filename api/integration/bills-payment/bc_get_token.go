package bills_payment

import (
	"context"
	"encoding/json"
)

type BCGetTokenRequest struct {
	GrantType string `json:"grant_type"`
	TpaID     string `json:"tpa_id"`
	Scope     string `json:"scope"`
}

type BCGetTokenResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Result  BCGetTokenResult `json:"result"`
	RemcoID int              `json:"remco_id"`
}

type BCGetTokenResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (c *Client) BCGetToken(ctx context.Context, req BCGetTokenRequest) (*BCGetTokenResponse, error) {
	res, err := c.phService.BillsPost(ctx, c.getUrl("bayad/bayad-center/token"), req)
	if err != nil {
		return nil, err
	}

	bcgt := &BCGetTokenResponse{}
	if err := json.Unmarshal(res, bcgt); err != nil {
		return nil, err
	}
	return bcgt, nil
}

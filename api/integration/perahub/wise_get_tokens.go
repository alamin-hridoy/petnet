package perahub

import (
	"context"
	"encoding/json"
)

type WISEGetTokensReq struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type WISEGetTokensResp struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    json.Number `json:"expires_in"`
	TokenType    string      `json:"token_type"`
}

func (s *Svc) WISEGetTokens(ctx context.Context, req WISEGetTokensReq) (*WISEGetTokensResp, error) {
	res, err := s.postNonex(ctx, s.nonexURL("transferwise/oauth/token"), req)
	if err != nil {
		return nil, err
	}

	rb := &WISEGetTokensResp{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

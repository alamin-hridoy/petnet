package perahub

import (
	"context"
	"encoding/json"
)

type WISECreateUserReq struct {
	Email string `json:"email"`
}

type WISECreateUserResp struct {
	Msg   string `json:"message"`
	Error string `json:"error"`
}

func (s *Svc) WISECreateUser(ctx context.Context, req WISECreateUserReq) (*WISECreateUserResp, error) {
	res, err := s.postNonex(ctx, s.nonexURL("transferwise/accounts"), req)
	if err != nil {
		return nil, err
	}

	rb := &WISECreateUserResp{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

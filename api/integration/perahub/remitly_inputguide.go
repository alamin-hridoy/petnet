package perahub

import (
	"context"
	"encoding/json"
)

type RMIDsRespBody struct {
	Code   json.Number `json:"code"`
	Msg    string      `json:"message"`
	Result []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"result"`
	RemcoID json.Number `json:"remco_id"`
}

func (s *Svc) RMIDs(ctx context.Context) (*RMIDsRespBody, error) {
	res, err := s.getNonex(ctx, s.nonexURL("remitly/ids"))
	if err != nil {
		return nil, err
	}

	rb := &RMIDsRespBody{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

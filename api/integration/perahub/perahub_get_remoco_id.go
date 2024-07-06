package perahub

import (
	"context"
	"encoding/json"
)

type PerahubGetRemcoIDResponse struct {
	Code    json.Number               `json:"code"`
	Message string                    `json:"message"`
	Result  []PerahubGetRemcoIDResult `json:"result"`
}

type PerahubGetRemcoIDResult struct {
	ID   json.Number `json:"id"`
	Name string      `json:"non_ex_name"`
}

func (s *Svc) PerahubGetRemcoID(ctx context.Context) (*PerahubGetRemcoIDResponse, error) {
	res, err := s.getNonex(ctx, s.phTransactURL("drp/remco"))
	if err != nil {
		return nil, err
	}

	remocoRes := &PerahubGetRemcoIDResponse{}
	if err := json.Unmarshal(res, remocoRes); err != nil {
		return nil, err
	}
	return remocoRes, nil
}

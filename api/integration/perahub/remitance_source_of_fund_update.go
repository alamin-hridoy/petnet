package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type SourceOfFundUpdateReq struct {
	SourceOfFund string `json:"source_of_fund"`
}

type SourceOfFundUpdateResult struct {
	ID           int       `json:"id"`
	SourceOfFund string    `json:"source_of_fund"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type SourceOfFundUpdateRes struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Result  SourceOfFundUpdateResult `json:"result"`
}

func (s *Svc) SourceOfFundUpdate(ctx context.Context, req SourceOfFundUpdateReq, id string) (*SourceOfFundUpdateRes, error) {
	res, err := s.remitancePut(ctx, s.remitanceURL("sourcefund/"+id), req)
	if err != nil {
		return nil, err
	}

	rb := &SourceOfFundUpdateRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

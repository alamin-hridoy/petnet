package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type SourceOfFundCreateReq struct {
	SourceOfFund string `json:"source_of_fund"`
}

type SourceOfFundCreateResult struct {
	ID           int       `json:"id"`
	SourceOfFund string    `json:"source_of_fund"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SourceOfFundCreateRes struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Result  SourceOfFundCreateResult `json:"result"`
}

func (s *Svc) SourceOfFundCreate(ctx context.Context, req SourceOfFundCreateReq) (*SourceOfFundCreateRes, error) {
	res, err := s.remitancePost(ctx, s.remitanceURL("sourcefund"), req)
	if err != nil {
		return nil, err
	}

	rb := &SourceOfFundCreateRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

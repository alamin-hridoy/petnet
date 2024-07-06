package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type SourceOfFundGridResult struct {
	ID           json.Number `json:"id"`
	SourceOfFund string      `json:"source_of_fund"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	DeletedAt    time.Time   `json:"deleted_at"`
}

type SourceOfFundGridRes struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Result  []SourceOfFundGridResult `json:"result"`
}

func (s *Svc) SourceOfFundGrid(ctx context.Context) (*SourceOfFundGridRes, error) {
	res, err := s.remitanceGet(ctx, s.remitanceURL("sourcefund"))
	if err != nil {
		return nil, err
	}

	rb := &SourceOfFundGridRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

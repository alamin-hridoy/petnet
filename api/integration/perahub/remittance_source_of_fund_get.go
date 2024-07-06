package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type SourceOfFundGetResult struct {
	ID           int       `json:"id"`
	SourceOfFund string    `json:"source_of_fund"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type SourceOfFundGetRes struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Result  *SourceOfFundGetResult `json:"result"`
}

func (s *Svc) SourceOfFundGet(ctx context.Context, id string) (*SourceOfFundGetRes, error) {
	res, err := s.remitanceGet(ctx, s.remitanceURL("sourcefund/"+id))
	if err != nil {
		return nil, err
	}

	rb := &SourceOfFundGetRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type SourceOfFundDeleteResult struct {
	ID           int       `json:"id"`
	SourceOfFund string    `json:"source_of_fund"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type SourceOfFundDeleteRes struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Result  SourceOfFundDeleteResult `json:"result"`
}

func (s *Svc) SourceOfFundDelete(ctx context.Context, id string) (*SourceOfFundDeleteRes, error) {
	res, err := s.remitanceDelete(ctx, s.remitanceURL("sourcefund/"+id))
	if err != nil {
		return nil, err
	}

	rb := &SourceOfFundDeleteRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

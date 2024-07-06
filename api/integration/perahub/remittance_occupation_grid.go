package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type OccupationGridResult struct {
	ID         int32     `json:"id"`
	Occupation string    `json:"occupation"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}

type OccupationGridRes struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Result  []OccupationGridResult `json:"result"`
}

func (s *Svc) OccupationGrid(ctx context.Context) (*OccupationGridRes, error) {
	res, err := s.remitanceGet(ctx, s.remitanceURL("occupation"))
	if err != nil {
		return nil, err
	}

	rb := &OccupationGridRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

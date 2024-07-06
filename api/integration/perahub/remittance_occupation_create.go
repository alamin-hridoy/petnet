package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type OccupationCreateReq struct {
	Occupation string `json:"occupation"`
}

type OccupationCreateResult struct {
	ID         int       `json:"id"`
	Occupation string    `json:"occupation"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type OccupationCreateRes struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Result  OccupationCreateResult `json:"result"`
}

func (s *Svc) OccupationCreate(ctx context.Context, req OccupationCreateReq) (*OccupationCreateRes, error) {
	res, err := s.remitancePost(ctx, s.remitanceURL("occupation"), req)
	if err != nil {
		return nil, err
	}

	rb := &OccupationCreateRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

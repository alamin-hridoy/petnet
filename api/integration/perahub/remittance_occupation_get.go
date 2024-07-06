package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type OccupationGetResult struct {
	ID         int       `json:"id"`
	Occupation string    `json:"occupation"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}

type OccupationGetRes struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Result  *OccupationGetResult `json:"result"`
}

func (s *Svc) OccupationGet(ctx context.Context, id string) (*OccupationGetRes, error) {
	res, err := s.remitanceGet(ctx, s.remitanceURL("occupation/"+id))
	if err != nil {
		return nil, err
	}

	rb := &OccupationGetRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

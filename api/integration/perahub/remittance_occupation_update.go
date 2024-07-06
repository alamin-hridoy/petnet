package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type OccupationUpdateReq struct {
	Occupation string `json:"occupation"`
}

type OccupationUpdateResult struct {
	ID         int       `json:"id"`
	Occupation string    `json:"occupation"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}

type OccupationUpdateRes struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Result  OccupationUpdateResult `json:"result"`
}

func (s *Svc) OccupationUpdate(ctx context.Context, req OccupationUpdateReq, id string) (*OccupationUpdateRes, error) {
	res, err := s.remitancePut(ctx, s.remitanceURL("occupation/"+id), req)
	if err != nil {
		return nil, err
	}

	rb := &OccupationUpdateRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

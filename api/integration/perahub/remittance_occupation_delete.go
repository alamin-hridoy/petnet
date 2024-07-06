package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type OccupationDeleteResult struct {
	ID         int       `json:"id"`
	Occupation string    `json:"occupation"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}

type OccupationDeleteRes struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Result  OccupationDeleteResult `json:"result"`
}

func (s *Svc) OccupationDelete(ctx context.Context, id string) (*OccupationDeleteRes, error) {
	res, err := s.remitanceDelete(ctx, s.remitanceURL("occupation/"+id))
	if err != nil {
		return nil, err
	}

	rb := &OccupationDeleteRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type RelationshipGetResult struct {
	ID           int       `json:"id"`
	Relationship string    `json:"relationship"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type RelationshipGetRes struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Result  *RelationshipGetResult `json:"result"`
}

func (s *Svc) RelationshipGet(ctx context.Context, id string) (*RelationshipGetRes, error) {
	res, err := s.remitanceGet(ctx, s.remitanceURL("relationship/"+id))
	if err != nil {
		return nil, err
	}

	rb := &RelationshipGetRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

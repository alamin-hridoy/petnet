package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type RelationshipUpdateReq struct {
	Relationship string `json:"relationship"`
}

type RelationshipUpdateResult struct {
	ID           int       `json:"id"`
	Relationship string    `json:"relationship"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type RelationshipUpdateRes struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Result  RelationshipUpdateResult `json:"result"`
}

func (s *Svc) RelationshipUpdate(ctx context.Context, req RelationshipUpdateReq, id string) (*RelationshipUpdateRes, error) {
	res, err := s.remitancePut(ctx, s.remitanceURL("relationship/"+id), req)
	if err != nil {
		return nil, err
	}

	rb := &RelationshipUpdateRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

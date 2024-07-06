package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type RelationshipDeleteResult struct {
	ID           int       `json:"id"`
	Relationship string    `json:"relationship"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type RelationshipDeleteRes struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Result  RelationshipDeleteResult `json:"result"`
}

func (s *Svc) RelationshipDelete(ctx context.Context, id string) (*RelationshipDeleteRes, error) {
	res, err := s.remitanceDelete(ctx, s.remitanceURL("relationship/"+id))
	if err != nil {
		return nil, err
	}

	rb := &RelationshipDeleteRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

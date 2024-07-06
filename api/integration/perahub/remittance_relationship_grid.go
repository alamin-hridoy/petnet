package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type RemittanceRelationshipGridResult struct {
	ID           int       `json:"id"`
	Relationship string    `json:"relationship"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type RemittanceRelationshipGridRes struct {
	Code    int                                `json:"code"`
	Message string                             `json:"message"`
	Result  []RemittanceRelationshipGridResult `json:"result"`
}

func (s *Svc) RemittanceRelationshiptGrid(ctx context.Context) (*RemittanceRelationshipGridRes, error) {
	res, err := s.remitanceGet(ctx, s.remitanceURL("relationship"))
	if err != nil {
		return nil, err
	}

	pg := &RemittanceRelationshipGridRes{}
	if err := json.Unmarshal(res, pg); err != nil {
		return nil, err
	}
	return pg, nil
}

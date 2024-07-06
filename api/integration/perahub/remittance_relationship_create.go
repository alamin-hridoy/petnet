package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type RemittanceRelationshipCreateReq struct {
	Relationship string `json:"relationship"`
}

type RemittanceRelationshipCreateResult struct {
	Relationship string    `json:"relationship"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ID           int       `json:"id"`
}

type RemittanceRelationshipCreateRes struct {
	Code    int                                `json:"code"`
	Message string                             `json:"message"`
	Result  RemittanceRelationshipCreateResult `json:"result"`
}

func (s *Svc) RemittanceRelationshipCreate(ctx context.Context, req RemittanceRelationshipCreateReq) (*RemittanceRelationshipCreateRes, error) {
	res, err := s.remitancePost(ctx, s.remitanceURL("relationship"), req)
	if err != nil {
		return nil, err
	}

	rc := &RemittanceRelationshipCreateRes{}
	if err := json.Unmarshal(res, rc); err != nil {
		return nil, err
	}
	return rc, nil
}

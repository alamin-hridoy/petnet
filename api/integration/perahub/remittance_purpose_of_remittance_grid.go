package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type PurposeOfRemittanceGridResult struct {
	ID                  json.Number `json:"id"`
	PurposeOfRemittance string      `json:"purpose_of_remittance"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
	DeletedAt           time.Time   `json:"deleted_at"`
}

type PurposeOfRemittanceGridRes struct {
	Code    int                             `json:"code"`
	Message string                          `json:"message"`
	Result  []PurposeOfRemittanceGridResult `json:"result"`
}

func (s *Svc) PurposeOfRemittanceGrid(ctx context.Context) (*PurposeOfRemittanceGridRes, error) {
	res, err := s.remitanceGet(ctx, s.remitanceURL("purpose"))
	if err != nil {
		return nil, err
	}

	rb := &PurposeOfRemittanceGridRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

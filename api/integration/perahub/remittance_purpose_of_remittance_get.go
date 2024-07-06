package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type PurposeOfRemittanceGetResult struct {
	ID                  int       `json:"id"`
	PurposeOfRemittance string    `json:"purpose_of_remittance"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	DeletedAt           time.Time `json:"deleted_at"`
}

type PurposeOfRemittanceGetRes struct {
	Code    int                           `json:"code"`
	Message string                        `json:"message"`
	Result  *PurposeOfRemittanceGetResult `json:"result"`
}

func (s *Svc) PurposeOfRemittanceGet(ctx context.Context, id string) (*PurposeOfRemittanceGetRes, error) {
	res, err := s.remitanceGet(ctx, s.remitanceURL("purpose/"+id))
	if err != nil {
		return nil, err
	}

	rb := &PurposeOfRemittanceGetRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

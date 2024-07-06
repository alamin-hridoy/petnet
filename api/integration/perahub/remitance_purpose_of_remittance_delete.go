package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type PurposeOfRemittanceDeleteResult struct {
	ID                  int       `json:"id"`
	PurposeOfRemittance string    `json:"purpose_of_remittance"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	DeletedAt           time.Time `json:"deleted_at"`
}

type PurposeOfRemittanceDeleteRes struct {
	Code    int                             `json:"code"`
	Message string                          `json:"message"`
	Result  PurposeOfRemittanceDeleteResult `json:"result"`
}

func (s *Svc) PurposeOfRemittanceDelete(ctx context.Context, id string) (*PurposeOfRemittanceDeleteRes, error) {
	res, err := s.remitanceDelete(ctx, s.remitanceURL("purpose/"+id))
	if err != nil {
		return nil, err
	}

	rb := &PurposeOfRemittanceDeleteRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type PurposeOfRemittanceUpdateReq struct {
	PurposeOfRemittance string `json:"purpose_of_remittance"`
}

type PurposeOfRemittanceUpdateResult struct {
	ID                  int       `json:"id"`
	PurposeOfRemittance string    `json:"purpose_of_remittance"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	DeletedAt           time.Time `json:"deleted_at"`
}

type PurposeOfRemittanceUpdateRes struct {
	Code    int                             `json:"code"`
	Message string                          `json:"message"`
	Result  PurposeOfRemittanceUpdateResult `json:"result"`
}

func (s *Svc) PurposeOfRemittanceUpdate(ctx context.Context, req PurposeOfRemittanceUpdateReq, id string) (*PurposeOfRemittanceUpdateRes, error) {
	res, err := s.remitancePut(ctx, s.remitanceURL("purpose/"+id), req)
	if err != nil {
		return nil, err
	}

	rb := &PurposeOfRemittanceUpdateRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

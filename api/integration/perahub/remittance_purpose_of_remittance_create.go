package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type PurposeOfRemittanceCreateReq struct {
	PurposeOfRemittance string `json:"purpose_of_remittance"`
}

type PurposeOfRemittanceCreateResult struct {
	ID                  int       `json:"id"`
	PurposeOfRemittance string    `json:"purpose_of_remittance"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type PurposeOfRemittanceCreateRes struct {
	Code    int                             `json:"code"`
	Message string                          `json:"message"`
	Result  PurposeOfRemittanceCreateResult `json:"result"`
}

func (s *Svc) PurposeOfRemittanceCreate(ctx context.Context, req PurposeOfRemittanceCreateReq) (*PurposeOfRemittanceCreateRes, error) {
	res, err := s.remitancePost(ctx, s.remitanceURL("purpose"), req)
	if err != nil {
		return nil, err
	}

	rb := &PurposeOfRemittanceCreateRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

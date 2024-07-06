package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type RemittanceEmploymentGridResult struct {
	ID               int       `json:"id"`
	EmploymentNature string    `json:"employment_nature"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeletedAt        time.Time `json:"deleted_at"`
}

type RemittanceEmploymentGridRes struct {
	Code    int                              `json:"code"`
	Message string                           `json:"message"`
	Result  []RemittanceEmploymentGridResult `json:"result"`
}

func (s *Svc) RemittanceEmploymentGrid(ctx context.Context) (*RemittanceEmploymentGridRes, error) {
	res, err := s.remitanceGet(ctx, s.remitanceURL("employment"))
	if err != nil {
		return nil, err
	}

	pg := &RemittanceEmploymentGridRes{}
	if err := json.Unmarshal(res, pg); err != nil {
		return nil, err
	}
	return pg, nil
}

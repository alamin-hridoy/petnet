package perahub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type RemittanceEmploymentUpdateReq struct {
	Employment       string `json:"employment"`
	EmploymentNature string `json:"employment_nature"`
}

type RemittanceEmploymentUpdateResult struct {
	ID               int       `json:"id"`
	EmploymentNature string    `json:"employment_nature"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeletedAt        time.Time `json:"deleted_at"`
}

type RemittanceEmploymentUpdateRes struct {
	Code    int                              `json:"code"`
	Message string                           `json:"message"`
	Result  RemittanceEmploymentUpdateResult `json:"result"`
}

func (s *Svc) RemittanceEmploymentUpdate(ctx context.Context, req RemittanceEmploymentUpdateReq, id string) (*RemittanceEmploymentUpdateRes, error) {
	url := fmt.Sprintf(`employment/%s`, id)
	res, err := s.remitancePut(ctx, s.remitanceURL(url), req)
	if err != nil {
		return nil, err
	}

	rb := &RemittanceEmploymentUpdateRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

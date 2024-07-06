package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type RemittanceEmploymentCreateReq struct {
	Employment       string `json:"employment"`
	EmploymentNature string `json:"employment_nature"`
}
type RemittanceEmploymentCreateResult struct {
	ID               int       `json:"id"`
	EmploymentNature string    `json:"employment_nature"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type RemittanceEmploymentCreateResp struct {
	Code    int                              `json:"code"`
	Message string                           `json:"message"`
	Result  RemittanceEmploymentCreateResult `json:"result"`
}

func (s *Svc) RemittanceEmploymentCreate(ctx context.Context, req RemittanceEmploymentCreateReq) (*RemittanceEmploymentCreateResp, error) {
	res, err := s.remitancePost(ctx, s.remitanceURL("employment"), req)
	if err != nil {
		return nil, err
	}

	rb := &RemittanceEmploymentCreateResp{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

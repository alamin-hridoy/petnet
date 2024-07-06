package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type RemittanceEmploymentGetResult struct {
	ID               int       `json:"id"`
	EmploymentNature string    `json:"employment_nature"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeletedAt        time.Time `json:"deleted_at"`
}

type RemittanceEmploymentGetRes struct {
	Code    int                            `json:"code"`
	Message string                         `json:"message"`
	Result  *RemittanceEmploymentGetResult `json:"result"`
}

func (s *Svc) RemittanceEmploymentGet(ctx context.Context, id string) (*RemittanceEmploymentGetRes, error) {
	res, err := s.remitanceGet(ctx, s.remitanceURL("employment/"+id))
	if err != nil {
		return nil, err
	}

	response := &RemittanceEmploymentGetRes{}
	if err := json.Unmarshal(res, response); err != nil {
		return nil, err
	}
	return response, nil
}

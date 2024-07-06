package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type RemittanceEmploymentDeleteResult struct {
	ID               int       `json:"id"`
	EmploymentNature string    `json:"employment_nature"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeletedAt        time.Time `json:"deleted_at"`
}

type RemittanceEmploymentDeleteRes struct {
	Code    int                               `json:"code"`
	Message string                            `json:"message"`
	Result  *RemittanceEmploymentDeleteResult `json:"result"`
}

func (s *Svc) RemittanceEmploymentDelete(ctx context.Context, id string) (*RemittanceEmploymentDeleteRes, error) {
	res, err := s.remitanceDelete(ctx, s.remitanceURL("employment/"+id))
	if err != nil {
		return nil, err
	}

	response := &RemittanceEmploymentDeleteRes{}
	if err := json.Unmarshal(res, response); err != nil {
		return nil, err
	}
	return response, nil
}

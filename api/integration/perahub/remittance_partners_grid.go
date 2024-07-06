package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type RemittancePartnersGridResult struct {
	ID           int       `json:"id"`
	PartnerCode  string    `json:"partner_code"`
	PartnerName  string    `json:"partner_name"`
	ClientSecret string    `json:"client_secret"`
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type RemittancePartnersGridRes struct {
	Code    int                            `json:"code"`
	Message string                         `json:"message"`
	Result  []RemittancePartnersGridResult `json:"result"`
}

func (s *Svc) RemittancePartnersGrid(ctx context.Context) (*RemittancePartnersGridRes, error) {
	res, err := s.remitanceGet(ctx, s.remitanceURL("partner"))
	if err != nil {
		return nil, err
	}

	pg := &RemittancePartnersGridRes{}
	if err := json.Unmarshal(res, pg); err != nil {
		return nil, err
	}
	return pg, nil
}

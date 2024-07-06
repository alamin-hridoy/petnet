package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type RemittancePartnersUpdateReq struct {
	ID          string `json:"id"`
	PartnerCode string `json:"partner_code"`
	PartnerName string `json:"partner_name"`
	Service     string `json:"service"`
}

type RemittancePartnersUpdateResult struct {
	ID           int       `json:"id"`
	PartnerCode  string    `json:"partner_code"`
	PartnerName  string    `json:"partner_name"`
	ClientSecret string    `json:"client_secret"`
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    string    `json:"deleted_at"`
}

type RemittancePartnersUpdateRes struct {
	Code    int                            `json:"code"`
	Message string                         `json:"message"`
	Result  RemittancePartnersUpdateResult `json:"result"`
}

func (s *Svc) RemittancePartnersUpdate(ctx context.Context, req RemittancePartnersUpdateReq, id string) (*RemittancePartnersUpdateRes, error) {
	res, err := s.remitancePut(ctx, s.remitanceURL("partner/"+id), req)
	if err != nil {
		return nil, err
	}

	pu := &RemittancePartnersUpdateRes{}
	if err := json.Unmarshal(res, pu); err != nil {
		return nil, err
	}
	return pu, nil
}

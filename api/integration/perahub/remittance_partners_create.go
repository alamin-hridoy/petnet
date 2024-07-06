package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type RemittancePartnersCreateReq struct {
	PartnerCode string `json:"partner_code"`
	PartnerName string `json:"partner_name"`
	Service     string `json:"service"`
}

type RemittancePartnersCreateResult struct {
	ID           int       `json:"id"`
	PartnerCode  string    `json:"partner_code"`
	PartnerName  string    `json:"partner_name"`
	ClientSecret string    `json:"client_secret"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RemittancePartnersCreateRes struct {
	Code    int                            `json:"code"`
	Message string                         `json:"message"`
	Result  RemittancePartnersCreateResult `json:"result"`
}

func (s *Svc) RemittancePartnersCreate(ctx context.Context, req RemittancePartnersCreateReq) (*RemittancePartnersCreateRes, error) {
	res, err := s.remitancePost(ctx, s.remitanceURL("partner"), req)
	if err != nil {
		return nil, err
	}

	pc := &RemittancePartnersCreateRes{}
	if err := json.Unmarshal(res, pc); err != nil {
		return nil, err
	}
	return pc, nil
}

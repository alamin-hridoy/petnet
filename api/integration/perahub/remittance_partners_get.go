package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type RemittancePartnersGetResult struct {
	ID           int       `json:"id"`
	PartnerCode  string    `json:"partner_code"`
	PartnerName  string    `json:"partner_name"`
	ClientSecret string    `json:"client_secret"`
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    string    `json:"deleted_at"`
}

type RemittancePartnersGetRes struct {
	Code    int                          `json:"code"`
	Message string                       `json:"message"`
	Result  *RemittancePartnersGetResult `json:"result"`
}

func (s *Svc) RemittancePartnersGet(ctx context.Context, id string) (*RemittancePartnersGetRes, error) {
	res, err := s.remitanceGet(ctx, s.remitanceURL("partner/"+id))
	if err != nil {
		return nil, err
	}

	rpg := &RemittancePartnersGetRes{}
	if err := json.Unmarshal(res, rpg); err != nil {
		return nil, err
	}
	return rpg, nil
}

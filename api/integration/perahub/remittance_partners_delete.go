package perahub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type RemittancePartnersDeleteReq struct {
	ID          string `json:"id"`
	PartnerCode string `json:"partner_code"`
	PartnerName string `json:"partner_name"`
}

type RemittancePartnersDeleteResult struct {
	ID           int       `json:"id"`
	PartnerCode  string    `json:"partner_code"`
	PartnerName  string    `json:"partner_name"`
	ClientSecret string    `json:"client_secret"`
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type RemittancePartnersDeleteRes struct {
	Code    int                            `json:"code"`
	Message string                         `json:"message"`
	Result  RemittancePartnersDeleteResult `json:"result"`
}

func (s *Svc) RemittancePartnersDelete(ctx context.Context, req RemittancePartnersDeleteReq) (*RemittancePartnersDeleteRes, error) {
	remitanceURL := s.remitanceURL(fmt.Sprintf("partner/%s", req.ID))
	decodedUrl, _ := url.QueryUnescape(remitanceURL)
	res, err := s.remitanceDelete(ctx, decodedUrl)
	if err != nil {
		return nil, err
	}
	pd := &RemittancePartnersDeleteRes{}
	if err := json.Unmarshal(res, pd); err != nil {
		return nil, err
	}
	return pd, nil
}

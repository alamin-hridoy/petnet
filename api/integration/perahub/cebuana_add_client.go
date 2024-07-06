package perahub

import (
	"context"
	"encoding/json"
)

type CebAddClientReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"`
	CpCtryID  string `json:"cellphone_country_id"`
	ContactNo string `json:"contact_number"`
	TpCtryID  string `json:"telephone_country_id"`
	TpArCode  string `json:"telephone_area_code"`
	CrtyAdID  string `json:"country_address_id"`
	PAdd      string `json:"province_address"`
	CAdd      string `json:"current_address"`
	UserID    string `json:"user_id"`
	SOFID     string `json:"source_of_fund_id"`
	Tin       string `json:"tin"`
	TpNo      string `json:"telephone_number"`
	AgentCode string `json:"agent_code"`
}

type CebAddClientResp struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Result  RUResult `json:"result"`
	RemcoID int      `json:"remco_id"`
}

type RUResult struct {
	ResultStatus string `json:"ResultStatus"`
	MessageID    int    `json:"MessageID"`
	LogID        int    `json:"LogID"`
	ClientID     int    `json:"ClientID"`
	ClientNo     string `json:"ClientNumber"`
}

func (s *Svc) CebAddClient(ctx context.Context, req CebAddClientReq) (*CebAddClientResp, error) {
	res, err := s.postNonex(ctx, s.nonexURL("cebuana/add-client"), req)
	if err != nil {
		return nil, err
	}

	rb := &CebAddClientResp{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

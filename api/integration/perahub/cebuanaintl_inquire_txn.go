package perahub

import (
	"context"
	"encoding/json"
)

const (
	CEBINTAwaitPayment = "Successful"
)

type CEBINTInquireRequest struct {
	ControlNumber            string      `json:"control_number"`
	LocationID               json.Number `json:"location_id"`
	UserID                   json.Number `json:"user_id"`
	LocationName             string      `json:"location_name"`
	InternationalPartnerCode string      `json:"international_partner_code"`
	DeviceID                 string      `json:"device_id"`
	AgentID                  string      `json:"agent_id"`
	AgentCode                string      `json:"agent_code"`
	BranchCode               string      `json:"branch_code"`
	LocationCode             string      `json:"location_code"`
	Branch                   string      `json:"branch"`
	OutletCode               string      `json:"outlet_code"`
	ReferenceNumber          string      `json:"reference_number"`
}

type CEBINTInquireResponseBody struct {
	Code    json.Number         `json:"code"`
	Message string              `json:"message"`
	Result  CEBINTInquireResult `json:"result"`
	RemcoID json.Number         `json:"remco_id"`
}

type CEBINTInquireResult struct {
	IsDomestic                  json.Number `json:"is_domestic"`
	ResultStatus                string      `json:"result_status"`
	MessageID                   json.Number `json:"message_id"`
	LogID                       json.Number `json:"log_id"`
	ClientReferenceNo           json.Number `json:"client_reference_no"`
	ControlNumber               string      `json:"control_number"`
	SenderName                  string      `json:"sender_name"`
	ReceiverName                string      `json:"receiver_name"`
	PrincipalAmount             json.Number `json:"principal_amount"`
	ServiceCharge               json.Number `json:"service_charge"`
	BirthDate                   string      `json:"birth_date"`
	Currency                    string      `json:"currency"`
	BeneficiaryID               json.Number `json:"beneficiary_id"`
	RemittanceStatusID          json.Number `json:"remittance_status_id"`
	RemittanceStatusDescription string      `json:"remittance_status_description"`
}

func (s *Svc) CEBINTInquire(ctx context.Context, sr CEBINTInquireRequest) (*CEBINTInquireResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("cebuana-international/inquire"), sr)
	if err != nil {
		return nil, err
	}

	cebRes := &CEBINTInquireResponseBody{}
	if err := json.Unmarshal(res, cebRes); err != nil {
		return nil, err
	}
	return cebRes, nil
}

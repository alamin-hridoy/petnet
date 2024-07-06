package perahub

import (
	"context"
	"encoding/json"
)

const (
	CEBAwaitPayment = "Successful"
)

type CEBInquireRequest struct {
	ControlNo    string      `json:"control_number"`
	LocationID   json.Number `json:"location_id"`
	UserID       json.Number `json:"user_id"`
	LocationName string      `json:"location_name"`
}

type CEBInquireResponseBody struct {
	Code    json.Number      `json:"code"`
	Message string           `json:"message"`
	Result  CEBInquireResult `json:"result"`
	RemcoID json.Number      `json:"remco_id"`
}

type CEBInquireResult struct {
	ResultStatus      string      `json:"result_status"`
	MessageID         json.Number `json:"message_id"`
	LogID             json.Number `json:"log_id"`
	ClientReferenceNo json.Number `json:"client_reference_no"`
	ControlNo         string      `json:"control_number"`
	SenderName        string      `json:"sender_name"`
	RcvName           string      `json:"receiver_name"`
	PnplAmt           json.Number `json:"principal_amount"`
	ServiceCharge     json.Number `json:"service_charge"`
	BirthDate         string      `json:"birth_date"`
	Currency          string      `json:"currency"`
	BeneficiaryID     json.Number `json:"beneficiary_id"`
	RemStatusID       json.Number `json:"remittance_status_id"`
	RemStatusDes      string      `json:"remittance_status_description"`
}

func (s *Svc) CEBInquire(ctx context.Context, sr CEBInquireRequest) (*CEBInquireResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("cebuana/inquire"), sr)
	if err != nil {
		return nil, err
	}

	cebRes := &CEBInquireResponseBody{}
	if err := json.Unmarshal(res, cebRes); err != nil {
		return nil, err
	}
	return cebRes, nil
}

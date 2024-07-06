package perahub

import (
	"context"
	"encoding/json"
)

const (
	USSCAwaitPayment = "OK"
)

type USSCInquireRequest struct {
	RefNo        string      `json:"reference_number"`
	ControlNo    string      `json:"control_number"`
	LocationID   json.Number `json:"location_id"`
	UserID       json.Number `json:"user_id"`
	LocationName string      `json:"location_name"`
	BranchCode   string      `json:"branch_code"`
}

type USSCInquireResponseBody struct {
	Code    string        `json:"code"`
	Msg     string        `json:"message"`
	Result  USSCInqResult `json:"result"`
	RemcoID json.Number   `json:"remco_id"`
}

type USSCInqResult struct {
	RcvName            string      `json:"receiver_name"`
	ControlNo          string      `json:"control_number"`
	PrincipalAmount    string      `json:"principal_amount"`
	ContactNumber      string      `json:"contact_number"`
	RefNo              string      `json:"reference_number"`
	SenderName         string      `json:"sender_name"`
	TrxDate            json.Number `json:"trx_date"`
	SenderLastName     string      `json:"sender_last_name"`
	SenderFirstName    string      `json:"sender_first_name"`
	SenderMiddleName   string      `json:"sender_middle_name"`
	ServiceCharge      string      `json:"service_charge"`
	TotalAmount        string      `json:"total_amount"`
	ReceiverFirstName  string      `json:"receiver_first_name"`
	ReceiverMiddleName string      `json:"receiver_middle_name"`
	ReceiverLastName   string      `json:"receiver_last_name"`
	RelationTo         string      `json:"relation_to"`
	PurposeTransaction string      `json:"purpose_transaction"`
}

func (s *Svc) USSCInquire(ctx context.Context, sr USSCInquireRequest) (*USSCInquireResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("ussc/inquire"), sr)
	if err != nil {
		return nil, err
	}
	rb := &USSCInquireResponseBody{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

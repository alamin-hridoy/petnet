package perahub

import (
	"context"
	"encoding/json"
)

const (
	MBAwaitPayment = "Available for pick-up"
)

type MBInquireRequest struct {
	RefNo        string      `json:"reference_number"`
	ControlNo    string      `json:"control_number"`
	LocationID   json.Number `json:"location_id"`
	UserID       json.Number `json:"user_id"`
	LocationName string      `json:"location_name"`
}

type MBInquireResponseBody struct {
	Code    json.Number `json:"code"`
	Msg     string      `json:"message"`
	Result  MBInqResult `json:"result"`
	RemcoID json.Number `json:"remco_id"`
}

type MBInqResult struct {
	RefNo           string      `json:"reference_number"`
	ControlNo       string      `json:"control_number"`
	StatusText      string      `json:"status_text"`
	PrincipalAmount json.Number `json:"principal_amount"`
	RcvName         string      `json:"receiver_name"`
	Address         string      `json:"address"`
	ContactNumber   string      `json:"contact_number"`
	Currency        string      `json:"currency"`
}

func (s *Svc) MBInquire(ctx context.Context, sr MBInquireRequest) (*MBInquireResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("metrobank/inquire"), sr)
	if err != nil {
		return nil, err
	}

	rb := &MBInquireResponseBody{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

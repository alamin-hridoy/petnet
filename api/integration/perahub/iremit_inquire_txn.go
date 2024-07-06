package perahub

import (
	"context"
	"encoding/json"
)

const (
	IRAwaitPayment = "Available for Pick-up"
)

type IRInquireRequest struct {
	Branch       string      `json:"branch"`
	RefNo        string      `json:"reference_number"`
	ControlNo    string      `json:"control_number"`
	LocationID   json.Number `json:"location_id"`
	UserID       json.Number `json:"user_id"`
	LocationName string      `json:"location_name"`
}

type IRInquireResponseBody struct {
	Code    json.Number `json:"code"`
	Msg     string      `json:"message"`
	Result  IRResult    `json:"result"`
	RemcoID json.Number `json:"remco_id"`
}

type IRResult struct {
	Status        json.Number `json:"status"`
	Desc          string      `json:"desc"`
	ControlNo     string      `json:"control_number"`
	RefNo         string      `json:"reference_number"`
	PnplAmt       json.Number `json:"principal_amount"`
	SenderName    string      `json:"sender_name"`
	RcvName       string      `json:"receiver_name"`
	Address       StringArray `json:"address"`
	CurrencyCode  string      `json:"currency"`
	ContactNumber StringArray `json:"contact_number"`
	RcvLastName   string      `json:"receiver_last_name"`
	RcvFirstName  string      `json:"receiver_first_name"`
	TxnDate       string      `json:"transaction_date"`
}

type StringArray string

func (w *StringArray) UnmarshalJSON(data []byte) (err error) {
	s := string(data)
	if s == "[]" {
		*w = ""
		return
	}
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	*w = StringArray(s)
	return nil
}

func (s *Svc) IRemitInquire(ctx context.Context, sr IRInquireRequest) (*IRInquireResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("iremit/inquire"), sr)
	if err != nil {
		return nil, err
	}

	irRes := &IRInquireResponseBody{}
	if err := json.Unmarshal(res, irRes); err != nil {
		return nil, err
	}
	return irRes, nil
}

package perahub

import (
	"context"
	"encoding/json"
)

const (
	TFAwaitPayment = "T"
)

type TFInquireRequest struct {
	Branch       string      `json:"branch"`
	RefNo        string      `json:"reference_number"`
	ControlNo    string      `json:"control_number"`
	LocationID   json.Number `json:"location_id"`
	UserID       json.Number `json:"user_id"`
	LocationName string      `json:"location_name"`
}

type TFInquireResponseBody struct {
	Code    json.Number `json:"code"`
	Msg     string      `json:"message"`
	Result  TFResult    `json:"result"`
	RemcoID json.Number `json:"remco_id"`
}

type TFResult struct {
	Status        string      `json:"Status"`
	Desc          string      `json:"Desc"`
	ControlNo     string      `json:"control_number"`
	RefNo         json.Number `json:"reference_number"`
	PnplAmt       json.Number `json:"principal_amount"`
	SenderName    string      `json:"sender_name"`
	RcvName       string      `json:"receiver_name"`
	Address       string      `json:"address"`
	CurrencyCode  string      `json:"currency"`
	ContactNumber string      `json:"contact_number"`
	RcvLastName   string      `json:"receiver_last_name"`
	RcvFirstName  string      `json:"receiver_first_name"`
	OrgnCtry      string      `json:"originating_country"`
	DestCtry      string      `json:"destination_country"`
	TxnDate       string      `json:"transaction_date"`
	IsDomestic    json.Number `json:"is_domestic"`
	IDType        json.Number `json:"id_type"`
	RcvCtryCode   string      `json:"receiver_country_iso_code"`
	RcvStateID    string      `json:"receiver_state_id"`
	RcvStateName  string      `json:"receiver_state_name"`
	RcvCityID     json.Number `json:"receiver_city_id"`
	RcvCityName   string      `json:"receiver_city_name"`
	RcvIDType     json.Number `json:"receiver_id_type"`
	RcvIsIndiv    string      `json:"receiver_is_individual"`
	PrpsOfRmtID   json.Number `json:"purpose_of_remittance_id"`
}

func (s *Svc) TFInquire(ctx context.Context, sr TFInquireRequest) (*TFInquireResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("transfast/inquire"), sr)
	if err != nil {
		return nil, err
	}

	tfRes := &TFInquireResponseBody{}
	if err := json.Unmarshal(res, tfRes); err != nil {
		return nil, err
	}
	return tfRes, nil
}

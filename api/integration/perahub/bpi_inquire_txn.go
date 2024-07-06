package perahub

import (
	"context"
	"encoding/json"
)

const (
	BPAwaitPayment  = "IN PROCESS:TRANSACTION PROCESS ONGOING"
	BPAwaitPaymentD = "Ready for payout"
)

type BPInquireRequest struct {
	RefNo        string      `json:"reference_number"`
	ControlNo    string      `json:"control_number"`
	LocationID   json.Number `json:"location_id"`
	UserID       json.Number `json:"user_id"`
	TrxDate      string      `json:"trx_date"`
	LocationName string      `json:"location_name"`
}

type BPInquireResponseBody struct {
	Code    json.Number     `json:"code"`
	Msg     string          `json:"message"`
	Result  BPInquireResult `json:"result"`
	RemcoID json.Number     `json:"remco_id"`
}

type BPInquireResult struct {
	Status            string `json:"Status"`
	Desc              string `json:"Desc"`
	ControlNo         string `json:"control_number"`
	RefNo             string `json:"reference_number"`
	ClientReferenceNo string `json:"client_reference_no"`
	PnplAmt           string `json:"principal_amount"`
	SenderName        string `json:"sender_name"`
	RcvName           string `json:"receiver_name"`
	Address           string `json:"address"`
	Currency          string `json:"currency"`
	ContactNumber     string `json:"contact_number"`
	OrgnCtry          string `json:"originating_country"`
	DestCtry          string `json:"destination_country"`
}

func (s *Svc) BPInquire(ctx context.Context, sr BPInquireRequest) (*BPInquireResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("bpi/inquire"), sr)
	if err != nil {
		return nil, err
	}

	bpRes := &BPInquireResponseBody{}
	if err := json.Unmarshal(res, bpRes); err != nil {
		return nil, err
	}
	return bpRes, nil
}

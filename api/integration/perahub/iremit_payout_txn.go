package perahub

import (
	"context"
	"encoding/json"
)

type IRPayoutRequest struct {
	LocationID    json.Number `json:"location_id"`
	LocationName  string      `json:"location_name"`
	UserID        json.Number `json:"user_id"`
	TxnDate       string      `json:"trx_date"`
	CustomerID    json.Number `json:"customer_id"`
	CurrencyID    json.Number `json:"currency_id"`
	RemcoID       json.Number `json:"remco_id"`
	TxnType       json.Number `json:"trx_type"`
	IsDomestic    json.Number `json:"is_domestic"`
	CustomerName  string      `json:"customer_name"`
	ServiceCharge json.Number `json:"service_charge"`
	RmtLocID      json.Number `json:"remote_location_id"`
	DstAmount     json.Number `json:"dst_amount"`
	TotalAmount   string      `json:"total_amount"`
	RmtUserID     json.Number `json:"remote_user_id"`
	RmtIPAddr     string      `json:"remote_ip_address"`
	PurposeTxn    string      `json:"purpose_transaction"`
	SourceFund    string      `json:"source_fund"`
	Occupation    string      `json:"occupation"`
	RelationTo    string      `json:"relation_to"`
	BirthDate     string      `json:"birth_date"`
	BirthPlace    string      `json:"birth_place"`
	BirthCountry  string      `json:"birth_country"`
	IDType        string      `json:"id_type"`
	IDNumber      string      `json:"id_number"`
	Address       string      `json:"address"`
	Barangay      string      `json:"barangay"`
	City          string      `json:"city"`
	Province      string      `json:"province"`
	ZipCode       string      `json:"zip_code"`
	Country       string      `json:"country"`
	ContactNumber string      `json:"contact_number"`
	RiskScore     json.Number `json:"risk_score"`
	RiskCriteria  json.Number `json:"risk_criteria"`
	PayoutType    json.Number `json:"payout_type"`
	SenderName    string      `json:"sender_name"`
	RcvName       string      `json:"receiver_name"`
	PnplAmt       json.Number `json:"principal_amount"`
	ControlNo     string      `json:"control_number"`
	RefNo         string      `json:"reference_number"`

	// new fields
	MCRate      string `json:"mc_rate"`
	BBAmt       string `json:"buy_back_amount"`
	RateCat     string `json:"rate_category"`
	MCRateID    string `json:"mc_rate_id"`
	OrgnCtry    string `json:"originating_country"`
	DestCtry    string `json:"destination_country"`
	ClientRefNo string `json:"client_reference_no"`
	FormType    string `json:"form_type"`
	FormNumber  string `json:"form_number"`
	DsaCode     string `json:"dsa_code"`
	DsaTrxType  string `json:"dsa_trx_type"`
}

type IRPayoutResponseBody struct {
	Code    json.Number `json:"code"`
	Msg     string      `json:"message"`
	RemcoID json.Number `json:"remco_id"`
}

func (s *Svc) IRemitPayout(ctx context.Context, sr IRPayoutRequest) (*IRPayoutResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("iremit/payout"), sr)
	if err != nil {
		return nil, err
	}

	irRes := &IRPayoutResponseBody{}
	if err := json.Unmarshal(res, irRes); err != nil {
		return nil, err
	}
	return irRes, nil
}

package perahub

import (
	"context"
	"encoding/json"
)

type BPPayoutRequest struct {
	BBAmt         json.Number `json:"buy_back_amount"`
	MCRate        string      `json:"mc_rate"`
	RateCat       json.Number `json:"rate_category"`
	MCRateID      json.Number `json:"mc_rate_id"`
	ControlNo     string      `json:"control_number"`
	LocationID    string      `json:"location_id"`
	UserID        string      `json:"user_id"`
	TxnDate       string      `json:"trx_date"`
	LocationName  string      `json:"location_name"`
	CustomerID    string      `json:"customer_id"`
	CurrencyID    string      `json:"currency_id"`
	RemcoID       string      `json:"remco_id"`
	TxnType       string      `json:"trx_type"`
	IsDomestic    string      `json:"is_domestic"`
	CustomerName  string      `json:"customer_name"`
	ServiceCharge string      `json:"service_charge"`
	RmtLocID      string      `json:"remote_location_id"`
	DstAmount     string      `json:"dst_amount"`
	TotalAmount   string      `json:"total_amount"`
	RmtUserID     string      `json:"remote_user_id"`
	OrgnCtry      string      `json:"originating_country"`
	DestCtry      string      `json:"destination_country"`
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
	CurAdd        string      `json:"current_address"`
	PerAdd        string      `json:"permanent_address"`
	RiskScore     json.Number `json:"risk_score"`
	RiskCriteria  string      `json:"risk_criteria"`
	FormType      string      `json:"form_type"`
	FormNumber    string      `json:"form_number"`
	PayoutType    string      `json:"payout_type"`
	SenderName    string      `json:"sender_name"`
	RcvName       string      `json:"receiver_name"`
	PnplAmt       string      `json:"principal_amount"`
	RefNo         string      `json:"reference_number"`
	ClientRefNo   string      `json:"client_reference_no"`
	DsaCode       string      `json:"dsa_code"`
	DsaTrxType    string      `json:"dsa_trx_type"`
}

type BPPayoutResponseBody struct {
	Code    json.Number    `json:"code"`
	Msg     string         `json:"message"`
	Result  BPPayoutResult `json:"result"`
	RemcoID json.Number    `json:"remco_id"`
}

type BPPayoutResult struct {
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
	RcvLastName       string `json:"receiver_last_name"`
	RcvFirstName      string `json:"receiver_first_name"`
	OrgnCtry          string `json:"originating_country"`
	DestCtry          string `json:"destination_country"`
	TxnDate           string `json:"transaction_date"`
	IsDomestic        string `json:"is_domestic"`
	IDType            string `json:"id_type"`
	RcvCtryCode       string `json:"receiver_country_iso_code"`
	RcvStateID        string `json:"receiver_state_id"`
	RcvStateName      string `json:"receiver_state_name"`
	RcvCityID         string `json:"receiver_city_id"`
	RcvCityName       string `json:"receiver_city_name"`
	RcvIDType         string `json:"receiver_id_type"`
	RcvIsIndiv        string `json:"receiver_is_individual"`
	PrpsOfRmtID       string `json:"purpose_of_remittance_id"`
	DsaCode           string `json:"dsa_code"`
	DsaTrxType        string `json:"dsa_trx_type"`
}

func (s *Svc) BPPayout(ctx context.Context, sr BPPayoutRequest) (*BPPayoutResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("bpi/payout"), sr)
	if err != nil {
		return nil, err
	}

	bpRes := &BPPayoutResponseBody{}
	if err := json.Unmarshal(res, bpRes); err != nil {
		return nil, err
	}
	return bpRes, nil
}

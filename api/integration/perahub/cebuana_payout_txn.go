package perahub

import (
	"context"
	"encoding/json"
)

type CEBPayoutRequest struct {
	LocationID       json.Number  `json:"location_id"`
	LocationName     string       `json:"location_name"`
	UserID           json.Number  `json:"user_id"`
	TxnDate          string       `json:"trx_date"`
	CustomerID       string       `json:"customer_id"`
	CurrencyID       json.Number  `json:"currency_id"`
	BeneficiaryID    json.Number  `json:"beneficiary_id"`
	IDTypeID         json.Number  `json:"identification_type_id"`
	RemcoID          json.Number  `json:"remco_id"`
	TxnType          string       `json:"trx_type"`
	IsDomestic       string       `json:"is_domestic"`
	CustomerName     string       `json:"customer_name"`
	ServiceCharge    json.Number  `json:"service_charge"`
	RmtLocID         json.Number  `json:"remote_location_id"`
	DstAmount        string       `json:"dst_amount"`
	TotalAmount      string       `json:"total_amount"`
	RmtUserID        json.Number  `json:"remote_user_id"`
	RmtIpADD         string       `json:"remote_ip_address"`
	OrgnCtry         string       `json:"originating_country"`
	DestCtry         string       `json:"destination_country"`
	PurposeTxn       string       `json:"purpose_transaction"`
	SourceFund       string       `json:"source_fund"`
	Occupation       string       `json:"occupation"`
	RelationTo       string       `json:"relation_to"`
	BirthDate        string       `json:"birth_date"`
	BirthPlace       string       `json:"birth_place"`
	BirthCountry     string       `json:"birth_country"`
	IDType           string       `json:"id_type"`
	IDNumber         string       `json:"id_number"`
	Address          string       `json:"address"`
	Barangay         string       `json:"barangay"`
	City             string       `json:"city"`
	Province         string       `json:"province"`
	ZipCode          string       `json:"zip_code"`
	Country          string       `json:"country"`
	ContactNumber    string       `json:"contact_number"`
	CurrentAddress   NonexAddress `json:"current_address"`
	PermanentAddress NonexAddress `json:"permanent_address"`
	RiskScore        json.Number  `json:"risk_score"`
	PayoutType       json.Number  `json:"payout_type"`
	SenderName       string       `json:"sender_name"`
	RcvName          string       `json:"receiver_name"`
	PnplAmt          json.Number  `json:"principal_amount"`
	ClientRefNo      json.Number  `json:"client_reference_no"`
	ControlNo        string       `json:"control_number"`
	RefNo            string       `json:"reference_number"`
	BuyBackAmount    string       `json:"buy_back_amount"`
	MCRate           string       `json:"mc_rate"`
	RateCat          string       `json:"rate_category"`
	MCRateID         string       `json:"mc_rate_id"`
	FormType         string       `json:"form_type"`
	FormNumber       string       `json:"form_number"`
	DsaCode          string       `json:"dsa_code"`
	DsaTrxType       string       `json:"dsa_trx_type"`
}

type CEBPayoutResponseBody struct {
	Code    json.Number     `json:"code"`
	Message string          `json:"message"`
	Result  CEBPayoutResult `json:"result"`
	RemcoID json.Number     `json:"remco_id"`
}

type CEBPayoutResult struct {
	Message string `json:"message"`
}

func (s *Svc) CEBPayout(ctx context.Context, sr CEBPayoutRequest) (*CEBPayoutResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("cebuana/payout"), sr)
	if err != nil {
		return nil, err
	}

	cebRes := &CEBPayoutResponseBody{}
	if err := json.Unmarshal(res, cebRes); err != nil {
		return nil, err
	}
	return cebRes, nil
}

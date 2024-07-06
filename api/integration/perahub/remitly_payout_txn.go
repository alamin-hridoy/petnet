package perahub

import (
	"context"
	"encoding/json"
)

type RMPayoutRequest struct {
	ClientRefNo   string      `json:"client_reference_no"`
	RefNo         string      `json:"reference_number"`
	LocationCode  string      `json:"location_code"`
	LocationID    string      `json:"location_id"`
	LocationName  string      `json:"location_name"`
	Gender        string      `json:"gender"`
	ControlNo     string      `json:"control_number"`
	CurrencyCode  string      `json:"currency"`
	PnplAmt       string      `json:"principal_amount"`
	IDNumber      string      `json:"id_number"`
	IDType        string      `json:"id_type"`
	IDIssBy       string      `json:"id_issued_by"`
	IDIssueDate   string      `json:"id_date_of_issue"`
	IDExpDate     string      `json:"id_expiration_date"`
	ContactNumber string      `json:"contact_number"`
	Address       string      `json:"address"`
	City          string      `json:"city"`
	Province      string      `json:"province"`
	Country       string      `json:"country"`
	ZipCode       string      `json:"zip_code"`
	State         string      `json:"state"`
	Natl          string      `json:"nationality"`
	BirthDate     string      `json:"birth_date"`
	BirthCountry  string      `json:"birth_country"`
	Occupation    string      `json:"occupation"`
	UserID        string      `json:"user_id"`
	TxnDate       string      `json:"trx_date"`
	CustomerID    json.Number `json:"customer_id"`
	CurrencyID    json.Number `json:"currency_id"`
	RemcoID       json.Number `json:"remco_id"`
	TxnType       json.Number `json:"trx_type"`
	IsDomestic    json.Number `json:"is_domestic"`
	CustomerName  string      `json:"customer_name"`
	ServiceCharge string      `json:"service_charge"`
	RmtLocID      json.Number `json:"remote_location_id"`
	DstAmount     json.Number `json:"dst_amount"`
	TotalAmount   string      `json:"total_amount"`
	BBAmt         string      `json:"buy_back_amount"`
	MCRateID      string      `json:"mc_rate_id"`
	MCRate        string      `json:"mc_rate"`
	RmtIPAddr     string      `json:"remote_ip_address"`
	RmtUserID     string      `json:"remote_user_id"`
	OrgnCtry      string      `json:"originating_country"`
	DestCtry      string      `json:"destination_country"`
	PurposeTxn    string      `json:"purpose_transaction"`
	SourceFund    string      `json:"source_fund"`
	RelationTo    string      `json:"relation_to"`
	BirthPlace    string      `json:"birth_place"`
	Barangay      string      `json:"barangay"`
	RiskScore     string      `json:"risk_score"`
	RiskCriteria  string      `json:"risk_criteria"`
	FormType      string      `json:"form_type"`
	FormNumber    string      `json:"form_number"`
	PayoutType    string      `json:"payout_type"`
	SenderName    string      `json:"sender_name"`
	RcvName       string      `json:"receiver_name"`
	SenderFName   string      `json:"sender_first_name"`
	SenderMName   string      `json:"sender_middle_name"`
	SenderLName   string      `json:"sender_last_name"`
	RcvFName      string      `json:"receiver_first_name"`
	RcvMName      string      `json:"receiver_middle_name"`
	RcvLName      string      `json:"receiver_last_name"`
	AgentID       string      `json:"agent_id"`
	AgentCode     string      `json:"agent_code"`
	OrderNumber   string      `json:"order_number"`
	IPAddr        string      `json:"ip_address"`
	RateCat       string      `json:"rate_category"`
	DsaCode       string      `json:"dsa_code"`
	DsaTrxType    string      `json:"dsa_trx_type"`
}

type RMPayoutResponseBody struct {
	Code    json.Number `json:"code"`
	Msg     string      `json:"message"`
	Result  RMPayResult `json:"result"`
	RemcoID json.Number `json:"remco_id"`
}

type RMPayResult struct {
	RefNo      string `json:"reference_number"`
	Created    string `json:"created_on"`
	State      string `json:"state"`
	Type       string `json:"type"`
	PayerCodes string `json:"payer_codes"`
}

func (s *Svc) RMPayout(ctx context.Context, sr RMPayoutRequest) (*RMPayoutResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("remitly/payout"), sr)
	if err != nil {
		return nil, err
	}

	rb := &RMPayoutResponseBody{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

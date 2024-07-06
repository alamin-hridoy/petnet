package perahub

import (
	"context"
	"encoding/json"
)

type TFPayoutRequest struct {
	LocationID      json.Number `json:"location_id"`
	LocationName    string      `json:"location_name"`
	UserID          json.Number `json:"user_id"`
	TxnDate         string      `json:"trx_date"`
	CustomerID      json.Number `json:"customer_id"`
	CurrencyID      json.Number `json:"currency_id"`
	RemcoID         json.Number `json:"remco_id"`
	TxnType         json.Number `json:"trx_type"`
	IsDomestic      json.Number `json:"is_domestic"`
	CustomerName    string      `json:"customer_name"`
	ServiceCharge   json.Number `json:"service_charge"`
	RmtLocID        json.Number `json:"remote_location_id"`
	DstAmount       json.Number `json:"dst_amount"`
	TotalAmount     string      `json:"total_amount"`
	RmtUserID       json.Number `json:"remote_user_id"`
	RmtIPAddr       string      `json:"remote_ip_address"`
	PurposeTxn      string      `json:"purpose_transaction"`
	SourceFund      string      `json:"source_fund"`
	Occupation      string      `json:"occupation"`
	RelationTo      string      `json:"relation_to"`
	BirthDate       string      `json:"birth_date"`
	BirthPlace      string      `json:"birth_place"`
	BirthCountry    string      `json:"birth_country"`
	IDType          string      `json:"id_type"`
	IDNumber        string      `json:"id_number"`
	Address         string      `json:"address"`
	Barangay        string      `json:"barangay"`
	City            string      `json:"city"`
	Province        string      `json:"province"`
	ZipCode         string      `json:"zip_code"`
	Country         string      `json:"country"`
	ContactNumber   string      `json:"contact_number"`
	RiskScore       json.Number `json:"risk_score"`
	RiskCriteria    json.Number `json:"risk_criteria"`
	PayoutType      json.Number `json:"payout_type"`
	SenderName      string      `json:"sender_name"`
	RcvName         string      `json:"receiver_name"`
	PnplAmt         json.Number `json:"principal_amount"`
	ControlNo       string      `json:"control_number"`
	RefNo           string      `json:"reference_number"`
	ClientRefNo     string      `json:"client_reference_no"`
	RcvOccupationID json.Number `json:"receiver_occupation_id"`
	RcvStateID      string      `json:"receiver_state_id"`
	RcvCityID       json.Number `json:"receiver_city_id"`
	KYCVerified     bool        `json:"kyc_verified"`
	IDExpDate       string      `json:"id_expiration_date"`
	Gender          string      `json:"gender"`

	// new fields, clarify
	RmtReasonID    json.Number `json:"remittance_reason_id"`
	ProofOfAddress string      `json:"proof_of_address_collected"`
	MCRate         string      `json:"mc_rate"`
	BBAmt          string      `json:"buy_back_amount"`
	RateCat        string      `json:"rate_category"`
	MCRateID       string      `json:"mc_rate_id"`
	OrgnCtry       string      `json:"originating_country"`
	DestCtry       string      `json:"destination_country"`
	FormType       string      `json:"form_type"`
	FormNumber     string      `json:"form_number"`
	IDIssueDate    string      `json:"id_date_of_issue"`
	IDIssBy        string      `json:"id_issued_by"`
	DsaCode        string      `json:"dsa_code"`
	DsaTrxType     string      `json:"dsa_trx_type"`
}

type TFPayoutResponseBody struct {
	Code    json.Number `json:"code"`
	Msg     string      `json:"message"`
	RemcoID json.Number `json:"remco_id"`
}

func (s *Svc) TFPayout(ctx context.Context, sr TFPayoutRequest) (*TFPayoutResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("transfast/payout"), sr)
	if err != nil {
		return nil, err
	}

	irRes := &TFPayoutResponseBody{}
	if err := json.Unmarshal(res, irRes); err != nil {
		return nil, err
	}
	return irRes, nil
}

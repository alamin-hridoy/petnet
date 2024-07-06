package perahub

import (
	"context"
	"encoding/json"
)

type InstaCashPayoutRequest struct {
	LocationID         json.Number  `json:"location_id"`
	UserID             int          `json:"user_id"`
	TrxDate            string       `json:"trx_date"`
	CurrencyID         string       `json:"currency_id"`
	RemcoID            json.Number  `json:"remco_id"`
	TrxType            json.Number  `json:"trx_type"`
	IsDomestic         json.Number  `json:"is_domestic"`
	CustomerID         string       `json:"customer_id"`
	CustomerName       string       `json:"customer_name"`
	ControlNumber      string       `json:"control_number"`
	SenderName         string       `json:"sender_name"`
	ReceiverName       string       `json:"receiver_name"`
	PrincipalAmount    string       `json:"principal_amount"`
	ServiceCharge      string       `json:"service_charge"`
	DstAmount          string       `json:"dst_amount"`
	TotalAmount        string       `json:"total_amount"`
	McRate             string       `json:"mc_rate"`
	BuyBackAmount      string       `json:"buy_back_amount"`
	McRateId           string       `json:"mc_rate_id"`
	OriginatingCountry string       `json:"originating_country"`
	DestinationCountry string       `json:"destination_country"`
	PurposeTransaction string       `json:"purpose_transaction"`
	SourceFund         string       `json:"source_fund"`
	Occupation         string       `json:"occupation"`
	RelationTo         string       `json:"relation_to"`
	BirthDate          string       `json:"birth_date"`
	BirthPlace         string       `json:"birth_place"`
	BirthCountry       string       `json:"birth_country"`
	IDType             string       `json:"id_type"`
	IDNumber           string       `json:"id_number"`
	Address            string       `json:"address"`
	Barangay           string       `json:"barangay"`
	City               string       `json:"city"`
	Province           string       `json:"province"`
	Country            string       `json:"country"`
	ContactNumber      string       `json:"contact_number"`
	CurrentAddress     NonexAddress `json:"current_address"`
	PermanentAddress   NonexAddress `json:"permanent_address"`
	RiskScore          string       `json:"risk_score"`
	RiskCriteria       string       `json:"risk_criteria"`
	ClientReferenceNo  string       `json:"client_reference_no"`
	FormType           string       `json:"form_type"`
	FormNumber         string       `json:"form_number"`
	PayoutType         json.Number  `json:"payout_type"`
	RemoteLocationID   json.Number  `json:"remote_location_id"`
	RemoteUserID       int          `json:"remote_user_id"`
	RemoteIPAddress    string       `json:"remote_ip_address"`
	IPAddress          string       `json:"ip_address"`
	ReferenceNumber    string       `json:"reference_number"`
	ZipCode            string       `json:"zip_code"`
	LocationName       string       `json:"location_name"`
	RateCategory       string       `json:"rate_category"`
	DsaCode            string       `json:"dsa_code"`
	DsaTrxType         string       `json:"dsa_trx_type"`
}

type InstaCashPayoutResponseBody struct {
	Code    json.Number           `json:"code"`
	Message string                `json:"message"`
	Result  InstaCashPayoutResult `json:"result"`
	RemcoID json.Number           `json:"remco_id"`
}

type InstaCashPayoutResult struct {
	ControlNumber string `json:"control_number"`
	Status        bool   `json:"status"`
	Remarks       string `json:"remarks"`
}

func (s *Svc) InstaCashPayout(ctx context.Context, sr InstaCashPayoutRequest) (*InstaCashPayoutResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("instacash/payout"), sr)
	if err != nil {
		return nil, err
	}

	ICPRes := &InstaCashPayoutResponseBody{}
	if err := json.Unmarshal(res, ICPRes); err != nil {
		return nil, err
	}
	return ICPRes, nil
}

package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type AYANNAHSendRequest struct {
	LocationID         json.Number `json:"location_id"`
	LocationName       string      `json:"location_name"`
	UserID             json.Number `json:"user_id"`
	TrxDate            string      `json:"trx_date"`
	CustomerID         string      `json:"customer_id"`
	CurrencyID         string      `json:"currency_id"`
	RemcoID            string      `json:"remco_id"`
	TrxType            string      `json:"trx_type"`
	IsDomestic         string      `json:"is_domestic"`
	CustomerName       string      `json:"customer_name"`
	ServiceCharge      string      `json:"service_charge"`
	RemoteLocationID   json.Number `json:"remote_location_id"`
	DstAmount          string      `json:"dst_amount"`
	TotalAmount        string      `json:"total_amount"`
	RemoteUserID       json.Number `json:"remote_user_id"`
	RemoteIPAddress    string      `json:"remote_ip_address"`
	PurposeTransaction string      `json:"purpose_transaction"`
	SourceFund         string      `json:"source_fund"`
	Occupation         string      `json:"occupation"`
	RelationTo         string      `json:"relation_to"`
	BirthDate          string      `json:"birth_date"`
	BirthPlace         string      `json:"birth_place"`
	BirthCountry       string      `json:"birth_country"`
	IDType             string      `json:"id_type"`
	IDNumber           string      `json:"id_number"`
	Address            string      `json:"address"`
	Barangay           string      `json:"barangay"`
	City               string      `json:"city"`
	Province           string      `json:"province"`
	ZipCode            string      `json:"zip_code"`
	Country            string      `json:"country"`
	ContactNumber      string      `json:"contact_number"`
	RiskScore          json.Number `json:"risk_score"`
	PayoutType         json.Number `json:"payout_type"`
	SenderName         string      `json:"sender_name"`
	ReceiverName       string      `json:"receiver_name"`
	PrincipalAmount    json.Number `json:"principal_amount"`
	ClientReferenceNo  string      `json:"client_reference_no"`
	ControlNumber      string      `json:"control_number"`
	McRate             string      `json:"mc_rate"`
	McRateID           string      `json:"mc_rate_id"`
	RateCategory       string      `json:"rate_category"`
	OriginatingCountry string      `json:"originating_country"`
	DestinationCountry string      `json:"destination_country"`
	FormType           string      `json:"form_type"`
	FormNumber         string      `json:"form_number"`
	IPAddress          string      `json:"ip_address"`
	ReferenceNumber    string      `json:"reference_number"`
	BuyBackAmount      string      `json:"buy_back_amount"`
	DsaCode            string      `json:"dsa_code"`
	DsaTrxType         string      `json:"dsa_trx_type"`
}

type AYANNAHSendResponseBody struct {
	Code    json.Number       `json:"code"`
	Message string            `json:"message"`
	Result  AYANNAHSendResult `json:"result"`
	RemcoID json.Number       `json:"remco_id"`
}

type AYANNAHSendResult struct {
	Message            string       `json:"message"`
	ID                 json.Number  `json:"id"`
	LocationID         json.Number  `json:"location_id"`
	UserID             json.Number  `json:"user_id"`
	TrxDate            string       `json:"trx_date"`
	CurrencyID         json.Number  `json:"currency_id"`
	RemcoID            json.Number  `json:"remco_id"`
	TrxType            json.Number  `json:"trx_type"`
	IsDomestic         json.Number  `json:"is_domestic"`
	CustomerID         json.Number  `json:"customer_id"`
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
	RateCategory       string       `json:"rate_category"`
	McRateID           json.Number  `json:"mc_rate_id"`
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
	RiskScore          json.Number  `json:"risk_score"`
	RiskCriteria       string       `json:"risk_criteria"`
	ClientReferenceNo  string       `json:"client_reference_no"`
	FormType           string       `json:"form_type"`
	FormNumber         string       `json:"form_number"`
	PayoutType         json.Number  `json:"payout_type"`
	RemoteLocationID   json.Number  `json:"remote_location_id"`
	RemoteUserID       json.Number  `json:"remote_user_id"`
	RemoteIPAddress    string       `json:"remote_ip_address"`
	IPAddress          string       `json:"ip_address"`
	CreatedAt          time.Time    `json:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at"`
	ReferenceNumber    string       `json:"reference_number"`
	ZipCode            string       `json:"zip_code"`
	Status             json.Number  `json:"status"`
	APIRequest         string       `json:"api_request"`
	SapForm            string       `json:"sap_form"`
	SapFormNumber      string       `json:"sap_form_number"`
	SapValidID1        string       `json:"sap_valid_id1"`
	SapValidID2        string       `json:"sap_valid_id2"`
	SapOboLastName     string       `json:"sap_obo_last_name"`
	SapOboFirstName    string       `json:"sap_obo_first_name"`
	SapOboMiddleName   string       `json:"sap_obo_middle_name"`
	AyannahStatus      string       `json:"ayannah_status"`
}

func (s *Svc) AYANNAHSendMoney(ctx context.Context, sr AYANNAHSendRequest) (*AYANNAHSendResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("ayannah/send"), sr)
	if err != nil {
		return nil, err
	}

	ayannahRes := &AYANNAHSendResponseBody{}
	if err := json.Unmarshal(res, ayannahRes); err != nil {
		return nil, err
	}
	return ayannahRes, nil
}

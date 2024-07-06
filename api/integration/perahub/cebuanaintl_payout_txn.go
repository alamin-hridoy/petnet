package perahub

import (
	"context"
	"encoding/json"
)

type CEBINTPayoutRequest struct {
	LocationID               json.Number  `json:"location_id"`
	LocationName             string       `json:"location_name"`
	UserID                   json.Number  `json:"user_id"`
	TrxDate                  string       `json:"trx_date"`
	CustomerID               string       `json:"customer_id"`
	CurrencyID               string       `json:"currency_id"`
	RemcoID                  json.Number  `json:"remco_id"`
	TrxType                  string       `json:"trx_type"`
	IsDomestic               string       `json:"is_domestic"`
	CustomerName             string       `json:"customer_name"`
	ServiceCharge            string       `json:"service_charge"`
	RemoteLocationID         json.Number  `json:"remote_location_id"`
	DstAmount                string       `json:"dst_amount"`
	TotalAmount              string       `json:"total_amount"`
	RemoteUserID             json.Number  `json:"remote_user_id"`
	RemoteIPAddress          string       `json:"remote_ip_address"`
	OriginatingCountry       string       `json:"originating_country"`
	DestinationCountry       string       `json:"destination_country"`
	PurposeTransaction       string       `json:"purpose_transaction"`
	SourceFund               string       `json:"source_fund"`
	Occupation               string       `json:"occupation"`
	RelationTo               string       `json:"relation_to"`
	BirthDate                string       `json:"birth_date"`
	BirthPlace               string       `json:"birth_place"`
	BirthCountry             string       `json:"birth_country"`
	IDType                   string       `json:"id_type"`
	IDNumber                 string       `json:"id_number"`
	Address                  string       `json:"address"`
	Barangay                 string       `json:"barangay"`
	City                     string       `json:"city"`
	Province                 string       `json:"province"`
	ZipCode                  string       `json:"zip_code"`
	Country                  string       `json:"country"`
	ContactNumber            string       `json:"contact_number"`
	CurrentAddress           NonexAddress `json:"current_address"`
	PermanentAddress         NonexAddress `json:"permanent_address"`
	RiskScore                json.Number  `json:"risk_score"`
	RiskCriteria             json.Number  `json:"risk_criteria"`
	FormType                 string       `json:"form_type"`
	FormNumber               string       `json:"form_number"`
	PayoutType               json.Number  `json:"payout_type"`
	SenderName               string       `json:"sender_name"`
	ReceiverName             string       `json:"receiver_name"`
	PrincipalAmount          json.Number  `json:"principal_amount"`
	ClientReferenceNo        string       `json:"client_reference_no"`
	ControlNumber            string       `json:"control_number"`
	ReferenceNumber          string       `json:"reference_number"`
	IdentificationTypeID     string       `json:"identification_type_id"`
	BeneficiaryID            string       `json:"beneficiary_id"`
	AgentCode                string       `json:"agent_code"`
	IDIssuedBy               string       `json:"id_issued_by"`
	IDIssuedState            string       `json:"id_issued_state"`
	IDIssuedCity             string       `json:"id_issued_city"`
	IDDateOfIssue            string       `json:"id_date_of_issue"`
	IDIssuingCountry         string       `json:"id_issuing_country"`
	IDIssuingState           string       `json:"id_issuing_state"`
	PassportIDIssuedCountry  string       `json:"passport_id_issued_country"`
	InternationalPartnerCode string       `json:"international_partner_code"`
	McRate                   string       `json:"mc_rate"`
	BuyBackAmount            string       `json:"buy_back_amount"`
	RateCategory             string       `json:"rate_category"`
	McRateID                 string       `json:"mc_rate_id"`
	DsaCode                  string       `json:"dsa_code"`
	DsaTrxType               string       `json:"dsa_trx_type"`
}

type CEBINTPayoutResponseBody struct {
	Code    json.Number        `json:"code"`
	Message string             `json:"message"`
	Result  CEBINTPayoutResult `json:"result"`
	RemcoID json.Number        `json:"remco_id"`
}

type CEBINTPayoutResult struct {
	Message string `json:"message"`
}

func (s *Svc) CEBINTPayout(ctx context.Context, sr CEBINTPayoutRequest) (*CEBINTPayoutResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("cebuana-international/payout"), sr)
	if err != nil {
		return nil, err
	}

	cebRes := &CEBINTPayoutResponseBody{}
	if err := json.Unmarshal(res, cebRes); err != nil {
		return nil, err
	}
	return cebRes, nil
}

package perahub

import (
	"context"
	"encoding/json"
)

type CebuanaSendRequest struct {
	LocationID         json.Number  `json:"location_id"`
	LocationName       string       `json:"location_name"`
	UserID             int          `json:"user_id"`
	TrxDate            string       `json:"trx_date"`
	CustomerID         string       `json:"customer_id"`
	CurrencyID         string       `json:"currency_id"`
	RemcoID            json.Number  `json:"remco_id"`
	TrxType            string       `json:"trx_type"`
	IsDomestic         string       `json:"is_domestic"`
	CustomerName       string       `json:"customer_name"`
	ServiceCharge      string       `json:"service_charge"`
	RemoteLocationID   json.Number  `json:"remote_location_id"`
	DstAmount          string       `json:"dst_amount"`
	TotalAmount        string       `json:"total_amount"`
	RemoteUserID       json.Number  `json:"remote_user_id"`
	RemoteIPAddress    string       `json:"remote_ip_address"`
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
	ZipCode            string       `json:"zip_code"`
	Country            string       `json:"country"`
	ContactNumber      string       `json:"contact_number"`
	CurrentAddress     NonexAddress `json:"current_address"`
	PermanentAddress   NonexAddress `json:"permanent_address"`
	RiskScore          json.Number  `json:"risk_score"`
	FormType           string       `json:"form_type"`
	FormNumber         string       `json:"form_number"`
	BeneficiaryID      string       `json:"beneficiary_id"`
	AgentCode          string       `json:"agent_code"`
	McRateID           string       `json:"mc_rate_id"`
	RateCtg            string       `json:"rate_category"`
	BbAmount           string       `json:"buy_back_amount"`
	McRate             string       `json:"mc_rate"`
	SendCurrencyID     string       `json:"send_currency_id"`
	PrincipalAmount    string       `json:"principal_amount"`
	SenderName         string       `json:"sender_name"`
	ReceiverName       string       `json:"receiver_name"`
	ClientRefNo        string       `json:"client_reference_no"`
	PayoutType         json.Number  `json:"payout_type"`
	DsaCode            string       `json:"dsa_code"`
	DsaTrxType         string       `json:"dsa_trx_type"`
}

type CebuanaSendResponseBody struct {
	Code    json.Number       `json:"code"`
	Message string            `json:"message"`
	Result  CebuanaSendResult `json:"result"`
	RemcoID json.Number       `json:"remco_id"`
}

type CebuanaSendResult struct {
	ResultStatus string      `json:"ResultStatus"`
	MessageID    json.Number `json:"MessageID"`
	LogID        json.Number `json:"LogID"`
	ControlNo    string      `json:"ControlNumber"`
	ServiceFee   string      `json:"ServiceFee"`
}

func (s *Svc) CebuanaSendMoney(ctx context.Context, sr CebuanaSendRequest) (*CebuanaSendResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("cebuana/send"), sr)
	if err != nil {
		return nil, err
	}

	cebRes := &CebuanaSendResponseBody{}
	if err := json.Unmarshal(res, cebRes); err != nil {
		return nil, err
	}
	return cebRes, nil
}

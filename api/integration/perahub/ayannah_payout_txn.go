package perahub

import (
	"context"
	"encoding/json"
)

type AYANNAHPayoutRequest struct {
	LocationID         json.Number `json:"location_id"`
	LocationName       string      `json:"location_name"`
	UserID             json.Number `json:"user_id"`
	TrxDate            string      `json:"trx_date"`
	CustomerID         string      `json:"customer_id"`
	CurrencyID         json.Number `json:"currency_id"`
	RemcoID            json.Number `json:"remco_id"`
	TrxType            json.Number `json:"trx_type"`
	IsDomestic         json.Number `json:"is_domestic"`
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
	ControlNumber      string      `json:"control_number"`
	ReferenceNumber    string      `json:"reference_number"`
	McRate             string      `json:"mc_rate"`
	McRateID           string      `json:"mc_rate_id"`
	RateCategory       string      `json:"rate_category"`
	OriginatingCountry string      `json:"originating_country"`
	DestinationCountry string      `json:"destination_country"`
	FormType           string      `json:"form_type"`
	FormNumber         string      `json:"form_number"`
	ClientReferenceNo  string      `json:"client_reference_no"`
	BuyBackAmount      string      `json:"buy_back_amount"`
	IPAddress          string      `json:"ip_address"`
	DsaCode            string      `json:"dsa_code"`
	DsaTrxType         string      `json:"dsa_trx_type"`
}

type AYANNAHPayoutResponseBody struct {
	Code    json.Number         `json:"code"`
	Message string              `json:"message"`
	Result  AYANNAHPayoutResult `json:"result"`
	RemcoID json.Number         `json:"remco_id"`
}

type AYANNAHPayoutResult struct {
	Message string `json:"message"`
}

func (s *Svc) AYANNAHPayout(ctx context.Context, sr AYANNAHPayoutRequest) (*AYANNAHPayoutResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("ayannah/payout"), sr)
	if err != nil {
		return nil, err
	}

	cebRes := &AYANNAHPayoutResponseBody{}
	if err := json.Unmarshal(res, cebRes); err != nil {
		return nil, err
	}
	return cebRes, nil
}

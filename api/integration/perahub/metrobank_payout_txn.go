package perahub

import (
	"context"
	"encoding/json"
)

type MBPayoutRequest struct {
	LocationId         json.Number `json:"location_id"`
	LocationName       string      `json:"location_name"`
	UserId             json.Number `json:"user_id"`
	TrxDate            string      `json:"trx_date"`
	CustomerId         string      `json:"customer_id"`
	CurrencyId         string      `json:"currency_id"`
	RemcoId            string      `json:"remco_id"`
	TrxType            string      `json:"trx_type"`
	IsDomestic         json.Number `json:"is_domestic"`
	CustomerName       string      `json:"customer_name"`
	ServiceCharge      string      `json:"service_charge"`
	RemoteLocationId   json.Number `json:"remote_location_id"`
	DstAmount          string      `json:"dst_amount"`
	TotalAmount        string      `json:"total_amount"`
	RemoteUserId       json.Number `json:"remote_user_id"`
	RemoteIpAddress    string      `json:"remote_ip_address"`
	PurposeTransaction string      `json:"purpose_transaction"`
	SourceFund         string      `json:"source_fund"`
	Occupation         string      `json:"occupation"`
	RelationTo         string      `json:"relation_to"`
	BirthDate          string      `json:"birth_date"`
	BirthPlace         string      `json:"birth_place"`
	BirthCountry       string      `json:"birth_country"`
	IdType             string      `json:"id_type"`
	IdNumber           string      `json:"id_number"`
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
	BuyBackAmount      string      `json:"buy_back_amount"`
	RateCategory       string      `json:"rate_category"`
	McRateId           string      `json:"mc_rate_id"`
	OriginatingCountry string      `json:"originating_country"`
	DestinationCountry string      `json:"destination_country"`
	ClientReferenceNo  string      `json:"client_reference_no"`
	FormType           string      `json:"form_type"`
	FormNumber         string      `json:"form_number"`
	Currency           string      `json:"currency"`
	DsaCode            string      `json:"dsa_code"`
	DsaTrxType         string      `json:"dsa_trx_type"`
}

type MBPayoutResponseBody struct {
	Code    json.Number `json:"code"`
	Msg     string      `json:"message"`
	Result  MBPayResult `json:"result"`
	RemcoID json.Number `json:"remco_id"`
}

type MBPayResult struct {
	RefNo           string `json:"reference_number"`
	ClientRefNo     string `json:"client_reference_no"`
	ControlNo       string `json:"control_number"`
	StatusText      string `json:"status_text"`
	PrincipalAmount string `json:"principal_amount"`
	RcvName         string `json:"receiver_name"`
	Address         string `json:"address"`
	ReceiptNo       string `json:"receipt_no"`
	ContactNumber   string `json:"contact_number"`
}

func (s *Svc) MBPayout(ctx context.Context, sr MBPayoutRequest) (*MBPayoutResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("metrobank/payout"), sr)
	if err != nil {
		return nil, err
	}

	rb := &MBPayoutResponseBody{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

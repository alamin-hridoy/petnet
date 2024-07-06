package perahub

import (
	"context"
	"encoding/json"
	"time"
)

type PerahubRemitValidateRequest struct {
	LocationID         string              `json:"location_id"`
	UserID             string              `json:"user_id"`
	TrxDate            string              `json:"trx_date"`
	CurrencyID         string              `json:"currency_id"`
	RemcoID            string              `json:"remco_id"`
	TrxType            string              `json:"trx_type"`
	IsDomestic         string              `json:"is_domestic"`
	CustomerID         string              `json:"customer_id"`
	ControlNumber      string              `json:"control_number"`
	ClientReferenceNo  string              `json:"client_reference_no"`
	CustomerName       string              `json:"customer_name"`
	SenderName         string              `json:"sender_name"`
	ReceiverName       string              `json:"receiver_name"`
	PrincipalAmount    int                 `json:"principal_amount"`
	ServiceCharge      string              `json:"service_charge"`
	DstAmount          string              `json:"dst_amount"`
	TotalAmount        string              `json:"total_amount"`
	McRate             string              `json:"mc_rate"`
	BuyBackAmount      string              `json:"buy_back_amount"`
	RateCategory       string              `json:"rate_category"`
	McRateID           string              `json:"mc_rate_id"`
	OriginatingCountry string              `json:"originating_country"`
	DestinationCountry string              `json:"destination_country"`
	PurposeTransaction string              `json:"purpose_transaction"`
	SourceFund         string              `json:"source_fund"`
	Occupation         string              `json:"occupation"`
	RelationTo         string              `json:"relation_to"`
	BirthPlace         string              `json:"birth_place"`
	BirthCountry       string              `json:"birth_country"`
	IDType             string              `json:"id_type"`
	IDNumber           string              `json:"id_number"`
	Address            string              `json:"address"`
	Barangay           string              `json:"barangay"`
	City               string              `json:"city"`
	Province           string              `json:"province"`
	ZipCode            string              `json:"zip_code"`
	Country            string              `json:"country"`
	ContactNumber      string              `json:"contact_number"`
	FormType           string              `json:"form_type"`
	FormNumber         string              `json:"form_number"`
	PayoutType         string              `json:"payout_type"`
	RemoteLocationID   string              `json:"remote_location_id"`
	RemoteUserID       string              `json:"remote_user_id"`
	RemoteIPAddress    string              `json:"remote_ip_address"`
	IPAddress          string              `json:"ip_address"`
	ReferenceNumber    string              `json:"reference_number"`
	CurrentAddress     PerahubRemitAddress `json:"current_address"`
	PermanentAddress   PerahubRemitAddress `json:"permanent_address"`
	SenderLastName     string              `json:"sender_last_name"`
	SenderFirstName    string              `json:"sender_first_name"`
	SenderMiddleName   string              `json:"sender_middle_name"`
	ReceiverLastName   string              `json:"receiver_last_name"`
	ReceiverFirstName  string              `json:"receiver_first_name"`
	ReceiverMiddleName string              `json:"receiver_middle_name"`
	PartnerCode        string              `json:"partner_code"`
	DsaCode            string              `json:"dsa_code"`
	DsaTrxType         string              `json:"dsa_trx_type"`
}

type PerahubRemitValidateResponseBody struct {
	Code    string                     `json:"code"`
	Message string                     `json:"message"`
	Result  PerahubRemitValidateResult `json:"result"`
}

type PerahubRemitValidateResult struct {
	LocationID         int                            `json:"location_id"`
	UserID             string                         `json:"user_id"`
	TrxDate            string                         `json:"trx_date"`
	CurrencyID         string                         `json:"currency_id"`
	RemcoID            string                         `json:"remco_id"`
	TrxType            string                         `json:"trx_type"`
	IsDomestic         string                         `json:"is_domestic"`
	CustomerID         string                         `json:"customer_id"`
	ControlNumber      string                         `json:"control_number"`
	ClientReferenceNo  string                         `json:"client_reference_no"`
	CustomerName       string                         `json:"customer_name"`
	SenderName         string                         `json:"sender_name"`
	ReceiverName       string                         `json:"receiver_name"`
	PrincipalAmount    int                            `json:"principal_amount"`
	ServiceCharge      string                         `json:"service_charge"`
	DstAmount          string                         `json:"dst_amount"`
	TotalAmount        string                         `json:"total_amount"`
	McRate             string                         `json:"mc_rate"`
	BuyBackAmount      string                         `json:"buy_back_amount"`
	RateCategory       string                         `json:"rate_category"`
	McRateID           string                         `json:"mc_rate_id"`
	OriginatingCountry string                         `json:"originating_country"`
	DestinationCountry string                         `json:"destination_country"`
	PurposeTransaction string                         `json:"purpose_transaction"`
	SourceFund         string                         `json:"source_fund"`
	Occupation         string                         `json:"occupation"`
	RelationTo         string                         `json:"relation_to"`
	BirthPlace         string                         `json:"birth_place"`
	BirthCountry       string                         `json:"birth_country"`
	IDType             string                         `json:"id_type"`
	IDNumber           string                         `json:"id_number"`
	Address            string                         `json:"address"`
	Barangay           string                         `json:"barangay"`
	City               string                         `json:"city"`
	Province           string                         `json:"province"`
	ZipCode            string                         `json:"zip_code"`
	Country            string                         `json:"country"`
	ContactNumber      string                         `json:"contact_number"`
	FormType           string                         `json:"form_type"`
	FormNumber         string                         `json:"form_number"`
	PayoutType         string                         `json:"payout_type"`
	RemoteLocationID   string                         `json:"remote_location_id"`
	RemoteUserID       string                         `json:"remote_user_id"`
	RemoteIPAddress    string                         `json:"remote_ip_address"`
	IPAddress          string                         `json:"ip_address"`
	ReferenceNumber    string                         `json:"reference_number"`
	CurrentAddress     PerahubRemitAddress            `json:"current_address"`
	PermanentAddress   PerahubRemitAddress            `json:"permanent_address"`
	APIRequest         PerahubRemitValidateAPIRequest `json:"api_request"`
	UpdatedAt          time.Time                      `json:"updated_at"`
	CreatedAt          time.Time                      `json:"created_at"`
	ID                 int                            `json:"id"`
}

type PerahubRemitValidateAddress struct {
	Address1 string `json:"address_1"`
	Barangay string `json:"barangay"`
	City     string `json:"city"`
	Province string `json:"province"`
	ZipCode  string `json:"zip_code"`
	Country  string `json:"country"`
}

type PerahubRemitValidateAPIRequest struct {
	LocationID         int                 `json:"location_id"`
	UserID             string              `json:"user_id"`
	TrxDate            string              `json:"trx_date"`
	CurrencyID         string              `json:"currency_id"`
	RemcoID            string              `json:"remco_id"`
	TrxType            string              `json:"trx_type"`
	IsDomestic         string              `json:"is_domestic"`
	CustomerID         string              `json:"customer_id"`
	ControlNumber      string              `json:"control_number"`
	ClientReferenceNo  string              `json:"client_reference_no"`
	CustomerName       string              `json:"customer_name"`
	SenderName         string              `json:"sender_name"`
	ReceiverName       string              `json:"receiver_name"`
	PrincipalAmount    int                 `json:"principal_amount"`
	ServiceCharge      string              `json:"service_charge"`
	DstAmount          string              `json:"dst_amount"`
	TotalAmount        string              `json:"total_amount"`
	McRate             string              `json:"mc_rate"`
	BuyBackAmount      string              `json:"buy_back_amount"`
	RateCategory       string              `json:"rate_category"`
	McRateID           string              `json:"mc_rate_id"`
	OriginatingCountry string              `json:"originating_country"`
	DestinationCountry string              `json:"destination_country"`
	PurposeTransaction string              `json:"purpose_transaction"`
	SourceFund         string              `json:"source_fund"`
	Occupation         string              `json:"occupation"`
	RelationTo         string              `json:"relation_to"`
	BirthPlace         string              `json:"birth_place"`
	BirthCountry       string              `json:"birth_country"`
	IDType             string              `json:"id_type"`
	IDNumber           string              `json:"id_number"`
	Address            string              `json:"address"`
	Barangay           string              `json:"barangay"`
	City               string              `json:"city"`
	Province           string              `json:"province"`
	ZipCode            string              `json:"zip_code"`
	Country            string              `json:"country"`
	ContactNumber      string              `json:"contact_number"`
	FormType           string              `json:"form_type"`
	FormNumber         string              `json:"form_number"`
	PayoutType         string              `json:"payout_type"`
	RemoteLocationID   string              `json:"remote_location_id"`
	RemoteUserID       string              `json:"remote_user_id"`
	RemoteIPAddress    string              `json:"remote_ip_address"`
	IPAddress          string              `json:"ip_address"`
	ReferenceNumber    string              `json:"reference_number"`
	CurrentAddress     PerahubRemitAddress `json:"current_address"`
	PermanentAddress   PerahubRemitAddress `json:"permanent_address"`
	SenderLastName     string              `json:"sender_last_name"`
	SenderFirstName    string              `json:"sender_first_name"`
	SenderMiddleName   string              `json:"sender_middle_name"`
	ReceiverLastName   string              `json:"receiver_last_name"`
	ReceiverFirstName  string              `json:"receiver_first_name"`
	ReceiverMiddleName string              `json:"receiver_middle_name"`
	PartnerCode        string              `json:"partner_code"`
}

func (s *Svc) PerahubRemitValidate(ctx context.Context, vr PerahubRemitValidateRequest) (*PerahubRemitValidateResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("/perahub-remit/payout/validate"), vr)
	if err != nil {
		return nil, err
	}

	prvRes := &PerahubRemitValidateResponseBody{}
	if err := json.Unmarshal(res, prvRes); err != nil {
		return nil, err
	}
	return prvRes, nil
}

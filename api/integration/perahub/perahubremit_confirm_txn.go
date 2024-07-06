package perahub

import (
	"context"
	"encoding/json"
)

type PerahubRemitConfirmRequest struct {
	ID              int    `json:"id"`
	ReferenceNumber string `json:"reference_number"`
	PartnerCode     string `json:"partner_code"`
}

type PerahubRemitConfirmResponseBody struct {
	Code    string             `json:"code"`
	Message string             `json:"message"`
	Result  PerahubRemitResult `json:"result"`
}

type PerahubRemitResult struct {
	ID                 int                           `json:"id"`
	LocationID         int                           `json:"location_id"`
	UserID             int                           `json:"user_id"`
	TrxDate            string                        `json:"trx_date"`
	CurrencyID         int                           `json:"currency_id"`
	RemcoID            int                           `json:"remco_id"`
	TrxType            int                           `json:"trx_type"`
	IsDomestic         int                           `json:"is_domestic"`
	CustomerID         int                           `json:"customer_id"`
	CustomerName       string                        `json:"customer_name"`
	ControlNumber      string                        `json:"control_number"`
	SenderName         string                        `json:"sender_name"`
	ReceiverName       string                        `json:"receiver_name"`
	PrincipalAmount    string                        `json:"principal_amount"`
	ServiceCharge      string                        `json:"service_charge"`
	DstAmount          string                        `json:"dst_amount"`
	TotalAmount        string                        `json:"total_amount"`
	McRate             string                        `json:"mc_rate"`
	BuyBackAmount      string                        `json:"buy_back_amount"`
	RateCategory       string                        `json:"rate_category"`
	McRateID           int                           `json:"mc_rate_id"`
	OriginatingCountry string                        `json:"originating_country"`
	DestinationCountry string                        `json:"destination_country"`
	PurposeTransaction string                        `json:"purpose_transaction"`
	SourceFund         string                        `json:"source_fund"`
	Occupation         string                        `json:"occupation"`
	RelationTo         string                        `json:"relation_to"`
	BirthPlace         string                        `json:"birth_place"`
	BirthCountry       string                        `json:"birth_country"`
	IDType             string                        `json:"id_type"`
	IDNumber           string                        `json:"id_number"`
	Address            string                        `json:"address"`
	Barangay           string                        `json:"barangay"`
	City               string                        `json:"city"`
	Province           string                        `json:"province"`
	Country            string                        `json:"country"`
	ContactNumber      string                        `json:"contact_number"`
	RiskScore          int                           `json:"risk_score"`
	RiskCriteria       string                        `json:"risk_criteria"`
	ClientReferenceNo  string                        `json:"client_reference_no"`
	FormType           string                        `json:"form_type"`
	FormNumber         string                        `json:"form_number"`
	PayoutType         int                           `json:"payout_type"`
	RemoteLocationID   int                           `json:"remote_location_id"`
	RemoteUserID       int                           `json:"remote_user_id"`
	RemoteIPAddress    string                        `json:"remote_ip_address"`
	IPAddress          string                        `json:"ip_address"`
	CreatedAt          string                        `json:"created_at"`
	UpdatedAt          string                        `json:"updated_at"`
	ReferenceNumber    string                        `json:"reference_number"`
	ZipCode            string                        `json:"zip_code"`
	Status             int                           `json:"status"`
	CurrentAddress     PerahubRemitConfirmAddress    `json:"current_address"`
	PermanentAddress   PerahubRemitConfirmAddress    `json:"permanent_address"`
	APIRequest         PerahubRemitConfirmAPIRequest `json:"api_request"`
}
type PerahubRemitConfirmAPIRequest struct {
	City                       string                     `json:"city"`
	Address                    string                     `json:"address"`
	Country                    string                     `json:"country"`
	IDType                     string                     `json:"id_type"`
	McRate                     string                     `json:"mc_rate"`
	UserID                     string                     `json:"user_id"`
	Barangay                   string                     `json:"barangay"`
	Province                   string                     `json:"province"`
	RemcoID                    string                     `json:"remco_id"`
	TrxDate                    string                     `json:"trx_date"`
	TrxType                    string                     `json:"trx_type"`
	ZipCode                    string                     `json:"zip_code"`
	FormType                   string                     `json:"form_type"`
	IDNumber                   string                     `json:"id_number"`
	DstAmount                  string                     `json:"dst_amount"`
	IPAddress                  string                     `json:"ip_address"`
	McRateID                   string                     `json:"mc_rate_id"`
	Occupation                 string                     `json:"occupation"`
	BirthPlace                 string                     `json:"birth_place"`
	CurrencyID                 string                     `json:"currency_id"`
	CustomerID                 string                     `json:"customer_id"`
	FormNumber                 string                     `json:"form_number"`
	IsDomestic                 string                     `json:"is_domestic"`
	LocationID                 int                        `json:"location_id"`
	PayoutType                 string                     `json:"payout_type"`
	RelationTo                 string                     `json:"relation_to"`
	SenderName                 string                     `json:"sender_name"`
	SourceFund                 string                     `json:"source_fund"`
	PartnerCode                string                     `json:"partner_code"`
	TotalAmount                string                     `json:"total_amount"`
	BirthCountry               string                     `json:"birth_country"`
	CustomerName               string                     `json:"customer_name"`
	RateCategory               string                     `json:"rate_category"`
	ReceiverName               string                     `json:"receiver_name"`
	ContactNumber              string                     `json:"contact_number"`
	ControlNumber              string                     `json:"control_number"`
	RemoteUserID               string                     `json:"remote_user_id"`
	ServiceCharge              string                     `json:"service_charge"`
	BuyBackAmount              string                     `json:"buy_back_amount"`
	PrincipalAmount            int                        `json:"principal_amount"`
	ReferenceNumber            string                     `json:"reference_number"`
	SenderLastName             string                     `json:"sender_last_name"`
	RemoteIPAddress            string                     `json:"remote_ip_address"`
	SenderFirstName            string                     `json:"sender_first_name"`
	ReceiverLastName           string                     `json:"receiver_last_name"`
	RemoteLocationID           string                     `json:"remote_location_id"`
	SenderMiddleName           string                     `json:"sender_middle_name"`
	ClientReferenceNo          string                     `json:"client_reference_no"`
	DestinationCountry         string                     `json:"destination_country"`
	OriginatingCountry         string                     `json:"originating_country"`
	PurposeTransaction         string                     `json:"purpose_transaction"`
	ReceiverFirstName          string                     `json:"receiver_first_name"`
	ReceiverMiddleName         string                     `json:"receiver_middle_name"`
	APIRequestCurrentAddress   PerahubRemitConfirmAddress `json:"current_address"`
	APIRequestPermanentAddress PerahubRemitConfirmAddress `json:"permanent_address"`
}

type PerahubRemitConfirmAddress struct {
	City     string `json:"city"`
	Country  string `json:"country"`
	Barangay string `json:"barangay"`
	Province string `json:"province"`
	ZipCode  string `json:"zip_code"`
	Address1 string `json:"address_1"`
}

func (s *Svc) PerahubRemitconfirm(ctx context.Context, sr PerahubRemitConfirmRequest) (*PerahubRemitConfirmResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("/perahub-remit/payout/confirm"), sr)
	if err != nil {
		return nil, err
	}

	PRCRes := &PerahubRemitConfirmResponseBody{}
	if err := json.Unmarshal(res, PRCRes); err != nil {
		return nil, err
	}
	return PRCRes, nil
}

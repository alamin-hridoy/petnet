package perahub

import (
	"context"
	"encoding/json"
)

type USSCSendRequest struct {
	ControlNo          string       `json:"control_number"`
	McRate             string       `json:"mc_rate"`
	BbAmount           string       `json:"buy_back_amount"`
	RateCtg            string       `json:"rate_category"`
	McRateID           string       `json:"mc_rate_id"`
	BranchCode         string       `json:"branch_code"`
	LocationID         json.Number  `json:"location_id"`
	LocationName       string       `json:"location_name"`
	UserID             json.Number  `json:"user_id"`
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
	RiskCrt            json.Number  `json:"risk_criteria"`
	FormType           string       `json:"form_type"`
	FormNumber         string       `json:"form_number"`
	PayoutType         json.Number  `json:"payout_type"`
	SenderName         string       `json:"sender_name"`
	ReceiverName       string       `json:"receiver_name"`
	PrincipalAmount    string       `json:"principal_amount"`
	ClientRefNo        string       `json:"client_reference_no"`
	ReferenceNo        string       `json:"reference_number"`
	SendFName          string       `json:"sender_first_name"`
	SendMName          string       `json:"sender_middle_name"`
	SendLName          string       `json:"sender_last_name"`
	RecFName           string       `json:"receiver_first_name"`
	RecMName           string       `json:"receiver_middle_name"`
	RecLName           string       `json:"receiver_last_name"`
	RecConNo           string       `json:"receiver_contact_number"`
	KycVer             bool         `json:"kyc_verified"`
	Gender             string       `json:"gender"`
	DsaCode            string       `json:"dsa_code"`
	DsaTrxType         string       `json:"dsa_trx_type"`
}

type USSCSendResponseBody struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Result  USSCSendResult `json:"result"`
	RemcoID json.Number    `json:"remco_id"`
}

type USSCSendResult struct {
	ControlNo          string `json:"control_number"`
	TrxDate            string `json:"trx_date"`
	SendFName          string `json:"sender_first_name"`
	SendMName          string `json:"sender_middle_name"`
	SendLName          string `json:"sender_last_name"`
	PrincipalAmount    string `json:"principal_amount"`
	ServiceCharge      string `json:"service_charge"`
	TotalAmount        string `json:"total_amount"`
	RecFName           string `json:"receiver_first_name"`
	RecMName           string `json:"receiver_middle_name"`
	RecLName           string `json:"receiver_last_name"`
	ContactNumber      string `json:"contact_number"`
	RelationTo         string `json:"relation_to"`
	PurposeTransaction string `json:"purpose_transaction"`
	ReferenceNo        string `json:"reference_number"`
}

func (s *Svc) USSCSendMoney(ctx context.Context, sr USSCSendRequest) (*USSCSendResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("ussc/send"), sr)
	if err != nil {
		return nil, err
	}

	usscRes := &USSCSendResponseBody{}
	if err := json.Unmarshal(res, usscRes); err != nil {
		return nil, err
	}
	return usscRes, nil
}

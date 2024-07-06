package perahub

import (
	"context"
	"encoding/json"
)

type UNTPayoutRequest struct {
	ClientReferenceNo  string      `json:"client_reference_no"`
	ReferenceNumber    string      `json:"reference_number"`
	LocationCode       string      `json:"location_code"`
	LocationID         string      `json:"location_id"`
	LocationName       string      `json:"location_name"`
	Gender             string      `json:"gender"`
	ControlNumber      string      `json:"control_number"`
	Currency           string      `json:"currency"`
	PrincipalAmount    json.Number `json:"principal_amount"`
	IDNumber           string      `json:"id_number"`
	IDType             string      `json:"id_type"`
	IDIssuedBY         string      `json:"id_issued_by"`
	IDDOIssue          string      `json:"id_date_of_issue"`
	IDExpDate          string      `json:"id_expiration_date"`
	ContactNumber      string      `json:"contact_number"`
	Address            string      `json:"address"`
	City               string      `json:"city"`
	Province           string      `json:"province"`
	Country            string      `json:"country"`
	ZipCode            string      `json:"zip_code"`
	State              string      `json:"state"`
	Nationality        string      `json:"nationality"`
	BirthDate          string      `json:"birth_date"`
	BirthCountry       string      `json:"birth_country"`
	Occupation         string      `json:"occupation"`
	UserID             int         `json:"user_id"`
	TrxDate            string      `json:"trx_date"`
	CustomerID         string      `json:"customer_id"`
	CurrencyID         string      `json:"currency_id"`
	RemcoID            string      `json:"remco_id"`
	TrxType            string      `json:"trx_type"`
	IsDomestic         string      `json:"is_domestic"`
	CustomerName       string      `json:"customer_name"`
	ServiceCharge      string      `json:"service_charge"`
	RemoteLocationID   string      `json:"remote_location_id"`
	DstAmount          string      `json:"dst_amount"`
	TotalAmount        string      `json:"total_amount"`
	BuyBackAmount      string      `json:"buy_back_amount"`
	McRateId           string      `json:"mc_rate_id"`
	McRate             string      `json:"mc_rate"`
	RemoteIPAddress    string      `json:"remote_ip_address"`
	RemoteUserID       int         `json:"remote_user_id"`
	OriginatingCountry string      `json:"originating_country"`
	DestinationCountry string      `json:"destination_country"`
	PurposeTransaction string      `json:"purpose_transaction"`
	SourceFund         string      `json:"source_fund"`
	RelationTo         string      `json:"relation_to"`
	BirthPlace         string      `json:"birth_place"`
	Barangay           string      `json:"barangay"`
	RiskScore          string      `json:"risk_score"`
	RiskCriteria       string      `json:"risk_criteria"`
	FormType           string      `json:"form_type"`
	FormNumber         string      `json:"form_number"`
	PayoutType         string      `json:"payout_type"`
	SenderName         string      `json:"sender_name"`
	ReceiverName       string      `json:"receiver_name"`
	SendFName          string      `json:"sender_first_name"`
	SendMName          string      `json:"sender_middle_name"`
	SendLName          string      `json:"sender_last_name"`
	RecFName           string      `json:"receiver_first_name"`
	RecMName           string      `json:"receiver_middle_name"`
	RecLName           string      `json:"receiver_last_name"`
	DeviceID           string      `json:"device_id"`
	AgentID            string      `json:"agent_id"`
	AgentCode          string      `json:"agent_code"`
	OrderNumber        string      `json:"order_number"`
	IPAddress          string      `json:"ip_address"`
	RateCategory       string      `json:"rate_category"`
	DsaCode            string      `json:"dsa_code"`
	DsaTrxType         string      `json:"dsa_trx_type"`
}

type UNTPayoutResponseBody struct {
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Result  UNTPayoutResult `json:"result"`
	RemcoID json.Number     `json:"remco_id"`
}

type UNTPayoutResult struct {
	ResponseCode       string `json:"response_code"`
	ControlNumber      string `json:"control_number"`
	PrincipalAmount    string `json:"principal_amount"`
	Currency           string `json:"currency"`
	CreationDate       string `json:"creation_date"`
	ReceiverName       string `json:"receiver_name"`
	Address            string `json:"address"`
	City               string `json:"city"`
	Country            string `json:"country"`
	SenderName         string `json:"sender_name"`
	ZipCode            string `json:"zip_code"`
	OriginatingCountry string `json:"originating_country"`
	DestinationCountry string `json:"destination_country"`
	ContactNumber      string `json:"contact_number"`
	FmtSenderName      string `json:"formatted_sender_name"`
	FmtReceiverName    string `json:"formatted_receiver_name"`
}

func (s *Svc) UNTPayout(ctx context.Context, sr UNTPayoutRequest) (*UNTPayoutResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("uniteller/payout"), sr)
	if err != nil {
		return nil, err
	}

	UNTPRes := &UNTPayoutResponseBody{}
	if err := json.Unmarshal(res, UNTPRes); err != nil {
		return nil, err
	}
	return UNTPRes, nil
}

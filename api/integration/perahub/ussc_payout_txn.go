package perahub

import (
	"context"
	"encoding/json"
)

type USSCPayoutRequest struct {
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
	DstAmount          json.Number  `json:"dst_amount"`
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
	PayoutType         json.Number  `json:"payout_type"`
	SenderName         string       `json:"sender_name"`
	ReceiverName       string       `json:"receiver_name"`
	PrincipalAmount    json.Number  `json:"principal_amount"`
	ControlNumber      string       `json:"control_number"`
	ReferenceNumber    string       `json:"reference_number"`
	SenderFirstName    string       `json:"sender_first_name"`
	SenderLastName     string       `json:"sender_last_name"`
	ReceiverFirstName  string       `json:"receiver_first_name"`
	ReceiverLastName   string       `json:"receiver_last_name"`
	BranchCode         string       `json:"branch_code"`
	ClientReferenceNo  string       `json:"client_reference_no"`
	FormType           string       `json:"form_type"`
	FormNumber         string       `json:"form_number"`
	McRate             string       `json:"mc_rate"`
	BuyBackAmount      string       `json:"buy_back_amount"`
	RateCategory       string       `json:"rate_category"`
	McRateId           string       `json:"mc_rate_id"`
	DsaCode            string       `json:"dsa_code"`
	DsaTrxType         string       `json:"dsa_trx_type"`
}

type USSCPayoutResponseBody struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Result  USSCPayResult `json:"result"`
	RemcoID int           `json:"remco_id"`
}

type USSCPayResult struct {
	Spcn            string      `json:"SPCN"`
	SendPk          string      `json:"SendPk"`
	SendPw          string      `json:"SendPw"`
	SendDate        json.Number `json:"SendDate"`
	SendJournalNo   json.Number `json:"SendJournalNo"`
	SendLastName    string      `json:"SendLastName"`
	SendFirstName   string      `json:"SendFirstName"`
	SendMiddleName  string      `json:"SendMiddleName"`
	PayAmount       string      `json:"PayAmount"`
	SendFee         string      `json:"SendFee"`
	SendVat         string      `json:"SendVat"`
	SendFeeAfterVat string      `json:"SendFeeAfterVat"`
	SendTotalAmount string      `json:"SendTotalAmount"`
	PayPk           string      `json:"PayPk"`
	PayPw           string      `json:"PayPw"`
	PayLastName     string      `json:"PayLastName"`
	PayFirstName    string      `json:"PayFirstName"`
	PayMiddleName   string      `json:"PayMiddleName"`
	Relationship    string      `json:"Relationship"`
	Purpose         string      `json:"Purpose"`
	PromoCode       string      `json:"PromoCode"`
	PayBranchCode   string      `json:"PayBranchCode"`
	Remarks         string      `json:"Remarks"`
	OrNo            string      `json:"OrNo"`
	OboBranchCode   string      `json:"OboBranchCode"`
	OboUserID       string      `json:"OboUserId"`
	Message         string      `json:"Message"`
	Code            string      `json:"Code"`
	NewScreen       string      `json:"NewScreen"`
	JournalNo       string      `json:"JournalNo"`
	ProcessDate     string      `json:"ProcessDate"`
	ReferenceNo     string      `json:"reference_no"`
}

func (s *Svc) USSCPayout(ctx context.Context, sr USSCPayoutRequest) (*USSCPayoutResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("ussc/payout"), sr)
	if err != nil {
		return nil, err
	}

	rb := &USSCPayoutResponseBody{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

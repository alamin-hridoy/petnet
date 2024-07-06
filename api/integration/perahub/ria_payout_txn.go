package perahub

import (
	"context"
	"encoding/json"
)

type RiaPayoutRequest struct {
	LocationID         json.Number `json:"location_id"`
	LocationName       string      `json:"location_name"`
	UserID             json.Number `json:"user_id"`
	TrxDate            string      `json:"trx_date"`
	CustomerID         string      `json:"customer_id"`
	CurrencyID         json.Number `json:"currency_id"`
	RemcoID            json.Number `json:"remco_id"`
	TrxType            string      `json:"trx_type"`
	IsDomestic         json.Number `json:"is_domestic"`
	CustomerName       string      `json:"customer_name"`
	ServiceCharge      string      `json:"service_charge"`
	RmtLocID           json.Number `json:"remote_location_id"`
	DstAmount          string      `json:"dst_amount"`
	TotalAmount        string      `json:"total_amount"`
	RmtUserID          json.Number `json:"remote_user_id"`
	RmtIPAddr          string      `json:"remote_ip_address"`
	PurposeTxn         string      `json:"purpose_transaction"`
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
	RcvName            string      `json:"receiver_name"`
	PnplAmt            string      `json:"principal_amount"`
	ControlNo          string      `json:"control_number"`
	RefNo              string      `json:"reference_number"`
	Natl               string      `json:"nationality"`
	Gender             string      `json:"gender"`
	FormType           string      `json:"form_type"`
	FormNumber         string      `json:"form_number"`
	IDIssueDate        string      `json:"id_date_of_issue"`
	IDExpDate          string      `json:"id_expiration_date"`
	IDIssBy            string      `json:"id_issued_by"`
	BBAmt              string      `json:"buy_back_amount"`
	RateCat            string      `json:"rate_category"`
	MCRateID           string      `json:"mc_rate_id"`
	MCRate             string      `json:"mc_rate"`
	OriginatingCountry string      `json:"originating_country"`
	DestinationCountry string      `json:"destination_country"`
	DeviceID           string      `json:"device_id"`
	AgentID            string      `json:"agent_id"`
	AgentCode          string      `json:"agent_code"`
	Currency           string      `json:"currency"`
	ClientReferenceNo  string      `json:"client_reference_no"`
	OrderNo            string      `json:"order_number"`
	DsaCode            string      `json:"dsa_code"`
	DsaTrxType         string      `json:"dsa_trx_type"`
}

type RiaPayoutResponseBody struct {
	Code    json.Number `json:"code"`
	Msg     string      `json:"message"`
	RemcoID json.Number `json:"remco_id"`
}

func (s *Svc) RiaPayout(ctx context.Context, sr RiaPayoutRequest) (*RiaPayoutResponseBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("ria/payout"), sr)
	if err != nil {
		return nil, err
	}

	RiaRes := &RiaPayoutResponseBody{}
	if err := json.Unmarshal(res, RiaRes); err != nil {
		return nil, err
	}
	return RiaRes, nil
}

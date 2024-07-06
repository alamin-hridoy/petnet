package perahub

import (
	"context"
	"encoding/json"
)

type RecMoneyPayRequest struct {
	FrgnRefNo    string `json:"foreign_reference_no"`
	UserCode     string `json:"user_code"`
	CustomerCode string `json:"customer_code"`

	ReceiverNameType   string `json:"receiver_name_type"`
	ReceiverFirstName  string `json:"receiver_first_name"`
	ReceiverMiddleName string `json:"receiver_middle_name"`
	ReceiverLastName   string `json:"receiver_last_name"`

	SenderFirstName string `json:"sender_first_name"`
	SenderLastName  string `json:"sender_last_name"`

	AddrLine1  string `json:"addr_line1"`
	AddrLine2  string `json:"addr_line2"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
	CurrCity   string `json:"curr_city"`
	CurrState  string `json:"curr_state_name"`

	ReceiverAddress1     string `json:"receiver_street"`
	ReceiverCity         string `json:"receiver_city"`
	ReceiverState        string `json:"receiver_state"`
	ReceiverStateZip     string `json:"receiver_state_zip"`
	ReceiverCountryCode  string `json:"receiver_country_code"`
	ReceiverCurrencyCode string `json:"receiver_currency_code"`

	ReceiverHasPhone  string `json:"Does_receiver_have_a_phone_number"`
	RecMobCountryCode string `json:"rec_mob_country_code"`
	PhoneNumber       string `json:"national_number"`
	PhoneCityCode     string `json:"telephone_city_code"`
	ContactPhone      string `json:"contact_phone"`
	Email             string `json:"email"`

	Gender       string `json:"gender"`
	Birthdate    string `json:"date_of_birth"`
	BirthCountry string `json:"Country_of_Birth"`
	Nationality  string `json:"nationality"`

	IDType           string `json:"id_type"`
	IDCountry        string `json:"id_country_of_issue"`
	IDNumber         string `json:"id_number"`
	IDIssueDate      string `json:"id_issue_date"`
	IDHasExpiry      string `json:"Does_the_ID_have_an_expiration_date"`
	IDExpirationDate string `json:"id_expiration_date"`

	Occupation    string `json:"occupation"`
	FundSource    string `json:"Source_of_Funds"`
	EmployerName  string `json:"Name_of_Employer_Business"`
	PositionLevel string `json:"employment_position_level"`

	TxnPurpose           string `json:"Purpose_of_Transaction"`
	ReceiverRelationship string `json:"Relationship_to_Receiver_Sender"`

	PdState            string      `json:"pd_state_code"`
	PdCity             string      `json:"pd_city"`
	PdDestCountry      string      `json:"pd_dest_country_code"`
	PdDestCurrency     string      `json:"pd_dest_currency_code"`
	PdOriginatingCity  string      `json:"pd_originating_city"`
	PdOrigCountryCode  string      `json:"pd_orig_country_code"`
	PdOrigCurrencyCode string      `json:"pd_orig_currency_code"`
	PdTransactionType  string      `json:"pd_transaction_type"`
	PdExchangeRate     json.Number `json:"pd_exchange_rate"`
	PdOrigDestCountry  string      `json:"pd_orig_dest_country_currency_country_code"`
	PdOrigDestCurrency string      `json:"pd_orig_dest_country_currency_currency_code"`

	GrossTotal    json.Number `json:"gross_total_amount"`
	PayAmount     json.Number `json:"pay_amount"`
	Principal     json.Number `json:"principal_amount"`
	Charges       json.Number `json:"charges"`
	Tolls         json.Number `json:"tolls"`
	RealPrincipal string      `json:"real_principal_amount"`
	DstAmount     string      `json:"dst_amount"`
	RealNet       string      `json:"real_net_amount"`

	FilingTime string `json:"filing_time"`
	FilingDate string `json:"filing_date"`

	MoneyTransferKey string `json:"money_transfer_key"`
	Mtcn             string `json:"mtcn"`
	NewMtcn          string `json:"new_mtcn"`
	PayIndicator     string `json:"pay_or_do_not_pay_indicator"`

	Message []string `json:"message"`
	AckFlag string   `json:"ack_flag"`

	TerminalID       string `json:"terminal_id"`
	OperatorID       string `json:"operator_id"`
	RemoteTerminalID string `json:"remote_terminal_id"`
	RemoteOperatorID string `json:"remote_operator_id"`

	GalacticID *string `json:"galactic_id"`
	MyWUNumber string  `json:"mywu_number"`
	MyWUPoints string  `json:"mywu_current_points"`
	HasLoyalty string  `json:"has_loyalty"`
}

type RMPConfirmedDetails struct {
	AdvisoryText    string      `json:"advisory_text"`
	NewPointsEarned json.Number `json:"new_points_earned"`
	PaidDateTime    string      `json:"paid_date_time"`
	HostMessageSet1 string      `json:"host_message_set1"`
	HostMessageSet2 string      `json:"host_message_set2"`
	HostMessageSet3 string      `json:"host_message_set3"`
	PeraCardPoints  string      `json:"pera_card_points"`
}

func (s *Svc) RecMoneyPay(ctx context.Context, pr RecMoneyPayRequest) (*RMPConfirmedDetails, error) {
	pr.ReceiverNameType = "D"
	pr.PayIndicator = "P"
	pr.AckFlag = "X"

	const mod, modReq = "wupo", "pay"
	req, err := s.newParahubRequest(ctx,
		mod, modReq, pr,
		WithUserCode(json.Number(pr.UserCode)),
		WithLocationCode(pr.UserCode))
	if err != nil {
		return nil, err
	}

	res, err := s.post(ctx, s.moduleURL(mod, modReq), *req)
	if err != nil {
		return nil, err
	}

	var smvRes RMPConfirmedDetails
	if err := json.Unmarshal(res, &struct {
		ConfirmedDetails *RMPConfirmedDetails `json:"confirmed_details"`
	}{ConfirmedDetails: &smvRes}); err != nil {
		return nil, err
	}

	return &smvRes, nil
}

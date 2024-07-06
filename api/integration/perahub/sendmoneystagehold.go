package perahub

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
)

type SMSHoldRequest struct {
	ForeignReferenceNo            string            `json:"foreign_reference_no"`
	SenderFirstName               string            `json:"sender_first_name"`
	SenderMiddleName              string            `json:"sender_middle_name"`
	SenderLastName                string            `json:"sender_last_name"`
	SenderNameType                string            `json:"sender_name_type"`
	SenderAddrLine1               string            `json:"sender_addr_line1"`
	SenderAddrLine2               string            `json:"sender_addr_line2"`
	SenderCity                    string            `json:"sender_city"`
	SenderState                   string            `json:"sender_state"`
	SenderPostalCode              string            `json:"sender_postal_code"`
	SenderAddrCountryCode         string            `json:"sender_addr_country_code"`
	SenderAddrCurrencyCode        string            `json:"sender_addr_currency_code"`
	SenderContactPhone            string            `json:"sender_contact_phone"`
	SenderMobileCountryCode       string            `json:"sender_mobile_country_code"`
	SenderMobileNo                string            `json:"sender_mobile_no"`
	MyWUNumber                    string            `json:"mywu_number"`
	MyWUEnrollTag                 string            `json:"mywu_enroll_tag"`
	MyWUCurrentPoints             string            `json:"mywu_current_points"`
	IDType                        string            `json:"id_type"`
	IDCountryOfIssue              string            `json:"id_country_of_issue"`
	IDNumber                      string            `json:"id_number"`
	IDIssueDate                   string            `json:"id_issue_date"`
	IDExpirationDate              string            `json:"id_expiration_date"`
	DateOfBirth                   string            `json:"date_of_birth"`
	Occupation                    string            `json:"occupation"`
	CountryOfBirth                string            `json:"Country_of_Birth"`
	Nationality                   string            `json:"nationality"`
	SourceOfFunds                 string            `json:"Source_of_Funds"`
	SenderEmployeer               string            `json:"Sender_Employeer"`
	RelationshipToReceiver        string            `json:"Relationship_to_Receiver"`
	Gender                        string            `json:"Gender"`
	ReasonForSend                 string            `json:"reason_for_send"`
	Email                         string            `json:"email"`
	ReceiverFirstName             string            `json:"receiver_first_name"`
	ReceiverMiddleName            string            `json:"receiver_middle_name"`
	ReceiverLastName              string            `json:"receiver_last_name"`
	ReceiverNameType              string            `json:"receiver_name_type"`
	ReceiverAddrLine1             string            `json:"receiver_addr_line1"`
	ReceiverAddrLine2             string            `json:"receiver_addr_line2"`
	ReceiverCity                  string            `json:"receiver_city"`
	ReceiverState                 string            `json:"receiver_state"`
	ReceiverPostalCode            string            `json:"receiver_postal_code"`
	ReceiverAddrCountryCode       string            `json:"receiver_addr_country_code"`
	ReceiverAddrCurrencyCode      string            `json:"receiver_addr_currency_code"`
	ReceiverContactPhone          string            `json:"receiver_contact_phone"`
	ReceiverMobileCountryCode     string            `json:"receiver_mobile_country_code"`
	ReceiverMobileNo              string            `json:"receiver_mobile_no"`
	DestinationCountryCode        string            `json:"destination_country_code"`
	DestinationCurrencyCode       string            `json:"destination_currency_code"`
	DestinationCityName           string            `json:"destination_city_name"`
	DestinationStateCode          string            `json:"destination_state_code"`
	TransactionType               string            `json:"transaction_type"`
	PrincipalAmount               string            `json:"principal_amount"`
	ExchangeRate                  string            `json:"exchange_rate"`
	FixedAmountFlag               string            `json:"fixed_amount_flag"`
	PromoCode                     string            `json:"promo_code"`
	Message                       []string          `json:"message"`
	BankName                      string            `json:"bank_name"`
	AccountNumber                 string            `json:"account_number"`
	BankCode                      string            `json:"bank_code"`
	BankLocation                  string            `json:"bank_location"`
	AccountSuffix                 string            `json:"account_suffix"`
	OTPPin                        string            `json:"otp_pin"`
	OTPCode                       string            `json:"otp_code"`
	TerminalID                    string            `json:"terminal_id"`
	OperatorID                    string            `json:"operator_id"`
	RemoteTerminalID              string            `json:"remote_terminal_id"`
	RemoteOperatorID              string            `json:"remote_operator_id"`
	AddlServiceCharges            string            `json:"addl_service_charges"`
	StagingBuffer                 string            `json:"staging_buffer"`
	ComplianceDataBuffer          string            `json:"compliance_data_buffer"`
	OriginatingCity               string            `json:"originating_city"`
	OriginatingState              string            `json:"originating_state"`
	Promotions                    Promotions        `json:"promotions"`
	Financials                    Financials        `json:"financials"`
	MTCN                          string            `json:"mtcn"`
	NewMTCN                       string            `json:"new_mtcn"`
	TestQuestion                  string            `json:"test_question"`
	Answer                        string            `json:"answer"`
	HasLoyalty                    string            `json:"has_loyalty"`
	SecondIDType                  string            `json:"second_id_type"`
	SecondIDNumber                string            `json:"second_id_number"`
	SecondIDCountryOfIssue        string            `json:"second_id_country_of_issue"`
	SecondIDIssueDate             string            `json:"second_id_issue_date"`
	SecondIDExpirationDate        string            `json:"second_id_expiration_date"`
	ThirdIDType                   string            `json:"third_id_type"`
	ThirdIDNumber                 string            `json:"third_id_number"`
	ThirdIDCountryOfIssue         string            `json:"third_id_country_of_issue"`
	ThirdIDIssueDate              string            `json:"third_id_issue_date"`
	ThirdIDExpirationDate         string            `json:"third_id_expiration_date"`
	EmploymentPositionLevel       string            `json:"employment_position_level"`
	PurposeOfTransaction          string            `json:"Purpose_of_Transaction"`
	IsCurrentAndPermanentAddrSame string            `json:"is_current_and_permanent_addr_same"`
	PermaAddrLine1                string            `json:"perma_addr_line1"`
	PermaAddrLine2                string            `json:"perma_addr_line2"`
	PermaCity                     string            `json:"perma_city"`
	PermaStateName                string            `json:"perma_state_name"`
	PermaPostalCode               string            `json:"perma_postal_code"`
	PermaCountry                  string            `json:"perma_country"`
	IsOnBehalf                    string            `json:"is_on_behalf"`
	SenderAddrCountryName         string            `json:"sender_addr_country_name"`
	ComplianceDetails             ComplianceDetails `json:"compliance_details"`
	CustomerCode                  string            `json:"customer_code"`
	ClientReferenceNo             string            `json:"client_reference_no"`
	UserCode                      string            `json:"user_code"`
}

type StagedDetails struct {
	AdvisoryText      string `json:"advisory_text"`
	MTCN              string `json:"mtcn"`
	NewMTCN           string `json:"new_mtcn"`
	FilingDate        string `json:"filing_date"`
	FilingTime        string `json:"filing_time"`
	MTRequestedStatus string `json:"mt_requested_status"`
}

type SMSHoldResponseBody struct {
	StagedDetails StagedDetails `json:"staged_details"`
}

type SMSHoldResponseWU struct {
	Header ResponseHeader      `json:"header"`
	Body   SMSHoldResponseBody `json:"body"`
}

type SMSHoldResponse struct {
	WU SMSHoldResponseWU `json:"uspwuapi"`
}

func (s *Svc) SendMoneyStageHold(ctx context.Context, smshReq SMSHoldRequest) (*SMSHoldResponse, error) {
	req, err := s.newParahubRequest(ctx, "wusostg", "SMSstore", smshReq)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := s.cl.Post(s.moduleURL("wusostg", ""), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var smshRes SMSHoldResponse
	if err := json.Unmarshal(body, &smshRes); err != nil {
		return nil, err
	}

	return &smshRes, nil
}

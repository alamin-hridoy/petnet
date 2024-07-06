package perahub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testSMSVBody struct {
	Module  string      `json:"module"`
	Request string      `json:"request"`
	Param   SMSVRequest `json:"param"`
}

type testSMSVWU struct {
	Header    RequestHeader `json:"header"`
	Body      testSMSVBody  `json:"body"`
	Signature string        `json:"signature"`
}

type testSMSVRequest struct {
	WU testSMSVWU `json:"uspwuapi"`
}

var smsvReq = SMSVRequest{
	ForeignReferenceNo:            "ZMp4gJ3vXmjNttvKsYwR",
	ReceiverCompany:               "",
	ReceiverAttention:             "",
	StagingBuffer:                 "0101C0202MB0303PIL",
	TestQue:                       "",
	Answer:                        "",
	SenderNameType:                "D",
	SenderFirstName:               "WINNIE",
	SenderMiddleName:              "",
	SenderLastName:                "CONSTANTINO",
	SenderAddrLine1:               "1953 PH 3B block 6 lot 9 Camarin",
	SenderAddrLine2:               "Ascension street",
	SenderCity:                    "Caloocan city",
	SenderState:                   "Metro Manila",
	SenderPostalCode:              "1400",
	SenderAddrCountryCode:         "PH",
	SenderAddrCurrencyCode:        "PHP",
	SenderContactPhone:            "9089488895",
	SenderMobileCountryCode:       "63",
	SenderMobileNo:                "9089488895",
	SenderAddrCountryName:         "PHILIPPINES",
	MyWUNumber:                    "185015541",
	IDType:                        "A",
	IDCountryOfIssue:              "PHILIPPINES",
	IDNumber:                      "987654321",
	IDIssueDate:                   "10082014",
	IDExpirationDate:              "12012019",
	DateOfBirth:                   "12011996",
	Occupation:                    "Airline/Maritime Employee",
	CountryOfBirth:                "PHILIPPINES",
	Nationality:                   "PHILIPPINES",
	Gender:                        "F",
	SourceOfFunds:                 "Gift",
	SenderEmployeer:               "Vox Dei Protocol System",
	RelationshipToReceiver:        "Donor/Receiver of Ch",
	GENIIIIndicator:               "",
	AckFlag:                       "X",
	ReasonForSend:                 "",
	MyWUEnrollTag:                 "mywu",
	Email:                         "winnieconstantino@gmail.com",
	ReceiverNameType:              "D",
	ReceiverFirstName:             "JERICA",
	ReceiverMiddleName:            "",
	ReceiverLastName:              "NAPARATE",
	ReceiverAddrLine1:             "Bulacan city",
	ReceiverAddrLine2:             "",
	ReceiverCity:                  "Bulacan city",
	ReceiverState:                 "Bulacan city",
	ReceiverPostalCode:            "",
	ReceiverAddrCountryCode:       "PH",
	ReceiverAddrCurrencyCode:      "PHP",
	ReceiverContactPhone:          "",
	ReceiverMobileCountryCode:     "",
	ReceiverMobileNo:              "",
	DestinationCountryCode:        "PH",
	DestinationCurrencyCode:       "PHP",
	DestinationStateCode:          "",
	DestinationCityName:           "",
	TransactionType:               "SO",
	PrincipalAmount:               "100000",
	FixedAmountFlag:               "N",
	PromoCode:                     "",
	Message:                       []string{},
	BankName:                      "",
	AccountNumber:                 "",
	BankCode:                      "",
	BankLocation:                  "",
	AccountSuffix:                 "",
	TerminalID:                    "PH259ART001A",
	OperatorID:                    "001",
	RemoteTerminalID:              "",
	RemoteOperatorID:              "",
	SecondIDType:                  "",
	SecondIDNumber:                "",
	SecondIDCountryOfIssue:        "",
	SecondIDIssueDate:             "",
	SecondIDExpirationDate:        "",
	ThirdIDType:                   "",
	ThirdIDNumber:                 "",
	ThirdIDCountryOfIssue:         "",
	ThirdIDIssueDate:              "",
	ThirdIDExpirationDate:         "",
	EmploymentPositionLevel:       "Mid-Level",
	PurposeOfTransaction:          "Gift",
	IsCurrentAndPermanentAddrSame: "Y",
	PermaAddrLine1:                "1953 PH 3B BLOCK 6 LOT 9 CAMARIN",
	PermaAddrLine2:                "ASCENSION STREET",
	PermaCity:                     "CALOOCAN CITY",
	PermaStateName:                "METRO MANILA",
	PermaPostalCode:               "1400",
	PermaCountry:                  "PHILIPPINES",
	IsOnBehalf:                    "N",
	GalacticID:                    StrPtr("1000000000028668380"),
}

func TestSendMoneyStageValidate(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          SMSVRequest
		expectedReq testSMSVRequest
		want        *SMSVResponse
		wantErr     bool
	}{
		{
			name: "Success",
			in:   smsvReq,
			expectedReq: testSMSVRequest{
				WU: testSMSVWU{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testSMSVBody{
						Module:  "wusostg",
						Request: "SMSvalidate",
						Param:   smsvReq,
					},
					Signature: "placeholderSignature",
				},
			},
			want: &SMSVResponse{
				WU: SMSVResponseWU{
					Header: ResponseHeader{
						ErrorCode: "1",
						Message:   "good",
					},
					Body: SMSVResponseBody{
						ServiceCode: VServiceCode{
							AddlServiceCharges: "1103FEE12041500130415000103MSG0201003010210300022010230109001096210104500002021003035009703NNN9806PILPIL",
						},
						Compliance: SCompliance{
							ComplianceDataBuffer: "0108UNI_01_S0201A03099876543210611PHILIPPINES0708100819960825Airline Maritime Employee0908100820141008100820192401N321090894888959901X3311PHILIPPINES3411PHILIPPINES 4404Gift2104Gift6723Vox Dei Protocol System7020Donor Receiver of Ch7401Y8201NC1109089488895C4026357321953 PH 3B BLOCK 6 LOT 9 CAMARIN5816ASCENSION STREET5913CALOOCAN CITY6012METRO MANILA610414006211PHILIPPINESF901YJ6191000000000028668380J701YM20263M732Mid-Level Supervisory Management",
						},
						PaymentDetails: SPaymentDetails{
							OriginatingCity:  "METRO MANILA",
							OriginatingState: "MA",
							StagingBuffer:    "0101C0202MB0303PIL0820CCUANA00WUCCMBAGCA",
						},
						Financials: SFinancials{
							Taxes: STaxes{
								MunicipalTax: "0",
								StateTax:     "0",
								CountyTax:    "0",
							},
							OriginatorsPrincipalAmount: "100000",
							DestinationPrincipalAmount: "100000",
							GrossTotalAmount:           "101500",
							PlusChargesAmount:          "0",
							Charges:                    "1500",
							MessageCharge:              "0",
							TUCharges:                  "1500",
							TDCharges:                  "1500",
						},
						Promotions: SPromotions{
							PromoCodeDescription: "",
							PromoMessage:         "",
							SenderPromoCode:      "",
						},
						NewDetails: SNewDetails{
							MTCN:       "8851333182",
							NewMTCN:    "1809588851333182",
							FilingDate: "04-05",
							FilingTime: "1240A EDT",
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(successSMSVHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.WU.Header.Token

			got, err := s.SendMoneyStageValidate(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("SendMoneyStageValidate() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if got.WU.Header.ErrorCode != "1" {
				t.Errorf("want error code 1 but got %q", got.WU.Header.ErrorCode)
			}

			if !cmp.Equal(test.want, got) {
				t.Fatal(cmp.Diff(test.want, got))
			}
		})
	}
}

func successSMSVHandler(t *testing.T, expectedReq testSMSVRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/wusostg" {
			t.Errorf("expected request to '/wusostg', got '%s'", req.URL.EscapedPath())
		}

		var newReq testSMSVRequest
		if err := json.NewDecoder(req.Body).Decode(&newReq); err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(expectedReq, newReq) {
			t.Error(cmp.Diff(expectedReq, newReq))
		}

		res.WriteHeader(200)
		res.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(res, `{
			"uspwuapi":{
				"header":{
					"errorcode":"1",
					"message":"good"
				},
				"body":{
					"service_code": {
						"addl_service_charges": "1103FEE12041500130415000103MSG0201003010210300022010230109001096210104500002021003035009703NNN9806PILPIL"
					},
					"compliance": {
						"compliance_data_buffer": "0108UNI_01_S0201A03099876543210611PHILIPPINES0708100819960825Airline Maritime Employee0908100820141008100820192401N321090894888959901X3311PHILIPPINES3411PHILIPPINES 4404Gift2104Gift6723Vox Dei Protocol System7020Donor Receiver of Ch7401Y8201NC1109089488895C4026357321953 PH 3B BLOCK 6 LOT 9 CAMARIN5816ASCENSION STREET5913CALOOCAN CITY6012METRO MANILA610414006211PHILIPPINESF901YJ6191000000000028668380J701YM20263M732Mid-Level Supervisory Management"
					},
					"payment_details": {
						"originating_city": "METRO MANILA",
						"originating_state": "MA",
						"staging_buffer": "0101C0202MB0303PIL0820CCUANA00WUCCMBAGCA"
					},
					"financials": {
						"taxes": {
							"municipal_tax": 0,
							"state_tax": 0,
							"county_tax": 0
						},
						"originators_principal_amount": 100000,
						"destination_principal_amount": 100000,
						"gross_total_amount": 101500,
						"plus_charges_amount": 0,
						"charges": 1500,
						"message_charge": 0,
						"total_undiscounted_charges": 1500,
						"total_discounted_charges": 1500
					},
					"promotions": {
						"promo_code_description": "",
						"promo_message": "",
						"sender_promo_code": ""
					},
					"new_details": {
						"mtcn": "8851333182",
						"new_mtcn": "1809588851333182",
						"filing_date": "04-05",
						"filing_time": "1240A EDT"
					}
				}
			}
		}`)
	}
}

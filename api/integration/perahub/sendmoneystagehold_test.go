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

type testSMSHRequestBody struct {
	Module  string         `json:"module"`
	Request string         `json:"request"`
	Param   SMSHoldRequest `json:"param"`
}

type testSMSHRequestWU struct {
	Header    RequestHeader       `json:"header"`
	Body      testSMSHRequestBody `json:"body"`
	Signature string              `json:"signature"`
}

type testSMSHRequest struct {
	WU testSMSHRequestWU `json:"uspwuapi"`
}

var smshReq = SMSHoldRequest{
	ForeignReferenceNo:            "cdffb28cea04d8ee9cb1",
	SenderNameType:                "D",
	SenderFirstName:               "WINNIE",
	SenderMiddleName:              "",
	SenderLastName:                "CONSTANTINO",
	SenderAddrLine1:               "1953 PH 3B BLOCK 6 LOT 9 CAMARIN",
	SenderAddrLine2:               "",
	SenderCity:                    "CALOOCAN CITY",
	SenderState:                   "METRO MANILA",
	SenderPostalCode:              "1400",
	SenderAddrCountryCode:         "PH",
	SenderAddrCurrencyCode:        "PHP",
	SenderContactPhone:            "9089488895",
	SenderMobileCountryCode:       "63",
	SenderMobileNo:                "9089488895",
	SenderAddrCountryName:         "Philippines",
	MyWUNumber:                    "185015541",
	IDType:                        "A",
	IDCountryOfIssue:              "Philippines",
	IDNumber:                      "987654321",
	IDIssueDate:                   "12012012",
	IDExpirationDate:              "12012019",
	DateOfBirth:                   "12011996",
	Occupation:                    "Airline Maritime Employee",
	CountryOfBirth:                "Philippines",
	Nationality:                   "Philippines",
	Gender:                        "",
	SourceOfFunds:                 "Salary Income",
	SenderEmployeer:               "VOX DEI PROTOCOL SYSTEMS",
	RelationshipToReceiver:        "",
	ReasonForSend:                 "",
	MyWUEnrollTag:                 "mywu",
	Email:                         "WINNIEANN.CONSTANTINO@VOXDEISYSTEMS.COM",
	ReceiverNameType:              "D",
	ReceiverFirstName:             "JERICA",
	ReceiverMiddleName:            "",
	ReceiverLastName:              "NAPARATE",
	ReceiverAddrLine1:             "",
	ReceiverAddrLine2:             "",
	ReceiverCity:                  "",
	ReceiverState:                 "",
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
	PrincipalAmount:               "600000",
	FixedAmountFlag:               "N",
	PromoCode:                     "",
	Message:                       []string{},
	BankName:                      "",
	AccountNumber:                 "",
	BankCode:                      "",
	BankLocation:                  "",
	AccountSuffix:                 "",
	TerminalID:                    "PH259ART001A",
	OperatorID:                    "5",
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
	EmploymentPositionLevel:       "Entry Level",
	PurposeOfTransaction:          "",
	IsCurrentAndPermanentAddrSame: "N",
	PermaAddrLine1:                "1953 PH 3B BLOCK 6 LOT 9 CAMARIN",
	PermaAddrLine2:                "",
	PermaCity:                     "CALOOCAN CITY",
	PermaStateName:                "METRO MANILA",
	PermaPostalCode:               "1400",
	PermaCountry:                  "Philippines",
	IsOnBehalf:                    "",
	AddlServiceCharges:            "1103FEE1205115001305115000103MSG020100301021030002201023009001096210104500002021003035009703NNN9806PILPIL",
	ComplianceDataBuffer:          "0108UNI_01_S0201A03099876543210611Philippines0708120119960825AirlineMaritimeEmployee0908120120121008120120199901X3311Philippines3411Philippines4413SalaryIncome6724VOX DEI PROTOCOL SYSTEMS7401YB401157321953 PH 3B BLOCK 6 LOT 9CAMARIN5913CALOOCAN CITY6012METROMANILA610414006211PhilippinesF901NJ6191000000000028668380J701YM711EntryLevel",
	OriginatingCity:               "METRO MANILA",
	OriginatingState:              "MA",
	Financials: Financials{
		Taxes: Taxes{
			// Currency:  "PHP",
			MuniTax:   "0",
			StateTax:  "0",
			CountyTax: "0",
			// Total:     "0",
		},
		OrigPnplAmt:   "600000",
		DestPcplAmt:   "600000",
		GrossTotal:    "611500",
		AddlCharges:   "0",
		Charges:       "11500",
		MessageCharge: "0",
	},
	Promotions: Promotions{
		PromoDesc:       "",
		PromoMessage:    "",
		SenderPromoCode: "",
	},
	MTCN:              "3010897610",
	NewMTCN:           "1809483010897610",
	ExchangeRate:      "1.0000000",
	OTPPin:            "",
	OTPCode:           "",
	UserCode:          "000005",
	MyWUCurrentPoints: "1670",
	StagingBuffer:     "0101C0202MB0303PIL0820CCUANA00WUCCMBAGCA",
	TestQuestion:      "",
	Answer:            "",
	HasLoyalty:        "0",
	ComplianceDetails: ComplianceDetails{
		IDDetails: IDDetails{
			IDType:    "A",
			IDCountry: "Philippines",
			IDNumber:  "987654321",
		},
		IDIssued:          "12012012",
		IDExpiry:          "12012019",
		Birthdate:         "12011996",
		Occupation:        "Airline Maritime Employee",
		TransactionReason: "",
		CurrentAddress: Address{
			AddrLine1:  "1953 PH 3B BLOCK 6 LOT 9 CAMARIN",
			AddrLine2:  "",
			City:       "CALOOCAN CITY",
			State:      "METRO MANILA",
			PostalCode: "1400",
			Country:    "Philippines",
		},
		TxnRelationship:         "",
		IActOnMyBehalf:          "",
		EmploymentPositionLevel: "Entry Level",
	},
	CustomerCode:      "WC1201962A0A",
	ClientReferenceNo: "PH259ART001A201804050011",
}

func TestSendMoneyStageHold(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          SMSHoldRequest
		expectedReq testSMSHRequest
		want        *SMSHoldResponse
		wantErr     bool
	}{
		{
			name: "Success",
			in:   smshReq,
			expectedReq: testSMSHRequest{
				WU: testSMSHRequestWU{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testSMSHRequestBody{
						Module:  "wusostg",
						Request: "SMSstore",
						Param:   smshReq,
					},
					Signature: "placeholderSignature",
				},
			},
			want: &SMSHoldResponse{
				WU: SMSHoldResponseWU{
					Header: ResponseHeader{
						ErrorCode: "1",
						Message:   "good",
					},
					Body: SMSHoldResponseBody{
						StagedDetails: StagedDetails{
							AdvisoryText:      "04/05/18 00:41",
							MTCN:              "9641478070",
							NewMTCN:           "1809589641478070",
							FilingDate:        "04-05",
							FilingTime:        "1241A EDT ",
							MTRequestedStatus: "HOLD",
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
			ts := httptest.NewServer(successSMSHoldHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.WU.Header.Token

			got, err := s.SendMoneyStageHold(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatal(err)
				t.Errorf("SendMoneyStore() error = %v, wantErr %v", err, test.wantErr)
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

func successSMSHoldHandler(t *testing.T, expectedReq testSMSHRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/wusostg" {
			t.Errorf("expected request to '/wusostg', got '%s'", req.URL.EscapedPath())
		}

		var newReq testSMSHRequest
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
					"staged_details":{
						"advisory_text":"04/05/18 00:41",
						"mtcn":"9641478070",
						"new_mtcn":"1809589641478070",
						"filing_date":"04-05",
						"filing_time":"1241A EDT ",
						"mt_requested_status":"HOLD"
					}
				}
			}
		}`)
	}
}

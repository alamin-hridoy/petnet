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

type testSMSBody struct {
	Module  string         `json:"module"`
	Request string         `json:"request"`
	Param   SMStoreRequest `json:"param"`
}

type testSMSWU struct {
	Header    RequestHeader `json:"header"`
	Body      testSMSBody   `json:"body"`
	Signature string        `json:"signature"`
}

type testSMSRequest struct {
	WU testSMSWU `json:"uspwuapi"`
}

var smsReq = SMStoreRequest{
	FrgnRefNo:                 "cdffb28cea04d8ee9cb1",
	UserCode:                  userCode,
	CustomerCode:              "WC1201962A0A",
	SenderNameType:            "D",
	SenderFirstName:           "WINNIE",
	SenderMiddleName:          "",
	SenderLastName:            "CONSTANTINO",
	SenderAddrCountryCode:     "PH",
	SenderAddrCurrencyCode:    "PHP",
	SenderContactPhone:        "9089488895",
	SenderMobileCountryCode:   "63",
	SenderMobileNo:            "9089488895",
	SenderReason:              "",
	Email:                     "WINNIEANN.CONSTANTINO@VOXDEISYSTEMS.COM",
	ReceiverNameType:          "D",
	ReceiverFirstName:         "JERICA",
	ReceiverMiddleName:        "",
	ReceiverLastName:          "NAPARATE",
	ReceiverAddrLine1:         "",
	ReceiverAddrLine2:         "",
	ReceiverCity:              "",
	ReceiverState:             "",
	ReceiverPostalCode:        "",
	ReceiverAddrCountryCode:   "PH",
	ReceiverAddrCurrencyCode:  "PHP",
	ReceiverContactPhone:      "",
	ReceiverMobileCountryCode: "",
	ReceiverMobileNo:          "",
	DestinationCountryCode:    "PH",
	DestinationCurrencyCode:   "PHP",
	DestinationStateCode:      "",
	DestinationCityName:       "",
	TransactionType:           "SO",
	PrincipalAmount:           "600000",
	FixedAmountFlag:           "N",
	PromoCode:                 "",
	Message:                   []string{},
	AddlServiceChg:            "1103FEE1205115001305115000103MSG020100301021030002201023009001096210104500002021003035009703NNN9806PILPIL",
	ComplianceDataBuffer:      "0108UNI_01_S0201A03099876543210611Philippines0708120119960825AirlineMaritimeEmployee0908120120121008120120199901X3311Philippines3411Philippines4413SalaryIncome6724VOX DEI PROTOCOL SYSTEMS7401YB401157321953 PH 3B BLOCK 6 LOT 9CAMARIN5913CALOOCAN CITY6012METROMANILA610414006211PhilippinesF901NJ6191000000000028668380J701YM711EntryLevel",
	OriginatingCity:           "METRO MANILA",
	OriginatingState:          "MA",
	MTCN:                      "3010897610",
	NewMTCN:                   "1809483010897610",
	ExchangeRate:              "1.0000000",
	Fin: Financials{
		Taxes: Taxes{
			MuniTax:   "0",
			StateTax:  "0",
			CountyTax: "0",
		},
		OrigPnplAmt:   "600000",
		DestPcplAmt:   "600000",
		GrossTotal:    "611500",
		AddlCharges:   "0",
		Charges:       "11500",
		MessageCharge: "0",
	},
	Promo: Promotions{
		PromoDesc:       "",
		PromoMessage:    "",
		SenderPromoCode: "",
	},
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
	BankName:         "",
	BankLocation:     "",
	AccountNumber:    "",
	BankCode:         "",
	AccountSuffix:    "",
	MyWUNumber:       "185015541",
	WUEnroll:         "none",
	TerminalID:       "PH259ART001A",
	OperatorID:       "5",
	RemoteTerminalID: "",
	RemoteOperatorID: "",
}

func TestSendMoneyStore(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          SMStoreRequest
		expectedReq testSMSRequest
		want        *ConfirmedDetails
		wantErr     bool
	}{
		{
			name: "Success",
			in:   smsReq,
			expectedReq: testSMSRequest{
				WU: testSMSWU{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: userCode,
						UserCode:     json.Number(userCode),
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testSMSBody{
						Module:  "wuso",
						Request: "SMstore",
						Param:   smsReq,
					},
					Signature: "placeholderSignature",
				},
			},
			want: &ConfirmedDetails{
				AdvisoryText: "04 04 18 22:57",
				MTCN:         "3010897610",
				NewMTCN:      "1809483010897610",
				FilingDate:   "04-04",
				FilingTime:   "1057P EDT ",
				PinMessage: []string{
					"Thank you for joining My WU program.",
					"Watch out for exciting benefits & exclusive member offers!",
					"Coming soon! Visit www.westernunion.com.ph for more information.",
				},
				PromoTextMessage: []string{
					"Mahalaga sa amin ang iyong opinyon!",
					"Pumunta sa westernunion.com nakikinig at ipaalam sa amin ang iyong",
					"masasabi tungkol sa aming serbisyo.",
				},
				MyWUNumber:      "185015541",
				NewPointsEarned: 30,
				OtherMessage1:   []string{"1", "2", "3"},
				OtherMessage2:   []string{"4", "5", "6"},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(successSMSHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.WU.Header.Token

			got, err := s.SendMoneyStore(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatal(err)
				t.Errorf("SendMoneyStore() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if !cmp.Equal(test.want, got) {
				t.Fatal(cmp.Diff(test.want, got))
			}
		})
	}
}

func successSMSHandler(t *testing.T, expectedReq testSMSRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/1.1/wuso-store" {
			t.Errorf("expected request to '/1.1/wuso-store', got '%s'", req.URL.EscapedPath())
		}

		var newReq testSMSRequest
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
					"confirmed_details":{
						"advisory_text":"04 04 18 22:57",
						"mtcn":"3010897610",
						"new_mtcn":"1809483010897610",
						"filing_date":"04-04",
						"filing_time":"1057P EDT ",
						"pin_message":[
							"Thank you for joining My WU program.",
							"Watch out for exciting benefits & exclusive member offers!",
							"Coming soon! Visit www.westernunion.com.ph for more information."
						],
						"promo_text_message":[
							"Mahalaga sa amin ang iyong opinyon!",
							"Pumunta sa westernunion.com nakikinig at ipaalam sa amin ang iyong",
							"masasabi tungkol sa aming serbisyo."
						],
						"mywu_number":"185015541",
						"new_points_earned":30,
						"other_message_1":[
							"1",
							"2",
							"3"
						],
						"other_message_2":[
							"4",
							"5",
							"6"
						]
					}
				}
			}
		}`)
	}
}

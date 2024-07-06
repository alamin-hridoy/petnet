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

type testRecMoneyPayBody struct {
	Module  string             `json:"module"`
	Request string             `json:"request"`
	Param   RecMoneyPayRequest `json:"param"`
}

type testRecMoneyPayData struct {
	Header    RequestHeader       `json:"header"`
	Body      testRecMoneyPayBody `json:"body"`
	Signature string              `json:"signature"`
}

type testRecMoneyPayRequest struct {
	Data testRecMoneyPayData `json:"uspwuapi"`
}

var reqRMPData = RecMoneyPayRequest{
	FrgnRefNo:            "13b6945f956cc3922ad8",
	UserCode:             userCode,
	CustomerCode:         "WC1201962A0A",
	ReceiverNameType:     "D",
	ReceiverFirstName:    "WINNIE",
	ReceiverMiddleName:   "",
	ReceiverLastName:     "CONSTANTINO",
	SenderFirstName:      "MICHAELLA CHI",
	SenderLastName:       "VEGA",
	AddrLine1:            "1953 PH 3B BLOCK 6 LOT 9 CAMARIN",
	AddrLine2:            "",
	PostalCode:           "1400",
	Country:              "PHILIPPINES",
	CurrCity:             "CALOOCAN CITY",
	CurrState:            "METRO MANILA",
	ReceiverAddress1:     "",
	ReceiverCity:         "CALOOCAN CITY",
	ReceiverState:        "METRO MANILA",
	ReceiverStateZip:     "1400",
	ReceiverCountryCode:  "PH",
	ReceiverCurrencyCode: "PHP",
	ReceiverHasPhone:     "N",
	RecMobCountryCode:    "63",
	PhoneNumber:          "639089488895",
	PhoneCityCode:        "",
	ContactPhone:         "639089488895",
	Email:                "",
	Gender:               "",
	Birthdate:            "12011992",
	BirthCountry:         "PHILIPPINES",
	Nationality:          "PHILIPPINES",
	IDType:               "A",
	IDCountry:            "PHILIPPINES",
	IDNumber:             "987654321",
	IDIssueDate:          "12012016",
	IDHasExpiry:          "Y",
	IDExpirationDate:     "12012019",
	Occupation:           "AirlineMaritime Employee",
	FundSource:           "",
	EmployerName:         "VOX DEI PROTOCOL SYSTEMS",
	PositionLevel:        "",
	TxnPurpose:           "",
	ReceiverRelationship: "",
	PdState:              "",
	PdCity:               "",
	PdDestCountry:        "PH",
	PdDestCurrency:       "PHP",
	PdOriginatingCity:    "METRO MANILAPI1",
	PdOrigCountryCode:    "PH",
	PdOrigCurrencyCode:   "PHP",
	PdTransactionType:    "WMF",
	PdExchangeRate:       "1",
	PdOrigDestCountry:    "PH",
	PdOrigDestCurrency:   "PHP",
	GrossTotal:           "304500",
	PayAmount:            "300000",
	Principal:            "300000",
	Charges:              "4500",
	Tolls:                "0",
	RealPrincipal:        "3000",
	DstAmount:            "0",
	RealNet:              "3000",
	FilingTime:           "0330A EDT",
	FilingDate:           "04-04-18 ",
	MoneyTransferKey:     "3624030255",
	Mtcn:                 "9370896950",
	NewMtcn:              "1809489370896950",
	PayIndicator:         "P",
	Message:              []string{},
	AckFlag:              "X",
	TerminalID:           "PH259ART001A",
	OperatorID:           "5",
	RemoteTerminalID:     "",
	RemoteOperatorID:     "",
	GalacticID:           StrPtr("1000000000028668380"),
	MyWUNumber:           "185015541",
	MyWUPoints:           "",
	HasLoyalty:           "",
}

func TestRecMoneyPayword(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          RecMoneyPayRequest
		expectedReq testRecMoneyPayRequest
		want        *RMPConfirmedDetails
		wantErr     bool
	}{
		{
			name: "Success",
			in:   reqRMPData,
			expectedReq: testRecMoneyPayRequest{
				Data: testRecMoneyPayData{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: userCode,
						UserCode:     json.Number(userCode),
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testRecMoneyPayBody{
						Module:  "wupo",
						Request: "pay",
						Param:   reqRMPData,
					},
					Signature: "placeholderSignature",
				},
			},
			want: &RMPConfirmedDetails{
				AdvisoryText:    "",
				NewPointsEarned: "10",
				PaidDateTime:    "04-04-2018 03:33",
				HostMessageSet1: "",
				HostMessageSet2: "Thank you for joining My WU program.Watch out for exciting benefits & exclusive member offers!Coming soon! Visit www.westernunion.com.ph for more information.",
				HostMessageSet3: "",
				PeraCardPoints:  "0",
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(RecMoneyPayHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.Data.Header.Token

			got, err := s.RecMoneyPay(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("RecMoneyPay() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if !cmp.Equal(test.want, got) {
				t.Fatal(cmp.Diff(test.want, got))
			}
		})
	}
}

func RecMoneyPayHandler(t *testing.T, expectedReq testRecMoneyPayRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/1.1/wupo-pay" {
			t.Errorf("expected request to '/1.1/wupo-pay', got '%s'", req.URL.EscapedPath())
		}

		var newReq testRecMoneyPayRequest
		if err := json.NewDecoder(req.Body).Decode(&newReq); err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(expectedReq, newReq) {
			t.Error(cmp.Diff(expectedReq, newReq))
		}

		res.WriteHeader(200)
		res.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(res, `{
			"uspwuapi": {
				"header": {
					"errorcode": "1",
					"message": "Good"
				},
				"body": {
					"confirmed_details": {
						"advisory_text": "",
						"new_points_earned": "10",
						"paid_date_time": "04-04-2018 03:33",
						"host_message_set1": "",
						"host_message_set2": "Thank you for joining My WU program.Watch out for exciting benefits & exclusive member offers!Coming soon! Visit www.westernunion.com.ph for more information.",
						"host_message_set3": "",
						"pera_card_points": "0"
					}
				}
			}
		}`)
	}
}

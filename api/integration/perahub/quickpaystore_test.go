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

type testQPSRequestBody struct {
	Module  string     `json:"module"`
	Request string     `json:"request"`
	Param   QPSRequest `json:"param"`
}

type testQPSRequestWU struct {
	Header    RequestHeader      `json:"header"`
	Body      testQPSRequestBody `json:"body"`
	Signature string             `json:"signature"`
}

type testQPSRequest struct {
	WU testQPSRequestWU `json:"uspwuapi"`
}

var (
	userCode = "5"
	QPSReq   = QPSRequest{
		ClientReferenceNo:       "PH259ART001A201804050011",
		FrgnRefNo:               "cdffb28cea04d8ee9cb1",
		UserCode:                userCode,
		CustomerCode:            "WC1201962A0A",
		SenderNameType:          "D",
		SenderFirstName:         "WINNIE",
		SenderMiddleName:        "",
		SenderLastName:          "CONSTANTINO",
		SenderAddrCountryCode:   "PH",
		SenderAddrCurrencyCode:  "PHP",
		SenderContactPhone:      "9089488895",
		SenderMobileCountryCode: "63",
		SenderMobileNo:          "9089488895",
		SendingReason:           "",
		Email:                   "WINNIEANN.CONSTANTINO@VOXDEISYSTEMS.COM",
		CompanyName:             "TESTING ONLY",
		CompanyCode:             "1A IA",
		CompanyAccountCode:      "121212121212A",
		ReferenceNo:             "",
		DestCountry:             "PH",
		DestCurrency:            "PHP",
		DestState:               "",
		DestCity:                "",
		TransactionType:         "SO",
		PrincipalAmount:         "600000",
		FixedAmountFlag:         "N",
		PromoCode:               "",
		Message:                 []string{},
		AddlServiceCharges:      "1103FEE1205115001305115000103MSG020100301021030002201023009001096210104500002021003035009703NNN9806PILPIL",
		ComplianceDataBuffer:    "0108UNI_01_S0201A03099876543210611Philippines0708120119960825AirlineMaritimeEmployee0908120120121008120120199901X3311Philippines3411Philippines4413SalaryIncome6724VOX DEI PROTOCOL SYSTEMS7401YB401157321953 PH 3B BLOCK 6 LOT 9CAMARIN5913CALOOCAN CITY6012METROMANILA610414006211PhilippinesF901NJ6191000000000028668380J701YM711EntryLevel",
		OrigCity:                "",
		OrigState:               "",
		MTCN:                    "3010897610",
		NewMTCN:                 "1809483010897610",
		ExchangeRate:            "1.0000000",
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
		Compliance: ComplianceDetails{
			IDDetails: IDDetails{
				IDType:    "A",
				IDCountry: "Philippines",
				IDNumber:  "987654321",
			},
			IDIssued:          "12012012",
			IDExpiry:          "12012019",
			Birthdate:         "12011996",
			BirthCountry:      "",
			Nationality:       "",
			SourceFunds:       "",
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
			GalacticID:              "",
			EmploymentPositionLevel: "Entry Level",
		},
		BankName:         "",
		BankLocation:     "",
		AccountNumber:    "",
		BankCode:         "",
		AccountSuffix:    "",
		MyWUNumber:       "185015541",
		MyWuPoints:       "",
		MyWUEnrollTag:    "mywu",
		TerminalID:       "PH259ART001A",
		OperatorID:       "5",
		RemoteTerminalID: "",
		RemoteOperatorID: "",
	}
)

func TestQuickPayStore(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          QPSRequest
		expectedReq testQPSRequest
		want        *ConfirmedDetails
		wantErr     bool
	}{
		{
			name: "Success",
			in:   QPSReq,
			expectedReq: testQPSRequest{
				WU: testQPSRequestWU{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: userCode,
						UserCode:     json.Number(userCode),
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testQPSRequestBody{
						Module:  "wuqp",
						Request: "wuqp-store",
						Param:   QPSReq,
					},
					Signature: "placeholderSignature",
				},
			},
			want: &ConfirmedDetails{
				AdvisoryText: "04 04 18 23:59",
				MTCN:         "6711768425",
				NewMTCN:      "1809486711768425",
				FilingDate:   "04-04",
				FilingTime:   "1159P EDT",
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
			ts := httptest.NewServer(successQPSHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(ts.Client(), "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.WU.Header.Token

			got, err := s.QuickPayStore(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatal(err)
				t.Errorf("QuickPayStore() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if !cmp.Equal(test.want, got) {
				t.Fatal(cmp.Diff(test.want, got))
			}
		})
	}
}

func successQPSHandler(t *testing.T, expectedReq testQPSRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/1.1/wuqp-store" {
			t.Errorf("expected request to '/1.1/wuqp-store', got '%s'", req.URL.EscapedPath())
		}

		var newReq testQPSRequest
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
						"advisory_text":"04 04 18 23:59",
						"mtcn":"6711768425",
						"new_mtcn":"1809486711768425",
						"filing_date":"04-04",
						"filing_time":"1159P EDT",
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
						"other_message_1": [
							"1",
							"2",
							"3"
						],
						"other_message_2": [
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

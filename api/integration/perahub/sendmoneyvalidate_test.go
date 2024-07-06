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

type testSMVBody struct {
	Module  string     `json:"module"`
	Request string     `json:"request"`
	Param   SMVRequest `json:"param"`
}

type testSMVWU struct {
	Header    RequestHeader `json:"header"`
	Body      testSMVBody   `json:"body"`
	Signature string        `json:"signature"`
}

type testSMVRequest struct {
	WU testSMVWU `json:"uspwuapi"`
}

var smvReq = SMVRequest{
	FrgnRefNo:                 "cdffb28cea04d8ee9cb1",
	SenderNameType:            "D",
	SenderFirstName:           "WINNIE",
	SenderMiddleName:          "",
	SenderLastName:            "CONSTANTINO",
	SenderAddrLine1:           "1953 PH 3B BLOCK 6 LOT 9 CAMARIN",
	SenderAddrLine2:           "",
	SenderCity:                "CALOOCAN CITY",
	SenderState:               "METRO MANILA",
	SenderPostalCode:          "1400",
	SenderAddrCountryCode:     "PH",
	SenderAddrCurrencyCode:    "PHP",
	SenderContactPhone:        "9089488895",
	SenderMobileCountryCode:   "63",
	SenderMobileNo:            "9089488895",
	SenderAddrCountryName:     "Philippines",
	MyWUNumber:                "185015541",
	IDType:                    "A",
	IDCountry:                 "Philippines",
	IDNumber:                  "987654321",
	IDIssued:                  "12012012",
	IDExpiry:                  "12012019",
	DateOfBirth:               "12011996",
	Occupation:                "Airline Maritime Employee",
	CountryOfBirth:            "Philippines",
	Nationality:               "Philippines",
	Gender:                    "",
	SourceOfFunds:             "Salary Income",
	ReceiverRelationship:      "",
	AckFlag:                   "X",
	SenderReason:              "",
	MyWUEnrollTag:             "none",
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
	DestCountry:               "PH",
	DestCurrency:              "PHP",
	DestStateCode:             "",
	DestCity:                  "",
	TransactionType:           "SO",
	PrincipalAmount:           "600000",
	FixedAmountFlag:           "N",
	PromoCode:                 "",
	Message:                   []string{},
	BankName:                  "",
	AccountNumber:             "",
	BankCode:                  "",
	BankLocation:              "",
	AccountSuffix:             "",
	TerminalID:                "PH259ART001A",
	OperatorID:                "5",
	RemoteTerminalID:          "",
	RemoteOperatorID:          "",
	EmploymentPositionLevel:   "Entry Level",
	TransactionPurpose:        "",
	IsOnBehalf:                "",
	GalacticID:                StrPtr("1000000000028668380"),
}

func TestSendMoneyValidate(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          SMVRequest
		expectedReq testSMVRequest
		want        *SMVResponseBody
		wantErr     bool
	}{
		{
			name: "Success",
			in:   smvReq,
			expectedReq: testSMVRequest{
				WU: testSMVWU{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testSMVBody{
						Module:  "wuso",
						Request: "SMvalidate",
						Param:   smvReq,
					},
					Signature: "placeholderSignature",
				},
			},
			want: &SMVResponseBody{
				ServiceCode: ServiceCode{
					AddlSvcChg: "1103FEE1205115001305115000103MSG0201003010210300022010230109001096210104500002021003035009703NNN9806PILPIL",
				},
				Compliance: Compliance{
					ComplianceBuf: "0108UNI_01_S0201A03099876543210611Philippines0708120119960825 AirlineMaritime Employee0908120120121008120120199901X3311Philippines3411Philippines4413Salary Income6724VOX DEI PROTOCOL SYSTEMS7401YB401157321953 PH 3B BLOCK 6 LOT 9 CAMARIN5913CALOOCAN CITY6012METRO MANILA610414006211PhilippinesF901NJ6191000000000028668380J701YM711Entry Level",
				},
				PaymentDetails: PaymentDetails{
					OrigCity:  "METRO MANILA",
					OrigState: "MA",
				},
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
				Promotions: Promotions{
					PromoDesc:       "",
					PromoMessage:    "",
					SenderPromoCode: "",
				},
				NewDetails: NewDetails{
					MTCN:       "3010897610",
					NewMTCN:    "1809483010897610",
					FilingDate: "04-04",
					FilingTime: "1056P EDT",
				},
				PreferredCustomer: PreferredCustomer{
					MyWUNumber: "185015541",
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(successSMVHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.WU.Header.Token

			got, err := s.SendMoneyValidate(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("SendMoneyValidate() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if !cmp.Equal(test.want, got) {
				t.Fatal(cmp.Diff(test.want, got))
			}
		})
	}
}

func successSMVHandler(t *testing.T, expectedReq testSMVRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/1.1/wuso-validate" {
			t.Errorf("expected request to '/1.1/wuso-validate', got '%s'", req.URL.EscapedPath())
		}

		var newReq testSMVRequest
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
					"service_code":{
						"addl_service_charges":"1103FEE1205115001305115000103MSG0201003010210300022010230109001096210104500002021003035009703NNN9806PILPIL"
					},
					"compliance":{
						"compliance_data_buffer":"0108UNI_01_S0201A03099876543210611Philippines0708120119960825 AirlineMaritime Employee0908120120121008120120199901X3311Philippines3411Philippines4413Salary Income6724VOX DEI PROTOCOL SYSTEMS7401YB401157321953 PH 3B BLOCK 6 LOT 9 CAMARIN5913CALOOCAN CITY6012METRO MANILA610414006211PhilippinesF901NJ6191000000000028668380J701YM711Entry Level"
					},
					"payment_details":{
						"originating_city":"METRO MANILA",
						"originating_state":"MA"
					},
					"financials":{
						"taxes":{
							"municipal_tax":0,
							"state_tax":0,
							"county_tax":0
						},
						"originators_principal_amount":600000,
						"destination_principal_amount":600000,
						"gross_total_amount":611500,
						"plus_charges_amount":0,
						"charges":11500,
						"message_charge":0
					},
					"promotions":{
						"promo_code_description":"",
						"promo_message":"",
						"sender_promo_code":""
					},
					"new_details":{
						"mtcn":"3010897610",
						"new_mtcn":"1809483010897610",
						"filing_date":"04-04",
						"filing_time":"1056P EDT"
					},
					"preferred_customer":{
						"mywu_number":"185015541"
					}
				}
			}
		}`)
	}
}

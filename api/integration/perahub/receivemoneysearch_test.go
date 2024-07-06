package perahub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type testRMSBody struct {
	Module  string          `json:"module"`
	Request string          `json:"request"`
	Param   RMSearchRequest `json:"param"`
}

type testRMSData struct {
	Header    RequestHeader `json:"header"`
	Body      testRMSBody   `json:"body"`
	Signature string        `json:"signature"`
}

type testRMSRequest struct {
	Data testRMSData `json:"uspwuapi"`
}

var rmsReq = RMSearchRequest{
	MTCN:         "9370896950",
	DestCurrency: "PHP",
	FrgnRefNo:    "13b6945f956cc3922ad8",
	TerminalID:   "PH259ART001A",
	OperatorID:   "5",
}

var rmsRes = RMSearchResponseBody{
	Txn: RMSPaymentTransaction{
		Sender: Contact{
			Name: Name{
				NameType:  "D",
				FirstName: "MICHAELLA CHI",
				LastName:  "DELA",
			},
			Address: RMSAddress{
				City:  "QUEZON CITY",
				State: "METRO MANILA",
				CountryCode: RMSCountryCode{
					IsoCode: RMSIsoCode{
						Country:  "PH",
						Currency: "PHP",
					},
				},
				PostalCode: "1400",
				Street:     "1900 MAGINHAWA STREET",
			},
			Phone: "90934140007",
			Mobile: RMSMobilePhone{
				RawPhone: json.RawMessage(`{"phone_number":{"country_code":"PHP","national_number":"1324423"}}`),
				Phone: RMSPhoneNumber{
					CountryCode: "63",
					Number:      "90934140007",
				},
			},
			MobileDetails: RMSMobileDetails{
				CountryCode: "63",
				Number:      "90934140007",
			},
		},
		Receiver: Contact{
			Name: Name{
				NameType:  "D",
				FirstName: "WINNIE",
				LastName:  "CONSTANTINO",
			},
			Address: RMSAddress{
				CountryCode: RMSCountryCode{
					IsoCode: RMSIsoCode{
						Country:  "PH",
						Currency: "PHP",
					},
				},
			},
			Mobile: RMSMobilePhone{
				RawPhone: json.RawMessage(`{"phone_number":{"country_code":"PHP","national_number":"3242432"}}`),
				Phone: RMSPhoneNumber{
					CountryCode: "63",
					Number:      "90934140007",
				},
			},
			MobileDetails: RMSMobileDetails{},
		},
		Financials: RMSFinancials{
			Taxes: RMSTaxes{
				TaxWorksheet: "",
			},
			GrossTotal: "304500",
			PayAmount:  "300000",
			Principal:  "300000",
			Charges:    "4500",
			Tolls:      "0",
		},
		Payment: RMSPaymentDetails{
			ExpectedPayoutLocation: RMSExpectedPayoutLocation{
				State: "",
				City:  "",
			},
			DestCountry: RMSCountryCurrency{
				IsoCode: RMSIsoCode{
					Country:  "PH",
					Currency: "PHP",
				},
			},
			OrigCountry: RMSCountryCurrency{
				IsoCode: RMSIsoCode{
					Country:  "PH",
					Currency: "PHP",
				},
			},
			OriginatingCity: "METRO MANILAPI1",
			TransactionType: "WMF",
			ExchangeRate:    "1",
			SenderDestCountry: RMSCountryCurrency{
				IsoCode: RMSIsoCode{
					Country:  "PH",
					Currency: "PHP",
				},
			},
		},
		FilingDate:       "04-04-18 ",
		FilingTime:       "0330A EDT",
		MoneyTransferKey: "3624030255",
		PayStatus:        "WC",
		Mtcn:             "9370896950",
		NewMtcn:          "1809489370896950",
		Fusion: RMSFusion{
			FusionStatus:  "WC",
			AccountNumber: "",
		},
		WuNetworkAgentIndicator: "",
	},
	GrossPayout: "3000",
	DST:         "0",
	NetPayout:   "3000",
}

func TestRMSword(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          RMSearchRequest
		expectedReq testRMSRequest
		want        *RMSearchResponseBody
		wantErr     bool
	}{
		{
			name: "Success",
			in:   rmsReq,
			expectedReq: testRMSRequest{
				Data: testRMSData{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testRMSBody{
						Module:  "wupo",
						Request: "search",
						Param:   rmsReq,
					},
					Signature: "placeholderSignature",
				},
			},
			want: &rmsRes,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(RMSSuccessHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.Data.Header.Token

			got, err := s.RMSearch(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("RMSword() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			o := cmp.Options{
				cmpopts.IgnoreFields(
					RMSMobilePhone{}, "RawPhone",
				),
				cmpopts.IgnoreFields(
					Contact{}, "RawMobileDetails",
				),
			}

			if !cmp.Equal(test.want, got, o) {
				t.Fatal(cmp.Diff(test.want, got, o))
			}
		})
	}
}

func RMSSuccessHandler(t *testing.T, expectedReq testRMSRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/1.1/wupo-search" {
			t.Errorf("expected request to '/1.1/wupo-search', got '%s'", req.URL.EscapedPath())
		}

		var newReq testRMSRequest
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
					"payment_transaction":{
					   "sender":{
						  "name":{
							 "name_type":"D",
							 "first_name":"MICHAELLA CHI",
							 "last_name":"DELA"
						  },
						  "address":{
							 "city":"QUEZON CITY",
							 "state":"METRO MANILA",
							 "country_code":{
								"iso_code":{
								   "country_code":"PH",
								   "currency_code":"PHP"
								}
							 },
							 "state_zip":"1400",
							 "street":"1900 MAGINHAWA STREET"
						  },
						  "contact_phone":"90934140007",
						  "mobile_phone":{
							 "phone_number":{
								"country_code":"63",
								"national_number":"90934140007"
							 }
						  },
						  "mobile_details":{
							 "city_code":"63",
							 "number":"90934140007"
						  }
					   },
					   "receiver":{
						  "name":{
							 "name_type":"D",
							 "first_name":"WINNIE",
							 "last_name":"CONSTANTINO"
						  },
						  "address":{
							 "country_code":{
								"iso_code":{
								   "country_code":"PH",
								   "currency_code":"PHP"
								}
							 }
						  },
						  "preferred_customer":{
							 "mywu_number":""
						  },
						  "mobile_phone":{
							 "phone_number":{
								"country_code":"63",
								"national_number":"90934140007"
							 }
						  },
						  "mobile_details":{}
					   },
					   "financials":{
						  "taxes":{
							 "tax_worksheet":""
						  },
						  "gross_total_amount":304500,
						  "pay_amount":300000,
						  "principal_amount":300000,
						  "charges":4500,
						  "tolls":0
					   },
					   "payment_details":{
						  "expected_payout_location":{
							 "state_code":"",
							 "city":""
						  },
						  "destination_country_currency":{
							 "iso_code":{
								"country_code":"PH",
								"currency_code":"PHP"
							 }
						  },
						  "originating_country_currency":{
							 "iso_code":{
								"country_code":"PH",
								"currency_code":"PHP"
							 }
						  },
						  "originating_city":"METRO MANILAPI1",
						  "transaction_type":"WMF",
						  "exchange_rate":1,
						  "original_destination_country_currency":{
							 "iso_code":{
								"country_code":"PH",
								"currency_code":"PHP"
							 }
						  }
					   },
					   "filing_date":"04-04-18 ",
					   "filing_time":"0330A EDT",
					   "money_transfer_key":"3624030255",
					   "pay_status_description":"WC",
					   "mtcn":"9370896950",
					   "new_mtcn":"1809489370896950",
					   "fusion":{
						  "fusion_status":"WC",
						  "account_number":""
					   },
					   "wu_network_agent_indicator":""
					},
					"gross_payout":3000,
					"dst":0,
					"net_payout":3000
				 }
			}
		}`)
	}
}

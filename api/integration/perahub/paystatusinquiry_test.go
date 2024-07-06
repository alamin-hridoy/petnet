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

type testPaySIRequestBody struct {
	Module  string       `json:"module"`
	Request string       `json:"request"`
	Param   PaySIRequest `json:"param"`
}

type testPaySIRequestWU struct {
	Header    RequestHeader        `json:"header"`
	Body      testPaySIRequestBody `json:"body"`
	Signature string               `json:"signature"`
}

type testPaySIRequest struct {
	WU testPaySIRequestWU `json:"uspwuapi"`
}

func TestPayStatusInquiry(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          PaySIRequest
		expectedReq testPaySIRequest
		want        *PaySIResponse
		wantErr     bool
	}{
		{
			name: "Success",
			in: PaySIRequest{
				MTCN:                    "9370896950",
				OriginatingCurrencyCode: "PHP",
				OriginatingCountryCode:  "PH",
				ForeignReferenceNo:      "13b6945f956cc3922ad8",
				TerminalID:              "PH259ART001A",
				OperatorID:              "5",
			},
			expectedReq: testPaySIRequest{
				WU: testPaySIRequestWU{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testPaySIRequestBody{
						Module:  "wupo",
						Request: "checkstat",
						Param: PaySIRequest{
							MTCN:                    "9370896950",
							OriginatingCurrencyCode: "PHP",
							OriginatingCountryCode:  "PH",
							ForeignReferenceNo:      "13b6945f956cc3922ad8",
							TerminalID:              "PH259ART001A",
							OperatorID:              "5",
						},
					},
					Signature: "placeholderSignature",
				},
			},
			want: &PaySIResponse{
				WU: PaySIResponseWU{
					Header: ResponseHeader{
						ErrorCode: "1",
						Message:   "good",
					},
					Body: PaySIResponseBody{
						PaymentTransactions: PaymentTransactions{
							PaymentTransaction: PaymentTransaction{
								Sender: Sender{
									Name: Name{
										NameType:  "D",
										FirstName: "JUAN",
										LastName:  "CRUZ",
									},
								},
								Receiver: Receiver{
									Name: Name{
										NameType:  "D",
										FirstName: "EDGAR",
										LastName:  "CRUZ",
									},
								},
								Financials: Financials{
									OrigPnplAmt: "550000",
								},
								PaymentDetails: InquiryPaymentDetails{
									OriginatingCountryCurrency: OriginatingCountryCurrency{
										IsoCode: IsoCode{
											CountryCode:  "JP",
											CurrencyCode: "JPY",
										},
									},
								},
								FilingDate:           "08 26 19",
								FilingTime:           "08:27:56",
								MoneyTransferKey:     "3626805515",
								PayStatusDescription: "W C",
							},
						},
						ForeignRemoteSystem: ForeignRemoteSystem{
							Identifier:  "WGPSPH2590T",
							ReferenceNo: "3894f2ccceb305460072",
							CounterID:   "PH259ART001A",
						},
						NumberMatches:     1,
						CurrentPageNumber: 1,
						LastPageNumber:    1,
					},
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(successPaySIHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.WU.Header.Token

			got, err := s.PayStatusInquiry(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatal(err)
				t.Errorf("PayStatusInquiry() error = %v, wantErr %v", err, test.wantErr)
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

func successPaySIHandler(t *testing.T, expectedReq testPaySIRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/1.1/wupo-checkstat" {
			t.Errorf("expected request to '/1.1/wupo-checkstat', got '%s'", req.URL.EscapedPath())
		}

		var newReq testPaySIRequest
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
					"payment_transactions":{
						"payment_transaction":{
							"sender":{
								"name":{
									"name_type":"D",
									"first_name":"JUAN",
									"last_name":"CRUZ"
								}
							},
							"receiver":{
								"name":{
									"name_type":"D",
									"first_name":"EDGAR",
									"last_name":"CRUZ"
								}
							},
							"financials":{
								"originators_principal_amount":550000
							},
							"payment_details":{
								"originating_country_currency":{
									"iso_code":{
										"country_code":"JP",
										"currency_code":"JPY"
									}
								}
							},
							"filing_date":"08 26 19",
							"filing_time":"08:27:56",
							"money_transfer_key":"3626805515",
							"pay_status_description":"W C"
						}
					},
					"foreign_remote_system":{
						"identifier":"WGPSPH2590T",
						"reference_no":"3894f2ccceb305460072",
						"counter_id":"PH259ART001A"
					},
					"number_matches":1,
					"current_page_number":1,
					"last_page_number":1
				}
			}
		}`)
	}
}

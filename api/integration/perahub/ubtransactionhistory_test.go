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

type testUBTHBody struct {
	Module  string      `json:"module"`
	Request string      `json:"request"`
	Param   UBTHRequest `json:"param"`
}

type testUBTH struct {
	Header    RequestHeader `json:"header"`
	Body      testUBTHBody  `json:"body"`
	Signature string        `json:"signature"`
}

type testUBTHRequest struct {
	WU testUBTH `json:"uspwuapi"`
}

var ubthReq = UBTHRequest{
	UserCode:     "MC0216929DRT",
	DateStart:    "2017-06-15 00:00:00",
	DateEnd:      "2017-08-14 23:59:59",
	OperatorCode: "",
	LocationCode: "",
	SessionID:    "",
}

func TestUBTransactionHistory(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          UBTHRequest
		expectedReq testUBTHRequest
		want        *UBTHResponse
		wantErr     bool
	}{
		{
			name: "Success",
			in:   ubthReq,
			expectedReq: testUBTHRequest{
				WU: testUBTH{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testUBTHBody{
						Module:  "Transaction",
						Request: "get_transaction",
						Param:   ubthReq,
					},
					Signature: "placeholderSignature",
				},
			},
			want: &UBTHResponse{
				WU: UBTHResponseData{
					Header: ResponseHeader{
						ErrorCode: "1",
						Message:   "Good",
					},
					Body: UBTHResponseBody{
						TranList: TranListDetails{
							Transaction: []Transaction{
								{
									ID:           "98328",
									TranDate:     "2017-08-14 12:28:07",
									Amount:       "1.00",
									Branch:       "NA",
									TranDesc:     "Fund Transfer To Third Party - Reversal",
									Balance:      "16149.80",
									TrxnCurrency: "608",
									PplCode:      "",
								},
								{
									ID:           "97974",
									TranDate:     "2017-08-14 12:28:07",
									Amount:       "-1.00",
									Branch:       "NA",
									TranDesc:     "Fund Transfer To Third Party",
									Balance:      "16147.80",
									TrxnCurrency: "608",
									PplCode:      "",
								},
								{
									ID:           "98156",
									TranDate:     "2017-08-14 12:28:06",
									Amount:       "1.00",
									Branch:       "NA",
									TranDesc:     "Fund Transfer To Third Party - Reversal",
									Balance:      "16148.80",
									TrxnCurrency: "608",
									PplCode:      "",
								},
							},
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
			ts := httptest.NewServer(successUBTHHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.WU.Header.Token

			got, err := s.UBTransactionHistory(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatal(err)
				t.Errorf("UBTransactionHistory() error = %v, wantErr %v", err, test.wantErr)
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

func successUBTHHandler(t *testing.T, expectedReq testUBTHRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/Transaction" {
			t.Errorf("expected request to '/Transaction', got '%s'", req.URL.EscapedPath())
		}

		var newReq testUBTHRequest
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
					"message":"Good"
				},
				"body":{
					"trxnList": {
						"trxn": [{
							"id": "98328",
							"tranDate": "2017-08-14 12:28:07",
							"amount": "1.00",
							"branch": "NA",
							"tranDesc": "Fund Transfer To Third Party - Reversal",
							"balance": "16149.80",
							"trxnCurrency": "608",
							"pplcode": ""
						},
						{
							"id": "97974",
							"tranDate": "2017-08-14 12:28:07",
							"amount": "-1.00",
							"branch": "NA",
							"tranDesc": "Fund Transfer To Third Party",
							"balance": "16147.80",
							"trxnCurrency": "608",
							"pplcode": ""
						},
						{
							"id": "98156",
							"tranDate": "2017-08-14 12:28:06",
							"amount": "1.00",
							"branch": "NA",
							"tranDesc": "Fund Transfer To Third Party - Reversal",
							"balance": "16148.80",
							"trxnCurrency": "608",
							"pplcode": ""
						}]
					}
				}
			}
		}`)
	}
}

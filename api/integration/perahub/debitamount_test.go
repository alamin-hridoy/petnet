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

type testDMCBody struct {
	Module  string             `json:"module"`
	Request string             `json:"request"`
	Param   DebitAmountRequest `json:"param"`
}

type testDMCWU struct {
	Header    RequestHeader `json:"header"`
	Body      testDMCBody   `json:"body"`
	Signature string        `json:"signature"`
}

type testDMCRequest struct {
	WU testDMCWU `json:"uspwuapi"`
}

func TestDebitAmount(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          DebitAmountRequest
		expectedReq testDMCRequest
		want        *DebitAmountResponse
		wantErr     bool
	}{
		{
			name: "Success",
			in: DebitAmountRequest{
				Username:      "RV101577JSCT",
				Amount:        "10068",
				Password:      "3627909A29C31381A071EC27F7C9CA97726182AED29A7DDD2E54353322CFB30ABB9E3A6DF2AC2C20FE23436311D678564D0C8D305930575F60E2D3D048184D79",
				OperationCode: "2221111",
				SessionID:     "111",
			},
			expectedReq: testDMCRequest{
				WU: testDMCWU{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testDMCBody{
						Module:  "DebitAmount",
						Request: "debit_amount",
						Param: DebitAmountRequest{
							Username:      "RV101577JSCT",
							Amount:        "10068",
							Password:      "3627909A29C31381A071EC27F7C9CA97726182AED29A7DDD2E54353322CFB30ABB9E3A6DF2AC2C20FE23436311D678564D0C8D305930575F60E2D3D048184D79",
							OperationCode: "2221111",
							SessionID:     "111",
						},
					},
					Signature: "placeholderSignature",
				},
			},
			want: &DebitAmountResponse{
				WU: DebitAmountResponseWU{
					Header: ResponseHeader{
						ErrorCode: "1",
						Message:   "Good",
					},

					Body: ResponseBody{
						Code:           "1",
						AmountDeducted: 10068,
					},
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(debitAmountHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.WU.Header.Token

			got, err := s.DebitAmount(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("DebitAmount() error = %v, wantErr %v", err, test.wantErr)
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

func debitAmountHandler(t *testing.T, expectedReq testDMCRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/DebitAmount" {
			t.Errorf("expected request to '/DebitAmount', got '%s'", req.URL.EscapedPath())
		}

		var newReq testDMCRequest
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
					"code": "1",
					"amount_deducted": 10068
				}
			}
		}`)
	}
}

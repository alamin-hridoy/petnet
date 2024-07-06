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

type testSendOTPBody struct {
	Module  string         `json:"module"`
	Request string         `json:"request"`
	Param   SendOTPRequest `json:"param"`
}

type testSendOTPData struct {
	Header    RequestHeader   `json:"header"`
	Body      testSendOTPBody `json:"body"`
	Signature string          `json:"signature"`
}

type testSendOTPRequest struct {
	Data testSendOTPData `json:"uspwuapi"`
}

func TestTransactionOTP(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          SendOTPRequest
		expectedReq testSendOTPRequest
		wantErr     bool
	}{
		{
			name: "Success",
			in: SendOTPRequest{
				UserCode:  "RV101577JSCT",
				SessionID: "213123",
			},
			expectedReq: testSendOTPRequest{
				Data: testSendOTPData{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testSendOTPBody{
						Module:  "SendSMSUser",
						Request: "send_sms_user",
						Param: SendOTPRequest{
							UserCode:  "RV101577JSCT",
							SessionID: "213123",
						},
					},
					Signature: "placeholderSignature",
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(transactionOTPHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.Data.Header.Token

			got, err := s.TransactionOTP(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if got.Data.Header.ErrorCode != "1" {
				t.Errorf("want error code 1 but got %q", got.Data.Header.ErrorCode)
			}
		})
	}
}

func transactionOTPHandler(t *testing.T, expectedReq testSendOTPRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/SendSMSUser" {
			t.Errorf("expected request to '/SendSMSUser', got '%s'", req.URL.EscapedPath())
		}

		var newReq testSendOTPRequest
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
					"message": "good"
				},
				"body": {}
			}
		}`)
	}
}

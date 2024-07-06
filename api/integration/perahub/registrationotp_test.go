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

type testOTPBody struct {
	Module  string     `json:"module"`
	Request string     `json:"request"`
	Param   OTPRequest `json:"param"`
}

type testOTPData struct {
	Header    RequestHeader `json:"header"`
	Body      testOTPBody   `json:"body"`
	Signature string        `json:"signature"`
}

type testOTPRequest struct {
	Data testOTPData `json:"uspwuapi"`
}

func TestRegOTP(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          OTPRequest
		expectedReq testOTPRequest
		wantErr     bool
	}{
		{
			name: "Success",
			in: OTPRequest{
				MobileNo: "09999378034",
			},
			expectedReq: testOTPRequest{
				Data: testOTPData{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testOTPBody{
						Module:  "SendSmsNewUser",
						Request: "send_sms_new_user",
						Param: OTPRequest{
							MobileNo: "09999378034",
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
			ts := httptest.NewServer(regOTPHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.Data.Header.Token

			got, err := s.RegistrationOTP(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("RegistrationOTP() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if got.Data.Header.ErrorCode != "1" {
				t.Errorf("want error code 1 but got %q", got.Data.Header.ErrorCode)
			}
		})
	}
}

func regOTPHandler(t *testing.T, expectedReq testOTPRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/SendSmsNewUser" {
			t.Errorf("expected request to '/SendSmsNewUser', got '%s'", req.URL.EscapedPath())
		}

		var newReq testOTPRequest
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
					"mobile": "09999378034"
				}
			}
		}`)
	}
}

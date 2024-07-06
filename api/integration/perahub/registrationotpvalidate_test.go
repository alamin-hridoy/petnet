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

type testOTPValidateBody struct {
	Module  string             `json:"module"`
	Request string             `json:"request"`
	Param   OTPValidateRequest `json:"param"`
}

type testOTPValidateData struct {
	Header    RequestHeader       `json:"header"`
	Body      testOTPValidateBody `json:"body"`
	Signature string              `json:"signature"`
}

type testOTPValidateRequest struct {
	Data testOTPValidateData `json:"uspwuapi"`
}

func TestRegOTPValidate(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          OTPValidateRequest
		expectedReq testOTPValidateRequest
		want        *OTPValidateResponse
		wantErr     bool
	}{
		{
			name: "Success",
			in: OTPValidateRequest{
				MobileNo: "09055195950",
				OTPCode:  182918,
			},
			expectedReq: testOTPValidateRequest{
				Data: testOTPValidateData{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testOTPValidateBody{
						Module:  "ValidateSMSNewUser",
						Request: "validate_otp_new_user",
						Param: OTPValidateRequest{
							MobileNo: "09055195950",
							OTPCode:  182918,
						},
					},
					Signature: "placeholderSignature",
				},
			},
			want: &OTPValidateResponse{
				Data: OTPResponseData{
					Header: ResponseHeader{
						ErrorCode: "1",
						Message:   "Good",
					},
					Body: OTPResponseBody{
						MobileNo: "09055195950",
					},
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(regOTPValidateHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.Data.Header.Token

			got, err := s.RegistrationOTPValidate(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("RegistrationOTPValidate() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if got.Data.Header.ErrorCode != "1" {
				t.Errorf("want error code 1 but got %q", got.Data.Header.ErrorCode)
			}

			if !cmp.Equal(test.want, got) {
				t.Fatal(cmp.Diff(test.want, got))
			}
		})
	}
}

func regOTPValidateHandler(t *testing.T, expectedReq testOTPValidateRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/ValidateSMSNewUser" {
			t.Errorf("expected request to '/ValidateSMSNewUser', got '%s'", req.URL.EscapedPath())
		}

		var newReq testOTPValidateRequest
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
					"mobile": "09055195950"
				}
			}
		}`)
	}
}

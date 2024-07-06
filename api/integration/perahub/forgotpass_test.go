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

type testForgotPassBody struct {
	Module  string            `json:"module"`
	Request string            `json:"request"`
	Param   ForgotPassRequest `json:"param"`
}

type testForgotPassData struct {
	Header    RequestHeader      `json:"header"`
	Body      testForgotPassBody `json:"body"`
	Signature string             `json:"signature"`
}

type testForgotPassRequest struct {
	Data testForgotPassData `json:"uspwuapi"`
}

var reqData = ForgotPassRequest{
	UserName: "reymar0524",
}

func TestForgotPassword(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          ForgotPassRequest
		expectedReq testForgotPassRequest
		want        *ForgotPassResponse
		wantErr     bool
	}{
		{
			name: "Success",
			in:   reqData,
			expectedReq: testForgotPassRequest{
				Data: testForgotPassData{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testForgotPassBody{
						Module:  "forgot_pwd_init",
						Request: "forgot_pwd_init",
						Param:   reqData,
					},
					Signature: "placeholderSignature",
				},
			},
			want: &ForgotPassResponse{
				Data: ForgotPassResponseData{
					Header: ResponseHeader{
						ErrorCode: "1",
						Message:   "Good",
					},
					Body: ForgotPassResponseBody{
						MobileNo: "********900",
					},
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(ForgotPassHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.Data.Header.Token

			got, err := s.ForgotPassword(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("ForgotPassword() error = %v, wantErr %v", err, test.wantErr)
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

func ForgotPassHandler(t *testing.T, expectedReq testForgotPassRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/forgot_pwd_init" {
			t.Errorf("expected request to '/forgot_pwd_init', got '%s'", req.URL.EscapedPath())
		}

		var newReq testForgotPassRequest
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
					"mobile": "********900"
				}
			}
		}`)
	}
}

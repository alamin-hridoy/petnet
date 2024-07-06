package perahub

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type testResetPassBody struct {
	Module  string           `json:"module"`
	Request string           `json:"request"`
	Param   ResetPassRequest `json:"param"`
}

type testResetPassData struct {
	Header    RequestHeader     `json:"header"`
	Body      testResetPassBody `json:"body"`
	Signature string            `json:"signature"`
}

type testResetPassRequest struct {
	Data testResetPassData `json:"uspwuapi"`
}

var reqRPData = ResetPassRequest{
	UserName:       "REYNOVFAD05241980RNF",
	Password:       "3627909A29C31381A071EC27F7C9CA97726182AED29A7DDD2E54353322CFB30ABB9E3A6DF2AC2C20FE23436311D678564D0C8D305930575F60E2D3D048184D79",
	SecretQuestion: "number 1?",
	Answer:         "4DFF4EA340F0A823F15D3F4F01AB62EAE0E5DA579CCB851F8DB9DFE84C58B2B37B89903A740E1EE172DA793A6E79D560E5F7F9BD058A12A280433ED6FA46510A",
}

func TestResetPassword(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          ResetPassRequest
		expectedReq testResetPassRequest
		want        *ResetPassResponse
		wantErr     bool
	}{
		{
			name: "Success",
			in:   reqRPData,
			expectedReq: testResetPassRequest{
				Data: testResetPassData{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testResetPassBody{
						Module:  "ResetPassword",
						Request: "reset_password",
						Param:   reqRPData,
					},
					Signature: "placeholderSignature",
				},
			},
			want: &ResetPassResponse{
				Data: ResetPassResponseData{
					Header: ResponseHeader{
						ErrorCode: "1",
						Message:   "Good",
					},
					Body: ResetPassResponseBody{
						UserName: "REYNOVFAD05241980RNF",
					},
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(resetPassHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.Data.Header.Token

			got, err := s.ResetPassword(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, test.wantErr)
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

func resetPassHandler(t *testing.T, expectedReq testResetPassRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/ResetPassword" {
			t.Errorf("expected request to '/ResetPassword', got '%s'", req.URL.EscapedPath())
		}

		var newReq testResetPassRequest
		if err := json.NewDecoder(req.Body).Decode(&newReq); err != nil {
			t.Fatal(err)
		}

		hashPByte := sha512.New().Sum([]byte(expectedReq.Data.Body.Param.Password))

		reqHashPByte, err := hex.DecodeString(newReq.Data.Body.Param.Password)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(hashPByte, reqHashPByte) {
			t.Error(cmp.Diff(hashPByte, reqHashPByte))
		}

		hashAByte := sha512.New().Sum([]byte(expectedReq.Data.Body.Param.Answer))

		reqHashAByte, err := hex.DecodeString(newReq.Data.Body.Param.Answer)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(hashAByte, reqHashAByte) {
			t.Error(cmp.Diff(hashAByte, reqHashAByte))
		}

		opts := cmpopts.IgnoreFields(testResetPassRequest{}, "Data.Body.Param.Password", "Data.Body.Param.Answer")
		if !cmp.Equal(expectedReq, newReq, opts) {
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
					"username": "REYNOVFAD05241980RNF"
				}
			}
		}`)
	}
}

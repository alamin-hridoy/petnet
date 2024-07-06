package perahub

import (
	"context"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testBody struct {
	Module  string       `json:"module"`
	Request string       `json:"request"`
	Param   LoginRequest `json:"param"`
}

type testWU struct {
	Header    RequestHeader `json:"header"`
	Body      testBody      `json:"body"`
	Signature string        `json:"signature"`
}

type testRequest struct {
	WU testWU `json:"uspwuapi"`
}

func TestLogin(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"
	pw := "password"

	tests := []struct {
		name        string
		in          LoginRequest
		expectedReq testRequest
		wantErr     bool
	}{
		{
			name: "Success",
			in: LoginRequest{
				Username: "RN052480W",
				Password: pw,
			},
			expectedReq: testRequest{
				WU: testWU{
					Header: RequestHeader{
						Coy:          "yondu",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testBody{
						Module:  "SignOn",
						Request: "login",
						Param: LoginRequest{
							Username: "RN052480W",
							Password: pw,
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
			ts := httptest.NewServer(handler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.WU.Header.Token

			got, err := s.Login(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if got.FrgnRefNo == "" {
				t.Error("want foreign reference no but got empty")
			}
		})
	}
}

func handler(t *testing.T, expectedReq testRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/signin" {
			t.Errorf("expected request to '/signin', got '%s'", req.URL.EscapedPath())
		}

		var newReq testRequest
		if err := json.NewDecoder(req.Body).Decode(&newReq); err != nil {
			t.Fatal(err)
		}

		opt := cmp.FilterPath(func(p cmp.Path) bool { return p.Last().String() == ".Password" },
			cmp.Comparer(func(a, b string) bool {
				return a == fmt.Sprintf("%X", sha512.New().Sum([]byte(b))) ||
					b == fmt.Sprintf("%X", sha512.New().Sum([]byte(a)))
			}))
		if !cmp.Equal(expectedReq, newReq, opt) {
			t.Error(cmp.Diff(expectedReq, newReq, opt))
		}

		res.WriteHeader(200)
		res.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(res, `{
			"uspwuapi": {
				"header": {
					"errorcode": "1",
					"message": "good"
				},
				"body": {
					"foreign_reference_no": "d13400918208a58b4e42869316e55963994b6031",
					"customer": {
						"surname": "Doe",
						"givenname": "John",
						"Birthdate": "1987-02-14",
						"Nationality": "Vinland",
						"PresentAddress": "26 Butterfly Road",
						"Occupation": "Engineer",
						"NameOfEmployer": "N/A",
						"ValidIdentification": "N/A",
						"WuCardNo": "27947290572",
						"DebitCardNo": "08295823095",
						"LoyaltyCardNo": "29035732085"
					}
				}
			}
		}`)
	}
}

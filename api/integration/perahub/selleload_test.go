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

type testSELRequestBody struct {
	Module  string        `json:"module"`
	Request string        `json:"request"`
	Param   SELoadRequest `json:"param"`
}

type testSELRequestWU struct {
	Header    RequestHeader      `json:"header"`
	Body      testSELRequestBody `json:"body"`
	Signature string             `json:"signature"`
}

type testSELRequest struct {
	WU testSELRequestWU `json:"uspwuapi"`
}

var selReq = SELoadRequest{
	SessionID:      "db49a1d1044b440aa14a7facac2cb138",
	ProductCode:    "SMGS50",
	Amount:         "50",
	Provider:       "New denoms (as of 14 Nov 2016)",
	ProductType:    "Eload",
	TargetMobileNo: "09170000000",
	Password:       "hashed password value",
}

func TestSellLoad(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          SELoadRequest
		expectedReq testSELRequest
		want        *SELoadResponse
		wantErr     bool
	}{
		{
			name: "Success",
			in:   selReq,
			expectedReq: testSELRequest{
				WU: testSELRequestWU{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testSELRequestBody{
						Module:  "eload",
						Request: "sell",
						Param:   selReq,
					},
					Signature: "placeholderSignature",
				},
			},
			want: &SELoadResponse{
				WU: SELResponseWU{
					Header: ResponseHeader{
						ErrorCode: "1",
						Message:   "good",
					},
					Body: SELResponseBody{
						RRN:       "db49a1",
						TID:       "WB3072702822",
						EPIN:      "CARD 86138664 PIN 3201142812949614",
						Timestamp: "2018-11-07 14:27:18",
					},
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(successSELHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.WU.Header.Token

			got, err := s.SellLoad(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatal(err)
				t.Errorf("SellLoad() error = %v, wantErr %v", err, test.wantErr)
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

func successSELHandler(t *testing.T, expectedReq testSELRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/eload" {
			t.Errorf("expected request to '/eload', got '%s'", req.URL.EscapedPath())
		}

		var newReq testSELRequest
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
					"RRN": "db49a1",
					"TID": "WB3072702822",
					"EPIN": "CARD 86138664 PIN 3201142812949614",
					"timestamp": "2018-11-07 14:27:18"
				}
			}
		}`)
	}
}

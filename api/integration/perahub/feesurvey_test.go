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

type testFSRequestBody struct {
	Module  string    `json:"module"`
	Request string    `json:"request"`
	Param   FSRequest `json:"param"`
}

type testFSRequestWU struct {
	Header    RequestHeader     `json:"header"`
	Body      testFSRequestBody `json:"body"`
	Signature string            `json:"signature"`
}

type testFSRequest struct {
	WU testFSRequestWU `json:"uspwuapi"`
}

func TestFS(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          FSRequest
		expectedReq testFSRequest
		want        string
		wantErr     bool
	}{
		{
			name: "Success",
			in: FSRequest{
				RefNo:                   "123123sadasd",
				PrincipalAmount:         "100000",
				FixedAmountFlag:         "N",
				DestinationCountryCode:  "PH",
				DestinationCurrencyCode: "PHP",
				TransactionType:         "SO",
				PromoCode:               "",
				Message:                 []string{},
				MessageLineCount:        "0",
				TerminalID:              "WBPt",
				OperatorID:              "001",
			},
			expectedReq: testFSRequest{
				WU: testFSRequestWU{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testFSRequestBody{
						Module:  "prereq",
						Request: "feesurvey",
						Param: FSRequest{
							RefNo:                   "123123sadasd",
							PrincipalAmount:         "100000",
							FixedAmountFlag:         "N",
							DestinationCountryCode:  "PH",
							DestinationCurrencyCode: "PHP",
							TransactionType:         "SO",
							PromoCode:               "",
							Message:                 []string{},
							MessageLineCount:        "0",
							TerminalID:              "WBPt",
							OperatorID:              "001",
						},
					},
					Signature: "placeholderSignature",
				},
			},
			want: "0101C0202MB0303PIL",
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(successFSHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.WU.Header.Token

			got, err := s.FeeSurvey(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("FS() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

func successFSHandler(t *testing.T, expectedReq testFSRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/prereq" {
			t.Errorf("expected request to '/prereq', got '%s'", req.URL.EscapedPath())
		}

		var newReq testFSRequest
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
				"body": {
					"staging_buffer": "0101C0202MB0303PIL"
				}
			}
		}`)
	}
}

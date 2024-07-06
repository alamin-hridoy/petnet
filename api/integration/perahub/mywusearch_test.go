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

type testMWSBody struct {
	Module  string          `json:"module"`
	Request string          `json:"request"`
	Param   WUSearchRequest `json:"param"`
}

type testMWSWU struct {
	Header    RequestHeader `json:"header"`
	Body      testMWSBody   `json:"body"`
	Signature string        `json:"signature"`
}

type testMWSRequest struct {
	WU testMWSWU `json:"uspwuapi"`
}

var mwsReq = WUSearchRequest{
	ForRefNo:    "LBCTrackingNO01",
	SearchParam: "",
	SearchType:  "name",
	SenderFName: "MIGUEL",
	SenderLName: "FADRIQUELA",
	OperatorID:  "001",
	TerminalID:  "WBPt",
}

func TestMyWUSearch(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          WUSearchRequest
		expectedReq testMWSRequest
		want        *ResponseDetails
		wantErr     bool
	}{
		{
			name: "Success",
			in:   mwsReq,
			expectedReq: testMWSRequest{
				WU: testMWSWU{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testMWSBody{
						Module:  "prereq",
						Request: "mywusearch",
						Param:   mwsReq,
					},
					Signature: "placeholderSignature",
				},
			},
			want: &ResponseDetails{
				ForRefNo:            "LBCTrackingNO01",
				Surname:             "Doe",
				GivenName:           "John",
				Birthdate:           "1987-02-14",
				Nationality:         "Vinland",
				PresentAddress:      "26 Butterfly Road",
				Occupation:          "Engineer",
				NameOfEmployer:      "N/A",
				ValidIdentification: "N/A",
				WuCardNo:            "27947290572",
				DebitCardNo:         "08295823095",
				LoyaltyCardNo:       "29035732085",
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(successMWSHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.WU.Header.Token

			got, err := s.MyWUSearch(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatal(err)
				t.Errorf("MyWUSearch() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if !cmp.Equal(test.want, got) {
				t.Fatal(cmp.Diff(test.want, got))
			}
		})
	}
}

func successMWSHandler(t *testing.T, expectedReq testMWSRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/prereq" {
			t.Errorf("expected request to '/prereq', got '%s'", req.URL.EscapedPath())
		}

		var newReq testMWSRequest
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
					"foreign_reference_no": "LBCTrackingNO01",
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

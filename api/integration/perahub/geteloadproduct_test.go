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

type testELPRequestBody struct {
	Module  string     `json:"module"`
	Request string     `json:"request"`
	Param   ELPRequest `json:"param"`
}

type testELPRequestWU struct {
	Header    RequestHeader      `json:"header"`
	Body      testELPRequestBody `json:"body"`
	Signature string             `json:"signature"`
}

type testELPRequest struct {
	WU testELPRequestWU `json:"uspwuapi"`
}

func TestGetELoadProducts(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          ELPRequest
		expectedReq testELPRequest
		want        *ELPResponse
		wantErr     bool
	}{
		{
			name: "Success",
			in: ELPRequest{
				CurrentRow: 0,
				GroupCount: 100,
			},
			expectedReq: testELPRequest{
				WU: testELPRequestWU{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testELPRequestBody{
						Module:  "eload",
						Request: "product_list",
						Param: ELPRequest{
							CurrentRow: 0,
							GroupCount: 100,
						},
					},
					Signature: "placeholderSignature",
				},
			},
			want: &ELPResponse{
				WU: ELPResponseWU{
					Header: ResponseHeader{
						ErrorCode: "1",
						Message:   "good",
					},
					Body: ELPResponseBody{
						Items: []Item{
							{
								EloadID:     "101",
								ProductName: "Smart Surf Max Plus 995",
								ProductCode: "SMSURFMAXPLS995",
								Amount:      "995",
								Commission:  "4.42",
								Provider:    "New denoms (as of 1 Feb 2016)",
								ProductType: "Eload",
								Remarks:     "",
								Additional:  "2",
								UpdatedBy:   "6",
								DateUpdated: "2017-03-03 17:34:22",
							},
							{
								EloadID:     "102",
								ProductName: "Giga Surf 50",
								ProductCode: "SMGS50",
								Amount:      "50",
								Commission:  "3.94",
								Provider:    "New denoms (as of 14 Nov 2016)",
								ProductType: "Eload",
								Remarks:     "",
								Additional:  "2",
								UpdatedBy:   "6",
								DateUpdated: "2017-03-03 17:34:22",
							},
						},
						ItemCounts: []ItemCount{
							{
								TotalRows: "533",
							},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(getELPSuccessHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.WU.Header.Token

			got, err := s.GetELoadProducts(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("GetELoadProducts() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if got.WU.Header.ErrorCode != "1" {
				t.Errorf("want error code 1 but got %q", got.WU.Header.ErrorCode)
			}

			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

func getELPSuccessHandler(t *testing.T, expectedReq testELPRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/eload" {
			t.Errorf("expected request to '/eload', got '%s'", req.URL.EscapedPath())
		}

		var newReq testELPRequest
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
					"1":[
						{
							"eload_id":"101",
							"ProductName":"Smart Surf Max Plus 995",
							"ProductCode":"SMSURFMAXPLS995",
							"Amount":"995",
							"Commission":"4.42",
							"Provider":"New denoms (as of 1 Feb 2016)",
							"ProductType":"Eload",
							"Remarks":"",
							"Additional":"2",
							"updated_by":"6",
							"date_updated":"2017-03-03 17:34:22"
						},
						{
							"eload_id":"102",
							"ProductName":"Giga Surf 50",
							"ProductCode":"SMGS50",
							"Amount":"50",
							"Commission":"3.94",
							"Provider":"New denoms (as of 14 Nov 2016)",
							"ProductType":"Eload",
							"Remarks":"",
							"Additional":"2",
							"updated_by":"6",
							"date_updated":"2017-03-03 17:34:22"
						}
					],
					"2":[
						{
							"total_rows":"533"
						}
					]
				}
			}
		}`)
	}
}

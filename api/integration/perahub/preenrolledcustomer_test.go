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

type testPECustomerRequestBody struct {
	Module  string            `json:"module"`
	Request string            `json:"request"`
	Param   PECustomerRequest `json:"param"`
}

type testPECustomerRequestWU struct {
	Header    RequestHeader             `json:"header"`
	Body      testPECustomerRequestBody `json:"body"`
	Signature string                    `json:"signature"`
}

type testPECustomerRequest struct {
	WU testPECustomerRequestWU `json:"uspwuapi"`
}

func TestPreEnrolledCustomer(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          PECustomerRequest
		expectedReq testPECustomerRequest
		want        *PECustomer
		wantErr     bool
	}{
		{
			name: "Success",
			in: PECustomerRequest{
				CustomerCode: "LT1710898BCF",
			},
			expectedReq: testPECustomerRequest{
				WU: testPECustomerRequestWU{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testPECustomerRequestBody{
						Module:  "Registration",
						Request: "getinfo",
						Param: PECustomerRequest{
							CustomerCode: "LT1710898BCF",
						},
					},
					Signature: "placeholderSignature",
				},
			},
			want: &PECustomer{
				CustomerID:     3047954,
				CustomerNumber: "LT1710898BCF",
				LastName:       "Tranate",
				FirstName:      "Loralyn",
				MiddleName:     "Caragay",
				CurrentAddress: Address{
					AddrLine1:  "1953 PH 3B BLOCK 6 LOT 9 CAMARIN",
					AddrLine2:  "",
					City:       "CALOOCAN CITY",
					State:      "METRO MANILA",
					PostalCode: "1400",
					Country:    "Philippines",
				},
				TelNo:      "",
				EmailAdd:   "",
				MobileNo:   "09984427677",
				BirthDate:  "1989-10-17",
				Occupation: "None",
				PermaAddress: Address{
					AddrLine1:  "1953 PH 3B BLOCK 6 LOT 9 CAMARIN",
					AddrLine2:  "",
					City:       "CALOOCAN CITY",
					State:      "METRO MANILA",
					PostalCode: "1400",
					Country:    "Philippines",
				},
				Nationality:  "Philippines",
				Gender:       "female",
				CivilStatus:  "",
				TINID:        "",
				SSSID:        "",
				GSISID:       "",
				DriverLic:    "",
				SourceFund:   "Remittance",
				EmployerName: "None",
				NatureWork:   "",
				Employment:   "",
				CardNo:       "",
				Img: map[string]Image{
					"1": {
						Img:      "20170829393.png",
						IDType:   "Voter'sID",
						IDNumber: "08030150CJ1789LCT20000",
						Country:  "Philippines",
						Expiry:   "None",
						Issue:    "None",
					},
				},
				CardPoints:  0,
				CardPesoVal: 0,
				WBCardNo:    "",
				UBCardNo:    "0612345678",
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(successPECustomerHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.WU.Header.Token

			got, err := s.PreEnrolledCustomer(context.Background(), test.in.CustomerCode)
			if (err != nil) != test.wantErr {
				t.Fatal(err)
				t.Errorf("PreEnrolledCustomer() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if !cmp.Equal(test.want, got) {
				t.Fatal(cmp.Diff(test.want, got))
			}
		})
	}
}

func successPECustomerHandler(t *testing.T, expectedReq testPECustomerRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/Registration" {
			t.Errorf("expected request to '/Registration', got '%s'", req.URL.EscapedPath())
		}

		var newReq testPECustomerRequest
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
				"body": [
					{
						"customer_id":3047954,
						"customer_number":"LT1710898BCF",
						"last_name":"Tranate",
						"first_name":"Loralyn",
						"middle_name":"Caragay",
						"Current_address":{
							"addr_line1":"1953 PH 3B BLOCK 6 LOT 9 CAMARIN",
							"addr_line2":"",
							"city":"CALOOCAN CITY",
							"state_name":"METRO MANILA",
							"postal_code":"1400",
							"country":"Philippines"
						},
						"tel_no":"",
						"email_add":"",
						"mobile_no":"09984427677",
						"birth_date":"1989-10-17",
						"occupation":"None",
						"perma_address":{
							"addr_line1":"1953 PH 3B BLOCK 6 LOT 9 CAMARIN",
							"addr_line2":"",
							"city":"CALOOCAN CITY",
							"state_name":"METRO MANILA",
							"postal_code":"1400",
							"country":"Philippines"
						},
						"nationality":"Philippines",
						"gender":"female",
						"civil_status":"",
						"tin_id":null,
						"sss_id":null,
						"gsis_id":null,
						"driver_lic":null,
						"source_fund":"Remittance",
						"employer_name":"None",
						"nature_work":"",
						"employment":"",
						"card_no":null,
						"img":{
							"1":{
								"img":"20170829393.png",
								"id_type":"Voter'sID",
								"id_number":"08030150CJ1789LCT20000",
								"country":"Philippines",
								"expiry":"None",
								"issue":"None"
							}
						},
						"card_points":0,
						"card_peso_val":0,
						"wu_card_no":null,
						"ub_card_no":"0612345678"
					}
				]
			}
		}`)
	}
}

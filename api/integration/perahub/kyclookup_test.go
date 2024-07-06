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

type testKYCLookupBody struct {
	Module  string           `json:"module"`
	Request string           `json:"request"`
	Param   KYCLookupRequest `json:"param"`
}

type testKYCLookupWU struct {
	Header    RequestHeader     `json:"header"`
	Body      testKYCLookupBody `json:"body"`
	Signature string            `json:"signature"`
}

type testKYCLookupRequest struct {
	WU testKYCLookupWU `json:"uspwuapi"`
}

func TestKYCLookup(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          KYCLookupRequest
		expectedReq testKYCLookupRequest
		want        *KYCLookupBody
		wantErr     bool
	}{
		{
			name: "Success",
			in: KYCLookupRequest{
				RefNo:        "67de2e73d50dccd0d0f4",
				SearchType:   "by_id",
				TrxType:      "PAY",
				IDType:       "A",
				IDNumber:     "987654321",
				ContactPhone: "",
				MyWUNumber:   "181595968",
				FirstName:    "MIRASOL",
				LastName:     "SOLOMON",
				IsMulti:      "S",
				OperatorID:   "5",
				TerminalID:   "PH259ART001A",
			},
			expectedReq: testKYCLookupRequest{
				WU: testKYCLookupWU{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testKYCLookupBody{
						Module:  "prereq",
						Request: "kyclookup",
						Param: KYCLookupRequest{
							RefNo:        "67de2e73d50dccd0d0f4",
							SearchType:   "by_id",
							TrxType:      "PAY",
							IDType:       "A",
							IDNumber:     "987654321",
							ContactPhone: "",
							MyWUNumber:   "181595968",
							FirstName:    "MIRASOL",
							LastName:     "SOLOMON",
							IsMulti:      "S",
							OperatorID:   "5",
							TerminalID:   "PH259ART001A",
						},
					},
					Signature: "placeholderSignature",
				},
			},
			want: &KYCLookupBody{
				Customer: KYCCustomer{
					Name: Name{
						NameType:  "D",
						FirstName: "MIRASOL",
						LastName:  "SOLOMON",
					},
					Address: KYCCustomerAddress{
						AddrType:   "PRIMARY",
						AddrLine1:  "1953 PH 3B BLOCK 6 LOT 9 CAMARIN",
						AddrLine2:  "1953 PH 3B BLOCK 6 LOT 9 CAMARIN",
						City:       "CALOOCAN CITY",
						StateName:  "METRO MANILA",
						PostalCode: "1400",
						CountryDetails: KYCCountryDetails{
							CtryCode: "PH",
							CtryName: "PHILIPPINES",
						},
					},
					MyWUDetails: KYCMyWUDetails{
						MyWUNumber:       "181595968",
						IsConvenience:    "N",
						LevelCode:        "PHC",
						CurrentYrPts:     "135",
						EnrollmentSource: "P",
					},
					ComplianceDetails: KYCComplianceDetails{
						ComplianceFlagsBuffer: "J630__",
						ComplianceDataBuffer:  "J6191000000000043726053",
					},
					Email:        "JERICA.REAS@GMAIL.COM",
					ContactPhone: "639089488895",
					MobileNumber: MobileNumber{
						CountryCode:    "63",
						NationalNumber: "9089488895",
					},
					DateOfBirth:      "",
					SuppressFlag:     "N",
					MarketingDetails: []byte("[]"),
					KYCDetails: KYCDetails{
						IsKyced: "N",
					},
					Preferences: []byte("[]"),
				},
				Receiver: KYCReceiver{
					Receiver: []ReceiverDetails{
						{
							Name: Name{
								NameType:  "D",
								FirstName: "EMRE UMUT",
								LastName:  "KEBABCI",
							},
							ReceiverType: "L",
							Address: KYCAddress{
								Country: KYCCountryDetails{
									CtryCode: "PH",
									CtryName: "",
								},
							},
							ReceiverIndexNo: "1",
						},
						{
							Name: Name{
								NameType:  "D",
								FirstName: "CZARLYN GANUT",
								LastName:  "CAJAYON",
							},
							ReceiverType: "L",
							Address: KYCAddress{
								Country: KYCCountryDetails{
									CtryCode: "PH",
									CtryName: "",
								},
							},
							ReceiverIndexNo: "2",
						},
						{
							Name: Name{
								NameType:  "D",
								FirstName: "WINNIE",
								LastName:  "CONSTANTINO",
							},
							ReceiverType: "L",
							Address: KYCAddress{
								Country: KYCCountryDetails{
									CtryCode: "PH",
									CtryName: "",
								},
							},
							ReceiverIndexNo: "3",
						},
						{
							Name: Name{
								NameType:  "D",
								FirstName: "CZARLYN",
								LastName:  "CAJAYON",
							},
							ReceiverType: "L",
							Address: KYCAddress{
								Country: KYCCountryDetails{
									CtryCode: "PH",
									CtryName: "",
								},
							},
							ReceiverIndexNo: "4",
						},
						{
							Name: Name{
								NameType:  "D",
								FirstName: "WINNIE ANN",
								LastName:  "CONSTANTINO",
							},
							ReceiverType: "L",
							Address: KYCAddress{
								Country: KYCCountryDetails{
									CtryCode: "PH",
									CtryName: "",
								},
							},
							ReceiverIndexNo: "5",
						},
						{
							Name: Name{
								NameType:  "D",
								FirstName: "JERICA",
								LastName:  "NAPARATE",
							},
							ReceiverType: "L",
							Address: KYCAddress{
								Country: KYCCountryDetails{
									CtryCode: "PH",
									CtryName: "",
								},
							},
							ReceiverIndexNo: "6",
						},
						{
							Name: Name{
								NameType:  "D",
								FirstName: "MAXPEIN ZIN",
								LastName:  "DEL VALLE",
							},
							ReceiverType: "D",
							Address: KYCAddress{
								Country: KYCCountryDetails{
									CtryCode: "US",
									CtryName: "",
								},
							},
							ReceiverIndexNo: "8",
						},
						{
							Name: Name{
								NameType:  "D",
								FirstName: "MAXWELL LAURENT",
								LastName:  "DEL VALLE",
							},
							ReceiverType: "M",
							Address: KYCAddress{
								Country: KYCCountryDetails{
									CtryCode: "MX",
									CtryName: "",
								},
							},
							ReceiverIndexNo: "9",
						},
						{
							Name: Name{
								NameType:  "D",
								FirstName: "SANCHEZ",
								LastName:  "GONZALEZ",
							},
							ReceiverType: "I",
							Address: KYCAddress{
								Country: KYCCountryDetails{
									CtryCode: "CO",
									CtryName: "",
								},
							},
							ReceiverIndexNo: "10",
						},
					},
					NumberMatches: "9",
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(kycLookupSuccessHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.WU.Header.Token

			got, err := s.KYCLookup(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("KYCLookup() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if !cmp.Equal(test.want, got) {
				t.Fatal(cmp.Diff(test.want, got))
			}
		})
	}
}

func kycLookupSuccessHandler(t *testing.T, expectedReq testKYCLookupRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/prereq" {
			t.Errorf("expected request to '/prereq', got '%s'", req.URL.EscapedPath())
		}

		var newReq testKYCLookupRequest
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
				  "message":"Good"
			   },
			   "body":{
				  "customer":{
					 "name":{
						"name_type":"D",
						"first_name":"MIRASOL",
						"last_name":"SOLOMON"
					 },
					 "address":{
						"addr_type":"PRIMARY",
						"addr_line1":"1953 PH 3B BLOCK 6 LOT 9 CAMARIN",
						"addr_line2":"1953 PH 3B BLOCK 6 LOT 9 CAMARIN",
						"city":"CALOOCAN CITY",
						"state_name":"METRO MANILA",
						"postal_code":"1400",
						"country_details":{
						   "ctry_code":"PH",
						   "ctry_name":"PHILIPPINES"
						}
					 },
					 "mywu_details":{
						"mywu_number":"181595968",
						"is_convenience":"N",
						"level_code":"PHC",
						"current_yr_pts":135,
						"enrollment_source":"P"
					 },
					 "compliance_details":{
						"compliance_flags_buffer":"J630__",
						"compliance_data_buffer":"J6191000000000043726053"
					 },
					 "email":"JERICA.REAS@GMAIL.COM",
					 "contact_phone":"639089488895",
					 "mobile_number":{
						"ctry_code":"63",
						"National_number":"9089488895"
					 },
					 "date_of_birth":"",
					 "suppress_flag":"N",
					 "mktng_details":[],
					 "kyc_details":{
						"is_kyced":"N"
					 },
					 "preferences":[]
				  },
				  "receiver":{
					 "receiver":[
						{
						   "name":{
							  "name_type":"D",
							  "first_name":"EMRE UMUT",
							  "last_name":"KEBABCI"
						   },
						   "receiver_type":"L",
						   "address":{
							  "country_details":{
								 "ctry_code":"PH",
								 "ctry_name":""
							  }
						   },
						   "receiver_index_no":1
						},
						{
						   "name":{
							  "name_type":"D",
							  "first_name":"CZARLYN GANUT",
							  "last_name":"CAJAYON"
						   },
						   "receiver_type":"L",
						   "address":{
							  "country_details":{
								 "ctry_code":"PH",
								 "ctry_name":""
							  }
						   },
						   "receiver_index_no":2
						},
						{
						   "name":{
							  "name_type":"D",
							  "first_name":"WINNIE",
							  "last_name":"CONSTANTINO"
						   },
						   "receiver_type":"L",
						   "address":{
							  "country_details":{
								 "ctry_code":"PH",
								 "ctry_name":""
							  }
						   },
						   "receiver_index_no":3
						},
						{
						   "name":{
							  "name_type":"D",
							  "first_name":"CZARLYN",
							  "last_name":"CAJAYON"
						   },
						   "receiver_type":"L",
						   "address":{
							  "country_details":{
								 "ctry_code":"PH",
								 "ctry_name":""
							  }
						   },
						   "receiver_index_no":4
						},
						{
						   "name":{
							  "name_type":"D",
							  "first_name":"WINNIE ANN",
							  "last_name":"CONSTANTINO"
						   },
						   "receiver_type":"L",
						   "address":{
							  "country_details":{
								 "ctry_code":"PH",
								 "ctry_name":""
							  }
						   },
						   "receiver_index_no":5
						},
						{
						   "name":{
							  "name_type":"D",
							  "first_name":"JERICA",
							  "last_name":"NAPARATE"
						   },
						   "receiver_type":"L",
						   "address":{
							  "country_details":{
								 "ctry_code":"PH",
								 "ctry_name":""
							  }
						   },
						   "receiver_index_no":6
						},
						{
						   "name":{
							  "name_type":"D",
							  "first_name":"MAXPEIN ZIN",
							  "last_name":"DEL VALLE"
						   },
						   "receiver_type":"D",
						   "address":{
							  "country_details":{
								 "ctry_code":"US",
								 "ctry_name":""
							  }
						   },
						   "receiver_index_no":8
						},
						{
						   "name":{
							  "name_type":"D",
							  "first_name":"MAXWELL LAURENT",
							  "last_name":"DEL VALLE"
						   },
						   "receiver_type":"M",
						   "address":{
							  "country_details":{
								 "ctry_code":"MX",
								 "ctry_name":""
							  }
						   },
						   "receiver_index_no":9
						},
						{
						   "name":{
							  "name_type":"D",
							  "first_name":"SANCHEZ",
							  "last_name":"GONZALEZ"
						   },
						   "receiver_type":"I",
						   "address":{
							  "country_details":{
								 "ctry_code":"CO",
								 "ctry_name":""
							  }
						   },
						   "receiver_index_no":10
						}
					 ],
					 "number_matches":9
				  }
			   }
			}
		 }`)
	}
}

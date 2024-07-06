package perahub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"brank.as/petnet/api/core/static"
	"github.com/google/go-cmp/cmp"
)

type testFIBody struct {
	Module  string    `json:"module"`
	Request string    `json:"request"`
	Param   FIRequest `json:"param"`
}

type testFIWU struct {
	Header    RequestHeader `json:"header"`
	Body      testFIBody    `json:"body"`
	Signature string        `json:"signature"`
}

type testFIRequest struct {
	WU testFIWU `json:"uspwuapi"`
}

func TestFeeInquiry(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          FIRequest
		expectedReq testFIRequest
		want        *FIResponseBody
		wantErr     bool
	}{
		{
			name: "Success",
			in: FIRequest{
				FrgnRefNo:       "123123",
				PrincipalAmount: "100000",
				FixedAmountFlag: "N",
				DestCountry:     "PH",
				DestCurrency:    "PHP",
				TransactionType: "SO",
				PromoCode:       "",
				Message:         []string{},
				MessageLen:      "0",
				TerminalID:      "WBPt",
				OperatorID:      "001",
			},
			expectedReq: testFIRequest{
				WU: testFIWU{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testFIBody{
						Module:  "wuso",
						Request: "feeinquiry",
						Param: FIRequest{
							FrgnRefNo:       "123123",
							PrincipalAmount: "100000",
							FixedAmountFlag: "N",
							DestCountry:     "PH",
							DestCurrency:    "PHP",
							TransactionType: static.WUSendMoney,
							PromoCode:       "",
							Message:         []string{},
							MessageLen:      "0",
							TerminalID:      "WBPt",
							OperatorID:      "001",
						},
					},
					Signature: "placeholderSignature",
				},
			},
			want: &FIResponseBody{
				Taxes: FITaxes{
					MuniTax:      "0",
					StateTax:     "0",
					CountyTax:    "0",
					TaxWorksheet: "",
				},
				OrigPrincipal:    "100000",
				OrigCurrency:     "",
				DestPrincipal:    "100000",
				ExchangeRate:     "1",
				GrossTotal:       "101500",
				PlusCharges:      "0",
				PayAmount:        "100000",
				Charges:          "1500",
				Tolls:            "0",
				CNDExgFee:        "0",
				MessageCharge:    "0",
				PromoCodeDesc:    "",
				PromoSequenceNo:  "",
				PromoName:        "",
				PromoDiscountAmt: "0",
				BaseMsgCharge:    "5000",
				BaseMsgLimit:     "10",
				IncMsgCharge:     "500",
				IncMsgLimit:      "10",
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(feeInquirySuccessHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.WU.Header.Token

			got, err := s.FeeInquiry(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("FeeInquiry() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

func feeInquirySuccessHandler(t *testing.T, expectedReq testFIRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/1.1/wuso-feeinquiry" {
			t.Errorf("expected request to '/1.1/wuso-feeinquiry', got '%s'", req.URL.EscapedPath())
		}

		var newReq testFIRequest
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
				  "taxes":{
					 "municipal_tax":0,
					 "state_tax":0,
					 "county_tax":0,
					 "tax_worksheet":""
				  },
				  "originators_principal_amount":100000,
				  "destination_principal_amount":100000,
				  "gross_total_amount":101500,
				  "plus_charges_amount":0,
				  "pay_amount":100000,
				  "charges":1500,
				  "tolls":0,
				  "originating_currency_principal":"",
				  "canadian_dollar_exchange_fee":0,
				  "message_charge":0,
				  "promo_code_description":"",
				  "promo_sequence_no":"",
				  "promo_name":"",
				  "promo_discount_amount":0,
				  "exchange_rate":1,
				  "base_message_charge":"5000",
				  "base_message_limit":"10",
				  "incremental_message_charge":"500",
				  "incremental_message_limit":"10"
			   }
			}
		 }`)
	}
}

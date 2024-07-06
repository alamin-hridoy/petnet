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

type testQPVRequestBody struct {
	Module  string     `json:"module"`
	Request string     `json:"request"`
	Param   QPVRequest `json:"param"`
}

type testQPVRequestWU struct {
	Header    RequestHeader      `json:"header"`
	Body      testQPVRequestBody `json:"body"`
	Signature string             `json:"signature"`
}

type testQPVRequest struct {
	WU testQPVRequestWU `json:"uspwuapi"`
}

var qpvReq = QPVRequest{
	FrgnRefNo:               "cdffb28cea04d8ee9cb1",
	SenderNameType:          "D",
	SenderFirstName:         "WINNIE",
	SenderMiddleName:        "",
	SenderLastName:          "CONSTANTINO",
	SenderAddrLine1:         "1953 PH 3B BLOCK 6 LOT 9 CAMARIN",
	SenderAddrLine2:         "",
	SenderCity:              "CALOOCAN CITY",
	SenderState:             "METRO MANILA",
	SenderPostalCode:        "1400",
	SenderAddrCountry:       "PH",
	SenderAddrCurrency:      "PHP",
	SenderContactPhone:      "9089488895",
	SenderMobileCountryCode: "63",
	SenderMobileNo:          "9089488895",
	SenderAddrCountryName:   "Philippines",
	MyWUNumber:              "185015541",
	IDType:                  "A",
	IDCountry:               "Philippines",
	IDNumber:                "987654321",
	IDIssued:                "12012012",
	IDExpiry:                "12012019",
	Birthdate:               "12011996",
	Occupation:              "Airline Maritime Employee",
	BirthCountry:            "Philippines",
	Nationality:             "Philippines",
	Gender:                  "",
	FundSource:              "Salary Income",
	ReceiverRelationship:    "",
	MyWUEnrollTag:           "mywu",
	Email:                   "WINNIEANN.CONSTANTINO@VOXDEISYSTEMS.COM",
	DestCountry:             "PH",
	DestCurrency:            "PHP",
	DestState:               "",
	DestCity:                "",
	TransactionType:         "SO",
	PrincipalAmount:         "600000",
	FixedAmountFlag:         "N",
	PromoCode:               "",
	Message:                 []string{},
	TerminalID:              "PH259ART001A",
	OperatorID:              "5",
	RemoteTerminalID:        "",
	RemoteOperatorID:        "",
	EmploymentPositionLevel: "Entry Level",
	TransactionPurpose:      "",
	IsOnBehalf:              "",
	GalacticID:              StrPtr("1000000000028668380"),
	AckFlag:                 "X",
	CompanyName:             "TESTING ONLY",
	CompanyCode:             "1A IA",
	CompanyAccountCode:      "121212121212A",
	ReferenceNo:             "",
}

func TestQuickPayValidate(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          QPVRequest
		expectedReq testQPVRequest
		want        *QPVResponse
		wantErr     bool
	}{
		{
			name: "Success",
			in:   qpvReq,
			expectedReq: testQPVRequest{
				WU: testQPVRequestWU{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testQPVRequestBody{
						Module:  "wuqp",
						Request: "wuqp-validate",
						Param:   qpvReq,
					},
					Signature: "placeholderSignature",
				},
			},
			want: &QPVResponse{
				ServiceCode: QPVServiceCode{
					AddlServiceCharges: "1103FEE1205375001305375000103MSG02010030109001096150101002010030109703NNN9807PILQQCU",
					AddlServiceBlock: QPVAddlServiceBlock{
						AddlServiceLength:     377,
						AddlServiceDataBuffer: "",
					},
				},
				Compliance: QPVCompliance{
					ComplianceFlagsBuffer: "",
					ComplianceDataBuffer:  "0108UNI_01_S0201A03099876543210611Philippines0708120119960825Airline Maritime Employee0908120120141008120120192011Philippines29321953 PH 3B BLOCK 6 LOT 9 CAMARIN3013CALOOCAN CITY310414009901X3311Philippines3411Philippines4413Salary Income6724VOX DEI PROTOCOL SYSTEMS7401YB401157321953 PH 3B BLOCK 6 LOT 9 CAMARIN5913CALOOCAN CITY6012METRO MANILA610414006211PhilippinesF901NJ6191000000000028668380J701YM711Entry Level",
				},
				PaymentDetails: QPVPaymentDetails{
					OrigCity:  "METRO MANILA",
					OrigState: "MA",
				},
				Fin: Financials{
					Taxes: Taxes{
						MuniTax:   "0",
						StateTax:  "0",
						CountyTax: "0",
					},
					OrigPnplAmt: "600000",
					DestPcplAmt: "11225",
					GrossTotal:  "637500",
					AddlCharges: "0",
					Charges:     "37500",
				},
				Promotions: QPVPromotions{
					PromoDesc:       "",
					Message:         "",
					DiscountAmount:  "0",
					PromotionError:  "",
					SenderPromoCode: "",
				},
				NewDetails: QPVNewDetails{
					MTCN:       "6711768425",
					NewMTCN:    "1809486711768425",
					FilingDate: "04-04",
					FilingTime: "1158P EDT",
				},
				PreferredCustomer: QPVPreferredCustomer{
					MyWUNumber: "185015541",
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(successQPVHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.WU.Header.Token

			got, err := s.QuickPayValidate(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatal(err)
				t.Errorf("QuickPayValidate() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if !cmp.Equal(test.want, got) {
				t.Fatal(cmp.Diff(test.want, got))
			}
		})
	}
}

func successQPVHandler(t *testing.T, expectedReq testQPVRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/1.1/wuqp-validate" {
			t.Errorf("expected request to '/1.1/wuqp-validate', got '%s'", req.URL.EscapedPath())
		}

		var newReq testQPVRequest
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
					"service_code":{
						"addl_service_charges":"1103FEE1205375001305375000103MSG02010030109001096150101002010030109703NNN9807PILQQCU",
						"addl_service_block":{
							"addl_service_length":377,
							"addl_service_data_buffer":""
						}
					},
					"compliance":{
						"compliance_flags_buffer":"",
						"compliance_data_buffer":"0108UNI_01_S0201A03099876543210611Philippines0708120119960825Airline Maritime Employee0908120120141008120120192011Philippines29321953 PH 3B BLOCK 6 LOT 9 CAMARIN3013CALOOCAN CITY310414009901X3311Philippines3411Philippines4413Salary Income6724VOX DEI PROTOCOL SYSTEMS7401YB401157321953 PH 3B BLOCK 6 LOT 9 CAMARIN5913CALOOCAN CITY6012METRO MANILA610414006211PhilippinesF901NJ6191000000000028668380J701YM711Entry Level"
					},
					"payment_details":{
						"originating_city":"METRO MANILA",
						"originating_state":"MA"
					},
					"financials":{
						"taxes":{
							"municipal_tax":0,
							"state_tax":0,
							"county_tax":0
						},
						"originators_principal_amount":600000,
						"destination_principal_amount":11225,
						"gross_total_amount":637500,
						"plus_charges_amount":0,
						"charges":37500
					},
					"promotions":{
						"promo_code_description":"",
						"promo_message":"",
						"promo_discount_amount":0,
						"promotion_error":"",
						"sender_promo_code":""
					},
					"new_details":{
						"mtcn":"6711768425",
						"new_mtcn":"1809486711768425",
						"filing_date":"04-04",
						"filing_time":"1158P EDT"
					},
					"preferred_customer":{
						"mywu_number":"185015541"
					}
				}
			}
		}`)
	}
}

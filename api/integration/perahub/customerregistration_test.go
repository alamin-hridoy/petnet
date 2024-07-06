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

type testCustomerRegRequestBody struct {
	Module  string             `json:"module"`
	Request string             `json:"request"`
	Param   CustomerRegRequest `json:"param"`
}

type testCustomerRegRequestWU struct {
	Header    RequestHeader              `json:"header"`
	Body      testCustomerRegRequestBody `json:"body"`
	Signature string                     `json:"signature"`
}

type testCustomerRegRequest struct {
	WU testCustomerRegRequestWU `json:"uspwuapi"`
}

var crReq = CustomerRegRequest{
	Username:            "wc_ac97",
	Mobile:              "09057649300",
	Email:               "jeanandrew.fuentes@whitecloak.com",
	Phone:               "09057649300",
	Surname:             "fuentes",
	Givenname:           "jean",
	Middlename:          "andrew",
	Password:            "f33fdc2ec2a654ba0af87c318770ef5eba1b752b18560ee96220e9813403c29c52018f1e0f7daf3598147f87fa510065217bfcb729e8cd5b83c847876a6c25d0",
	Birthdate:           "06/12/1994",
	Nationality:         "Angola",
	Occupation:          "12312312",
	NameOfEmployer:      "White Cloak",
	SecurityQuestion1:   "V2hhdCBpcyBteSBzcG91c2Uncy9wYXJ0bmVyJ3MgbWlkZGxlIG5hbWU",
	Answer1:             "A450D42CA6F5A5EEB6D0ED79B01129F96CAC4B82FD059D14C9837FF4A971A1AEA3CC68CB02F37E2B43631C99953A9490E39F5CEB4C847EE797638078231E5D94",
	SecurityQuestion2:   "SW4gd2hpY2ggY2l0eSB3YXMgbXkgbW90aGVyIGJvcm4",
	Answer2:             "C611F7EBDDB95AEEAB4DA7AB7B5283FD413D0C571455C5A8B2CE99DA6A96C6B4E3FFDAF45A14BA254679E58DD5FADD1E6CEF221A46ABCC6F60DA3EED1E60F5E1",
	SecurityQuestion3:   "V2hhdCBpcyB0aGUgbmFtZSBåvZiB0aGUgZmlyc3Qgc2Nob29sIEkgYXR0ZW5kZWQ",
	Answer3:             "A450D42CA6F5A5EEB6D0ED79B01129F96CAC4B82FD059D14C9837FF4A971A1AEA3CC68CB02F37E2B43631C99953A9490E39F5CEB4C847EE797638078231E5D94",
	ValidIdentification: "015",
	IDImage:             "",
	WuCardNo:            "",
	DebitCardNo:         "",
	LoyaltyCardNo:       "",
	CustomerIDNumber:    1231312,
	IDCountryIssue:      "Aland Islands",
	IDIssueDate:         "01/01/1700",
	Gender:              "N",
	City:                "asdas",
	State:               "asdsadas",
	PostalCode:          4000,
	IsWalkIn:            "1",
	CustomerCode:        "123456789",
	IDType:              "O",
	IDExpirationDate:    "10 05 2017",
	CountryOfBirth:      "Antarctica",
	TIN:                 "none",
	SSS:                 13,
	PresentAddress:      "asdadas",
	PresentCity:         "asdas",
	PresentState:        "asdsadas",
	PresentProvince:     "asdsadas",
	PresentRegion:       "asdsadas",
	PresentCountry:      "AF",
	PresentPostalcode:   4000,
	PermanentAddress:    "adsad",
	PermanentCity:       "dsadsad",
	PermanentState:      "asdsadas",
	PermanentProvince:   "asdsadas",
	PermanentRegion:     "asdsadas",
	PermanentCountry:    "XP",
	PermanentPostalcode: 4000,
	ACountry:            "AF",
}

func TestCustomerRegistration(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          CustomerRegRequest
		expectedReq testCustomerRegRequest
		want        *CustomerRegResponse
		wantErr     bool
	}{
		{
			name: "Success",
			in:   crReq,
			expectedReq: testCustomerRegRequest{
				WU: testCustomerRegRequestWU{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testCustomerRegRequestBody{
						Module:  "Register",
						Request: "Register",
						Param:   crReq,
					},
					Signature: "placeholderSignature",
				},
			},
			want: &CustomerRegResponse{
				WU: CustomerRegResponseWU{
					Header: ResponseHeader{
						ErrorCode: "1",
						Message:   "good",
					},
					Body: CustomerRegResponseBody{
						Code:         "1",
						Key:          "489e0e74-8b56-4bb2-a58a-01a42e114d0f",
						FType:        "CreateSession",
						CustomerCode: "JF061294UC7V",
					},
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(successCustomerRegHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.WU.Header.Token

			got, err := s.CustomerRegistration(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatal(err)
				t.Errorf("CustomerRegistration() error = %v, wantErr %v", err, test.wantErr)
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

func successCustomerRegHandler(t *testing.T, expectedReq testCustomerRegRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/Register" {
			t.Errorf("expected request to '/Register', got '%s'", req.URL.EscapedPath())
		}

		var newReq testCustomerRegRequest
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
					"code": "1",
					"Key": "489e0e74-8b56-4bb2-a58a-01a42e114d0f",
					"fType": "CreateSession",
					"customer_code ": "JF061294UC7V"
				}
			}
		}`)
	}
}

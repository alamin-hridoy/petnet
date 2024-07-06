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

type testCustomerIURequestBody struct {
	Module  string            `json:"module"`
	Request string            `json:"request"`
	Param   CustomerIURequest `json:"param"`
}

type testCustomerIURequestWU struct {
	Header    RequestHeader             `json:"header"`
	Body      testCustomerIURequestBody `json:"body"`
	Signature string                    `json:"signature"`
}

type testCustomerIURequest struct {
	WU testCustomerIURequestWU `json:"uspwuapi"`
}

var ciuReq = CustomerIURequest{
	Nationality:         "Filipino",
	Usercode:            "QS060699MEFH",
	Surname:             "Six",
	Givenname:           "QA",
	Middlename:          "Yondu",
	Mobile:              "09166585050",
	Email:               "qsix@mailinator.com",
	NewEmail:            "qsix@mailinator.com",
	Presentaddress:      "Panorama",
	Permanentaddress:    "Taguig City",
	Secretquestion1:     "V2hhdCBpcyB0aGUgbmFtZSBvZiB0aGUgZmlyc3Qgc2Nob29sIEkgYXR0ZW5kZWQ",
	Answer1:             "047B10FE577A23EFD96546DCFCE8485FC4AA8AE84DD3BF0C435A294CF318C7A260418DD96A97FEB0AD7AED90FF011620CFA5D7B3CDC8AEA4C4E81E56A0FC9934",
	Secretquestion2:     "V2hhdCB3YXMgdGhlIG5hbWUgb2YgeW91ciBlbGVtZW50YXJ5IHNjaG9vbD8=",
	Answer2:             "3C9909AFEC25354D551DAE21590BB26E38D53F2173B8D3DC3EEE4C047E7AB1C1EB8B85103E3BE7BA613B31BB5C9C36214DC9F14A42FD7A2FDB84856BCA5C44C2",
	Secretquestion3:     "V2hhdCB3YXMgdGhlIGZpcnN0IGNvbXBhbnkgSSBldmVyIHdvcmtlZCBmb3I",
	Answer3:             "21B4F4BD9E64ED355C3EB676A28EBEDAF6D8F17BDC365995B319097153044080516BD083BFCCE66121A3072646994C8430CC382B8DC543E84880183BF856CFF5",
	Presentcity:         "Taguig City",
	Presentstate:        "Metro Manila",
	Presentprovince:     "Metro Manila",
	Presentregion:       "Metro Manila",
	Presentcountry:      "US",
	Presentpostalcode:   "1630",
	Permanentcity:       "none",
	Permanentstate:      "Metro Manila",
	Permanentprovince:   "Metro Manila",
	Permanentregion:     "none",
	Permanentcountry:    "none",
	Permanentpostalcode: "none",
	IDType:              "L",
	IDImage:             "",
	IDExpirationDate:    "09 30 2017",
	CustomerIDNumber:    "123456789",
	IDCountryIssue:      "XP",
	ProfileImage:        []byte(""),
}

func TestCustomerInfoUpdate(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          CustomerIURequest
		expectedReq testCustomerIURequest
		want        *CustomerIUResponse
		wantErr     bool
	}{
		{
			name: "Success",
			in:   ciuReq,
			expectedReq: testCustomerIURequest{
				WU: testCustomerIURequestWU{
					Header: RequestHeader{
						Coy:          "yondu",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testCustomerIURequestBody{
						Module:  "UpdateInfo",
						Request: "update_info",
						Param:   ciuReq,
					},
					Signature: "placeholderSignature",
				},
			},
			want: &CustomerIUResponse{
				WU: CustomerIUResponseWU{
					Header: ResponseHeader{
						ErrorCode: "1",
						Message:   "good",
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(successCustomerIUHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.WU.Header.Token
			// s.signKey = test.expectedReq.WU.Signature

			got, err := s.CustomerInfoUpdate(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatal(err)
				t.Errorf("CustomerInfoUpdate() error = %v, wantErr %v", err, test.wantErr)
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

func successCustomerIUHandler(t *testing.T, expectedReq testCustomerIURequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/UpdateInfo" {
			t.Errorf("expected request to '/UpdateInfo', got '%s'", req.URL.EscapedPath())
		}

		var newReq testCustomerIURequest
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
				"body": null
			}
		}`)
	}
}

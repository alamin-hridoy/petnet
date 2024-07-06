package perahub

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type testForgotPassResetBody struct {
	Module  string                 `json:"module"`
	Request string                 `json:"request"`
	Param   ForgotPassResetRequest `json:"param"`
}

type testForgotPassResetData struct {
	Header    RequestHeader           `json:"header"`
	Body      testForgotPassResetBody `json:"body"`
	Signature string                  `json:"signature"`
}

type testForgotPassResetRequest struct {
	Data testForgotPassResetData `json:"uspwuapi"`
}

var fprReq = ForgotPassResetRequest{
	UserName: "reymar0524",
	Password: "secret",
	Category: "login",
	OTPCode:  "256191",
}

var fprResBody = ForgotPassResetResponseBody{
	ForeignRefNo:        "f83346f2dcfe4a169d54a3b5855b896f",
	FirstName:           "Reymar",
	LastName:            "Fadriquela",
	MiddleName:          "gs",
	CustomerCode:        "DA021692XRZ2",
	Mobile:              "09175337900",
	Email:               "a@y.com",
	Birthdate:           "02/16/1992",
	Nationality:         "Philippines",
	PresentAddress:      "hdhd",
	PermanentAddress:    "hdhd",
	Occupation:          "gsys",
	NameOfEmployer:      "",
	Wucardno:            "",
	DebitCardNo:         "060000309793",
	LoyaltyCardNo:       "",
	PasswordExpiration:  "2019-04-30 17:04:21",
	CardPoints:          0,
	CustomerIDNo:        "6564",
	CountryIDIssue:      "Philippines",
	IDIssueDate:         "01/01/1700",
	Gender:              "N",
	City:                "hdh",
	State:               "hdhd",
	PostalCode:          "64",
	IDType:              "D",
	Acountry:            "AS",
	CountryOfBirth:      "Philippines",
	IDExpiration:        "12/31/2999",
	Tin:                 "none",
	Sss:                 "none",
	SecretQuestion1:     "V2hhdCBpcyB0aGUgbmFtZSBvZiB0aGUgZmlyc3Qgc2Nob29sIEkgYXR0ZW5kZWQ/",
	Answer1:             "2A7DEBE713CACC8C99D0C38EDFB852E82927489EC81D2B2D1A02E99B6CCF83DFE792B79144CEA03486735BC3E675E7F9311ADF02A119BF68377112A4EF29E74C",
	SecretQuestion2:     "V2hhdCB3YXMgdGhlIGZpcnN0IGNvbXBhbnkgSSBldmVyIHdvcmtlZCBmb3I/",
	Answer2:             "5EFFBE7C8675304032269A719E547467D15BE5A76BC8227ED4A3D92949876FC08BFFCBC810CF9B4565B975EC7ECF3993A7EDF330662F17F15CF041CA5A2B40E7",
	SecretQuestion3:     "V2hhdCBpcyBteSBzcG91c2Uncy9wYXJ0bmVyJ3MgbWlkZGxlIG5hbWU/",
	Answer3:             "C9BD429C7DF48F1BDFCF1BF7ACC1B7B816A424601C52F05917931DA10DA58124FCA1FA1948274A31D2F1B3ED40EA080F83B5D0CFB755BC7D753298D85B476A9D",
	PresentCity:         "hdh",
	PresentState:        "hdhd",
	PresentProvince:     "hdhd",
	PresentRegion:       "hdhd",
	PresentCountry:      "AS",
	PresentPostalCode:   "64",
	PermanentCity:       "gss",
	PermanentState:      "hdhd",
	PermanentProvince:   "hdhd",
	PermanentRegion:     "AS",
	PermanentCountry:    "AS",
	PermanentPostalCode: "67",
	IdImage:             "data:image/jpeg;base64",
	IsTemporaryPassword: "1",
	ProfileImage:        nil,
	IDRows:              []IDRows{},
}

func TestForgotPassResetword(t *testing.T) {
	t.Parallel()
	// baseUrl := "http://kycdevgateway.perahub.com.ph/gateway"
	partnerID := "client_01"
	clientKey := "de9462831b"
	serverIP := "127.0.0.1"

	tests := []struct {
		name        string
		in          ForgotPassResetRequest
		expectedReq testForgotPassResetRequest
		want        *ForgotPassResetResponse
		wantErr     bool
	}{
		{
			name: "Success",
			in:   fprReq,
			expectedReq: testForgotPassResetRequest{
				Data: testForgotPassResetData{
					Header: RequestHeader{
						Coy:          "usp",
						Token:        "test-token",
						LocationCode: "1",
						UserCode:     "1",
						ClientIP:     "127.0.0.1",
						IsWeb:        "1",
					},
					Body: testForgotPassResetBody{
						Module:  "forgot_pwd_commit",
						Request: "forgot_pwd_commit",
						Param:   fprReq,
					},
					Signature: "placeholderSignature",
				},
			},
			want: &ForgotPassResetResponse{
				Data: ForgotPassResetResponseData{
					Header: ResponseHeader{
						ErrorCode: "1",
						Message:   "Good",
					},
					Body: fprResBody,
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(ForgotPassResetHandler(t, test.expectedReq))
			t.Cleanup(func() { ts.Close() })

			s, err := New(nil, "dev", ts.URL, "", "", "", "", partnerID, clientKey, "api-key", serverIP, "", nil)
			if err != nil {
				t.Fatal(err)
			}

			s.token = test.expectedReq.Data.Header.Token

			got, err := s.ForgotPasswordReset(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("ForgotPasswordReset() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if got.Data.Header.ErrorCode != "1" {
				t.Errorf("want error code 1 but got %q", got.Data.Header.ErrorCode)
			}

			if !cmp.Equal(test.want, got) {
				t.Fatal(cmp.Diff(test.want, got))
			}
		})
	}
}

func ForgotPassResetHandler(t *testing.T, expectedReq testForgotPassResetRequest) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected 'POST' request, got '%s'", req.Method)
		}
		if req.URL.EscapedPath() != "/forgot_pwd_commit" {
			t.Errorf("expected request to '/forgot_pwd_commit', got '%s'", req.URL.EscapedPath())
		}

		var newReq testForgotPassResetRequest
		if err := json.NewDecoder(req.Body).Decode(&newReq); err != nil {
			t.Fatal(err)
		}
		hashByte := sha512.New().Sum([]byte(expectedReq.Data.Body.Param.Password))

		reqHashByte, err := hex.DecodeString(newReq.Data.Body.Param.Password)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(hashByte, reqHashByte) {
			t.Error(cmp.Diff(hashByte, reqHashByte))
		}

		opts := cmpopts.IgnoreFields(testForgotPassResetRequest{}, "Data.Body.Param.Password")
		if !cmp.Equal(expectedReq, newReq, opts) {
			t.Error(cmp.Diff(expectedReq, newReq))
		}

		res.WriteHeader(200)
		res.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(res, `{
			"uspwuapi": {
				"header": {
					"errorcode": "1",
					"message": "Good"
				},
				"body": {
					"foreign_reference_no": "f83346f2dcfe4a169d54a3b5855b896f",
					"first_name": "Reymar",
					"last_name": "Fadriquela",
					"middle_name": "gs",
					"customer_code": "DA021692XRZ2",
					"mobile": "09175337900",
					"email": "a@y.com",
					"birthdate": "02/16/1992",
					"nationality": "Philippines",
					"presentaddress": "hdhd",
					"permanentaddress": "hdhd",
					"occupation": "gsys",
					"nameofemployer": "",
					"wucardno": "",
					"debitcardno": "060000309793",
					"loyaltycardno": "",
					"password_expiration": "2019-04-30 17:04:21",
					"card_points": 0,
					"customer_id_no": "6564",
					"country_id_issue": "Philippines",
					"id_issue_date": "01/01/1700",
					"gender": "N",
					"city": "hdh",
					"state": "hdhd",
					"postal_code": "64",
					"id_type": "D",
					"acountry": "AS",
					"country_of_birth": "Philippines",
					"id_expiration": "12/31/2999",
					"tin": "none",
					"sss": "none",
					"secretquestion1":
					"V2hhdCBpcyB0aGUgbmFtZSBvZiB0aGUgZmlyc3Qgc2Nob29sIEkgYXR0ZW5kZWQ/",
					"answer1": "2A7DEBE713CACC8C99D0C38EDFB852E82927489EC81D2B2D1A02E99B6CCF83DFE792B79144CEA03486735BC3E675E7F9311ADF02A119BF68377112A4EF29E74C",
					"secretquestion2": "V2hhdCB3YXMgdGhlIGZpcnN0IGNvbXBhbnkgSSBldmVyIHdvcmtlZCBmb3I/",
					"answer2": "5EFFBE7C8675304032269A719E547467D15BE5A76BC8227ED4A3D92949876FC08BFFCBC810CF9B4565B975EC7ECF3993A7EDF330662F17F15CF041CA5A2B40E7",
					"secretquestion3": "V2hhdCBpcyBteSBzcG91c2Uncy9wYXJ0bmVyJ3MgbWlkZGxlIG5hbWU/",
					"answer3": "C9BD429C7DF48F1BDFCF1BF7ACC1B7B816A424601C52F05917931DA10DA58124FCA1FA1948274A31D2F1B3ED40EA080F83B5D0CFB755BC7D753298D85B476A9D",
					"presentcity": "hdh",
					"presentstate": "hdhd",
					"presentprovince": "hdhd",
					"presentregion": "hdhd",
					"presentcountry": "AS",
					"presentpostalcode": "64",
					"permanentcity": "gss",
					"permanentstate": "hdhd",
					"permanentprovince": "hdhd",
					"permanentregion": "AS",
					"permanentcountry": "AS",
					"permanentpostalcode": "67",
					"IdImage": "data:image/jpeg;base64",
					"isTemporaryPassword": "1",
					"ProfileImage": null,
					"idrows": []
				}
			}
		}`)
	}
}

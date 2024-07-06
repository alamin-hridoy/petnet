package perahub

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ForgotPassResetRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Category string `json:"category"`
	OTPCode  string `json:"otp_code"`
}

type IDRows interface{}

type ForgotPassResetResponseBody struct {
	ForeignRefNo        string   `json:"foreign_reference_no"`
	FirstName           string   `json:"first_name"`
	LastName            string   `json:"last_name"`
	MiddleName          string   `json:"middle_name"`
	CustomerCode        string   `json:"customer_code"`
	Mobile              string   `json:"mobile"`
	Email               string   `json:"email"`
	Birthdate           string   `json:"birthdate"`
	Nationality         string   `json:"nationality"`
	PresentAddress      string   `json:"presentaddress"`
	PermanentAddress    string   `json:"permanentaddress"`
	Occupation          string   `json:"occupation"`
	NameOfEmployer      string   `json:"nameofemployer"`
	Wucardno            string   `json:"wucardno"`
	DebitCardNo         string   `json:"debitcardno"`
	LoyaltyCardNo       string   `json:"loyaltycardno"`
	PasswordExpiration  string   `json:"password_expiration"`
	CardPoints          int      `json:"card_points"`
	CustomerIDNo        string   `json:"customer_id_no"`
	CountryIDIssue      string   `json:"country_id_issue"`
	IDIssueDate         string   `json:"id_issue_date"`
	Gender              string   `json:"gender"`
	City                string   `json:"city"`
	State               string   `json:"state"`
	PostalCode          string   `json:"postal_code"`
	IDType              string   `json:"id_type"`
	Acountry            string   `json:"acountry"`
	CountryOfBirth      string   `json:"country_of_birth"`
	IDExpiration        string   `json:"id_expiration"`
	Tin                 string   `json:"tin"`
	Sss                 string   `json:"sss"`
	SecretQuestion1     string   `json:"secretquestion1"`
	Answer1             string   `json:"answer1"`
	SecretQuestion2     string   `json:"secretquestion2"`
	Answer2             string   `json:"answer2"`
	SecretQuestion3     string   `json:"secretquestion3"`
	Answer3             string   `json:"answer3"`
	PresentCity         string   `json:"presentcity"`
	PresentState        string   `json:"presentstate"`
	PresentProvince     string   `json:"presentprovince"`
	PresentRegion       string   `json:"presentregion"`
	PresentCountry      string   `json:"presentcountry"`
	PresentPostalCode   string   `json:"presentpostalcode"`
	PermanentCity       string   `json:"permanentcity"`
	PermanentState      string   `json:"permanentstate"`
	PermanentProvince   string   `json:"permanentprovince"`
	PermanentRegion     string   `json:"permanentregion"`
	PermanentCountry    string   `json:"permanentcountry"`
	PermanentPostalCode string   `json:"permanentpostalcode"`
	IdImage             string   `json:"IdImage"`
	IsTemporaryPassword string   `json:"isTemporaryPassword"`
	ProfileImage        []byte   `json:"ProfileImage"`
	IDRows              []IDRows `json:"idrows"`
}

type ForgotPassResetResponseData struct {
	Header ResponseHeader              `json:"header"`
	Body   ForgotPassResetResponseBody `json:"body"`
}

type ForgotPassResetResponse struct {
	Data ForgotPassResetResponseData `json:"uspwuapi"`
}

func (s *Svc) ForgotPasswordReset(ctx context.Context, pr ForgotPassResetRequest) (*ForgotPassResetResponse, error) {
	pr.Password = fmt.Sprintf("%x", sha512.New().Sum([]byte(pr.Password)))

	req, err := s.newParahubRequest(ctx, "forgot_pwd_commit", "forgot_pwd_commit", pr)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := s.cl.Post(s.moduleURL("forgot_pwd_commit", ""), "aplication/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var rpRes ForgotPassResetResponse
	if err := json.Unmarshal(body, &rpRes); err != nil {
		return nil, err
	}

	return &rpRes, nil
}

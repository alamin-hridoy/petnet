package perahub

import (
	"context"
	"crypto/sha512"
	"encoding/json"
	"fmt"

	"brank.as/petnet/api/storage"
	"brank.as/petnet/serviceutil/logging"
)

// Type conversion check
var _ storage.Customer = storage.Customer(Customer{})

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Category string `json:"category"`
}

type Customer struct {
	FrgnRefNo    string `json:"foreign_reference_no"`
	CustomerCode string `json:"customer_code"`

	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`

	Phone  string `json:"phone"`
	Mobile string `json:"mobile"`
	Email  string `json:"email"`

	Birthdate    string `json:"birthdate"`
	BirthCountry string `json:"country_of_birth"`
	Nationality  string `json:"nationality"`
	Gender       string `json:"gender"`

	Occupation   string `json:"occupation"`
	EmployerName string `json:"nameofemployer"`

	Wucardno      string `json:"wucardno"`
	Debitcardno   string `json:"debitcardno"`
	Loyaltycardno string `json:"loyaltycardno"`

	IsTemporaryPassword string `json:"isTemporaryPassword"`
	PasswordExpiration  string `json:"password_expiration"`

	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`

	CardPoints        json.Number `json:"card_points"`
	OneSignalPlayerID string      `json:"one_signal_player_id"`
	Tin               string      `json:"tin"`
	SSS               string      `json:"sss"`

	Secretquestion1 string `json:"secretquestion1"`
	Answer1         string `json:"answer1"`
	Secretquestion2 string `json:"secretquestion2"`
	Answer2         string `json:"answer2"`
	Secretquestion3 string `json:"secretquestion3"`
	Answer3         string `json:"answer3"`

	PresentAddress    string `json:"presentaddress"`
	PresentCity       string `json:"presentcity"`
	PresentState      string `json:"presentstate"`
	PresentProvince   string `json:"presentprovince"`
	PresentRegion     string `json:"presentregion"`
	PresentCountry    string `json:"presentcountry"`
	PresentPostalCode string `json:"presentpostalcode"`

	PermanentAddress    string `json:"permanentaddress"`
	PermanentCity       string `json:"permanentcity"`
	PermanentState      string `json:"permanentstate"`
	PermanentProvince   string `json:"permanentprovince"`
	PermanentRegion     string `json:"permanentregion"`
	PermanentCountry    string `json:"permanentcountry"`
	PermanentPostalCode string `json:"permanentpostalcode"`

	ProfileImage string `json:"ProfileImage"`
	IDImage      string `json:"IdImage"`

	IDType         string `json:"id_type"`
	CustomerIDNo   string `json:"customer_id_no"`
	CountryIDIssue string `json:"country_id_issue"`
	IDIssueDate    string `json:"id_issue_date"`
	IDExpiration   string `json:"id_expiration"`
	Acountry       string `json:"acountry"`

	Idrows []storage.IDEntry `json:"idrows"`
}

func (s *Svc) Login(ctx context.Context, lr LoginRequest) (*Customer, error) {
	r := loginRequest{
		Username: lr.Username,
		Password: fmt.Sprintf("%X", sha512.New().Sum([]byte(lr.Password))),
		Category: "login",
	}
	const mod, reqMod, modReq = "SignOn", "signin", "login"
	req, err := s.newParahubRequest(ctx, mod, modReq, r)
	if err != nil {
		return nil, err
	}

	resp, err := s.post(ctx, s.moduleURL(reqMod, ""), *req)
	if err != nil {
		return nil, err
	}

	res := Customer{}
	if err := json.Unmarshal(resp, &res); err != nil {
		logging.FromContext(ctx).WithField("body", string(resp)).Error("unmarshal failed")
		return nil, err
	}
	return &res, nil
}

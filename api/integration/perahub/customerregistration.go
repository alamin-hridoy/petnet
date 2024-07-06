package perahub

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
)

type CustomerRegRequest struct {
	Username            string `json:"username"`
	Mobile              string `json:"mobile"`
	Email               string `json:"email"`
	Phone               string `json:"phone"`
	Surname             string `json:"surname"`
	Givenname           string `json:"givenname"`
	Middlename          string `json:"middlename"`
	Password            string `json:"password"`
	Birthdate           string `json:"Birthdate"`
	Nationality         string `json:"Nationality"`
	Occupation          string `json:"Occupation"`
	NameOfEmployer      string `json:"NameOfEmployer"`
	SecurityQuestion1   string `json:"SecurityQuestion1"`
	Answer1             string `json:"Answer1"`
	SecurityQuestion2   string `json:"SecurityQuestion2"`
	Answer2             string `json:"Answer2"`
	SecurityQuestion3   string `json:"SecurityQuestion3"`
	Answer3             string `json:"Answer3"`
	ValidIdentification string `json:"ValidIdentification"`
	IDImage             string `json:"IdImage"`
	WuCardNo            string `json:"WuCardNo"`
	DebitCardNo         string `json:"DebitCardNo"`
	LoyaltyCardNo       string `json:"LoyaltyCardNo"`
	CustomerIDNumber    int    `json:"CustomerIdNumber"`
	IDCountryIssue      string `json:"IdCountryIssue"`
	IDIssueDate         string `json:"IdIssueDate"`
	Gender              string `json:"Gender"`
	City                string `json:"City"`
	State               string `json:"State"`
	PostalCode          int    `json:"PostalCode"`
	IsWalkIn            string `json:"IsWalkIn"`
	CustomerCode        string `json:"customer_code"`
	IDType              string `json:"IdType"`
	IDExpirationDate    string `json:"IdExpirationDate"`
	CountryOfBirth      string `json:"CountryOfBirth"`
	TIN                 string `json:"TIN"`
	SSS                 int    `json:"SSS"`
	PresentAddress      string `json:"PresentAddress"`
	PresentCity         string `json:"PresentCity"`
	PresentState        string `json:"PresentState"`
	PresentProvince     string `json:"PresentProvince"`
	PresentRegion       string `json:"PresentRegion"`
	PresentCountry      string `json:"PresentCountry"`
	PresentPostalcode   int    `json:"PresentPostalcode"`
	PermanentAddress    string `json:"PermanentAddress"`
	PermanentCity       string `json:"PermanentCity"`
	PermanentState      string `json:"PermanentState"`
	PermanentProvince   string `json:"PermanentProvince"`
	PermanentRegion     string `json:"PermanentRegion"`
	PermanentCountry    string `json:"PermanentCountry"`
	PermanentPostalcode int    `json:"PermanentPostalcode"`
	ACountry            string `json:"ACountry"`
}

type CustomerRegResponseBody struct {
	Code         string `json:"code"`
	Key          string `json:"Key"`
	FType        string `json:"fType"`
	CustomerCode string `json:"customer_code "`
}

type CustomerRegResponseWU struct {
	Header ResponseHeader          `json:"header"`
	Body   CustomerRegResponseBody `json:"body"`
}

type CustomerRegResponse struct {
	WU CustomerRegResponseWU `json:"uspwuapi"`
}

func (s *Svc) CustomerRegistration(ctx context.Context, crReq CustomerRegRequest) (*CustomerRegResponse, error) {
	req, err := s.newParahubRequest(ctx, "Register", "Register", crReq)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := s.cl.Post(s.moduleURL("Register", ""), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var crRes CustomerRegResponse
	if err := json.Unmarshal(body, &crRes); err != nil {
		return nil, err
	}

	return &crRes, nil
}

package perahub

import (
	"context"
	"encoding/json"
)

const (
	AwaitPayment = "W/C"
)

type WUSearchRequest struct {
	ForRefNo    string `json:"foreign_reference_no"`
	SearchParam string `json:"searchparam"`
	SearchType  string `json:"searchtype"`
	SenderFName string `json:"sender_first_name"`
	SenderLName string `json:"sender_last_name"`
	OperatorID  string `json:"operator_id"`
	TerminalID  string `json:"terminal_id"`
}

type ResponseDetails struct {
	ForRefNo            string `json:"-"`
	Surname             string `json:"surname"`
	GivenName           string `json:"givenname"`
	Birthdate           string `json:"Birthdate"`
	Nationality         string `json:"Nationality"`
	PresentAddress      string `json:"PresentAddress"`
	Occupation          string `json:"Occupation"`
	NameOfEmployer      string `json:"NameOfEmployer"`
	ValidIdentification string `json:"ValidIdentification"`
	WuCardNo            string `json:"WuCardNo"`
	DebitCardNo         string `json:"DebitCardNo"`
	LoyaltyCardNo       string `json:"LoyaltyCardNo"`
}

type WUSearchResponseBody struct {
	ForRefNo        string          `json:"foreign_reference_no"`
	ResponseDetails ResponseDetails `json:"customer"`
}

func (s *Svc) MyWUSearch(ctx context.Context, sr WUSearchRequest) (*ResponseDetails, error) {
	const mod, modReq = "prereq", "mywusearch"
	req, err := s.newParahubRequest(ctx, mod, modReq, sr)
	if err != nil {
		return nil, err
	}

	resp, err := s.post(ctx, s.moduleURL(mod, ""), *req)
	if err != nil {
		return nil, err
	}

	var res WUSearchResponseBody
	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, err
	}

	cust := res.ResponseDetails
	cust.ForRefNo = res.ForRefNo
	return &cust, nil
}

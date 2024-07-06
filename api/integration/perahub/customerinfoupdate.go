package perahub

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
)

type CustomerIURequest struct {
	Nationality         string `json:"nationality"`
	Usercode            string `json:"usercode"`
	Surname             string `json:"surname"`
	Givenname           string `json:"givenname"`
	Middlename          string `json:"middlename"`
	Mobile              string `json:"mobile"`
	Email               string `json:"email"`
	NewEmail            string `json:"new_email"`
	Presentaddress      string `json:"presentaddress"`
	Permanentaddress    string `json:"permanentaddress"`
	Secretquestion1     string `json:"secretquestion1"`
	Answer1             string `json:"answer1"`
	Secretquestion2     string `json:"secretquestion2"`
	Answer2             string `json:"answer2"`
	Secretquestion3     string `json:"secretquestion3"`
	Answer3             string `json:"answer3"`
	Presentcity         string `json:"presentcity"`
	Presentstate        string `json:"presentstate"`
	Presentprovince     string `json:"presentprovince"`
	Presentregion       string `json:"presentregion"`
	Presentcountry      string `json:"presentcountry"`
	Presentpostalcode   string `json:"presentpostalcode"`
	Permanentcity       string `json:"permanentcity"`
	Permanentstate      string `json:"permanentstate"`
	Permanentprovince   string `json:"permanentprovince"`
	Permanentregion     string `json:"permanentregion"`
	Permanentcountry    string `json:"permanentcountry"`
	Permanentpostalcode string `json:"permanentpostalcode"`
	IDType              string `json:"IdType"`
	IDImage             string `json:"IdImage"`
	IDExpirationDate    string `json:"IdExpirationDate"`
	CustomerIDNumber    string `json:"CustomerIdNumber"`
	IDCountryIssue      string `json:"IdCountryIssue"`
	ProfileImage        []byte `json:"ProfileImage"`
}

type CustomerIUResponseBody struct{}

type CustomerIUResponseWU struct {
	Header ResponseHeader         `json:"header"`
	Body   CustomerIUResponseBody `json:"body"`
}

type CustomerIUResponse struct {
	WU CustomerIUResponseWU `json:"uspwuapi"`
}

func (s *Svc) CustomerInfoUpdate(ctx context.Context, ciuReq CustomerIURequest) (*CustomerIUResponse, error) {
	req, err := s.newParahubRequest(ctx, "UpdateInfo", "update_info", ciuReq)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := s.cl.Post(s.moduleURL("UpdateInfo", ""), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var crRes CustomerIUResponse
	if err := json.Unmarshal(body, &crRes); err != nil {
		return nil, err
	}

	return &crRes, nil
}

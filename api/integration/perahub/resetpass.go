package perahub

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ResetPassRequest struct {
	UserName       string `json:"username"`
	Password       string `json:"password"`
	SecretQuestion string `json:"secret_question"`
	Answer         string `json:"answer"`
}

type ResetPassResponseBody struct {
	UserName string `json:"username"`
}

type ResetPassResponseData struct {
	Header ResponseHeader        `json:"header"`
	Body   ResetPassResponseBody `json:"body"`
}

type ResetPassResponse struct {
	Data ResetPassResponseData `json:"uspwuapi"`
}

func (s *Svc) ResetPassword(ctx context.Context, pr ResetPassRequest) (*ResetPassResponse, error) {
	pr.Password = fmt.Sprintf("%x", sha512.New().Sum([]byte(pr.Password)))
	pr.Answer = fmt.Sprintf("%x", sha512.New().Sum([]byte(pr.Answer)))

	req, err := s.newParahubRequest(ctx, "ResetPassword", "reset_password", pr)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := s.cl.Post(s.moduleURL("ResetPassword", ""), "aplication/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var rpRes ResetPassResponse
	if err := json.Unmarshal(body, &rpRes); err != nil {
		return nil, err
	}

	return &rpRes, nil
}

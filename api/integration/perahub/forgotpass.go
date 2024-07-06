package perahub

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
)

type ForgotPassRequest struct {
	UserName string `json:"username"`
}

type ForgotPassResponseBody struct {
	MobileNo string `json:"mobile"`
}

type ForgotPassResponseData struct {
	Header ResponseHeader         `json:"header"`
	Body   ForgotPassResponseBody `json:"body"`
}

type ForgotPassResponse struct {
	Data ForgotPassResponseData `json:"uspwuapi"`
}

func (s *Svc) ForgotPassword(ctx context.Context, pr ForgotPassRequest) (*ForgotPassResponse, error) {
	req, err := s.newParahubRequest(ctx, "forgot_pwd_init", "forgot_pwd_init", pr)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := s.cl.Post(s.moduleURL("forgot_pwd_init", ""), "aplication/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var rpRes ForgotPassResponse
	if err := json.Unmarshal(body, &rpRes); err != nil {
		return nil, err
	}

	return &rpRes, nil
}

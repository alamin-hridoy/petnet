package perahub

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
)

type SendOTPRequest struct {
	UserCode  string `json:"usercode"`
	SessionID string `json:"session_id"`
}

type SendOTPResponseBody struct{}

type SendOTPResponseData struct {
	Header ResponseHeader  `json:"header"`
	Body   OTPResponseBody `json:"body"`
}

type SendOTPResponse struct {
	Data OTPResponseData `json:"uspwuapi"`
}

func (s *Svc) TransactionOTP(ctx context.Context, tr SendOTPRequest) (*SendOTPResponse, error) {
	req, err := s.newParahubRequest(ctx, "SendSMSUser", "send_sms_user", tr)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := s.cl.Post(s.moduleURL("SendSMSUser", ""), "aplication/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var otpres SendOTPResponse
	if err := json.Unmarshal(body, &otpres); err != nil {
		return nil, err
	}

	return &otpres, nil
}

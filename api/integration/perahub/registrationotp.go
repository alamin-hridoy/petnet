package perahub

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
)

type OTPRequest struct {
	MobileNo string `json:"mobile_no"`
}

type OTPResponseBody struct {
	MobileNo string `json:"mobile"`
}

type OTPResponseData struct {
	Header ResponseHeader  `json:"header"`
	Body   OTPResponseBody `json:"body"`
}

type OTPResponse struct {
	Data OTPResponseData `json:"uspwuapi"`
}

func (s *Svc) RegistrationOTP(ctx context.Context, ores OTPRequest) (*OTPResponse, error) {
	req, err := s.newParahubRequest(ctx, "SendSmsNewUser", "send_sms_new_user", ores)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := s.cl.Post(s.moduleURL("SendSmsNewUser", ""), "aplication/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var otpRes OTPResponse
	if err := json.Unmarshal(body, &otpRes); err != nil {
		return nil, err
	}

	return &otpRes, nil
}

package perahub

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
)

type OTPValidateRequest struct {
	MobileNo string `json:"mobile_no"`
	OTPCode  int    `json:"otp_code"`
}

type OTPVResponseBody struct {
	MobileNo string `json:"mobile"`
}

type OTPVResponseData struct {
	Header ResponseHeader   `json:"header"`
	Body   OTPVResponseBody `json:"body"`
}

type OTPValidateResponse struct {
	Data OTPResponseData `json:"uspwuapi"`
}

func (s *Svc) RegistrationOTPValidate(ctx context.Context, ores OTPValidateRequest) (*OTPValidateResponse, error) {
	req, err := s.newParahubRequest(ctx, "ValidateSMSNewUser", "validate_otp_new_user", ores)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := s.cl.Post(s.moduleURL("ValidateSMSNewUser", ""), "aplication/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var otpRes OTPValidateResponse
	if err := json.Unmarshal(body, &otpRes); err != nil {
		return nil, err
	}

	return &otpRes, nil
}

package perahub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	errCicoValidateSendMoney = errors.New("cico validate send money error")
	errCicoConfirmSendMoney  = errors.New("cico validate send money error")
)
var cicoDynamicURL = []string{"purpose", "partner", "occupation", "employment", "sourcefund", "relationship"}

func cicoAndAction(path string) string {
	return strings.ReplaceAll(path, "/v1/cico/wrapper/", "")
}

func isCicoAction(path string) bool {
	return strings.Contains(path, "/cico/wrapper/")
}

func ciCoDynamicUrlModify(req *http.Request, act string) (rt string) {
	rt = act
	suf := fmt.Sprintf("%s_", req.Method)
	for _, v := range cicoDynamicURL {
		if strings.Contains(act, v) {
			sact := strings.Split(act, "/")
			if len(sact) != 2 {
				rt = suf + rt
				return
			}
			rt = suf + sact[0] + "/{ID}"
			return
		}
	}
	rt = suf + rt
	return
}

func (m *HTTPMock) CicoReq(req *http.Request) (*http.Response, error) {
	m.httpHeaders = req.Header
	var rbb []byte
	var err error
	act := cicoAndAction(req.URL.Path)
	act = ciCoDynamicUrlModify(req, act)
	switch act {
	case "POST_inquiry":
		rbb, err = m.CicoInquire(req)
	case "POST_execute":
		rbb, err = m.CicoExecute(req)
	case "POST_retry":
		rbb, err = m.CicoRetry(req)
	case "POST_otp":
		rbb, err = m.CicoOTPConfirm(req)
	case "POST_validate":
		rbb, err = m.CicoValidate(req)
	}
	sc := 200
	switch err {
	case errCicoValidateSendMoney, errRemitConfirmSendMoney:
		sc = 400
	case conflictErr:
		sc = 409
	case authErr:
		sc = 401
	}
	return &http.Response{
		StatusCode: sc,
		Body:       ioutil.NopCloser(bytes.NewReader(rbb)),
	}, nil
}

func (m *HTTPMock) CicoInquire(req *http.Request) ([]byte, error) {
	if !m.cicoErr {
		rb := &CicoInquireResponse{
			Code:    200,
			Message: "Successful",
			Result: &CicoInquireResult{
				StatusMessage:    "SUCCESSFUL CASHIN",
				PetnetTrackingno: "238a8006885b57765cd8",
				TrxType:          "Cash In",
				ReferenceNumber:  "09654767706",
			},
		}
		return json.Marshal(rb)
	}
	rb := &remitanceError{
		Code:    "404",
		Message: "Transaction can’t be found.",
		Error: map[string]string{
			"status_message":    "Transaction can’t be found",
			"petnet_trackingno": "3115bc3f587d747cf8f0",
			"trx_type":          "Cash In",
			"reference_number":  "09654767706",
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) CicoExecute(req *http.Request) ([]byte, error) {
	if !m.cicoErr {
		rb := &CicoExecuteResponse{
			Code:    200,
			Message: "Successful",
			Result: &CicoExecuteResult{
				PartnerCode:        "DSA",
				Provider:           "GCASH",
				PetnetTrackingno:   "5a269417e107691f3d7c",
				TrxDate:            "2022-05-17",
				TrxType:            "Cash In",
				ProviderTrackingno: "7000001521345",
				ReferenceNumber:    "09654767706",
				PrincipalAmount:    10,
				Charges:            0,
				TotalAmount:        10,
			},
		}
		return json.Marshal(rb)
	}
	rb := &remitanceError{
		Code:    "404",
		Message: "Transaction can’t be found.",
		Error: map[string]string{
			"status_message":    "Transaction can’t be found",
			"petnet_trackingno": "3115bc3f587d747cf8f0",
			"trx_type":          "Cash In",
			"reference_number":  "09654767706",
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) CicoRetry(req *http.Request) ([]byte, error) {
	if !m.cicoErr {
		rb := &CicoRetryResponse{
			Code:    200,
			Message: "SUCCESS TRANSACTION.",
			Result: &CicoRetryResult{
				PartnerCode:        "DSA",
				Provider:           "GCASH",
				PetnetTrackingno:   "5a269417e107691f3d7c",
				TrxDate:            "2022-05-17",
				TrxType:            "Cash In",
				ProviderTrackingno: "09654767706",
				ReferenceNumber:    "09654767706",
				PrincipalAmount:    10,
				Charges:            0,
				TotalAmount:        10,
			},
		}
		return json.Marshal(rb)
	}
	rb := &remitanceError{
		Code:    "404",
		Message: "CICO - EXECUTE ERROR",
		Error: map[string]string{
			"message": "UNRECOGNIZABLE TRACKING NUMBER",
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) CicoOTPConfirm(req *http.Request) ([]byte, error) {
	if !m.cicoErr {
		rb := &CicoOTPConfirmResponse{
			Code:    200,
			Message: "SUCCESS TRANSACTION.",
			Result: &CicoOTPConfirmResult{
				PartnerCode:        "DSA",
				Provider:           "DiskarTech",
				PetnetTrackingno:   "8340b7bdb171cbdf6350",
				TrxDate:            "2022-06-03",
				TrxType:            "Cash Out",
				ProviderTrackingno: "",
				ReferenceNumber:    "220603-000003-1",
				PrincipalAmount:    200,
				Charges:            0,
				TotalAmount:        200,
			},
		}
		return json.Marshal(rb)
	}
	rb := &remitanceError{
		Code:    "500",
		Message: "CICO - OTP ERROR - DISKARTECH",
		Error: map[string]string{
			"message": "Unable to process your transaction at this moment",
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) CicoValidate(req *http.Request) ([]byte, error) {
	if !m.cicoErr {
		rb := &CicoValidateResponse{
			Code:    200,
			Message: "Successful",
			Result: &CicoValidateResult{
				PetnetTrackingno:   "5a269417e107691f3d7c",
				TrxDate:            "2022-05-17",
				TrxType:            "Cash In",
				Provider:           "GCASH",
				ProviderTrackingno: "7000001521345",
				ReferenceNumber:    "09654767706",
				PrincipalAmount:    10,
				Charges:            0,
				TotalAmount:        10,
				Timestamp:          "2022-05-17 09:35:18",
			},
		}
		return json.Marshal(rb)
	}

	rb := &remitanceError{
		Code:    "404",
		Message: "Transaction can’t be found.",
		Error: map[string]string{
			"status_message":    "Transaction can’t be found",
			"petnet_trackingno": "3115bc3f587d747cf8f0",
			"trx_type":          "Cash In",
			"reference_number":  "09654767706",
		},
	}
	return json.Marshal(rb)
}

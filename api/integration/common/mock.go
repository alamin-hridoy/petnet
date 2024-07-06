package common

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	revenueCommissionPath = "/v1/drp/"
	microinsurancePath    = "/v1/insurance/ruralnet/"
)

var revComDynamicURL = []string{}

func IsRevenueCommissionAction(path string) bool {
	return strings.Contains(path, revenueCommissionPath)
}

func RevenueCommissionRequest(req *http.Request) (*http.Response, error) {
	act := dynamicUrlModify(req, strings.ReplaceAll(req.URL.Path, revenueCommissionPath, ""))
	var res []byte
	var err error
	switch act {
	case "POST_dsa":
		res, err = mockDSA(req)
	case "DELETE_dsa/{ID}":
		res, err = mockDSA(req)
	case "POST_remco-sf":
		res, err = mockRemcoServiceFee(req)
	case "GET_remco-sf":
		res, err = listRemcoServiceFee(req)
	case "DELETE_remco-sf/{ID}":
		res, err = mockRemcoServiceFee(req)
	}

	status := 200
	if err != nil {
		status = 500
		res = []byte("{\"code\": 500, \"message\": \"mock error\"}")
	}

	return &http.Response{
		StatusCode: status,
		Body:       ioutil.NopCloser(bytes.NewReader(res)),
	}, nil
}

func dynamicUrlModify(req *http.Request, act string) (rt string) {
	rt = act
	suf := fmt.Sprintf("%s_", req.Method)
	for _, v := range revComDynamicURL {
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

func mockDSA(req *http.Request) ([]byte, error) {
	now := time.Now()
	res := `{"code": 200, "message": "Good", "result":{"id": %d, "dsa_code": "TESTDRP6","dsa_name": "Test DRP", "email_address": "dsaworkshop@yopmail.com", "vatable": "1", "address": "San Mateo, Rizal", "tin": "1234567890", "updated_by": "DRP", "contact_person": "DRP", "city": "SAN MATEO", "province": "RIZAL", "zipcode": "1850", "president": "DRP", "general_manager": "DRP", "updated_at": "2022-09-06T03:35:45.000000Z", "created_at": "2022-09-06T03:35:45.000000Z"}}`
	res = fmt.Sprintf(res, now.Unix())
	return []byte(res), nil
}

func mockRemcoServiceFee(req *http.Request) ([]byte, error) {
	now := time.Now()
	res := fmt.Sprintf(mockRemcoServiceFeeResBody(), now.Unix())
	return []byte(res), nil
}

func listRemcoServiceFee(req *http.Request) ([]byte, error) {
	res := `{
		  "code": 200,
		  "message": "Good",
		  "result": [
			{
			  "id": %d,
			  "remco_id": 16,
			  "min_amount": 0.01,
			  "max_amount": 10000,
			  "service_fee": 0,
			  "commission_amount": 35,
			  "commission_amount_otc": 35,
			  "commission_type": "absolute",
			  "trx_type": "inbound",
			  "updated_by": "SONNY",
			  "created_at": "2022-07-25T09:50:28.000000Z",
			  "updated_at": "2022-07-25T09:50:28.000000Z"
			}
		  ]
		}`

	now := time.Now()
	res = fmt.Sprintf(res, now.Unix())
	return []byte(res), nil
}

func mockRemcoServiceFeeResBody() string {
	return `{
	  "code": 200,
	  "message": "Good",
	  "result": {
		"id": %d,
		"remco_id": "16",
		"min_amount": "250000.01",
		"max_amount": "350000",
		"service_fee": 0,
		"commission_amount": "140",
		"commission_amount_otc": "140",
		"commission_type": "absolute",
		"trx_type": "inbound",
		"updated_by": "DRP",
		"updated_at": "2022-07-25T09:52:25.000000Z",
		"created_at": "2022-07-25T09:52:25.000000Z"
	  }
	}`
}

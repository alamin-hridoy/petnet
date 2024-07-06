package perahub

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/common"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/storage/postgres"
	"brank.as/petnet/svcutil/random"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HTTPMock struct {
	st           *postgres.Storage
	partnerErr   bool
	drpErr       bool
	nonexErr     bool
	crPtnrErr    bool
	dsbPtnrErr   bool
	remitanceErr bool
	cfDsbPtnrErr bool
	cfCrPtnrErr  bool
	conflictErr  bool
	crUsrErr     bool
	authErr      bool
	saveReq      bool
	cicoErr      bool
	miErr        bool
	reqBody      []byte
	httpHeaders  http.Header
	reqOrder     []string
}

func NewHTTPMock(st *postgres.Storage) *HTTPMock {
	return &HTTPMock{
		st: st,
	}
}

type MockConfig struct {
	CrPtnrErr    bool
	DsbPtnrErr   bool
	CfCrPtnrErr  bool
	CfDsbPtnrErr bool
	MIPtnrErr    bool
	SaveReq      bool
}

func NewTestHTTPMock(st *postgres.Storage, c MockConfig) *HTTPMock {
	return &HTTPMock{
		st:           st,
		crPtnrErr:    c.CrPtnrErr,
		dsbPtnrErr:   c.DsbPtnrErr,
		cfCrPtnrErr:  c.CfCrPtnrErr,
		cfDsbPtnrErr: c.CfDsbPtnrErr,
		miErr:        c.MIPtnrErr,
		saveReq:      c.SaveReq,
	}
}

var BillerAPiUrls = []string{"biller-category", "biller-by-category", "validate-account", "transact", "retry"}

type NonexReq struct {
	RefNo string `json:"reference_number"`
}

func (m *HTTPMock) GetMockRequest() []byte {
	return m.reqBody
}

func (m *HTTPMock) SetConflictError() {
	m.conflictErr = true
}

func (m *HTTPMock) SetAuthError() {
	m.authErr = true
}

func (m *HTTPMock) UnsetAuthError() {
	m.authErr = false
}

func (m *HTTPMock) ResetReqOrder() {
	m.reqOrder = []string{}
}

func (m *HTTPMock) Do(req *http.Request) (*http.Response, error) {
	if common.IsMicroInsuranceRequest(req.URL.Path) {
		return common.MicroInsuranceRequest(req, m.miErr)
	}

	r := &PerahubRequest{}
	if req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(body, r); err != nil {
			return nil, err
		}
		if m.saveReq {
			m.reqBody = body
		}
	}

	if isPerahubGetRemocoIDAction(req.URL.Path) {
		return m.GetRemocoIdReq(req)
	}

	if isRemitanceAction(req.URL.Path) {
		return m.RemitanceReq(req)
	}

	if isCicoAction(req.URL.Path) {
		return m.CicoReq(req)
	}

	isBillsPayment := strings.Contains(req.URL.Path, "billspay")
	oldBills := strings.Contains(req.URL.Path, "wrapper/api")
	if isBillsPayment && oldBills {
		bu := billerAndAction(req.URL.Path)
		if IsThere(bu, BillerAPiUrls) {
			return m.BillerReq(req)
		}
	} else if isBillsPayment {
		return m.BillsPayReq(req)
	}

	if common.IsRevenueCommissionAction(req.URL.Path) {
		return common.RevenueCommissionRequest(req)
	}

	ptnr, act, errp := partnerAndAction(req.URL.Path)
	if errp != nil {
		return nil, errp
	}
	if ptnr == "cebuana" {
		switch act {
		case "get-service-fee":
			var newReq CebuanaSFInquiryRequest
			reqData := strings.Split(req.URL.RawQuery, "&")
			if len(reqData) > 0 {
				for _, dt := range reqData {
					sdt := strings.Split(dt, "=")
					if len(sdt) > 0 {
						switch sdt[0] {
						case "principal_amount":
							newReq.PrincipalAmount = json.Number(sdt[1])
						case "currency_id":
							newReq.CurrencyID = json.Number(sdt[1])
						case "agent_code":
							newReq.AgentCode = sdt[1]
						}
					}
				}
			}
			reqBody, err := json.Marshal(newReq)
			if err != nil {
				return nil, err
			}
			m.reqBody = reqBody
		case "find-client":
			var newReq CebFindClientRequest
			reqData := strings.Split(req.URL.RawQuery, "&")
			if len(reqData) > 0 {
				for _, dt := range reqData {
					sdt := strings.Split(dt, "=")
					if len(sdt) > 0 {
						switch sdt[0] {
						case "first_name":
							newReq.FirstName = sdt[1]
						case "last_name":
							newReq.LastName = sdt[1]
						case "birth_date":
							newReq.BirthDate = sdt[1]
						case "client_number":
							newReq.ClientNumber = sdt[1]
						}
					}
				}
			}
			reqBody, err := json.Marshal(newReq)
			if err != nil {
				return nil, err
			}
			m.reqBody = reqBody
		case "send-currency-collection":
			var newReq CEBCurrencyReq
			reqData := strings.Split(req.URL.RawQuery, "&")
			if len(reqData) > 0 {
				for _, dt := range reqData {
					sdt := strings.Split(dt, "=")
					if len(sdt) > 0 {
						switch sdt[0] {
						case "agent_code":
							newReq.AgentCode = sdt[1]
						}
					}
				}
			}
			reqBody, err := json.Marshal(newReq)
			if err != nil {
				return nil, err
			}
			m.reqBody = reqBody
		case "beneficiary-by-sender":
			var newReq CebFindBFReq
			reqData := strings.Split(req.URL.RawQuery, "&")
			if len(reqData) > 0 {
				for _, dt := range reqData {
					sdt := strings.Split(dt, "=")
					if len(sdt) > 0 {
						switch sdt[0] {
						case "sender_client_id":
							newReq.SenderClientId = sdt[1]
						}
					}
				}
			}
			reqBody, err := json.Marshal(newReq)
			if err != nil {
				return nil, err
			}
			m.reqBody = reqBody
		}
	}
	if ptnr == "perahub-remit" {
		switch act {
		case "address":
			var newReq GetProvincesCityListRequest
			reqData := strings.Split(req.URL.RawQuery, "&")
			if len(reqData) > 0 {
				for _, dt := range reqData {
					sdt := strings.Split(dt, "=")
					if len(sdt) > 0 {
						switch sdt[0] {
						case "id":
							id, err := strconv.Atoi(sdt[1])
							if err != nil {
								return nil, err
							}
							newReq.ID = id
						case "partner_code":
							newReq.PartnerCode = sdt[1]
						}
					}
				}
			}
			reqBody, err := json.Marshal(newReq)
			if err != nil {
				return nil, err
			}
			m.reqBody = reqBody
		}
	}
	var res *http.Response
	var err error
	if isNonexRequest(r) {
		res, err = m.nonexResp(req)
		if err != nil {
			return nil, err
		}
	} else {
		res, err = m.wuResp(req, r)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (m *HTTPMock) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	bd, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	r := &PerahubRequest{}
	if err := json.Unmarshal(bd, r); err != nil {
		return nil, err
	}

	var b []byte
	switch r.WU.Body.Request {
	case "todo":
	}
	res := ioutil.NopCloser(bytes.NewReader(b))
	return &http.Response{
		StatusCode: 200,
		Body:       res,
	}, nil
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

func partnerAndAction(path string) (string, string, error) {
	str := strings.ReplaceAll(path, "/v1/remit/nonex/", "")
	a := strings.SplitN(str, "/", 2)
	if len(a) < 2 {
		return "", "", fmt.Errorf("should have partner and action, got: %v", a)
	}
	return a[0], a[1], nil
}

func billerAndAction(path string) (str string) {
	strS := strings.ReplaceAll(path, "/v1/billspay/wrapper/api/", "")
	strM := strings.Split(strS, "/")
	str = strM[0]
	return
}

func billsPayAndAction(path string) (str string) {
	strS := strings.ReplaceAll(path, "/v1/billspay/", "")
	strM := strings.Split(strS, "/")
	if len(strM) == 3 && strM[1] == "biller-by-category" {
		strM[2] = "{ID}"
	}
	str = strings.Join(strM, "_")
	return str
}

func (m *HTTPMock) nonexResp(req *http.Request) (*http.Response, error) {
	ptnr, act, err := partnerAndAction(req.URL.Path)
	if err != nil {
		return nil, err
	}

	m.httpHeaders = req.Header

	switch {
	case m.nonexErr:
		rbb := []byte("<html>server-error<html>")
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader(rbb)),
		}, nil
	case m.drpErr:
		rbb := []byte(`{
				"non": {
					"standard": "error"
				}
		}`)
		return &http.Response{
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewReader(rbb)),
		}, nil
	}

	var rbb []byte
	switch ptnr {
	case "iremit":
		rbb, err = m.irResp(act)
	case "transfast":
		rbb, err = m.tfResp(act)
	case "ria":
		rbb, err = m.riaResp(act)
	case "remitly":
		rbb, err = m.rmResp(act)
	case "metrobank":
		rbb, err = m.mbResp(act)
	case "bpi":
		rbb, err = m.bpResp(act)
	case "ussc":
		rbb, err = m.usscResp(act)
	case "transferwise":
		rbb, err = m.wiseResp(act)
	case "instacash":
		rbb, err = m.icResp(act)
	case "japanremit":
		rbb, err = m.jprResp(act)
	case "uniteller":
		rbb, err = m.untResp(act)
	case "cebuana":
		rbb, err = m.cebResp(act)
	case "cebuana-international":
		rbb, err = m.cebIntResp(act)
	case "ayannah":
		rbb, err = m.ayannahResp(act)
	case "intelexpress":
		rbb, err = m.ieResp(act)
	case "perahub-remit":
		rbb, err = m.peraHubRemitResp(req)
	default:
		return nil, fmt.Errorf("no response for partner: %v", ptnr)
	}

	sc := 200
	switch err {
	case crPtnrErr, dsbPtnrErr, cfDsbPtnrErr, cfCrPtnrErr:
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

func (m *HTTPMock) BillerReq(req *http.Request) (*http.Response, error) {
	m.httpHeaders = req.Header
	var rbb []byte
	var err error
	act := billerAndAction(req.URL.Path)
	switch act {
	}
	sc := 200
	switch err {
	case crPtnrErr, dsbPtnrErr, cfDsbPtnrErr, cfCrPtnrErr:
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

func (m *HTTPMock) BillsPayErr() []byte {
	rb := &nonexError{
		Code: "1",
		Msg:  "some error",
		Error: ErrorType{
			Type: string(BillerError),
			Msg:  "internal error",
		},
	}
	rbb, err := json.Marshal(rb)
	if err != nil {
		return []byte{}
	}
	return rbb
}

func (m *HTTPMock) BillsPayReq(req *http.Request) (*http.Response, error) {
	m.httpHeaders = req.Header
	var rbb []byte
	var err error
	act := billsPayAndAction(req.URL.Path)

	switch act {
	case "multipay_biller-inquire":
		rbb = []byte(`{
			"status": 200,
			"reason": "OK",
			"data": {
			  "account_number": "MP1MN7JMJV",
			  "amount": 51,
			  "biller": "MSYS_TEST_BILLER"
			}
		  }`)
	case "multipay_biller-process":
		rbb = []byte(`{
				"status": 200,
				"reason": "You have successfully processed this transaction.",
				"data": {
				  "refno": "MP1MN7JMJV",
				  "txnid": "TEST-62F0802B98B6E",
				  "biller": "MSYS_TEST_BILLER",
				  "meta": []
				}
			  }`)
	case "multipay_biller-by-category_{ID}":
		rbb = []byte(`{
					"code": 200,
					"message": "Good",
					"result": [
					  {
						"partner_id": 3,
						"BillerTag": "MULTIPAY-AppendPay",
						"Description": "AppendPay",
						"Category": 1,
						"FieldList": [
						  {
							"id": "amount",
							"type": "numeric",
							"label": "Amount",
							"order": 1,
							"rules": [
							  {
								"code": 1,
								"type": "required",
								"value": "",
								"format": "",
								"message": "Please provide the amount.",
								"options": ""
							  }
							],
							"description": "Amount to be paid",
							"placeholder": "Insert Amount"
						  }
						],
						"ServiceCharge": 0
					  }
					]
				  }`)
	case "multipay_biller-category":
		rbb = []byte(`{
					"code": 200,
					"message": "Good",
					"result": [
						{
						"id": 6,
						"bill_id": 1,
						"category_name": "Airlines",
						"created_at": "2022-04-08T04:26:01.000000Z",
						"updated_at": "2022-04-08T04:26:01.000000Z"
						}
					]
					}`)
	case "multipay_biller-list":
		rbb = []byte(`{
					"code": 200,
					"message": "Good",
					"result": [
						{
						"partner_id": 3,
						"BillerTag": "MULTIPAY-NBI",
						"Description": "NBI",
						"Category": 4,
						"FieldList": [
							{
							"id": "amount",
							"type": "numeric",
							"label": "Amount",
							"order": 1,
							"rules": [
								{
								"code": 1,
								"type": "required",
								"value": "",
								"format": "",
								"message": "Please provide the amount.",
								"options": ""
								}
							],
							"description": "Amount to be paid",
							"placeholder": "Insert Amount"
							}
						],
						"ServiceCharge": 0
						}
					]
					}`)
	case "multipay_search-transaction":
		rbb = []byte(`{
					"code": 200,
					"message": "Good",
					"result": {
						"txnid": "TEST-62F0802B98B6E",
						"refno": "MP9UYHSEHSIKRTY4RHC",
						"amount": "51.00",
						"fee": "0.00",
						"status": "V",
						"payment_channel": "PGI",
						"is_transaction_expired": true,
						"created_at": "2022-08-09 16:48:58",
						"expires_at": "2022-08-12 16:48:58"
					}
					}`)
	case "multipay_generate-transaction":
		rbb = []byte(`{
					"data": {
						"url": "https://pgi-staging.multipay.ph/MP0XQPUWDTVSRUONACH"
					}
					}`)
	case "multipay_void-transaction":
		rbb = []byte(`{
					"data": {
						"txnid": "TEST-62F0802B98B6E",
						"refno": "MP0XQPUWDTVSRUONACH",
						"amount": "51.00",
						"fee": "0.00",
						"status": "V",
						"payment_channel": "PGI",
						"is_transaction_expired": true,
						"created_at": "2022-09-05 07:41:36",
						"expires_at": "2022-09-08 07:41:36"
					}
					}`)
	case "ecpay_biller-by-category_{ID}":
		rbb = []byte(`{
			"code": 200,
			"message": "Success",
			"result": [
			  {
				"BillerTag": "MANILAWATER",
				"Description": "MANILA WATER COMPANY",
				"FirstField": "8 Digit Contract Account Number",
				"FirstFieldFormat": "Numeric",
				"FirstFieldWidth": "8",
				"SecondField": "Account Name",
				"SecondFieldFormat": "Alphanumeric",
				"SecondFieldWidth": "30",
				"ServiceCharge": 2
			  }
			]
		  }`)
	case "ecpay_biller-category":
		rbb = []byte(`[
			{
			  "id": 6,
			  "bill_id": 1,
			  "category_name": "Airlines",
			  "created_at": "2022-04-08T04:26:01.000000Z",
			  "updated_at": "2022-04-08T04:26:01.000000Z"
			}
		  ]`)
	case "ecpay_biller-list":
		rbb = []byte(`{
			"code": 200,
			"message": "Success",
			"result": [
			  {
				"BillerTag": "MANILAWATER",
				"Description": "MANILA WATER COMPANY",
				"FirstField": "8 Digit Contract Account Number",
				"FirstFieldFormat": "Numeric",
				"FirstFieldWidth": "8",
				"SecondField": "Account Name",
				"SecondFieldFormat": "Alphanumeric",
				"SecondFieldWidth": "30",
				"ServiceCharge": 2
			  }
			],
			"remco_id": 1
		  }`)
	case "ecpay_check-balance":
		rbb = []byte(`{
			"code": 200,
			"message": "Success",
			"result": {
			  "RemBal": "10000.00"
			},
			"remco_id": 1
		  }`)
	case "ecpay_validate-account":
		rbb = []byte(`{
			"code": 200,
			"message": "Success",
			"result": "DAVAOLIGHT",
			"remco_id": 1
		  }`)
	case "ecpay_transact":
		rbb = []byte(`{
			"code": "0",
			"message": "Success",
			"result": {
			  "Status": "0",
			  "Message": "SUCCESS! REF #F6L2S84JN00M",
			  "ServiceCharge": 10,
			  "timestamp": "2022-05-19 06:46:05",
			  "referenceNumber": "F6L2S84JN00M"
			},
			"remco_id": 1
		  }`)
	case "ecpay_retry":
		rbb = []byte(`{
			"code": "0",
			"message": "Success",
			"result": {
			  "Status": "0",
			  "Message": "SUCCESS! REF #72482FD0A467",
			  "ServiceCharge": 10,
			  "timestamp": "2021-03-28 08:58:28",
			  "referenceNumber": "72482FD0A467"
			},
			"remco_id": 1
		  }`)
	case "bayad_bayad-center_biller-info":
		rbb = []byte(`{
				"code": 200,
				"message": "Success",
				"result": {
					"code": "MWCOM",
					"isCde": 0,
					"isAsync": 0,
					"name": "MANILA WATER",
					"description": "Manila Water Company",
					"logo": "https://stg-bc-api-images.s3-ap-southeast-1.amazonaws.com/biller-logos/250/MWCOM.png",
					"category": "Water",
					"type": "Batch",
					"isMultipleBills": 0,
					"parameters": {
						"verify": [
							{
								"referenceNumber": {
									"description": "Account Number",
									"rules": {
										"digits:8": {
											"message": "The account number must be 8 digits.",
											"code": 5
										},
										"required": {
											"message": "Please provide the account number.",
											"code": 4
										}
									},
									"label": "Account Number"
								}
							}
						],
						"transact": [
							{
								"clientReference": {
									"description": "Client unique transaction reference number",
									"rules": {
										"alpha_dash": {
											"message": "Please make sure that the client reference number is in alpha dash format.",
											"code": 9
										},
										"required": {
											"message": "Please provide the client reference number.",
											"code": 4
										},
										"unique_crn": {
											"message": "This client reference number already exists.",
											"code": 11
										}
									},
									"label": "Client Reference Number"
								}
							}
						]
					}
				},
				"remco_id": 2
			}`)
	case "bayad_bayad-center_billers":
		rbb = []byte(`{
				"code": 200,
				"message": "Success",
				"result": [
					{
						"name": "MERALCO",
						"code": "MECOR",
						"description": "Meralco Real-Time Posting",
						"category": "Electricity",
						"type": "RTP",
						"logo": "https://stg-bc-api-images.s3-ap-southeast-1.amazonaws.com/biller-logos/250/MECOR.png",
						"isMultipleBills": 1,
						"isCde": 0,
						"isAsync": 1
					}
				]
			}`)
	case "bayad_bayad-center_token":
		rbb = []byte(`{
				"code": 200,
				"message": "Success",
				"result": {
					"access_token": "eyJraWQiOiJraWoydFFERTZiSWxnOFE3enZMSmFZaE5jNXdlWHRzaVM0OW1vYVR4YWs0PSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiI0M2JoZThlcjk1YzM4ajg3MDdtbWpsYWxkIiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJtZWNvbS1hdXRoXC9hbGwiLCJhdXRoX3RpbWUiOjE2NjY3NTEwMjgsImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC5hcC1zb3V0aGVhc3QtMS5hbWF6b25hd3MuY29tXC9hcC1zb3V0aGVhc3QtMV9aZkJqVWVTeTMiLCJleHAiOjE2NjY3NTQ2MjgsImlhdCI6MTY2Njc1MTAyOCwidmVyc2lvbiI6MiwianRpIjoiMGRjNjZjNDUtODY2My00NWY5LWIyNWItNDVkMDg0YzZiMDY2IiwiY2xpZW50X2lkIjoiNDNiaGU4ZXI5NWMzOGo4NzA3bW1qbGFsZCJ9.dJpsEgqPQRSLrjNlz9QS74-gX3m-DNqiDfyBSl_NtaEupCguzFP_G0pRn1I2oxdfNsAUQog0a-NA2KYSQqA_CHmo81JoSPVaXmc7EhlPWl2ANkn7brVVSroRn3Of_cktS_gsxMWDe2t7Wxb8cSt5sFKaT2USA2fMqB1r4RoCmz3s6k9yvs_Niukmmzw_o6_4brte0-gxW6F1Jfx8dF2M27RXvMBVZ3fJQYHz_Njq-pfhuY-xlbUZvgyCywGOihZ7V-jvxp2p9vh3mIgYVtrEkY08vJeHdXpv2VxbZYd3snqdI4CO68jNONpqOqzZHSsRyE3hV6JNGtJThWZ0ezRqFw",
					"expires_in": 3600,
					"token_type": "Bearer"
				},
				"remco_id": 2
			}`)
	case "bayad_bayad-center_wallets":
		rbb = []byte(`{
				"code": 200,
				"message": "Success",
				"result": {
					"balance": "0.00"
				},
				"remco_id": 2
			}`)
	case "bayad_bayad-center_fees":
		rbb = []byte(`{
				"code": 200,
				"message": "Success",
				"result": {
					"otherCharges": "5.00"
				},
				"remco_id": 2
			}`)
	case "bayad_bayad-center_retry":
		rbb = []byte(`{
				"code": 200,
				"message": "Success",
				"result": {
					"transactionId": "21187PP01024FF38I",
					"referenceNumber": "511111",
					"clientReference": "6217d18d-5cff-4f1e-affd-6503883dflk0",
					"billerReference": "PP012118724FF38I",
					"paymentMethod": "CASH",
					"amount": "100.00",
					"otherCharges": "0.00",
					"status": "PENDING",
					"message": "The payment was successfully created.",
					"details": [],
					"createdAt": "2021-07-06 16:57:29",
					"timestamp": "2021-07-06 08:57:25"
				},
				"remco_id": 2
			}`)
	case "bayad_bayad-center_transact-inquire":
		rbb = []byte(`{
				"code": 200,
				"message": "Success",
				"result": {
					"transactionId": "21152PP0120000006",
					"referenceNumber": "0136402637",
					"clientReference": "a1mef46c-d994-487e-a086-oaesdc78a42f",
					"billerReference": "21152PP0120000006",
					"paymentMethod": "CASH",
					"amount": "3974.83",
					"otherCharges": "0.00",
					"status": "PENDING",
					"message": {
						"header": "Payment Receipt",
						"message": "Sweet! We have received your MERALCO bill payment and are currently processing it. Thank you. Have a great day ahead!",
						"footer": "Please note that payments made after 7PM will be posted 7AM the next day."
					},
					"details": [],
					"createdAt": "2021-06-01 16:41:41"
				},
				"remco_id": 2
			}`)
	case "bayad_bayad-center_transact":
		rbb = []byte(`{
				"code": 200,
				"message": "Success",
				"result": {
					"transactionId": "21187PP01024FF38I",
					"referenceNumber": "511111",
					"clientReference": "6217d18d-5cff-4f1e-affd-6503883dflk0",
					"billerReference": "PP012118724FF38I",
					"paymentMethod": "CASH",
					"amount": "100.00",
					"otherCharges": "0.00",
					"status": "PENDING",
					"message": "The payment was successfully created.",
					"details": [],
					"createdAt": "2021-07-06 16:57:29",
					"timestamp": "2021-07-06 08:57:25"
				},
				"remco_id": 2
			}`)
	case "bayad_bayad-center_validate":
		rbb = []byte(`{
				"code": 200,
				"message": "Success",
				"result": {
					"valid": true,
					"code": 0,
					"account": "200820248",
					"details": [],
					"validationNumber": "48f2b647-0ff8-4ace-9928-346258f08df5"
				},
				"remco_id": 2
			}`)
	}
	sc := 200
	switch err {
	case crPtnrErr, dsbPtnrErr, cfDsbPtnrErr, cfCrPtnrErr:
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

var (
	crPtnrErr  = errors.New("create remit error")
	dsbPtnrErr = errors.New("disburse error")

	cfDsbPtnrErr = errors.New("confirm disburse error")
	cfCrPtnrErr  = errors.New("confirm create remit error")
	conflictErr  = errors.New("conflict error")
	authErr      = errors.New("unauthorized error")
)

func (m *HTTPMock) irResp(act string) ([]byte, error) {
	switch act {
	case "inquire":
		if !m.dsbPtnrErr {
			rb := &IRInquireResponseBody{
				Code: "0",
				Msg:  "Available for Pick-up",
				Result: IRResult{
					Status:        "0",
					Desc:          "Available for Pick-up",
					ControlNo:     "CTRL1",
					RefNo:         "REF1",
					PnplAmt:       "1000.00",
					SenderName:    "John, Michael Doe",
					RcvName:       "Jane, Emily Doe",
					Address:       "PLA",
					CurrencyCode:  "PHP",
					ContactNumber: "09190000000",
					RcvLastName:   "Doe",
					RcvFirstName:  "Jane",
				},
				RemcoID: "1",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: IRErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, dsbPtnrErr
		}
	case "payout":
		if !m.cfDsbPtnrErr {
			rb := &IRPayoutResponseBody{
				Code:    "200",
				Msg:     "Successful.",
				RemcoID: "1",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: IRErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfDsbPtnrErr
		}
	}
	return nil, fmt.Errorf("no response for action: %v", act)
}

func (m *HTTPMock) tfResp(act string) ([]byte, error) {
	switch act {
	case "inquire":
		if !m.dsbPtnrErr {
			rb := &TFInquireResponseBody{
				Code: "200",
				Msg:  "Success!",
				Result: TFResult{
					Status:        "T",
					Desc:          "TRANSMIT",
					ControlNo:     "CTRL1",
					RefNo:         "1",
					PnplAmt:       "1000.00",
					SenderName:    "John, Michael Doe",
					RcvName:       "Jane, Emily Doe",
					Address:       "PLA",
					CurrencyCode:  "PHP",
					ContactNumber: "09190000000",
					RcvLastName:   "Doe",
					RcvFirstName:  "Jane",
					OrgnCtry:      "UNITED ARAB EMIRATES",
					DestCtry:      "PHILIPPINES",
					TxnDate:       "2021-07-15T13:26:01.493-04:00",
					IsDomestic:    "0",
					IDType:        "0",
					RcvCtryCode:   "PH",
					RcvStateID:    "PH023",
					RcvStateName:  "METRO MANILA",
					RcvCityID:     "1",
					RcvCityName:   "METRO MANILA",
					RcvIDType:     "1",
					RcvIsIndiv:    "True",
					PrpsOfRmtID:   "1",
				},
				RemcoID: "1",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: TFErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, dsbPtnrErr
		}
	case "payout":
		if !m.cfDsbPtnrErr {
			rb := &IRPayoutResponseBody{
				Code:    "200",
				Msg:     "Successful.",
				RemcoID: "7",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: TFErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfDsbPtnrErr
		}
	case "receiver-ids":
		rb := &TFIDsRespBody{
			TFBaseResponseBody: TFBaseResponseBody{
				Code:    "200",
				Msg:     "Success!",
				RemcoID: "7",
			},
			Result: TFIDsResult{
				IDs: []TFID{
					{
						ID:             "1",
						Name:           "FAMILY MAINTENANCE",
						CountryIsoCode: "PH",
					},
					{
						ID:             "2",
						Name:           "EDUCATION",
						CountryIsoCode: "PH",
					},
				},
			},
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "relationships":
		rb := &TFRelationsRespBody{
			TFBaseResponseBody: TFBaseResponseBody{
				Code:    "200",
				Msg:     "Success!",
				RemcoID: "7",
			},
			Result: TFRelationsResult{
				Relations: []TFIDName{
					{
						ID:   "8",
						Name: "FRIEND",
					},
					{
						ID:   "16",
						Name: "Self",
					},
				},
			},
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "beneficiary-occupations":
		rb := &TFOccupsRespBody{
			TFBaseResponseBody: TFBaseResponseBody{
				Code:    "200",
				Msg:     "Success!",
				RemcoID: "7",
			},
			Result: TFOccupsResult{
				Occups: []TFIDName{
					{
						ID:   "1",
						Name: "HOUSEWIFE",
					},
					{
						ID:   "2",
						Name: "STUDENT",
					},
				},
			},
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "remittance-purposes":
		rb := &TFPrpsRespBody{
			TFBaseResponseBody: TFBaseResponseBody{
				Code:    "200",
				Msg:     "Success!",
				RemcoID: "7",
			},
			Result: TFPrpsResult{
				Prps: []TFPrp{
					{
						ID:             "1",
						Name:           "FAMILY MAINTENANCE",
						CountryIsoCode: "PH",
					},
					{
						ID:             "2",
						Name:           "EDUCATION",
						CountryIsoCode: "PH",
					},
				},
			},
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	}
	return nil, fmt.Errorf("no response for action: %v", act)
}

func (m *HTTPMock) riaResp(act string) ([]byte, error) {
	switch act {
	case "inquire":
		switch {
		default:
			rb := &RiaInquireResponseBody{
				Code: "200",
				Msg:  "Order is available for payout.",
				Result: RiaResult{
					ControlNumber:      "CTRL1",
					ClientReferenceNo:  "1",
					OriginatingCountry: "TH",
					DestinationCountry: "PH",
					SenderName:         "John, Michael Doe",
					ReceiverName:       "Jane, Emily Doe",
					PrincipalAmount:    "1000.00",
					Currency:           "PHP",
					IsDomestic:         "1",
					OrderNo:            "TH1950882455",
				},
				RemcoID: "12",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		case m.dsbPtnrErr:
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: RMErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, dsbPtnrErr
		}
	case "payout":
		if !m.cfDsbPtnrErr {
			rb := &RiaPayoutResponseBody{
				Code:    "200",
				Msg:     "Successful.",
				RemcoID: "12",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: RMErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfDsbPtnrErr
		}
	}
	return nil, fmt.Errorf("no response for action: %v", act)
}

func (m *HTTPMock) rmResp(act string) ([]byte, error) {
	switch act {
	case "inquire":
		if !m.dsbPtnrErr {
			rb := &RMInquireResponseBody{
				Code: "200",
				Msg:  "PAYABLE",
				Result: RMInqResult{
					ControlNo:     "CTRL1",
					PnplAmt:       "1000.00",
					CurrencyCode:  "PHP",
					RcvName:       "Jane, Emily Doe",
					Address:       "7118 Street",
					City:          "Manila",
					Country:       "PHL",
					SenderName:    "John, Michael Doe",
					ContactNumber: "9162427505",
				},
				RemcoID: "21",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: RMErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, dsbPtnrErr
		}
	case "payout":
		if !m.cfDsbPtnrErr {
			rb := &RMPayoutResponseBody{
				Code: "200",
				Msg:  "Success",
				Result: RMPayResult{
					RefNo:      "REF1",
					Created:    "2021-10-21",
					State:      "PAID",
					Type:       "CASH_PICKUP",
					PayerCodes: "",
				},
				RemcoID: "21",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: RMErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfDsbPtnrErr
		}
	case "ids":
		rb := &RMIDsRespBody{
			Code: "200",
			Msg:  "Success",
			Result: []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			}{
				{
					Name:  "AFP ID",
					Value: "GOVERNMENT_ISSUED_ID",
				},
				{
					Name:  "Driver License",
					Value: "DRIVERS_LICENSE",
				},
			},
			RemcoID: "21",
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	}
	return nil, fmt.Errorf("no response for action: %v", act)
}

func (m *HTTPMock) mbResp(act string) ([]byte, error) {
	switch act {
	case "inquire":
		if !m.dsbPtnrErr {
			rb := &MBInquireResponseBody{
				Code: "0",
				Msg:  "Available for pick-up",
				Result: MBInqResult{
					RefNo:           "REF1",
					ControlNo:       "CTRL1",
					StatusText:      "0",
					PrincipalAmount: "1000.00",
					RcvName:         "Jane, Emily Doe",
					Address:         "7118 Street",
					ContactNumber:   "9162427505",
					Currency:        "PHP",
				},
				RemcoID: "8",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: MBErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, dsbPtnrErr
		}
	case "payout":
		if !m.cfDsbPtnrErr {
			rb := &MBPayoutResponseBody{
				Code: "200",
				Msg:  "Successful.",
				Result: MBPayResult{
					RefNo:           "REF1",
					ClientRefNo:     "1",
					ControlNo:       "CTRL1",
					StatusText:      "0",
					PrincipalAmount: "1000.00",
					RcvName:         "rcv-nam",
					Address:         "addr",
					ReceiptNo:       "",
					ContactNumber:   "201-20211126-000288",
				},
				RemcoID: "8",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: MBErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfDsbPtnrErr
		}
	}
	return nil, fmt.Errorf("no response for action: %v", act)
}

func (m *HTTPMock) wuResp(httpreq *http.Request, phreq *PerahubRequest) (*http.Response, error) {
	var err error
	var rbb []byte
	modReq := phreq.WU.Body.Request
	switch modReq {
	// Send money flow
	// 1. feeinquiry
	// 2. SMvalidate
	// 3. SendMoneyStore
	// CreateRemit = feeinquiry, SMvalidate (caches remit)
	// ConfirmRemit = SMstore
	case "feeinquiry":
		d := &FIRequest{}
		if err := json.Unmarshal(phreq.WU.Body.Param, d); err != nil {
			return nil, err
		}

		rb := &FIResponseBody{
			OrigPrincipal: d.PrincipalAmount,
			DestPrincipal: d.PrincipalAmount,
			ExchangeRate:  json.Number("1.0000000"),
			GrossTotal:    d.PrincipalAmount,
			PayAmount:     d.PrincipalAmount,
			Charges:       json.Number("300"),

			BaseMsgCharge: "5000",
			BaseMsgLimit:  "10",
			IncMsgCharge:  "500",
			IncMsgLimit:   "10",
		}
		rbb, err = json.Marshal(rb)
		if err != nil {
			return nil, err
		}
	case "SMvalidate":
		d := &SMVRequest{}
		if err := json.Unmarshal(phreq.WU.Body.Param, d); err != nil {
			return nil, err
		}

		rb := &SMVResponseBody{
			ServiceCode: ServiceCode{
				AddlSvcChg: "100",
			},
			Compliance: Compliance{
				ComplianceBuf: "100",
			},
			PaymentDetails: PaymentDetails{
				OrigCity:  "METRO MANILA",
				OrigState: "MA",
			},
			Fin: Financials{
				OrigPnplAmt: "100000",
				DestPcplAmt: "100000",
				GrossTotal:  "100100",
				Charges:     "100",
			},
			NewDetails: NewDetails{
				MTCN:       random.InvitationCode(10),
				NewMTCN:    random.InvitationCode(10),
				FilingDate: "09-27",
				FilingTime: "0810P EDT",
			},
			PreferredCustomer: PreferredCustomer{
				MyWUNumber: random.InvitationCode(10),
			},
		}
		rbb, err = json.Marshal(rb)
		if err != nil {
			return nil, err
		}
	case "SMstore":
		d := &SMStoreRequest{}
		if err := json.Unmarshal(phreq.WU.Body.Param, d); err != nil {
			return nil, err
		}

		rb := &struct {
			ConfirmedDetails *ConfirmedDetails `json:"confirmed_details"`
		}{
			ConfirmedDetails: &ConfirmedDetails{
				AdvisoryText:     "09/27/21 20:31",
				MTCN:             d.MTCN,
				NewMTCN:          d.NewMTCN,
				FilingDate:       "09-27",
				FilingTime:       "0831P EDT",
				PromoTextMessage: []string{"Mahalaga sa amin ang iyong opinyon!"},
			},
		}
		rbb, err = json.Marshal(rb)
		if err != nil {
			return nil, err
		}
	// Recieve money flow
	// 1. search
	// 2. pay
	// LookupRemit = search (without caching, only for lookup)
	// Disburse = search (caches the remit)
	// ConfirmRemit = pay
	//
	case "search":
		d := &RMSearchRequest{}
		if err := json.Unmarshal(phreq.WU.Body.Param, d); err != nil {
			return nil, err
		}

		rs, err := m.st.ListRemitHistory(httpreq.Context(),
			storage.LRHFilter{
				TxnStep:   string(storage.ConfirmStep),
				TxnStatus: string(storage.SuccessStatus),
				ControlNo: []string{d.MTCN},
			},
		)
		if err != nil {
			return nil, err
		}
		rm := rs[0]

		rb := &RMSearchResponseBody{
			GrossPayout: "1000",
			Txn: RMSPaymentTransaction{
				Sender: Contact{
					Name: Name{
						NameType:  "D",
						FirstName: rm.Remittance.Remitter.FirstName,
						LastName:  rm.Remittance.Remitter.LastName,
					},
					Address: RMSAddress{
						Street:     rm.Remittance.Remitter.Address1,
						City:       rm.Remittance.Remitter.City,
						State:      rm.Remittance.Remitter.State,
						PostalCode: rm.Remittance.Remitter.PostalCode,
						CountryCode: RMSCountryCode{
							IsoCode: RMSIsoCode{
								Country:  rm.Remittance.Remitter.Country,
								Currency: rm.Remittance.SourceAmt.CurrencyCode(),
							},
						},
					},
					Mobile: RMSMobilePhone{
						RawPhone: json.RawMessage(`{"phone_number":{"country_code":"PHP","national_number":"1324423"}}`),
						Phone: RMSPhoneNumber{
							CountryCode: "PHP",
							Number:      "12345",
						},
					},
				},
				Receiver: Contact{
					Name: Name{
						NameType:  "D",
						FirstName: rm.Remittance.Receiver.FirstName,
						LastName:  rm.Remittance.Receiver.LastName,
					},
					Address: RMSAddress{
						Street:     rm.Remittance.Receiver.Address1,
						City:       rm.Remittance.Receiver.City,
						State:      rm.Remittance.Receiver.State,
						PostalCode: rm.Remittance.Receiver.PostalCode,
						CountryCode: RMSCountryCode{
							IsoCode: RMSIsoCode{
								Country:  rm.Remittance.Receiver.Country,
								Currency: rm.Remittance.DestAmt.CurrencyCode(),
							},
						},
					},
					Mobile: RMSMobilePhone{
						RawPhone: json.RawMessage(`{"phone_number":{"country_code":"PHP","national_number":"3242432"}}`),
						Phone: RMSPhoneNumber{
							CountryCode: "PHP",
							Number:      "12345",
						},
					},
				},
				Financials: RMSFinancials{
					Taxes: RMSTaxes{
						TaxWorksheet: "100000",
					},
					GrossTotal: json.Number("100000"),
					PayAmount:  json.Number("100000"),
					Principal:  json.Number("100000"),
					Charges:    json.Number("100000"),
					Tolls:      json.Number("100000"),
				},
				Payment: RMSPaymentDetails{
					SenderDestCountry: RMSCountryCurrency{
						IsoCode: RMSIsoCode{
							Currency: "PHP",
						},
					},
					OrigCountry: RMSCountryCurrency{
						IsoCode: RMSIsoCode{
							Currency: "PHP",
						},
					},
					DestCountry: RMSCountryCurrency{
						IsoCode: RMSIsoCode{
							Currency: "PHP",
						},
					},
					OriginatingCity: "city",
					TransactionType: "WMF",
					ExchangeRate:    json.Number("100000"),
				},
				FilingDate:       "09-27-21 ",
				FilingTime:       "0124A EDT",
				MoneyTransferKey: "3623945513",
				PayStatus:        "W/C",
				Mtcn:             rm.RemcoControlNo,
				NewMtcn:          random.InvitationCode(10),
				Fusion: RMSFusion{
					FusionStatus: "W/C",
				},
			},
		}
		rbb, err = json.Marshal(rb)
		if err != nil {
			return nil, err
		}
	case "pay":
		rb := &struct {
			RMPConfirmedDetails *RMPConfirmedDetails `json:"confirmed_details"`
		}{
			RMPConfirmedDetails: &RMPConfirmedDetails{
				PaidDateTime: "2021-09-28 10:28",
			},
		}
		rbb, err = json.Marshal(rb)
		if err != nil {
			return nil, err
		}
	case "login":
		rb := &Customer{
			FrgnRefNo:      "frgnRefNo",
			CustomerCode:   "customerCode",
			LastName:       "Doe",
			FirstName:      "John",
			Birthdate:      "2020/12/12",
			Nationality:    "phillipine",
			PresentAddress: "address",
			Occupation:     "occupation",
			EmployerName:   "empName",
			CustomerIDNo:   "custIdNo",
			Wucardno:       "wuCardNo",
			Debitcardno:    "debitCardNo",
			Loyaltycardno:  "Loyaltycardno",
		}
		rbb, err = json.Marshal(rb)
		if err != nil {
			return nil, err
		}
	case "report":
		rs, err := m.st.ListRemitHistory(httpreq.Context(), storage.LRHFilter{
			Partner: static.WUCode,
		})
		if err != nil {
			return nil, err
		}
		ts := []WUTransaction{}
		for _, r := range rs {
			var txndate string
			if r.TxnCompletedTime.Valid {
				txndate = r.TxnCompletedTime.Time.Format("2006-01-02 15:04:05")
			}
			t := WUTransaction{
				TxnDate:     txndate,
				MTCN:        r.RemcoControlNo,
				Principal:   r.Remittance.SourceAmt.Number(),
				SvcFee:      json.Number("300"),
				Currency:    r.Remittance.SourceAmt.CurrencyCode(),
				TxnType:     r.Remittance.TxnType,
				DateClaimed: txndate,
				OrderID:     r.DsaOrderID,
			}
			ts = append(ts, t)
		}
		p := &wuHistory{
			Status: "1",
			Msg:    "Success",
			Data:   ts,
		}
		rbb, err = json.Marshal(p)
		if err != nil {
			return nil, err
		}
	case "sdq":
		d := &SDQsRequest{}
		if err := json.Unmarshal(phreq.WU.Body.Param, d); err != nil {
			return nil, err
		}
		switch d.SDQType {
		case "id":
			rbb = idJSON
		case "occupation":
			rbb = ocJSON
		case "position":
			rbb = poJSON
		case "purpose":
			rbb = puJSON
		case "relationship":
			rbb = reJSON
		case "source_of_fund":
			rbb = sfJSON
		}
	case "GetCountriesCurrencies":
		rbb = ccJSON
	default:
		return nil, fmt.Errorf("mock has not been implemented for mod request: %s", modReq)
	}

	if phreq.WU.Body.Request == "report" {
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader(rbb)),
		}, nil
	}

	re := &struct {
		WU *response `json:"uspwuapi"`
	}{
		WU: &response{
			Body: rbb,
			Header: struct {
				ErrorCode string `json:"errorcode"`
				Message   string `json:"message"`
			}{
				// success
				ErrorCode: "1",

				// for testing using partner error
				// ErrorCode: "WU",
				// Message:   "E9205 INVALID REQUEST",
			},
		},
	}
	sc := 200
	if phreq.WU.Body.Request == "feeinquiry" && m.crPtnrErr ||
		phreq.WU.Body.Request == "SMstore" && m.cfCrPtnrErr ||
		phreq.WU.Body.Request == "pay" && m.cfDsbPtnrErr ||
		phreq.WU.Body.Request == "search" && m.dsbPtnrErr {
		re.WU.Header.ErrorCode = "WU"
		re.WU.Header.Message = "E0000 Transaction does not exists"
		sc = 400
	}
	b, err := json.Marshal(re)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		StatusCode: sc,
		Body:       ioutil.NopCloser(bytes.NewReader(b)),
	}, nil
}

func (m *HTTPMock) bpResp(act string) ([]byte, error) {
	switch act {
	case "inquire":
		if !m.dsbPtnrErr {
			rb := &BPInquireResponseBody{
				Code: "200",
				Msg:  "IN PROCESS:TRANSACTION PROCESS ONGOING",
				Result: BPInquireResult{
					Status:            "T",
					Desc:              "TRANSMIT",
					ControlNo:         "CTRL1",
					RefNo:             "1",
					ClientReferenceNo: "CL1",
					PnplAmt:           "1000.00",
					SenderName:        "John, Michael Doe",
					RcvName:           "Jane, Emily Doe",
					Address:           "PLA",
					Currency:          "PHP",
					ContactNumber:     "09190000000",
					OrgnCtry:          "SINGAPORE",
					DestCtry:          "PHILIPPINES",
				},
				RemcoID: "2",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: BPErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, dsbPtnrErr
		}
	case "payout":
		if !m.cfDsbPtnrErr {
			rb := &BPPayoutResponseBody{
				Code: "200",
				Msg:  "Success",
				Result: BPPayoutResult{
					Status:            "T",
					Desc:              "TRANSMIT",
					ControlNo:         "CTRL1",
					RefNo:             "1",
					ClientReferenceNo: "CL1",
					PnplAmt:           "1000.00",
					SenderName:        "MARLON GAVINO REYES VILLA",
					RcvName:           "DANIEL, JENI",
					Address:           "null",
					Currency:          "null",
					ContactNumber:     "null",
					RcvLastName:       "null",
					RcvFirstName:      "null",
					OrgnCtry:          "SINGAPORE",
					DestCtry:          "PHILIPPINES",
					TxnDate:           "null",
					IsDomestic:        "null",
					IDType:            "null",
					RcvCtryCode:       "null",
					RcvStateID:        "null",
					RcvStateName:      "null",
					RcvCityID:         "null",
					RcvCityName:       "null",
					RcvIDType:         "null",
					RcvIsIndiv:        "null",
					PrpsOfRmtID:       "null",
					DsaCode:           "TEST_DSA",
					DsaTrxType:        "digital",
				},
				RemcoID: "2",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: BPErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfDsbPtnrErr
		}
	}
	return nil, fmt.Errorf("no response for action: %v", act)
}

func (m *HTTPMock) icResp(act string) ([]byte, error) {
	switch act {
	case "inquire":
		if !m.dsbPtnrErr {
			rb := &InstaCashInquireResponseBody{
				Code:    "1",
				Message: "Client Details",
				Result: InstaCashResult{
					ControlNumber:      "CTRL1",
					ReferenceNumber:    "REF1",
					OriginatingCountry: "UNITED ARAB EMIRATES",
					DestinationCountry: "PHILIPPINES",
					SenderName:         "John Michael Doe",
					ReceiverName:       "Jane Emily Doe",
					PrincipalAmount:    "1000.00",
					Currency:           "PHP",
					Purpose:            "FAMILY MAINTENANCE",
					Status:             "For Pick Up",
				},
				RemcoID: "16",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: ICErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, dsbPtnrErr
		}
	case "payout":
		if !m.cfDsbPtnrErr {
			rb := &InstaCashPayoutResponseBody{
				Code:    "1",
				Message: "Transaction Status",
				Result: InstaCashPayoutResult{
					ControlNumber: "CTRL1",
					Status:        true,
					Remarks:       "Succesful Payout",
				},
				RemcoID: "16",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: ICErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfDsbPtnrErr
		}
	}
	return nil, fmt.Errorf("no response for action: %v", act)
}

func (m *HTTPMock) jprResp(act string) ([]byte, error) {
	switch act {
	case "inquire":
		if !m.dsbPtnrErr {
			rb := &JPRInquireResponseBody{
				Code:    "0",
				Message: "Available For Pickup",
				Result: JPRResult{
					ControlNumber:      "CTRL1",
					ReferenceNumber:    "REF1",
					OriginatingCountry: "UNITED ARAB EMIRATES",
					DestinationCountry: "PHILIPPINES",
					SenderName:         "John Michael Doe",
					ReceiverName:       "Jane Emily Doe",
					PrincipalAmount:    "1000.00",
					Currency:           "PHP",
					PayTokenId:         "729764",
				},
				RemcoID: "17",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: JPRErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, dsbPtnrErr
		}
	case "payout":
		if !m.cfDsbPtnrErr {
			rb := &JPRPayoutResponseBody{
				Code:    "1",
				Message: "Success",
				Result: JPRPayoutResult{
					ControlNumber: "CTRL1",
					Status:        "paid",
				},
				RemcoID: "17",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: JPRErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfDsbPtnrErr
		}
	}
	return nil, fmt.Errorf("no response for action: %v", act)
}

func (m *HTTPMock) usscResp(act string) ([]byte, error) {
	switch act {
	case "inquire":
		if !m.dsbPtnrErr {
			rb := &USSCInquireResponseBody{
				Code: "000000",
				Msg:  "OK",
				Result: USSCInqResult{
					RcvName:            "Iglesia, Julius James",
					ControlNo:          "CTRL1",
					PrincipalAmount:    "1000.00",
					ContactNumber:      "0922261616161",
					RefNo:              "1",
					SenderName:         "John Michael Doe",
					TrxDate:            "20211214",
					SenderLastName:     "Doe",
					SenderFirstName:    "John",
					SenderMiddleName:   "Michael",
					ServiceCharge:      "1.00",
					TotalAmount:        "1001.00",
					ReceiverFirstName:  "CAMITAN",
					ReceiverMiddleName: "ALVIN",
					ReceiverLastName:   "JOMAR TE TEST",
					RelationTo:         "Family",
					PurposeTransaction: "Family Support/Living Expenses",
				},
				RemcoID: "10",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: USSCErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, dsbPtnrErr
		}
	case "payout":
		if !m.cfDsbPtnrErr {
			rb := &USSCPayoutResponseBody{
				Code:    "000000",
				Message: "Successfully Payout!",
				Result: USSCPayResult{
					Spcn:            "SP697342605141",
					SendPk:          "",
					SendPw:          "",
					SendDate:        "0",
					SendJournalNo:   "0",
					SendLastName:    "Doe",
					SendFirstName:   "John",
					SendMiddleName:  "Michael",
					PayAmount:       "1000.00",
					SendFee:         "0.00",
					SendVat:         "0.00",
					SendFeeAfterVat: "0.00",
					SendTotalAmount: "0.00",
					PayPk:           "",
					PayPw:           "",
					PayLastName:     "CAMITAN",
					PayFirstName:    "ALVIN",
					PayMiddleName:   "",
					Relationship:    "",
					Purpose:         "",
					PromoCode:       "",
					PayBranchCode:   "",
					Remarks:         "",
					OrNo:            "",
					OboBranchCode:   "",
					OboUserID:       "",
					Message:         "0000: ACCEPTED - PAID OUT SUCCESSFULLY",
					Code:            "0",
					NewScreen:       "0",
					JournalNo:       "011330407",
					ProcessDate:     "20200703",
					ReferenceNo:     "1",
				},
				RemcoID: 10,
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: USSCErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfDsbPtnrErr
		}
	case "fee-inquiry":
		if !m.crPtnrErr {
			rb := &USSCFeeInquiryRespBody{
				Code:    "200",
				Message: "OK",
				Result: USSCFeeInquiryResult{
					PnplAmount:    "1000.00",
					ServiceCharge: "1.00",
					Msg:           "",
					Code:          "0",
					NewScreen:     "0",
					JournalNo:     "000000202",
					ProcessDate:   "null",
					RefNo:         "1",
					TotAmount:     "1001.00",
					SendOTP:       "Y",
				},
				RemcoID: 10,
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: USSCErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, crPtnrErr
		}
	case "send":
		if !m.cfCrPtnrErr {
			rb := &USSCSendResponseBody{
				Code:    "000000",
				Message: "OK",
				Result: USSCSendResult{
					ControlNo:          "CTRL1",
					TrxDate:            "2021-12-22",
					SendFName:          "John",
					SendMName:          "Michael",
					SendLName:          "Doe",
					PrincipalAmount:    "1000.00",
					ServiceCharge:      "1.00",
					TotalAmount:        "1001.00",
					RecFName:           "ALVIN",
					RecMName:           "JOMAR TE TEST",
					RecLName:           "CAMITAN",
					ContactNumber:      "0922261616161",
					RelationTo:         "Family",
					PurposeTransaction: "Family Support/Living Expenses",
					ReferenceNo:        "1",
				},
				RemcoID: "10",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: USSCErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfDsbPtnrErr
		}
	}
	return nil, fmt.Errorf("no response for action: %v", act)
}

func (m *HTTPMock) ayannahResp(act string) ([]byte, error) {
	switch act {
	case "send":
		if !m.cfCrPtnrErr {
			rb := &AYANNAHSendResponseBody{
				Code:    "200",
				Message: "Success",
				Result: AYANNAHSendResult{
					Message:            "Successfully Sendout.",
					ID:                 "7048",
					LocationID:         "191",
					UserID:             "1893",
					TrxDate:            "2021-07-13",
					CurrencyID:         "1",
					RemcoID:            "22",
					TrxType:            "1",
					IsDomestic:         "1",
					CustomerID:         "7712780",
					CustomerName:       "Cortez, Fernando",
					ControlNumber:      "CTRL1",
					SenderName:         "Mercado, Marites Cueto",
					ReceiverName:       "Cortez, Fernando",
					PrincipalAmount:    "1000.00",
					ServiceCharge:      "1.00",
					DstAmount:          "0.00",
					TotalAmount:        "1001.00",
					McRate:             "1.00",
					BuyBackAmount:      "1.00",
					RateCategory:       "required",
					McRateID:           "1",
					OriginatingCountry: "Philippines",
					DestinationCountry: "Philippines",
					PurposeTransaction: "Family Support/Living Expenses",
					SourceFund:         "Salary/Income",
					Occupation:         "Unemployed",
					RelationTo:         "Family",
					BirthDate:          "2000-12-15",
					BirthPlace:         "MALOLOS,BULACAN",
					BirthCountry:       "Philippines",
					IDType:             "PASSPORT",
					IDNumber:           "PRND32200265569P",
					Address:            "18 SITIO PULO",
					Barangay:           "BULIHAN",
					City:               "MALOLOS",
					Province:           "BULACAN",
					Country:            "PH",
					ContactNumber:      "09265454935",
					CurrentAddress: NonexAddress{
						Address1: "Marcos Highway",
						Address2: "null",
						Barangay: "Mayamot",
						City:     "ERMITA",
						Province: "MANILA METROPOLITAN",
						ZipCode:  "1000",
						Country:  "PH",
					},
					PermanentAddress: NonexAddress{
						Address1: "Marcos Highway",
						Address2: "null",
						Barangay: "Mayamot",
						City:     "ERMITA",
						Province: "MANILA METROPOLITAN",
						ZipCode:  "1000",
						Country:  "PH",
					},
					RiskScore:         "1",
					RiskCriteria:      "0",
					ClientReferenceNo: "1",
					FormType:          "0",
					FormNumber:        "0",
					PayoutType:        "1",
					RemoteLocationID:  "191",
					RemoteUserID:      "1893",
					RemoteIPAddress:   "130.211.2.203",
					IPAddress:         "::1",
					CreatedAt:         time.Now(),
					UpdatedAt:         time.Now(),
					ReferenceNumber:   "REF1",
					ZipCode:           "3000A",
					Status:            "1",
					APIRequest:        "null",
					SapForm:           "null",
					SapFormNumber:     "null",
					SapValidID1:       "null",
					SapValidID2:       "null",
					SapOboLastName:    "null",
					SapOboFirstName:   "null",
					SapOboMiddleName:  "null",
					AyannahStatus:     "NEW",
				},
				RemcoID: "22",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: AYANNAHErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfDsbPtnrErr
		}
	case "inquire":
		if !m.dsbPtnrErr {
			rb := &AYANNAHInquireResponseBody{
				Code:    "200",
				Message: "Success",
				Result: AYANNAHInquireResult{
					ResponseCode:       "AVAILABLE",
					ResponseMessage:    "Transaction Available For Payout.",
					ControlNumber:      "CTRL1",
					PrincipalAmount:    "1000.00",
					Currency:           "PHP",
					CreationDate:       "2022-01-14 14:57:26 ",
					ReceiverName:       "Octaviano Luis Rafael",
					SenderName:         "John, Michael Doe",
					Address:            "IMORTALONE",
					City:               "ANTIPOLO,RIZAL",
					Country:            "PH",
					ZipCode:            "null",
					OriginatingCountry: "PH",
					DestinationCountry: "PH",
					ContactNumber:      "null",
					ReferenceNumber:    "1",
				},
				RemcoID: "22",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: AYANNAHErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, dsbPtnrErr
		}
	case "payout":
		if !m.cfDsbPtnrErr {
			rb := &AYANNAHPayoutResponseBody{
				Code:    "200",
				Message: "Success",
				Result: AYANNAHPayoutResult{
					Message: "Successfully Payout.",
				},
				RemcoID: "22",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: AYANNAHErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfDsbPtnrErr
		}
	}
	return nil, fmt.Errorf("no response for action: %v", act)
}

func (m *HTTPMock) ieResp(act string) ([]byte, error) {
	switch act {
	case "send":
		if !m.cfCrPtnrErr {
			rb := &IESendResponse{
				Code:    "200",
				Message: "Good",
				Result: IESendResult{
					ID:                 "7288",
					LocationID:         "192",
					UserID:             "1894",
					TrxDate:            "2021-07-14",
					CurrencyID:         "1",
					RemcoID:            "25",
					TrxType:            "1",
					IsDomestic:         "1",
					CustomerID:         "7712781",
					CustomerName:       "Cortez, Fernandos",
					ControlNumber:      "CTRL1",
					SenderName:         "John, Michael Doe",
					ReceiverName:       "Jane, Emily Doe",
					PrincipalAmount:    "1000.00",
					ServiceCharge:      "1.00",
					DstAmount:          "0.00",
					TotalAmount:        "1001.00",
					McRate:             "1.00",
					BuyBackAmount:      "0.00",
					RateCategory:       "required",
					McRateID:           "1",
					OriginatingCountry: "Philippines",
					DestinationCountry: "Philippines",
					PurposeTransaction: "Family Support/Living Expenses",
					SourceFund:         "Salary/Income",
					Occupation:         "Unemployed",
					RelationTo:         "Family",
					BirthDate:          "2000-12-16",
					BirthPlace:         "MALOLOS,BULACANS",
					BirthCountry:       "Philippines",
					IDType:             "Postal ID",
					IDNumber:           "PRND32200265569Q",
					Address:            "18 SITIO PULO",
					Barangay:           "BULIHAN",
					City:               "MALOLOS",
					Province:           "BULIHAN",
					Country:            "Philippines",
					ContactNumber:      "09265454936",
					CurrentAddress:     "null",
					PermanentAddress:   "null",
					RiskScore:          "1",
					RiskCriteria:       "1",
					ClientReferenceNo:  "REF1",
					FormType:           "0",
					FormNumber:         "0",
					PayoutType:         "1",
					RemoteLocationID:   "192",
					RemoteUserID:       "1894",
					RemoteIPAddress:    "130.211.2.204",
					IPAddress:          "::1",
					CreatedAt:          time.Now(),
					UpdatedAt:          time.Now(),
					ReferenceNumber:    "1",
					ZipCode:            "3000B",
					Status:             "1",
					APIRequest:         "null",
					SapForm:            "null",
					SapFormNumber:      "null",
					SapValidID1:        "null",
					SapValidID2:        "null",
					SapOboLastName:     "null",
					SapOboFirstName:    "null",
					SapOboMiddleName:   "null",
				},
				RemcoID: "24",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: IEErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfDsbPtnrErr
		}
	case "inquire":
		if !m.dsbPtnrErr {
			rb := &IEInquireResponse{
				Code:    "200",
				Message: "Success",
				Result: IEInquireResult{
					ControlNumber:      "CTRL1",
					TrxDate:            "07/03/2021",
					PrincipalAmount:    "1000.00",
					Currency:           "PHP",
					ReceiverName:       "Jane Emily Doe",
					SenderName:         "John Michael Doe",
					Address:            "TBILISI,KOKHREIDZIS - 8",
					Country:            "GEO",
					OriginatingCountry: "GEO",
					DestinationCountry: "PHL",
					ContactNumber:      "",
					ReferenceNumber:    "1",
				},
				RemcoID: "24",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: IEErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, dsbPtnrErr
		}
	case "payout":
		if !m.cfDsbPtnrErr {
			rb := &IEPayoutResponse{
				Code:    "200",
				Message: "Good",
				Result: IEPayoutResult{
					ID:                 "7287",
					LocationID:         "371",
					UserID:             "5500",
					TrxDate:            "2021-06-03",
					CurrencyID:         "1",
					RemcoID:            "24",
					TrxType:            "2",
					IsDomestic:         "1",
					CustomerID:         "6925594",
					CustomerName:       "Levy Robert Sogocio",
					ControlNumber:      "CTRL1",
					SenderName:         "Reyes, Ana,",
					ReceiverName:       "Bayuga, Mary Monica",
					PrincipalAmount:    "1000.00",
					ServiceCharge:      "1.00",
					DstAmount:          "0.00",
					TotalAmount:        "1001.00",
					McRate:             "0.00",
					BuyBackAmount:      "0.00",
					RateCategory:       "0",
					McRateID:           "0",
					OriginatingCountry: "Philippines",
					DestinationCountry: "PH",
					PurposeTransaction: "Family Support/Living Expenses",
					SourceFund:         "Savings",
					Occupation:         "OTH",
					RelationTo:         "Family",
					BirthDate:          "1995-04-14",
					BirthPlace:         "TAGUDIN,ILOCOS SUR",
					BirthCountry:       "PH",
					IDType:             "LICENSE",
					IDNumber:           "B83180608851",
					Address:            "MAIN ST",
					Barangay:           "TALLAOEN",
					City:               "AKLAN CITY",
					Province:           "AKLAN CITY",
					Country:            "PH",
					ContactNumber:      "09516738640",
					CurrentAddress: NonexAddress{
						Address1: "#32 Griffin",
						Address2: "",
						Barangay: "Pinagbuhatan",
						City:     "ERMITA",
						Province: "MANILA",
						ZipCode:  "1000A",
						Country:  "Philippines",
					},
					PermanentAddress: NonexAddress{
						Address1: "#32 Griffin",
						Address2: "",
						Barangay: "Pinagbuhatan",
						City:     "ERMITA",
						Province: "MANILA",
						ZipCode:  "1000A",
						Country:  "Philippines",
					},
					RiskScore:         "0",
					RiskCriteria:      "0",
					ClientReferenceNo: "7884447474",
					FormType:          "0",
					FormNumber:        "0",
					PayoutType:        "1",
					RemoteLocationID:  "371",
					RemoteUserID:      "5684",
					RemoteIPAddress:   "130.211.2.187",
					IPAddress:         "130.211.2.187",
					CreatedAt:         time.Now(),
					UpdatedAt:         time.Now(),
					ReferenceNumber:   "848jjfu23u333",
					ZipCode:           "36989",
					Status:            "1",
					APIRequest:        "null",
					SapForm:           "null",
					SapFormNumber:     "null",
					SapValidID1:       "null",
					SapValidID2:       "null",
					SapOboLastName:    "null",
					SapOboFirstName:   "null",
					SapOboMiddleName:  "null",
				},
				RemcoID: "24",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: IEErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfDsbPtnrErr
		}
	}
	return nil, fmt.Errorf("no response for action: %v", act)
}

func (m *HTTPMock) untResp(act string) ([]byte, error) {
	switch act {
	case "inquire":
		if !m.dsbPtnrErr {
			rb := &UNTInquireResponseBody{
				Code:    "00000000",
				Message: "Success",
				Result: UNTResult{
					ResponseCode:       "00000000",
					ControlNumber:      "CTRL1",
					PrincipalAmount:    "1000.00",
					Currency:           "PHP",
					CreationDate:       "2021-06-15T19:10:16.000-0400",
					ReceiverName:       "Jane Emily Doe",
					Address:            "TEST",
					City:               "MORONG",
					Country:            "PH",
					SenderName:         "John Michael Doe",
					ZipCode:            "36978",
					OriginatingCountry: "US",
					DestinationCountry: "PH",
					ContactNumber:      "1540254852",
					FmtSenderName:      "John, Michael Doe",
					FmtReceiverName:    "Jane, Emily Doe",
				},
				RemcoID: "20",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: UNTErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, dsbPtnrErr
		}
	case "payout":
		if !m.cfDsbPtnrErr {
			rb := &UNTPayoutResponseBody{
				Code:    "00000000",
				Message: "Success",
				Result: UNTPayoutResult{
					ResponseCode:       "00000000",
					ControlNumber:      "CTRL1",
					PrincipalAmount:    "1000.00",
					Currency:           "PHP",
					CreationDate:       "2021-06-15T19:10:16.000-0400",
					ReceiverName:       "RAMSES COMODO",
					Address:            "TEST",
					City:               "MORONG",
					Country:            "PH",
					SenderName:         "FERDINAND CORTES",
					ZipCode:            "36978",
					OriginatingCountry: "US",
					DestinationCountry: "PH",
					ContactNumber:      "1540254852",
					FmtSenderName:      "CORTES, FERDINAND, ",
					FmtReceiverName:    "COMODO, RAMSES, ",
				},
				RemcoID: "20",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: UNTErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfDsbPtnrErr
		}
	case "countries":
		rb := &UNTBaseResponseList{
			Data: []UNTBaseResponse{
				{
					Code: "AS",
					Name: "AMERICAN SAMOA",
				},
				{
					Code: "AU",
					Name: "AUSTRALIAN",
				},
			},
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "currencies":
		rb := &UNTBaseResponseList{
			Data: []UNTBaseResponse{
				{
					Code: "USD",
					Name: "US DOLLAR",
				},
				{
					Code: "PHP",
					Name: "PHILIPPINE PESOS",
				},
			},
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "occupations":
		rb := &UNTCommonResponseList{
			Data: []UNTCommonResponse{
				{
					Country: "OTHER",
					Code:    "OTH",
				},
				{
					Country: "HOUSEWIFE",
					Code:    "HW",
				},
			},
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "ids":
		rb := &UNTCommonResponseList{
			Data: []UNTCommonResponse{
				{
					Country: "LICENSE",
					Code:    "PHILIPPINES",
				},
				{
					Country: "PASSPORT",
					Code:    "PHILIPPINES",
				},
			},
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "ph-states":
		rb := &UNTStateResponseList{
			Data: []UNTStatesResponse{
				{
					Country:   "PHILIPPINES",
					StateName: "DAVAO OCCIDENTAL",
					UtlCode:   "PH-DVO",
				},
				{
					Country:   "PHILIPPINES",
					StateName: "PH-NA",
					UtlCode:   "PH-NA",
				},
			},
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "usa-states":
		rb := &UNTUsStateResponseList{
			Data: []UNTUsStatesResponse{
				{
					Country:   "USA",
					StateName: "NEW JERSEY",
					UtlCode:   "NJ",
				},
				{
					Country:   "USA",
					StateName: "PENNSYLVANIA",
					UtlCode:   "PA",
				},
			},
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	}
	return nil, fmt.Errorf("no response for action: %v", act)
}

func (m *HTTPMock) cebResp(act string) ([]byte, error) {
	if strings.HasPrefix(act, "get-service-fee") {
		act = "get-service-fee"
	}
	if strings.HasPrefix(act, "find-client") {
		act = "find-client"
	}
	if strings.HasPrefix(act, "send-currency-collection") {
		act = "send-currency-collection"
	}
	if strings.HasPrefix(act, "beneficiary-by-sender") {
		act = "beneficiary-by-sender"
	}
	switch act {
	case "add-beneficiary":
		rb := &CebAddBfResp{
			Code:    0,
			Message: "Successful",
			Result: ABResult{
				ResultStatus:  "Successful",
				MessageID:     0,
				LogID:         0,
				BeneficiaryID: 8595,
			},
			RemcoID: 9,
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "beneficiary-by-sender":
		rb := &CebFindBFRes{
			Code:    "0",
			Message: "Successful",
			Result: FBFResult{
				Beneficiary: []Beneficiary{
					{
						BeneficiaryID:  "3663",
						FirstName:      "GREMAR",
						MiddleName:     "REAS",
						LastName:       "NAPARATE",
						BirthDate:      "1994-11-04T00:00:00",
						StateIDAddress: "0",
						CPCountry: CrtyID{
							CountryID: 0,
						},
						TPCountry: CrtyID{
							CountryID: 0,
						},
						CtryAddress: CrtyID{
							CountryID: 0,
						},
						BirthCountry: CrtyID{
							CountryID: 0,
						},
					},
				},
			},
			RemcoID: "9",
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "add-client":
		rb := &CebAddClientResp{
			Code:    0,
			Message: "Successful",
			Result: RUResult{
				ResultStatus: "Successful",
				MessageID:    0,
				LogID:        0,
				ClientID:     3673,
				ClientNo:     "EWFHM0000070155828",
			},
			RemcoID: 9,
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "find-client":
		rb := &CebFindClientRespBody{
			Code:    0,
			Message: "Successful",
			Result: GUResult{
				Client: Client{
					ClientID:     3663,
					ClientNumber: "EWFHL0000070155824",
					FirstName:    "GREMAR",
					MiddleName:   "REAS",
					LastName:     "NAPARATE",
					BirthDate:    "1994-11-04T00:00:00",
					CPCountry: CrtyID{
						CountryID: 0,
					},
					TPCountry: CrtyID{
						CountryID: 0,
					},
					CtryAddress: CrtyID{
						CountryID: 0,
					},
					CSOfFund: CSOfFund{
						SourceOfFundID: 0,
					},
				},
			},
			RemcoID: 9,
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "get-service-fee":
		if !m.crPtnrErr {
			rb := &CebuanaSFInquiryRespBody{
				Code:    "0",
				Message: "Successful",
				Result: CebuanaSFInquiryResult{
					ResultStatus: "Successful",
					MessageID:    "0",
					LogID:        "0",
					ServiceFee:   "1.00",
				},
				RemcoID: "9",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: CEBErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, crPtnrErr
		}
	case "send":
		if !m.cfCrPtnrErr {
			rb := &CebuanaSendResponseBody{
				Code:    "0",
				Message: "Successful",
				Result: CebuanaSendResult{
					ResultStatus: "Successful",
					MessageID:    "0",
					LogID:        "0",
					ControlNo:    "CTRL1",
					ServiceFee:   "1.00",
				},
				RemcoID: "9",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: CEBErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfCrPtnrErr
		}
	case "inquire":
		if !m.dsbPtnrErr {
			rb := &CEBInquireResponseBody{
				Code:    "1",
				Message: "Successful",
				Result: CEBInquireResult{
					ResultStatus:      "Successful",
					MessageID:         "125",
					LogID:             "0",
					ClientReferenceNo: "1",
					ControlNo:         "CTRL1",
					SenderName:        "John Michael Doe",
					RcvName:           "ESTOCAPIO, SHAIRA MIKA, MADJALIS",
					PnplAmt:           "1000",
					ServiceCharge:     "1",
					BirthDate:         "1980-08-10T00:00:00",
					Currency:          "PHP",
					BeneficiaryID:     "5342",
					RemStatusID:       "1",
					RemStatusDes:      "Outstanding",
				},
				RemcoID: "9",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: CEBErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, dsbPtnrErr
		}
	case "payout":
		if !m.cfDsbPtnrErr {
			rb := &CEBPayoutResponseBody{
				Code:    "0",
				Message: "Successful",
				Result: CEBPayoutResult{
					Message: "Successfully Payout!",
				},
				RemcoID: "9",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: CEBErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfDsbPtnrErr
		}
	case "get-country":
		rb := &CEBCountryBaseResponse{
			Code:    "0",
			Message: "Successful",
			Result: CEBCountryListResponse{
				Country: []CEBCountryResponse{
					{
						CountryID:         "1",
						CountryName:       "AFGHANISTAN",
						CountryCodeAplha2: "AF",
						CountryCodeAlpha3: "AFG",
						PhoneCode:         "93",
					},
					{
						CountryID:         "2",
						CountryName:       "ALBANIA",
						CountryCodeAplha2: "AL",
						CountryCodeAlpha3: "ALB",
						PhoneCode:         "355",
					},
				},
			},
			RemcoID: "9",
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "send-currency-collection":
		rb := &CEBCurrencyBaseResponse{
			Code:    "0",
			Message: "Successful",
			Result: CEBCurrencyListResponse{
				Currency: CEBCurrencyResponse{
					CurrencyID:  "6",
					Code:        "Php",
					Description: "PHILIPPINE PESO",
				},
			},
			RemcoID: "9",
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "get-source-of-fund":
		rb := &CEBSourceFundsBaseResponse{
			Code:    "0",
			Message: "Successful",
			Result: CEBSourceFundsListResponse{
				ClientSourceOfFund: []CEBSourceFundsResponse{
					{
						SourceOfFundID: "1",
						SourceOfFund:   "Employed",
					},
					{
						SourceOfFundID: "2",
						SourceOfFund:   "Self-Employed",
					},
				},
			},
			RemcoID: "9",
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "identification-types":
		rb := &CEBIdTypesResponseList{
			{
				IdentificationTypeID: "1",
				Description:          "Bank Credit Card",
				SmsCode:              "CR",
			},
			{
				IdentificationTypeID: "3",
				Description:          "Barangay ID",
				SmsCode:              "BY",
			},
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	}
	return nil, fmt.Errorf("no response for action: %v", act)
}

func (m *HTTPMock) cebIntResp(act string) ([]byte, error) {
	switch act {
	case "inquire":
		if !m.dsbPtnrErr {
			rb := &CEBINTInquireResponseBody{
				Code:    "1",
				Message: "Successful",
				Result: CEBINTInquireResult{
					IsDomestic:                  "0",
					ResultStatus:                "Successful",
					MessageID:                   "0",
					LogID:                       "0",
					ClientReferenceNo:           "1",
					ControlNumber:               "CTRL1",
					SenderName:                  "John Michael Doe",
					ReceiverName:                "ESTOCAPIO, SHAIRA MIKA, MADJALIS",
					PrincipalAmount:             "1000",
					ServiceCharge:               "1",
					BirthDate:                   "1980-08-10T00:00:00",
					Currency:                    "PHP",
					BeneficiaryID:               "5342",
					RemittanceStatusID:          "1",
					RemittanceStatusDescription: "Outstanding",
				},
				RemcoID: "9",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: CEBINTErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, dsbPtnrErr
		}
	case "payout":
		if !m.cfDsbPtnrErr {
			rb := &CEBINTPayoutResponseBody{
				Code:    "0",
				Message: "Successful",
				Result: CEBINTPayoutResult{
					Message: "Successfully Payout!",
				},
				RemcoID: "9",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: CEBINTErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfDsbPtnrErr
		}
	}
	return nil, fmt.Errorf("no response for action: %v", act)
}

func (m *HTTPMock) wiseResp(act string) ([]byte, error) {
	cntryCode := ""
	if strings.HasPrefix(act, "utils/state/") {
		cntryCode = strings.Replace(act, "utils/state/", "", 1)
	}
	switch act {
	case "oauth/token":
		m.reqOrder = append(m.reqOrder, "auth")
		rb := &WISEGetTokensResp{
			AccessToken:  "token",
			RefreshToken: "refr-token",
			ExpiresIn:    "43200",
			TokenType:    "Bearer",
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "accounts":
		m.reqOrder = append(m.reqOrder, "create")
		if !m.conflictErr && !m.authErr {
			rb := &WISECreateUserResp{
				Msg: "Success! User Account Created.",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else if m.conflictErr {
			rb := &WISECreateUserResp{
				Error: "Conflict",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, conflictErr
		} else if m.authErr {
			m.authErr = false
			rb := &WISECreateUserResp{
				Error: "Unauthorized",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, authErr
		}
	case "profiles":
		rb := &WISECreateProfileResp{
			Msg: "Success! User Profile Created.",
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "profiles/personal":
		rb := &WISEGetProfileResp{
			Type: "personal",
			Details: WISEGetPFDetails{
				FirstName:   "Brankas",
				LastName:    "Sender",
				BirthDate:   "1990-01-10",
				PhoneNumber: "+639999999999",
				Occupations: []WISEOccupation{
					{
						Code:   "Software Engineer",
						Format: "FREE_FORM",
					},
				},
				PrimaryAddress: json.Number("37325852"),
			},
			ProfileID: json.Number("16325688"),
			Address: WISEPFAddress{
				Country:   "ph",
				FirstLine: "East Offices Bldg., 114 Aguirre St.,Legaspi Village,",
				PostCode:  "1229",
				City:      "Makati",
			},
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "recipients":
		rb := &WISECreateRecipientResp{
			RecipientID: "148142936",
			Details: WISECRDetails{
				Address: WISECRAddress{
					Country:     "GB",
					CountryCode: "GB",
					FirstLine:   "10 Downing Street",
					PostCode:    "SW1A 2AA",
					City:        "London",
				},
				LegalType:     "PRIVATE",
				AccountNumber: "28821822",
				SortCode:      "231470",
			},
			AccountHolderName: "Brankas Receiver",
			Currency:          "GBP",
			OwnedByCustomer:   false,
			Country:           "GB",
			Msg:               "Success! Recipient Account Created",
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "recipients/refresh-onchange":
		rb := &WISERefreshRecipientResp{
			Requirements: []WISERequirementsResp{
				{
					Type:      "sort_code",
					Title:     "Local bank account",
					UsageInfo: "usage",
					Fields: []WISEField{
						{
							Name: "Recipient type",
							Group: []WISEGroup{
								{
									Key:                "legalType",
									Name:               "Recipient type",
									Type:               "select",
									RefreshReqOnChange: true,
									Required:           true,
									Example:            "example",
									MinLength:          "1",
									MaxLength:          "2",
									ValidationAsync: WISEValidationAsync{
										URL: "url",
										Params: []WISEParam{
											{
												Key:       "key",
												ParamName: "paramname",
												Required:  true,
											},
										},
									},
									ValuesAllowed: []WISEValueAllowed{
										{
											Key:  "PRIVATE",
											Name: "Person",
										},
										{
											Key:  "BUSINESS",
											Name: "Business",
										},
									},
								},
							},
						},
					},
				},
			},
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "recipients/list":
		rb := []WISERecipient{
			{
				RecipientID: "148142936",
				Details: WISEGRDetails{
					AccountNumber:    "28821822",
					SortCode:         "231470",
					HashedByLooseAlg: "d641e415d0503d966ff4a5a7d246ba9511430a258e8a993b53059defb256c448",
				},
				AccountSummary:     "(23-14-70) 28821822",
				LongAccountSummary: "GBP account ending in 1822",
				DisplayFields: []WISEDisplayField{
					{
						Label: "UK Sort code",
						Value: "23-14-70",
					},
					{
						Label: "Account number",
						Value: "28821822",
					},
				},
				FullName: "Brankas Receiver",
				Currency: "GBP",
				Country:  "GB",
			},
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "recipients/12345":
		rb := &WISEDeleteRecipientResp{
			Msg: "Successfully Deleted Recipient Account!",
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "quotes/inquiry":
		rb := &WISEQuoteInquiryResp{
			SourceCurrency: "PHP",
			TargetCurrency: "GBP",
			SourceAmount:   "1500",
			TargetAmount:   "70.21",
			FeeBreakdown: WISEFeeBreakdown{
				Transferwise: "79.66",
				PayIn:        "0",
				Discount:     "0",
				Total:        "79.66",
				PriceSetID:   "132",
				Partner:      "0",
			},
			TotalFee:       "2.49",
			TransferAmount: "97.51",
			PayOut:         "BANK_TRANSFER",
			Rate:           "0.0148552",
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "quotes":
		rb := &WISECreateQuoteResp{
			Requirements: []WISERequirementsResp{
				{
					Type:      "sort_code",
					Title:     "Local bank account",
					UsageInfo: "usage",
					Fields: []WISEField{
						{
							Name: "Recipient type",
							Group: []WISEGroup{
								{
									Key:                "legalType",
									Name:               "Recipient type",
									Type:               "select",
									RefreshReqOnChange: true,
									Required:           true,
									Example:            "example",
									MinLength:          "1",
									MaxLength:          "2",
									ValidationAsync: WISEValidationAsync{
										URL: "url",
										Params: []WISEParam{
											{
												Key:       "key",
												ParamName: "paramname",
												Required:  true,
											},
										},
									},
									ValuesAllowed: []WISEValueAllowed{
										{
											Key:  "PRIVATE",
											Name: "Person",
										},
										{
											Key:  "BUSINESS",
											Name: "Business",
										},
									},
								},
							},
						},
					},
				},
			},
			QuoteSummary: WISEQuoteInquiryResp{
				SourceCurrency: "PHP",
				TargetCurrency: "GBP",
				SourceAmount:   "1500",
				TargetAmount:   "21.1",
				FeeBreakdown: WISEFeeBreakdown{
					Transferwise: "79.66",
					PayIn:        "0",
					Discount:     "0",
					Total:        "79.66",
					PriceSetID:   "132",
					Partner:      "0",
				},
				TotalFee:       "79.66",
				TransferAmount: "1420.34",
				PayOut:         "BANK_TRANSFER",
				Rate:           "0.0148552",
			},
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "quotes/requirements":
		rb := &WISEGetQuoteRequirementsResp{
			Requirements: []WISERequirementsResp{
				{
					Type:      "sort_code",
					Title:     "Local bank account",
					UsageInfo: "usage",
					Fields: []WISEField{
						{
							Name: "Recipient type",
							Group: []WISEGroup{
								{
									Key:                "legalType",
									Name:               "Recipient type",
									Type:               "select",
									RefreshReqOnChange: true,
									Required:           true,
									Example:            "example",
									MinLength:          "1",
									MaxLength:          "2",
									ValidationAsync: WISEValidationAsync{
										URL: "url",
										Params: []WISEParam{
											{
												Key:       "key",
												ParamName: "paramname",
												Required:  true,
											},
										},
									},
									ValuesAllowed: []WISEValueAllowed{
										{
											Key:  "PRIVATE",
											Name: "Person",
										},
										{
											Key:  "BUSINESS",
											Name: "Business",
										},
									},
								},
							},
						},
					},
				},
			},
			Quote: WISEQuote{
				SourceCurrency: "PHP",
				TargetCurrency: "GBP",
				SourceAmount:   "1500",
			},
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, nil
	case "transfers/prepare":
		if !m.crPtnrErr {
			rb := &WISEPrepareTransferResp{
				Requirements: []WISERequirementsResp{
					{
						Type:      "sort_code",
						Title:     "Local bank account",
						UsageInfo: "usage",
						Fields: []WISEField{
							{
								Name: "Recipient type",
								Group: []WISEGroup{
									{
										Key:                "legalType",
										Name:               "Recipient type",
										Type:               "select",
										RefreshReqOnChange: true,
										Required:           true,
										Example:            "example",
										MinLength:          "1",
										MaxLength:          "2",
										ValidationAsync: WISEValidationAsync{
											URL: "url",
											Params: []WISEParam{
												{
													Key:       "key",
													ParamName: "paramname",
													Required:  true,
												},
											},
										},
										ValuesAllowed: []WISEValueAllowed{
											{
												Key:  "PRIVATE",
												Name: "Person",
											},
											{
												Key:  "BUSINESS",
												Name: "Business",
											},
										},
									},
								},
							},
						},
					},
				},
				UpdatedQuoteSummary: WISEQuoteInquiryResp{
					SourceCurrency: "PHP",
					TargetCurrency: "GBP",
					SourceAmount:   "1500",
					TargetAmount:   "21.1",
					FeeBreakdown: WISEFeeBreakdown{
						Transferwise: "79.66",
						PayIn:        "0",
						Discount:     "0",
						Total:        "79.66",
						PriceSetID:   "132",
						Partner:      "0",
					},
					TotalFee:       "79.66",
					TransferAmount: "1420.34",
					PayOut:         "BANK_TRANSFER",
					Rate:           "0.0148552",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: WISEErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, crPtnrErr
		}
	case "transfers/proceed":
		if !m.cfCrPtnrErr {
			rb := &WISEProceedTransferResp{
				TransferID: "12345",
				Details: WISEPCDetails{
					Reference: "54321",
				},
				CustomerTxnID:  "aecd179d",
				RecipientID:    "447769582",
				Status:         "incoming_payment_waiting",
				SourceCurrency: "PHP",
				TargetCurrency: "GBP",
				SourceAmount:   "1500",
				DateCreated:    "2021-03-05T05:38:36.677Z",
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: WISEErr,
					Msg:  "Transaction does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfCrPtnrErr
		}
	case "utils/country":
		if !m.cfCrPtnrErr {
			rb := WISEBaseResponseList{
				{
					Key:  "AX",
					Name: "Åland Islands",
				},
				{
					Key:  "AL",
					Name: "Albania",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: WISEErr,
					Msg:  "Get country does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfCrPtnrErr
		}
	case "utils/state/" + cntryCode:
		if !m.cfCrPtnrErr {
			rb := WISEBaseResponseList{
				{
					Key:  "AL",
					Name: "Alabama",
				},
				{
					Key:  "AK",
					Name: "Alaska",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: WISEErr,
					Msg:  "get state does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfCrPtnrErr
		}
	case "utils/currency":
		if !m.cfCrPtnrErr {
			rb := WISECurrencyResponseList{
				{
					Currency:    "AED",
					Description: "UAE Dirham",
				},
				{
					Currency:    "ARS",
					Description: "Argentine peso",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, nil
		} else {
			rb := &nonexError{
				Code: "1",
				Msg:  "some error",
				Error: ErrorType{
					Type: WISEErr,
					Msg:  "get state does not exists",
				},
			}
			rbb, err := json.Marshal(rb)
			if err != nil {
				return nil, err
			}
			return rbb, cfCrPtnrErr
		}
	}
	return nil, fmt.Errorf("no response for action: %v", act)
}

func (m *HTTPMock) peraHubRemitResp(req *http.Request) ([]byte, error) {
	str := strings.ReplaceAll(req.URL.Path, "/v1/remit/nonex/", "")
	a := strings.Split(str, "/")
	if len(a) < 3 {
		return nil, status.Error(codes.NotFound, "end point not found")
	}
	act := fmt.Sprintf("%s_%s", a[1], a[2])

	switch act {
	case "address_provinces":
		return m.addressProvinces()
	case "address_brgy":
		return m.getBrgyList()
	case "utility_purpose":
		return m.getUtilityPurpose()
	case "utility_relationship":
		return m.getUtilityRelationship()
	case "utility_partner":
		return m.getUtilityPartner()
	case "utility_occupation":
		return m.getUtilityOccupation()
	case "utility_employment":
		return m.getUtilityEmployment()
	case "utility_sourcefund":
		return m.getUtilitySourceFund()
	case "payout_retry":
		return m.retryTransaction()
	case "payout_confirm":
		return m.confirmTransaction()
	case "payout_validate":
		return m.validateTransaction()
	case "payout_inquire":
		return m.inquireTransaction()
	}

	return nil, fmt.Errorf("no response for action: %v", act)
}

func (m *HTTPMock) commonNonexError() ([]byte, error) {
	rb := &nonexError{
		Code: "1",
		Msg:  "some error",
		Error: ErrorType{
			Type: PerahubRemitErr,
			Msg:  "data not found",
		},
	}
	rbb, err := json.Marshal(rb)
	if err != nil {
		return nil, err
	}
	return rbb, dsbPtnrErr
}

func (m *HTTPMock) addressProvinces() ([]byte, error) {
	if m.dsbPtnrErr {
		return m.commonNonexError()
	}

	rb := &GetProvincesCityListResponse{
		Code:    200,
		Message: "Success - PH All Provinces w/ City List",
		Result: []GetProvincesCityListResults{
			{
				Province: "METRO MANILA",
				CityList: []string{
					"MANILA",
					"CITY OF MAKATI",
					"CITY OF MUNTINLUPA",
					"CITY OF PARANAQUE",
					"PASAY CITY",
				},
			},
		},
	}
	rbb, err := json.Marshal(rb)
	if err != nil {
		return nil, err
	}
	return rbb, nil
}

func (m *HTTPMock) getBrgyList() ([]byte, error) {
	if m.dsbPtnrErr {
		return m.commonNonexError()
	}

	rb := &GetBrgyListResponse{
		Code:    200,
		Message: "Success - MANILA Barangays w/ Zipcode",
		Result: []GetBrgyListResults{
			{
				Barangay: "Barangay 1",
				Zipcode:  "1013",
			},
		},
	}
	rbb, err := json.Marshal(rb)
	if err != nil {
		return nil, err
	}
	return rbb, nil
}

func (m *HTTPMock) getUtilityPurpose() ([]byte, error) {
	if m.dsbPtnrErr {
		return m.commonNonexError()
	}

	rb := &GetUtilityPurposeResponse{
		Code:    200,
		Message: "Success - All PURPOSE/s",
		Result: []string{
			"Family Support/Living Expenses",
			"Saving/Investments",
			"Gift",
		},
	}
	rbb, err := json.Marshal(rb)
	if err != nil {
		return nil, err
	}
	return rbb, nil
}

func (m *HTTPMock) getUtilityRelationship() ([]byte, error) {
	if m.dsbPtnrErr {
		return m.commonNonexError()
	}

	rb := &GetUtilityRelationshipResponse{
		Code:    200,
		Message: "Success - All RELATIONSHIP/s",
		Result: []string{
			"Family",
			"Friend",
		},
	}
	rbb, err := json.Marshal(rb)
	if err != nil {
		return nil, err
	}
	return rbb, nil
}

func (m *HTTPMock) getUtilityPartner() ([]byte, error) {
	if m.dsbPtnrErr {
		return m.commonNonexError()
	}

	rb := &GetUtilityPartnerResponse{
		Code:    200,
		Message: "Success - All PARTNER/s",
		Result: []GetUtilityPartnerResult{
			{
				PartnerCode: "DRP",
				PartnerName: "BRANKAS",
				Status:      1,
			},
		},
	}

	rbb, err := json.Marshal(rb)
	if err != nil {
		return nil, err
	}
	return rbb, nil
}

func (m *HTTPMock) getUtilityOccupation() ([]byte, error) {
	if m.dsbPtnrErr {
		return m.commonNonexError()
	}

	rb := &GetUtilityOccupationResponse{
		Code:    200,
		Message: "Success - All OCCUPATION/s",
		Result: []string{
			"Airline/Maritime Employee",
			"Art/Entertainment/Media/Sports Professional",
			"Civil/Government Employee",
			"Domestic Helper",
			"Driver",
		},
	}

	rbb, err := json.Marshal(rb)
	if err != nil {
		return nil, err
	}
	return rbb, nil
}

func (m *HTTPMock) getUtilityEmployment() ([]byte, error) {
	if m.dsbPtnrErr {
		return m.commonNonexError()
	}

	rb := &GetUtilityEmploymentResponse{
		Code:    200,
		Message: "Success - All EMPLOYMENT/s",
		Result: []string{
			"Administrative/Human Resources",
			"Agriculture",
			"Banking /Financial Services",
			"Computer and Information Tech Services",
			"Construction/Contractors",
		},
	}

	rbb, err := json.Marshal(rb)
	if err != nil {
		return nil, err
	}
	return rbb, nil
}

func (m *HTTPMock) getUtilitySourceFund() ([]byte, error) {
	if m.dsbPtnrErr {
		return m.commonNonexError()
	}

	rb := &GetUtilitySourceFundResponse{
		Code:    200,
		Message: "Success - All SOURCEFUND/s",
		Result: []string{
			"Salary",
			"Savings",
			"Borrowed Funds/Loan",
			"Pension/Government/Welfare",
			"Gift",
		},
	}

	rbb, err := json.Marshal(rb)
	if err != nil {
		return nil, err
	}
	return rbb, nil
}

func (m *HTTPMock) retryTransaction() ([]byte, error) {
	if m.dsbPtnrErr {
		return m.commonNonexError()
	}
	rb := &PerahubRemitRetryResponseBody{
		Code:    "200",
		Message: "Good",
		Result: PerahubRemitRetryResult{
			ID:                 7461,
			LocationID:         0,
			UserID:             5188,
			TrxDate:            "2022-06-15",
			CurrencyID:         1,
			RemcoID:            25,
			TrxType:            2,
			IsDomestic:         1,
			CustomerID:         6925902,
			CustomerName:       "Soto, Blanche, G",
			ControlNumber:      "PH1655176065",
			SenderName:         "Sauer, Mittie, O",
			ReceiverName:       "Soto, Blanche, G",
			PrincipalAmount:    "179.00",
			ServiceCharge:      "50.00",
			DstAmount:          "1.00",
			TotalAmount:        "229.00",
			McRate:             "1.00",
			BuyBackAmount:      "1.00",
			RateCategory:       "OPTIONAL",
			McRateID:           1,
			OriginatingCountry: "PH",
			DestinationCountry: "PH",
			PurposeTransaction: "REQUIRED",
			SourceFund:         "REQUIRED",
			Occupation:         "EMPLOYED",
			RelationTo:         "REQUIRED",
			BirthPlace:         "MANILA",
			BirthCountry:       "PH",
			IDType:             "PASSPORT",
			IDNumber:           "0001292",
			Address:            "ADDRESS",
			Barangay:           "1",
			City:               "MANILA",
			Province:           "METRO",
			Country:            "PH",
			ContactNumber:      "09999999999",
			RiskScore:          0,
			RiskCriteria:       "",
			ClientReferenceNo:  "U306WYJXYCPF18",
			FormType:           "OPTIONAL",
			FormNumber:         "OPTIONAL",
			PayoutType:         1,
			RemoteLocationID:   1,
			RemoteUserID:       1,
			RemoteIPAddress:    "OPTIONAL",
			IPAddress:          "OPTIONAL",
			CreatedAt:          "2022-06-15T02:41:40.000000Z",
			UpdatedAt:          "2022-06-15T02:43:09.000000Z",
			ReferenceNumber:    "796e98e940b5982ffa5cee4f89816ad8",
			ZipCode:            "1013",
			Status:             1,
			CurrentAddress: PerahubRemitAddress{
				City:     "MANILA",
				Country:  "PH",
				Barangay: "Barangay 1",
				Province: "METRO MANILA",
				ZipCode:  "1013",
				Address1: "Marcos Highway",
			},
			PermanentAddress: PerahubRemitAddress{
				City:     "MANILA",
				Country:  "PH",
				Barangay: "Barangay 1",
				Province: "METRO MANILA",
				ZipCode:  "1013",
				Address1: "Marcos Highway",
			},
			APIRequest: PerahubRemitAPIRequest{
				City:               "MANILA",
				Address:            "ADDRESS",
				Country:            "PH",
				IDType:             "PASSPORT",
				McRate:             "1",
				UserID:             "5188",
				Barangay:           "Barangay 1",
				Province:           "METRO MANILA",
				RemcoID:            "25",
				TrxDate:            "2022-06-15",
				TrxType:            "2",
				ZipCode:            "1013",
				FormType:           "OPTIONAL",
				IDNumber:           "0001292",
				DstAmount:          "1",
				IPAddress:          "OPTIONAL",
				McRateID:           "1",
				Occupation:         "EMPLOYED",
				BirthPlace:         "MANILA",
				CurrencyID:         "1",
				CustomerID:         "6925902",
				FormNumber:         "OPTIONAL",
				IsDomestic:         "1",
				LocationID:         0,
				PayoutType:         "1",
				RelationTo:         "REQUIRED",
				SenderName:         "Sauer, Mittie, O",
				SourceFund:         "REQUIRED",
				PartnerCode:        "DRP",
				TotalAmount:        "229",
				BirthCountry:       "PH",
				CustomerName:       "Soto, Blanche, G",
				RateCategory:       "OPTIONAL",
				ReceiverName:       "Soto, Blanche, G",
				ContactNumber:      "09999999999",
				ControlNumber:      "PH1655176065",
				RemoteUserID:       "1",
				ServiceCharge:      "50",
				BuyBackAmount:      "1",
				PrincipalAmount:    179,
				ReferenceNumber:    "796e98e940b5982ffa5cee4f89816ad8",
				SenderLastName:     "Sauer",
				RemoteIPAddress:    "OPTIONAL",
				SenderFirstName:    "Mittie",
				ReceiverLastName:   "Soto",
				RemoteLocationID:   "1",
				SenderMiddleName:   "O",
				ClientReferenceNo:  "U306WYJXYCPF18",
				DestinationCountry: "PH",
				OriginatingCountry: "PH",
				PurposeTransaction: "REQUIRED",
				ReceiverFirstName:  "Blanche",
				ReceiverMiddleName: "G",
				APIRequestCurrentAddress: PerahubRemitAddress{
					City:     "MANILA",
					Country:  "PH",
					Barangay: "Barangay 1",
					Province: "METRO MANILA",
					ZipCode:  "1013",
					Address1: "Marcos Highway",
				},
				APIRequestPermanentAddress: PerahubRemitAddress{
					City:     "MANILA",
					Country:  "PH",
					Barangay: "Barangay 1",
					Province: "METRO MANILA",
					ZipCode:  "1013",
					Address1: "Marcos Highway",
				},
			},
		},
	}

	rbb, err := json.Marshal(rb)
	if err != nil {
		return nil, err
	}
	return rbb, nil
}

func (m *HTTPMock) confirmTransaction() ([]byte, error) {
	rbe := &nonexError{
		Code: "1",
		Msg:  "some error",
		Error: ErrorType{
			Type: IRErr,
			Msg:  "Transaction does not exists",
		},
	}
	rbbe, err := json.Marshal(rbe)
	if err != nil {
		return nil, err
	}
	if m.cfCrPtnrErr {
		return rbbe, cfCrPtnrErr
	}
	if m.cfDsbPtnrErr {
		return rbbe, cfDsbPtnrErr
	}

	rb := &PerahubRemitConfirmResponseBody{
		Code:    "200",
		Message: "Good",
		Result: PerahubRemitResult{
			ID:                 7461,
			LocationID:         0,
			UserID:             5188,
			TrxDate:            "2022-06-15",
			CurrencyID:         1,
			RemcoID:            25,
			TrxType:            2,
			IsDomestic:         1,
			CustomerID:         6925902,
			CustomerName:       "Soto, Blanche, G",
			ControlNumber:      "PH1655176065",
			SenderName:         "Sauer, Mittie, O",
			ReceiverName:       "Soto, Blanche, G",
			PrincipalAmount:    "179.00",
			ServiceCharge:      "50.00",
			DstAmount:          "1.00",
			TotalAmount:        "229.00",
			McRate:             "1.00",
			BuyBackAmount:      "1.00",
			RateCategory:       "OPTIONAL",
			McRateID:           1,
			OriginatingCountry: "PH",
			DestinationCountry: "PH",
			PurposeTransaction: "REQUIRED",
			SourceFund:         "REQUIRED",
			Occupation:         "EMPLOYED",
			RelationTo:         "REQUIRED",
			BirthPlace:         "MANILA",
			BirthCountry:       "PH",
			IDType:             "PASSPORT",
			IDNumber:           "0001292",
			Address:            "ADDRESS",
			Barangay:           "1",
			City:               "MANILA",
			Province:           "METRO",
			Country:            "PH",
			ContactNumber:      "09999999999",
			RiskScore:          0,
			RiskCriteria:       "",
			ClientReferenceNo:  "U306WYJXYCPF18",
			FormType:           "OPTIONAL",
			FormNumber:         "OPTIONAL",
			PayoutType:         1,
			RemoteLocationID:   1,
			RemoteUserID:       1,
			RemoteIPAddress:    "OPTIONAL",
			IPAddress:          "OPTIONAL",
			CreatedAt:          "2022-06-15T02:41:40.000000Z",
			UpdatedAt:          "2022-06-15T02:43:09.000000Z",
			ReferenceNumber:    "796e98e940b5982ffa5cee4f89816ad8",
			ZipCode:            "1013",
			Status:             1,
			CurrentAddress: PerahubRemitConfirmAddress{
				City:     "MANILA",
				Country:  "PH",
				Barangay: "Barangay 1",
				Province: "METRO MANILA",
				ZipCode:  "1013",
				Address1: "Marcos Highway",
			},
			PermanentAddress: PerahubRemitConfirmAddress{
				City:     "MANILA",
				Country:  "PH",
				Barangay: "Barangay 1",
				Province: "METRO MANILA",
				ZipCode:  "1013",
				Address1: "Marcos Highway",
			},
			APIRequest: PerahubRemitConfirmAPIRequest{
				City:               "MANILA",
				Address:            "ADDRESS",
				Country:            "PH",
				IDType:             "PASSPORT",
				McRate:             "1",
				UserID:             "5188",
				Barangay:           "Barangay 1",
				Province:           "METRO MANILA",
				RemcoID:            "25",
				TrxDate:            "2022-06-15",
				TrxType:            "2",
				ZipCode:            "1013",
				FormType:           "OPTIONAL",
				IDNumber:           "0001292",
				DstAmount:          "1",
				IPAddress:          "OPTIONAL",
				McRateID:           "1",
				Occupation:         "EMPLOYED",
				BirthPlace:         "MANILA",
				CurrencyID:         "1",
				CustomerID:         "6925902",
				FormNumber:         "OPTIONAL",
				IsDomestic:         "1",
				LocationID:         0,
				PayoutType:         "1",
				RelationTo:         "REQUIRED",
				SenderName:         "Sauer, Mittie, O",
				SourceFund:         "REQUIRED",
				PartnerCode:        "DRP",
				TotalAmount:        "229",
				BirthCountry:       "PH",
				CustomerName:       "Soto, Blanche, G",
				RateCategory:       "OPTIONAL",
				ReceiverName:       "Soto, Blanche, G",
				ContactNumber:      "09999999999",
				ControlNumber:      "PH1655176065",
				RemoteUserID:       "1",
				ServiceCharge:      "50",
				BuyBackAmount:      "1",
				PrincipalAmount:    179,
				ReferenceNumber:    "796e98e940b5982ffa5cee4f89816ad8",
				SenderLastName:     "Sauer",
				RemoteIPAddress:    "OPTIONAL",
				SenderFirstName:    "Mittie",
				ReceiverLastName:   "Soto",
				RemoteLocationID:   "1",
				SenderMiddleName:   "O",
				ClientReferenceNo:  "U306WYJXYCPF18",
				DestinationCountry: "PH",
				OriginatingCountry: "PH",
				PurposeTransaction: "REQUIRED",
				ReceiverFirstName:  "Blanche",
				ReceiverMiddleName: "G",
				APIRequestCurrentAddress: PerahubRemitConfirmAddress{
					City:     "MANILA",
					Country:  "PH",
					Barangay: "Barangay 1",
					Province: "METRO MANILA",
					ZipCode:  "1013",
					Address1: "Marcos Highway",
				},
				APIRequestPermanentAddress: PerahubRemitConfirmAddress{
					City:     "MANILA",
					Country:  "PH",
					Barangay: "Barangay 1",
					Province: "METRO MANILA",
					ZipCode:  "1013",
					Address1: "Marcos Highway",
				},
			},
		},
	}

	rbb, err := json.Marshal(rb)
	if err != nil {
		return nil, err
	}
	return rbb, nil
}

func (m *HTTPMock) inquireTransaction() ([]byte, error) {
	if m.crPtnrErr {
		rb := &nonexError{
			Code: "1",
			Msg:  "some error",
			Error: ErrorType{
				Type: IRErr,
				Msg:  "Transaction does not exists",
			},
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, crPtnrErr
	}
	rb := &PerahubRemitInquireResponse{
		Code:    200,
		Message: "PeraHUB Reference Number (PHRN) is available for Payout",
		Result: PerahubRemitInquireResult{
			PrincipalAmount:    179,
			IsoCurrency:        "PHP",
			ConversionRate:     1,
			SenderLastName:     "Sauer",
			SenderFirstName:    "Mittie",
			SenderMiddleName:   "O",
			ReceiverLastName:   "Soto",
			ReceiverFirstName:  "Blanche",
			ReceiverMiddleName: "G",
			ControlNumber:      "PH1655176065",
			OriginatingCountry: "PH",
			DestinationCountry: "PH",
			SenderName:         "Sauer, Mittie, O",
			ReceiverName:       "Soto, Blanche, G",
			PartnerCode:        "DRP",
		},
	}

	rbb, err := json.Marshal(rb)
	if err != nil {
		return nil, err
	}
	return rbb, nil
}

func (m *HTTPMock) validateTransaction() ([]byte, error) {
	if m.dsbPtnrErr {
		rb := &nonexError{
			Code: "1",
			Msg:  "some error",
			Error: ErrorType{
				Type: IRErr,
				Msg:  "Transaction does not exists",
			},
		}
		rbb, err := json.Marshal(rb)
		if err != nil {
			return nil, err
		}
		return rbb, dsbPtnrErr
	}
	vr := &PerahubRemitValidateResponseBody{
		Code:    "200",
		Message: "Good",
		Result: PerahubRemitValidateResult{
			LocationID:         0,
			UserID:             "5188",
			TrxDate:            "2022-06-15",
			CurrencyID:         "1",
			RemcoID:            "25",
			TrxType:            "2",
			IsDomestic:         "1",
			CustomerID:         "6925902",
			ControlNumber:      "PH1655176065",
			ClientReferenceNo:  "U306WYJXYCPF18",
			CustomerName:       "Soto, Blanche, G",
			SenderName:         "Sauer, Mittie, O",
			ReceiverName:       "Soto, Blanche, G",
			PrincipalAmount:    179,
			ServiceCharge:      "50",
			DstAmount:          "1",
			TotalAmount:        "229",
			McRate:             "1",
			BuyBackAmount:      "1",
			RateCategory:       "OPTIONAL",
			McRateID:           "1",
			OriginatingCountry: "PH",
			DestinationCountry: "PH",
			PurposeTransaction: "REQUIRED",
			SourceFund:         "REQUIRED",
			Occupation:         "EMPLOYED",
			RelationTo:         "REQUIRED",
			BirthPlace:         "MANILA",
			BirthCountry:       "PH",
			IDType:             "PASSPORT",
			IDNumber:           "0001292",
			Address:            "ADDRESS",
			Barangay:           "Barangay 1",
			City:               "MANILA",
			Province:           "METRO MANILA",
			ZipCode:            "1013",
			Country:            "PH",
			ContactNumber:      "09999999999",
			FormType:           "OPTIONAL",
			FormNumber:         "OPTIONAL",
			PayoutType:         "1",
			RemoteLocationID:   "1",
			RemoteUserID:       "1",
			RemoteIPAddress:    "OPTIONAL",
			IPAddress:          "OPTIONAL",
			ReferenceNumber:    "796e98e940b5982ffa5cee4f89816ad8",
			CurrentAddress: PerahubRemitAddress{
				Address1: "Marcos Highway",
				Barangay: "Barangay 1",
				City:     "MANILA",
				Province: "METRO MANILA",
				ZipCode:  "1013",
				Country:  "PH",
			},
			PermanentAddress: PerahubRemitAddress{
				Address1: "Marcos Highway",
				Barangay: "Barangay 1",
				City:     "MANILA",
				Province: "METRO MANILA",
				ZipCode:  "1013",
				Country:  "PH",
			},
			APIRequest: PerahubRemitValidateAPIRequest{
				LocationID:         0,
				UserID:             "5188",
				TrxDate:            "2022-06-15",
				CurrencyID:         "1",
				RemcoID:            "25",
				TrxType:            "2",
				IsDomestic:         "1",
				CustomerID:         "6925902",
				ControlNumber:      "PH1655176065",
				ClientReferenceNo:  "U306WYJXYCPF18",
				CustomerName:       "Soto, Blanche, G",
				SenderName:         "Sauer, Mittie, O",
				ReceiverName:       "Soto, Blanche, G",
				PrincipalAmount:    179,
				ServiceCharge:      "50",
				DstAmount:          "1",
				TotalAmount:        "229",
				McRate:             "1",
				BuyBackAmount:      "1",
				RateCategory:       "OPTIONAL",
				McRateID:           "1",
				OriginatingCountry: "PH",
				DestinationCountry: "PH",
				PurposeTransaction: "REQUIRED",
				SourceFund:         "REQUIRED",
				Occupation:         "EMPLOYED",
				RelationTo:         "REQUIRED",
				BirthPlace:         "MANILA",
				BirthCountry:       "PH",
				IDType:             "PASSPORT",
				IDNumber:           "0001292",
				Address:            "ADDRESS",
				Barangay:           "Barangay 1",
				City:               "MANILA",
				Province:           "METRO MANILA",
				ZipCode:            "1013",
				Country:            "PH",
				ContactNumber:      "09999999999",
				FormType:           "OPTIONAL",
				FormNumber:         "OPTIONAL",
				PayoutType:         "1",
				RemoteLocationID:   "1",
				RemoteUserID:       "1",
				RemoteIPAddress:    "OPTIONAL",
				IPAddress:          "OPTIONAL",
				ReferenceNumber:    "796e98e940b5982ffa5cee4f89816ad8",
				CurrentAddress: PerahubRemitAddress{
					Address1: "Marcos Highway",
					Barangay: "Barangay 1",
					City:     "MANILA",
					Province: "METRO MANILA",
					ZipCode:  "1013",
					Country:  "PH",
				},
				PermanentAddress: PerahubRemitAddress{
					Address1: "Marcos Highway",
					Barangay: "Barangay 1",
					City:     "MANILA",
					Province: "METRO MANILA",
					ZipCode:  "1013",
					Country:  "PH",
				},
				SenderLastName:     "Sauer",
				SenderFirstName:    "Mittie",
				SenderMiddleName:   "O",
				ReceiverLastName:   "Soto",
				ReceiverFirstName:  "Blanche",
				ReceiverMiddleName: "G",
				PartnerCode:        "DRP",
			},
			UpdatedAt: time.Time{},
			CreatedAt: time.Time{},
			ID:        7461,
		},
	}

	pvr, err := json.Marshal(vr)
	if err != nil {
		return nil, err
	}
	return pvr, nil
}

func (s *HTTPMock) Login(ctx context.Context, lr LoginRequest) (*Customer, error) {
	return &Customer{
		FrgnRefNo:      "frgnRefNo",
		CustomerCode:   "customerCode",
		LastName:       "Doe",
		FirstName:      "John",
		Birthdate:      "2020/12/12",
		Nationality:    "phillipine",
		PresentAddress: "address",
		Occupation:     "occupation",
		EmployerName:   "empName",
		CustomerIDNo:   "custIdNo",
		Wucardno:       "wuCardNo",
		Debitcardno:    "debitCardNo",
		Loyaltycardno:  "Loyaltycardno",
	}, nil
}

func isNonexRequest(r *PerahubRequest) bool {
	if r.WU.Body.Request == "" {
		return true
	}
	return false
}

var idJSON = []byte(`[{"index":"1","template_value":"A","document_type":"Passport","document_desc_eng":"Local and Foreign Passports","document_desc_fil":"Pasaporte na ibinibigay ng Pamahalan ng Pilipinas o ng ibang bansa.","hasIssueDate":"Y","hasExpiration":"Y"},{"index":"2","template_value":"B","document_type":"Gov't Office\/GOCC ID","document_desc_eng":"ID issued by Government and its agencies (digitized hard plastic type) including Firearms License","document_desc_fil":"ID na ibinibigay ng Pamahalaan ng Pilipinas o ng mga kagawaran\/tanggapan nito katulad ng Firearms License mula sa Pambansang Pulisya ng Pilipinas (PNP).  Ito ay dapat digitize at matigas na plastic type at hindi papel lamang o naka-laminate na papel.","hasIssueDate":"N","hasExpiration":"N"},{"index":"3","template_value":"C","document_type":"Driver's license","document_desc_eng":"Philippine Driver's License","document_desc_fil":"Lisensiya sa pagmamaneho na ibinibigay ng Opisina ng Transportasyong Panlupa or LTO (Land Transportation Office)","hasIssueDate":"Y","hasExpiration":"Y"},{"index":"4","template_value":"D","document_type":"Postal ID","document_desc_eng":"Philippine Postal Corporation (PPC) Issued Identity Card","document_desc_fil":"ID na ibinibigay ng Opisina ng Koreo sa Pilipinas o Philippine Postal Corporation.","hasIssueDate":"N","hasExpiration":"Y"},{"index":"5","template_value":"E","document_type":"OWWA ID (Overseas Workers)","document_desc_eng":"Overseas Workers Welfare Admin (OWWA) ID","document_desc_fil":"ID na ibinibigay ng Overseas Workers Welfare Admin (OWWA)","hasIssueDate":"Y","hasExpiration":"Y"},{"index":"8","template_value":"H","document_type":"Seaman's Book","document_desc_eng":"Seaman's Passport","document_desc_fil":"Pasaporte ng mga Marino o Seaman","hasIssueDate":"Y","hasExpiration":"Y"},{"index":"9","template_value":"I","document_type":"NBI Clearance","document_desc_eng":"National Bureau of Investigation (NBI) Clearance","document_desc_fil":"Clearance na ibinibigay ng National Bureau of Investigation (NBI) o Kagawaran ng Pagsisiyasat","hasIssueDate":"Y","hasExpiration":"Y"},{"index":"12","template_value":"M","document_type":"Voter's ID","document_desc_eng":"Commission on Election (Comelec) Issued Voter's ID","document_desc_fil":"ID na ibinibigay ng Local ng Pamahalaan para sa mga Senior Citizens","hasIssueDate":"N","hasExpiration":"N"},{"index":"13","template_value":"N","document_type":"GSIS e-Card","document_desc_eng":"Government Service Insurance System (GSIS) e-Card","document_desc_fil":"e-Card na ibinibigay ng GSIS o Government Service Insurance System","hasIssueDate":"N","hasExpiration":"Y"},{"index":"17","template_value":"R","document_type":"Disabled Person Cert\/ID NCWDP","document_desc_eng":"Certification issued by National Council for the Welfare of Disabled Persons (NCWDP)","document_desc_fil":"Certificate para sa mga may kapansanan na ibinibigay ng National Council for the Welfare of Disabled Persons (NCWDP)","hasIssueDate":"N","hasExpiration":"N"},{"index":"20","template_value":"U","document_type":"DSWD Certificate","document_desc_eng":"Travel Clearance issued by Department of Social Welfare and Development (DSWD) for minors traveling without parents","document_desc_fil":"Certificate na ibinibigay ng DSWD o Kagawaran ng Kagalingang Panlipunan at Kaunlaran para pahintulutan ang pagbiyahe ng mga menor de edad na hindi kasama ang mga magulang","hasIssueDate":"N","hasExpiration":"N"},{"index":"21","template_value":"V","document_type":"PRC ID (Prof. Reg. Commission)","document_desc_eng":"Regulated Professionals' Identification Card for Doctors, Engineers and other professions","document_desc_fil":"ID na ibinibigay ng PRC (Philippine Regulations Commission) o Komisyon na Namamahala sa mga Propesyonal","hasIssueDate":"Y","hasExpiration":"Y"},{"index":"22","template_value":"W","document_type":"Senior Citizen Card","document_desc_eng":"Senior Citizen's Card issued by Local Government Unit","document_desc_fil":"ID na ibinibigay ng Local ng Pamahalaan o Munisipyo para sa mga Senior Citizens","hasIssueDate":"Y","hasExpiration":"N"},{"index":"23","template_value":"X","document_type":"Student's ID","document_desc_eng":"School Issued ID","document_desc_fil":"e-Card na ibinibigay ng GSIS o Government Service Insurance System","hasIssueDate":"N","hasExpiration":"N"},{"index":"24","template_value":"1","document_type":"Integrated Bar (IBP) ID","document_desc_eng":"IBP is a Membership for Licensed Lawyers","document_desc_fil":"ID na ibinibigay ng Integrated Bar of the Philippines, ang opisyal na Samahan ng mga Abogado sa Pilipinas","hasIssueDate":"N","hasExpiration":"N"},{"index":"25","template_value":"2","document_type":"Police Clearance","document_desc_eng":"Police Clearance Certificate or Clearance Identification Card","document_desc_fil":"Clearance Certificate o pagpapatunay ng kawalan ng criminal record na ibinibigay ng Pambansang Pulisya ng Pilipinas o PNP","hasIssueDate":"Y","hasExpiration":"Y"},{"index":"26","template_value":"3","document_type":"SSS Card (Social Security)","document_desc_eng":"Social Security System (SSS) Card","document_desc_fil":"ID Card na ibinibigay ng Social Security System (SSS)","hasIssueDate":"N","hasExpiration":"N"},{"index":"27","template_value":"4","document_type":"OFW ID","document_desc_eng":"Overseas Employment Certificate (OEC) issued by Philippine Overseas & Employment Agency (POEA)","document_desc_fil":"Certificate para sa mga manggagawa sa ibang bansa, kilala sa tawag na Overseas Employment Certificate (OEC) na ibinibigay ng Philippine Overseas & Employment Agency (POEA)","hasIssueDate":"Y","hasExpiration":"Y"},{"index":"28","template_value":"5","document_type":"Alien Cert. of Registration","document_desc_eng":"Alien Certificate of Registration (ACR) i-card issued by Bureau of Immigration","document_desc_fil":"Alien Certificate of Registration (ACR) i-card na ibinibigay ng Bureau of Immigration (BOI) o Kagawaran ng Imigrasyon","hasIssueDate":"Y","hasExpiration":"Y"},{"index":"29","template_value":"6","document_type":"Employment ID","document_desc_eng":"Company Issued IDs, Employment IDs (Government and Private)","document_desc_fil":"ID na ibinibigay ng kumpanyang pinapasukan (pampubliko o pribado)","hasIssueDate":"N","hasExpiration":"N"},{"index":"30","template_value":"7","document_type":"Taxpayer's ID (TIN ID)","document_desc_eng":"Taxpayer's ID issued by Bureau of Internal Revenue (BIR)","document_desc_fil":"Tax Payer's ID na ibinibigay ng Kawanihan ng Rentas Internas o BIR (Bureau of Internal Revenue)","hasIssueDate":"N","hasExpiration":"N"},{"index":"31","template_value":"8","document_type":"Barangay Certification","document_desc_eng":"Certification or ID issued by Barangay","document_desc_fil":"Certificate o ID na ibinibigay ng Pamahalaang Barangay para sa pagpapatunay kung saan nakatira and indibidwal.","hasIssueDate":"Y","hasExpiration":"Y"},{"index":"32","template_value":"9","document_type":"Unified MultiPurpose ID (UMID)","document_desc_eng":"Unified Multi-Purpose ID (UMID)","document_desc_fil":"Unified Multi-Purpose ID (UMID) ay ID na maaaring ibinibigay ng SSS o Pag-Ibig\/HDMF","hasIssueDate":"N","hasExpiration":"N"}]`)

var ocJSON = []byte(`[{"occupation":"Airline\/Maritime Employee","occupation_value":"Airline\/Maritime Employee"},{"occupation":"Art\/Entertainment\/Media\/Sports Professional","occupation_value":"Art\/Entertainment\/Media\/Sports"},{"occupation":"Civil\/Government Employee","occupation_value":"Civil\/Government Employee"},{"occupation":"Domestic Helper","occupation_value":"Domestic Helper"},{"occupation":"Driver","occupation_value":"Driver"},{"occupation":"Teacher\/Educator","occupation_value":"Teacher\/Educator"},{"occupation":"Hotel\/Restaurant\/Leisure Services Employee","occupation_value":"Hotel\/Restaurant\/Leisure"},{"occupation":"Housewife\/Child Care","occupation_value":"Housewife\/Child Care"},{"occupation":"IT and Tech Professional","occupation_value":"IT and Tech Professional"},{"occupation":"Laborer-Agriculture","occupation_value":"Laborer-Agriculture"},{"occupation":"Laborer-Construction","occupation_value":"Laborer-Construction"},{"occupation":"Laborer-Manufacturing","occupation_value":"Laborer-Manufacturing"},{"occupation":"Laborer- Oil\/Gas\/Mining\/Forestry","occupation_value":"Laborer- Oil\/Gas\/Mining"},{"occupation":"Medical and Health Care Professional","occupation_value":"Medical\/Health Care"},{"occupation":"Non-profit, Volunteer","occupation_value":"Non-profit\/Volunteer"},{"occupation":"Cosmetic\/Personal Care Services","occupation_value":"Cosmetic\/Personal Care"},{"occupation":"Law Enforcement\/Military Professional","occupation_value":"Law Enforcement\/Military"},{"occupation":"Office Professional","occupation_value":"Office Professional"},{"occupation":"Professional Service Practitioner","occupation_value":"Prof Svc Practitioner"},{"occupation":"Religious\/Church Servant","occupation_value":"Religious\/Church Servant"},{"occupation":"Retail Sales","occupation_value":"Retail Sales"},{"occupation":"Retired","occupation_value":"Retired"},{"occupation":"Sales\/Insurance\/Real Estate Professional","occupation_value":"Sales\/Insurance\/Real Estate"},{"occupation":"Science\/Research Professional","occupation_value":"Science\/Research Professional"},{"occupation":"Security Guard","occupation_value":"Security Guard"},{"occupation":"Self-Employed","occupation_value":"Self-Employed"},{"occupation":"Skilled Trade\/Specialist","occupation_value":"Skilled Trade\/Specialist"},{"occupation":"Student","occupation_value":"Student"},{"occupation":"Unemployed","occupation_value":"Unemployed"}]`)

var poJSON = []byte(`[{"position":"Entry Level","position_value":"Entry Level"},{"position":"Mid-Level\/Supervisory\/Management","position_value":"Mid-Level\/Supervisory\/Management"},{"position":"Senior Level\/Executive","position_value":"Senior Level\/Executive"},{"position":"Owner","position_value":"Owner"}]`)

var puJSON = []byte(`[{"purpose":"Family Support\/Living Expenses","purpose_value":"Family Support\/Living Expenses"},{"purpose":"Saving\/Investments","purpose_value":"Saving\/Investments"},{"purpose":"Gift","purpose_value":"Gift"},{"purpose":"Goods & Services payment","purpose_value":"Goods & Services payment"},{"purpose":"Travel expenses","purpose_value":"Travel expenses"},{"purpose":"Education\/School Fee","purpose_value":"Education\/School Fee"},{"purpose":"Rent\/Mortgage","purpose_value":"Rent\/Mortgage"},{"purpose":"Emergency\/Medical Aid","purpose_value":"Emergency\/Medical Aid"},{"purpose":"Charity\/Aid Payment","purpose_value":"Charity\/Aid Payment"},{"purpose":"Employee Payroll\/Employee Expense","purpose_value":"Employee Payroll\/Employee Expense"},{"purpose":"Prize or Lottery Fees\/Taxes","purpose_value":"Prize or Lottery Fees\/Taxes"}]`)

var reJSON = []byte(`[{"relationship":"Family","relationship_value":"Family"},{"relationship":"Friend","relationship_value":"Friend"},{"relationship":"Trade\/Business Partner","relationship_value":"Trade\/BusinesPartner"},{"relationship":"Employee\/Employer","relationship_value":"Employee\/Employer"},{"relationship":"Donor\/Receiver of Charitable Funds","relationship_value":"Donor\/Receiver of Ch"},{"relationship":"Purchaser\/Seller","relationship_value":"Purchaser\/Seller"}]`)

var sfJSON = []byte(`[{"source_of_funds":"Salary","source_of_funds_value":"Salary"},{"source_of_funds":"Savings","source_of_funds_value":"Savings"},{"source_of_funds":"Borrowed Funds\/Loan","source_of_funds_value":"Borrowed Funds\/Loan"},{"source_of_funds":"Gift","source_of_funds_value":"Gift"},{"source_of_funds":"Pension\/Government\/Welfare","source_of_funds_value":"Pension\/Government\/Welfare"},{"source_of_funds":"Inheritance","source_of_funds_value":"Inheritance"},{"source_of_funds":"Charitable Donations","source_of_funds_value":"Charitable Donations"},{"source_of_funds":"Cash Tips","source_of_funds_value":"Cash Tips"},{"source_of_funds":"Sale of Goods\/Property\/Services","source_of_funds_value":"Sale of Goods\/Property\/Services"},{"source_of_funds":"Investment Income","source_of_funds_value":"Investment Income"}]`)

var ccJSON = []byte(`[{"COUNTRY_LONG":"Afghanistan","ISO_COUNTRY_NUM_CD":"004","ISO_COUNTRY_CD":"AF","CURRENCY_CD":"AFN","ISO_CURRENCY_NUM_CD":"971","CURRENCY_NAME":"Afghanistan Afghani"},{"COUNTRY_LONG":"Afghanistan","ISO_COUNTRY_NUM_CD":"004","ISO_COUNTRY_CD":"AF","CURRENCY_CD":"USD","ISO_CURRENCY_NUM_CD":"840","CURRENCY_NAME":"US Dollar"},{"COUNTRY_LONG":"Afghanistan US Military Base","ISO_COUNTRY_NUM_CD":"840","ISO_COUNTRY_CD":"XP","CURRENCY_CD":"USD","ISO_CURRENCY_NUM_CD":"840","CURRENCY_NAME":"US Dollar"},{"COUNTRY_LONG":"Albania","ISO_COUNTRY_NUM_CD":"008","ISO_COUNTRY_CD":"AL","CURRENCY_CD":"ALL","ISO_CURRENCY_NUM_CD":"008","CURRENCY_NAME":"Albanian Lek"},{"COUNTRY_LONG":"Albania","ISO_COUNTRY_NUM_CD":"008","ISO_COUNTRY_CD":"AL","CURRENCY_CD":"EUR","ISO_CURRENCY_NUM_CD":"978","CURRENCY_NAME":"Euro"},{"COUNTRY_LONG":"Algeria","ISO_COUNTRY_NUM_CD":"012","ISO_COUNTRY_CD":"DZ","CURRENCY_CD":"DZD","ISO_CURRENCY_NUM_CD":"012","CURRENCY_NAME":"Algerian Dinar"},{"COUNTRY_LONG":"American Samoa","ISO_COUNTRY_NUM_CD":"016","ISO_COUNTRY_CD":"AS","CURRENCY_CD":"USD","ISO_CURRENCY_NUM_CD":"840","CURRENCY_NAME":"US Dollar"},{"COUNTRY_LONG":"Angola","ISO_COUNTRY_NUM_CD":"024","ISO_COUNTRY_CD":"AO","CURRENCY_CD":"AOA","ISO_CURRENCY_NUM_CD":"973","CURRENCY_NAME":"Angolan  Kwanza"}]`)

func ErrStruct(errStruct PtnrErrStruct, errType string) PtnrErr {
	e := PtnrErrStructs[errStruct]
	et := e.Error.(*nonexError)
	a := et.Error.(ErrorType)
	a.Type = RIAErr
	e.Error = et
	return e
}

type PtnrErrStruct string

// A list of different partner error structures, this is added so we can manage handling different types of
// error structures in a more organized way
// If finding a new error structure add it here and create a test for it
const (
	E1 PtnrErrStruct = "E1"
	E2 PtnrErrStruct = "E2"
	E3 PtnrErrStruct = "E3"
)

type PtnrErr struct {
	ErrStruct PtnrErrStruct
	Error     interface{}
}

var PtnrErrStructs = map[PtnrErrStruct]PtnrErr{
	E1: {
		ErrStruct: E1,
		Error: &nonexError{
			Code: "2001",
			Msg:  "M1",
			Error: ErrorType{
				Msg: "M1",
			},
		},
	},
	E2: {
		ErrStruct: E2,
		Error: &nonexError{
			Msg: "The given data was invalid.",
			Errors: map[string][]string{
				"agent_code": {
					"The agent code field is required.",
				},
			},
		},
	},
	E3: {
		ErrStruct: E3,
		Error: &nonexError{
			Code:  "2003",
			Msg:   "M3",
			Error: "E3",
		},
	},
}

func IsThere(val interface{}, array interface{}) (exists bool) {
	exists = false
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				exists = true
				return
			}
		}
	}
	return
}

func ErrStructBillPay(errStruct BillPayErrStruct, errType string) BillPayErr {
	e := BillPayErrStructs[errStruct]
	et := e.Error.(*billerError)
	a := et.Error.(ErrorType)
	a.Type = BillPayError
	e.Error = et
	return e
}

type BillPayErrStruct string

// A list of different partner error structures, this is added so we can manage handling different types of
// error structures in a more organized way
// If finding a new error structure add it here and create a test for it
const (
	BP1 BillPayErrStruct = "BP1"
	BP2 BillPayErrStruct = "BP2"
	BP3 BillPayErrStruct = "BP3"
)

type BillPayErr struct {
	ErrStruct BillPayErrStruct
	Error     interface{}
}

var BillPayErrStructs = map[BillPayErrStruct]BillPayErr{
	BP1: {
		ErrStruct: BP1,
		Error: &billerError{
			Code: "2001",
			Msg:  "M1",
			Error: ErrorType{
				Msg: "M1",
			},
		},
	},
	BP2: {
		ErrStruct: BP2,
		Error: &billerError{
			Msg: "validation_error",
			Errors: map[string][]string{
				"code": {
					"The code is required",
				},
			},
		},
	},
	BP3: {
		ErrStruct: BP3,
		Error: &billerError{
			Code:  "2003",
			Msg:   "M3",
			Error: "BP3",
		},
	},
}

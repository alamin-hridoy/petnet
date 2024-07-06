package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func IsMicroInsuranceRequest(path string) bool {
	return strings.Contains(path, microinsurancePath)
}

func MicroInsuranceRequest(req *http.Request, isErr bool) (*http.Response, error) {
	var res []byte
	var err error

	if isErr {
		res = []byte(`{
			"code": "07",
			"message": "Failed",
			"remco_id": null,
			"error": {
			  "type": "RuralNet Error",
			  "message": " You have already exceeded the maximum units allowed in purchasing this product!",
			  "line": null
			}
		}`)

		return &http.Response{
			StatusCode: 422,
			Body:       ioutil.NopCloser(bytes.NewReader(res)),
		}, nil
	}

	act := dynamicUrlModify(req, strings.ReplaceAll(req.URL.Path, microinsurancePath, ""))
	status := 200
	switch act {
	case "POST_transact", "POST_retry":
		res, err = mockTransactResult(req)
	case "POST_get-reprint":
		res, err = mockGetReprint()
	case "POST_get-transactions-list":
		res, err = mockTransactionList()
	case "POST_get-product":
		res, err = mockGetProduct(req)
	case "POST_get-offer-product":
		res, err = mockGetOfferProduct()
	case "POST_check-active-product":
		res, err = mockCheckActiveProduct(req)
	case "GET_get-all-cities":
		res, err = mockGetAllCities()
	case "GET_product-code-list":
		res, err = mockGetProductCodeList()
	case "GET_relationships":
		res, err = mockGetRelationships()
	default:
		res = []byte("{\"code\": 404, \"message\": \"not found\"}")
		status = 404
	}

	if err != nil {
		status = 500
		res = []byte("{\"code\": 500, \"message\": \"mock error\"}")
	}

	return &http.Response{
		StatusCode: status,
		Body:       ioutil.NopCloser(bytes.NewReader(res)),
	}, nil
}

func mockTransactResult(req *http.Request) ([]byte, error) {
	now := time.Now()
	res := `{
	  "code": "00",
	  "message": "Success",
	  "result": {
		"sessionID": "%d",
		"statusCode": "00",
		"statusDesc": "Success",
		"insProductID": "UCPB04",
		"insProductDesc": "UCPBGEN PA PAWNSHOP",
		"trnDate": "%s",
		"trnAmount": "150.00",
		"traceNo": "INS%d",
		"clientNo": "18155",
		"numUnits": "1",
		"begPolicyNo": "INSPNIHEDK6795317-1",
		"endPolicyNo": "INSPNIHEDK6795317-1",
		"effectiveDate": "06/07/2022",
		"expiryDate": "06/07/2023",
		"pocPDFLink": "https://uat.coopnet.ph:8443/static/PETNET-INS-POC-INSPNIHEDK6795317.pdf",
		"cocPDFLink": "https://uat.coopnet.ph:8443/static/PETNET-INS-COC-INSPNIHEDK6795317.pdf",
		"partnerCommission": "70.00",
		"tellerCommission": "0.00",
		"timestamp": "%s"
	  }
	}`

	res = fmt.Sprintf(res, now.Unix(), now.Format("01/02/2006"), now.Unix(), now.Format("2006-01-02 15:04:05"))
	return []byte(res), nil
}

func mockGetReprint() ([]byte, error) {
	res := `{
	  "code": "00",
	  "message": "Success",
	  "result": {
		"sessionID": "f625da4bfd67491e92349122b366c2b2",
		"statusCode": "00",
		"statusDesc": "Success",
		"trnAmount": "150.00",
		"trnDate": "06/02/2022",
		"pocPDFLink": "https://uat.coopnet.ph:8443/static/PETNET-INS-POC-INSPNIHEDK6294410.pdf",
		"cocPDFLink": "https://uat.coopnet.ph:8443/static/PETNET-INS-COC-INSPNIHEDK6294410.pdf"
	  }
	}`

	return []byte(res), nil
}

func mockTransactionList() ([]byte, error) {
	res := `{
	  "code": "00",
	  "message": "Success",
	  "result": {
		"sessionID": "f625da4bfd67491e92349122b366c2b2",
		"statusCode": "00",
		"statusDesc": "Success",
		"transactions": [
		  {
			"sysCode": "RNT",
			"clientCode": "PNI",
			"brhCode": "HED",
			"trnDate": "06/02/2022",
			"traceNo": "INSPNIHEDK6293243",
			"clientNo": "17063",
			"lastName": "van der Linden",
			"firstName": "Seth",
			"middleName": "Cheng",
			"gender": "M",
			"birthDate": "02/26/2000",
			"cellNo": "09691531894",
			"maritalStatus": "M",
			"occupation": "VICE PRESIDENT",
			"insGroupID": "IPAPF",
			"insProductID": "UCPB04",
			"insProductDesc": "UCPBGEN PA PAWNSHOP (FAMILY)",
			"insurerCode": "UCP",
			"insuranceType": "PA",
			"begPolicyNo": "INSPNIHEDK6293243-1",
			"endPolicyNo": "INSPNIHEDK6293243-1",
			"coverageInMos": "12",
			"contestabilityInMos": "12",
			"effectiveDate": "06/02/2022",
			"expiryDate": "06/02/2023",
			"insCardNo": "",
			"ben1LastName": "Clark",
			"ben1FirstName": "Blanche",
			"ben1MiddleName": "Sample",
			"ben1Relationship": "PAR",
			"ben2LastName": "",
			"ben2FirstName": "",
			"ben2MiddleName": "",
			"ben2Relationship": "",
			"ben3LastName": "",
			"ben3FirstName": "",
			"ben3MiddleName": "",
			"ben3Relationship": "",
			"ben4LastName": "",
			"ben4FirstName": "",
			"ben4MiddleName": "",
			"ben4Relationship": "",
			"loanAmt": "0.0",
			"loanTerm": "",
			"noOfMonths": "0",
			"numUnits": "1",
			"perUnitFee": "150.0",
			"trnAmt": "100.0",
			"trnFee": "0.0",
			"totAmt": "150.0",
			"staffLoginName": "000:designex",
			"staffCommission": "0.0",
			"insTrnStatus": "ACT",
			"provinceCode": "1",
			"cityCode": "1",
			"address": "1610 Ijan Pike"
		  }
		]
	  }
	}`

	return []byte(res), nil
}

func mockGetProduct(req *http.Request) ([]byte, error) {
	res := `{
		  "code": "00",
		  "message": "Success",
		  "result": {
			"sessionID": "f1949eb858414be895d672dfade8cccd",
			"statusCode": "00",
			"statusDesc": "Success",
			"product": {
			  "insGroupID": "IPAH",
			  "insProductID": "%s",
			  "insProductDesc": "%s COVID-19 PRIME -YEAR",
			  "insurerCode": "UCP",
			  "insuranceType": "PA",
			  "insuranceCategory": "IND",
			  "policyNo": "PA-MCR-HO-21-0000007-00",
			  "minAge": "18",
			  "maxAge": "60",
			  "coverageInMos": "12.00",
			  "contestabilityInMos": "12.00",
			  "activationDelay": "2.00",
			  "maxUnits": "1.00",
			  "perUnitFee": "450.00",
			  "product_name": "PETNET COVID-19 Assist %s"
			},
			"coverages": [
			  {
				"insCoverageID": "1",
				"insCoverageDesc": "Accidental Death & Disablement",
				"insCoverageIconID": "cross",
				"insCoverageType1": "AM",
				"insCoverageAmt1": "30,000.00",
				"insCoverageType2": "",
				"insCoverageAmt2": "0.00",
				"insCoverageType3": "",
				"insCoverageAmt3": "0.00",
				"insCoverageType4": "",
				"insCoverageAmt4": "0.00",
				"insCoverageType5": "",
				"insCoverageAmt5": "0.00"
			  }
			]
		  }
		}`

	var b struct {
		ProductCode string `json:"product_code"`
	}

	err := json.NewDecoder(req.Body).Decode(&b)
	if err != nil {
		return nil, err
	}

	res = fmt.Sprintf(res, b.ProductCode, b.ProductCode, b.ProductCode)
	return []byte(res), nil
}

func mockGetOfferProduct() ([]byte, error) {
	res := `{
		  "code": "0",
		  "message": "Success",
		  "result": {
			"product_name": "PETNET COVID-19 Assist",
			"product_code": "UCPB15",
			"product_type": "PRIME 12 mos",
			"dependents": 0,
			"beneficiary": 1,
			"beneficiary_policy": {
			  "max": 4,
			  "min": 1
			},
			"age_policy": {
			  "insurer": {
				"max_age": 60,
				"min_age": 18
			  }
			},
			"end_spiels_title": "Mag-avail na ng PETNET Covid-19 Assist (Prime)",
			"end_spiels_description": "Php 1.25/day lang, may panlaban na sa Covid-19!",
			"sales_pitch": "1 taong tanggal-pangamba sa Covid-19 sa halagang Php 450 lang.",
			"terms_and_condition": "https://register.cashko-insurance.com/terms-and-conditions/petnet450",
			"data_privacy": "https://register.cashko-insurance.com/terms-and-conditions/Privacy"
		  }
		}`

	return []byte(res), nil
}

func mockCheckActiveProduct(req *http.Request) ([]byte, error) {
	res := `{
		  "code": "0",
		  "message": "Success",
		  "result": {
			"product_name": "PETNET MedPamilya %s",
			"product_code": "%s",
			"dependents": 1,
			"beneficiary": 1,
			"dependents_policy": {
			  "max": 4,
			  "min": 1
			},
			"beneficiary_policy": {
			  "max": 4,
			  "min": 1
			},
			"age_policy": {
			  "insurer": {
				"max_age": 70,
				"min_age": 18
			  },
			  "dependents": {
				"CHI": {
				  "max_age": 21,
				  "min_age": 0
				},
				"PAR": {
				  "max_age": 70,
				  "min_age": 18
				},
				"SIB": {
				  "max_age": 21,
				  "min_age": 0
				},
				"SPS": {
				  "max_age": 70,
				  "min_age": 18
				}
			  }
			},
			"end_spiels_title": "Mag-avail na ng PETNET MedPAmilya",
			"end_spiels_description": "Php 1.39/day lang, protektado ang buong pamilya sa di inaasahang aksidente o sakit!",
			"sales_pitch": "Depensa ng buong pamilya laban sa di inaasahang sakit o aksidente, sa Php 500 lang sa 1 taon.",
			"terms_and_condition": "https://register.cashko-insurance.com/terms-and-conditions/petnet500",
			"data_privacy": "https://register.cashko-insurance.com/terms-and-conditions/Privacy"
		  }
		}`

	var b struct {
		LastName    string `json:"last_name"`
		FirstName   string `json:"first_name"`
		MiddleName  string `json:"middle_name"`
		Birthdate   string `json:"birthdate"`
		Gender      string `json:"gender"`
		ProductCode string `json:"product_code"`
	}

	err := json.NewDecoder(req.Body).Decode(&b)
	if err != nil {
		return nil, err
	}

	res = fmt.Sprintf(res, b.ProductCode, b.ProductCode)
	return []byte(res), nil
}

func mockGetAllCities() ([]byte, error) {
	res := `{
	  "code": "00",
	  "message": "Success",
	  "result": {
		"sessionID": "f1949eb858414be895d672dfade8cccd",
		"statusCode": "00",
		"statusDesc": "Success",
		"allcities": [
		  {
			"cityCode": "1",
			"cityName": "BANGUED",
			"provinceCode": "1",
			"provinceName": "ABRA"
		  }
		]
	  }
	}`

	return []byte(res), nil
}

func mockGetProductCodeList() ([]byte, error) {
	res := `{
	  "code": "0",
	  "message": "Success",
	  "result": [
		{
		  "product_name": "PETNET MedPamilya",
		  "product_code": "UCPB02",
		  "dependents": 1,
		  "beneficiary": 1,
		  "dependents_policy": {
			"max": 4,
			"min": 1
		  },
		  "beneficiary_policy": {
			"max": 4,
			"min": 1
		  },
		  "age_policy": {
			"insurer": {
			  "max_age": 70,
			  "min_age": 18
			},
			"dependents": {
			  "CHI": {
				"max_age": 21,
				"min_age": 0
			  },
			  "PAR": {
				"max_age": 70,
				"min_age": 18
			  },
			  "SIB": {
				"max_age": 21,
				"min_age": 0
			  },
			  "SPS": {
				"max_age": 70,
				"min_age": 18
			  }
			}
		  },
		  "end_spiels_title": "Mag-avail na ng PETNET MedPAmilya",
		  "end_spiels_description": "Php 1.39/day lang, protektado ang buong pamilya sa di inaasahang aksidente o sakit!",
		  "sales_pitch": "Depensa ng buong pamilya laban sa di inaasahang sakit o aksidente, sa Php 500 lang sa 1 taon.",
		  "terms_and_condition": "https://register.cashko-insurance.com/terms-and-conditions/petnet500",
		  "data_privacy": "https://register.cashko-insurance.com/terms-and-conditions/Privacy"
		}
	  ]
	}`

	return []byte(res), nil
}

func mockGetRelationships() ([]byte, error) {
	res := `{
	  "code": "0",
	  "message": "Success",
	  "result": [
		{
		  "relationship": "Spouse",
		  "relationship_value": "SPS"
		}
	  ]
	}`

	return []byte(res), nil
}

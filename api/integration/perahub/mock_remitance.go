package perahub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc/status"
)

var (
	errRemitValidateSendMoney         = errors.New("remitance validate send money error")
	errRemitConfirmSendMoney          = errors.New("remitance confirm send money error")
	errRemitCancelSendMoney           = errors.New("remitance cancel send money error")
	errRemitanceInquire               = errors.New("remitance inquire error")
	errRemitConfirmReceiveMoney       = errors.New("remitance confirm receive money error")
	errRemitancePartnersGrid          = errors.New("remitance partners grid error")
	errRemitancePartnersCreate        = errors.New("remitance partners create error")
	errRemitPurposeOfRemittanceGrid   = errors.New("remitance purpose of remittance grid error")
	errPurposeOfRemitGet              = errors.New("purpose of remittance get error")
	errRemitPurposeOfRemittanceUpdate = errors.New("remitance purpose of remittance update error")
	errRemitPurposeOfRemitCreate      = errors.New("purpose of remittance create error")
	errRemitSourceOfFundGrid          = errors.New("source of fund grid error")
	errRemitSourceOfFundCreate        = errors.New("source of fund create error")
	errRemitSourceOfFundGet           = errors.New("source of fund get error")
	errRemittanceEmploymentCreate     = errors.New("employment of remittance create error")
	errRemittanceRelationshipGet      = errors.New("relationship get error")
	errRemitSourceOfFundUpdate        = errors.New("source of fund update error")
	errRemitSourceOfFundDelete        = errors.New("source of fund delete error")
	errPurposeOfRemittanceDelete      = errors.New("remitance purpose of remittance delete error")
	errRemitRelationshipDelete        = errors.New("relationship delete error")
	errRemittanceEmploymentGet        = errors.New("remittance employment get error")
	errRemittanceEmploymentUpdate     = errors.New("remitance Employment update error")
	errRemitEmploymentDelete          = errors.New("employment delete error")
	errRemittanceRelationshipGrid     = errors.New("remitance Relationship grid error")
	errRemittanceRelationshipUpdate   = errors.New("remitance Relationship update error")
)

var remittanceDynamicURL = []string{"purpose", "partner", "occupation", "employment", "sourcefund", "relationship"}

func remitanceAndAction(path string) string {
	return strings.ReplaceAll(path, "/v1/remit/dmt/", "")
}

func dynamicUrlModify(req *http.Request, act string) (rt string) {
	rt = act
	suf := fmt.Sprintf("%s_", req.Method)
	for _, v := range remittanceDynamicURL {
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

func isRemitanceAction(path string) bool {
	return strings.Contains(path, "/remit/dmt/")
}

func isPerahubGetRemocoIDAction(path string) bool {
	return strings.Contains(path, "/transactions/api/drp/remco")
}

func (m *HTTPMock) RemitanceReq(req *http.Request) (*http.Response, error) {
	m.httpHeaders = req.Header
	var rbb []byte
	var err error
	act := remitanceAndAction(req.URL.Path)
	act = dynamicUrlModify(req, act)
	switch act {
	case "POST_send/validate":
		rbb, err = m.SendValidate(req)
	case "POST_send/confirm":
		rbb, err = m.ConfirmMoney(req)
	case "POST_send/cancel":
		rbb, err = m.CancelSendMoney(req)
	case "POST_receive/validate":
		rbb, err = m.ValidateReceiveMoney(req)
	case "POST_inquire":
		rbb, err = m.Inquire(req)
	case "POST_receive/confirm":
		rbb, err = m.ConfirmReceiveMoney(req)
	case "GET_partner":
		rbb, err = m.RemittancePartnersGrid(req)
	case "POST_partner":
		rbb, err = m.RemittancePartnersCreate(req)
	case "GET_purpose":
		rbb, err = m.PurposeOfRemittanceGrid(req)
	case "GET_purpose/{ID}":
		rbb, err = m.PurposeOfRemittanceGet(req)
	case "PUT_purpose/{ID}":
		rbb, err = m.PurposeOfRemittanceUpdate(req)
	case "POST_purpose":
		rbb, err = m.PurposeOfRemittanceCreate(req)
	case "GET_sourcefund":
		rbb, err = m.SourceOFFundGrid(req)
	case "POST_sourcefund":
		rbb, err = m.SourceOfFundCreate(req)
	case "GET_sourcefund/{ID}":
		rbb, err = m.SourceOfFundGet(req)
	case "GET_employment":
		rbb, err = m.RemittanceEmploymentGrid(req)
	case "POST_employment":
		rbb, err = m.RemittanceEmploymentCreate(req)
	case "PUT_employment/{ID}":
		rbb, err = m.RemittanceEmploymentUpdate(req)
	case "GET_occupation":
		rbb, err = m.OccupationGrid(req)
	case "GET_occupation/{ID}":
		rbb, err = m.OccupationGet(req)
	case "POST_occupation":
		rbb, err = m.OccupationCreate(req)
	case "PUT_occupation/{ID}":
		rbb, err = m.OccupationUpdate(req)
	case "DELETE_occupation/{ID}":
		rbb, err = m.OccupationDelete(req)
	case "GET_relationship/{ID}":
		rbb, err = m.RelationshipGet(req)
	case "PUT_sourcefund/{ID}":
		rbb, err = m.SourceOfFundUpdate(req)
	case "DELETE_sourcefund/{ID}":
		rbb, err = m.SourceOfFundDelete(req)
	case "DELETE_purpose/{ID}":
		rbb, err = m.PurposeOfRemittanceDelete(req)
	case "DELETE_relationship/{ID}":
		rbb, err = m.RelationshipDelete(req)
	case "GET_employment/{ID}":
		rbb, err = m.RemittanceEmploymentGet(req)
	case "DELETE_employment/{ID}":
		rbb, err = m.RemittanceEmploymentDelete(req)
	case "GET_relationship":
		rbb, err = m.RemittanceRelationshiptGrid(req)
	case "PUT_relationship/{ID}":
		rbb, err = m.RemittanceRelationshipUpdate(req)
	case "DELETE_partner/{ID}":
		rbb, err = m.RemittancePartnerDelete(req)
	case "GET_partner/{ID}":
		rbb, err = m.RemittancePartnersGet(req)
	case "PUT_partner/{ID}":
		rbb, err = m.RemittancePartnersUpdate(req)
	case "POST_relationship":
		rbb, err = m.RelationshipCreate(req)
	}
	sc := 200
	switch err {
	case errRemitValidateSendMoney, errRemitConfirmSendMoney, errRemitCancelSendMoney, errRemitanceInquire, errRemitConfirmReceiveMoney, errRemitancePartnersGrid, errRemitPurposeOfRemittanceGrid, errPurposeOfRemitGet, errRemitPurposeOfRemittanceUpdate, errRemitPurposeOfRemitCreate, errRemitancePartnersCreate, errRemitSourceOfFundGrid, errRemitSourceOfFundCreate, errRemitSourceOfFundGet, errRemittanceEmploymentCreate, errRemittanceRelationshipGet, errRemitSourceOfFundUpdate, errRemitSourceOfFundDelete, errPurposeOfRemittanceDelete, errRemitRelationshipDelete, errRemittanceEmploymentGet, errRemittanceEmploymentUpdate, errRemitEmploymentDelete, errRemittanceRelationshipGrid, errRemittanceRelationshipUpdate:
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

func (m *HTTPMock) SendValidate(req *http.Request) ([]byte, error) {
	if !m.remitanceErr {
		rb := &RemitanceValidateSendMoneyRes{
			Code:    200,
			Message: "Good",
			Result: RemitanceValidateSendMoneyResult{
				SendValidateReferenceNumber: "1653296685161",
			},
		}
		return json.Marshal(rb)

	}

	rb := &remitanceError{
		Code:    "422",
		Message: "The given data was invalid.",
		Error: map[string][]string{
			"partner_reference_number": {
				"The partner reference number field is required.",
			},
			"principal_amount": {
				"The principal amount field is required.",
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) ConfirmMoney(req *http.Request) ([]byte, error) {
	if !m.remitanceErr {
		rb := &RemitanceConfirmSendMoneyRes{
			Code:    200,
			Message: "Successful",
			Result: RemitanceConfirmSendMoneyResult{
				Phrn: "PH1654787564",
			},
		}
		return json.Marshal(rb)

	}
	rb := &remitanceError{
		Code:    "422",
		Message: "The given data was invalid.",
		Error: map[string][]string{
			"send_validate_reference_number": {
				"The send validate reference number field is required.",
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) CancelSendMoney(req *http.Request) ([]byte, error) {
	if !m.remitanceErr {
		rb := &RemitanceCancelSendMoneyRes{
			Code:    200,
			Message: "Cancel Send Remittance successful",
			Result: RemitanceCancelSendMoneyResult{
				Phrn:                      "PH1654789142",
				CancelSendDate:            "2022-06-14 19:35:02",
				CancelSendReferenceNumber: "6a6e74400561a300d627aba12107bb6c",
			},
		}
		return json.Marshal(rb)
	}

	rb := &remitanceError{
		Code:    "422",
		Message: "The given data was invalid.",
		Error: map[string][]string{
			"partner_code": {
				"The partner_code field is required.",
			},
			"trx_type": {
				"The trx_type field is required.",
			},
			"account_number": {
				"The account_number field is required.",
			},
			"amount": {
				"The amount field is required.",
			},
			"service_charge": {
				"The service_charge field is required.",
			},
			"partner_reference_no": {
				"The partner_reference_no field is required.",
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) ValidateReceiveMoney(req *http.Request) ([]byte, error) {
	if !m.remitanceErr {
		rb := &RemitanceValidateReceiveMoneyRes{
			Code:    200,
			Message: "Successful",
			Result: RemitanceValidateReceiveMoneyResult{
				PayoutValidateReferenceNumber: "4f8a09d3b293807aa50305f66d6cc73c",
			},
		}
		return json.Marshal(rb)

	}

	rb := &remitanceError{
		Code:    "422",
		Message: "The given data was invalid.",
		Error: map[string][]string{
			"phrn": {
				"The phrn field is required.",
			},
			"payout_partner_code": {
				"The payout partner code field is required.",
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) Inquire(req *http.Request) ([]byte, error) {
	if !m.remitanceErr {
		rb := &RemitanceInquireRes{
			Code:    200,
			Message: "PeraHUB Reference Number (PHRN) is available for Payout",
			Result: RemitanceInquireResult{
				Phrn:                  "PH1658296732",
				PrincipalAmount:       10000,
				IsoCurrency:           "PHP",
				ConversionRate:        1,
				IsoOriginatingCountry: "PHP",
				IsoDestinationCountry: "PHP",
				SenderLastName:        "HERMO",
				SenderFirstName:       "IRENE",
				SenderMiddleName:      "M",
				ReceiverLastName:      "HERMO",
				ReceiverFirstName:     "SONNY",
				ReceiverMiddleName:    "D",
			},
		}
		return json.Marshal(rb)
	}

	rb := &remitanceError{
		Code:    "422",
		Message: "The given data was invalid.",
		Error: map[string][]string{
			"phrn": {
				"The phrn field is required.",
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) ConfirmReceiveMoney(req *http.Request) ([]byte, error) {
	if !m.remitanceErr {
		rb := &RemitanceConfirmReceiveMoneyRes{
			Code:    200,
			Message: "Successful",
			Result: RemitanceConfirmReceiveMoneyResult{
				Phrn:                  "PH1654789142",
				PrincipalAmount:       10000,
				IsoOriginatingCountry: "PHP",
				IsoDestinationCountry: "PHP",
				SenderLastName:        "HERMO",
				SenderFirstName:       "IRENE",
				SenderMiddleName:      "M",
				ReceiverLastName:      "HERMO",
				ReceiverFirstName:     "SONNY",
				ReceiverMiddleName:    "D",
			},
		}
		return json.Marshal(rb)
	}

	rb := &remitanceError{
		Code:    "422",
		Message: "The given data was invalid.",
		Error: map[string][]string{
			"payout_validate_reference_number": {
				"The payout validate reference number field is required.",
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) RemittancePartnersGrid(req *http.Request) ([]byte, error) {
	rb := &RemittancePartnersGridRes{
		Code:    200,
		Message: "Good",
		Result: []RemittancePartnersGridResult{
			{
				ID:           1,
				PartnerCode:  "DRP",
				PartnerName:  "PERA HUB",
				ClientSecret: "26da230221d9e506b1fd823df1869875",
				Status:       1,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				DeletedAt:    time.Now(),
			},
			{
				ID:           2,
				PartnerCode:  "USP",
				PartnerName:  "PERA HUB",
				ClientSecret: "12358fbef0bb08d7a7bab57df956a335",
				Status:       1,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				DeletedAt:    time.Now(),
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) PurposeOfRemittanceGrid(req *http.Request) ([]byte, error) {
	rb := &PurposeOfRemittanceGridRes{
		Code:    200,
		Message: "Good",
		Result: []PurposeOfRemittanceGridResult{
			{
				ID:                  "1",
				PurposeOfRemittance: "Gift",
				CreatedAt:           time.Now(),
				UpdatedAt:           time.Now(),
				DeletedAt:           time.Now(),
			},
			{
				ID:                  "2",
				PurposeOfRemittance: "Fund",
				CreatedAt:           time.Now(),
				UpdatedAt:           time.Now(),
				DeletedAt:           time.Now(),
			},
			{
				ID:                  "3",
				PurposeOfRemittance: "Allowance",
				CreatedAt:           time.Now(),
				UpdatedAt:           time.Now(),
				DeletedAt:           time.Now(),
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) RemittanceEmploymentGrid(req *http.Request) ([]byte, error) {
	rb := &RemittanceEmploymentGridRes{
		Code:    200,
		Message: "Good",
		Result: []RemittanceEmploymentGridResult{
			{
				ID:               1,
				EmploymentNature: "REGULAR",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
				DeletedAt:        time.Now(),
			},
			{
				ID:               2,
				EmploymentNature: "PROBATIONARY",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
				DeletedAt:        time.Now(),
			},
			{
				ID:               3,
				EmploymentNature: "CONTRACTUAL",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
				DeletedAt:        time.Now(),
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) PurposeOfRemittanceGet(req *http.Request) ([]byte, error) {
	rb := &PurposeOfRemittanceGetRes{
		Code:    200,
		Message: "Good",
		Result: &PurposeOfRemittanceGetResult{
			ID:                  1,
			PurposeOfRemittance: "Donation",
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			DeletedAt:           time.Now(),
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) PurposeOfRemittanceUpdate(req *http.Request) ([]byte, error) {
	if !m.remitanceErr {
		rb := &PurposeOfRemittanceUpdateRes{
			Code:    200,
			Message: "Good",
			Result: PurposeOfRemittanceUpdateResult{
				ID:                  1,
				PurposeOfRemittance: "Donation",
				CreatedAt:           time.Now(),
				UpdatedAt:           time.Now(),
				DeletedAt:           time.Now(),
			},
		}
		return json.Marshal(rb)
	}

	rb := &remitanceError{
		Code:    "422",
		Message: "The given data was invalid.",
		Error: map[string][]string{
			"purpose_of_remittance": {
				"The purpose of remittance field is required.",
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) PurposeOfRemittanceCreate(req *http.Request) ([]byte, error) {
	if !m.remitanceErr {
		rb := &PurposeOfRemittanceCreateRes{
			Code:    200,
			Message: "Good",
			Result: PurposeOfRemittanceCreateResult{
				ID:                  1,
				PurposeOfRemittance: "USP",
				CreatedAt:           time.Now(),
				UpdatedAt:           time.Now(),
			},
		}
		return json.Marshal(rb)
	}

	rb := &remitanceError{
		Code:    "422",
		Message: "The given data was invalid.",
		Error: map[string][]string{
			"purpose_of_remittance": {
				"The purpose of remittance field is required.",
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) RemittancePartnersCreate(req *http.Request) ([]byte, error) {
	if !m.remitanceErr {
		rb := &RemittancePartnersCreateRes{
			Code:    200,
			Message: "Good",
			Result: RemittancePartnersCreateResult{
				ID:           1,
				PartnerCode:  "USP",
				PartnerName:  "PERA HUB",
				ClientSecret: "adawdawdawd",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
		}
		return json.Marshal(rb)
	}

	rb := &remitanceError{
		Code:    "422",
		Message: "The given data was invalid.",
		Error: map[string][]string{
			"partner_code": {
				"The partner code field is required.",
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) SourceOFFundGrid(req *http.Request) ([]byte, error) {
	rb := &SourceOfFundGridRes{
		Code:    200,
		Message: "Good",
		Result: []SourceOfFundGridResult{
			{
				ID:           "1",
				SourceOfFund: "SALARY",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				DeletedAt:    time.Now(),
			},
			{
				ID:           "2",
				SourceOfFund: "BUSINESS",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				DeletedAt:    time.Now(),
			},
			{
				ID:           "3",
				SourceOfFund: "REMITTANCE",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				DeletedAt:    time.Now(),
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) SourceOfFundCreate(req *http.Request) ([]byte, error) {
	if !m.remitanceErr {
		rb := &SourceOfFundCreateRes{
			Code:    200,
			Message: "Good",
			Result: SourceOfFundCreateResult{
				ID:           1,
				SourceOfFund: "SALARY",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
		}
		return json.Marshal(rb)
	}

	rb := &remitanceError{
		Code:    "422",
		Message: "The given data was invalid.",
		Error: map[string][]string{
			"purpose_of_remittance": {
				"The source of fund field is required.",
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) SourceOfFundGet(req *http.Request) ([]byte, error) {
	rb := &SourceOfFundGetRes{
		Code:    200,
		Message: "Good",
		Result: &SourceOfFundGetResult{
			ID:           1,
			SourceOfFund: "SALARY",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			DeletedAt:    time.Now(),
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) OccupationGrid(req *http.Request) ([]byte, error) {
	rb := &OccupationGridRes{
		Code:    200,
		Message: "Good",
		Result: []OccupationGridResult{
			{
				ID:         1,
				Occupation: "Programmer",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
				DeletedAt:  time.Now(),
			},
			{
				ID:         2,
				Occupation: "Engineer",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
				DeletedAt:  time.Now(),
			},
			{
				ID:         3,
				Occupation: "Doctor",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
				DeletedAt:  time.Now(),
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) OccupationGet(req *http.Request) ([]byte, error) {
	rb := &OccupationGetRes{
		Code:    200,
		Message: "Good",
		Result: &OccupationGetResult{
			ID:         1,
			Occupation: "Programmer",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			DeletedAt:  time.Now(),
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) OccupationCreate(req *http.Request) ([]byte, error) {
	if !m.remitanceErr {
		rb := &OccupationCreateRes{
			Code:    200,
			Message: "Good",
			Result: OccupationCreateResult{
				ID:         1,
				Occupation: "Programmer",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		}
		return json.Marshal(rb)
	}

	rb := &remitanceError{
		Code:    "422",
		Message: "The given data was invalid.",
		Error: map[string][]string{
			"partner_code": {
				"The occupation field is required.",
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) OccupationUpdate(req *http.Request) ([]byte, error) {
	if !m.remitanceErr {
		rb := &OccupationUpdateRes{
			Code:    200,
			Message: "Good",
			Result: OccupationUpdateResult{
				ID:         1,
				Occupation: "Programmer",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
				DeletedAt:  time.Now(),
			},
		}
		return json.Marshal(rb)
	}

	rb := &remitanceError{
		Code:    "422",
		Message: "The given data was invalid.",
		Error: map[string][]string{
			"partner_code": {
				"The occupation field is required.",
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) OccupationDelete(req *http.Request) ([]byte, error) {
	rb := &OccupationDeleteRes{
		Code:    200,
		Message: "Good",
		Result: OccupationDeleteResult{
			ID:         1,
			Occupation: "Programmer",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			DeletedAt:  time.Now(),
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) RemittanceEmploymentCreate(req *http.Request) ([]byte, error) {
	if !m.remitanceErr {
		rb := &RemittanceEmploymentCreateResp{
			Code:    200,
			Message: "Good",
			Result: RemittanceEmploymentCreateResult{
				ID:               1,
				EmploymentNature: "REGULAR",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
		}
		return json.Marshal(rb)
	}
	rb := &remitanceError{
		Code:    "422",
		Message: "The given data was invalid.",
		Error: map[string][]string{
			"remittance_employment_create": {
				"Remittance Employment Create field is required.",
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) RelationshipGet(req *http.Request) ([]byte, error) {
	rb := &RelationshipGetRes{
		Code:    200,
		Message: "Good",
		Result: &RelationshipGetResult{
			ID:           1,
			Relationship: "Friend",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			DeletedAt:    time.Now(),
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) SourceOfFundUpdate(req *http.Request) ([]byte, error) {
	if !m.remitanceErr {
		rb := &SourceOfFundUpdateRes{
			Code:    200,
			Message: "Good",
			Result: SourceOfFundUpdateResult{
				ID:           1,
				SourceOfFund: "SALARY",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				DeletedAt:    time.Now(),
			},
		}
		return json.Marshal(rb)
	}
	rb := &remitanceError{
		Code:    "422",
		Message: "The given data was invalid.",
		Error: map[string][]string{
			"source_of_fund": {
				"The source of fund field is required.",
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) SourceOfFundDelete(req *http.Request) ([]byte, error) {
	rb := &SourceOfFundDeleteRes{
		Code:    200,
		Message: "Good",
		Result: SourceOfFundDeleteResult{
			ID:           1,
			SourceOfFund: "SALARY",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			DeletedAt:    time.Now(),
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) PurposeOfRemittanceDelete(req *http.Request) ([]byte, error) {
	var rbb []byte
	var err error
	rb := &PurposeOfRemittanceDeleteRes{
		Code:    200,
		Message: "Good",
		Result: PurposeOfRemittanceDeleteResult{
			ID:                  1,
			PurposeOfRemittance: "Gift",
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			DeletedAt:           time.Now(),
		},
	}
	rbb, err = json.Marshal(rb)
	if err != nil {
		return nil, err
	}
	return rbb, err
}

func (m *HTTPMock) RelationshipDelete(req *http.Request) ([]byte, error) {
	rb := &RelationshipDeleteRes{
		Code:    200,
		Message: "Good",
		Result: RelationshipDeleteResult{
			ID:           1,
			Relationship: "Friend",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			DeletedAt:    time.Now(),
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) RemittanceEmploymentGet(req *http.Request) ([]byte, error) {
	rb := &RemittanceEmploymentGetRes{
		Code:    200,
		Message: "Good",
		Result: &RemittanceEmploymentGetResult{
			ID:               1,
			EmploymentNature: "REGULAR",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			DeletedAt:        time.Now(),
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) RemittanceEmploymentUpdate(req *http.Request) ([]byte, error) {
	if !m.remitanceErr {
		rb := &RemittanceEmploymentUpdateRes{
			Code:    200,
			Message: "Good",
			Result: RemittanceEmploymentUpdateResult{
				ID:               1,
				EmploymentNature: "REGULAR",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
				DeletedAt:        time.Now(),
			},
		}
		return json.Marshal(rb)
	}
	rb := &remitanceError{
		Code:    "422",
		Message: "The given data was invalid.",
		Error: map[string][]string{
			"employment": {
				"The employment field is required.",
			},
			"employment_nature": {
				"The employment_nature field is required.",
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) RemittanceEmploymentDelete(req *http.Request) ([]byte, error) {
	rb := &RemittanceEmploymentDeleteRes{
		Code:    200,
		Message: "Good",
		Result: &RemittanceEmploymentDeleteResult{
			ID:               1,
			EmploymentNature: "REGULAR",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			DeletedAt:        time.Now(),
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) RemittanceRelationshiptGrid(req *http.Request) ([]byte, error) {
	if !m.remitanceErr {
		rb := &RemittanceRelationshipGridRes{
			Code:    200,
			Message: "Good",
			Result: []RemittanceRelationshipGridResult{
				{
					ID:           1,
					Relationship: "Friend",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
					DeletedAt:    time.Now(),
				},
				{
					ID:           2,
					Relationship: "Father",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
					DeletedAt:    time.Now(),
				},
				{
					ID:           3,
					Relationship: "Mother",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
					DeletedAt:    time.Now(),
				},
			},
		}
		return json.Marshal(rb)
	}
	rb := &remitanceError{
		Code:    "422",
		Message: "The given data was invalid.",
		Error: map[string][]string{
			"code": {
				"relationship is required.",
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) RemittanceRelationshipUpdate(req *http.Request) ([]byte, error) {
	if !m.remitanceErr {
		rb := &RelationshipUpdateRes{
			Code:    200,
			Message: "Good",
			Result: RelationshipUpdateResult{
				ID:           1,
				Relationship: "Friend",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				DeletedAt:    time.Now(),
			},
		}
		return json.Marshal(rb)
	}
	rb := &remitanceError{
		Code:    "422",
		Message: "The given data was invalid.",
		Error: map[string][]string{
			"relationship": {
				"The relationship field is required.",
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) RemittancePartnerDelete(req *http.Request) ([]byte, error) {
	rb := &RemittancePartnersDeleteRes{
		Code:    200,
		Message: "Good",
		Result: RemittancePartnersDeleteResult{
			ID:           1,
			PartnerCode:  "USP",
			PartnerName:  "PERA HUB",
			ClientSecret: "adawdawdawd",
			Status:       1,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			DeletedAt:    time.Now(),
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) RemittancePartnersGet(req *http.Request) ([]byte, error) {
	rb := &RemittancePartnersGetRes{
		Code:    200,
		Message: "Good",
		Result: &RemittancePartnersGetResult{
			ID:           1,
			PartnerCode:  "DRP",
			PartnerName:  "BRANKAS",
			ClientSecret: "4fab1de660a6b7faef0168ca4788408a",
			Status:       1,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			DeletedAt:    "null",
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) RemittancePartnersUpdate(req *http.Request) ([]byte, error) {
	if !m.remitanceErr {
		rb := &RemittancePartnersUpdateRes{
			Code:    200,
			Message: "Good",
			Result: RemittancePartnersUpdateResult{
				ID:           1,
				PartnerCode:  "USP",
				PartnerName:  "PERA HUB",
				ClientSecret: "adawdawdawd",
				Status:       1,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				DeletedAt:    "null",
			},
		}
		return json.Marshal(rb)
	}
	rb := &remitanceError{
		Code:    "422",
		Message: "The given data was invalid.",
		Error: map[string][]string{
			"partner_code": {
				"The partner code field is required.",
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) RelationshipCreate(req *http.Request) ([]byte, error) {
	if !m.remitanceErr {
		rb := &RemittanceRelationshipCreateRes{
			Code:    200,
			Message: "Good",
			Result: RemittanceRelationshipCreateResult{
				Relationship: "Friend",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				ID:           1,
			},
		}
		return json.Marshal(rb)
	}

	rb := &remitanceError{
		Code:    "422",
		Message: "The given data was invalid.",
		Error: map[string][]string{
			"relationship": {
				"The relationship field is required.",
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) GetRemoco(req *http.Request) ([]byte, error) {
	rb := &PerahubGetRemcoIDResponse{
		Code:    "200",
		Message: "Good",
		Result: []PerahubGetRemcoIDResult{
			{
				ID:   "1",
				Name: "iRemit",
			},
			{
				ID:   "2",
				Name: "BPI",
			},
		},
	}
	return json.Marshal(rb)
}

func (m *HTTPMock) GetRemocoIdReq(req *http.Request) (*http.Response, error) {
	m.httpHeaders = req.Header
	code := 200
	rbb, err := m.GetRemoco(req)
	if err != nil {
		code = int(status.Code(err))
	}
	return &http.Response{
		StatusCode: code,
		Body:       ioutil.NopCloser(bytes.NewReader(rbb)),
	}, nil
}

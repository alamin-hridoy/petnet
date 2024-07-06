package perahub

import (
	"context"
	"encoding/json"
	"fmt"

	"brank.as/petnet/serviceutil/logging"
)

type PECustomerRequest struct {
	CustomerCode string `json:"customer_code"`
}

type Image struct {
	Img      string `json:"img"`
	IDType   string `json:"id_type"`
	IDNumber string `json:"id_number"`
	Country  string `json:"country"`
	Expiry   string `json:"expiry"`
	Issue    string `json:"issue"`
}

type PECustomer struct {
	CustomerID     int              `json:"customer_id"`
	CustomerNumber string           `json:"customer_number"`
	LastName       string           `json:"last_name"`
	FirstName      string           `json:"first_name"`
	MiddleName     string           `json:"middle_name"`
	CurrentAddress Address          `json:"current_address"`
	TelNo          string           `json:"tel_no"`
	EmailAdd       string           `json:"email_add"`
	MobileNo       string           `json:"mobile_no"`
	BirthDate      string           `json:"birth_date"`
	Occupation     string           `json:"occupation"`
	PermaAddress   Address          `json:"perma_address"`
	Nationality    string           `json:"nationality"`
	Gender         string           `json:"gender"`
	CivilStatus    string           `json:"civil_status"`
	TINID          string           `json:"tin_id"`
	SSSID          string           `json:"sss_id"`
	GSISID         string           `json:"gsis_id"`
	DriverLic      string           `json:"driver_lic"`
	SourceFund     string           `json:"source_fund"`
	EmployerName   string           `json:"employer_name"`
	NatureWork     string           `json:"nature_work"`
	Employment     string           `json:"employment"`
	CardNo         string           `json:"card_no"`
	Img            map[string]Image `json:"img"`
	CardPoints     int              `json:"card_points"`
	CardPesoVal    int              `json:"card_peso_val"`
	WBCardNo       string           `json:"wu_card_no"`
	UBCardNo       string           `json:"ub_card_no"`
}

type PECustomerResponseWU struct {
	Header ResponseHeader `json:"header"`
	Body   []PECustomer   `json:"body"`
}

type PECustomerResponse struct {
	WU PECustomerResponseWU `json:"uspwuapi"`
}

func (s *Svc) PreEnrolledCustomer(ctx context.Context, usrCd string) (*PECustomer, error) {
	log := logging.FromContext(ctx)

	const mod, modReq = "Registration", "getinfo"
	req, err := s.newParahubRequest(ctx, mod, modReq, PECustomerRequest{CustomerCode: usrCd})
	if err != nil {
		return nil, err
	}

	res, err := s.post(ctx, s.moduleURL(mod, ""), *req)
	if err != nil {
		return nil, err
	}

	var resp []PECustomer
	if err := json.Unmarshal(res, &resp); err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, fmt.Errorf("customer not found")
	}
	if len(resp) > 1 {
		log.WithField("customer_code", usrCd).WithField("records", resp).Error("multiple records")
		return nil, err
	}

	return &resp[0], nil
}

package perahub

import (
	"context"
	"encoding/json"
)

type KYCLookupRequest struct {
	RefNo        string `json:"foreign_reference_no"`
	SearchType   string `json:"searchtype"`
	TrxType      string `json:"trx_type"`
	IDType       string `json:"kyc_id_type"`
	IDNumber     string `json:"kyc_id_number"`
	ContactPhone string `json:"kyc_contact_phone"`
	MyWUNumber   string `json:"kyc_mywu_number"`
	FirstName    string `json:"kyc_first_name"`
	LastName     string `json:"kyc_last_name"`
	IsMulti      string `json:"is_multi"`
	OperatorID   string `json:"operator_id"`
	TerminalID   string `json:"terminal_id"`
	UserCode     string `json:"-"`
}

type KYCCustomer struct {
	Name              Name                 `json:"name"`
	Address           KYCCustomerAddress   `json:"address"`
	MyWUDetails       KYCMyWUDetails       `json:"mywu_details"`
	ComplianceDetails KYCComplianceDetails `json:"compliance_details"`
	Email             string               `json:"email"`
	ContactPhone      string               `json:"contact_phone"`
	MobileNumber      MobileNumber         `json:"mobile_number"`
	DateOfBirth       string               `json:"date_of_birth"`
	SuppressFlag      string               `json:"suppress_flag"`
	MarketingDetails  json.RawMessage      `json:"mktng_details"`
	KYCDetails        KYCDetails           `json:"kyc_details"`
	Preferences       json.RawMessage      `json:"preferences"`
}

type KYCCountryDetails struct {
	CtryCode string `json:"ctry_code"`
	CtryName string `json:"ctry_name"`
}

type KYCCustomerAddress struct {
	AddrType       string            `json:"addr_type"`
	AddrLine1      string            `json:"addr_line1"`
	AddrLine2      string            `json:"addr_line2"`
	City           string            `json:"city"`
	StateName      string            `json:"state_name"`
	PostalCode     string            `json:"postal_code"`
	CountryDetails KYCCountryDetails `json:"country_details"`
}

type KYCMyWUDetails struct {
	MyWUNumber       string      `json:"mywu_number"`
	IsConvenience    string      `json:"is_convenience"`
	LevelCode        string      `json:"level_code"`
	CurrentYrPts     json.Number `json:"current_yr_pts"`
	EnrollmentSource string      `json:"enrollment_source"`
}

type KYCComplianceDetails struct {
	ComplianceFlagsBuffer string `json:"compliance_flags_buffer"`
	ComplianceDataBuffer  string `json:"compliance_data_buffer"`
}

type MobileNumber struct {
	CountryCode    string `json:"ctry_code"`
	NationalNumber string `json:"National_number"`
}

type KYCDetails struct {
	IsKyced string `json:"is_kyced"`
}

type KYCAddress struct {
	AddrType   string            `json:"addr_type"`
	AddrLine1  string            `json:"addr_line1"`
	AddrLine2  string            `json:"addr_line2"`
	AddrLine3  string            `json:"addr_line3"`
	City       string            `json:"city"`
	StateName  string            `json:"state_name"`
	PostalCode string            `json:"postal_code"`
	LocalArea  string            `json:"local_Area"`
	Country    KYCCountryDetails `json:"country_details"`
}

type ReceiverDetails struct {
	Name            Name        `json:"name"`
	ReceiverType    string      `json:"receiver_type"`
	Address         KYCAddress  `json:"address"`
	ReceiverIndexNo json.Number `json:"receiver_index_no"`
}

type KYCReceiver struct {
	Receiver      []ReceiverDetails `json:"receiver"`
	NumberMatches json.Number       `json:"number_matches"`
}

type KYCLookupBody struct {
	Customer KYCCustomer `json:"customer"`
	Receiver KYCReceiver `json:"receiver"`
}

func (s *Svc) KYCLookup(ctx context.Context, kycReq KYCLookupRequest) (*KYCLookupBody, error) {
	kycReq.IsMulti = "S"
	const mod, modReq = "prereq", "kyclookup"
	req, err := s.newParahubRequest(ctx,
		mod, modReq, kycReq,
		WithUserCode(json.Number(kycReq.UserCode)),
		WithLocationCode(kycReq.UserCode))
	if err != nil {
		return nil, err
	}

	res, err := s.post(ctx, s.moduleURL("prereq", ""), *req)
	if err != nil {
		return nil, err
	}

	var kycLookupResponse KYCLookupBody
	if err := json.Unmarshal(res, &kycLookupResponse); err != nil {
		return nil, err
	}

	return &kycLookupResponse, nil
}

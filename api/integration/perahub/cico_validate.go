package perahub

import (
	"context"
	"encoding/json"
)

type CicoValidateRequest struct {
	PartnerCode string               `json:"partner_code"`
	Trx         CicoValidateTrx      `json:"trx"`
	Customer    CicoValidateCustomer `json:"customer"`
}

type CicoValidateCustomer struct {
	CustomerID        string `json:"customer_id"`
	CustomerFirstname string `json:"customer_firstname"`
	CustomerLastname  string `json:"customer_lastname"`
	CurrAddress       string `json:"curr_address"`
	CurrBarangay      string `json:"curr_barangay"`
	CurrCity          string `json:"curr_city"`
	CurrProvince      string `json:"curr_province"`
	CurrCountry       string `json:"curr_country"`
	BirthDate         string `json:"birth_date"`
	BirthPlace        string `json:"birth_place"`
	BirthCountry      string `json:"birth_country"`
	ContactNo         string `json:"contact_no"`
	IDType            string `json:"id_type"`
	IDNumber          string `json:"id_number"`
}

type CicoValidateTrx struct {
	Provider        string `json:"provider"`
	ReferenceNumber string `json:"reference_number"`
	TrxType         string `json:"trx_type"`
	PrincipalAmount int    `json:"principal_amount"`
}

type CicoValidateResult struct {
	PetnetTrackingno   string `json:"petnet_trackingno"`
	TrxDate            string `json:"trx_date"`
	TrxType            string `json:"trx_type"`
	Provider           string `json:"provider"`
	ProviderTrackingno string `json:"provider_trackingno"`
	ReferenceNumber    string `json:"reference_number"`
	PrincipalAmount    int    `json:"principal_amount"`
	Charges            int    `json:"charges"`
	TotalAmount        int    `json:"total_amount"`
	Timestamp          string `json:"timestamp"`
}

type CicoValidateResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Result  *CicoValidateResult `json:"result"`
}

func (s *Svc) CicoValidate(ctx context.Context, sr CicoValidateRequest) (*CicoValidateResponse, error) {
	res, err := s.cicoPost(ctx, s.cicoURL("validate"), sr)
	if err != nil {
		return nil, err
	}

	CicoRes := &CicoValidateResponse{}
	if err := json.Unmarshal(res, CicoRes); err != nil {
		return nil, err
	}
	return CicoRes, nil
}

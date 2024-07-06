package perahub

import (
	"context"
	"encoding/json"
	"net/url"
)

type CebFindBFReq struct {
	SenderClientId string `json:"sender_client_id"`
}

type CebFindBFRes struct {
	Code    json.Number `json:"code"`
	Message string      `json:"message"`
	Result  FBFResult   `json:"result"`
	RemcoID json.Number `json:"remco_id"`
}

type FBFResult struct {
	Beneficiary []Beneficiary `json:"Beneficiary"`
}

type Beneficiary struct {
	BeneficiaryID  json.Number `json:"BeneficiaryID"`
	FirstName      string      `json:"FirstName"`
	MiddleName     string      `json:"MiddleName"`
	LastName       string      `json:"LastName"`
	BirthDate      string      `json:"BirthDate"`
	StateIDAddress json.Number `json:"StateIDAddress"`
	CPCountry      CrtyID      `json:"CellphoneCountry"`
	TPCountry      CrtyID      `json:"TelephoneCountry"`
	CtryAddress    CrtyID      `json:"CountryAddress"`
	BirthCountry   CrtyID      `json:"BirthCountry"`
}

func (s *Svc) CebFindBF(ctx context.Context, sr CebFindBFReq) (*CebFindBFRes, error) {
	nonexUrl := s.nonexURL("cebuana/beneficiary-by-sender?sender_client_id=" + sr.SenderClientId)
	decodedUrl, _ := url.QueryUnescape(nonexUrl)
	res, err := s.getNonex(ctx, decodedUrl)
	if err != nil {
		return nil, err
	}

	rb := &CebFindBFRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

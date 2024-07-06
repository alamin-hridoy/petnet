package perahub

import (
	"context"
	"encoding/json"
	"net/url"
)

type CebFindClientRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	BirthDate    string `json:"birth_date"`
	ClientNumber string `json:"client_number"`
}

type CebFindClientRespBody struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Result  GUResult `json:"result"`
	RemcoID int      `json:"remco_id"`
}

type GUResult struct {
	Client Client `json:"Client"`
}

type Client struct {
	ClientID     int      `json:"ClientID"`
	ClientNumber string   `json:"ClientNumber"`
	FirstName    string   `json:"FirstName"`
	MiddleName   string   `json:"MiddleName"`
	LastName     string   `json:"LastName"`
	BirthDate    string   `json:"BirthDate"`
	CPCountry    CrtyID   `json:"CellphoneCountry"`
	TPCountry    CrtyID   `json:"TelephoneCountry"`
	CtryAddress  CrtyID   `json:"CountryAddress"`
	CSOfFund     CSOfFund `json:"ClientSourceOfFund"`
}

type CrtyID struct {
	CountryID int `json:"CountryID"`
}

type CSOfFund struct {
	SourceOfFundID int `json:"SourceOfFundID"`
}

func (s *Svc) CebFindClient(ctx context.Context, sr CebFindClientRequest) (*CebFindClientRespBody, error) {
	nonexUrl := s.nonexURL("cebuana/find-client?first_name=" + sr.FirstName + "&last_name=" + sr.LastName + "&birth_date=" + sr.BirthDate + "&client_number=" + sr.ClientNumber)
	decodedUrl, _ := url.QueryUnescape(nonexUrl)
	res, err := s.getNonex(ctx, decodedUrl)
	if err != nil {
		return nil, err
	}

	rb := &CebFindClientRespBody{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

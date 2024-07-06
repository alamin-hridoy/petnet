package perahub

import (
	"context"
	"encoding/json"
)

type CebAddBftReq struct {
	FirstName          string `json:"first_name"`
	MiddleName         string `json:"middle_name"`
	LastName           string `json:"last_name"`
	SenderClientID     int    `json:"sender_client_id"`
	BirthDate          string `json:"birth_date"`
	CellphoneCountryID string `json:"cellphone_country_id"`
	ContactNumber      string `json:"contact_number"`
	TelephoneCountryID string `json:"telephone_country_id"`
	TelephoneAreaCode  string `json:"telephone_area_code"`
	TelephoneNumber    string `json:"telephone_number"`
	CountryAddressID   string `json:"country_address_id"`
	BirthCountryID     string `json:"birth_country_id"`
	ProvinceAddress    string `json:"province_address"`
	Address            string `json:"address"`
	UserID             int    `json:"user_id"`
	Occupation         string `json:"occupation"`
	ZipCode            string `json:"zip_code"`
	StateIDAddress     string `json:"state_id_address"`
	Tin                string `json:"tin"`
}

type CebAddBfResp struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Result  ABResult `json:"result"`
	RemcoID int      `json:"remco_id"`
}

type ABResult struct {
	ResultStatus  string `json:"ResultStatus"`
	MessageID     int    `json:"MessageID"`
	LogID         int    `json:"LogID"`
	BeneficiaryID int    `json:"BeneficiaryID"`
}

func (s *Svc) CebAddBf(ctx context.Context, req CebAddBftReq) (*CebAddBfResp, error) {
	res, err := s.postNonex(ctx, s.nonexURL("cebuana/add-beneficiary"), req)
	if err != nil {
		return nil, err
	}

	rb := &CebAddBfResp{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

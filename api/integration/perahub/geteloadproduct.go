package perahub

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
)

type ELPRequest struct {
	CurrentRow int `json:"currow"`
	GroupCount int `json:"groupcount"`
}

type Item struct {
	EloadID     string `json:"eload_id"`
	ProductName string `json:"ProductName"`
	ProductCode string `json:"ProductCode"`
	Amount      string `json:"Amount"`
	Commission  string `json:"Commission"`
	Provider    string `json:"Provider"`
	ProductType string `json:"ProductType"`
	Remarks     string `json:"Remarks"`
	Additional  string `json:"Additional"`
	UpdatedBy   string `json:"updated_by"`
	DateUpdated string `json:"date_updated"`
}

type ItemCount struct {
	TotalRows string `json:"total_rows"`
}

type ELPResponseBody struct {
	Items      []Item      `json:"1"`
	ItemCounts []ItemCount `json:"2"`
}

type ELPResponseWU struct {
	Header ResponseHeader  `json:"header"`
	Body   ELPResponseBody `json:"body"`
}

type ELPResponse struct {
	WU ELPResponseWU `json:"uspwuapi"`
}

func (s *Svc) GetELoadProducts(ctx context.Context, elpReq ELPRequest) (*ELPResponse, error) {
	req, err := s.newParahubRequest(ctx, "eload", "product_list", elpReq)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := s.cl.Post(s.moduleURL("eload", ""), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var elpRes ELPResponse
	if err := json.Unmarshal(body, &elpRes); err != nil {
		return nil, err
	}

	return &elpRes, nil
}

package perahub

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
)

type SELoadRequest struct {
	SessionID      string `json:"session_id"`
	ProductCode    string `json:"ProductCode"`
	Amount         string `json:"Amount"`
	Provider       string `json:"Provider"`
	ProductType    string `json:"ProductType"`
	TargetMobileNo string `json:"Target_mobileno"`
	Password       string `json:"password"`
}

type SELResponseBody struct {
	RRN       string `json:"RRN"`
	TID       string `json:"TID"`
	EPIN      string `json:"EPIN"`
	Timestamp string `json:"timestamp"`
}

type SELResponseWU struct {
	Header ResponseHeader  `json:"header"`
	Body   SELResponseBody `json:"body"`
}

type SELoadResponse struct {
	WU SELResponseWU `json:"uspwuapi"`
}

func (s *Svc) SellLoad(ctx context.Context, selReq SELoadRequest) (*SELoadResponse, error) {
	req, err := s.newParahubRequest(ctx, "eload", "sell", selReq)
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

	var selRes SELoadResponse
	if err := json.Unmarshal(body, &selRes); err != nil {
		return nil, err
	}

	return &selRes, nil
}

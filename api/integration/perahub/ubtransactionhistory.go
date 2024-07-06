package perahub

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
)

type UBTHRequest struct {
	UserCode     string `json:"usercode"`
	DateStart    string `json:"date_start"`
	DateEnd      string `json:"date_end"`
	OperatorCode string `json:"operator_code"`
	LocationCode string `json:"location_code"`
	SessionID    string `json:"session_id"`
}

type Transaction struct {
	ID           string `json:"id"`
	TranDate     string `json:"tranDate"`
	Amount       string `json:"amount"`
	Branch       string `json:"branch"`
	TranDesc     string `json:"tranDesc"`
	Balance      string `json:"balance"`
	TrxnCurrency string `json:"trxnCurrency"`
	PplCode      string `json:"pplcode"`
}

type TranListDetails struct {
	Transaction []Transaction `json:"trxn"`
}

type UBTHResponseBody struct {
	TranList TranListDetails `json:"trxnList"`
}

type UBTHResponseData struct {
	Header ResponseHeader   `json:"header"`
	Body   UBTHResponseBody `json:"body"`
}

type UBTHResponse struct {
	WU UBTHResponseData `json:"uspwuapi"`
}

func (s *Svc) UBTransactionHistory(ctx context.Context, tr UBTHRequest) (*UBTHResponse, error) {
	req, err := s.newParahubRequest(ctx, "Transaction", "get_transaction", tr)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := s.cl.Post(s.moduleURL("Transaction", ""), "aplication/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var tranHisUBRes UBTHResponse
	if err := json.Unmarshal(body, &tranHisUBRes); err != nil {
		return nil, err
	}

	return &tranHisUBRes, nil
}

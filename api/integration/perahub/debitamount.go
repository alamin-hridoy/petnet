package perahub

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
)

type DebitAmountRequest struct {
	Username      string `json:"usercode"`
	Amount        string `json:"amount"`
	Password      string `json:"password"`
	OperationCode string `json:"Operator_code"`
	SessionID     string `json:"Session_id"`
}

type ResponseBody struct {
	Code           string `json:"code"`
	AmountDeducted int    `json:"amount_deducted"`
}

type DebitAmountResponseWU struct {
	Header ResponseHeader `json:"header"`
	Body   ResponseBody   `json:"body"`
}

type DebitAmountResponse struct {
	WU DebitAmountResponseWU `json:"uspwuapi"`
}

func (s *Svc) DebitAmount(ctx context.Context, dr DebitAmountRequest) (*DebitAmountResponse, error) {
	req, err := s.newParahubRequest(ctx, "DebitAmount", "debit_amount", dr)
	if err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := s.cl.Post(s.moduleURL("DebitAmount", ""), "aplication/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var amountRes DebitAmountResponse
	if err := json.Unmarshal(body, &amountRes); err != nil {
		return nil, err
	}

	return &amountRes, nil
}

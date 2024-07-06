package perahub

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"brank.as/petnet/serviceutil/logging"
)

type WUTHRequest struct {
	DateStart string `json:"date_from"`
	DateEnd   string `json:"date_to"`
}

type wuHistory struct {
	Status string          `json:"status"`
	Msg    string          `json:"errmsg"`
	Data   []WUTransaction `json:"data"`
}

type WUTransaction struct {
	TxnDate      string      `json:"transactiondate"`
	MTCN         string      `json:"MTCN"`
	Principal    string      `json:"principalamount"`
	SvcFee       json.Number `json:"servicefee"`
	Currency     string      `json:"currency"`
	TxnType      string      `json:"TransactionType"`
	DateClaimed  string      `json:"DateClaimed"`
	CustomerCode string      `json:"CustomerCode"`
	OrderID      string      `json:"OrderID"`
}

func (s *Svc) WUTransactionHistory(ctx context.Context, tr WUTHRequest) ([]WUTransaction, error) {
	log := logging.FromContext(ctx)

	const mod, modReq = "wurpt", "report"
	const reqMod, reqReq = "wu", "rpt"
	body, err := s.newParahubRequest(ctx, mod, modReq, tr)
	if err != nil {
		return nil, err
	}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	url := s.moduleURL(reqMod, reqReq)
	log.WithField("url", url).WithField("request body", json.RawMessage(reqBody)).Debug("sending")
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	resp, err := s.cl.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer(nil)
	resp.Body = io.NopCloser(io.TeeReader(resp.Body, buf))

	res := wuHistory{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.WithField("payload:", buf.String()).Debug("invalid json")
		return nil, err
	}

	log.WithField("response payload:", json.RawMessage(buf.String())).Debug("parsed")

	return res.Data, nil
}

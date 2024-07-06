package perahub

import (
	"context"
	"encoding/json"
)

type RemitanceConfirmSendMoneyReq struct {
	SendValidateReferenceNumber string `json:"send_validate_reference_number"`
}

type RemitanceConfirmSendMoneyResult struct {
	Phrn string `json:"phrn"`
}

type RemitanceConfirmSendMoneyRes struct {
	Code    int                             `json:"code"`
	Message string                          `json:"message"`
	Result  RemitanceConfirmSendMoneyResult `json:"result"`
}

func (s *Svc) RemitanceConfirmSendMoney(ctx context.Context, req RemitanceConfirmSendMoneyReq) (*RemitanceConfirmSendMoneyRes, error) {
	res, err := s.remitancePost(ctx, s.remitanceURL("send/confirm"), req)
	if err != nil {
		return nil, err
	}

	rb := &RemitanceConfirmSendMoneyRes{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

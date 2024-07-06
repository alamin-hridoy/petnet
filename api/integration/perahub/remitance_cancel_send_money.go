package perahub

import (
	"context"
	"encoding/json"
)

type RemitanceCancelSendMoneyReq struct {
	Phrn        string `json:"phrn"`
	PartnerCode string `json:"partner_code"`
	Remarks     string `json:"remarks"`
}

type RemitanceCancelSendMoneyResult struct {
	Phrn                      string `json:"phrn"`
	CancelSendDate            string `json:"cancel_send_date"`
	CancelSendReferenceNumber string `json:"cancel_send_reference_number"`
}

type RemitanceCancelSendMoneyRes struct {
	Code    int                            `json:"code"`
	Message string                         `json:"message"`
	Result  RemitanceCancelSendMoneyResult `json:"result"`
}

func (s *Svc) RemitanceCancelSendMoney(ctx context.Context, req RemitanceCancelSendMoneyReq) (*RemitanceCancelSendMoneyRes, error) {
	res, err := s.remitancePost(ctx, s.remitanceURL("send/cancel"), req)
	if err != nil {
		return nil, err
	}

	rcsm := &RemitanceCancelSendMoneyRes{}
	if err := json.Unmarshal(res, rcsm); err != nil {
		return nil, err
	}
	return rcsm, nil
}

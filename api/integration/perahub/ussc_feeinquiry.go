package perahub

import (
	"context"
	"encoding/json"
)

const (
	USSCFeeInquiryOK = "OK"
)

type USSCFeeInquiryRequest struct {
	Panalokard string `json:"panalokard"`
	Amount     string `json:"amount"`
	USSCPromo  string `json:"ussc_promo"`
	BranchCode string `json:"branch_code"`
}

type USSCFeeInquiryRespBody struct {
	Code    json.Number          `json:"code"`
	Message string               `json:"message"`
	Result  USSCFeeInquiryResult `json:"result"`
	RemcoID int                  `json:"remco_id"`
}

type USSCFeeInquiryResult struct {
	PnplAmount    string `json:"principal_amount"`
	ServiceCharge string `json:"service_charge"`
	Msg           string `json:"message"`
	Code          string `json:"code"`
	NewScreen     string `json:"new_screen"`
	JournalNo     string `json:"journal_no"`
	ProcessDate   string `json:"process_date"`
	RefNo         string `json:"reference_number"`
	TotAmount     string `json:"total_amount"`
	SendOTP       string `json:"send_otp"`
}

func (s *Svc) USSCFeeInquiry(ctx context.Context, sr USSCFeeInquiryRequest) (*USSCFeeInquiryRespBody, error) {
	res, err := s.postNonex(ctx, s.nonexURL("ussc/fee-inquiry"), sr)
	if err != nil {
		return nil, err
	}

	rb := &USSCFeeInquiryRespBody{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}

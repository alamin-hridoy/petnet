package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWISEQuoteInquiry(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	s, m := newTestSvc(t, st)
	wantReq := WISEQuoteInquiryReq{
		SourceCurrency: "PHP",
		TargetCurrency: "GBP",
		SourceAmount:   "1500",
	}
	wantResp := &WISEQuoteInquiryResp{
		SourceCurrency: "PHP",
		TargetCurrency: "GBP",
		SourceAmount:   "1500",
		TargetAmount:   "70.21",
		FeeBreakdown: WISEFeeBreakdown{
			Transferwise: "79.66",
			PayIn:        "0",
			Discount:     "0",
			Total:        "79.66",
			PriceSetID:   "132",
			Partner:      "0",
		},
		TotalFee:       "2.49",
		TransferAmount: "97.51",
		PayOut:         "BANK_TRANSFER",
		Rate:           "0.0148552",
	}
	gotResp, err := s.WISEQuoteInquiry(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}

	var gotReq WISEQuoteInquiryReq
	if err := json.Unmarshal(m.GetMockRequest(), &gotReq); err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(wantReq, gotReq) {
		t.Error(cmp.Diff(wantReq, gotReq))
	}
	if !cmp.Equal(wantResp, gotResp) {
		t.Error(cmp.Diff(wantResp, gotResp))
	}
}

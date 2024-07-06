package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWISEProceedTransfer(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	s, m := newTestSvc(t, st)
	wantReq := WISEProceedTransferReq{
		Email:       "sender@brankas.com",
		RecipientID: "12345",
		Details: WISEPCDetails{
			Reference: "54321",
		},
	}
	wantResp := &WISEProceedTransferResp{
		TransferID: "12345",
		Details: WISEPCDetails{
			Reference: "54321",
		},
		CustomerTxnID:  "aecd179d",
		RecipientID:    "447769582",
		Status:         "incoming_payment_waiting",
		SourceCurrency: "PHP",
		TargetCurrency: "GBP",
		SourceAmount:   "1500",
		DateCreated:    "2021-03-05T05:38:36.677Z",
	}
	gotResp, err := s.WISEProceedTransfer(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}

	var gotReq WISEProceedTransferReq
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

package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWISEGetRecipients(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	s, m := newTestSvc(t, st)
	wantReq := WISEGetRecipientsReq{
		Email:    "user@email.com",
		Currency: "GBP",
	}
	wantResp := &WISEGetRecipientsResp{
		Recipients: []WISERecipient{
			{
				RecipientID: "148142936",
				Details: WISEGRDetails{
					AccountNumber:    "28821822",
					SortCode:         "231470",
					HashedByLooseAlg: "d641e415d0503d966ff4a5a7d246ba9511430a258e8a993b53059defb256c448",
				},
				AccountSummary:     "(23-14-70) 28821822",
				LongAccountSummary: "GBP account ending in 1822",
				DisplayFields: []WISEDisplayField{
					{
						Label: "UK Sort code",
						Value: "23-14-70",
					},
					{
						Label: "Account number",
						Value: "28821822",
					},
				},
				FullName: "Brankas Receiver",
				Currency: "GBP",
				Country:  "GB",
			},
		},
	}
	gotResp, err := s.WISEGetRecipients(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}

	var gotReq WISEGetRecipientsReq
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

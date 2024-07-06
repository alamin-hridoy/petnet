package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWISEDeleteRecipient(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	s, m := newTestSvc(t, st)
	wantReq := WISEDeleteRecipientReq{
		RecipientID: "12345",
		Email:       "user@email.com",
	}
	wantResp := &WISEDeleteRecipientResp{
		Msg: "Successfully Deleted Recipient Account!",
	}
	gotResp, err := s.WISEDeleteRecipient(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}

	var gotReq WISEDeleteRecipientReq
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

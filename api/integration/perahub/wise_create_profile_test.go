package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWISECreateProfile(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	s, m := newTestSvc(t, st)
	wantReq := WISECreateProfileReq{
		Email: "user@email.com",
	}
	wantResp := &WISECreateProfileResp{
		Msg: "Success! User Profile Created.",
	}
	gotResp, err := s.WISECreateProfile(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}

	var gotReq WISECreateProfileReq
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

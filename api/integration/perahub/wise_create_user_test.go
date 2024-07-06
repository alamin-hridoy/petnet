package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWISECreateUser(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	s, m := newTestSvc(t, st)
	wantReq := WISECreateUserReq{
		Email: "user@email.com",
	}
	wantResp := &WISECreateUserResp{
		Msg: "Success! User Account Created.",
	}
	gotResp, err := s.WISECreateUser(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}

	var gotReq WISECreateUserReq
	if err := json.Unmarshal(m.GetMockRequest(), &gotReq); err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(wantReq, gotReq) {
		t.Error(cmp.Diff(wantReq, gotReq))
	}
	if !cmp.Equal(wantResp, gotResp) {
		t.Error(cmp.Diff(wantResp, gotResp))
	}

	m.SetConflictError()
	wantResp.Error = "Conflict"
	wantResp.Msg = ""
	if _, err = s.WISECreateUser(context.Background(), wantReq); err == nil {
		t.Fatal("got nil, want error")
	}
}

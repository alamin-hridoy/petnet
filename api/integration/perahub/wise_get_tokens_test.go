package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWISEGetTokens(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	s, m := newTestSvc(t, st)
	wantReq := WISEGetTokensReq{
		ClientID:     "clientID",
		ClientSecret: "clientSecret",
	}
	wantResp := &WISEGetTokensResp{
		AccessToken:  "token",
		RefreshToken: "refr-token",
		ExpiresIn:    "43200",
		TokenType:    "Bearer",
	}
	gotResp, err := s.WISEGetTokens(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}

	var gotReq WISEGetTokensReq
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

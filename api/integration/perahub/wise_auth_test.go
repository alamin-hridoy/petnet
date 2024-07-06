package perahub

import (
	"context"
	"testing"
	"time"

	"brank.as/petnet/api/core/static"
	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/metadata"
)

func TestWISEAuth(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	s, m := newTestSvc(t, st)

	nmd := metautils.NiceMD(metadata.Pairs("partner", static.WISECode))
	ctx := nmd.ToIncoming(context.Background())

	t.Run("No Auth Cache", func(t *testing.T) {
		if _, err := s.WISECreateUser(ctx, WISECreateUserReq{
			Email: "user@email.com",
		}); err != nil {
			t.Fatal(err)
		}
		wantReqOrder := []string{"auth", "create"}
		if !cmp.Equal(wantReqOrder, m.reqOrder) {
			t.Error(cmp.Diff(wantReqOrder, m.reqOrder))
		}
		if h := m.httpHeaders.Get("Authorization"); h != "Bearer token" {
			t.Error("want Authorization header Bearer token, got: ", h)
		}
		m.ResetReqOrder()
	})

	t.Run("Has Auth Cache", func(t *testing.T) {
		if _, err := s.WISECreateUser(ctx, WISECreateUserReq{
			Email: "user@email.com",
		}); err != nil {
			t.Fatal(err)
		}
		wantReqOrder := []string{"create"}
		if !cmp.Equal(wantReqOrder, m.reqOrder) {
			t.Error(cmp.Diff(wantReqOrder, m.reqOrder))
		}
		if h := m.httpHeaders.Get("Authorization"); h != "Bearer token" {
			t.Error("want Authorization header Bearer token, got: ", h)
		}
		m.ResetReqOrder()
	})

	t.Run("Create Auth Error", func(t *testing.T) {
		m.SetAuthError()
		if _, err := s.WISECreateUser(ctx, WISECreateUserReq{
			Email: "user@email.com",
		}); err != nil {
			t.Fatal(err)
		}
		wantReqOrder := []string{"create", "auth", "create"}
		if !cmp.Equal(wantReqOrder, m.reqOrder) {
			t.Error(cmp.Diff(wantReqOrder, m.reqOrder))
		}
		if h := m.httpHeaders.Get("Authorization"); h != "Bearer token" {
			t.Error("want Authorization header Bearer token, got: ", h)
		}
		m.ResetReqOrder()
	})

	t.Run("Token Expired", func(t *testing.T) {
		setPartnerAuth(static.WISECode, oauth{
			AccessToken:  "token",
			RefreshToken: "refr-token",
			Expiry:       time.Now(),
		})
		if _, err := s.WISECreateUser(ctx, WISECreateUserReq{
			Email: "user@email.com",
		}); err != nil {
			t.Fatal(err)
		}
		wantReqOrder := []string{"auth", "create"}
		if !cmp.Equal(wantReqOrder, m.reqOrder) {
			t.Error(cmp.Diff(wantReqOrder, m.reqOrder))
		}
		if h := m.httpHeaders.Get("Authorization"); h != "Bearer token" {
			t.Error("want Authorization header Bearer token, got: ", h)
		}
	})
}

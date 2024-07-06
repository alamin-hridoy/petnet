package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"brank.as/rbac/svcutil/random"
	"brank.as/rbac/usermgm/storage"
)

func TestOAuth2(t *testing.T) {
	t.Parallel()
	ts := newTestStorage(t)
	opts := cmp.Options{cmp.FilterValues(func(a, b time.Time) bool { return true }, cmp.Ignore())}
	id, err := random.String(16)
	if err != nil {
		t.Fatal(err)
	}
	cl := storage.OAuthClient{
		OrgID:        uuid.New().String(),
		ClientID:     id,
		ClientName:   "test-name",
		CreateUserID: uuid.New().String(),
		Environment:  "sandbox",
	}
	t.Run("Create", func(t *testing.T) {
		got, err := ts.CreateOauthClient(context.TODO(), cl)
		if err != nil {
			t.Fatalf("CreateOauthClient() = got error %v, want nil", err)
		}
		if !cmp.Equal(cl, *got, opts...) {
			t.Error(cmp.Diff(cl, *got, opts...))
		}
		cl.Created = got.Created
		cl.Updated = got.Updated
	})
	t.Run("Update", func(t *testing.T) {
		cl.UpdateUserID = uuid.New().String()
		got, err := ts.UpdateOauthClient(context.TODO(), cl)
		if err != nil {
			t.Fatalf("UpdateOauthClient() = got error %v, want nil", err)
		}
		if !cmp.Equal(cl, *got, opts...) {
			t.Error(cmp.Diff(cl, *got, opts...))
		}
		if got.Updated.Equal(cl.Updated) {
			t.Error(cmp.Diff(cl, got))
		}
	})
	t.Run("Get", func(t *testing.T) {
		got, err := ts.GetOauthClientByID(context.TODO(), cl.ClientID)
		if err != nil {
			t.Fatalf("GetOauthClientByID() = got error %v, want nil", err)
		}
		if !cmp.Equal(cl, *got, opts...) {
			t.Error(cmp.Diff(cl, *got, opts...))
		}
		if got.Updated.Equal(cl.Updated) {
			t.Error(cmp.Diff(cl, got))
		}
	})
	t.Run("Get By OrgID", func(t *testing.T) {
		got, err := ts.GetOauthClientByOrgID(context.TODO(), cl.OrgID, false)
		if err != nil {
			t.Fatalf("GetOauthClientByOrgID() = got error %v, want nil", err)
		}
		gots := got[0]
		if !cmp.Equal(cl, gots, opts...) {
			t.Error(cmp.Diff(cl, gots, opts...))
		}
		if gots.Updated.Equal(cl.Updated) {
			t.Error(cmp.Diff(cl, gots))
		}
	})
	t.Run("Delete", func(t *testing.T) {
		cl.DeleteUserID = uuid.New().String()
		tm, err := ts.DeleteOauthClient(context.TODO(), cl)
		if err != nil {
			t.Error(err)
		}
		if tm.After(time.Now()) || tm.IsZero() {
			t.Error("invalid delete timestamp", tm)
		}
	})
	t.Run("Get Oauth By OrgID Deleted", func(t *testing.T) {
		got, err := ts.GetOauthClientByOrgID(context.TODO(), cl.OrgID, true)
		if err != nil {
			t.Fatalf("GetOauthClientByOrgID() = got error %v, want nil", err)
		}
		gots := got[0]
		if len(got) != 1 {
			t.Error("expected 1 deleted oauth clientId, got", len(got))
		}
		if gots.Updated.Equal(cl.Updated) {
			t.Error(cmp.Diff(cl, gots))
		}
	})
	t.Run("Get Oauth By OrgID Exclude Deleted", func(t *testing.T) {
		got, err := ts.GetOauthClientByOrgID(context.TODO(), cl.OrgID, false)
		if err != nil {
			t.Fatalf("GetOauthClientByOrgID() = got error %v, want nil", err)
		}
		if len(got) != 0 {
			t.Error("expected 0 oauth clientId, got", len(got))
		}
	})
}

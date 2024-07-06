package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"brank.as/rbac/usermgm/storage"
)

func TestUpsertScope(t *testing.T) {
	t.Parallel()
	ts := newTestStorage(t)

	ignUpd := cmp.FilterValues(func(a, b time.Time) bool { return true }, cmp.Ignore())
	t.Run("new", func(t *testing.T) {
		sc := storage.Scope{
			ID:    uuid.NewString(),
			Name:  "testnew",
			Group: "firstgroup",
			Desc:  "new scope insert",
		}
		got, err := ts.UpsertScope(context.TODO(), sc)
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(&sc, got, ignUpd) {
			t.Error(cmp.Diff(&sc, got))
		}
	})

	t.Run("update", func(t *testing.T) {
		sc := storage.Scope{
			ID:    uuid.NewString(),
			Name:  "testupdate",
			Group: "updategroup",
			Desc:  "update scope insert",
		}
		ctx := context.TODO()
		got, err := ts.UpsertScope(ctx, sc)
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(&sc, got, ignUpd) {
			t.Error(cmp.Diff(&sc, got))
		}
		sc.Desc = "updated scope desc"
		got, err = ts.UpsertScope(ctx, sc)
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(&sc, got, ignUpd) {
			t.Error(cmp.Diff(&sc, got))
		}

		list, err := ts.GetScopes(ctx, []string{sc.ID})
		if err != nil {
			t.Error(err)
		}
		if !cmp.Equal([]storage.Scope{sc}, list, ignUpd) {
			t.Error(cmp.Diff([]storage.Scope{sc}, list))
		}
	})

	t.Run("group rename", func(t *testing.T) {
		sc := storage.Scope{
			ID:    uuid.NewString(),
			Name:  "testgroupupdate",
			Group: "updategroupname",
			Desc:  "update group insert",
		}
		ctx := context.TODO()
		got, err := ts.UpsertScope(ctx, sc)
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(&sc, got, ignUpd) {
			t.Error(cmp.Diff(&sc, got))
		}
		grInit, err := ts.GetScopeGroups(ctx, []string{sc.Group})
		if err != nil {
			t.Error(err)
		}
		if want := []storage.ScopeGroup{{Name: sc.Group}}; !cmp.Equal(want, grInit, ignUpd) {
			t.Error(cmp.Diff(want, got))
		}

		gr := storage.ScopeGroup{
			Name: sc.Group,
			Desc: "New group description",
		}
		gotGr, err := ts.UpdateGroup(ctx, gr)
		if err != nil {
			t.Error(err)
		}
		if !cmp.Equal(&gr, gotGr, ignUpd) {
			t.Error(cmp.Diff(&gr, gotGr))
		}
	})
}

func TestGrantConsent(t *testing.T) {
	t.Parallel()
	ts := newTestStorage(t)

	ignUpd := cmp.FilterValues(func(a, b time.Time) bool { return true }, cmp.Ignore())
	grantID := cmp.FilterPath(func(p cmp.Path) bool { return p.Last().String() == ".ID" },
		cmp.Comparer(func(a, b string) bool {
			_, aerr := uuid.Parse(a)
			_, berr := uuid.Parse(b)
			return aerr == nil || berr == nil
		}))
	t.Run("record grant", func(t *testing.T) {
		gr := storage.ConsentGrant{
			UserID:   uuid.NewString(),
			ClientID: randomString(20),
			OwnerID:  uuid.NewString(),
			Scopes: []string{
				"openid", "offline_access", "https://service.bnk.to/object.operation",
			},
		}
		ctx := context.TODO()
		g, err := ts.RecordGrant(ctx, gr)
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(&gr, g, ignUpd, grantID) {
			t.Error(cmp.Diff(&gr, g, ignUpd, grantID))
		}
	})
}

package postgres_test

import (
	"context"
	"database/sql"
	"sort"
	"testing"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sirupsen/logrus"
)

func TestUserProfile(t *testing.T) {
	ts := newTestStorage(t)

	oid := "10000000-0000-0000-0000-000000000000"
	oid2 := "20000000-0000-0000-0000-000000000000"
	uid := "11000000-0000-0000-0000-000000000000"
	uid2 := "22000000-0000-0000-0000-000000000000"
	want := []*storage.UserProfile{
		{
			OrgID:          oid,
			UserID:         uid,
			ProfilePicture: "test",
			Email:          "email@example.com",
		},
		{
			OrgID:          oid2,
			UserID:         uid2,
			ProfilePicture: "test",
			Email:          "email2@example.com",
		},
	}

	logr := logging.NewLogger(nil)
	logr.SetFormatter(&logrus.JSONFormatter{})
	ctx := logging.WithLogger(context.TODO(), logr)

	pid, err := ts.CreateUserProfile(ctx, want[0])
	if err != nil {
		t.Fatal(err)
	}
	if pid == "" {
		t.Error("id should not be empty")
	}

	_, err = ts.CreateUserProfile(ctx, want[0])
	if err != storage.Conflict {
		t.Error("error should be conflict")
	}
	pid2, err := ts.CreateUserProfile(ctx, want[1])
	if err != nil {
		t.Fatal(err)
	}
	if pid2 == "" {
		t.Error("id should not be empty")
	}

	got, err := ts.GetUserProfile(ctx, want[0].UserID)
	if err != nil {
		t.Fatal(err)
	}

	tOps := []cmp.Option{
		cmpopts.IgnoreFields(storage.UserProfile{}, "ID", "Created", "Updated", "Deleted"),
	}
	if !cmp.Equal(want[0], got, tOps...) {
		t.Error("(-want +got): ", cmp.Diff(&want[0], got, tOps...))
	}
	if got.ID == "" {
		t.Error("id should not be empty")
	}
	if got.Created.IsZero() || got.Updated.IsZero() {
		t.Error("created and updated shouldn't be empty")
	}
	if got.Deleted.Valid {
		t.Error("deleted should be null")
	}

	gotlist, err := ts.GetUserProfiles(ctx, oid)
	if err != nil {
		t.Fatal(err)
	}

	sort.Slice(want, func(i, j int) bool {
		return want[i].OrgID < want[j].OrgID
	})
	sort.Slice(gotlist, func(i, j int) bool {
		return gotlist[i].OrgID < gotlist[j].OrgID
	})
	for i, pf := range gotlist {
		if !cmp.Equal(want[i], &pf, tOps...) {
			t.Error("(-want +got): ", cmp.Diff(want[i], pf, tOps...))
		}
		if pf.ID == "" {
			t.Error("id should not be empty")
		}
		if pf.Created.IsZero() || pf.Updated.IsZero() {
			t.Error("created and updated shouldn't be empty")
		}
		if pf.Deleted.Valid {
			t.Error("deleted should be null")
		}
	}

	wantup := &storage.UserProfile{
		ID:             pid,
		UserID:         uid,
		OrgID:          oid,
		ProfilePicture: "test",
		Email:          want[0].Email,
		Deleted:        sql.NullTime{},
	}
	upid, err := ts.UpdateUserProfile(ctx, wantup)
	if err != nil {
		t.Fatal(err)
	}
	if upid != pid {
		t.Error("id mismatch")
	}

	got, err = ts.GetUserProfile(ctx, uid)
	if err != nil {
		t.Fatal(err)
	}
	wantupp := &storage.UpdateOrgProfileOrgIDUserID{
		OldOrgID: oid,
		NewOrgID: "40000000-0000-0000-0000-000000000000",
		UserID:   "41000000-0000-0000-0000-000000000000",
	}

	_, err = ts.UpdateUserProfileByOrgID(ctx, wantupp)
	if err != nil {
		t.Fatal(err)
	}
	tOps = []cmp.Option{
		cmpopts.IgnoreFields(storage.UserProfile{}, "ID", "Created", "Updated"),
	}
	if !cmp.Equal(wantup, got, tOps...) {
		t.Error("(-want +got): ", cmp.Diff(wantup, got, tOps...))
	}
}

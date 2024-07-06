package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sirupsen/logrus"
)

func TestSession(t *testing.T) {
	ts := newTestStorage(t)

	uid := "11000000-0000-0000-0000-000000000000"
	want := &storage.Session{
		UserID: uid,
		Expiry: sql.NullTime{Time: time.Unix(1414141414, 0), Valid: true},
	}

	logr := logging.NewLogger(nil)
	logr.SetFormatter(&logrus.JSONFormatter{})
	ctx := logging.WithLogger(context.TODO(), logr)

	sid, err := ts.UpsertSession(ctx, want)
	if err != nil {
		t.Fatal(err)
	}
	if sid == "" {
		t.Error("id should not be empty")
	}

	got, err := ts.GetSession(ctx, uid)
	if err != nil {
		t.Fatal(err)
	}
	tOps := []cmp.Option{
		cmpopts.IgnoreFields(storage.UserProfile{}, "ID", "Created", "Updated", "Deleted"),
	}
	if !cmp.Equal(want, got, tOps...) {
		t.Error("(-want +got): ", cmp.Diff(&want, got, tOps...))
	}
	if got.ID == "" {
		t.Error("session id should not be empty")
	}
	if got.Created.IsZero() || got.Updated.IsZero() {
		t.Error("created and updated shouldn't be empty")
	}
	if got.Deleted.Valid {
		t.Error("deleted should be null")
	}

	wantup := &storage.Session{
		ID:      sid,
		UserID:  uid,
		Expiry:  sql.NullTime{Time: time.Unix(1616161616, 0), Valid: true},
		Deleted: sql.NullTime{Time: time.Unix(1616161616, 0), Valid: true},
	}
	usid, err := ts.UpsertSession(ctx, wantup)
	if err != nil {
		t.Fatal(err)
	}
	if usid != sid {
		t.Error("id mismatch")
	}

	got, err = ts.GetSession(ctx, uid)
	if err != nil {
		t.Fatal("GetUserProfile: ", err)
	}
	tOps = []cmp.Option{
		cmpopts.IgnoreFields(storage.UserProfile{}, "ID", "Created", "Updated"),
	}
	if !cmp.Equal(wantup, got, tOps...) {
		t.Error("GetUserProfile (-want +got): ", cmp.Diff(wantup, got, tOps...))
	}
}

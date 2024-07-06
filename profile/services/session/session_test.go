package session

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus/hooks/test"

	"brank.as/petnet/profile/core/session"
	"brank.as/petnet/profile/storage/postgres"

	spb "brank.as/petnet/gunk/v1/session"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func TestSessionProfile(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	want := &spb.SetSessionExpiryRequest{
		IDType: spb.IDType_USERID,
		ID:     uuid.New().String(),
		Expiry: tspb.New(time.Unix(1414141414, 0)),
	}
	test.NewNullLogger()
	st, cleanup := postgres.NewTestStorage(os.Getenv("DATABASE_CONNECTION"), filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(cleanup)

	h := New(session.New(st))
	if _, err := h.SetSessionExpiry(ctx, want); err != nil {
		t.Fatal(err)
	}

	sess, err := st.GetSession(ctx, want.ID)
	if err != nil {
		t.Error(err)
	}
	if sess.Expiry.Time != want.Expiry.AsTime() {
		t.Error("expiry doesn't match")
	}
	if sess.UserID != want.ID {
		t.Error("user ids don't match")
	}

	if _, err := h.ExpireSession(ctx,
		&spb.ExpireSessionRequest{
			IDType: spb.IDType_USERID,
			ID:     want.ID,
		}); err != nil {
		t.Fatal(err)
	}
	sess, err = st.GetSession(ctx, want.ID)
	if err != nil {
		t.Error(err)
	}
	ts := sql.NullTime{}
	if sess.Expiry.Time != ts.Time {
		t.Error("expiry should be empty")
	}
}

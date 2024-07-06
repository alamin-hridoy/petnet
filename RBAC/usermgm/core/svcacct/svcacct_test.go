package svcacct

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"brank.as/rbac/usermgm/storage"
	"brank.as/rbac/usermgm/storage/postgres"
)

func TestAPIKey(t *testing.T) {
	conn := os.Getenv("DATABASE_CONNECTION")
	switch "" {
	case conn:
		t.Skip("missing env 'DATABASE_CONNECTION'")
	}
	t.Parallel()
	tmOpt := cmp.FilterValues(func(a, b time.Time) bool { return true }, cmp.Ignore())
	db, clean := postgres.NewTestStorage(conn, filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(clean)

	s := &Svc{store: db}

	ctx := context.Background()
	sa := storage.SvcAccount{
		AuthType:     storage.APIKey,
		Environment:  "sandbox",
		ClientName:   "test-client",
		ClientID:     "",
		Challenge:    "",
		CreateUserID: "someuser",
	}

	oid, err := db.CreateOrg(ctx, storage.Organization{
		ID:      sa.OrgID,
		OrgName: "test apikey",
		Active:  true,
	})
	sa.OrgID = oid

	id, key, err := s.CreateSvcAccount(ctx, sa)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(key, id) {
		t.Error("invalid prefix:", id, key)
	}
	sa.ClientID = id

	got, err := s.ValidateSvcAccount(ctx, key)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(&sa, got, tmOpt) {
		t.Error(cmp.Diff(&sa, got, tmOpt))
	}
}

package postgres

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

const dbConnEnv = "DATABASE_CONNECTION"

func TestMain(m *testing.M) {
	connStr := os.Getenv(dbConnEnv)
	if connStr == "" {
		log.Printf("%s is not set, skipping", dbConnEnv)
		return
	}
	exitCode := m.Run()
	os.Exit(exitCode)
}

func newTestStorage(tb testing.TB) *Storage {
	if testing.Short() {
		tb.Skip("skipping tests that use postgres on -short")
	}
	connStr := os.Getenv(dbConnEnv)
	st, cl := NewTestStorage(connStr, filepath.Join("..", "..", "migrations", "sql"))
	tb.Cleanup(cl)
	return st
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/pressly/goose"

	"brank.as/petnet/serviceutil/storage/postgres"

	// Postgres storage driver
	_ "github.com/lib/pq"
)

const (
	// https://www.postgresql.org/docs/9.6/errcodes-appendix.html
	pqUnique   = "23505"
	pqNotFound = "42703"
)

// Storage provides a wrapper around an sql database and provides
// required methods for interacting with the database
type Storage struct {
	db *sqlx.DB
}

// NewStorage returns a new Storage from the provides psql databse string
func NewStorage(dbstring string) (*Storage, error) {
	db, err := sqlx.Connect("postgres", dbstring)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to postgres '%s'", dbstring)
	}
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Hour)
	return &Storage{db: db}, nil
}

func NewStorageDB(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

// RunMigration runs the migrations in the dir using the goose package
func (s *Storage) RunMigration(dir string) error {
	return goose.Run("up", s.db.DB, dir)
}

// NewTestStorage returns a Storage that uses an isolated database for testing purposes
// and a teardown function
func NewTestStorage(dbstring string, migrationDir string) (*Storage, func()) {
	db, teardown := postgres.MustNewDevelopmentDB(dbstring, migrationDir)
	db.SetMaxOpenConns(5)
	db.SetConnMaxLifetime(time.Hour)

	return NewStorageDB(db), teardown
}

type pgTx struct{}

type tx struct {
	*sqlx.Tx
	committed *bool
}

func (s *Storage) NewTransacton(ctx context.Context) (context.Context, error) {
	t, err := s.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, pgTx{}, &tx{
		Tx:        t,
		committed: new(bool),
	}), nil
}

func (s *Storage) Commit(ctx context.Context) error {
	t := getTx(ctx)
	if t == nil {
		return fmt.Errorf("not a transaction")
	}
	if *t.committed {
		return nil
	}
	if err := t.Commit(); err != nil {
		return err
	}
	*t.committed = true
	return nil
}

func (s *Storage) Rollback(ctx context.Context) error {
	t := getTx(ctx)
	if t == nil {
		return fmt.Errorf("not a transaction")
	}
	if *t.committed {
		return nil
	}
	return t.Rollback()
}

func getTx(ctx context.Context) *tx {
	if t, ok := ctx.Value(pgTx{}).(*tx); ok && !*t.committed {
		return t
	}
	return nil
}

// PrepareNamed prepares a named query in the current transaction (if begun) or with the db.
func (s *Storage) prepareNamed(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	if tx := getTx(ctx); tx != nil {
		return tx.PrepareNamedContext(ctx, query)
	}
	return s.db.PrepareNamedContext(ctx, query)
}

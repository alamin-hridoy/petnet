package postgres

import (
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
	"github.com/spf13/viper"

	"brank.as/petnet/serviceutil/storage/postgres"
)

const (
	// https://www.postgresql.org/docs/9.6/errcodes-appendix.html
	pqUnique   = "23505"
	pqNotFound = "42703"
)

type Storage struct {
	db *sqlx.DB
}

func NewDB(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

func New(config *viper.Viper) (*Storage, error) {
	db, err := postgres.Connectx(config)
	if err != nil {
		return nil, err
	}
	return NewDB(db), nil
}

func NewStorageDB(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

// NewTestStorage returns a Storage that uses an isolated database for testing purposes
// and a teardown function
func NewTestStorage(dbstring string, migrationDir string) (*Storage, func()) {
	db, teardown := postgres.MustNewDevelopmentDB(dbstring, migrationDir)
	db.SetMaxOpenConns(5)
	db.SetConnMaxLifetime(time.Hour)

	return NewStorageDB(db), teardown
}

// RunMigration runs the migrations in the dir using the goose package
func (s *Storage) RunMigration(dir string) error {
	return goose.Run("up", s.db.DB, dir)
}

// stringToSlice is used for format string to slice
func stringToSlice(v string) []string {
	exc := []string{}
	excS := []string{}
	if v == "" {
		return excS
	}
	if v != "" {
		exc = strings.Split(v, ",")
	}
	if len(exc) > 0 {
		excS = append(excS, exc...)
	}
	return excS
}

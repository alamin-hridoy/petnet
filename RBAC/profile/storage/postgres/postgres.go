package postgres

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/pressly/goose"

	"brank.as/rbac/serviceutil/storage/postgres"

	// Postgres storage driver
	_ "github.com/lib/pq"
)

const (
	// https://www.postgresql.org/docs/9.6/errcodes-appendix.html
	pqUnique = "23505"

	userEmailDup  = "user_account_email_key"
	usernameDup   = "user_account_username_key"
	orgInvCodeDup = "org_invitation_code_key"
	usrInvCodeDup = "user_account_invite_code_key"
)

// validationField is used to determine validation rules for
// storage objects
type validationField struct {
	field, value string
	maxLength    int
	required     bool
}

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

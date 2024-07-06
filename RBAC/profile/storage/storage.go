package storage

import (
	"database/sql"
	"errors"
	"time"
)

var (
	// NotFound is returned when the requested resource does not exist.
	NotFound = errors.New("not found")
	// Conflict is returned when trying to create the same resource twice.
	Conflict = errors.New("conflict")
	// EmailExists is returned when signup email already exists in storage.
	EmailExists = errors.New("email already exists")
)

type User struct {
	ID      string       `db:"id"`
	OrgID   string       `db:"org_id"`
	Created time.Time    `db:"created"`
	Updated time.Time    `db:"updated"`
	Deleted sql.NullTime `db:"deleted"`
}

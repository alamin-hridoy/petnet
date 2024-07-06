package storage

import (
	"database/sql"
	"time"
)

type Permission struct {
	ID          string `db:"id"`
	OrgID       string `db:"org_id"`
	SvcPermID   string `db:"service_permission_id"`
	Name        string `db:"permission_name"`
	Description string `db:"description"`

	CreateUID string         `db:"create_user_id"`
	Created   time.Time      `db:"created"`
	Updated   time.Time      `db:"updated"`
	DeleteUID sql.NullString `db:"delete_user_id"`
	Delete    sql.NullTime   `db:"deleted"`
}

type Role struct {
	ID          string `db:"id"`
	OrgID       string `db:"org_id"`
	Name        string `db:"role_name"`
	Description string `db:"description"`

	CreateUID  string         `db:"create_user_id"`
	Created    time.Time      `db:"created"`
	Updated    time.Time      `db:"updated"`
	DeleteUID  sql.NullString `db:"delete_user_id"`
	Delete     sql.NullTime   `db:"deleted"`
	UpdatedUID string         `db:"updateduid"`
	Count      int
}

type Scope struct {
	ID      string    `db:"id"`
	Name    string    `db:"name"`
	Group   string    `db:"group_name"`
	Desc    string    `db:"description"`
	Updated time.Time `db:"updated"`
}

type ScopeGroup struct {
	Name    string    `db:"name"`
	Desc    string    `db:"description"`
	Updated time.Time `db:"updated"`
}

type ConsentGrant struct {
	ID        string    `db:"grant_id"`
	UserID    string    `db:"user_id"`
	ClientID  string    `db:"client_id"`
	OwnerID   string    `db:"owner_id"`
	Scopes    []string  `db:"scopes"`
	Timestamp time.Time `db:"timestamp"`
}

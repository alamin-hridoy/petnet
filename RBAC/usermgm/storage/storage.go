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
	// UsernameExists is returned when the username already exists in storage.
	UsernameExists = errors.New("username already exists")
	// EmailExists is returned when signup email already exists in storage.
	EmailExists = errors.New("email already exists")
	// InvCodeExists is returned when invitation code already exists in storage.
	InvCodeExists = errors.New("invitation code already exists")
)

type Organization struct {
	ID           string       `db:"id"`
	SysRole      string       `db:"sys_role"`
	OrgName      string       `db:"org_name"`
	ContactEmail string       `db:"contact_email"`
	ContactPhone string       `db:"contact_phone"`
	Active       bool         `db:"active"`
	MFALogin     sql.NullBool `db:"mfa_login"`
	Created      time.Time    `db:"created"`
	Updated      time.Time    `db:"updated"`
	Deleted      sql.NullTime `db:"deleted"`
}

type Credential struct {
	ID            string       `db:"id"`
	OrgID         string       `db:"org_id"`
	Username      string       `db:"username"`
	Password      string       `db:"password"`
	PreferredMFA  string       `db:"preferred_mfa"`
	MFALogin      bool         `db:"mfa_login"`
	Deleted       sql.NullTime `db:"deleted"`
	Locked        sql.NullTime `db:"locked"`
	LastLogin     sql.NullTime `db:"last_login"`
	LastFailed    sql.NullTime `db:"last_failed"`
	ResetRequired sql.NullTime `db:"reset_required"`
	FailCount     int          `db:"fail_count"`
	EmailVerified bool         `db:"email_verified"`
}

type User struct {
	ID            string       `db:"id"`
	OrgID         string       `db:"org_id"`
	Username      string       `db:"username"`
	FirstName     string       `db:"first_name"`
	LastName      string       `db:"last_name"`
	Email         string       `db:"email"`
	EmailVerified bool         `db:"email_verified"`
	InviteSender  string       `db:"invite_sender"`
	InviteStatus  InviteStatus `db:"invite_status"`
	InviteCode    string       `db:"invite_code"`
	InviteExpiry  time.Time    `db:"invite_expiry"`
	PreferredMFA  string       `db:"preferred_mfa"`
	MFALogin      bool         `db:"mfa_login"`
	Created       time.Time    `db:"created"`
	Updated       time.Time    `db:"updated"`
	Deleted       sql.NullTime `db:"deleted"`
	Locked        sql.NullTime `db:"locked"`
	ResetRequired sql.NullTime `db:"reset_required"`
	LastLogin     sql.NullTime `db:"last_login"`
	LastFailed    sql.NullTime `db:"last_failed"`
	FailCount     int          `db:"fail_count"`
	Count         int
}

type MFAType = string

const (
	Pass     string = "PASS"
	TOTP     string = "TOTP"
	PINCode  string = "CODE"
	SMS      string = "SMS"
	Recovery string = "RECOVERY"
	EMail    string = "EMAIL"
)

type MFA struct {
	ID        string       `db:"id"`
	UserID    string       `db:"user_id"`
	MFAType   string       `db:"mfa_type"`
	Token     string       `db:"token"`
	Active    bool         `db:"active"`
	Created   time.Time    `db:"created"`
	Updated   time.Time    `db:"updated"`
	Deadline  time.Time    `db:"deadline"`
	Confirmed sql.NullTime `db:"confirmed"`
	Revoked   sql.NullTime `db:"revoked"`
}

type MFAEvent struct {
	EventID    string       `db:"event_id"`
	UserID     string       `db:"user_id"`
	MFAID      string       `db:"mfa_id"`
	MFAType    string       `db:"mfa_type"`
	Active     bool         `db:"active"`
	Expired    bool         `db:"expired"`
	Validation bool         `db:"validation"`
	Token      string       `db:"token"`
	Desc       string       `db:"description"`
	Attempt    int          `db:"attempt"`
	Initiated  time.Time    `db:"initiated"`
	Deadline   time.Time    `db:"deadline"`
	Confirmed  sql.NullTime `db:"confirmed"`
}

type PasswordReset struct {
	ID      string    `db:"id"`
	UserID  string    `db:"user_id"`
	Expiry  time.Time `db:"expiry"`
	Created time.Time `db:"created"`
}

type InviteStatus = string

const (
	Invited    InviteStatus = "Invited"
	InviteSent InviteStatus = "Invite Sent"
	Expired    InviteStatus = "Expired"
	InProgress InviteStatus = "In-Progress"
	Revoked    InviteStatus = "Revoked"
	Approved   InviteStatus = "Approved"
)

var (
	ValidUserStatus = []InviteStatus{Invited, InviteSent, Expired, Revoked}
	ValidOrgStatus  = []InviteStatus{Invited, InviteSent, Expired, InProgress, Revoked, Approved}
)

type AuthType string

const (
	OAuth2 AuthType = "oauth"
	APIKey AuthType = "apikey"
)

type SvcAccount struct {
	AuthType AuthType `db:"auth_type"`
	// ID of the platform associated with the service account.
	OrgID string `db:"org_id"`
	// Initial Role Assignment
	Role string `db:"-"`
	// Environment authorized for service account.
	Environment string `db:"environment"`
	// Hydra client name.
	ClientName string `db:"client_name"`
	// Service Account ID.
	ClientID string `db:"client_id"`
	// Service Account Challenge for API key challenges.
	Challenge string `db:"challenge"`
	// ID of the user that created the service account.
	CreateUserID string    `db:"create_user_id"`
	Created      time.Time `db:"created"`
	// ID of the user that disabled the service account.
	DisableUserID string       `db:"disable_user_id"`
	Disabled      sql.NullTime `db:"disabled"`
}

type OAuthClient struct {
	OrgID        string       `db:"org_id"`
	ClientID     string       `db:"client_id"`
	ClientName   string       `db:"client_name"`
	Environment  string       `db:"environment"`
	CreateUserID string       `db:"created_user_id"`
	UpdateUserID string       `db:"updated_user_id"`
	DeleteUserID string       `db:"deleted_user_id"`
	Created      time.Time    `db:"created"`
	Updated      time.Time    `db:"updated"`
	Deleted      sql.NullTime `db:"deleted"`
}

// Deployed service or product.
type Service struct {
	ID          string    `db:"service_id"`
	Name        string    `db:"service_name"`
	Description string    `db:"description"`
	Default     bool      `db:"assign_default"`
	Created     time.Time `db:"created"`
	Updated     time.Time `db:"updated"`
}

// Permission defined by a service.
type ServicePermission struct {
	ID          string    `db:"id"`
	ServiceID   string    `db:"service_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Resource    string    `db:"resource"`
	Action      string    `db:"action"`
	Created     time.Time `db:"created"`
	Updated     time.Time `db:"updated"`
}

type DefaultService struct {
	GrantID         string       `db:"grant_id"`
	ServiceID       string       `db:"service_id"`
	Environment     string       `db:"environment"`
	Published       time.Time    `db:"published"`
	PublishedUserID string       `db:"published_user_id"`
	Retracted       sql.NullTime `db:"retracted"`
	RetractedUserID string       `db:"retracted_user_id"`
	Created         time.Time    `db:"created"`
}

// Grant to an org for all permissions associated with a given service.
type ServiceAssignment struct {
	GrantID      string         `db:"grant_id"`
	ServiceID    string         `db:"service_id"`
	Environment  string         `db:"environment"`
	OrgID        string         `db:"org_id"`
	Default      bool           `db:"assign_default"`
	AssignUserID string         `db:"assign_user_id"`
	Assigned     time.Time      `db:"assigned"`
	RevokeUserID sql.NullString `db:"revoke_user_id"`
	Revoked      sql.NullTime   `db:"revoked"`
	Created      time.Time      `db:"created"`
}

// Granular permission granted by a service assignment.
type OrgPermission struct {
	ID           string    `db:"id"`
	OrgID        string    `db:"org_id"`
	GrantID      string    `db:"grant_id"`
	ServiceID    string    `db:"service_id"`
	PermissionID string    `db:"permission_id"`
	Name         string    `db:"name"`
	Description  string    `db:"description"`
	Resource     string    `db:"resource"`
	Action       string    `db:"action"`
	Environment  string    `db:"environment"`
	Created      time.Time `db:"created"`
	Updated      time.Time `db:"updated"`
}

type FilterList struct {
	OrgID        string
	ID           []string
	Name         string
	SortBy       string
	SortByColumn string
	Status       []string
	Limit        int32
	Offset       int32
}

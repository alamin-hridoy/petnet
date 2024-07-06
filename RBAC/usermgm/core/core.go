package core

import (
	"time"
)

type Identity struct {
	ID         string
	Name       string
	OrgID      string
	EventID    string
	MFA        string
	MFATrial   int
	Token      string
	PWExpiry   time.Time
	ForceReset bool

	Locked       bool
	Retries      int
	TrackRetries bool
}

type User struct {
	ID           string
	FName        string
	LName        string
	Email        string
	PreferredMFA string

	EnableMFA  bool
	DisableMFA bool
}

type AuthCredential struct {
	Username     string
	Password     string
	AuthClientID string
	MFA          *MFAChallenge
}

type AuthCodeClient struct {
	OrgID          string
	ClientID       string
	ClientName     string
	ClientSecret   string
	Environment    string
	CORS           []string
	RedirectURLs   []string
	LogoutRedirect string
	Logo           string
	GrantTypes     []string
	ResponseTypes  []string
	Scopes         []string
	Audience       []string
	SubjectType    string
	AuthMethod     string
	AuthBackend    string
	AuthConfig     AuthClientConfig
	CreatedBy      string
	UpdatedBy      string
	DeletedBy      string
	Created        time.Time
	Updated        time.Time
	Deleted        time.Time
}

type AuthClientConfig struct {
	LoginTmpl       string
	OTPTmpl         string
	ConsentTmpl     string
	ForceConsent    bool
	SessionDuration time.Duration
	IdentitySource  string
}

type MFA struct {
	UserID    string
	Type      string
	Source    string
	MFAID     string
	ConfirmID string
	Codes     []string
	Confirmed time.Time
	Updated   time.Time
	Revoked   time.Time
}

type MFAChallenge struct {
	EventID    string
	ExternalID string
	EventDesc  string
	UserID     string
	SourceID   string
	Type       string
	Token      string
	Attempt    int
	Sources    []MFA
}

type Invite struct {
	// Sender info
	InvOrgID  string
	InvUserID string
	AppUserID string

	// Invite info
	ID              string
	Code            string
	OrgID           string
	OrgName         string
	FName           string
	LName           string
	RoleID          string
	Email           string
	CustomEmailData map[string]string

	// Auth info
	Username string
	Password string

	// Status info
	Status string
	Sent   time.Time
	Expiry time.Time
}

type UserActivation struct {
	ID              string
	CustomEmailData map[string]string
}

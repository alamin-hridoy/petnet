package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ory/hydra-client-go/models"

	client "brank.as/rbac/svcutil/hydraclient"
	"brank.as/rbac/usermgm/errors/session"
)

// Identity provides details about users
type Identity struct {
	UserID     string            `json:"user_id"`
	OrgID      string            `json:"org_id"`
	MFAEventID string            `json:"mfa_event_id"`
	MFAType    string            `json:"mfa_type"`
	PWExpiry   time.Time         `json:"pw_expiry"`
	PWReset    bool              `json:"required"`
	ForceLogin bool              `json:"force_login"`
	Session    map[string]string `json:"session"`
	OpenID     map[string]string `json:"open_id"`
}

// LoginContext is the context from authenticated identity to set in accept login request.
// The context will be available later in OIDC /userinfo endpoint.
func (id *Identity) LoginContext() map[string]string {
	c := make(map[string]string, len(id.OpenID)+2)
	for k, v := range id.OpenID {
		c[k] = v
	}
	c["userid"] = id.UserID
	c["orgid"] = id.OrgID
	return c
}

// SessionContext is the internal context from authenticated identity to set in accept login request.
// The context will be available later in auth middleware.
func (id *Identity) SessionContext() map[string]string {
	c := make(map[string]string, len(id.OpenID)+2)
	for k, v := range id.Session {
		c[k] = v
	}
	c["userid"] = id.UserID
	c["orgid"] = id.OrgID
	return c
}

type ErrCode int

const (
	Unknown ErrCode = iota
	NotFound
	PermissionDenied
	ExistingSession
	InvalidRecord
	OTPInvalid
	OTPExpired
	ExpiredPassword
)

// Error gives a specific http error for validation or authentication errors.
type Error struct {
	Err        error
	Code       ErrCode
	StatusCode int

	Message       string
	TrackAttempts bool
	AttemptRemain int
	Errors        map[string]string
}

// MergeDetails adds defaults from dt for missing messages in Errors.
func (e *Error) MergeDetails(dt map[string]string) {
	for k, v := range dt {
		if _, ok := e.Errors[k]; !ok {
			e.Errors[k] = v
		}
	}
}

func IsError(err error) bool { return err != nil && errors.Is(err, Error{}) }
func FromError(err error) Error {
	e := &Error{Errors: map[string]string{}}
	errors.As(err, e)
	return *e
}

func FromStatus(err error) Error {
	if deets := session.FromError(err); deets != nil {
		e := deets.GetErrorDetails()
		if e == nil {
			e = map[string]string{}
		}
		// Error details from session service.
		return Error{
			Err:           err,
			StatusCode:    http.StatusUnauthorized,
			Message:       deets.GetMessage(),
			TrackAttempts: deets.GetTrackingAttempts(),
			AttemptRemain: int(deets.GetRemainingAttempts()),
			Errors:        e,
		}
	}
	return Error{
		StatusCode: http.StatusUnauthorized,
		Message:    "account does not exist or credentials are invalid",
		Errors:     map[string]string{},
	}
}

// Error fulfils the error interface.
func (ae Error) Error() string { return ae.Message }

// Cause of the underlying error.  Used for internal logging.
func (ae Error) Cause() error {
	if ae.Err == nil {
		return errors.New(ae.Message)
	}
	return ae.Err
}

type Challenge struct {
	ID          string
	Username    string
	Password    string
	HydraClient string
	Extra       map[string]string
}

type OTPChallenge struct {
	Code  string
	Type  string
	Event string
}

// Authenticator interface for receiving identities
type Authenticator interface {
	// Authenticate tries to authenticate with given username and password in parameters.
	Authenticate(context.Context, Challenge, *OTPChallenge) (*Identity, error)

	// ResetMFA will re-genenrate and re-send (if necessary) the MFA token.
	ResetMFA(context.Context, Challenge, OTPChallenge) (*Identity, error)

	// Lookup returns the identity with given userID.
	// It is intended to use when user is using remember me feature and
	// has a valid authenticated session.
	Lookup(context.Context, Challenge) (*Identity, error)

	// Remember returns the session duration.
	Remember() int64

	// Consent grantor for this authenticator.
	Consent() ConsentGrantor
}

func ParseClientConfig(hcl *models.OAuth2Client) (*client.AuthConfig, error) {
	b, err := json.Marshal(hcl.Metadata)
	if err != nil {
		return nil, err
	}
	c := &client.AuthConfig{}
	if err := json.Unmarshal(b, c); err != nil {
		return nil, err
	}
	return c, nil
}

package dummy

import (
	"context"
	"net/http"

	"brank.as/rbac/idp/auth"
)

var _ auth.Authenticator = (*Dummy)(nil)

type Dummy struct {
	identity *auth.Identity
	password string
}

func New() *Dummy {
	return &Dummy{
		identity: &auth.Identity{
			UserID: "123e4567-e89b-12d3-a456-426614174000",
			OrgID:  "223e4567-e89b-12d3-a456-426614174000",
		},
		password: "dummy",
	}
}

func (d *Dummy) Lookup(_ context.Context, u auth.Challenge) (*auth.Identity, error) {
	if u.ID == d.identity.UserID {
		return d.identity, nil
	}
	return nil, auth.Error{
		Code:       auth.NotFound,
		StatusCode: http.StatusBadRequest,
	}
}

func (d *Dummy) Authenticate(_ context.Context, c auth.Challenge, o *auth.OTPChallenge) (*auth.Identity, error) {
	if c.Username == d.identity.UserID && c.Password == d.password || o != nil {
		return d.identity, nil
	}
	return nil, auth.Error{
		Code:       auth.NotFound,
		StatusCode: http.StatusBadRequest,
	}
}

func (d *Dummy) ResetMFA(_ context.Context, c auth.Challenge, o auth.OTPChallenge) (*auth.Identity, error) {
	return nil, auth.Error{Code: auth.InvalidRecord, StatusCode: http.StatusBadRequest}
}

// Remember returns the session duration.
func (*Dummy) Remember() int64 { return 86400 }

func (*Dummy) Consent() auth.ConsentGrantor { return nil }

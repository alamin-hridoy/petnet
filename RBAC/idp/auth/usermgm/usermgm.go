package usermgm

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/idp/auth"

	upb "brank.as/rbac/gunk/v1/user"
)

var _ auth.Authenticator = (*Auth)(nil)

// NewAuthenticator instantiate implementation of auth.Authenticator
// that searches user credentials from the usermgm database.
func NewAuthenticator(uaCl upb.UserAuthServiceClient, uCL upb.UserServiceClient, dur time.Duration) (*Auth, error) {
	return &Auth{uaCl: uaCl, uCL: uCL, dur: int64(dur / time.Second)}, nil
}

// Auth implements Authenticator that calls the AuthenticateUser
// service on usermgm.
type Auth struct {
	uaCl upb.UserAuthServiceClient
	uCL  upb.UserServiceClient
	dur  int64
}

func toIdentity(aur *upb.AuthenticateUserResponse) *auth.Identity {
	return &auth.Identity{
		UserID: aur.UserID,
		OrgID:  aur.OrgID,
	}
}

// Authenticate looks up the user with given email and returns user information if the password is match.
func (a *Auth) Authenticate(ctx context.Context, c auth.Challenge, _ *auth.OTPChallenge) (*auth.Identity, error) {
	r, err := a.uaCl.AuthenticateUser(ctx, &upb.AuthenticateUserRequest{
		Username: c.Username,
		Password: c.Password,
	})
	if err != nil {
		e := auth.FromStatus(err)
		switch status.Code(err) {
		case codes.NotFound:
			e.Code = auth.NotFound
		case codes.PermissionDenied:
			e.Code = auth.PermissionDenied
		case codes.AlreadyExists:
			e.Code = auth.ExistingSession
		}
		return nil, e
	}
	return toIdentity(r), nil
}

func (a *Auth) ResetMFA(ctx context.Context, c auth.Challenge, _ auth.OTPChallenge) (*auth.Identity, error) {
	// Deprecated backend.  No OTP support.
	return nil, auth.Error{Code: auth.PermissionDenied, StatusCode: http.StatusBadRequest}
}

// Lookup returns the identity with matching given user id in parameter.
func (a *Auth) Lookup(ctx context.Context, c auth.Challenge) (*auth.Identity, error) {
	r, err := a.uCL.GetUser(ctx, &upb.GetUserRequest{ID: c.ID})
	if err != nil {
		e := auth.FromStatus(err)
		switch status.Code(err) {
		case codes.NotFound:
			e.Code = auth.NotFound
		case codes.PermissionDenied:
			e.Code = auth.OTPInvalid
		}
		return nil, e
	}
	return &auth.Identity{
		UserID: r.User.ID,
		OrgID:  r.User.OrgID,
	}, nil
}

// Remember returns the session duration.
func (a *Auth) Remember() int64              { return a.dur }
func (a *Auth) Consent() auth.ConsentGrantor { return nil }

package usermgm

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/idp/auth"

	apb "brank.as/rbac/gunk/v1/authenticate"
	mpb "brank.as/rbac/gunk/v1/mfa"
)

var _ auth.Authenticator = (*Session)(nil)

func NewSessionAuth(s apb.SessionServiceClient, dur time.Duration, c auth.ConsentGrantor) *Session {
	return &Session{sCl: s, dur: int64(dur / time.Second), c: c}
}

// Session implements Authenticator that calls a Login service implementation.
type Session struct {
	sCl apb.SessionServiceClient
	dur int64
	c   auth.ConsentGrantor
}

// Authenticate looks up the user information using the user-presented challenge for authentication.
func (s *Session) Authenticate(ctx context.Context, c auth.Challenge, o *auth.OTPChallenge) (*auth.Identity, error) {
	if o != nil { // Login MFA validation
		r, err := s.sCl.Login(ctx, &apb.LoginRequest{
			ClientID:   c.HydraClient,
			Subject:    c.ID,
			MFAEventID: o.Event,
			MFAType:    mpb.MFA(mpb.MFA_value[o.Type]),
			MFAToken:   o.Code,
		})
		if err != nil {
			e := auth.FromStatus(err)
			switch status.Code(err) {
			case codes.NotFound:
				e.Code = auth.NotFound
			case codes.PermissionDenied:
				e.Code = auth.OTPInvalid
			case codes.ResourceExhausted:
				e.Code = auth.OTPInvalid
			default:
				e.Code = auth.NotFound
			}
			return nil, e
		}
		return s.toIdentity(r), nil
	}

	r, err := s.sCl.Login(ctx, &apb.LoginRequest{
		Username: c.Username,
		Password: c.Password,
		ClientID: c.HydraClient,
		Extra:    c.Extra,
	})
	if err != nil {
		e := auth.FromStatus(err)
		switch status.Code(err) {
		case codes.NotFound:
			e.Code = auth.NotFound
		case codes.ResourceExhausted:
			e.Code = auth.ExpiredPassword
		case codes.PermissionDenied:
			e.Code = auth.PermissionDenied
		case codes.AlreadyExists:
			e.Code = auth.ExistingSession
		case codes.FailedPrecondition:
			e.Code = auth.InvalidRecord
		default:
			e.Code = auth.NotFound
		}
		return nil, e
	}
	return s.toIdentity(r), err
}

// ResetMFA ...
func (s *Session) ResetMFA(ctx context.Context, c auth.Challenge, o auth.OTPChallenge) (*auth.Identity, error) {
	r, err := s.sCl.RetryMFA(ctx, &apb.RetryMFARequest{
		Subject:    c.ID,
		ClientID:   c.HydraClient,
		MFAEventID: o.Event,
		MFAType:    mpb.MFA(mpb.MFA_value[o.Type]),
	})
	if err != nil {
		e := auth.FromStatus(err)
		switch status.Code(err) {
		case codes.NotFound:
			e.Code = auth.NotFound
		case codes.PermissionDenied:
			e.Code = auth.OTPInvalid
		case codes.ResourceExhausted:
			e.Code = auth.OTPInvalid
		default:
			e.Code = auth.NotFound
		}
		return nil, e
	}
	return s.toIdentity(r), nil
}

// Lookup returns the identity with matching given user id in parameter.
func (s *Session) Lookup(ctx context.Context, a auth.Challenge) (*auth.Identity, error) {
	r, err := s.sCl.GetSession(ctx, &apb.GetSessionRequest{
		UserID:   a.ID,
		ClientID: a.HydraClient,
		Extra:    a.Extra,
	})
	if err != nil {
		e := auth.FromStatus(err)
		switch status.Code(err) {
		case codes.NotFound:
			e.Code = auth.NotFound
		case codes.PermissionDenied:
			e.Code = auth.OTPInvalid
		case codes.ResourceExhausted:
			e.Code = auth.OTPInvalid
		default:
			e.Code = auth.NotFound
		}
		return nil, e
	}
	return s.toIdentity(r), nil
}

func (*Session) toIdentity(r *apb.Session) *auth.Identity {
	if r == nil {
		return nil
	}
	return &auth.Identity{
		UserID:     r.GetUserID(),
		OrgID:      r.GetOrgID(),
		MFAEventID: r.GetMFAEventID(),
		MFAType:    r.GetMFAType().String(),
		PWExpiry:   r.GetPasswordExpiry().AsTime(),
		PWReset:    r.GetResetRequired(),
		ForceLogin: r.GetForceLogin(),
		Session:    r.GetSession(),
		OpenID:     r.GetOpenID(),
	}
}

// Remember returns the session duration.
func (s *Session) Remember() int64 { return s.dur }

// Consent backend.
func (s *Session) Consent() auth.ConsentGrantor { return s.c }

package session

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"

	"brank.as/rbac/usermgm/errors/session"

	authpb "brank.as/rbac/gunk/v1/authenticate"
	mpb "brank.as/rbac/gunk/v1/mfa"
)

func (s *Svc) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.Session, error) {
	switch req.GetUsername() {
	case "retry":
		return nil, session.Error(codes.InvalidArgument, "retry", &authpb.SessionError{
			Message:           "retry mock",
			TrackingAttempts:  true,
			RemainingAttempts: 3,
			ErrorDetails:      map[string]string{},
		})
	case "bademail":
		return nil, session.Error(codes.InvalidArgument, "bademail", &authpb.SessionError{
			Message:           "error details for testing",
			TrackingAttempts:  true,
			RemainingAttempts: 1,
			ErrorDetails:      map[string]string{"emailError": "Invalid username (reason)"},
		})
	case "badpass":
		return nil, session.Error(codes.InvalidArgument, "badpass", &authpb.SessionError{
			Message:           "error details for testing",
			TrackingAttempts:  true,
			RemainingAttempts: 1,
			ErrorDetails: map[string]string{
				"passwordError": "invalid password (must contain one upper, one lower, and one emoji)",
			},
		})
	case "expired":
		return nil, session.Error(codes.ResourceExhausted, "expired", &authpb.SessionError{
			Message: "password expired",
			ErrorDetails: map[string]string{
				"passwordError": "invalid password (must contain one upper, one lower, and one emoji)",
				"username":      "Password has expired. Please contact your organization administrator.",
			},
		})
	case "mfa":
		return &authpb.Session{
			MFAEventID: uuid.NewString(),
			MFAType:    mpb.MFA_CODE,
		}, nil
	case "valid":
		return &authpb.Session{
			UserID:  uuid.NewString(),
			OrgID:   uuid.NewString(),
			Session: map[string]string{"session-test": "value"},
			OpenID:  map[string]string{"opeid-test": "vals"},
		}, nil
	default:
		return nil, session.Error(codes.OutOfRange, "no matching mock", &authpb.SessionError{
			Message: "no mock",
			ErrorDetails: map[string]string{
				"emailError": "invalid email/username - no matching test mock",
			},
		})
	}
}

func (s *Svc) GetSession(context.Context, *authpb.GetSessionRequest) (*authpb.Session, error) {
	return &authpb.Session{
		UserID:  uuid.NewString(),
		OrgID:   uuid.NewString(),
		Session: map[string]string{"session-test": "value"},
		OpenID:  map[string]string{"opeid-test": "vals"},
	}, nil
}

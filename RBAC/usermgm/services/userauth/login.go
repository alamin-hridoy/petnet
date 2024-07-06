package userauth

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/errors/session"

	apb "brank.as/rbac/gunk/v1/authenticate"
	mpb "brank.as/rbac/gunk/v1/mfa"
)

func (s *Svc) Login(ctx context.Context, req *apb.LoginRequest) (*apb.Session, error) {
	log := logging.FromContext(ctx).WithField("method", "service.auth.Authenticate")
	log.Trace("request received")

	ev := req.GetMFAEventID() != ""
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Username, validation.When(!ev, validation.Required)),
		validation.Field(&req.Password, validation.When(!ev, validation.Required)),
		validation.Field(&req.MFAEventID, is.UUIDv4),
		validation.Field(&req.MFAToken,
			validation.When(ev, validation.Required, is.Alphanumeric),
		),
		validation.Field(&req.MFAType,
			validation.When(ev, validation.By(func(interface{}) error {
				return validation.Validate(mpb.MFA_name[mpb.MFA_value[req.MFAType.String()]],
					validation.Required)
			})),
		),
	); err != nil {
		logging.WithError(err, log).Info("invalid request")
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	var m *core.MFAChallenge
	if req.MFAEventID != "" {
		m = &core.MFAChallenge{
			EventID: req.MFAEventID,
			Type:    req.MFAType.String(),
			Token:   req.MFAToken,
		}
	}

	usr, err := s.auth.AuthUser(ctx, core.AuthCredential{
		Username:     req.GetUsername(),
		Password:     req.GetPassword(),
		AuthClientID: req.GetClientID(),
		MFA:          m,
	})
	if err != nil {
		logging.WithError(err, log).Error("failed to authenticate user")
		switch status.Code(err) {
		case codes.NotFound:
			return nil, status.Error(codes.InvalidArgument, "user not found")
		case codes.ResourceExhausted:
			return nil, session.Error(codes.ResourceExhausted, "account locked",
				&apb.SessionError{
					Message: "attempts exceeded - account locked",
					ErrorDetails: map[string]string{
						"username": "Login attempts exceeded - account locked.",
					},
				})
		case codes.FailedPrecondition:
			if s := session.FromError(err); s != nil {
				return nil, session.Error(codes.InvalidArgument, s.Message, s)
			}
		}

		if usr != nil {
			return nil, session.Error(codes.InvalidArgument, "invalid username or password",
				&apb.SessionError{
					Message:           "invalid login attempt",
					RemainingAttempts: int32(usr.Retries),
					TrackingAttempts:  usr.TrackRetries,
				})
		}
		return nil, status.Error(codes.InvalidArgument, "invalid argument given")
	}

	return &apb.Session{
		UserID:     usr.ID,
		OrgID:      usr.OrgID,
		MFAEventID: usr.EventID,
		MFAType:    mpb.MFA(mpb.MFA_value[usr.MFA]),
	}, nil
}

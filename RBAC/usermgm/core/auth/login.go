package auth

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/errors/session"
	"brank.as/rbac/usermgm/storage"

	apb "brank.as/rbac/gunk/v1/authenticate"
)

const loginEvent = "User Login"

// AuthUser ...
func (s *Svc) AuthUser(ctx context.Context, c core.AuthCredential) (id *core.Identity, err error) {
	log := logging.FromContext(ctx).WithField("method", "core.user.authenticateuser")

	ctx, err = s.db.NewTransacton(ctx)
	if err != nil {
		logging.WithError(err, log).Error("db transaction")
		return nil, status.Error(codes.Internal, "user authentication failed")
	}
	defer func() {
		if err := s.db.Rollback(ctx); err != nil {
			logging.WithError(err, log).Error("db rollback")
		}
	}()
	var u *storage.User
	if c.Username != "" && c.Password != "" {
		usr, err := s.usr.GetUser(ctx, c.Username, c.Password)
		if err != nil {
			logging.WithError(err, log).Error("user auth")
			if err := s.db.Commit(ctx); err != nil { // record failed attempt and/or lock.
				logging.WithError(err, log).Error("commit user login failed")
			}
			switch {
			case err != storage.NotFound:
				return nil, status.Error(codes.Internal, "user authentication failed")
			case usr == nil:
				return &core.Identity{}, status.Error(codes.InvalidArgument, "username or password invalid")
			case s.lock == 0:
				// s.lock == 0 means that locking is disabled
				return &core.Identity{},
					status.Error(codes.InvalidArgument, "username or password invalid")
			case usr.Locked.Valid:
				return &core.Identity{Locked: true}, nil
			case usr.FailCount >= s.lock:
				if err := s.db.LockUser(ctx, usr.ID); err != nil {
					logging.WithError(err, log).Error("locking user record")
				}
				return &core.Identity{Locked: true},
					status.Error(codes.ResourceExhausted, "username or password invalid - account locked")
			}
			return &core.Identity{Retries: s.lock - usr.FailCount, TrackRetries: true},
				status.Error(codes.InvalidArgument, "username or password invalid")
		}
		u = usr
	} else {
		ev, err := s.usr.GetMFAEventByID(ctx, c.MFA.EventID)
		if err != nil {
			logging.WithError(err, log).Error("mfa event")
			if err == storage.NotFound {
				return nil, status.Error(codes.NotFound, "login mfa event not found")
			}
			return nil, status.Error(codes.Internal, "mfa authentication failed")
		}
		if ev.Desc != loginEvent {
			return nil, status.Error(codes.Internal, "mfa authentication failed")
		}
		usr, err := s.usr.GetUserByID(ctx, ev.UserID)
		if err != nil {
			logging.WithError(err, log).Error("mfa event")
			if err == storage.NotFound {
				return nil, status.Error(codes.NotFound, "user not found")
			}
			return nil, status.Error(codes.Internal, "mfa authentication failed")
		}
		u = usr
	}
	if u.Deleted.Valid {
		return nil, status.Error(codes.FailedPrecondition, "user account not active")
	}

	if s.requireEmail && !u.EmailVerified {
		return nil, session.Error(codes.FailedPrecondition, "user's e-mail not verified", &apb.SessionError{
			Message: "user's e-mail not verified",
			ErrorDetails: map[string]string{
				"username": "Your e-mail is not verified. Please, check your inbox.",
			},
		})
	}

	org, err := s.db.GetOrgByID(ctx, u.OrgID)
	if err != nil {
		logging.WithError(err, log).Error("fetch org")
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "org record unavailable")
	}

	defer func(e *error) { // record login success
		err := *e
		if err != nil {
			return
		}
		if err := s.db.Commit(ctx); err != nil {
			logging.WithError(err, log).Error("commit user login")
		}
	}(&err)

	// Login MFA validation
	if c.MFA != nil {
		c.MFA.UserID = u.ID
		e, err := s.mfa.MFAuth(ctx, *c.MFA)
		if err != nil {
			logging.WithError(err, log).Error("mfa validation")
			if status.Code(err) != codes.Unknown {
				return nil, err
			}
			return nil, status.Error(codes.Internal, "mfa validation failed")
		}
		return &core.Identity{ID: e.UserID, Name: u.Username, OrgID: u.OrgID}, nil
	}

	// TODO: "initiation" duration, for registering an MFA
	if org.MFALogin.Bool || u.MFALogin {
		if u.PreferredMFA == "" {
			return nil, status.Error(codes.FailedPrecondition, "no MFA sources registered")
		}
		e, err := s.mfa.InitiateMFA(ctx, core.MFAChallenge{
			EventDesc: loginEvent,
			UserID:    u.ID,
			Type:      u.PreferredMFA,
		})
		if err != nil {
			logging.WithError(err, log).Error("mfa initiate")
			if status.Code(err) != codes.Unknown {
				return nil, err
			}
			return nil, status.Error(codes.Internal, "2FA procedure failed")
		}
		return &core.Identity{
			ID:      u.ID,
			Name:    c.Username,
			OrgID:   u.OrgID,
			EventID: e.EventID,
			MFA:     e.Type,
			Token:   e.Token,
		}, nil
	}

	return &core.Identity{
		ID:    u.ID,
		Name:  c.Username,
		OrgID: u.OrgID,
		// TODO: PW Reset schedule/force reset
		PWExpiry:   time.Time{},
		ForceReset: false,
	}, nil
}

func (s *Svc) UserSession(ctx context.Context, userID string) (*core.Identity, error) {
	log := logging.FromContext(ctx).WithField("method", "core.auth.usersession")
	log.Info("received")

	u, err := s.usr.GetUserByID(ctx, userID)
	if err != nil {
		logging.WithError(err, log).Error("user auth")
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "user authentication failed")
	}
	if u.Deleted.Valid {
		return nil, status.Error(codes.FailedPrecondition, "user account not active")
	}
	o, err := s.db.GetOrgByID(ctx, u.OrgID)
	if err != nil {
		logging.WithError(err, log).Error("fetch org")
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "organization validation failed")
	}
	if !o.Active {
		return nil, status.Error(codes.FailedPrecondition, "organization account not active")
	}

	return &core.Identity{
		ID:    userID,
		OrgID: u.OrgID,
	}, nil
}

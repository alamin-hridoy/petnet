package mfauth

import (
	"context"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/svcutil/random"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"

	"brank.as/rbac/serviceutil/logging"
)

func (s *Svc) RegisterMFA(ctx context.Context, c core.MFA) (*core.MFA, error) {
	log := logging.FromContext(ctx).WithField("method", "core.mfauth.register")
	log.WithField("request", c).Trace("received")

	u, err := s.st.GetUserByID(ctx, c.UserID)
	if err != nil {
		logging.WithError(err, log).Error("user from storage")
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, err
	}
	mfas, err := s.st.GetMFAByType(ctx, c.UserID, c.Type)
	if err != nil && err != storage.NotFound {
		logging.WithError(err, log).Error("mfas from storage")
		return nil, status.Error(codes.Internal, "failed to read mfa")
	}

	switch c.Type {
	case storage.EMail:
		log.WithFields(logrus.Fields{
			"req_source": c.Source,
			"user":       u.Email,
			"verified":   u.EmailVerified,
		}).Info("received")
		if c.Source == u.Email && u.EmailVerified {
			for _, m := range mfas {
				if m.Active && m.Token == u.Email {
					c.MFAID = m.ID
					c.Confirmed = m.Confirmed.Time
					return &c, nil
				}
			}
			m, err := s.st.CreateMFA(ctx, storage.MFA{
				UserID:  c.UserID,
				MFAType: c.Type,
				Token:   c.Source,
				Active:  true,
			})
			if err != nil {
				logging.WithError(err, log).Error("storage create mfa")
				if err == storage.NotFound {
					return nil, status.Error(codes.Internal, "mfa registration failed")
				}
				return nil, err
			}
			c.MFAID = m.ID
			c.Confirmed = time.Now()
			if u.PreferredMFA == "" {
				u.PreferredMFA = c.Type
				if _, err := s.st.UpdateUserByID(ctx, *u); err != nil {
					logging.WithError(err, log).Error("update preference")
				}
			}
			return &c, nil
		}
		m, err := s.initMFA(ctx, c)
		if err != nil {
			logging.WithError(err, log).Error("mfa init")
			return nil, err
		}
		if err := s.em.EmailMFA(c.Source, m.Source); err != nil {
			if err != nil {
				logging.WithError(err, log).Error("failed to send mfa email")
				return nil, status.Error(codes.InvalidArgument, "failed to send email")
			}
		}
		m.Source = ""
		return m, nil
	case storage.PINCode, storage.SMS:
		m, err := s.initMFA(ctx, c)
		if err != nil {
			logging.WithError(err, log).Error("mfa init")
			return nil, err
		}
		if c.Type != storage.SMS {
			m.Source = ""
		}
		if u.PreferredMFA == "" {
			u.PreferredMFA = c.Type
			if _, err := s.st.UpdateUserByID(ctx, *u); err != nil {
				logging.WithError(err, log).Error("update preference")
			}
		}
		return m, nil
	case storage.TOTP:
		k, err := totp.Generate(totp.GenerateOpts{
			Issuer:      s.svcName,
			AccountName: u.Username,
			SecretSize:  32,
		})
		if err != nil {
			logging.WithError(err, log).Error("totp generate")
			return nil, status.Error(codes.Internal, "mfa registration failed")
		}
		c.Source = k.Secret()
		m, err := s.initMFA(ctx, c)
		if err != nil {
			logging.WithError(err, log).Error("totp init")
			if status.Code(err) != codes.Unknown {
				return nil, err
			}
			return nil, status.Error(codes.Internal, "initialization failed")
		}
		m.Source = k.String()
		return m, nil
	case storage.Recovery:
		c.Codes = make([]string, 6)
		txCtx, err := s.st.NewTransacton(ctx)
		if err != nil {
			logging.WithError(err, log).Error("create transaction")
			return nil, status.Error(codes.Internal, "generating recovery codes failed")
		}
		defer s.st.Rollback(txCtx)
		for i := range c.Codes {
			c.Codes[i] = random.NumString(8)
			c.Source = c.Codes[i]
			if _, err := s.initMFA(txCtx, c); err != nil {
				if err != nil {
					return nil, status.Error(codes.Internal, "generating recovery codes failed")
				}
			}
		}
		if err := s.st.Commit(txCtx); err != nil {
			logging.WithError(err, log).Error("commit transaction")
			return nil, status.Error(codes.Internal, "generating recovery codes failed")
		}
		return &c, nil
	default:
		return nil, status.Error(codes.InvalidArgument, "mfa type not supported")
	}
}

func (s *Svc) initMFA(ctx context.Context, c core.MFA) (*core.MFA, error) {
	log := logging.FromContext(ctx).WithField("method", "core.mfauth.initmfa")

	m, err := s.st.CreateMFA(ctx, storage.MFA{
		UserID:   c.UserID,
		MFAType:  c.Type,
		Token:    c.Source,
		Deadline: time.Now().Add(s.dur),
	})
	if err != nil {
		logging.WithError(err, log).Error("storage create mfa")
		if err == storage.NotFound {
			return nil, status.Error(codes.Internal, "mfa registration failed")
		}
		return nil, err
	}
	c.MFAID = m.ID
	switch c.Type {
	case storage.PINCode, storage.Recovery:
		return &c, nil
	}

	ev, err := s.st.CreateMFAEvent(ctx, storage.MFAEvent{
		UserID:     c.UserID,
		MFAID:      m.ID,
		MFAType:    c.Type,
		Desc:       "mfa activation",
		Validation: true,
		Deadline:   m.Deadline,
	})
	if err != nil {
		return nil, err
	}
	c.ConfirmID = ev.EventID
	switch c.Type {
	case storage.EMail, storage.SMS:
		c.Source = ev.Token
	}
	return &c, nil
}

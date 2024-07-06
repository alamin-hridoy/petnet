package mfauth

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/svcutil/random"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

func (s *Svc) InitiateMFA(ctx context.Context, c core.MFAChallenge) (*core.MFAChallenge, error) {
	log := logging.FromContext(ctx).WithField("method", "core.mfauth.initiate")
	log.WithField("request", c).Trace("received")

	u, err := s.st.GetUserByID(ctx, c.UserID)
	if err != nil {
		logging.WithError(err, log).Error("fetch user account")
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, err
	}
	if u.Deleted.Valid {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	// default to user preference
	if c.Type == "" && c.SourceID == "" {
		c.Type = u.PreferredMFA
	}

	src := storage.MFA{}
	if c.SourceID != "" {
		sc, err := s.st.GetMFAByID(ctx, c.SourceID)
		if err != nil {
			logging.WithError(err, log).Error("mfa from storage")
			if err == storage.NotFound {
				return nil, status.Error(codes.NotFound, "mfa source not found")
			}
			return nil, err
		}
		if !sc.Active {
			return nil, status.Error(codes.InvalidArgument, "mfa source not active")
		}
		src = *sc
	} else {
		sc, err := s.st.GetActiveMFAByType(ctx, c.UserID, c.Type)
		if err != nil {
			logging.WithError(err, log).Error("list mfa type from storage")
			if err == storage.NotFound {
				return nil, status.Error(codes.NotFound, "mfa source not found")
			}
			return nil, err
		}
		src = sc[0]
	}

	switch src.MFAType {
	case storage.TOTP, storage.PINCode, storage.Recovery:
		c.Token = src.Token
	case storage.EMail, storage.SMS:
		c.Token = random.NumString(6)
	}
	if src.MFAType == storage.EMail {
		if err := s.em.EmailMFA(src.Token, c.Token); err != nil {
			logging.WithError(err, log).Error("failed to send mfa email")
			return nil, status.Error(codes.InvalidArgument, "failed to send email")
		}
	}

	ev, err := s.st.CreateMFAEvent(ctx, storage.MFAEvent{
		UserID:    c.UserID,
		MFAID:     src.ID,
		MFAType:   src.MFAType,
		Token:     c.Token,
		Desc:      c.EventDesc,
		Initiated: time.Now(),
		Deadline:  time.Now().Add(s.dur),
	})
	if err != nil {
		logging.WithError(err, log).Error("storage create event")
		return nil, status.Error(codes.Internal, "failed to initiate mfa")
	}
	switch src.MFAType {
	case storage.TOTP, storage.PINCode, storage.Recovery:
		src.Token = ""
		fallthrough
	case storage.EMail:
		ev.Token = ""
	case storage.SMS:
	}
	log.WithField("event", ev).WithField("src", src).Trace("mfa event generated")

	return &core.MFAChallenge{
		EventID:  ev.EventID,
		UserID:   c.UserID,
		SourceID: ev.MFAID,
		Type:     ev.MFAType,
		Token:    ev.Token,
		Sources: []core.MFA{{
			UserID:    src.UserID,
			Type:      src.MFAType,
			Source:    src.Token,
			MFAID:     src.ID,
			Confirmed: src.Confirmed.Time,
			Updated:   src.Updated,
			Revoked:   src.Revoked.Time,
		}},
	}, nil
}

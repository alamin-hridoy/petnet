package postgres

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/pquerna/otp/totp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/svcutil/random"
	"brank.as/rbac/usermgm/storage"
)

const mfaEventInsert = `
INSERT INTO mfa_event (
	user_id,
	mfa_id,
	mfa_type,
	active,
	validation,
	token,
	description,
	deadline,
	attempt
) VALUES (
	:user_id,
	:mfa_id,
	:mfa_type,
	True,
	:validation,
	:token,
	:description,
	:deadline,
	:attempt
)
RETURNING *`

// CreateMFAEvent creates new MFAEvent account returns the created MFAEvent's ID.
func (s *Storage) CreateMFAEvent(ctx context.Context, ma storage.MFAEvent) (*storage.MFAEvent, error) {
	switch "" {
	case ma.UserID, ma.MFAID:
		return nil, fmt.Errorf("invalid mfaEvent: user_id and mfa_id are required")
	}
	if time.Since(ma.Deadline) >= 0 {
		return nil, fmt.Errorf("deadline must be in the future")
	}

	m, err := s.GetMFAByID(ctx, ma.MFAID)
	if err != nil || m.UserID != ma.UserID {
		return nil, storage.NotFound
	}

	tok := ""
	switch ma.MFAType {
	case storage.TOTP, storage.Recovery, storage.PINCode:
		// Codes stored in mfa entry.
	case storage.SMS, storage.EMail:
		// Return plaintext token for SMS/Email only.
		tok = ma.Token
		if ma.Token == "" {
			tok = random.NumString(6)
		}
		h, err := s.hashPassword(tok, "")
		if err != nil {
			return nil, err
		}
		ma.Token = h
	default:
		return nil, fmt.Errorf("invalid mfaEvent type: %q", ma.MFAType)
	}

	stmt, err := s.prepareNamed(ctx, mfaEventInsert)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&ma, ma); err != nil {
		return nil, fmt.Errorf("executing MFAEvent insert: %w", err)
	}
	ma.Token = tok
	return &ma, nil
}

// GetMFAEventByID return MFAEvent account matched against ID primary key column
func (s *Storage) GetMFAEventByID(ctx context.Context, id string) (*storage.MFAEvent, error) {
	const mfaEventSelectByID = `SELECT * FROM mfa_event WHERE event_id = $1`
	sa := &storage.MFAEvent{}
	if err := s.db.Get(sa, mfaEventSelectByID, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	sa.Token = ""
	return sa, nil
}

// UpdateMFAEventToken sets a new token for the event.
func (s *Storage) UpdateMFAEventToken(ctx context.Context, e storage.MFAEvent) (*storage.MFAEvent, error) {
	const mfaEventToken = `UPDATE mfa_event SET token=:token WHERE event_id=:event_id RETURNING mfa_id`
	switch e.MFAType {
	case storage.SMS, storage.EMail:
		h, err := s.hashPassword(e.Token, "")
		if err != nil {
			return nil, err
		}
		e.Token = h
	default:
		return nil, fmt.Errorf("invalid MFA type - cannot be updated")
	}
	stmt, err := s.prepareNamed(ctx, mfaEventToken)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&e, e); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return &e, nil
}

// DisableUserMFAEvent will disable an event belonging to the user.
func (s *Storage) DisableUserMFAEvent(ctx context.Context, usrID, eventID string) ([]storage.MFAEvent, error) {
	const disableMFAEvent = `UPDATE mfa_event SET
(active, expired) = (False, now())
WHERE event_id=$1 AND confirmed IS NULL AND user_id=$2
RETURNING *`
	var lst []storage.MFAEvent
	if err := s.db.Select(&lst, disableMFAEvent, eventID, usrID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return lst, nil
}

// DisableMFAEvent expires the event without confirmation.
func (s *Storage) DisableMFAEvent(ctx context.Context, ma storage.MFAEvent) (*storage.MFAEvent, error) {
	const disableMFAEvent = `UPDATE mfa_event SET
(active, expired, confirmed) = (False, :expired, :confirmed)
WHERE event_id=:event_id AND confirmed IS NULL
RETURNING *`
	switch "" {
	case ma.EventID:
		return nil, fmt.Errorf("missing mfaEvent id")
	}
	stmt, err := s.prepareNamed(ctx, disableMFAEvent)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&ma, ma); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("executing MFAEvent disable: %w", err)
	}
	ma.Token = ""
	return &ma, nil
}

// equalConstTime does a constant time comparison on the two strings, returning true when equal.
func equalConstTime(a, b string) bool {
	ah, bh := sha256.New(), sha256.New()
	fmt.Fprint(ah, a)
	fmt.Fprint(bh, b)
	return subtle.ConstantTimeCompare(ah.Sum(nil), bh.Sum(nil)) == 1
}

func (s *Storage) ConfirmMFAEvent(ctx context.Context, ma storage.MFAEvent) (*storage.MFAEvent, error) {
	const validateSA = `SELECT * FROM mfa_event WHERE user_id=$1 AND mfa_type=$2 AND active order by initiated asc`
	switch "" {
	case ma.UserID, ma.MFAType:
		return nil, fmt.Errorf("user id and mfaEvent type are required")
	case ma.Token:
		return nil, fmt.Errorf("missing mfaEvent confirm token")
	}

	sa := []storage.MFAEvent{}
	if err := s.db.Select(&sa, validateSA, ma.UserID, ma.MFAType); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, fmt.Errorf("executing validate MFAEvent get: %w", err)
	}
	ev := storage.MFAEvent{}
	for _, e := range sa {
		if e.EventID == ma.EventID {
			ev = e
		}
		if time.Since(e.Deadline) > 0 {
			e.Expired, e.Active = true, false
			if _, err := s.DisableMFAEvent(ctx, e); err != nil {
				return nil, err
			}
			if e.EventID == ma.EventID {
				return nil, status.Error(codes.DeadlineExceeded, "mfa event expired")
			}
		}
	}
	if ma.EventID != "" {
		if ev.EventID == "" {
			return nil, storage.NotFound
		}
		// only keep identified event
		sa = []storage.MFAEvent{ev}
	}

	switch ma.MFAType {
	case storage.TOTP:
		mfas, err := s.GetMFAByType(ctx, ma.UserID, storage.TOTP)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, storage.NotFound
			}
			return nil, fmt.Errorf("executing validate MFAEvent get: %w", err)
		}
		e := sa[0]
		for _, m := range mfas {
			if !m.Active || m.Revoked.Valid {
				continue
			}
			if ok := totp.Validate(ma.Token, m.Token); ok {
				e.Token = ""
				e.Confirmed = sql.NullTime{Time: time.Now(), Valid: true}
				return s.DisableMFAEvent(ctx, e)
			}
		}

	case storage.Recovery, storage.PINCode:
		if len(sa) == 0 {
			return nil, storage.NotFound
		}
		e := sa[0]
		cd, err := s.GetMFAByType(ctx, ma.UserID, ma.MFAType)
		if err != nil {
			return nil, err
		}
		for _, m := range cd {
			if !m.Active || m.Revoked.Valid {
				continue
			}
			h, err := s.hashPassword(ma.Token, m.Token[:saltLen])
			if err != nil || !equalConstTime(h, m.Token) {
				continue
			}
			if m.MFAType == storage.Recovery {
				if _, err := s.DisableMFA(ctx, m.ID); err != nil {
					return nil, err
				}
			}
			m.Token = ""
			e.Confirmed = sql.NullTime{Time: time.Now(), Valid: true}
			return s.DisableMFAEvent(ctx, e)
		}

	case storage.SMS, storage.EMail:
		for _, m := range sa {
			if time.Since(m.Deadline) > 0 {
				continue
			}
			h, err := s.hashPassword(ma.Token, m.Token[:saltLen])
			if err != nil || !equalConstTime(h, m.Token) {
				continue
			}
			m.Token = ""
			m.Confirmed = sql.NullTime{Time: time.Now(), Valid: true}
			return s.DisableMFAEvent(ctx, m)
		}

	}
	return nil, storage.NotFound
}

func (s *Storage) ExpireMFAEvents(ctx context.Context) error {
	const expire = `
UPDATE mfa_event SET (active,expired)=(False,True) WHERE deadline < NOW() AND active`
	if _, err := s.db.ExecContext(ctx, expire); err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

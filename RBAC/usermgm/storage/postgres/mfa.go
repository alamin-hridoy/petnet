package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"brank.as/rbac/usermgm/storage"
)

const mfaInsert = `
INSERT INTO mfa (
	user_id,
	mfa_type,
	active,
	deadline,
	confirmed,
	token
) VALUES (
	:user_id,
	:mfa_type,
	:active,
	:deadline,
	:confirmed,
	:token
)
RETURNING *`

// CreateMFA creates new MFA account returns the created MFA's ID.
func (s *Storage) CreateMFA(ctx context.Context, ma storage.MFA) (*storage.MFA, error) {
	switch "" {
	case ma.UserID, ma.Token:
		return nil, fmt.Errorf("invalid mfa: user_id and token are required")
	}
	switch ma.MFAType {
	case storage.TOTP, storage.SMS, storage.EMail:
	case storage.PINCode, storage.Recovery:
		ma.Active = true                       // PIN/RECOVERY codes are active when set.
		h, err := s.hashPassword(ma.Token, "") // salt/hash the plaintext tokens
		if err != nil {
			return nil, err
		}
		ma.Token = h
	default:
		return nil, fmt.Errorf("invalid mfa type: %q", ma.MFAType)
	}
	if ma.Active {
		ma.Confirmed = sql.NullTime{Time: time.Now(), Valid: true}
	}
	stmt, err := s.prepareNamed(ctx, mfaInsert)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&ma, ma); err != nil {
		return nil, fmt.Errorf("executing MFA insert: %w", err)
	}
	return &ma, nil
}

// GetMFAByID return MFA account matched against ID primary key column
func (s *Storage) GetMFAByID(ctx context.Context, id string) (*storage.MFA, error) {
	const mfaSelectByID = `SELECT * FROM mfa WHERE id = $1`
	sa := &storage.MFA{}
	if err := sqlx.GetContext(ctx, s.queryer(ctx), sa, mfaSelectByID, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return sa, nil
}

// GetMFAByType return MFA account.
func (s *Storage) GetActiveMFAByType(ctx context.Context, user, mType string) ([]storage.MFA, error) {
	const mfaSelectByType = `SELECT * FROM mfa WHERE user_id = $1 and mfa_type = $2 and active`
	sa := []storage.MFA{}
	if err := s.db.Select(&sa, mfaSelectByType, user, mType); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	// select doesn't return error if resultset is empty
	if len(sa) == 0 {
		return nil, storage.NotFound
	}
	return sa, nil
}

// GetMFAByType return MFA account.
func (s *Storage) GetMFAByType(ctx context.Context, user, mType string) ([]storage.MFA, error) {
	const mfaSelectByType = `SELECT * FROM mfa WHERE user_id = $1 and mfa_type = $2`
	sa := []storage.MFA{}
	if err := s.db.Select(&sa, mfaSelectByType, user, mType); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return sa, nil
}

// GetMFAByUserID return MFA account associated with the given user.
func (s *Storage) GetMFAByUserID(ctx context.Context, usrID string) ([]storage.MFA, error) {
	const mfaSelectByUserID = `SELECT * FROM mfa WHERE user_id = $1 ORDER BY created desc`
	var lst []storage.MFA
	if err := s.db.Select(&lst, mfaSelectByUserID, usrID); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return lst, nil
}

const enableMFA = `
UPDATE mfa SET
    active=True,
    confirmed=NOW()
WHERE
    user_id=$1 AND id=$2
    and not active and revoked is null and confirmed is null
RETURNING *`

// EnableMFA permanently disables the service account credential.
func (s *Storage) EnableMFA(ctx context.Context, user, id string) (*time.Time, error) {
	switch "" {
	case id:
		return nil, fmt.Errorf("missing confirmation event id")
	}
	sa := storage.MFA{}
	if err := s.db.Get(&sa, enableMFA, user, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, fmt.Errorf("executing signup MFA insert: %w", err)
	}

	return &sa.Confirmed.Time, nil
}

const disableMFA = `
UPDATE mfa SET
    active=False,
    revoked=NOW()
WHERE
    id=$1 and revoked is null
RETURNING *`

// DisableMFA permanently disables the service account credential.
func (s *Storage) DisableMFA(ctx context.Context, id string) (time.Time, error) {
	t := time.Time{}
	switch "" {
	case id:
		return t, fmt.Errorf("missing mfa id")
	}
	sa := storage.MFA{}
	if err := s.db.Get(&sa, disableMFA, id); err != nil {
		if err == sql.ErrNoRows {
			return t, storage.NotFound
		}
		return t, fmt.Errorf("executing signup MFA insert: %w", err)
	}

	return sa.Revoked.Time, nil
}

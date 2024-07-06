package postgres

import (
	"context"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"brank.as/rbac/svcutil/random"
	"brank.as/rbac/usermgm/storage"
)

func hashChallenge(key, salt string) (string, error) {
	if salt == "" || len(salt) < saltLen {
		st, err := random.String(saltLen)
		if err != nil {
			return "", err
		}
		salt = st
	}
	if key == "" || len(key) < saltLen {
		return "", fmt.Errorf("key must be minimum length %d", saltLen)
	}
	h := sha512.New()
	h.Write([]byte(salt[:saltLen/2] + key + salt[saltLen/2:]))
	hbts := h.Sum(nil)

	hashHex := make([]byte, hex.EncodedLen(len(hbts)))
	hex.Encode(hashHex, hbts)
	return salt + string(hashHex), nil
}

const svcAccountInsert = `
INSERT INTO service_account (
	auth_type,
	org_id,
    environment,
	client_name,
	client_id,
	challenge,
	create_user_id
) VALUES (
	:auth_type,
	:org_id,
    :environment,
	:client_name,
	:client_id,
	:challenge,
	:create_user_id
)
RETURNING created`

// CreateSvcAccount creates new SvcAccount account returns the created SvcAccount's ID.
func (s *Storage) CreateSvcAccount(ctx context.Context, sa storage.SvcAccount) (string, error) {
	switch "" {
	case sa.ClientID, sa.CreateUserID, sa.ClientName, sa.OrgID:
		return "", fmt.Errorf("invalid account client_id, client_name, org_id, and create_user_id are required")
	}
	if sa.Challenge != "" {
		// salt/hash the api key
		h, err := hashChallenge(sa.Challenge, "")
		if err != nil {
			return "", err
		}
		sa.Challenge = h
	}
	stmt, err := s.db.PrepareNamedContext(ctx, svcAccountInsert)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	if err := stmt.Get(&sa.Created, sa); err != nil {
		return "", fmt.Errorf("executing SvcAccount insert: %w", err)
	}
	return sa.ClientID, nil
}

const svcAccountSelectByID = `
SELECT
	auth_type,
	org_id,
    environment,
	client_name,
	client_id,
	create_user_id,
    disable_user_id,
    created,
    disabled
FROM service_account
WHERE
	client_id = :client_id
`

// GetSvcAccountByID return SvcAccount account matched against ID primary key column
func (s *Storage) GetSvcAccountByID(ctx context.Context, clientID string) (*storage.SvcAccount, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, svcAccountSelectByID)
	if err != nil {
		return nil, fmt.Errorf("preparing named query GetSvcAccountByID: %w", err)
	}
	defer stmt.Close()
	sa := &storage.SvcAccount{ClientID: clientID}
	if err := stmt.Get(sa, sa); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return sa, nil
}

const svcAccountSelectByOrgID = `
SELECT
	auth_type,
	org_id,
    environment,
	client_name,
	client_id,
	create_user_id,
    disable_user_id,
    created,
    disabled
FROM service_account
WHERE
	org_id=:org_id
ORDER BY created desc`

// GetSvcAccountByOrgID return SvcAccount account associated with the given org.
func (s *Storage) GetSvcAccountByOrgID(ctx context.Context, orgID string) ([]storage.SvcAccount, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, svcAccountSelectByOrgID)
	if err != nil {
		return nil, fmt.Errorf("preparing named query GetSvcAccountByID: %w", err)
	}
	defer stmt.Close()
	sa := storage.SvcAccount{OrgID: orgID}
	var lst []storage.SvcAccount
	if err := stmt.Select(&lst, &sa); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return lst, nil
}

const disableSvcAcctUpdate = `
UPDATE service_account SET
    disable_user_id=:disable_user_id,
    disabled=NOW()
WHERE
    client_name=:client_name
RETURNING disabled`

// DisableSvcAccount permanently disables the service account credential.
func (s *Storage) DisableSvcAccount(ctx context.Context, sa storage.SvcAccount) (*time.Time, error) {
	switch "" {
	case sa.ClientID:
		return nil, fmt.Errorf("missing client id")
	case sa.DisableUserID:
		return nil, fmt.Errorf("missing user id for disable action")
	}
	stmt, err := s.db.PrepareNamedContext(ctx, disableSvcAcctUpdate)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&sa, &sa); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, fmt.Errorf("executing signup SvcAccount insert: %w", err)
	}

	return &sa.Disabled.Time, nil
}

func (s *Storage) ValidateSvcAccount(ctx context.Context, id, key string) (*storage.SvcAccount, error) {
	const validateSA = `SELECT * FROM service_account WHERE client_id=$1 AND disabled IS NULL`
	switch "" {
	case id:
		return nil, fmt.Errorf("missing account id")
	case key:
		return nil, fmt.Errorf("missing service account key")
	}

	sa := storage.SvcAccount{}
	if err := s.db.Get(&sa, validateSA, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, fmt.Errorf("executing validate SvcAccount get: %w", err)
	}
	h, err := hashChallenge(key, sa.Challenge[:saltLen])
	if err != nil {
		return nil, storage.NotFound
	}
	if h != sa.Challenge {
		return nil, storage.NotFound
	}
	sa.Challenge = ""
	return &sa, err
}

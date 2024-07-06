package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"github.com/lib/pq"
)

const insertUserProfile = `
INSERT INTO user_profile (
    org_id,
    user_id,
    email,
		profile_picture
) VALUES (
	 :org_id,
	 :user_id,
	 :email,
	 :profile_picture
) RETURNING
    id,created,updated
`

// CreateUserProfile creates new user profile and returns the created profile's ID.
func (s *Storage) CreateUserProfile(ctx context.Context, pf *storage.UserProfile) (string, error) {
	log := logging.FromContext(ctx)
	pstmt, err := s.db.PrepareNamedContext(ctx, insertUserProfile)
	if err != nil {
		logging.WithError(err, log).Error("insert user profile")
		return "", err
	}
	defer pstmt.Close()
	if err := pstmt.Get(pf, pf); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return "", storage.Conflict
		}
		return "", fmt.Errorf("executing user profile insert: %w", err)
	}
	return pf.ID, nil
}

const userProfileUpdateByUserID = `
UPDATE
	user_profile
SET
	profile_picture = :profile_picture,
	deleted= COALESCE(:deleted, deleted)
WHERE
	id = :id
RETURNING
   id,created,updated
`

// UpdateUserProfile updates the db values for a given user profile using the user ID
func (s *Storage) UpdateUserProfile(ctx context.Context, pf *storage.UserProfile) (string, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, userProfileUpdateByUserID)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	if err := stmt.Get(pf, pf); err != nil {
		if err == sql.ErrNoRows {
			return "", storage.NotFound
		}
		return "", fmt.Errorf("executing user profile update: %w", err)
	}
	return pf.ID, nil
}

// GetUserProfiles return all profiles
func (s *Storage) GetUserProfiles(ctx context.Context, oid string) ([]storage.UserProfile, error) {
	const getUserProfiles = `SELECT * FROM user_profile WHERE org_id = $1`
	var pfs []storage.UserProfile
	// sqlvet: ignore
	if err := s.db.Select(&pfs, getUserProfiles, oid); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return pfs, nil
}

// GetUserProfile return profile matched against user ID
func (s *Storage) GetUserProfile(ctx context.Context, id string) (*storage.UserProfile, error) {
	const userProfileID = `SELECT * FROM user_profile WHERE user_id = $1 and deleted is null`
	var pf storage.UserProfile
	if err := s.db.Get(&pf, userProfileID, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return &pf, nil
}

// GetUserProfileByEmail return profile matched against user email
func (s *Storage) GetUserProfileByEmail(ctx context.Context, email string) (*storage.UserProfile, error) {
	const userProfileEmail = `SELECT * FROM user_profile WHERE email = $1`
	var pf storage.UserProfile
	if err := s.db.Get(&pf, userProfileEmail, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return &pf, nil
}

const deleteUserProfile = `
UPDATE
	user_profile
SET
	deleted=NOW()
WHERE
	user_id = :id
`

// DeleteUserProfile delete the user profile by user id
func (s *Storage) DeleteUserProfile(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("user id is required")
	}
	r, err := s.db.NamedExecContext(ctx, deleteUserProfile, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}
	if c, err := r.RowsAffected(); err == nil && c == 0 {
		return nil
	}
	return nil
}

const enableUserProfile = `
UPDATE
	user_profile
SET
	deleted=null
WHERE
	user_id = :id
`

func (s *Storage) EnableUserProfile(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("user id is required")
	}
	r, err := s.db.NamedExecContext(ctx, enableUserProfile, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}
	if c, err := r.RowsAffected(); err == nil && c == 0 {
		return nil
	}
	return nil
}

const userProfileUpdateByOrgID = `
UPDATE
	user_profile
SET
	user_id = :user_id,
	org_id = :newOrgID
WHERE
	org_id = :oldOrgID
RETURNING
   id
`

func (s *Storage) UpdateUserProfileByOrgID(ctx context.Context, pf *storage.UpdateOrgProfileOrgIDUserID) (string, error) {
	if pf.OldOrgID == "" {
		return "", fmt.Errorf("invalid org id")
	}
	stmt, err := s.db.PrepareNamedContext(ctx, userProfileUpdateByOrgID)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var id string
	if err := stmt.Get(&id, pf); err != nil {
		if err == sql.ErrNoRows {
			return "", storage.NotFound
		}
		return "", fmt.Errorf("executing user profile update: %w", err)
	}
	return id, nil
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"brank.as/rbac/usermgm/storage"
)

const passwordResetInsert = `
INSERT INTO password_reset (
	"user_id",
	"expiry"
)
VALUES (
	:user_id,
	now() + :expiry * INTERVAL '1 second'
) RETURNING id
`

const passwordResetSelect = `
SELECT
   user_id,
   expiry,
   created
FROM password_reset
WHERE
	id = :id 
`

// CreatePasswordReset generate and store an password reset code
// return the created code
func (s *Storage) CreatePasswordReset(ctx context.Context, userID string, expiry time.Duration) (string, error) {
	var confCode string
	stmt, err := s.prepareNamed(ctx, passwordResetInsert)
	if err != nil {
		return "", err
	}
	arg := map[string]interface{}{
		"user_id": userID,
		"expiry":  expiry.Seconds(),
	}
	if err := stmt.Get(&confCode, arg); err != nil {
		if err == sql.ErrNoRows {
			return "", storage.NotFound
		}
		return "", fmt.Errorf("executing password reset insert: %w", err)
	}
	return confCode, nil
}

// GetUserIDByResetCode gets a password reset by reset code
func (s *Storage) GetResetCodeByID(ctx context.Context, resetCode string) (storage.PasswordReset, error) {
	var pr storage.PasswordReset
	stmt, err := s.db.PrepareNamedContext(ctx, passwordResetSelect)
	if err != nil {
		return pr, fmt.Errorf("preparing named query GetUserIDByResetCode: %w", err)
	}
	arg := map[string]interface{}{
		"id": resetCode,
	}
	if err := stmt.Get(&pr, arg); err != nil {
		if err == sql.ErrNoRows {
			return pr, storage.NotFound
		}
		return pr, err
	}
	return pr, nil
}

const passwordResetUpdate = `
WITH deleted AS(
	DELETE FROM password_reset
	WHERE id = :reset_code
	AND now() <= expiry
	RETURNING user_id
)
UPDATE user_account
SET password = :new_password
WHERE id = (SELECT user_id FROM deleted)
`

// PasswordReset validate code and set new password
func (s *Storage) PasswordReset(ctx context.Context, code, newPassword string) error {
	hashPw, err := s.hashPassword(newPassword, "")
	if err != nil {
		return err
	}

	arg := map[string]interface{}{
		"reset_code":   code,
		"new_password": hashPw,
	}
	res, err := s.db.NamedExecContext(ctx, passwordResetUpdate, arg)
	if err != nil {
		return fmt.Errorf("executing password reset update: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("password reset fetching number of affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return storage.NotFound
	}
	return nil
}

const setPasswordUpdate = `
UPDATE user_account
SET password = :new_password
WHERE invite_code = :invitation_code
`

// SetPassword sets new password for a user
func (s *Storage) SetPassword(ctx context.Context, inviteCode, password string) error {
	hashPw, err := s.hashPassword(password, "")
	if err != nil {
		return err
	}

	arg := map[string]interface{}{
		"invitation_code": inviteCode,
		"new_password":    hashPw,
	}
	res, err := sqlx.NamedExecContext(ctx, s.execer(ctx), setPasswordUpdate, arg)
	if err != nil {
		return fmt.Errorf("executing set password update: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("set password fetching number of affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return storage.NotFound
	}
	return nil
}

// SetUsername sets new username for a user.
func (s *Storage) SetUsername(ctx context.Context, inviteCode, username string) error {
	const setUsernameUpdate = `UPDATE user_account SET username = $1
WHERE invite_code = $2
RETURNING username`

	tx, err := s.NewTransacton(ctx)
	if err != nil {
		return err
	}
	defer s.Rollback(tx)
	const countUsername = `SELECT COUNT(*) FROM user_account WHERE username = $1`
	count := 0
	if err := sqlx.GetContext(tx, s.queryer(tx), &count, countUsername, username); err != nil {
		return fmt.Errorf("executing get username count: %w", err)
	}
	if count != 0 {
		return storage.UsernameExists
	}

	usr := ""
	if err := sqlx.GetContext(tx, s.queryer(tx), &usr, setUsernameUpdate, username, inviteCode); err != nil {
		if err == sql.ErrNoRows {
			return storage.NotFound
		}
		return fmt.Errorf("executing set username update: %w", err)
	}
	if usr != username {
		return fmt.Errorf("failed to update username")
	}
	return s.Commit(tx)
}

const updatePasswordByID = `
UPDATE user_account
SET password = :new_password
WHERE id = :id
`

// SetPassword sets new password for a user
func (s *Storage) SetPasswordByID(ctx context.Context, id, newPassword string) error {
	hashPw, err := s.hashPassword(newPassword, "")
	if err != nil {
		return err
	}

	arg := map[string]interface{}{
		"id":           id,
		"new_password": hashPw,
	}
	res, err := s.db.NamedExecContext(ctx, updatePasswordByID, arg)
	if err != nil {
		return fmt.Errorf("executing set password update: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("set password fetching number of affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return storage.NotFound
	}
	return nil
}

const passwordChangeInsert = `
INSERT INTO change_password (
	"user_id",
	"event_id",
	"new_password"
)
VALUES (
	:user_id,
	:event_id,
	:new_password
) RETURNING id
`

// CreateChangePassword store a change password record
func (s *Storage) CreateChangePassword(ctx context.Context, userID, eventID, newPass string) (string, error) {
	var id string

	hashNewPw, err := s.hashPassword(newPass, "")
	if err != nil {
		return "", err
	}

	stmt, err := s.db.PrepareNamedContext(ctx, passwordChangeInsert)
	if err != nil {
		return "", err
	}
	arg := map[string]interface{}{
		"user_id":      userID,
		"event_id":     eventID,
		"new_password": hashNewPw,
	}
	if err := stmt.Get(&id, arg); err != nil {
		if err == sql.ErrNoRows {
			return "", storage.NotFound
		}
		return "", fmt.Errorf("executing change password insert: %w", err)
	}
	return id, nil
}

const passwordChangeUpdate = `
WITH password AS(
	DELETE FROM change_password
	WHERE event_id = :event_id
	RETURNING user_id, new_password
)
UPDATE user_account
SET password = (SELECT new_password FROM password)
WHERE id = (SELECT user_id FROM password)
`

// ChangePassword execute change password
func (s *Storage) ChangePassword(ctx context.Context, userID, eventID string) error {
	arg := map[string]interface{}{
		"user_id":  userID,
		"event_id": eventID,
	}
	res, err := s.db.NamedExecContext(ctx, passwordChangeUpdate, arg)
	if err != nil {
		return fmt.Errorf("executing change password update: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("password change fetching number of affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return storage.NotFound
	}
	return nil
}

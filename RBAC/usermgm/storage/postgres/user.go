package postgres

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/crypto/argon2"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/svcutil/random"
	"brank.as/rbac/usermgm/storage"
)

const saltLen = 32

func (s *Storage) hashPassword(password, salt string) (string, error) {
	if password == "" {
		pw, err := random.String(20)
		if err != nil {
			return "", err
		}
		password = pw
	}
	if salt == "" {
		st, err := random.String(saltLen)
		if err != nil {
			return "", err
		}
		salt = st
	}
	pwHash := argon2.IDKey([]byte(password), []byte(salt), 5, 16*1024, 4, 32)
	return salt + hex.EncodeToString(pwHash), nil
}

const userInsert = `
INSERT INTO user_account (
	email,
	username,
	password,
	first_name,
	last_name,
	org_id,
	invite_status,
	invite_code,
	invite_expiry
) VALUES (
	:email,
	:username,
	'',
	:first_name,
	:last_name,
	:org_id,
	:invite_status,
	:invite_code,
	:invite_expiry
) RETURNING 
	id,
	org_id,
	username,
	first_name,
	last_name,
	email,
	email_verified,
	invite_code,
	invite_status,
	invite_expiry,
	preferred_mfa,
	mfa_login,
	created,
	updated
`

// CreateUser creates new user account returns the created user's ID.
func (s *Storage) CreateUser(ctx context.Context, user storage.User, cred storage.Credential) (*storage.User, error) {
	const setpass = `UPDATE user_account SET password=:password where id=:id`
	pwd, err := s.hashPassword(cred.Password, "")
	if err != nil {
		return nil, err
	}
	var success bool
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		if success {
			return
		}
		tx.Rollback()
	}()
	stmt, err := tx.PrepareNamed(userInsert)
	if err != nil {
		return nil, err
	}
	pw, err := tx.PrepareNamed(setpass)
	if err != nil {
		return nil, err
	}

	if err := stmt.Get(&user, user); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			switch pErr.Constraint {
			case userEmailDup:
				return nil, storage.EmailExists
			case usernameDup:
				return nil, storage.UsernameExists
			case usrInvCodeDup:
				return nil, storage.InvCodeExists
			}
		}
		return nil, fmt.Errorf("executing user insert: %w", err)
	}
	cred.ID = user.ID
	cred.Password = string(pwd)
	if _, err := pw.Exec(cred); err != nil {
		return nil, fmt.Errorf("executing password insert: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	success = true
	return &user, nil
}

const userSelectAuth = `
select
	id,
	org_id,
	mfa_login,
	preferred_mfa,
	password,
	locked,
	last_login,
	last_failed,
	fail_count,
	deleted,
	email_verified
FROM user_account
WHERE
	username = $1
`

func (s *Storage) GetUser(ctx context.Context, username, password string) (*storage.User, error) {
	cred := &storage.Credential{}
	if err := sqlx.GetContext(ctx, s.queryer(ctx), cred, userSelectAuth, username); err != nil {
		return nil, storage.NotFound
	}
	if cred.Password == "" || len(cred.Password) < 2*saltLen {
		return nil, storage.NotFound
	}
	pwH, _ := s.hashPassword(password, cred.Password[:saltLen])
	if !equalConstTime(pwH, cred.Password) {
		tm, ct, err := s.loginFail(ctx, cred.ID)
		if err != nil {
			return nil, storage.NotFound
		}
		return &storage.User{
			ID:         cred.ID,
			Deleted:    cred.Deleted,
			Locked:     cred.Locked,
			LastLogin:  cred.LastLogin,
			LastFailed: sql.NullTime{Time: tm, Valid: true},
			FailCount:  ct,
		}, storage.NotFound
	}

	tm, err := s.loginSuccess(ctx, cred.ID)
	if err != nil {
		logging.FromContext(ctx).WithField("method", "postgres.loginsuccess").
			WithError(err).Error("record failed")
	}
	return &storage.User{
		ID:            cred.ID,
		OrgID:         cred.OrgID,
		Username:      username,
		PreferredMFA:  cred.PreferredMFA,
		MFALogin:      cred.MFALogin,
		Deleted:       cred.Deleted,
		Locked:        cred.Locked,
		LastLogin:     sql.NullTime{Time: tm, Valid: true},
		LastFailed:    cred.LastFailed,
		FailCount:     cred.FailCount,
		EmailVerified: cred.EmailVerified,
	}, nil
}

// loginSuccess recorded for user id.
func (s *Storage) loginSuccess(ctx context.Context, id string) (time.Time, error) {
	const successSQL = `UPDATE user_account SET fail_count=0, last_login=now() WHERE id = $1 RETURNING last_login`
	tm := time.Time{}
	err := sqlx.GetContext(ctx, s.queryer(ctx), &tm, successSQL, id)
	return tm, err
}

// loginFail recorded for user id.
func (s *Storage) loginFail(ctx context.Context, id string) (time.Time, int, error) {
	const failedSQL = `UPDATE user_account SET fail_count=fail_count+1, last_failed=now() WHERE id = $1 RETURNING last_failed, fail_count`
	cred := &storage.Credential{}
	err := sqlx.GetContext(ctx, s.queryer(ctx), cred, failedSQL, id)
	return cred.LastFailed.Time, cred.FailCount, err
}

// Locks user account.
func (s *Storage) LockUser(ctx context.Context, id string) error {
	const lockUser = `UPDATE user_account SET locked=now(), reset_required=now() where id = :id`
	if id == "" {
		return fmt.Errorf("user id is required")
	}
	r, err := sqlx.NamedExecContext(ctx, s.execer(ctx), lockUser, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}
	if c, err := r.RowsAffected(); err == nil && c == 0 {
		return storage.NotFound
	}
	return nil
}

// Unlocks user account.
func (s *Storage) UnlockUser(ctx context.Context, id string) error {
	const unlockUser = `UPDATE user_account SET locked=NULL, fail_count=0 where id=:id`
	if id == "" {
		return fmt.Errorf("user id is required")
	}
	r, err := sqlx.NamedExecContext(ctx, s.execer(ctx), unlockUser, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}
	if c, err := r.RowsAffected(); err == nil && c == 0 {
		return storage.NotFound
	}
	return nil
}

func (s *Storage) ValidateUserPass(ctx context.Context, id, password string) (string, string, error) {
	const userAuth = `SELECT id, org_id, password FROM user_account WHERE id = $1`
	cred := &storage.Credential{}
	if err := s.db.GetContext(ctx, cred, userAuth, id); err != nil {
		return "", "", storage.NotFound
	}
	if cred.Password == "" || len(cred.Password) < 2*saltLen {
		return "", "", storage.NotFound
	}
	pwH, _ := s.hashPassword(password, cred.Password[:saltLen])
	if !equalConstTime(pwH, cred.Password) {
		return "", "", storage.NotFound
	}
	return cred.ID, cred.OrgID, nil
}

const userSelectByID = `
SELECT
	id,
	org_id,
	username,
	first_name,
	last_name,
	email,
	email_verified,
	invite_code,
	invite_status,
	invite_expiry,
	preferred_mfa,
	mfa_login,
	reset_required,
	created,
	updated
FROM user_account
WHERE
	id = :id
`

// GetUserByID return user account matched against ID primary key column
func (s *Storage) GetUserByID(ctx context.Context, id string) (*storage.User, error) {
	var user storage.User
	stmt, err := s.db.PrepareNamed(userSelectByID)
	if err != nil {
		return nil, fmt.Errorf("preparing named query GetUserByID: %w", err)
	}
	arg := map[string]interface{}{
		"id": id,
	}
	if err := stmt.Get(&user, arg); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return &user, nil
}

const userSelectByEmail = `
SELECT
	id,
	org_id,
	username,
	first_name,
	last_name,
	email,
	email_verified,
	preferred_mfa,
	mfa_login,
	invite_status,
	invite_code,
	invite_expiry,
	reset_required,
	deleted,
	created,
	updated
FROM user_account
WHERE
	email = $1
`

// GetUserByEmail return user account matched against ID primary key column
func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*storage.User, error) {
	var user storage.User
	if err := s.db.GetContext(ctx, &user, userSelectByEmail, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return &user, nil
}

const userSelectByInvite = `
SELECT
	id,
	org_id,
	username,
	first_name,
	last_name,
	email,
	email_verified,
	created,
	updated,
	invite_status,
	invite_code,
	invite_expiry,
	reset_required
FROM user_account
WHERE
	invite_code = :invite_code
`

// GetUserByInvite return user account matched against ID primary key column
func (s *Storage) GetUserByInvite(ctx context.Context, code string) (*storage.User, error) {
	var user storage.User
	stmt, err := s.db.PrepareNamed(userSelectByInvite)
	if err != nil {
		return nil, fmt.Errorf("preparing named query GetUserByInvite: %w", err)
	}
	arg := map[string]interface{}{
		"invite_code": code,
	}
	if err := stmt.Get(&user, arg); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}

	exp := time.Now().After(user.InviteExpiry)
	if user.InviteStatus == storage.InviteSent && exp {
		user.InviteStatus = storage.Expired
		if _, err := s.UpdateUserByID(ctx, user); err != nil {
			return nil, err
		}
	}
	return &user, nil
}

const usersSelectByOrg = `
SELECT
	id,
	org_id,
	first_name,
	last_name,
	email,
	email_verified,
	invite_code,
	invite_status,
	invite_expiry,
	preferred_mfa,
	mfa_login,
	created,
	updated
FROM user_account
WHERE
	org_id = $1
`

const userDeleteExpired = `
UPDATE 
	user_account
SET 	
	deleted= :deleted 
WHERE
	invite_status = :invite_status
`

const usersSelectExpiredInv = `
SELECT
	id,
	org_id,
	name,
	first_name,
	last_name,
	email,
	email_verified,
	role,
	invite_status,
	invite_code,
	invite_expiry,
	created,
	updated,
FROM user_account
WHERE
	invite_status IN ($1, $2) AND deleted IS NULL
`

const userUpdateByID = `
UPDATE 
	user_account
SET 	
	first_name= COALESCE(NULLIF(:first_name, ''), first_name),
	last_name= COALESCE(NULLIF(:last_name, ''), last_name),
	email= COALESCE(NULLIF(:email, ''), email),
	email_verified= COALESCE(NULLIF(:email_verified, FALSE), email_verified),
	invite_status= COALESCE(NULLIF(:invite_status, ''), invite_status),
	invite_code= COALESCE(NULLIF(:invite_code, ''), invite_code),
	invite_expiry= CASE
           WHEN :invite_expiry > invite_expiry then :invite_expiry
           ELSE invite_expiry 
       END,
	mfa_login= COALESCE(NULLIF(:mfa_login, FALSE), mfa_login),
	preferred_mfa= COALESCE(NULLIF(:preferred_mfa, ''), preferred_mfa),
	deleted= COALESCE(:deleted, deleted)
WHERE
	id = :id
RETURNING
	org_id,
	username,
	first_name,
	last_name,
	email,
	email_verified,
	mfa_login,
	preferred_mfa,
	invite_status,
	invite_code,
	invite_expiry,
	deleted,
	created,
	updated
`

// UpdateUserByID updates the db values for a given user using the ID as primary key
func (s *Storage) UpdateUserByID(ctx context.Context, user storage.User) (*storage.User, error) {
	if user.ID == "" {
		return nil, fmt.Errorf("invalid user id is required by UpdateUserByID")
	}

	stmt, err := s.db.PrepareNamedContext(ctx, userUpdateByID)
	if err != nil {
		return nil, fmt.Errorf("preparing named query userUpdateByID: %w", err)
	}
	defer stmt.Close()

	if err := stmt.Get(&user, user); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique && pErr.Constraint == usrInvCodeDup {
			return nil, storage.InvCodeExists
		}
		return nil, fmt.Errorf("executing query: %w", err)
	}

	return &user, nil
}

// DeleteExpiredUsers deletes all expired users
func (s *Storage) DeleteExpiredUsers(ctx context.Context) error {
	_, err := s.db.NamedExecContext(ctx, userDeleteExpired, map[string]interface{}{
		"deleted":       time.Now(),
		"invite_status": storage.Expired,
	})
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}
	return nil
}

// GetUsersByOrg return user account matched against ID primary key column
func (s *Storage) GetUsersByOrg(ctx context.Context, orgID string) ([]storage.User, error) {
	usrs := []storage.User{}
	if err := s.db.SelectContext(ctx, &usrs, usersSelectByOrg, orgID); err != nil {
		if err == sql.ErrNoRows {
			return usrs, storage.NotFound
		}
		return usrs, err
	}
	return usrs, nil
}

const emailConfirmationInsert = `
INSERT INTO email_confirmation (
	user_id
) VALUES (
	:user_id
) RETURNING id
`

// CreateConfirmationCode generate and store an email confirmation code
// return the created verification code
func (s *Storage) CreateConfirmationCode(ctx context.Context, userID string) (string, error) {
	var confCode string
	stmt, err := s.db.PrepareNamedContext(ctx, emailConfirmationInsert)
	if err != nil {
		return "", err
	}
	arg := map[string]interface{}{
		"user_id": userID,
	}
	if err := stmt.Get(&confCode, arg); err != nil {
		if err == sql.ErrNoRows {
			return "", storage.NotFound
		}
		return "", fmt.Errorf("executing confirm code insert: %w", err)
	}
	return confCode, nil
}

const verifyConfirmationCodeUpdate = `
WITH deleted AS (
	DELETE FROM email_confirmation
	WHERE id = $1
	RETURNING user_id
)
UPDATE user_account
SET email_verified = true
WHERE id = (SELECT user_id FROM deleted)
RETURNING
	id,
	org_id,
	username,
	first_name,
	last_name,
	email;
`

// VerifyConfirmationCode set user_account.email_verified if confirmation code
// exists and return email address of the verified account
func (s *Storage) VerifyConfirmationCode(ctx context.Context, code string) (*storage.User, error) {
	u := storage.User{}
	if err := s.db.Get(&u, verifyConfirmationCodeUpdate, code); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, fmt.Errorf("executing signup user insert: %w", err)
	}

	return &u, nil
}

const confirmCodeSelect = `
select
	id
FROM email_confirmation
WHERE
	user_id = $1
`

func (s *Storage) GetConfirmationCode(ctx context.Context, uid string) (string, error) {
	var code string
	if err := s.db.Get(&code, confirmCodeSelect, uid); err != nil {
		if err == sql.ErrNoRows {
			return "", storage.NotFound
		}
		return "", err
	}
	return code, nil
}

// GetExpiredUserInvites return expired user invites
func (s *Storage) GetExpiredUserInvites(ctx context.Context) ([]storage.User, error) {
	usrs := []storage.User{}
	if err := s.db.SelectContext(ctx, &usrs, usersSelectExpiredInv, storage.Expired, storage.InviteSent); err != nil {
		if err == sql.ErrNoRows {
			return usrs, storage.NotFound
		}
		return usrs, err
	}

	expUsers := []storage.User{}
	for _, usr := range usrs {
		if usr.InviteStatus == storage.Expired {
			expUsers = append(expUsers, usr)
		}
		exp := time.Now().After(usr.InviteExpiry)
		if usr.InviteStatus == storage.InviteSent && exp {
			usr.InviteStatus = storage.Expired
			expUsers = append(expUsers, usr)
			if _, err := s.UpdateUserByID(ctx, usr); err != nil {
				return usrs, err
			}
		}
	}
	return expUsers, nil
}

// ReviveUserByID removes deletion on user by id
func (s *Storage) ReviveUserByID(ctx context.Context, id string) error {
	const usrReviveByID = `UPDATE user_account SET deleted = NULL, reset_required = now() WHERE id = :id`
	if id == "" {
		return fmt.Errorf("invalid user id is required by ReviveUserByID")
	}

	r, err := sqlx.NamedExecContext(ctx, s.execer(ctx), usrReviveByID, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}
	if c, err := r.RowsAffected(); err == nil && c == 0 {
		return storage.NotFound
	}
	return nil
}

// DeleteUserByID deletes a user by ID
func (s *Storage) DisableUser(ctx context.Context, id string) error {
	const disableUser = `UPDATE user_account SET deleted=now(), reset_required=now() where id = :id`
	if id == "" {
		return fmt.Errorf("user id is required")
	}
	r, err := sqlx.NamedExecContext(ctx, s.execer(ctx), disableUser, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}
	if c, err := r.RowsAffected(); err == nil && c == 0 {
		return storage.NotFound
	}
	return nil
}

const getUsersFilter = `
WITH cnt AS (select count(*) as count FROM user_account WHERE org_id = ?)
SELECT
	id,
	org_id,
	first_name,
	last_name,
	username,
	email,
	email_verified,
	invite_code,
	invite_status,
	invite_expiry,
	mfa_login,
	created,
	updated,
	deleted,
	cnt.count
FROM user_account as p left join cnt on true 
WHERE
	org_id = ? 
`

// Get users with filter
func (s *Storage) GetUsers(ctx context.Context, f storage.FilterList) ([]storage.User, error) {
	searchQ := getUsersFilter
	inp := []interface{}{f.OrgID}
	inp = append(inp, f.OrgID)
	if f.Name != "" {
		names := strings.Split(f.Name, " ")
		searchQL := []string{}
		for _, n := range names {
			searchQL = append(searchQL, " (first_name ILIKE ? ) OR (last_name ILIKE ?) ")
			nm := fmt.Sprintf("%%%s%%", n)
			inp = append(inp, nm)
			inp = append(inp, nm)
		}
		searchQ += " AND (" + strings.Join(searchQL, " OR ") + ") "
	}
	if f.Status != nil {
		searchQ += " AND invite_status IN(?)"
		inp = append(inp, f.Status)
	}
	if len(f.ID) != 0 {
		searchQ += " AND id IN(?)"
		inp = append(inp, f.ID)
	}
	if f.SortBy != "ASC" { // default to descending on empty or invalid input
		f.SortBy = "DESC"
	}
	// DB query string cannot process sort column/order as prepared statement parameters.
	// Enforce acceptable sort columns explicitly.
	sortCol := map[string]string{
		"CreatedDate": "created",
		"UserName":    "first_name",
		"created":     "created",
		"first_name":  "first_name",
	}
	if sb := sortCol[f.SortByColumn]; sb != "" {
		searchQ += fmt.Sprintf(" ORDER BY %s %s", sb, f.SortBy)
	}
	if f.Limit > 0 {
		searchQ += " LIMIT ?"
		inp = append(inp, f.Limit)
	}
	if f.Offset > 0 {
		searchQ += " OFFSET ?"
		inp = append(inp, f.Offset)
	}

	fullQuery, args, err := sqlx.In(searchQ, inp...)
	if err != nil {
		return nil, err
	}

	var usrs []storage.User
	if err := s.db.Select(&usrs, s.db.Rebind(fullQuery), args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return usrs, nil
}

// EnableUsrByID removes deletion on user by id
func (s *Storage) EnableUser(ctx context.Context, id string) error {
	const enableUsrByID = `UPDATE user_account SET deleted = NULL WHERE id = :id`
	if id == "" {
		return fmt.Errorf("user id is required by EnableUserByID")
	}

	r, err := s.db.NamedExecContext(ctx, enableUsrByID, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}
	if c, err := r.RowsAffected(); err == nil && c == 0 {
		return storage.NotFound
	}
	return nil
}

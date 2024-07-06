package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"brank.as/rbac/profile/storage"
)

const insertUser = `
INSERT INTO user_account (
    org_id
) VALUES (
    :org_id
) RETURNING
    id,created,updated
`

// CreateUser creates new user account returns the created user's ID.
func (s *Storage) CreateUser(ctx context.Context, u storage.User) (string, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, insertUser)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	if err := stmt.Get(&u, u); err != nil {
		return "", fmt.Errorf("executing user insert: %w", err)
	}
	return u.ID, nil
}

const userSelectByID = `
SELECT
	id,
	org_id,
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

const usersSelectByOrg = `
SELECT
	id,
	org_id,
	created,
	updated
FROM user_account
WHERE
	org_id = $1
`

// GetUsersByOrg returns user accounts associated with an org ID
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

const deleteUser = `
UPDATE user_account 
SET
deleted = :deleted 
WHERE id = :id
RETURNING *
`

// DeleteUserByID deletes a user by ID
func (s *Storage) DeleteUserByID(ctx context.Context, id string) (*storage.User, error) {
	u := storage.User{
		ID: id,
		Deleted: sql.NullTime{
			Valid: true,
			Time:  time.Now().UTC(),
		},
	}
	stmt, err := s.db.PrepareNamedContext(ctx, deleteUser)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&u, &u); err != nil {
		return nil, fmt.Errorf("executing user delete: %w", err)
	}
	return &u, nil
}

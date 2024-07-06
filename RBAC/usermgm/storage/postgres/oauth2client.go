package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"brank.as/rbac/usermgm/storage"
)

const oauthClientInsert = `
INSERT INTO oauth2_client (
	org_id,
	client_id,
	client_name,
	environment,
    updated_user_id,
	created_user_id
) VALUES (
	:org_id,
	:client_id,
	:client_name,
	:environment,
	:created_user_id,
	:created_user_id
)
RETURNING created, updated`

// CreateOauthClient creates new OauthClient account returns the created OauthClient's ID.
func (s *Storage) CreateOauthClient(ctx context.Context, sa storage.OAuthClient) (*storage.OAuthClient, error) {
	switch "" {
	case sa.ClientID, sa.CreateUserID, sa.ClientName, sa.OrgID:
		return nil, fmt.Errorf("invalid account client_id, client_name, org_id, and create_user_id are required")
	}
	stmt, err := s.db.PrepareNamedContext(ctx, oauthClientInsert)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&sa, sa); err != nil {
		return nil, fmt.Errorf("executing OauthClient insert: %w", err)
	}
	return &sa, nil
}

const oauthClientUpdate = `
UPDATE oauth2_client SET
    updated_user_id = :updated_user_id
WHERE
    client_id = :client_id
RETURNING created, updated`

// UpdateOauthClient updates OauthClient account returns the updated OauthClient's ID.
func (s *Storage) UpdateOauthClient(ctx context.Context, sa storage.OAuthClient) (*storage.OAuthClient, error) {
	switch "" {
	case sa.ClientID, sa.UpdateUserID:
		return nil, fmt.Errorf("invalid account client_id, or update_user_id are required")
	}
	stmt, err := s.db.PrepareNamedContext(ctx, oauthClientUpdate)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&sa, sa); err != nil {
		return nil, fmt.Errorf("executing OauthClient update: %w", err)
	}
	return &sa, nil
}

// GetOauthClientByID return OauthClient account matched against ID primary key column
func (s *Storage) GetOauthClientByID(ctx context.Context, clientID string) (*storage.OAuthClient, error) {
	const getOauthClient = `SELECT * FROM oauth2_client WHERE client_id = :client_id`
	stmt, err := s.db.PrepareNamedContext(ctx, getOauthClient)
	if err != nil {
		return nil, fmt.Errorf("preparing named query GetOauthClientByID: %w", err)
	}
	defer stmt.Close()
	sa := &storage.OAuthClient{ClientID: clientID}
	if err := stmt.Get(sa, sa); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return sa, nil
}

// GetOauthClientByOrgID return OauthClient account associated with the given org.
func (s *Storage) GetOauthClientByOrgID(ctx context.Context, orgID string, listDisable bool) ([]storage.OAuthClient, error) {
	getOauthClientByOrg := `
	SELECT * FROM oauth2_client
	WHERE org_id=:org_id`

	if !listDisable {
		getOauthClientByOrg += ` AND deleted IS NULL`
	}

	getOauthClientByOrg += ` ORDER BY created desc`
	stmt, err := s.db.PrepareNamedContext(ctx, getOauthClientByOrg)
	if err != nil {
		return nil, fmt.Errorf("preparing named query GetOauthClientByID: %w", err)
	}
	defer stmt.Close()
	sa := storage.OAuthClient{OrgID: orgID}
	var lst []storage.OAuthClient
	if err := stmt.Select(&lst, &sa); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return lst, nil
}

const disableOauthClientUpdate = `
UPDATE oauth2_client SET
    deleted_user_id=:deleted_user_id,
    deleted=NOW()
WHERE
    client_id=:client_id
RETURNING deleted`

// DeleteOauthClient permanently disables the service account credential.
func (s *Storage) DeleteOauthClient(ctx context.Context, sa storage.OAuthClient) (*time.Time, error) {
	switch "" {
	case sa.ClientID:
		return nil, fmt.Errorf("missing client id")
	case sa.DeleteUserID:
		return nil, fmt.Errorf("missing user id for disable action")
	}
	stmt, err := s.db.PrepareNamedContext(ctx, disableOauthClientUpdate)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&sa, &sa); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, fmt.Errorf("executing signup OauthClient insert: %w", err)
	}

	return &sa.Deleted.Time, nil
}

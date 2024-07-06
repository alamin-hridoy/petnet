package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"brank.as/rbac/usermgm/storage"
)

const orgInsert = `
INSERT INTO organization_information (
	org_name,
	contact_email,
	contact_phone,
	active
)
VALUES (
	:org_name,
	:contact_email,
	:contact_phone,
	:active
) RETURNING id
`

// CreateOrg creates new Organization returns the created Organization's ID.
func (s *Storage) CreateOrg(ctx context.Context, org storage.Organization) (string, error) {
	if org.OrgName == "" {
		return "", fmt.Errorf("invalid organization org_name is required")
	}

	stmt, err := s.db.PrepareNamedContext(ctx, orgInsert)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var orgID string
	if err := stmt.Get(&orgID, org); err != nil {
		return "", fmt.Errorf("executing organization insert: %w", err)
	}
	return orgID, nil
}

const orgSelectByID = `
SELECT
	id,
	org_name,
	contact_email,
	contact_phone,
	active,
	mfa_login,
	created,
	updated,
	deleted
FROM organization_information
WHERE
	id = :id
`

// GetOrgByID return org matched against ID primary key column
func (s *Storage) GetOrgByID(ctx context.Context, id string) (*storage.Organization, error) {
	var org storage.Organization
	stmt, err := s.db.PrepareNamedContext(ctx, orgSelectByID)
	if err != nil {
		return nil, fmt.Errorf("preparing named query GetOrgByID: %w", err)
	}
	arg := map[string]interface{}{
		"id": id,
	}
	if err := stmt.Get(&org, arg); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}

	return &org, nil
}

const orgList = `
SELECT
	id,
	org_name,
	contact_email,
	contact_phone,
	active,
	mfa_login,
	created,
	updated,
	deleted
FROM organization_information
ORDER BY created DESC
`

// GetOrgs return org matched against ID primary key column
func (s *Storage) GetOrgs(ctx context.Context) ([]storage.Organization, error) {
	var orgs []storage.Organization
	if err := s.db.SelectContext(ctx, &orgs, orgList); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}

	return orgs, nil
}

const orgUpdateByID = `
UPDATE 
	organization_information
SET (
	org_name,
	active,
	mfa_login,
	contact_email,
	contact_phone
) = (
	COALESCE(NULLIF(:org_name, ''), org_name),
	COALESCE(NULLIF(:active, FALSE), active),
	COALESCE(:mfa_login, mfa_login),
	COALESCE(NULLIF(:contact_email, ''), contact_email),
	COALESCE(NULLIF(:contact_phone, ''), contact_phone)
)  
WHERE
	id = :id
RETURNING *
`

// UpdateOrgByID updates the db values for a given organization using the ID as primary key
func (s *Storage) UpdateOrgByID(ctx context.Context, org storage.Organization) (*storage.Organization, error) {
	if org.ID == "" {
		return nil, fmt.Errorf("invalid organization id is required by UpdateOrgByID")
	}
	stmt, err := s.db.PrepareNamedContext(ctx, orgUpdateByID)
	if err != nil {
		return nil, fmt.Errorf("preparing query: %w", err)
	}
	defer stmt.Close()

	if err := stmt.Get(&org, org); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return &org, nil
}

const orgRevive = `
UPDATE 
	organization_information
SET 	
	active = TRUE,
	deleted = NULL
WHERE
	id = :id
`

// ReviveOrg removes deletion on org by id
func (s *Storage) ReviveOrgByID(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("invalid organization id is required by ReviveOrgByID")
	}

	r, err := s.db.NamedExecContext(ctx, orgRevive, map[string]interface{}{
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

const orgActivate = `
UPDATE 
	organization_information
SET 	
	active = TRUE,
	deleted = NULL
WHERE
	id = $1
`

func (s *Storage) ActivateOrg(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("invalid organization id is required by ActivateOrg")
	}
	r, err := s.db.ExecContext(ctx, orgActivate, id)
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}
	if c, err := r.RowsAffected(); err == nil && c == 0 {
		return storage.NotFound
	}
	return nil
}

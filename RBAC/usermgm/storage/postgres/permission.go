package postgres

import (
	"context"
	"fmt"

	"brank.as/rbac/usermgm/storage"
)

const createPermission = `
INSERT INTO permissions (
    id,
    org_id,
    service_permission_id,
    permission_name,
    description,
    create_user_id
) VALUES (
    :id,
    :org_id,
    :service_permission_id,
    :permission_name,
    :description,
    :create_user_id
) RETURNING
    created,updated
`

func (s *Storage) CreatePermission(ctx context.Context, p storage.Permission) (*storage.Permission, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, createPermission)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&p, p); err != nil {
		return nil, fmt.Errorf("executing permission insert: %w", err)
	}
	return &p, nil
}

const getPermission = `
SELECT *
FROM permissions
WHERE id = $1
`

func (s *Storage) GetPermission(ctx context.Context, id string) (*storage.Permission, error) {
	var p storage.Permission
	if err := s.db.Get(&p, getPermission, id); err != nil {
		return nil, fmt.Errorf("executing permission details: %w", err)
	}
	return &p, nil
}

const deletePermission = `
UPDATE permissions
SET
(
    delete_user_id,
    deleted
) = (
    :delete_user_id,
    NOW()
) WHERE id = :id
RETURNING *
`

func (s *Storage) DeletePermission(ctx context.Context, p storage.Permission) (*storage.Permission, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, deletePermission)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&p, &p); err != nil {
		return nil, fmt.Errorf("executing permission delete: %w", err)
	}
	return &p, nil
}

const updatePermission = `
UPDATE permissions SET
    description = :description
WHERE id = :id
RETURNING updated
`

func (s *Storage) UpdatePermission(ctx context.Context, p storage.Permission) (*storage.Permission, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updatePermission)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&p, p); err != nil {
		return nil, fmt.Errorf("executing permission update: %w", err)
	}
	return &p, nil
}

const listPermission = `
SELECT
 * 
FROM permissions 
WHERE org_id = $1
`

func (s *Storage) ListPermission(ctx context.Context, orgID string) ([]storage.Permission, error) {
	var perms []storage.Permission
	if err := s.db.Select(&perms, listPermission, orgID); err != nil {
		return nil, fmt.Errorf("executing permission list: %w", err)
	}
	return perms, nil
}

const listPermBySvcPermID = `
SELECT
 * 
FROM permissions 
WHERE service_permission_id = $1
`

func (s *Storage) ListPermBySvcPermID(ctx context.Context, svcPermID string) ([]storage.Permission, error) {
	var perms []storage.Permission
	if err := s.db.Select(&perms, listPermBySvcPermID, svcPermID); err != nil {
		return nil, fmt.Errorf("executing permission list: %w", err)
	}
	return perms, nil
}

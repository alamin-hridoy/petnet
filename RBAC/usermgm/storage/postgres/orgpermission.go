package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"brank.as/rbac/usermgm/storage"
)

const upsertOrgPermission = `
INSERT INTO org_permission(
	id,
	org_id,
	grant_id,
	service_id,
	name,
	description,
	resource,
	action,
	environment,
	permission_id
) VALUES (
	:id,
	:org_id,
	:grant_id,
	:service_id,
	:name,
	:description,
	:resource,
	:action,
	:environment,
	:permission_id
)
RETURNING id;
`

func (s *Storage) UpsertOrgPermission(ctx context.Context, ps storage.OrgPermission) (string, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, upsertOrgPermission)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var id string
	if err := stmt.Get(&id, ps); err != nil {
		return "", err
	}

	return id, nil
}

func (s *Storage) GetOrgPermission(ctx context.Context, id string) (*storage.OrgPermission, error) {
	const getOrgPermission = ` SELECT * FROM org_permission WHERE id = $1`
	var ps storage.OrgPermission
	if err := s.db.Get(&ps, getOrgPermission, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, fmt.Errorf("executing org_permission permission details: %w", err)
	}
	return &ps, nil
}

func (s *Storage) ListOrgPermission(ctx context.Context, org string) ([]storage.OrgPermission, error) {
	const getOrgPermission = `SELECT * FROM org_permission WHERE org_id = $1`
	ps := []storage.OrgPermission{}
	if err := s.db.Select(&ps, getOrgPermission, org); err != nil {
		return nil, err
	}
	return ps, nil
}

func (s *Storage) ListOrgPermissionID(ctx context.Context, permission string) ([]storage.OrgPermission, error) {
	const getOrgPermission = ` SELECT * FROM org_permission WHERE permission_id = $1`
	ps := []storage.OrgPermission{}
	if err := s.db.Select(&ps, getOrgPermission, permission); err != nil {
		return nil, err
	}
	return ps, nil
}

func (s *Storage) ListOrgPermissionGrant(ctx context.Context, grant string) ([]storage.OrgPermission, error) {
	const getOrgPermission = ` SELECT * FROM org_permission WHERE grant_id = $1`
	ps := []storage.OrgPermission{}
	if err := s.db.Select(&ps, getOrgPermission, grant); err != nil {
		return nil, err
	}
	return ps, nil
}

const deleteOrgPermission = `
DELETE from org_permission
WHERE id = $1
`

func (s *Storage) DeleteOrgPermission(ctx context.Context, id string) error {
	if res, err := s.db.Exec(deleteOrgPermission, id); err != nil {
		return fmt.Errorf("executing org_permission permission delete: %w", err)
	} else if _, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("executing org_permission permission delete: %w", err)
	}
	return nil
}

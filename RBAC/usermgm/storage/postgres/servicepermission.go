package postgres

import (
	"context"
	"fmt"

	"brank.as/rbac/usermgm/storage"
)

const upsertServicePerm = `
INSERT INTO service_permissions(
	service_id,
	name,
	description,
	resource,
	action
) VALUES (
	:service_id,
	:name,
	:description,
	:resource,
	:action
)
ON CONFLICT (service_id, resource, action)
DO UPDATE SET
(name, description) = (:name,:description)
RETURNING *;
`

func (s *Storage) UpsertServicePermission(ctx context.Context, ps storage.ServicePermission) (*storage.ServicePermission, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, upsertServicePerm)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&ps, ps); err != nil {
		return nil, err
	}

	return &ps, nil
}

const getSvcPerm = `
	SELECT *
	FROM service_permissions
	WHERE id = $1
`

func (s *Storage) GetServicePermission(ctx context.Context, id string) (*storage.ServicePermission, error) {
	var ps storage.ServicePermission
	if err := s.db.Get(&ps, getSvcPerm, id); err != nil {
		return nil, fmt.Errorf("executing service permission details: %w", err)
	}
	return &ps, nil
}

func (s *Storage) AllServicePermission(ctx context.Context) ([]storage.ServicePermission, error) {
	const getSvcPerm = `SELECT * FROM service_permissions`
	ps := []storage.ServicePermission{}
	if err := s.db.Select(&ps, getSvcPerm); err != nil {
		return nil, err
	}
	return ps, nil
}

func (s *Storage) ListServicePermission(ctx context.Context, service string) ([]storage.ServicePermission, error) {
	const getSvcPerm = `SELECT sp.*
FROM service_permissions as sp inner join service on sp.service_id=service.service_id
where sp.service_id = $1
`
	ps := []storage.ServicePermission{}
	if err := s.db.Select(&ps, getSvcPerm, service); err != nil {
		return nil, err
	}
	return ps, nil
}

const deleteSvcPerm = `
DELETE from service_permissions
WHERE id = $1
`

func (s *Storage) DeleteServicePermission(ctx context.Context, id string) error {
	if res, err := s.db.Exec(deleteSvcPerm, id); err != nil {
		return fmt.Errorf("executing service permission delete: %w", err)
	} else if _, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("executing service permission delete: %w", err)
	}
	return nil
}

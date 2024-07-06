package postgres

import (
	"context"
	"fmt"

	"brank.as/rbac/usermgm/storage"
)

const upsertService = `
INSERT INTO service(
	service_name,
	description,
	assign_default
) VALUES (
	:service_name,
	:description,
	:assign_default
)
ON CONFLICT (service_name)
DO UPDATE 
	SET description = :description
RETURNING *;
`

func (s *Storage) UpsertService(ctx context.Context, ps storage.Service) (*storage.Service, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, upsertService)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&ps, ps); err != nil {
		return nil, err
	}

	return &ps, nil
}

const updateService = `
UPDATE service SET (
	service_name,
	description,
	assign_default
) = (
	:service_name,
	:description,
	:assign_default
)
WHERE service_id = :service_id
RETURNING *
`

func (s *Storage) UpdateService(ctx context.Context, ps storage.Service) (*storage.Service, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateService)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&ps, ps); err != nil {
		return nil, err
	}

	return &ps, nil
}

func (s *Storage) GetService(ctx context.Context, id string) (*storage.Service, error) {
	const getSvc = `SELECT * FROM service WHERE service_id = $1`
	var ps storage.Service
	if err := s.db.Get(&ps, getSvc, id); err != nil {
		return nil, fmt.Errorf("executing service permission details: %w", err)
	}
	return &ps, nil
}

func (s *Storage) ListService(ctx context.Context) ([]storage.Service, error) {
	const getSvc = ` SELECT * FROM service`
	ps := []storage.Service{}
	if err := s.db.Select(&ps, getSvc); err != nil {
		return nil, err
	}
	return ps, nil
}

func (s *Storage) DeleteService(ctx context.Context, id string) error {
	const deleteSvc = `DELETE from service WHERE service_id = $1`
	if res, err := s.db.Exec(deleteSvc, id); err != nil {
		return fmt.Errorf("executing service permission delete: %w", err)
	} else if _, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("executing service permission delete: %w", err)
	}
	return nil
}

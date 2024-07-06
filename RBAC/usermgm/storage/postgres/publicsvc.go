package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"brank.as/rbac/usermgm/storage"
)

const upsertServicePublic = `
INSERT INTO public_service (
	service_id,
	published_user_id,
	environment
) VALUES (
	:service_id,
	:published_user_id,
	:environment
)
RETURNING grant_id;
`

func (s *Storage) UpsertPublicService(ctx context.Context, ps storage.DefaultService) (*storage.DefaultService, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, upsertServicePublic)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&ps, ps); err != nil {
		return nil, err
	}

	return &ps, nil
}

func (s *Storage) GetPublicServiceEnv(ctx context.Context, env string) ([]storage.DefaultService, error) {
	const getSvcPublicSvc = `SELECT * FROM public_service WHERE environment = $1`
	ds := []storage.DefaultService{}
	if err := s.db.Select(&ds, getSvcPublicSvc, env); err != nil {
		return nil, fmt.Errorf("executing servicePublic permission details: %w", err)
	}
	return ds, nil
}

func (s *Storage) GetPublicService(ctx context.Context, serviceID string) ([]storage.DefaultService, error) {
	const getEnvPublic = `SELECT * FROM public_service WHERE service_id = $1`
	ps := []storage.DefaultService{}
	if err := s.db.Select(&ps, getEnvPublic, serviceID); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, fmt.Errorf("executing servicePublic permission details: %w", err)
	}
	return ps, nil
}

func (s *Storage) ListServicePublic(ctx context.Context) ([]storage.DefaultService, error) {
	const getSvc = `SELECT * FROM public_service`
	ps := []storage.DefaultService{}
	if err := s.db.Select(&ps, getSvc); err != nil {
		return nil, err
	}
	return ps, nil
}

const deleteSvcPublic = `
WITH del_id AS (
    UPDATE
        public_service
    SET
        (retracted_user_id,
            retracted) = (:retracted_user_id,
            NOW())
    WHERE
        grant_id = :grant_id 
    RETURNING
        grant_id)
DELETE FROM public_service
WHERE grant_id = (
        SELECT
            grant_id
        FROM
            del_id
        LIMIT 1)
`

func (s *Storage) RetractService(ctx context.Context, asn storage.DefaultService) error {
	stmt, err := s.db.PrepareNamedContext(ctx, deleteSvcPublic)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, asn)
	return err
}

func (s *Storage) AuditPublicServices(ctx context.Context, orgID string) ([]storage.DefaultService, error) {
	const auditPS = `
SELECT
    *
FROM
    public_service,
    public_service_audit
WHERE
    org_id = $1;
`
	sa := []storage.DefaultService{}
	if err := s.db.Select(&sa, auditPS, orgID); err != nil {
		return nil, err
	}
	return sa, nil
}

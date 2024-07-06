package postgres

import (
	"context"
	"fmt"

	"brank.as/rbac/usermgm/storage"
)

const upsertServiceAssign = `
INSERT INTO service_assignment (
	service_id,
	org_id,
	assign_user_id,
	assign_default,
	environment
) VALUES (
	:service_id,
	:org_id,
	:assign_user_id,
	:assign_default,
	:environment
)
RETURNING grant_id;
`

func (s *Storage) AssignService(ctx context.Context, ps storage.ServiceAssignment) (string, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, upsertServiceAssign)
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

const getSvcAssign = `
	SELECT *
	FROM service_assignment 
	WHERE org_id = $1
`

func (s *Storage) GetAssignedServices(ctx context.Context, orgID string) ([]storage.ServiceAssignment, error) {
	ps := []storage.ServiceAssignment{}
	if err := s.db.Select(&ps, getSvcAssign, orgID); err != nil {
		return nil, fmt.Errorf("executing serviceAssign permission details: %w", err)
	}
	return ps, nil
}

const getSvcAssignSvc = `
	SELECT *
	FROM service_assignment 
	WHERE org_id = :org_id AND service_id = :service_id AND environment = :environment
`

func (s *Storage) GetAssignedService(ctx context.Context, gr storage.ServiceAssignment) (*storage.ServiceAssignment, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, getSvcAssignSvc)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&gr, gr); err != nil {
		return nil, fmt.Errorf("executing serviceAssign permission details: %w", err)
	}
	return &gr, nil
}

const getOrgAssign = `
	SELECT *
	FROM service_assignment 
	WHERE service_id = $1
`

func (s *Storage) GetServiceOrgs(ctx context.Context, serviceID string) ([]storage.ServiceAssignment, error) {
	ps := []storage.ServiceAssignment{}
	if err := s.db.Select(&ps, getOrgAssign, serviceID); err != nil {
		return nil, fmt.Errorf("executing serviceAssign permission details: %w", err)
	}
	return ps, nil
}

func (s *Storage) ListServiceAssign(ctx context.Context) ([]storage.ServiceAssignment, error) {
	const getSvc = ` SELECT * FROM service_assignment `
	ps := []storage.ServiceAssignment{}
	if err := s.db.Select(&ps, getSvc); err != nil {
		return nil, err
	}
	return ps, nil
}

func (s *Storage) ListServiceAssignID(ctx context.Context, service string) ([]storage.ServiceAssignment, error) {
	const getSvc = ` SELECT * FROM service_assignment where service_id = $1`
	ps := []storage.ServiceAssignment{}
	if err := s.db.Select(&ps, getSvc, service); err != nil {
		return nil, err
	}
	return ps, nil
}

func (s *Storage) ListServiceAssignOrg(ctx context.Context, org string) ([]storage.ServiceAssignment, error) {
	const getSvc = ` SELECT * FROM service_assignment where org_id = $1`
	ps := []storage.ServiceAssignment{}
	if err := s.db.Select(&ps, getSvc, org); err != nil {
		return nil, err
	}
	return ps, nil
}

const deleteSvcAssign = `
WITH del_id AS (
    UPDATE
        service_assignment
    SET
        (revoke_user_id,
            revoked) = (:revoke_user_id,
            NOW())
    WHERE
        grant_id = :grant_id 
    RETURNING
        grant_id)
DELETE FROM service_assignment
WHERE grant_id = (
        SELECT
            grant_id
        FROM
            del_id
        LIMIT 1)
`

func (s *Storage) RevokeService(ctx context.Context, asn storage.ServiceAssignment) error {
	stmt, err := s.db.PrepareNamedContext(ctx, deleteSvcAssign)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, asn)
	return err
}

func (s *Storage) AuditOrgServices(ctx context.Context, orgID string) ([]storage.ServiceAssignment, error) {
	const auditSA = `
SELECT
    *
FROM
    service_assignment,
    service_assignment_audit
WHERE
    org_id = $1;
`
	sa := []storage.ServiceAssignment{}
	if err := s.db.Select(&sa, auditSA, orgID); err != nil {
		return nil, err
	}
	return sa, nil
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"brank.as/petnet/profile/storage"
	"github.com/lib/pq"
)

// CreatePartner ...
func (s *Storage) CreatePartner(ctx context.Context, pnr *storage.Partner) (string, error) {
	if pnr.OrgID == "" {
		return "", fmt.Errorf("org_id cannot be empty")
	}
	if pnr.Type == "" {
		return "", fmt.Errorf("type cannot be empty")
	}

	const insertPartners = `
INSERT INTO partner (
    org_id,
	 type,
	 service,
	 updated_by,
	 status
) VALUES (
	 :org_id,
	 :type,
	 :service,
	 :updated_by,
	 :status
)
RETURNING id
`
	pstmt, err := s.db.PrepareNamedContext(ctx, insertPartners)
	if err != nil {
		return "", err
	}
	defer pstmt.Close()
	if err := pstmt.Get(pnr, pnr); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return "", storage.Conflict
		}
		return "", fmt.Errorf("executing partner insert: %w", err)
	}
	return pnr.ID, nil
}

// UpdatePartner updates the db values for a given partner using the ID
func (s *Storage) UpdatePartner(ctx context.Context, pnr *storage.Partner) (string, error) {
	const partnerUpdate = `
UPDATE
	partner
SET
	service = :service,
	updated_by = :updated_by,
	deleted = NULL
WHERE
	id = :id
RETURNING id
`
	stmt, err := s.db.PrepareNamedContext(ctx, partnerUpdate)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	if err := stmt.Get(pnr, pnr); err != nil {
		if err == sql.ErrNoRows {
			return "", storage.NotFound
		}
		return "", fmt.Errorf("executing partner update: %w", err)
	}
	return pnr.ID, nil
}

// GetPartners return partners matched against orgID
func (s *Storage) GetPartners(ctx context.Context, oid string) ([]storage.Partner, error) {
	const partnerSelect = `
SELECT *
FROM partner
WHERE
	org_id = :org_id AND deleted IS NULL
`
	var pnr []storage.Partner
	stmt, err := s.db.PrepareNamed(partnerSelect)
	if err != nil {
		return nil, fmt.Errorf("preparing named query GetPartners: %w", err)
	}
	arg := map[string]interface{}{
		"org_id": oid,
	}
	if err := stmt.Select(&pnr, arg); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return pnr, nil
}

// DeletePartner tombstones given partner by ID
func (s *Storage) DeletePartner(ctx context.Context, id string) (string, error) {
	const partnerUpdate = `
UPDATE
	partner
SET
	deleted = :deleted
WHERE
	id = :id
`
	stmt, err := s.db.PrepareNamedContext(ctx, partnerUpdate)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	arg := map[string]interface{}{
		"id":      id,
		"deleted": time.Now(),
	}
	if _, err := stmt.Exec(arg); err != nil {
		if err == sql.ErrNoRows {
			return "", storage.NotFound
		}
		return "", fmt.Errorf("executing partner update: %w", err)
	}
	return id, nil
}

// HaspnrAccess checks if a partner is enabled or not for org
func (s *Storage) ValidatePartnerAccess(ctx context.Context, oid, pnr string) error {
	const q = `
SELECT *
FROM partner
WHERE
	org_id = $1 AND type = $2 AND deleted IS NULL AND status = 'ENABLED'
`
	d := &storage.Partner{}
	err := s.db.Get(d, q, oid, pnr)
	if err != nil {
		if err == sql.ErrNoRows {
			return storage.NotFound
		}
		return err
	}
	return nil
}

// EnablePartner enables a partner for an org
func (s *Storage) EnablePartner(ctx context.Context, oid, pnr string) error {
	const q = `
UPDATE
	partner
SET
	status = 'ENABLED'
WHERE
	org_id = $1 AND type = $2
RETURNING id
`
	_, err := s.db.Exec(q, oid, pnr)
	if err != nil {
		if err == sql.ErrNoRows {
			return storage.NotFound
		}
		return err
	}
	return nil
}

// DisablePartner disables a partner for an org
func (s *Storage) DisablePartner(ctx context.Context, oid, pnr string) error {
	const q = `
UPDATE
	partner
SET
	status = 'DISABLED'
WHERE
	org_id = $1 AND type = $2
RETURNING id
`
	_, err := s.db.Exec(q, oid, pnr)
	if err != nil {
		if err == sql.ErrNoRows {
			return storage.NotFound
		}
		return err
	}
	return nil
}

// GetPartner return partner matched against orgID and Type
func (s *Storage) GetPartner(ctx context.Context, oid string, tp string) (*storage.Partner, error) {
	const partnerSelect = `
SELECT *
FROM partner
WHERE
	org_id = :org_id AND type = :type
`
	var pnr storage.Partner
	stmt, err := s.db.PrepareNamed(partnerSelect)
	if err != nil {
		return nil, fmt.Errorf("preparing named query GetPartner: %w", err)
	}
	arg := map[string]interface{}{
		"org_id": oid,
		"type":   tp,
	}
	if err := stmt.Get(&pnr, arg); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return &pnr, nil
}

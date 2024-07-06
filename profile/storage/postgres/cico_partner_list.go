package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"brank.as/petnet/profile/storage"
	"github.com/lib/pq"
)

// CreateCICOPartnerList ...
func (s *Storage) CreateCICOPartnerList(ctx context.Context, pnr *storage.CICOPartnerList) (*storage.CICOPartnerList, error) {
	if pnr.Stype == "" {
		return nil, fmt.Errorf("partner type cannot be empty")
	}
	if pnr.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if pnr.Status == "" {
		return nil, fmt.Errorf("status cannot be empty")
	}

	const insertCICOPartnerList = `
INSERT INTO cico_partner_list (
	stype,
	name,
	status
) VALUES (
	:stype,
	:name,
	:status
)
RETURNING *
`
	pstmt, err := s.db.PrepareNamedContext(ctx, insertCICOPartnerList)
	if err != nil {
		return nil, err
	}
	defer pstmt.Close()
	if err := pstmt.Get(pnr, pnr); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, fmt.Errorf("executing cico partner insert: %w", err)
	}
	return pnr, nil
}

// UpdateCICOPartnerList updates the db values for a given partner using the ID
func (s *Storage) UpdateCICOPartnerList(ctx context.Context, pnr *storage.CICOPartnerList) (*storage.CICOPartnerList, error) {
	const cicoPartnerUpdate = `
UPDATE
	cico_partner_list
SET
	name = :name,
	status = :status,
	deleted = NULL
WHERE
	stype = :stype
RETURNING *
`
	stmt, err := s.db.PrepareNamedContext(ctx, cicoPartnerUpdate)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(pnr, pnr); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, fmt.Errorf("executing cico partner update: %w", err)
	}
	return pnr, nil
}

// DeleteCICOPartnerList tombstones given partner by ID
func (s *Storage) DeleteCICOPartnerList(ctx context.Context, stype string) (string, error) {
	const cicoPartnerUpdate = `
UPDATE
	cico_partner_list
SET
	deleted = :deleted
WHERE
	stype = :stype
`
	stmt, err := s.db.PrepareNamedContext(ctx, cicoPartnerUpdate)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	arg := map[string]interface{}{
		"stype":   stype,
		"deleted": time.Now(),
	}
	if _, err := stmt.Exec(arg); err != nil {
		if err == sql.ErrNoRows {
			return "", storage.NotFound
		}
		return "", fmt.Errorf("executing cico partner update: %w", err)
	}
	return stype, nil
}

// EnablePartnerList enables a partner for an org
func (s *Storage) EnableCICOPartnerList(ctx context.Context, stype string) (string, error) {
	const q = `
UPDATE
	cico_partner_list
SET
	status = 'ENABLED'
WHERE
	 stype = $1
RETURNING stype
`
	_, err := s.db.Exec(q, stype)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", storage.NotFound
		}
		return "", err
	}
	return stype, nil
}

// DisableCICOPartnerList disables a partner for an org
func (s *Storage) DisableCICOPartnerList(ctx context.Context, stype string) (string, error) {
	const q = `
UPDATE
	cico_partner_list
SET
	status = 'DISABLED'
WHERE
	stype = $1
RETURNING stype
`
	_, err := s.db.Exec(q, stype)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", storage.NotFound
		}
		return "", err
	}
	return stype, nil
}

// GetCICOPartnerList return partner matched against orgID and Type
func (s *Storage) GetCICOPartnerList(ctx context.Context, req *storage.CICOPartnerList) ([]*storage.CICOPartnerList, error) {
	partnerSelect := `SELECT * FROM cico_partner_list WHERE deleted IS NULL`
	if req.ID != "" {
		partnerSelect += " AND id = :id"
	}
	if req.Stype != "" {
		stsM := strings.Split(req.Stype, ",")
		stsTemp := []string{}
		reqStype := ""
		if len(stsM) > 0 {
			for _, v := range stsM {
				stsTemp = append(stsTemp, "'"+v+"'")
			}
			reqStype = strings.Join(stsTemp, ",")
		}
		if reqStype != "" {
			partnerSelect += fmt.Sprintf(" AND stype IN(%s) ", reqStype)
		}
	}
	if req.Status != "" {
		partnerSelect += " AND status = :status"
	}
	if req.Name != "" {
		partnerSelect += " AND name ILIKE '%" + req.Name + "%'"
	}
	partnerSelect += " ORDER BY id ASC"
	var pnr []*storage.CICOPartnerList
	stmt, err := s.db.PrepareNamed(partnerSelect)
	if err != nil {
		return nil, fmt.Errorf("preparing named query Get CICO Partner List: %w", err)
	}
	if err := stmt.Select(&pnr, req); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return pnr, nil
}

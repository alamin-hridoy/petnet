package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/svcutil/mw"
	"github.com/lib/pq"
)

// CreatePartnerList ...
func (s *Storage) CreatePartnerList(ctx context.Context, pnr *storage.PartnerList) (*storage.PartnerList, error) {
	if pnr.Stype == "" {
		return nil, fmt.Errorf("partner type cannot be empty")
	}
	if pnr.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if pnr.Status == "" {
		return nil, fmt.Errorf("status cannot be empty")
	}
	if pnr.ServiceName == "" {
		return nil, fmt.Errorf("service name cannot be empty")
	}
	if pnr.Platform == "" {
		pnr.Platform = "Perahub"
	}

	const insertPartnerList = `
INSERT INTO partner_list (
	stype,
	name,
	status,
	service_name,
	platform,
	is_provider,
	updated_by,
	perahub_partner_id,
	remco_id
) VALUES (
	:stype,
	:name,
	:status,
	:service_name,
	:platform,
	:is_provider,
	:updated_by,
	:perahub_partner_id,
	:remco_id
)
RETURNING *
`
	pstmt, err := s.db.PrepareNamedContext(ctx, insertPartnerList)
	if err != nil {
		return nil, err
	}
	defer pstmt.Close()
	if err := pstmt.Get(pnr, pnr); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, fmt.Errorf("executing partner insert: %w", err)
	}
	return pnr, nil
}

// UpdatePartnerList updates the db values for a given partner using the ID
func (s *Storage) UpdatePartnerList(ctx context.Context, pnr *storage.PartnerList) (*storage.PartnerList, error) {
	const partnerUpdate = `
UPDATE
	partner_list
SET
	name = COALESCE(NULLIF(:name, ''), name),
	status = COALESCE(NULLIF(:status, ''), status),
	deleted = NULL,
	service_name = COALESCE(NULLIF(:service_name, ''), service_name),
	updated_by = :updated_by,
	disable_reason = :disable_reason,
	platform = :platform,
	is_provider = :is_provider,
	perahub_partner_id = COALESCE(NULLIF(:perahub_partner_id, ''), perahub_partner_id),
	remco_id = COALESCE(NULLIF(:remco_id, ''), remco_id)
WHERE
	stype = :stype
RETURNING *
`
	stmt, err := s.db.PrepareNamedContext(ctx, partnerUpdate)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(pnr, pnr); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, fmt.Errorf("executing partner update: %w", err)
	}
	return pnr, nil
}

// DeletePartnerList tombstones given partner by ID
func (s *Storage) DeletePartnerList(ctx context.Context, stype string) (string, error) {
	const partnerUpdate = `
UPDATE
	partner_list
SET
	deleted = :deleted
WHERE
	stype = :stype
`
	stmt, err := s.db.PrepareNamedContext(ctx, partnerUpdate)
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
		return "", fmt.Errorf("executing partner update: %w", err)
	}
	return stype, nil
}

// EnablePartnerList enables a partner for an org
func (s *Storage) EnablePartnerList(ctx context.Context, stype string) (string, error) {
	uid := mw.GetUserID(ctx)
	const q = `
UPDATE
	partner_list
SET
	status = 'ENABLED',
	updated_by = $2
WHERE
	 stype = $1
RETURNING stype
`
	_, err := s.db.Exec(q, stype, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", storage.NotFound
		}
		return "", err
	}
	return stype, nil
}

// DisablePartnerList disables a partner for an org
func (s *Storage) DisablePartnerList(ctx context.Context, stype string, disableReason string) (string, error) {
	uid := mw.GetUserID(ctx)
	const q = `
UPDATE
	partner_list
SET
	status = 'DISABLED',
	updated_by = $2,
	disable_reason = $3
WHERE
	stype = $1
RETURNING stype
`
	_, err := s.db.Exec(q, stype, uid, disableReason)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", storage.NotFound
		}
		return "", err
	}
	return stype, nil
}

// GetPartnerList return partner matched against orgID and Type
func (s *Storage) GetPartnerList(ctx context.Context, req *storage.PartnerList) ([]*storage.PartnerList, error) {
	partnerSelect := `SELECT * FROM partner_list WHERE deleted IS NULL`
	if req.ID != "" {
		partnerSelect += " AND id = :id"
	}
	if req.Stype != "" {
		stsM := strings.Split(req.Stype, ",")
		stsTemp := []string{}
		reqStype := ""
		if len(stsM) > 0 {
			for _, v := range stsM {
				stsTemp = append(stsTemp, "'"+strings.TrimSpace(v)+"'")
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

	if req.ServiceName != "" && req.ServiceName != "EMPTYSERVICETYPE" {
		partnerSelect += " AND service_name = :service_name"
	}
	if req.Name != "" {
		partnerSelect += " AND name ILIKE '%" + req.Name + "%'"
	}
	if req.IsProvider {
		partnerSelect += " AND is_provider = :is_provider"
	}
	partnerSelect += " ORDER BY name ASC"
	var pnr []*storage.PartnerList
	stmt, err := s.db.PrepareNamed(partnerSelect)
	if err != nil {
		return nil, fmt.Errorf("preparing named query Get Partner List: %w", err)
	}
	if err := stmt.Select(&pnr, req); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return pnr, nil
}

// GetPartnerList return partner matched against org Transaction Type
func (s *Storage) GetDSAPartnerList(ctx context.Context, req *storage.GetDSAPartnerListRequest) ([]storage.DSAPartnerList, error) {
	transactionTypesS := []string{}
	transactionTypes := ""
	if len(req.TransactionTypes) > 0 {
		for _, v := range req.TransactionTypes {
			transactionTypesS = append(transactionTypesS, "'"+v+"'")
		}
	}
	transactionTypes = strings.Join(transactionTypesS, ",")
	partnerSelect := `SELECT partner, transaction_type FROM partner_commission_config  `
	if transactionTypes != "" {
		partnerSelect += fmt.Sprintf(` WHERE transaction_type IN(%s)`, transactionTypes)
	}
	partnerSelect += " Group by partner, transaction_type"
	var pnr []storage.DSAPartnerList
	stmt, err := s.db.PrepareNamed(partnerSelect)
	if err != nil {
		return nil, fmt.Errorf("preparing named query Get DSA Partner List: %w", err)
	}
	if err := stmt.Select(&pnr, req); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return pnr, nil
}

// GetPartnerTransactionType return partner and transaction type
func (s *Storage) GetPartnerTransactionType(ctx context.Context, req *storage.GetAllPartnerListReq) ([]storage.GetAllPartnerListReq, error) {
	PartnerS := []string{}
	Partners := ""
	tTSlce := strings.Split(req.Partner, ",")
	if len(tTSlce) > 0 {
		for _, v := range tTSlce {
			PartnerS = append(PartnerS, "'"+strings.TrimSpace(v)+"'")
		}
		Partners = strings.Join(PartnerS, ",")
	}
	partnerSelect := `SELECT partner, transaction_type FROM partner_commission_config  `
	if Partners != "" {
		partnerSelect += fmt.Sprintf(` WHERE partner IN(%s)`, Partners)
	}
	partnerSelect += " Group by partner, transaction_type"
	var ptnr []storage.GetAllPartnerListReq
	stmt, err := s.db.PrepareNamed(partnerSelect)
	if err != nil {
		return nil, fmt.Errorf("preparing named query Get DSA Partner List: %w", err)
	}
	if err := stmt.Select(&ptnr, req); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return ptnr, nil
}

func (s *Storage) EnableMultiplePartnerList(ctx context.Context, stypes []string, updatedBy string) ([]string, error) {
	updated := time.Now()
	stypeStr := sliceFormator(stypes)

	const q = `UPDATE partner_list SET
				status = 'ENABLED',
				updated = $1,
				updated_by = $2
				WHERE`

	query := fmt.Sprintf("%s stype IN(%s)", q, stypeStr)
	_, err := s.db.Exec(query, updated, updatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return stypes, nil
}

func (s *Storage) DisableMultiplePartnerList(ctx context.Context, stypes []string, disableReason string, updatedBy string) ([]string, error) {
	updated := time.Now()
	stypeStr := sliceFormator(stypes)

	const q = `UPDATE partner_list SET
				status = 'DISABLED',
				updated = $1,
				updated_by = $2,
				disable_reason = $3
				WHERE`

	query := fmt.Sprintf("%s stype IN(%s)", q, stypeStr)
	_, err := s.db.Exec(query, updated, updatedBy, disableReason)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return stypes, nil
}

func (s *Storage) GetPartnerByStype(ctx context.Context, stype string) (*storage.PartnerList, error) {
	if stype == "" {
		return nil, fmt.Errorf("stype can't be empty")
	}
	const partnerListStype = `SELECT * FROM partner_list WHERE stype = $1 OR stype = $2`
	var pf storage.PartnerList
	if err := s.db.Get(&pf, partnerListStype, strings.ToLower(stype), strings.ToTitle(stype)); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return &pf, nil
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"github.com/lib/pq"
)

const insertOrgProfile = `
INSERT INTO org_profile (
    org_id,
    user_id,
	org_type,
    status,
    risk_score,
    date_applied,
    bus_info_company_name,
    bus_info_store_name,
    bus_info_phone_number,
    bus_info_fax_number,
    bus_info_website,
    bus_info_company_email,
    bus_info_contact_person,
    bus_info_position,
    bus_info_address1,
    bus_info_city,
    bus_info_state,
    bus_info_postal_code,
    acc_info_bank,
    acc_info_bank_account_number,
    acc_info_bank_account_holder,
    acc_info_agree_terms_conditions,
    acc_info_agree_online_supplier_form,
    acc_info_currency,
    transaction_types,
	reminder_sent,
	dsa_code,
	terminal_id_otc,
	terminal_id_digital,
	is_provider,
	partner
) VALUES (
	 :org_id,
	 :user_id,
	 :org_type,
	 :status,
	 :risk_score,
	 :date_applied,
	 :bus_info_company_name,
	 :bus_info_store_name,
	 :bus_info_phone_number,
	 :bus_info_fax_number,
	 :bus_info_website,
	 :bus_info_company_email,
	 :bus_info_contact_person,
	 :bus_info_position,
	 :bus_info_address1,
	 :bus_info_city,
	 :bus_info_state,
	 :bus_info_postal_code,
	 :acc_info_bank,
	 :acc_info_bank_account_number,
	 :acc_info_bank_account_holder,
	 :acc_info_agree_terms_conditions,
	 :acc_info_agree_online_supplier_form,
	 :acc_info_currency,
	 :transaction_types,
	 :reminder_sent,
	 :dsa_code,
	 :terminal_id_otc,
	 :terminal_id_digital,
	 :is_provider,
	 :partner
) RETURNING
    id,created,updated
`

// CreateOrgProfile creates new org profile and returns the created profile's ID.
func (s *Storage) CreateOrgProfile(ctx context.Context, pf *storage.OrgProfile) (string, error) {
	log := logging.FromContext(ctx)
	pstmt, err := s.db.PrepareNamedContext(ctx, insertOrgProfile)
	if err != nil {
		logging.WithError(err, log).Error("insert org")
		return "", err
	}
	defer pstmt.Close()
	if err := pstmt.Get(pf, pf); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return "", storage.Conflict
		}
		return "", fmt.Errorf("executing org profile insert: %w", err)
	}
	return pf.ID, nil
}

const orgProfileUpdateByOrgID = `
UPDATE
	org_profile
SET
	org_type= COALESCE(NULLIF(:org_type, 0), org_type),
	status= COALESCE(NULLIF(:status, 0), status),
	risk_score= COALESCE(NULLIF(:risk_score, 0), risk_score),
	date_applied= COALESCE(:date_applied, date_applied),
	bus_info_company_name= COALESCE(NULLIF(:bus_info_company_name, ''), bus_info_company_name),
	bus_info_store_name= COALESCE(NULLIF(:bus_info_store_name, ''), bus_info_store_name),
	bus_info_phone_number= COALESCE(NULLIF(:bus_info_phone_number, ''), bus_info_phone_number),
	bus_info_fax_number= COALESCE(NULLIF(:bus_info_fax_number, ''), bus_info_fax_number),
	bus_info_website= COALESCE(NULLIF(:bus_info_website, ''), bus_info_website),
	bus_info_company_email= COALESCE(NULLIF(:bus_info_company_email, ''), bus_info_company_email),
	bus_info_contact_person= COALESCE(NULLIF(:bus_info_contact_person, ''), bus_info_contact_person),
	bus_info_position= COALESCE(NULLIF(:bus_info_position, ''), bus_info_position),
	bus_info_address1= COALESCE(NULLIF(:bus_info_address1, ''), bus_info_address1),
	bus_info_city= COALESCE(NULLIF(:bus_info_city, ''), bus_info_city),
	bus_info_state= COALESCE(NULLIF(:bus_info_state, ''), bus_info_state),
	bus_info_postal_code= COALESCE(NULLIF(:bus_info_postal_code, ''), bus_info_postal_code),
	acc_info_bank= COALESCE(NULLIF(:acc_info_bank, ''), acc_info_bank),
	acc_info_bank_account_number= COALESCE(NULLIF(:acc_info_bank_account_number, ''), acc_info_bank_account_number),
	acc_info_bank_account_holder= COALESCE(NULLIF(:acc_info_bank_account_holder, ''), acc_info_bank_account_holder),
	acc_info_agree_terms_conditions= COALESCE(NULLIF(:acc_info_agree_terms_conditions, 0), acc_info_agree_terms_conditions),
	acc_info_agree_online_supplier_form= COALESCE(NULLIF(:acc_info_agree_online_supplier_form, 0), acc_info_agree_online_supplier_form),
	acc_info_currency= COALESCE(NULLIF(:acc_info_currency, 0), acc_info_currency),
	reminder_sent= COALESCE(NULLIF(:reminder_sent, 0), reminder_sent),
	transaction_types= COALESCE(NULLIF(:transaction_types, ''), transaction_types),
	deleted= COALESCE(:deleted, deleted),
	dsa_code= COALESCE(NULLIF(:dsa_code, ''), dsa_code),
	terminal_id_otc= COALESCE(NULLIF(:terminal_id_otc, ''), terminal_id_otc),
	terminal_id_digital= COALESCE(NULLIF(:terminal_id_digital, ''), terminal_id_digital),
	is_provider= COALESCE(NULLIF(:is_provider, FALSE), is_provider),
	partner= COALESCE(NULLIF(:partner, ''), partner)
WHERE
	org_id = :org_id
RETURNING
    id,created,updated
`

// UpdateOrgProfile updates the db values for a given org profile using the org ID
func (s *Storage) UpdateOrgProfile(ctx context.Context, pf *storage.OrgProfile) (string, error) {
	if pf.OrgID == "" {
		return "", fmt.Errorf("invalid org id")
	}

	stmt, err := s.db.PrepareNamedContext(ctx, orgProfileUpdateByOrgID)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	if err := stmt.Get(pf, pf); err != nil {
		if err == sql.ErrNoRows {
			return "", storage.NotFound
		}
		return "", fmt.Errorf("executing org profile update: %w", err)
	}

	return pf.ID, nil
}

// GetProfiles return all profiles
func (s *Storage) GetOrgProfiles(ctx context.Context, f storage.FilterList) ([]storage.OrgProfile, error) {
	if f.SortBy == "" {
		f.SortBy = "ASC"
	}
	if f.SortByColumn == "" {
		f.SortByColumn = "date_applied"
	}

	submitted_document := []string{"'IDPhoto'", "'Picture'", "'NBIClearance'", "'CourtClearance'", "'IncorporationPapers'", "'MayorsPermit'", "'FinancialStatement'", "'BankStatement'"}

	orgcondition := ""
	conditions := `((partner ILIKE '%%' || $3 || '%%') OR (bus_info_company_name ILIKE '%%' || $3 || '%%') OR (bus_info_company_email ILIKE '%%' || $3 || '%%')  OR (u.email ILIKE '%%' || $3 || '%%')) AND COALESCE(risk_score = ANY($4), TRUE)
    AND COALESCE(status = ANY($5), TRUE)
    `
	if f.OrgType != "" {
		orgcondition = fmt.Sprintf(` AND (org_type = '%s') `, f.OrgType)
	}

	submittedQ := ""
	submittedJoin := ""
	submittedWhere := ""
	submitted_document_cnt := len(submitted_document)
	if f.SubmittedDocument != "" {
		submittedQ = fmt.Sprintf(` doc_count as (select sum(f.submitted) as doc_cnt, org_id
		from file_upload f where upload_type in (%s)
		group by org_id), `, strings.Join(submitted_document, ","))

		submittedJoin = `left join doc_count on p.org_id = doc_count.org_id `

		submittedWhere = fmt.Sprintf(" and doc_count.doc_cnt = %d ", submitted_document_cnt)
		if f.SubmittedDocument == "not-submitted" {
			submittedWhere = fmt.Sprintf(" and doc_count.doc_cnt < %d ", submitted_document_cnt)
		}
	}

	isProviderQ := fmt.Sprintf(` AND (is_provider = %t) `, f.IsProvider)

	getOrgProfiles := fmt.Sprintf(`
		WITH %s cnt AS (select count(*) as count
		FROM org_profile as p inner join user_profile u on (u.user_id::text = p.user_id) %s WHERE %s %s %s %s)
		SELECT p.*, cnt.count FROM org_profile as p 
		left join cnt on true inner join user_profile u on (u.user_id::text = p.user_id)
		%s
		WHERE %s %s
		%s %s
		ORDER BY p.%s %s
		LIMIT NULLIF($1, 0)
		OFFSET $2;
    `, submittedQ, submittedJoin, conditions, orgcondition, submittedWhere, isProviderQ, submittedJoin, conditions, orgcondition, submittedWhere, isProviderQ, f.SortByColumn, f.SortBy)
	var pfs []storage.OrgProfile
	// sqlvet: ignore
	if err := s.db.Select(&pfs, getOrgProfiles, f.Limit, f.Offset, strings.TrimSpace(f.CompanyName), pq.Int32Array(f.RiskScore), pq.Int32Array(f.Status)); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return pfs, nil
}

// GetOrgProfile return profile matched against org ID
func (s *Storage) GetOrgProfile(ctx context.Context, id string) (*storage.OrgProfile, error) {
	const orgProfileID = `SELECT * FROM org_profile WHERE org_id = $1`
	var pf storage.OrgProfile
	if err := s.db.Get(&pf, orgProfileID, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return &pf, nil
}

// GetOrgProfile return profile matched against DSA Code
func (s *Storage) GetProfileByDsaCode(ctx context.Context, dsaCode string) (*storage.OrgProfile, error) {
	const orgProfileID = `SELECT * FROM org_profile WHERE dsa_code = $1`
	var pf storage.OrgProfile
	if err := s.db.Get(&pf, orgProfileID, dsaCode); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return &pf, nil
}

// GetProfileID returns the OrgProfileID associated with the org id.
func (s *Storage) GetProfileID(ctx context.Context, orgID string) (string, error) {
	const getID = `select id from org_profile where org_id = $1 limit 1`
	pid := ""
	if err := s.db.Get(&pid, getID, orgID); err != nil {
		if err == sql.ErrNoRows {
			return "", storage.NotFound
		}
		return "", err
	}
	return pid, nil
}

const orgProfileUpdateUserIDByOrgID = `
UPDATE
	org_profile
SET
	org_id = :newOrgID,
	user_id = :user_id
WHERE
	org_id = :oldOrgID
RETURNING
    id
`

// UpdateOrgProfile updates the db values for a given org profile using the org ID
func (s *Storage) UpdateOrgProfileUserID(ctx context.Context, req *storage.UpdateOrgProfileOrgIDUserID) (string, error) {
	if req.OldOrgID == "" {
		return "", fmt.Errorf("invalid org id")
	}

	stmt, err := s.db.PrepareNamedContext(ctx, orgProfileUpdateUserIDByOrgID)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var id string
	if err := stmt.Get(&id, req); err != nil {
		if err == sql.ErrNoRows {
			return "", storage.NotFound
		}
		return "", fmt.Errorf("executing org profile update: %w", err)
	}
	return id, nil
}

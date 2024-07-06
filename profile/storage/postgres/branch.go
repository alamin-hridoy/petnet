package postgres

import (
	"context"
	"database/sql"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"github.com/lib/pq"
)

const insertBranch = `
INSERT INTO branch (
    org_id,
    org_profile_id,
    title,
    address1,
    city,
    state,
    postal_code,
    phone_number,
    fax_number,
    contact_person
) VALUES (
	 :org_id,
	 :org_profile_id,
	 :title,
	 :address1,
	 :city,
	 :state,
	 :postal_code,
	 :phone_number,
	 :fax_number,
	 :contact_person
) RETURNING
    id,created,updated
`

// CreateBranch inserts/updates a branch record.
func (s *Storage) CreateBranch(ctx context.Context, br storage.Branch) (*storage.Branch, error) {
	log := logging.FromContext(ctx)

	stmt, err := s.prepareNamed(ctx, insertBranch)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&br, br); err != nil {
		logging.WithError(err, log).Error("insert branch")
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, err
	}
	return &br, nil
}

const updateBranch = `
    UPDATE
        branch
    SET
        title = COALESCE(NULLIF (:title, ''), title),
        address1 = COALESCE(NULLIF (:address1, ''), address1),
        city = COALESCE(NULLIF (:city, ''), city),
        state = COALESCE(NULLIF (:state, ''), state),
        postal_code = COALESCE(NULLIF (:postal_code, ''), postal_code),
        phone_number = COALESCE(NULLIF (:phone_number, ''), phone_number),
        fax_number = :fax_number,
        contact_person = :contact_person,
		updated = COALESCE(:updated, updated),
        deleted = COALESCE(:deleted, deleted)
    WHERE
        id = :id
    RETURNING id, created;
`

// UpsertBranch inserts/updates a branch record.
func (s *Storage) UpsertBranch(ctx context.Context, br storage.Branch) (*storage.Branch, error) {
	log := logging.FromContext(ctx)
	if br.ID == "" {
		return s.CreateBranch(ctx, br)
	}

	stmt, err := s.prepareNamed(ctx, updateBranch)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&br, br); err != nil {
		logging.WithError(err, log).Error("upsert branch")
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, err
	}
	return &br, nil
}

// ListBranches return branch records matched against org ID
func (s *Storage) ListBranches(ctx context.Context, oid string, f storage.LimitOffsetFilter, title string) ([]storage.Branch, error) {
	const getBranches = `
	WITH cnt AS (select count(*) as count FROM branch WHERE org_id = $1 AND (title ILIKE '%%' || $4 || '%%'))
	SELECT *, cnt.count
	FROM branch as p left join cnt on true
	WHERE (title ILIKE '%%' || $4 || '%%') AND org_id = $1
	AND deleted is null
	ORDER BY created ASC
	LIMIT NULLIF($2, 0)
	OFFSET $3;
	`
	var fs []storage.Branch
	if err := s.db.Select(&fs, getBranches, oid, f.Limit, f.Offset, title); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return fs, nil
}

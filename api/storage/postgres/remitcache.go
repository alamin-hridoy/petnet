package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"brank.as/petnet/api/storage"
	"brank.as/petnet/serviceutil/logging"
)

const createRemitCache = `
INSERT INTO remit_cache (
   dsa_id,
	user_id,
	remco_id,
	remco_member_id,
	remco_control_number,
	remco_alternate_control_number,
	remit_type,
	partner_remit_type,
	status,
	remit
) VALUES (
	:dsa_id,
	:user_id,
	:remco_id,
	:remco_member_id,
	:remco_control_number,
	:remco_alternate_control_number,
	:remit_type,
	:partner_remit_type,
	:status,
	:remit
) RETURNING
transaction_id, created,updated
`

func (s *Storage) CreateRemitCache(ctx context.Context, r storage.RemitCache) (*storage.RemitCache, error) {
	log := logging.FromContext(ctx)
	log.WithField("remit", r).Trace("storing")

	stmt, err := s.db.PrepareNamedContext(ctx, createRemitCache)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		return nil, fmt.Errorf("executing remit cache insert: %w", err)
	}
	return &r, nil
}

const updateRemitCache = `
	UPDATE remit_cache
	SET
		remco_control_number = :remco_control_number,
		remco_alternate_control_number = :remco_alternate_control_number,
		status = :status
	WHERE transaction_id = :transaction_id
	RETURNING updated`

func (s *Storage) UpdateRemitCache(ctx context.Context, r storage.RemitCache) (*storage.RemitCache, error) {
	if r.TxnID == "" {
		return nil, fmt.Errorf("transaction ID cannot be empty")
	}
	stmt, err := s.db.PrepareNamedContext(ctx, updateRemitCache)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		return nil, fmt.Errorf("executing remit cache update: %w", err)
	}
	return &r, nil
}

const getRemit = `
SELECT *
FROM remit_cache
WHERE transaction_id = $1
`

func (s *Storage) GetRemitCache(ctx context.Context, txnID string) (*storage.RemitCache, error) {
	var r storage.RemitCache
	if err := s.db.Get(&r, getRemit, txnID); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}
	return &r, nil
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"github.com/lib/pq"
)

const insertFee = `
INSERT INTO fee_commission (
    org_id,
    fee_commision_type,
    org_profile_id,
    fee_amount,
    commission_amount,
	 start_date,
	 end_date
) VALUES (
	 :org_id,
	 :fee_commision_type,
	 :org_profile_id,
	 :fee_amount,
	 :commission_amount,
	 :start_date,
	 :end_date
) RETURNING
    id,created,updated
`

// CreateOrgFees ...
func (s *Storage) CreateOrgFees(ctx context.Context, fee storage.FeeCommission) (*storage.FeeCommission, error) {
	log := logging.FromContext(ctx)
	stmt, err := s.prepareNamed(ctx, insertFee)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&fee, fee); err != nil {
		logging.WithError(err, log)
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, fmt.Errorf("executing fee insert: %w", err)
	}
	return &fee, nil
}

// GetOrgFees return fee records matched against org ID
func (s *Storage) GetOrgFees(ctx context.Context, id string) ([]storage.FeeCommission, error) {
	const getFees = `SELECT * FROM fee_commission WHERE org_profile_id = $1 AND deleted is null`
	var fs []storage.FeeCommission
	if err := s.db.Select(&fs, getFees, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	for _, f := range fs {
		f.FeeStatus = 1 // Active
		if time.Now().Before(f.StartDate.Time) || time.Now().After(f.EndDate.Time) {
			f.FeeStatus = 2 // Disabled
		}
	}
	return fs, nil
}

// ListOrgFees return fee records matched against org ID
func (s *Storage) ListOrgFees(ctx context.Context, oid string, f storage.LimitOffsetFilter) ([]storage.FeeCommission, error) {
	const getFees = `
	WITH cnt AS (select count(*) as count FROM fee_commission WHERE org_id = $1)
	SELECT *, cnt.count
	FROM fee_commission as p left join cnt on true
	WHERE org_id = $1 AND fee_commision_type = $4
	AND deleted is null
	ORDER BY created ASC
	LIMIT NULLIF($2, 0)
	OFFSET $3;
	`

	var fs []storage.FeeCommission
	if err := s.db.Select(&fs, getFees, oid, f.Limit, f.Offset, f.Type); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}

	for i := range fs {
		fs[i].FeeStatus = 1 // Active
		if time.Now().Before(fs[i].StartDate.Time) || time.Now().After(fs[i].EndDate.Time) {
			fs[i].FeeStatus = 2 // Disabled
		}
	}
	return fs, nil
}

const insertRate = `
INSERT INTO fee_commission_rate (
    fee_commission_id,
	min_volume,
	max_volume,
	txn_rate
) VALUES (
	 :fee_commission_id,
	 :min_volume,
	 :max_volume,
	 :txn_rate
) RETURNING
    id
`

// CreateFeeCommissionRate ...
func (s *Storage) CreateFeeCommissionRate(ctx context.Context, rate storage.Rate) (*storage.Rate, error) {
	log := logging.FromContext(ctx)
	stmt, err := s.prepareNamed(ctx, insertRate)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&rate, rate); err != nil {
		logging.WithError(err, log)
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, fmt.Errorf("executing rate insert: %w", err)
	}
	return &rate, nil
}

// ListFeesCommissionsRate return rate records matched against feescommison ID
func (s *Storage) ListFeesCommissionRate(ctx context.Context, fcid string) ([]storage.Rate, error) {
	const getRates = `SELECT * FROM fee_commission_rate WHERE fee_commission_id = $1`

	var fs []storage.Rate
	if err := s.db.Select(&fs, getRates, fcid); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return fs, nil
}

const feeUpdate = `
UPDATE
	fee_commission
SET
	fee_amount= COALESCE(NULLIF(:fee_amount, ''), fee_amount),
	commission_amount= COALESCE(NULLIF(:commission_amount, ''), commission_amount),
	fee_commision_type= COALESCE(:fee_commision_type, fee_commision_type),
	start_date= COALESCE(:start_date, start_date),
	end_date= COALESCE(:end_date, end_date),
	deleted= COALESCE(:deleted, deleted)
WHERE
	id = :id
RETURNING
   id,created,updated
`

func (s *Storage) UpsertOrgFees(ctx context.Context, fee storage.FeeCommission) (*storage.FeeCommission, error) {
	if fee.ID == "" {
		return s.CreateOrgFees(ctx, fee)
	}
	log := logging.FromContext(ctx)
	stmt, err := s.prepareNamed(ctx, feeUpdate)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&fee, fee); err != nil {
		logging.WithError(err, log)
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, fmt.Errorf("executing fee insert: %w", err)
	}
	return &fee, nil
}

const rateUpdate = `
UPDATE
	fee_commission_rate
SET
	min_volume= COALESCE(NULLIF(:min_volume, ''), min_volume),
	max_volume= COALESCE(NULLIF(:max_volume, ''), max_volume),
	txn_rate= COALESCE(NULLIF(:txn_rate, ''), txn_rate)
WHERE
	id = :id
RETURNING
   id
`

func (s *Storage) UpsertRate(ctx context.Context, rate storage.Rate) (*storage.Rate, error) {
	if rate.ID == "" {
		return s.CreateFeeCommissionRate(ctx, rate)
	}
	log := logging.FromContext(ctx)
	stmt, err := s.prepareNamed(ctx, rateUpdate)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&rate, rate); err != nil {
		logging.WithError(err, log)
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, fmt.Errorf("executing rate insert: %w", err)
	}
	return &rate, nil
}

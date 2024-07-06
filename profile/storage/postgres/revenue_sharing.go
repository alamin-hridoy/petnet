package postgres

import (
	"context"
	"database/sql"
	"fmt"

	revsha "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const insertRevenueSharing = `
INSERT INTO revenue_sharing (
	org_id,
    user_id,
	partner,
	bound_type,
	remit_type,
	transaction_type,
	tier_type,
	amount,
	created_by,
	updated_by
) VALUES (
	:org_id,
	:user_id,
	:partner,
	:bound_type,
	:remit_type,
	:transaction_type,
	:tier_type,
	:amount,
	:created_by,
	:updated_by
) RETURNING
    id
`

func (s *Storage) CreateRevenueSharing(ctx context.Context, req storage.RevenueSharing) (*storage.RevenueSharing, error) {
	switch {
	case req.OrgID == "":
		return nil, fmt.Errorf("org id cannot be empty")
	case req.UserID == "":
		return nil, fmt.Errorf("user id cannot be empty")
	case req.BoundType == "":
		return nil, fmt.Errorf("bound type cannot be empty")
	case req.RemitType == "":
		return nil, fmt.Errorf("remit type cannot be empty")
	case req.Partner == "":
		return nil, fmt.Errorf("partner cannot be empty")
	case req.TransactionType == "":
		return nil, fmt.Errorf("transaction type cannot be empty")
	case req.TierType == "":
		return nil, fmt.Errorf("tier type type cannot be empty")
	}
	log := logging.FromContext(ctx)
	stmt, err := s.prepareNamed(ctx, insertRevenueSharing)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&req, req); err != nil {
		logging.WithError(err, log)
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, fmt.Errorf("executing revenue sharing insert: %w", err)
	}
	return &req, nil
}

const revenueSharingUpdate = `
UPDATE
	revenue_sharing
SET
	tier_type = COALESCE(NULLIF(:tier_type, ''), tier_type),
	amount = :amount,
	updated = :updated,
	updated_by = COALESCE(NULLIF(:updated_by, ''), updated_by)
WHERE
	id = :id AND org_id = :org_id AND user_id = :user_id AND partner = :partner AND transaction_type = :transaction_type AND remit_type = :remit_type AND bound_type = :bound_type
RETURNING
   id
`

func (s *Storage) UpdateRevenueSharing(ctx context.Context, req storage.RevenueSharing) (*storage.RevenueSharing, error) {
	switch {
	case req.OrgID == "":
		return nil, fmt.Errorf("org id cannot be empty")
	case req.UserID == "":
		return nil, fmt.Errorf("user id cannot be empty")
	case req.BoundType == "":
		return nil, fmt.Errorf("bound type cannot be empty")
	case req.ID == "":
		return nil, fmt.Errorf("id cannot be empty")
	case req.Partner == "":
		return nil, fmt.Errorf("partner cannot be empty")
	case req.TransactionType == "":
		return nil, fmt.Errorf("transaction type cannot be empty")
	case req.RemitType == "":
		return nil, fmt.Errorf("remit type cannot be empty")
	}
	log := logging.FromContext(ctx)
	stmt, err := s.prepareNamed(ctx, revenueSharingUpdate)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&req, req); err != nil {
		logging.WithError(err, log)
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, fmt.Errorf("executing revenue sharing update: %w", err)
	}
	return &req, nil
}

const insertRevenueSharingTier = `
INSERT INTO revenue_sharing_tier (
	revenue_sharing_id,
	min_value,
	max_value,
	amount
) VALUES (
	:revenue_sharing_id,
	:min_value,
	:max_value,
	:amount
) RETURNING
    id
`

func (s *Storage) CreateRevenueSharingTier(ctx context.Context, req storage.RevenueSharingTier) (*storage.RevenueSharingTier, error) {
	switch {
	case req.RevenueSharingID == "":
		return nil, fmt.Errorf("revenue sharing id cannot be empty")
	case req.MinValue == "":
		return nil, fmt.Errorf("min value cannot be empty")
	case req.MaxValue == "":
		return nil, fmt.Errorf("max value cannot be empty")
	case req.Amount == "":
		return nil, fmt.Errorf("amount cannot be empty")
	}
	log := logging.FromContext(ctx)
	stmt, err := s.prepareNamed(ctx, insertRevenueSharingTier)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&req, req); err != nil {
		logging.WithError(err, log)
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, fmt.Errorf("executing revenue sharing tier insert: %w", err)
	}
	return &req, nil
}

const revenueSharingTierUpdate = `
UPDATE
	revenue_sharing_tier
SET
	min_value= COALESCE(NULLIF(:min_value, ''), min_value),
	max_value= COALESCE(NULLIF(:max_value, ''), max_value),
	amount= COALESCE(NULLIF(:amount, ''), amount)
WHERE
	id = :id AND revenue_sharing_id = :revenue_sharing_id
RETURNING
	id, revenue_sharing_id
`

func (s *Storage) UpdateRevenueSharingTier(ctx context.Context, req storage.RevenueSharingTier) (*storage.RevenueSharingTier, error) {
	switch {
	case req.ID == "":
		return nil, fmt.Errorf("id cannot be empty")
	case req.RevenueSharingID == "":
		return nil, fmt.Errorf("revenue sharing id cannot be empty")
	case req.MinValue == "":
		return nil, fmt.Errorf("min value cannot be empty")
	case req.MaxValue == "":
		return nil, fmt.Errorf("max value cannot be empty")
	case req.Amount == "":
		return nil, fmt.Errorf("amount cannot be empty")
	}
	log := logging.FromContext(ctx)
	stmt, err := s.prepareNamed(ctx, revenueSharingTierUpdate)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&req, req); err != nil {
		logging.WithError(err, log)
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, fmt.Errorf("executing revenue sharing tier update: %w", err)
	}
	return &req, nil
}

func (s *Storage) GetRevenueSharingList(ctx context.Context, req storage.RevenueSharing) ([]storage.RevenueSharing, error) {
	switch {
	case req.OrgID == "":
		return nil, fmt.Errorf("org id cannot be empty")
	case req.RemitType == "":
		return nil, fmt.Errorf("remit type cannot be empty")
	}

	boundQ := ""
	if req.BoundType != "" && req.BoundType != revsha.BoundType_EMPTYBOUNDTYPE.String() {
		boundQ = " AND bound_type = :bound_type "
	}
	partnerQ := ""
	if req.Partner != "" {
		partnerQ = " AND partner = :partner "
	}
	revenueSharingByRemitType := `
	WITH cnt AS (select count(*) as count FROM revenue_sharing WHERE org_id = :org_id AND remit_type = :remit_type %s %s)
	SELECT *, cnt.count
	FROM revenue_sharing as p left join cnt on true
	WHERE org_id = :org_id AND remit_type = :remit_type %s %s
	ORDER BY created ASC
	`
	revenueSharingByRemitType = fmt.Sprintf(revenueSharingByRemitType, boundQ, partnerQ, boundQ, partnerQ)
	var com []storage.RevenueSharing
	stmt, err := s.db.PrepareNamed(revenueSharingByRemitType)
	if err != nil {
		return nil, fmt.Errorf("preparing named query Get revenue sharing List: %w", err)
	}
	if err := stmt.Select(&com, req); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return com, nil
}

func (s *Storage) GetRevenueSharingTierList(ctx context.Context, req storage.RevenueSharingTier) ([]storage.RevenueSharingTier, error) {
	switch {
	case req.RevenueSharingID == "":
		return nil, fmt.Errorf("revenue sharing id cannot be empty")
	}
	const revenueSharingTier = `SELECT * FROM revenue_sharing_tier WHERE revenue_sharing_id = :revenue_sharing_id ORDER BY cast(min_value AS float) ASC`
	var comTier []storage.RevenueSharingTier
	stmt, err := s.db.PrepareNamed(revenueSharingTier)
	if err != nil {
		return nil, fmt.Errorf("preparing named query Get revenue sharing tier List: %w", err)
	}
	if err := stmt.Select(&comTier, req); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return comTier, nil
}

func (s *Storage) DeleteRevenueSharing(ctx context.Context, req storage.RevenueSharing) error {
	switch {
	case req.OrgID == "":
		return fmt.Errorf("org id cannot be empty")
	case req.UserID == "":
		return fmt.Errorf("user id cannot be empty")
	case req.Partner == "":
		return fmt.Errorf("partner cannot be empty")
	case req.TransactionType == "":
		return fmt.Errorf("transaction type cannot be empty")
	case req.RemitType == "":
		return fmt.Errorf("remit type cannot be empty")
	case req.BoundType == "":
		return fmt.Errorf("bound type cannot be empty")
	}
	if _, err := s.db.Exec("DELETE FROM revenue_sharing WHERE org_id = $1 AND user_id = $2 AND partner = $3 AND transaction_type = $4 AND remit_type = $5 AND bound_type = $6", req.OrgID, req.UserID, req.Partner, req.TransactionType, req.RemitType, req.BoundType); err != nil {
		if err == sql.ErrNoRows {
			return storage.NotFound
		}
		return status.Errorf(codes.Internal, "delete revenue sharing failed")
	}
	return nil
}

func (s *Storage) DeleteRevenueSharingTier(ctx context.Context, req storage.RevenueSharingTier) error {
	switch {
	case req.RevenueSharingID == "":
		return fmt.Errorf("revenue sharing id cannot be empty")
	}
	if _, err := s.db.Exec("DELETE FROM revenue_sharing_tier WHERE revenue_sharing_id = $1", req.RevenueSharingID); err != nil {
		if err == sql.ErrNoRows {
			return storage.NotFound
		}
		return status.Errorf(codes.Internal, "delete revenue sharing tier failed")
	}
	return nil
}

func (s *Storage) DeleteRevenueSharingTierById(ctx context.Context, req storage.RevenueSharingTier) error {
	switch {
	case req.ID == "":
		return fmt.Errorf("revenue sharing tier id cannot be empty")
	}
	if _, err := s.db.Exec("DELETE FROM revenue_sharing_tier WHERE id = $1", req.ID); err != nil {
		if err == sql.ErrNoRows {
			return storage.NotFound
		}
		return status.Errorf(codes.Internal, "delete revenue sharing tier failed")
	}
	return nil
}

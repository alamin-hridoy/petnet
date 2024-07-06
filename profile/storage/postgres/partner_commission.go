package postgres

import (
	"context"
	"database/sql"
	"fmt"

	ptnrcom "brank.as/petnet/gunk/dsa/v2/partnercommission"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const insertPartnerCommission = `
INSERT INTO partner_commission_config (
	partner,
	bound_type,
	remit_type,
	transaction_type,
	tier_type,
	amount,
	start_date,
	end_date,
	created_by,
	updated_by
) VALUES (
	:partner,
	:bound_type,
	:remit_type,
	:transaction_type,
	:tier_type,
	:amount,
	:start_date,
	:end_date,
	:created_by,
	:updated_by
) RETURNING
    id
`

// CreatePartnerCommission ...
func (s *Storage) CreatePartnerCommission(ctx context.Context, req storage.PartnerCommission) (*storage.PartnerCommission, error) {
	switch {
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
	stmt, err := s.prepareNamed(ctx, insertPartnerCommission)
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
		return nil, fmt.Errorf("executing partner commission insert: %w", err)
	}
	return &req, nil
}

const partnerCommissionUpdate = `
UPDATE
	partner_commission_config
SET
	tier_type = COALESCE(NULLIF(:tier_type, ''), tier_type),
	amount = :amount,
	start_date = :start_date,
	end_date = :end_date,
	updated = :updated,
	updated_by = COALESCE(NULLIF(:updated_by, ''), updated_by)
WHERE
	id = :id AND partner = :partner AND transaction_type = :transaction_type AND remit_type = :remit_type AND bound_type = :bound_type
RETURNING
   id
`

func (s *Storage) UpdatePartnerCommission(ctx context.Context, req storage.PartnerCommission) (*storage.PartnerCommission, error) {
	switch {
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
	stmt, err := s.prepareNamed(ctx, partnerCommissionUpdate)
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
		return nil, fmt.Errorf("executing partner commission update: %w", err)
	}
	return &req, nil
}

const insertPartnerCommissionTier = `
INSERT INTO partner_commission_tier (
	partner_commission_config_id,
	min_value,
	max_value,
	amount
) VALUES (
	:partner_commission_config_id,
	:min_value,
	:max_value,
	:amount
) RETURNING
    id
`

// InsertCommissionTier ...
func (s *Storage) CreatePartnerCommissionTier(ctx context.Context, req storage.PartnerCommissionTier) (*storage.PartnerCommissionTier, error) {
	switch {
	case req.PartnerCommissionID == "":
		return nil, fmt.Errorf("partner commission id cannot be empty")
	case req.MinValue == "":
		return nil, fmt.Errorf("min value cannot be empty")
	case req.MaxValue == "":
		return nil, fmt.Errorf("max value cannot be empty")
	case req.Amount == "":
		return nil, fmt.Errorf("amount cannot be empty")
	}
	log := logging.FromContext(ctx)
	stmt, err := s.prepareNamed(ctx, insertPartnerCommissionTier)
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
		return nil, fmt.Errorf("executing commision tier insert: %w", err)
	}
	return &req, nil
}

const partnerCommissionTierUpdate = `
UPDATE
	partner_commission_tier
SET
	min_value= COALESCE(NULLIF(:min_value, ''), min_value),
	max_value= COALESCE(NULLIF(:max_value, ''), max_value),
	amount= COALESCE(NULLIF(:amount, ''), amount)
WHERE
	id = :id AND partner_commission_config_id = :partner_commission_config_id
RETURNING
	id, partner_commission_config_id
`

func (s *Storage) UpdatePartnerCommissionTier(ctx context.Context, req storage.PartnerCommissionTier) (*storage.PartnerCommissionTier, error) {
	switch {
	case req.ID == "":
		return nil, fmt.Errorf("id cannot be empty")
	case req.PartnerCommissionID == "":
		return nil, fmt.Errorf("partner commission id cannot be empty")
	case req.MinValue == "":
		return nil, fmt.Errorf("min value cannot be empty")
	case req.MaxValue == "":
		return nil, fmt.Errorf("max value cannot be empty")
	case req.Amount == "":
		return nil, fmt.Errorf("amount cannot be empty")
	}
	log := logging.FromContext(ctx)
	stmt, err := s.prepareNamed(ctx, partnerCommissionTierUpdate)
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
		return nil, fmt.Errorf("executing partner commission tier update: %w", err)
	}
	return &req, nil
}

// ListPartnerCommissions
func (s *Storage) GetPartnerCommissionsList(ctx context.Context, req storage.PartnerCommission) ([]storage.PartnerCommission, error) {
	switch {
	case req.RemitType == "":
		return nil, fmt.Errorf("remit type cannot be empty")
	}
	boundQ := ""
	if req.BoundType != "" && req.BoundType != ptnrcom.BoundType_EMPTYBOUNDTYPE.String() {
		boundQ = " AND bound_type = :bound_type "
	}
	partnerQ := ""
	if req.Partner != "" {
		partnerQ = " AND partner = :partner "
	}
	partnerCommissionByRemitType := `
	WITH cnt AS (select count(*) as count FROM partner_commission_config WHERE remit_type = :remit_type %s %s)
	SELECT *, cnt.count
	FROM partner_commission_config as p left join cnt on true
	WHERE remit_type = :remit_type %s %s
	ORDER BY created ASC
	`
	partnerCommissionByRemitType = fmt.Sprintf(partnerCommissionByRemitType, boundQ, partnerQ, boundQ, partnerQ)
	var com []storage.PartnerCommission
	stmt, err := s.db.PrepareNamed(partnerCommissionByRemitType)
	if err != nil {
		return nil, fmt.Errorf("preparing named query Get commission List: %w", err)
	}
	if err := stmt.Select(&com, req); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return com, nil
}

// GetPartnerCommissionsTierList
func (s *Storage) GetPartnerCommissionsTierList(ctx context.Context, req storage.PartnerCommissionTier) ([]storage.PartnerCommissionTier, error) {
	switch {
	case req.PartnerCommissionID == "":
		return nil, fmt.Errorf("partner commission id cannot be empty")
	}
	const partnerCommissionTier = `SELECT * FROM partner_commission_tier WHERE partner_commission_config_id = :partner_commission_config_id ORDER BY cast(min_value AS float) ASC`
	var comTier []storage.PartnerCommissionTier
	stmt, err := s.db.PrepareNamed(partnerCommissionTier)
	if err != nil {
		return nil, fmt.Errorf("preparing named query Get partner commission tier List: %w", err)
	}
	if err := stmt.Select(&comTier, req); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return comTier, nil
}

// GetPartnerCommission
func (s *Storage) GetPartnerCommission(ctx context.Context, req storage.PartnerCommission) (*storage.PartnerCommission, error) {
	switch {
	case req.Partner == "":
		return nil, fmt.Errorf("partner cannot be empty")
	case req.TransactionType == "":
		return nil, fmt.Errorf("transaction type cannot be empty")
	case req.RemitType == "":
		return nil, fmt.Errorf("remit type cannot be empty")
	case req.BoundType == "":
		return nil, fmt.Errorf("bound type cannot be empty")
	}
	const getPartnerCommission = `Select * FROM partner_commission_config WHERE partner = $1 AND transaction_type = $2 AND remit_type = $3 AND bound_type = $4`
	if err := s.db.Get(&req, getPartnerCommission, req.Partner, req.TransactionType, req.RemitType, req.BoundType); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return &req, nil
}

// PartnerCommission
func (s *Storage) DeletePartnerCommission(ctx context.Context, req storage.PartnerCommission) error {
	switch {
	case req.Partner == "":
		return fmt.Errorf("partner cannot be empty")
	case req.TransactionType == "":
		return fmt.Errorf("transaction type cannot be empty")
	case req.RemitType == "":
		return fmt.Errorf("remit type cannot be empty")
	case req.BoundType == "":
		return fmt.Errorf("bound type cannot be empty")
	}
	gpc, err := s.GetPartnerCommission(ctx, req)
	if err != nil {
		if err == storage.NotFound {
			return status.Error(codes.NotFound, "partner commission not found")
		}
		return status.Errorf(codes.Internal, "get partner commission failed")
	}
	if gpc.ID != "" {
		err := s.DeletePartnerCommissionTier(ctx, storage.PartnerCommissionTier{
			PartnerCommissionID: gpc.ID,
		})
		if err != nil {
			if err != storage.NotFound {
				return err
			}
		}
	}
	if _, err := s.db.Exec("DELETE FROM partner_commission_config WHERE partner = $1 AND transaction_type = $2 AND remit_type = $3 AND bound_type = $4", req.Partner, req.TransactionType, req.RemitType, req.BoundType); err != nil {
		if err == sql.ErrNoRows {
			return storage.NotFound
		}
		return status.Errorf(codes.Internal, "delete partner commission failed")
	}
	return nil
}

// deleteCommissionTier
func (s *Storage) DeletePartnerCommissionTier(ctx context.Context, req storage.PartnerCommissionTier) error {
	switch {
	case req.PartnerCommissionID == "":
		return fmt.Errorf("partner commission id cannot be empty")
	}
	if _, err := s.db.Exec("DELETE FROM partner_commission_tier WHERE partner_commission_config_id = $1", req.PartnerCommissionID); err != nil {
		if err == sql.ErrNoRows {
			return storage.NotFound
		}
		return status.Errorf(codes.Internal, "delete partner commission tier failed")
	}
	return nil
}

// deleteCommissionTierById
func (s *Storage) DeletePartnerCommissionTierById(ctx context.Context, req storage.PartnerCommissionTier) error {
	switch {
	case req.ID == "":
		return fmt.Errorf("partner commission tier id cannot be empty")
	}
	if _, err := s.db.Exec("DELETE FROM partner_commission_tier WHERE id = $1", req.ID); err != nil {
		if err == sql.ErrNoRows {
			return storage.NotFound
		}
		return status.Errorf(codes.Internal, "delete partner commission tier failed")
	}
	return nil
}

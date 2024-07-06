package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"brank.as/petnet/api/storage"
	"github.com/lib/pq"
)

const createRemittanceHistory = `
INSERT INTO remittance_history (
	dsa_id,
	user_id,
	phrn,
	send_validate_reference_number,
	cancel_send_reference_number,
	payout_validate_reference_number,
	txn_status,
	error_code,
	error_message,
	error_time,
	error_type,
	details,
	remarks,
	txn_created_time,
	txn_updated_time,
	txn_confirm_time
) VALUES (
	:dsa_id,
	:user_id,
	:phrn,
	:send_validate_reference_number,
	:cancel_send_reference_number,
	:payout_validate_reference_number,
	'VALIDATE_SEND',
	:error_code,
	:error_message,
	:error_time,
	:error_type,
	:details,
	:remarks,
	:txn_created_time,
	:txn_updated_time,
	:txn_confirm_time
) RETURNING *`

func (s *Storage) CreateRemittanceHistory(ctx context.Context, r storage.PerahubRemittanceHistory) (*storage.PerahubRemittanceHistory, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, createRemittanceHistory)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, fmt.Errorf("executing remittance history insert: %w", err)
	}
	return &r, nil
}

const updateRemittanceHistory = `
UPDATE remittance_history
SET
	phrn = COALESCE(NULLIF(:phrn, ''), phrn),
	send_validate_reference_number = COALESCE(NULLIF(:send_validate_reference_number, ''), send_validate_reference_number),
	cancel_send_reference_number = COALESCE(NULLIF(:cancel_send_reference_number, ''), cancel_send_reference_number),
	payout_validate_reference_number = COALESCE(NULLIF(:payout_validate_reference_number, ''), payout_validate_reference_number),
	txn_status = COALESCE(NULLIF(:txn_status, ''), txn_status),
	error_code = COALESCE(NULLIF(:error_code, ''), error_code),
	error_message = COALESCE(NULLIF(:error_message, ''), error_message),
	error_time = COALESCE(NULLIF(:error_time, ''), error_time),
	error_type = COALESCE(NULLIF(:error_type, ''), error_type),
	remarks = COALESCE(NULLIF(:remarks, ''), remarks),
	txn_created_time = :txn_created_time,
	txn_updated_time = :txn_updated_time,
	txn_confirm_time = :txn_confirm_time
	WHERE remittance_history_id = :remittance_history_id
RETURNING *
`

func (s *Storage) UpdateRemittanceHistory(ctx context.Context, r storage.PerahubRemittanceHistory) (*storage.PerahubRemittanceHistory, error) {
	if r.RemittanceHistoryID == "" {
		return nil, fmt.Errorf("remittance history id cannot be empty")
	}

	stmt, err := s.db.PrepareNamedContext(ctx, updateRemittanceHistory)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		return nil, fmt.Errorf("executing remittance history update: %w", err)
	}
	return &r, nil
}

func (s *Storage) GetRemittanceHistory(ctx context.Context, RemittanceHistoryDetailsID string) (*storage.PerahubRemittanceHistory, error) {
	const getRemittanceHistory = `SELECT * FROM remittance_history WHERE remittance_history_id = $1`
	var r storage.PerahubRemittanceHistory
	if err := s.db.Get(&r, getRemittanceHistory, RemittanceHistoryDetailsID); err != nil {
		if err == sql.ErrNoRows {
			return &r, nil
		}
		return nil, fmt.Errorf("executing remittance history get: %w", err)
	}
	return &r, nil
}

func (s *Storage) ListRemittanceHistory(ctx context.Context, f storage.RemittanceHistoryFilter) ([]storage.PerahubRemittanceHistory, error) {
	ft := ""
	ut := ""
	tr := ""
	tc := ":total_count_where:"
	b := NewBuilder(fmt.Sprintf(`WITH cnt AS (select count(*) as total FROM remittance_history %s) SELECT *, cnt.total FROM remittance_history left join cnt on true`, tc))
	if !f.From.IsZero() {
		ft = f.From.Format("2006-01-02")
	}
	if !f.Until.IsZero() {
		ut = f.Until.Format("2006-01-02")
	}
	if !f.TxnConfirmTime.IsZero() {
		tr = f.TxnConfirmTime.Format("2006-01-02")
	}

	var scol string
	switch string(f.SortByColumn) {
	case string(storage.RemittanceHistoryIDCol):
		scol = "remittance_history_id"
	case string(storage.DsaColID):
		scol = "dsa_id"
	case string(storage.PhrnIDCol):
		scol = "phrn"
	case string(storage.UserColID):
		scol = "user_id"
	default:
		scol = ""
	}

	b.Where("remittance_history_id", eq, f.RemittanceHistoryID).
		Where("phrn", eq, f.Phrn).
		Where("txn_status", eq, f.TxnStatus).
		Where("user_id", eq, f.UserID).
		Where("txn_confirm_time", gtOrEq, ft, CompareDate()).
		Where("txn_confirm_time", ltOrEq, ut, CompareDate()).
		Where("txn_confirm_time", eq, tr, CompareDate()).
		SortByColumn(scol, f.SortOrder).
		Limit(f.Limit).
		Offset(f.Offset).AddTotalQuery(tc)

	stmt, err := s.db.PrepareNamed(b.query)
	if err != nil {
		return nil, err
	}

	r := []storage.PerahubRemittanceHistory{}
	if err := stmt.Select(&r, b.args); err != nil {
		return nil, fmt.Errorf("executing remittance history list: %w", err)
	}
	return r, nil
}

const confirmRemittanceHistory = `
UPDATE remittance_history
SET
	phrn = :phrn,
	txn_status = 'CONFIRM_SEND',
	txn_updated_time = now()
	WHERE send_validate_reference_number = :send_validate_reference_number
RETURNING *
`

func (s *Storage) ConfirmRemittanceHistory(ctx context.Context, r storage.PerahubRemittanceHistory) (*storage.PerahubRemittanceHistory, error) {
	if r.SendValidateReferenceNumber == "" {
		return nil, fmt.Errorf("send validate reference number cannot be empty")
	}
	if r.Phrn == "" {
		return nil, fmt.Errorf("phrn cannot be empty")
	}
	stmt, err := s.db.PrepareNamedContext(ctx, confirmRemittanceHistory)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		return nil, fmt.Errorf("executing remittance history confirm: %w", err)
	}
	return &r, nil
}

const cancelRemittanceHistory = `
UPDATE remittance_history
SET
	cancel_send_reference_number = :cancel_send_reference_number,
	txn_status = 'CANCEL_SEND',
	remarks = :remarks,
	txn_updated_time = now()
	WHERE phrn = :phrn
RETURNING *
`

func (s *Storage) CancelRemittanceHistory(ctx context.Context, r storage.PerahubRemittanceHistory) (*storage.PerahubRemittanceHistory, error) {
	if r.Phrn == "" {
		return nil, fmt.Errorf("phrn cannot be empty")
	}
	if r.CancelSendReferenceNumber == "" {
		return nil, fmt.Errorf("cancel send reference number cannot be empty")
	}
	stmt, err := s.db.PrepareNamedContext(ctx, cancelRemittanceHistory)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		return nil, fmt.Errorf("executing remittance history cancel: %w", err)
	}
	return &r, nil
}

const validateReceiveRemittanceHistory = `
UPDATE remittance_history
SET
	payout_validate_reference_number = :payout_validate_reference_number,
	txn_status = 'VALIDATE_RECEIVE',
	txn_updated_time = now()
	WHERE phrn = :phrn
RETURNING *
`

func (s *Storage) ValidateReceiveRemittanceHistory(ctx context.Context, r storage.PerahubRemittanceHistory) (*storage.PerahubRemittanceHistory, error) {
	if r.Phrn == "" {
		return nil, fmt.Errorf("phrn cannot be empty")
	}
	if r.PayoutValidateReferenceNumber == "" {
		return nil, fmt.Errorf("payout validate reference number cannot be empty")
	}
	stmt, err := s.db.PrepareNamedContext(ctx, validateReceiveRemittanceHistory)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		return nil, fmt.Errorf("executing remittance history cancel: %w", err)
	}
	return &r, nil
}

const confirmReceiveRemittanceHistory = `
UPDATE remittance_history
SET
	txn_status = 'CONFIRM_RECEIVE',
	txn_confirm_time = now()
	WHERE payout_validate_reference_number = :payout_validate_reference_number
RETURNING *
`

func (s *Storage) ConfirmReceiveRemittanceHistory(ctx context.Context, r storage.PerahubRemittanceHistory) (*storage.PerahubRemittanceHistory, error) {
	if r.PayoutValidateReferenceNumber == "" {
		return nil, fmt.Errorf("payout validate reference number cannot be empty")
	}
	stmt, err := s.db.PrepareNamedContext(ctx, confirmReceiveRemittanceHistory)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		return nil, fmt.Errorf("executing remittance history cancel: %w", err)
	}
	return &r, nil
}

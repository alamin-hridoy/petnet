package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"brank.as/petnet/api/storage"
	"brank.as/petnet/serviceutil/logging"
	"github.com/lib/pq"
)

const createRemitHistory = `
INSERT INTO remit_history (
	remit_id,
	dsa_id,
	dsa_order_id,
	user_id,
	remco_id,
	sender_member_id,
	receiver_member_id,
	remco_control_number,
	remittance,
	remit_type,
	txn_staged_time,
	txn_status,
	txn_step,
	error_code,
	error_message,
	error_time,
	error_type,
	transaction_type
) VALUES (
	:remit_id,
	:dsa_id,
	:dsa_order_id,
	:user_id,
	:remco_id,
	:sender_member_id,
	:receiver_member_id,
	:remco_control_number,
	:remittance,
	:remit_type,
	:txn_staged_time,
	:txn_status,
	:txn_step,
	:error_code,
	:error_message,
	:error_time,
	:error_type,
	:transaction_type
) RETURNING
updated
`

func (s *Storage) CreateRemitHistory(ctx context.Context, r storage.RemitHistory) (*storage.RemitHistory, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, createRemitHistory)
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

const updateRemitHistory = `
UPDATE remit_history
SET
	txn_status = :txn_status,
	txn_step = :txn_step,
	txn_completed_time = :txn_completed_time,
	remco_control_number = :remco_control_number,
	remittance = :remittance,
	error_code = :error_code,
	error_message = :error_message,
	error_time = :error_time,
	error_type = :error_type,
	transaction_type = :transaction_type
WHERE remit_id = :remit_id
RETURNING updated
`

func (s *Storage) UpdateRemitHistory(ctx context.Context, r storage.RemitHistory) (*storage.RemitHistory, error) {
	if r.TxnID == "" {
		return nil, fmt.Errorf("transaction ID cannot be empty")
	}
	stmt, err := s.db.PrepareNamedContext(ctx, updateRemitHistory)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		return nil, fmt.Errorf("executing remittance history update: %w", err)
	}
	return &r, nil
}

const updateRemitHistoryDate = `
UPDATE remit_history
SET
	txn_completed_time = $2
WHERE remit_id = $1
`

// UpdateRemitHistoryDate used for testing only
func (s *Storage) UpdateRemitHistoryDate(ctx context.Context, txnID string, t time.Time) error {
	if txnID == "" {
		return fmt.Errorf("transaction ID cannot be empty")
	}
	if _, err := s.db.Exec(updateRemitHistoryDate, txnID, t); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetRemitHistory(ctx context.Context, txnID string) (*storage.RemitHistory, error) {
	const getRemitHistory = `SELECT * FROM remit_history WHERE remit_id = $1`
	var r storage.RemitHistory
	if err := s.db.Get(&r, getRemitHistory, txnID); err != nil {
		if err == sql.ErrNoRows {
			return &r, nil
		}
		return nil, fmt.Errorf("executing remittance get history: %w", err)
	}
	return &r, nil
}

func (s *Storage) ListRemitHistory(ctx context.Context, f storage.LRHFilter) ([]storage.RemitHistory, error) {
	b := NewBuilder("SELECT *, count(*) OVER() AS total FROM remit_history")
	ft := ""
	ut := ""
	if !f.From.IsZero() {
		ft = f.From.Format("2006-01-02")
	}
	if !f.Until.IsZero() {
		ut = f.Until.Format("2006-01-02")
	}

	var scol string
	switch f.SortByColumn {
	case storage.ControlNumberCol:
		scol = "remco_control_number"
	case storage.RemittedToCol:
		scol = "remittance->'receiver'->>'first_name'"
	case storage.TotalRemittedAmountCol:
		scol = "remittance->'source_amt'->>'amount'"
	case storage.TransactionTimeCol:
		scol = "txn_completed_time"
	case storage.TransactionCompletedTime:
		scol = "txn_completed_time"
	case storage.UserIDCol:
		scol = "user_id"
	case storage.PartnerCol:
		scol = "remco_id"
	default:
		scol = ""
	}
	exPtnrs := stringToSlice(f.ExcludePartner)
	extyps := stringToSlice(f.ExcludeType)
	b.Any("remco_control_number", f.ControlNo).
		Where("dsa_id", eq, f.DsaOrgID).
		Where("remco_id", eq, f.Partner).
		NotAny("remco_id", exPtnrs).
		NotAny("remit_type", extyps).
		Where("remco_id", notEq, f.ExcludePartner).
		Where("txn_step", eq, f.TxnStep).
		Where("txn_status", eq, f.TxnStatus).
		Where("txn_completed_time", gtOrEq, ft, CompareDate()).
		Where("txn_completed_time", ltOrEq, ut, CompareDate()).
		Where("transaction_type", eq, f.Transactiontype).
		SortByColumn(scol, f.SortOrder).
		Limit(f.Limit).
		Offset(f.Offset)

	stmt, err := s.db.PrepareNamed(b.query)
	if err != nil {
		return nil, err
	}

	r := []storage.RemitHistory{}
	if err := stmt.Select(&r, b.args); err != nil {
		return nil, fmt.Errorf("executing remittance list history: %w", err)
	}
	return r, nil
}

func (s *Storage) OrderIDExists(ctx context.Context, ordID string) bool {
	log := logging.FromContext(ctx)
	const getRemitHistory = `SELECT * FROM remit_history WHERE dsa_order_id = $1`
	var r storage.RemitHistory
	if err := s.db.Get(&r, getRemitHistory, ordID); err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		logging.WithError(err, log).Error("selecting order id")
		return false
	}
	return true
}

func (s *Storage) GetTransactionReport(ctx context.Context, pf *storage.LRHFilter) (*storage.RemitHistory, error) {
	switch {
	case pf.TxnStatus == "":
		return nil, fmt.Errorf("transaction status cannot be empty")
	case pf.TxnStep == "":
		return nil, fmt.Errorf("transaction step cannot be empty")
	case pf.From.IsZero():
		return nil, fmt.Errorf("from date cannot empty")
	case pf.Until.IsZero():
		return nil, fmt.Errorf("until date cannot empty")
	}
	ft := ""  // from time
	ut := ""  // until time
	ftQ := "" // from time query
	utQ := "" // until time query
	ttQ := "" // transaction type query
	if !pf.From.IsZero() {
		ft = pf.From.Format("2006-01-02")
	}
	if !pf.Until.IsZero() {
		ut = pf.Until.Format("2006-01-02")
	}
	if ft != "" {
		ftQ = fmt.Sprintf(" and cast(txn_completed_time as date) >= '%s' ", ft)
	}
	if ut != "" {
		utQ = fmt.Sprintf(" and cast(txn_completed_time as date) <= '%s' ", ut)
	}
	if pf.Transactiontype != "" {
		ttQ = fmt.Sprintf(" and transaction_type = '%s' ", pf.Transactiontype)
	}
	getTransactionReport := fmt.Sprintf(`select count(remit_id) as remit_transaction_count FROM  remit_history WHERE txn_status = $1 and txn_step = $2 and dsa_id = $3 %s %s %s`, ftQ, utQ, ttQ)
	var pfs storage.RemitHistory
	if err := s.db.Get(&pfs, getTransactionReport, pf.TxnStatus, pf.TxnStep, pf.DsaID); err != nil {
		if err == sql.ErrNoRows {
			return &pfs, nil
		}
		return nil, err
	}
	return &pfs, nil
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"brank.as/petnet/serviceutil/logging"
	"github.com/lib/pq"

	"brank.as/petnet/api/storage"
)

const createMicroInsuranceHistory = `
INSERT INTO micro_insurance_history (
	coy,
    dsa_id,
	location_id,
	user_code,
	trx_date,
	promo_amount,
	promo_code, 
    amount,
	coverage_count, 
	product_code, 
	processing_branch, 
	processed_by, 
	user_email,
	last_name, 
	first_name, 
	middle_name, 
	gender, 
	birthdate, 
	mobile_number,
	province_code, 
	city_code, 
	address, 
	marital_status, 
	occupation, 
	number_units,
	beneficiaries, 
	dependents,
    trx_status,
	trace_number,
    insurance_details,
	error_code, 
	error_message, 
	error_type,
	error_time,
	org_id
) VALUES (
	:coy,
    :dsa_id,
	:location_id,
	:user_code,
	:trx_date,
	:promo_amount,
	:promo_code, 
    :amount,
	:coverage_count, 
	:product_code, 
	:processing_branch, 
	:processed_by, 
	:user_email,
	:last_name, 
	:first_name, 
	:middle_name, 
	:gender, 
	:birthdate, 
	:mobile_number,
	:province_code, 
	:city_code, 
	:address, 
	:marital_status, 
	:occupation, 
	:number_units,
	:beneficiaries, 
	:dependents,
	:trx_status, 
	:trace_number,
    :insurance_details,
	:error_code, 
	:error_message, 
	:error_type, 
	:error_time,
	:org_id
)
RETURNING *;
`

// CreateMicroInsuranceHistory ...
func (s *Storage) CreateMicroInsuranceHistory(ctx context.Context,
	r storage.MicroInsuranceHistory,
) (*storage.MicroInsuranceHistory, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, createMicroInsuranceHistory)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err = stmt.Get(&r, r); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}

		return nil, fmt.Errorf("executing micro insurance history insert: %w", err)
	}

	return &r, nil
}

// GetMicroInsuranceHistoryByID ...
func (s *Storage) GetMicroInsuranceHistoryByID(ctx context.Context, id string) (*storage.MicroInsuranceHistory, error) {
	const getMIHistory = `SELECT * FROM micro_insurance_history WHERE id = $1`
	var r storage.MicroInsuranceHistory
	if err := s.db.Get(&r, getMIHistory, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrNotFound
		}

		return nil, fmt.Errorf("executing micro insurance history get: %w", err)
	}

	return &r, nil
}

func (s *Storage) UpdateMicroInsuranceHistoryStatusByTraceNumber(ctx context.Context, r storage.MicroInsuranceHistory) (*storage.MicroInsuranceHistory, error) {
	log := logging.FromContext(ctx)
	if !r.TraceNumber.Valid || r.TraceNumber.String == "" {
		log.Error("trace number is empty for update microinsurance history")
		return nil, storage.ErrInvalid
	}

	updateQ := `
	UPDATE micro_insurance_history
	SET
		trx_status = :trx_status,
		insurance_details = :insurance_details,
		error_code = :error_code,
		error_message = :error_message,
		error_time = :error_time,
		error_type = :error_type
	WHERE trace_number = :trace_number
	RETURNING *
	`

	stmt, err := s.db.PrepareNamedContext(ctx, updateQ)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err = stmt.Get(&r, r); err != nil {
		logging.WithError(err, log).Error("update microinsurance history")
		if err == sql.ErrNoRows {
			return nil, storage.ErrNotFound
		}

		return nil, fmt.Errorf("executing microinsurance history update: %w", err)
	}

	return &r, nil
}

// GetMicroInsuranceHistoryByTraceNumber ...
func (s *Storage) GetMicroInsuranceHistoryByTraceNumber(ctx context.Context, traceNumber string) (*storage.MicroInsuranceHistory, error) {
	const getMIHistory = `SELECT * FROM micro_insurance_history WHERE trace_number = $1`
	var r storage.MicroInsuranceHistory
	if err := s.db.Get(&r, getMIHistory, traceNumber); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrNotFound
		}

		return nil, fmt.Errorf("executing micro insurance history get: %w", err)
	}

	return &r, nil
}

func (s *Storage) ListMicroInsuranceHistory(ctx context.Context,
	f storage.MicroInsuranceFilter,
) ([]storage.MicroInsuranceHistory, error) {
	ft := ""
	ut := ""
	tr := ""
	tc := ":total_count_where:"
	b := NewBuilder(fmt.Sprintf(`WITH cnt AS (select count(*) as total FROM micro_insurance_history %s) 
			SELECT *, cnt.total FROM micro_insurance_history left join cnt on true`, tc))
	if !f.From.IsZero() {
		ft = f.From.Format("2006-01-02")
	}

	if !f.Until.IsZero() {
		ut = f.Until.Format("2006-01-02")
	}

	if !f.TrxDate.IsZero() {
		tr = f.TrxDate.Format("2006-01-02")
		ft = ""
		ut = ""
	}

	var scol string
	switch string(f.SortByColumn) {
	case string(storage.TotalAmount):
		scol = "amount"
	case string(storage.TranTime):
		scol = "trx_date"
	default:
		scol = ""
	}

	if string(f.SortOrder) == "" {
		f.SortOrder = storage.Desc
	}

	b.Where("trace_number", eq, f.TraceNumber).
		Where("dsa_id", eq, f.DsaID).
		Where("user_code", eq, f.UserCode).
		Where("trx_status", eq, f.TrxStatus).
		Where("trx_date", gtOrEq, ft, CompareDate()).
		Where("trx_date", ltOrEq, ut, CompareDate()).
		Where("trx_date", eq, tr, CompareDate()).
		Where("org_id", eq, f.OrgID).
		SortByColumn(scol, f.SortOrder).
		Limit(f.Limit).
		Offset(f.Offset).AddTotalQuery(tc)

	stmt, err := s.db.PrepareNamed(b.query)
	if err != nil {
		return nil, err
	}

	var res []storage.MicroInsuranceHistory
	if err = stmt.Select(&res, b.args); err != nil {
		return nil, fmt.Errorf("executing bill payment list: %w", err)
	}
	return res, nil
}

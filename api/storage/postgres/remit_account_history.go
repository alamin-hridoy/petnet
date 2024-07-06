package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"brank.as/petnet/api/storage"
	"brank.as/petnet/serviceutil/logging"
	"github.com/lib/pq"
)

const createRTAHistory = `
INSERT INTO remit_to_acc_history (
	org_id,
	partner,
	reference_number,
	trx_date,
	account_number,
	currency,
	service_charge,
	remarks,
	particulars,
	merchant_name,
	bank_id,
	location_id,
	user_id,
	currency_id,
	customer_id,
	form_type,
	form_number,
	trx_type,
	remote_location_id,
	remote_user_id,
	biller_name,
	trx_time,
	total_amount,
	account_name,
	beneficiary_address,
	beneficiary_birthdate,
	beneficiary_city,
	beneficiary_civil,
	beneficiary_country,
	beneficiary_customertype,
	beneficiary_firstname,
	beneficiary_lastname,
	beneficiary_middlename,
	beneficiary_tin,
	beneficiary_sex,
	beneficiary_state,
	currency_code_principal_amount,
	principal_amount,
	record_type,
	remitter_address,
	remitter_birthdate,
	remitter_city,
	remitter_civil,
	remitter_country,
	remitter_customer_type,
	remitter_firstname,
	remitter_gender,
	remitter_id,
	remitter_lastname,
	remitter_middlename,
	remitter_state,
	settlement_mode,
	notification,
	bene_zip_code,
	info,
	details,
	txn_status,
	error_code,
	error_message,
	error_time,
	error_type,
	updated_by,
	created_by
) VALUES (
	:org_id,
	:partner,
	:reference_number,
	:trx_date,
	:account_number,
	:currency,
	:service_charge,
	:remarks,
	:particulars,
	:merchant_name,
	:bank_id,
	:location_id,
	:user_id,
	:currency_id,
	:customer_id,
	:form_type,
	:form_number,
	:trx_type,
	:remote_location_id,
	:remote_user_id,
	:biller_name,
	:trx_time,
	:total_amount,
	:account_name,
	:beneficiary_address,
	:beneficiary_birthdate,
	:beneficiary_city,
	:beneficiary_civil,
	:beneficiary_country,
	:beneficiary_customertype,
	:beneficiary_firstname,
	:beneficiary_lastname,
	:beneficiary_middlename,
	:beneficiary_tin,
	:beneficiary_sex,
	:beneficiary_state,
	:currency_code_principal_amount,
	:principal_amount,
	:record_type,
	:remitter_address,
	:remitter_birthdate,
	:remitter_city,
	:remitter_civil,
	:remitter_country,
	:remitter_customer_type,
	:remitter_firstname,
	:remitter_gender,
	:remitter_id,
	:remitter_lastname,
	:remitter_middlename,
	:remitter_state,
	:settlement_mode,
	:notification,
	:bene_zip_code,
	:info,
	:details,
	:txn_status,
	:error_code,
	:error_message,
	:error_time,
	:error_type,
	:updated_by,
	:created_by
) RETURNING *`

func (s *Storage) CreateRTAHistory(ctx context.Context, req storage.RemitToAccountHistory) (*storage.RemitToAccountHistory, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, createRTAHistory)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&req, req); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, fmt.Errorf("executing remit to account insert: %w", err)
	}
	return &req, nil
}

func (s *Storage) UpdateRTAHistory(ctx context.Context, req storage.RemitToAccountHistory) (*storage.RemitToAccountHistory, error) {
	updateRTAHistory := `
	UPDATE remit_to_acc_history SET
		trx_type = :trx_type,
		reference_number = :reference_number,
		location_id = :location_id,
		form_number = :form_number,
		principal_amount = :principal_amount,
		total_amount = :total_amount,
		trx_date = :trx_date,
		txn_status = :txn_status,
		error_code = :error_code,
		error_message = :error_message,
		error_time = :error_time,
		error_type = :error_type
	%s
	RETURNING *
	`

	wqM := []string{}
	wqS := ""
	if req.OrgID != "" {
		wqM = append(wqM, "org_id = :org_id")
	}
	if req.ReferenceNumber != "" {
		wqM = append(wqM, "reference_number = :reference_number")
	}
	if req.LocationID != 0 {
		wqM = append(wqM, "location_id = :location_id")
	}
	if req.BankID != 0 {
		wqM = append(wqM, "bank_id = :bank_id")
	}
	if req.FormNumber != "" {
		wqM = append(wqM, "form_number = :form_number")
	}
	if req.Partner != "" {
		wqM = append(wqM, "partner = :partner")
	}
	if len(wqM) > 0 {
		wqS = fmt.Sprintf(" WHERE %s", strings.Join(wqM, " AND "))
	}
	updateRTAHistory = fmt.Sprintf(updateRTAHistory, wqS)
	log := logging.FromContext(ctx)

	stmt, err := s.db.PrepareNamedContext(ctx, updateRTAHistory)
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
		return nil, fmt.Errorf("executing RTA history update: %w", err)
	}
	return &req, nil
}

func (s *Storage) GetRTAHistory(ctx context.Context, RTAID string) (*storage.RemitToAccountHistory, error) {
	const getRTAHistory = `SELECT * FROM remit_to_acc_history WHERE id = $1`
	var r storage.RemitToAccountHistory
	if err := s.db.Get(&r, getRTAHistory, RTAID); err != nil {
		if err == sql.ErrNoRows {
			return &r, nil
		}
		return nil, fmt.Errorf("executing remit to account history get: %w", err)
	}
	return &r, nil
}

func (s *Storage) ListRTAHistory(ctx context.Context, f storage.RemitToAccountHistory) ([]storage.RemitToAccountHistory, error) {
	tc := ":total_count_where:"
	b := NewBuilder(fmt.Sprintf(`WITH cnt AS (select count(*) as total FROM remit_to_acc_history %s) SELECT *, cnt.total FROM remit_to_acc_history left join cnt on true`, tc))

	var scol string
	switch string(f.SortByColumn) {
	case string(storage.RTAHistoryIDCol):
		scol = "id"
	case string(storage.OrgIDRTACol):
		scol = "org_id"
	case string(storage.TrxTypeRTACol):
		scol = "trx_type"
	default:
		scol = ""
	}

	b.Where("id", eq, f.ID).
		Where("org_id", eq, f.OrgID).
		Where("trx_type", eq, f.TrxType).
		SortByColumn(scol, f.SortOrder).
		Limit(f.Limit).
		Offset(f.Offset).AddTotalQuery(tc)

	stmt, err := s.db.PrepareNamed(b.query)
	if err != nil {
		return nil, err
	}

	r := []storage.RemitToAccountHistory{}
	if err := stmt.Select(&r, b.args); err != nil {
		return nil, fmt.Errorf("executing RTA history list: %w", err)
	}
	return r, nil
}

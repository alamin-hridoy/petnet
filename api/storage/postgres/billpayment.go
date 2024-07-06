package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"brank.as/petnet/api/storage"
	"github.com/lib/pq"
)

const createBillPayment = `
INSERT INTO bill_payment (
	bill_id,
	biller_tag,
	location_id,
	user_id,
	sender_member_id,
	currency_id,
	account_number,
	amount,
	identifier,
	coy,
	service_charge,
	total_amount,
	bill_payment_status,
	error_code,
	error_message,
	error_type,
	bills,
	bill_payment_date,
	partner_id,
	biller_name,
	trx_date,
	remote_user_id,
	customer_id,
	remote_location_id,
	location_name,
	form_type,
	form_number,
	payment_method,
	other_info,
	client_reference_number,
	bill_partner_id,
	partner_charge,
	reference_number,
	validation_number,
	receipt_validation_number,
	tpa_id,
	type,
	txnid,
	org_id
) VALUES (
	:bill_id,
	:biller_tag,
	:location_id,
	:user_id,
	:sender_member_id,
	:currency_id,
	:account_number,
	:amount,
	:identifier,
	:coy,
	:service_charge,
	:total_amount,
	:bill_payment_status,
	:error_code,
	:error_message,
	:error_type,
	:bills,
	:bill_payment_date,
	:partner_id,
	:biller_name,
	:trx_date,
	:remote_user_id,
	:customer_id,
	:remote_location_id,
	:location_name,
	:form_type,
	:form_number,
	:payment_method,
	:other_info,
	:client_reference_number,
	:bill_partner_id,
	:partner_charge,
	:reference_number,
	:validation_number,
	:receipt_validation_number,
	:tpa_id,
	:type,
	:txnid,
	:org_id
) RETURNING
*
`

func (s *Storage) CreateBillPayment(ctx context.Context, r storage.BillPayment) (*storage.BillPayment, error) {
	if r.PartnerID == "" {
		return nil, storage.ErrInvalid
	}

	stmt, err := s.db.PrepareNamedContext(ctx, createBillPayment)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, fmt.Errorf("executing bill payment insert: %w", err)
	}
	return &r, nil
}

const updateBillPayment = `
UPDATE bill_payment
SET
	bill_payment_status = :bill_payment_status,
	bills = :bills,
	error_code = :error_code,
	error_message = :error_message,
	error_type = :error_type 
	WHERE bill_payment_id = :bill_payment_id
RETURNING bill_payment_id
`

func (s *Storage) UpdateBillPayment(ctx context.Context, r storage.BillPayment) (*storage.BillPayment, error) {
	if r.BillPaymentID == "" {
		return nil, fmt.Errorf("bill Payment ID cannot be empty")
	}
	stmt, err := s.db.PrepareNamedContext(ctx, updateBillPayment)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		return nil, fmt.Errorf("executing Bill Payment update: %w", err)
	}
	return &r, nil
}

func (s *Storage) GetBillPayment(ctx context.Context, billPaymentID string) (*storage.BillPayment, error) {
	const getBillPayment = `SELECT * FROM bill_payment WHERE bill_payment_id = $1`
	var r storage.BillPayment
	if err := s.db.Get(&r, getBillPayment, billPaymentID); err != nil {
		if err == sql.ErrNoRows {
			return &r, nil
		}
		return nil, fmt.Errorf("executing bill payment get: %w", err)
	}
	return &r, nil
}

func (s *Storage) ListBillPayment(ctx context.Context, f storage.BillPaymentFilter) ([]storage.BillPayment, error) {
	ft := ""
	ut := ""
	tr := ""
	tc := ":total_count_where:"
	b := NewBuilder(fmt.Sprintf(`WITH cnt AS (select count(*) as total FROM bill_payment %s) SELECT *, cnt.total FROM bill_payment left join cnt on true`, tc))
	if !f.From.IsZero() {
		ft = f.From.Format("2006-01-02")
	}
	if !f.Until.IsZero() {
		ut = f.Until.Format("2006-01-02")
	}
	if !f.TrxDate.IsZero() {
		tr = f.TrxDate.Format("2006-01-02")
	}

	var scol string
	var scolAmount string
	switch string(f.SortByColumn) {
	case string(storage.BillPaymentIDCol):
		scol = "bill_payment_id"
	case string(storage.BillIDCol):
		scol = "bill_id"
	case string(storage.DsaIDCol):
		scol = "dsa_id"
	case string(storage.SenderMemberIDCol):
		scol = "sender_member_id"
	case string(storage.UserCol):
		scol = "user_id"
	case string(storage.FeesCol):
		scolAmount = "service_charge"
	case string(storage.AmountCol):
		scolAmount = "total_amount"
	case string(storage.TransactionCompletedTime):
		scol = "trx_date"
	case string(storage.TransactionTimeCol):
		scol = "trx_date"
	default:
		scol = ""
	}

	b.Where("bill_payment_id", eq, f.BillPaymentID).
		Where("bill_id", eq, f.BillID).
		Where("bill_payment_status", eq, f.BillPaymentStatus).
		Where("sender_member_id", eq, f.SenderMemberID).
		Where("user_id", eq, f.UserID).
		Where("trx_date", gtOrEq, ft, CompareDate()).
		Where("trx_date", ltOrEq, ut, CompareDate()).
		Where("trx_date", eq, tr, CompareDate()).
		Where("org_id", eq, f.OrgID).
		Where("reference_number", eq, f.ReferenceNumber).
		WhereNotIn("partner_id", f.ExcludePartners).
		SortByColumn(scol, f.SortOrder).
		SortByColumnTextToJson(scolAmount, "amount", f.SortOrder).
		Limit(f.Limit).
		Offset(f.Offset).AddTotalQuery(tc)
	stmt, err := s.db.PrepareNamed(b.query)
	if err != nil {
		return nil, err
	}

	r := []storage.BillPayment{}
	if err := stmt.Select(&r, b.args); err != nil {
		return nil, fmt.Errorf("executing bill payment list: %w", err)
	}
	return r, nil
}

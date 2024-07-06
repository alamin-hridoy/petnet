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

const createCICOHistory = `
INSERT INTO cico_history (
	org_id,
	svc_provider,
	partner_code,
	trx_provider,
	trx_type,
	reference_number,
	petnet_trackingno,
	provider_trackingno,
	principal_amount,
	charges,
	total_amount,
	trx_date,
	details,
	txn_status,
	error_code,
	error_message,
	error_time,
	error_type,
	created_by,
	created
) VALUES (
	:org_id,
	:svc_provider,
	:partner_code,
	:trx_provider,
	:trx_type,
	:reference_number,
	:petnet_trackingno,
	:provider_trackingno,
	:principal_amount,
	:charges,
	:total_amount,
	:trx_date,
	:details,
	:txn_status,
	:error_code,
	:error_message,
	:error_time,
	:error_type,
	:created_by,
	:created
) RETURNING *`

func (s *Storage) CreateCICOHistory(ctx context.Context, req storage.CashInCashOutHistory) (*storage.CashInCashOutHistory, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, createCICOHistory)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&req, req); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, fmt.Errorf("executing cashin cashout insert: %w", err)
	}
	return &req, nil
}

func (s *Storage) UpdateCICOHistory(ctx context.Context, req storage.CashInCashOutHistory) (*storage.CashInCashOutHistory, error) {
	updateCICOHistory := `
	UPDATE cico_history
	SET
		partner_code = :partner_code,
		trx_provider = :trx_provider,
		trx_type = :trx_type,
		reference_number = :reference_number,
		petnet_trackingno = :petnet_trackingno,
		provider_trackingno = :provider_trackingno,
		principal_amount = :principal_amount,
		charges = :charges,
		total_amount = :total_amount,
		trx_date = :trx_date,
		txn_status = :txn_status,
		error_code = :error_code,
		error_message = :error_message,
		error_time = :error_time,
		error_type = :error_type
	%s
	RETURNING id
	`
	wqM := []string{}
	wqS := ""
	if req.OrgID != "" {
		wqM = append(wqM, "org_id = :org_id")
	}
	if req.PartnerCode != "" {
		wqM = append(wqM, "partner_code = :partner_code")
	}
	if req.ReferenceNumber != "" {
		wqM = append(wqM, "reference_number = :reference_number")
	}
	if req.SvcProvider != "" {
		wqM = append(wqM, "svc_provider = :svc_provider")
	}
	if req.Provider != "" {
		wqM = append(wqM, "trx_provider = :trx_provider")
	}
	if req.PetnetTrackingNo != "" {
		wqM = append(wqM, "petnet_trackingno = :petnet_trackingno")
	}
	if req.ProviderTrackingNo != "" {
		wqM = append(wqM, "provider_trackingno = :provider_trackingno")
	}
	if len(wqM) > 0 {
		wqS = fmt.Sprintf(" WHERE %s", strings.Join(wqM, " AND "))
	}
	updateCICOHistory = fmt.Sprintf(updateCICOHistory, wqS)
	log := logging.FromContext(ctx)
	stmt, err := s.db.PrepareNamedContext(ctx, updateCICOHistory)
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
		return nil, fmt.Errorf("executing CICO history update: %w", err)
	}
	return &req, nil
}

func (s *Storage) GetCICOHistory(ctx context.Context, CICOID string) (*storage.CashInCashOutHistory, error) {
	const getCICOHistory = `SELECT * FROM cico_history WHERE id = $1`
	var r storage.CashInCashOutHistory
	if err := s.db.Get(&r, getCICOHistory, CICOID); err != nil {
		if err == sql.ErrNoRows {
			return &r, nil
		}
		return nil, fmt.Errorf("executing CICO history get: %w", err)
	}
	return &r, nil
}

func (s *Storage) ListCICOHistory(ctx context.Context, f storage.CashInCashOutHistory) ([]storage.CashInCashOutHistory, error) {
	tc := ":total_count_where:"
	b := NewBuilder(fmt.Sprintf(`WITH cnt AS (select count(*) as total FROM cico_history %s) SELECT *, cnt.total FROM cico_history left join cnt on true`, tc))

	var scol string
	switch string(f.SortByColumn) {
	case string(storage.CICOHistoryIDCol):
		scol = "id"
	case string(storage.OrgIDCol):
		scol = "org_id"
	case string(storage.PartnerCodeCol):
		scol = "partner_code"
	case string(storage.ProviderCol):
		scol = "provider"
	case string(storage.TrxTypeCol):
		scol = "trx_type"
	default:
		scol = ""
	}

	b.Where("id", eq, f.ID).
		Where("org_id", eq, f.OrgID).
		Where("partner_code", eq, f.PartnerCode).
		Where("svc_provider", eq, f.SvcProvider).
		Where("trx_provider", eq, f.Provider).
		Where("trx_type", eq, f.TrxType).
		SortByColumn(scol, f.SortOrder).
		Limit(f.Limit).
		Offset(f.Offset).AddTotalQuery(tc)

	stmt, err := s.db.PrepareNamed(b.query)
	if err != nil {
		return nil, err
	}

	r := []storage.CashInCashOutHistory{}
	if err := stmt.Select(&r, b.args); err != nil {
		return nil, fmt.Errorf("executing CICO history list: %w", err)
	}
	return r, nil
}

func (s *Storage) ListCICOTrx(ctx context.Context, f storage.CashInCashOutTrxListFilter) ([]storage.CashInCashOutHistory, error) {
	tc := ":total_count_where:"
	b := NewBuilder(fmt.Sprintf(`WITH cnt AS (select count(*) as total FROM cico_history %s) SELECT *, cnt.total FROM cico_history left join cnt on true`, tc))
	ft := ""
	ut := ""
	if !f.From.IsZero() {
		ft = f.From.Format("2006-01-02")
	}
	if !f.Until.IsZero() {
		ut = f.Until.Format("2006-01-02")
	}

	var scol string
	var scolAmount string
	switch string(f.SortByColumn) {
	case string(storage.CICOHistoryIDCol):
		scol = "id"
	case string(storage.OrgIDCol):
		scol = "org_id"
	case string(storage.PartnerCodeCol):
		scol = "partner_code"
	case string(storage.ProviderCol):
		scol = "trx_provider"
	case string(storage.TrxTypeCol):
		scol = "trx_type"
	case string(storage.FeeCol):
		scol = "charges"
	case string(storage.TotalAmountCol):
		scol = "total_amount"
	case string(storage.TranTimeCol):
		scol = "trx_date"
	case string(storage.TransactionTimeCol):
		scol = "trx_date"
	default:
		scol = ""
	}

	b.Where("reference_number", eq, f.ReferenceNumber).
		Where("org_id", eq, f.OrgID).
		Where("partner_code", eq, f.PartnerCode).
		Where("svc_provider", eq, f.SvcProvider).
		Where("trx_provider", eq, f.Provider).
		Where("trx_type", eq, f.TrxType).
		Where("txn_status", eq, f.TxnStatus).
		Where("petnet_trackingno", eq, f.PetnetTrackingNo).
		Where("provider_trackingno", eq, f.ProviderTrackingNo).
		WhereNotIn("trx_provider", f.ExcludeProviders).
		Where("created", gtOrEq, ft, CompareDate()).
		Where("created", ltOrEq, ut, CompareDate()).
		SortByColumn(scol, f.SortOrder).
		SortByColumnTextToJson(scolAmount, "amount", f.SortOrder).
		Limit(f.Limit).
		Offset(f.Offset).AddTotalQuery(tc)

	stmt, err := s.db.PrepareNamed(b.query)
	if err != nil {
		return nil, err
	}

	r := []storage.CashInCashOutHistory{}
	if err := stmt.Select(&r, b.args); err != nil {
		return nil, fmt.Errorf("executing cico trx list: %w", err)
	}

	return r, nil
}

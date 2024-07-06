package postgres

import (
	"context"
	"fmt"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"github.com/lib/pq"
)

const insertRevenueSharingReport = `
INSERT INTO revenue_sharing_report (
	org_id,
	dsa_code,
	year_month,
	remittance_count,
	cico_count,
	bills_payment_count,
	insurance_count,
	dsa_commission,
	dsa_commission_type,
	status
) VALUES (
	:org_id,
	:dsa_code,
	:year_month,
	:remittance_count,
	:cico_count,
	:bills_payment_count,
	:insurance_count,
	:dsa_commission,
	:dsa_commission_type,
	:status
) RETURNING
    id
`

func (s *Storage) CreateRevenueSharingReport(ctx context.Context, req storage.RevenueSharingReport) (*storage.RevenueSharingReport, error) {
	switch {
	case req.OrgID == "":
		return nil, fmt.Errorf("org id cannot be empty")
	}
	log := logging.FromContext(ctx)
	stmt, err := s.prepareNamed(ctx, insertRevenueSharingReport)
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
		return nil, fmt.Errorf("executing revenue sharing report insert: %w", err)
	}
	return &req, nil
}

func (s *Storage) GetRevenueSharingReportList(ctx context.Context, f storage.RevenueSharingReport) ([]storage.RevenueSharingReport, error) {
	switch {
	case f.OrgID == "":
		return nil, fmt.Errorf("org id cannot be empty")
	}
	tc := ":total_count_where:"
	b := NewBuilder(fmt.Sprintf(`WITH cnt AS (select count(*) as total FROM revenue_sharing_report %s) SELECT *, cnt.total FROM revenue_sharing_report left join cnt on true`, tc))
	var scol string
	switch string(f.SortByColumn) {
	case string(storage.OrgIDCol):
		scol = "org_id"
	case string(storage.StatusCol):
		scol = "status"
	default:
		scol = ""
	}
	b = b.Where("org_id", eq, f.OrgID).
		SortByColumn(scol, string(f.SortOrder)).
		Limit(f.Limit).
		Offset(f.Offset).AddTotalQuery(tc)

	stmt, err := s.db.PrepareNamed(b.query)
	if err != nil {
		return nil, err
	}
	r := []storage.RevenueSharingReport{}
	if err := stmt.Select(&r, b.args); err != nil {
		return nil, fmt.Errorf("executing revenue sharing report list: %w", err)
	}
	return r, nil
}

const revenueSharingReportUpdate = `
UPDATE
	revenue_sharing_report
SET
	remittance_count = :remittance_count,
	cico_count = :cico_count,
	bills_payment_count = :bills_payment_count,
	insurance_count = :insurance_count,
	dsa_commission = :dsa_commission,
	dsa_commission_type = :dsa_commission_type
WHERE
	 org_id = :org_id AND year_month = :year_month
RETURNING
	id, org_id, year_month
`

func (s *Storage) UpdateRevenueSharingReport(ctx context.Context, req storage.RevenueSharingReport) (*storage.RevenueSharingReport, error) {
	switch {
	case req.OrgID == "":
		return nil, fmt.Errorf("org id cannot be empty")
	case req.YearMonth == "":
		return nil, fmt.Errorf("year month cannot be empty")
	}
	log := logging.FromContext(ctx)
	stmt, err := s.prepareNamed(ctx, revenueSharingReportUpdate)
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
		return nil, fmt.Errorf("executing revenue sharing report update: %w", err)
	}
	return &req, nil
}

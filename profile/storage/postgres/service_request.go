package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	spb "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/profile/storage"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const svcRequestCreate = `
INSERT INTO service_request (
	service_name,
	partner,
   org_id,
	company_name
) VALUES (
	:service_name,
	:partner,
	:org_id,
	:company_name
) RETURNING
updated, created
`

func (s *Storage) CreateSvcRequest(ctx context.Context, r storage.ServiceRequest) (*storage.ServiceRequest, error) {
	switch {
	case r.SvcName == "":
		return nil, fmt.Errorf("SvcName cannot be empty")
	case r.Partner == "":
		return nil, fmt.Errorf("Partner cannot be empty")
	case r.OrgID == "":
		return nil, fmt.Errorf("OrgID cannot be empty")
	}
	stmt, err := s.db.PrepareNamedContext(ctx, svcRequestCreate)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, fmt.Errorf("executing service request create insert: %w", err)
	}
	return &r, nil
}

func (s *Storage) ApplySvcRequest(ctx context.Context, r storage.ServiceRequest) error {
	switch {
	case r.OrgID == "":
		return fmt.Errorf("OrgID cannot be empty")
	case r.SvcName == "":
		return fmt.Errorf("SvcName cannot be empty")
	}

	const q = `
			UPDATE service_request
			SET
				status= 'PENDING',
				applied= now(),
				updated= now()
			WHERE org_id = :org_id AND service_name = :service_name AND status != 'ACCEPTED' AND status != 'PENDING'
			RETURNING updated
			`

	// todo create if statement for if status is approved then enabled should be turned on
	stmt, err := s.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		if err == sql.ErrNoRows {
			return storage.NotFound
		}
		return fmt.Errorf("executing service accept: %w", err)
	}
	return nil
}

func (s *Storage) AcceptSvcRequest(ctx context.Context, r storage.ServiceRequest) error {
	switch {
	case r.UpdatedBy == "":
		return fmt.Errorf("UpdatedBy cannot be empty")
	case r.OrgID == "":
		return fmt.Errorf("OrgID cannot be empty")
	case r.Partner == "":
		return fmt.Errorf("Partner cannot be empty")
	case r.SvcName == "":
		return fmt.Errorf("SvcName cannot be empty")
	}

	q := `
UPDATE service_request
SET
   remarks= COALESCE(NULLIF(:remarks, ''), remarks),
   status= 'ACCEPTED',
   updated_by= COALESCE(NULLIF(:updated_by, ''), updated_by),
	enabled=true,
	updated= now()
WHERE org_id = :org_id AND partner = :partner AND service_name = :service_name
RETURNING updated
`

	// todo create if statement for if status is approved then enabled should be turned on
	stmt, err := s.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		return fmt.Errorf("executing service accept: %w", err)
	}
	return nil
}

func (s *Storage) RejectSvcRequest(ctx context.Context, r storage.ServiceRequest) error {
	switch {
	case r.UpdatedBy == "":
		return fmt.Errorf("UpdatedBy cannot be empty")
	case r.OrgID == "":
		return fmt.Errorf("OrgID cannot be empty")
	case r.Partner == "":
		return fmt.Errorf("Partner cannot be empty")
	case r.SvcName == "":
		return fmt.Errorf("SvcName cannot be empty")
	}

	q := `
UPDATE service_request
SET
   remarks= COALESCE(NULLIF(:remarks, ''), remarks),
   status= 'REJECTED',
   updated_by= COALESCE(NULLIF(:updated_by, ''), updated_by),
	updated= now()
WHERE org_id = :org_id AND partner = :partner AND service_name = :service_name
RETURNING updated
`

	// todo create if statement for if status is approved then enabled should be turned on
	stmt, err := s.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		return fmt.Errorf("executing service reject: %w", err)
	}
	return nil
}

func (s *Storage) EnableSvcRequest(ctx context.Context, r storage.ServiceRequest) error {
	switch {
	case r.UpdatedBy == "":
		return fmt.Errorf("UpdatedBy cannot be empty")
	case r.OrgID == "":
		return fmt.Errorf("OrgID cannot be empty")
	case r.Partner == "":
		return fmt.Errorf("Partner cannot be empty")
	case r.SvcName == "":
		return fmt.Errorf("SvcName cannot be empty")
	}

	q := `
UPDATE service_request
SET
   enabled=true,
   updated_by= :updated_by,
	updated= now()
WHERE org_id = :org_id AND partner = :partner AND service_name = :service_name
RETURNING updated
`

	// todo create if statement for if status is approved then enabled should be turned on
	stmt, err := s.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		return fmt.Errorf("executing service reject: %w", err)
	}
	return nil
}

func (s *Storage) DisableSvcRequest(ctx context.Context, r storage.ServiceRequest) error {
	switch {
	case r.UpdatedBy == "":
		return fmt.Errorf("UpdatedBy cannot be empty")
	case r.OrgID == "":
		return fmt.Errorf("OrgID cannot be empty")
	case r.Partner == "":
		return fmt.Errorf("Partner cannot be empty")
	case r.SvcName == "":
		return fmt.Errorf("SvcName cannot be empty")
	}

	const q = `
UPDATE service_request
SET
   enabled=false,
   updated_by= :updated_by,
	updated= now()
WHERE org_id = :org_id AND partner = :partner AND service_name = :service_name
RETURNING updated
`

	// todo create if statement for if status is approved then enabled should be turned on
	stmt, err := s.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		return fmt.Errorf("executing service reject: %w", err)
	}
	return nil
}

func (s *Storage) ListSvcRequest(ctx context.Context, f storage.SvcRequestFilter) ([]storage.ServiceRequest, error) {
	stsQ, paramsVal := generateAndWhere(f)
	postQ := generateLimitOffset(f)
	qry := fmt.Sprintf("SELECT *, (SELECT count(distinct org_id) FROM service_request %s ) as total FROM service_request %s %s", stsQ, stsQ, postQ)
	stmt, err := s.db.PrepareNamed(qry)
	if err != nil {
		return nil, err
	}
	r := []storage.ServiceRequest{}
	if err := stmt.Select(&r, paramsVal); err != nil {
		return nil, fmt.Errorf("executing remittance list history: %w", err)
	}
	return r, nil
}

func (s *Storage) GetAllServiceRequest(ctx context.Context, f storage.SvcRequestFilter) ([]storage.ServiceRequest, error) {
	stsQ, paramsVal := generateAndWhere(f)
	postQ := generateLimitOffset(f)
	serviceReques := fmt.Sprintf(`SELECT *, (SELECT count(id) FROM service_request %s  
	) as total FROM service_request %s %s`, stsQ, stsQ, postQ)
	var pnr []storage.ServiceRequest
	stmt, err := s.db.PrepareNamed(serviceReques)
	if err != nil {
		return nil, fmt.Errorf("preparing named query Get srvice request List: %w", err)
	}
	if err := stmt.Select(&pnr, paramsVal); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return pnr, nil
}

// generateLimitOffset is return limit offset query
func generateLimitOffset(f storage.SvcRequestFilter) string {
	sortOrder := getShortOrder(f.SortOrder)
	scol := getShortColumnName(f.SortByColumn)
	lmt, ofst := "", ""
	ordrby := fmt.Sprintf(" ORDER BY %s %s", scol, sortOrder)
	if f.Limit > 0 {
		lmt = fmt.Sprintf(" LIMIT %d", f.Limit)
	}
	if f.Offset > 0 {
		ofst = fmt.Sprintf(" OFFSET %d", f.Offset)
	}
	mQ := fmt.Sprintf("%s%s%s", ordrby, lmt, ofst)
	return mQ
}

// getShortOrder is return exact order type
func getShortOrder(sortOrder string) string {
	sortOrder = strings.ToUpper(sortOrder)
	if sortOrder == "DESC" || sortOrder == "ASC" {
		return sortOrder
	}
	return "DESC"
}

// getShortColumnName is return exact db column name
func getShortColumnName(sortByColumn string) string {
	var scol string
	switch sortByColumn {
	case "CREATED":
		scol = "created"
	case "COMPANYNAME":
		scol = "company_name"
	case "SERVICENAME":
		scol = "service_name"
	case "STATUS":
		scol = "status"
	case "LASTUPDATED":
		scol = "updated"
	case "UPDATEDBY":
		scol = "updated_by"
	case "PARTNER":
		scol = "partner"
	case "APPLIED":
		scol = "applied"
	default:
		scol = "created"
	}
	return scol
}

// generateAndWhere is generate all the necessery and condition
func generateAndWhere(f storage.SvcRequestFilter) (string, map[string]interface{}) {
	stsQ, quryL, orQuryL := "", []string{}, []string{}
	paramsVal := make(map[string]interface{})
	if f.CompanyName != "" {
		quryL = append(quryL, " company_name ILIKE :company_name")
		paramsVal["company_name"] = fmt.Sprintf("%%%s%%", f.CompanyName)
	}
	if len(f.Status) > 0 {
		quryL = append(quryL, " status = ANY(:status)")
		paramsVal["status"] = pq.Array(f.Status)
	}
	if len(f.SvcName) > 0 {
		quryL = append(quryL, " service_name = ANY(:service_name)")
		paramsVal["service_name"] = pq.Array(f.SvcName)
	}
	if len(f.Partner) > 0 {
		quryL = append(quryL, " partner = ANY(:partners)")
		paramsVal["partners"] = pq.Array(f.Partner)
	}
	if len(f.OrgID) > 0 {
		quryL = append(quryL, " org_id = ANY(:org_id)")
		paramsVal["org_id"] = pq.Array(f.OrgID)
	}
	if len(orQuryL) > 0 {
		quryL = append(quryL, strings.Join(orQuryL, " AND "))
	}
	if len(quryL) > 0 {
		stsQ = fmt.Sprintf(" WHERE %s", strings.Join(quryL, " AND "))
	}
	return stsQ, paramsVal
}

// sliceFormator is used for format slice for postgress query
func sliceFormator(slce []string) string {
	stsL := ""
	if len(slce) > 0 {
		sts := []string{}
		for _, v := range slce {
			sts = append(sts, "'"+v+"'")
		}
		stsL = strings.Join(sts, ",")
	}
	return stsL
}

func (s *Storage) ValidateSvcRequest(ctx context.Context, f storage.ValidateSvcRequestFilter) (*storage.ValidateSvcResponse, error) {
	if f.OrgID == "" {
		return nil, fmt.Errorf("OrgID cannot be empty")
	}
	if !f.IsAnyPartnerEnabled && f.Partner == "" {
		return nil, fmt.Errorf("partner cannot be empty if IsAnyPartnerEnabled is false")
	}
	b := NewBuilder("SELECT * FROM service_request")
	b.Where("org_id", "=", f.OrgID).
		Where("status", "=", "ACCEPTED").
		IsWhere("enabled", "=", true)
	if !f.IsAnyPartnerEnabled {
		b.Where("partner", "=", f.Partner)
	}
	stmt, err := s.db.PrepareNamed(b.query)
	if err != nil {
		return nil, err
	}
	res := []storage.ServiceRequest{}
	ret := &storage.ValidateSvcResponse{
		Enabled: false,
	}
	if err := stmt.Select(&res, b.args); err != nil {
		return nil, fmt.Errorf("executing Validate Svc Request: %w", err)
	}
	if len(res) > 0 {
		ret.Enabled = true
	}
	return ret, nil
}

func (s *Storage) RemoveSvcRequest(ctx context.Context, f storage.ServiceRequest) error {
	switch {
	case f.OrgID == "":
		return status.Errorf(codes.InvalidArgument, "OrgID cannot be empty")
	case f.Partner == "":
		return status.Errorf(codes.InvalidArgument, "Partner cannot be empty")
	case f.SvcName == "":
		return status.Errorf(codes.InvalidArgument, "SvcName cannot be empty")
	}
	stsM := []string{"ACCEPTED", "PENDING"}
	stsL := sliceFormator(stsM)
	if _, err := s.db.Exec("DELETE FROM service_request WHERE org_id=$1 AND partner=$2 AND service_name=$3 AND status NOT IN($4)", f.OrgID, f.Partner, f.SvcName, stsL); err != nil {
		if err == sql.ErrNoRows {
			return storage.NotFound
		}
		return status.Errorf(codes.Internal, "Remove Svc Request failed")
	}
	return nil
}

func (s *Storage) SetStatusSvcRequest(ctx context.Context, r storage.ServiceRequest) error {
	switch {
	case r.OrgID == "":
		return fmt.Errorf("OrgID cannot be empty")
	case r.Partner == "":
		return fmt.Errorf("Partner cannot be empty")
	case r.SvcName == "":
		return fmt.Errorf("SvcName cannot be empty")
	}

	if r.Status == spb.ServiceRequestStatus_NOSTATUS.String() {
		r.Status = ""
	}

	q := `UPDATE service_request SET
			status= :status
			WHERE org_id = :org_id AND partner = :partner AND service_name = :service_name 
			RETURNING id
			`
	// todo create if statement for if status is approved then enabled should be turned on
	stmt, err := s.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("executing service set status: %w", err)
	}
	return nil
}

func (s *Storage) AddRemarkSvcRequest(ctx context.Context, r storage.ServiceRequest) error {
	switch {
	case r.OrgID == "":
		return fmt.Errorf("OrgID cannot be empty")
	case r.UpdatedBy == "":
		return fmt.Errorf("UpdatedBy cannot be empty")
	case r.SvcName == "":
		return fmt.Errorf("SvcName cannot be empty")
	}

	q := `
UPDATE service_request
SET
   remarks= COALESCE(NULLIF(:remarks, ''), remarks),
   updated_by= COALESCE(NULLIF(:updated_by, ''), updated_by),
	updated= now()
WHERE org_id = :org_id AND service_name = :service_name
RETURNING updated
`
	// todo create if statement for if status is approved then enabled should be turned on
	stmt, err := s.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		return fmt.Errorf("executing add remark service accept: %w", err)
	}
	return nil
}

const serviceRequestUpdateByOrgID = `
UPDATE
	service_request
SET
	org_id = :newOrgID,
	status = :status
WHERE
	org_id = :oldOrgID
RETURNING
    id
`

func (s *Storage) UpdateServiceRequestByOrgID(ctx context.Context, pf storage.UpdateServiceRequestOrgID) (string, error) {
	if pf.OldOrgID == "" {
		return "", fmt.Errorf("invalid org id")
	}

	stmt, err := s.db.PrepareNamedContext(ctx, serviceRequestUpdateByOrgID)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var id string
	if err := stmt.Get(&id, pf); err != nil {
		if err == sql.ErrNoRows {
			return "", storage.NotFound
		}
		return "", fmt.Errorf("executing service request update: %w", err)
	}

	return id, nil
}

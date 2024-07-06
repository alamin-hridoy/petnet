package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"brank.as/petnet/profile/storage"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Storage) CreateUploadSvcRequest(ctx context.Context, r storage.UploadServiceRequest) (*storage.UploadServiceRequest, error) {
	switch {
	case r.SvcName == "":
		return nil, fmt.Errorf("SvcName cannot be empty")
	case r.Partner == "":
		return nil, fmt.Errorf("Partner cannot be empty")
	case r.OrgID == "":
		return nil, fmt.Errorf("OrgID cannot be empty")
	case r.FileType == "":
		return nil, fmt.Errorf("file type cannot be empty")
	case r.FileID == "":
		return nil, fmt.Errorf("file id cannot be empty")
	}

	const uploadSvcRequestCreate = `
	INSERT INTO upload_service_request (
		org_id,
		partner,
		service_name,
		file_type,
		file_id,
		create_by
	) VALUES (
		:org_id,
		:partner,
		:service_name,
		:file_type,
		:file_id,
		:create_by
	) RETURNING *`

	stmt, err := s.db.PrepareNamedContext(ctx, uploadSvcRequestCreate)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, fmt.Errorf("executing upload service request insert: %w", err)
	}
	return &r, nil
}

func (s *Storage) UpdateUploadSvcRequest(ctx context.Context, r storage.UploadServiceRequest) (*storage.UploadServiceRequest, error) {
	switch {
	case r.VerifyBy == "":
		return nil, fmt.Errorf("VerifyBy cannot be empty")
	case r.OrgID == "":
		return nil, fmt.Errorf("OrgID cannot be empty")
	case r.Partner == "":
		return nil, fmt.Errorf("Partner cannot be empty")
	case r.SvcName == "":
		return nil, fmt.Errorf("SvcName cannot be empty")
	case r.FileType == "":
		return nil, fmt.Errorf("FileType cannot be empty")
	case r.FileID == "":
		return nil, fmt.Errorf("FileID cannot be empty")
	}

	q := `UPDATE upload_service_request SET 
	file_id= :file_id,
	verify_by= COALESCE(NULLIF(:verify_by, ''), verify_by),
	status= '',
	verified= now()
	WHERE org_id = :org_id AND partner = :partner AND service_name = :service_name AND file_type = :file_type
	RETURNING *`

	stmt, err := s.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		return nil, fmt.Errorf("executing upload service  request accept: %w", err)
	}
	return &r, nil
}

func (s *Storage) AcceptUploadSvcRequest(ctx context.Context, r storage.UploadServiceRequest) error {
	switch {
	case r.VerifyBy == "":
		return fmt.Errorf("VerifyBy cannot be empty")
	case r.OrgID == "":
		return fmt.Errorf("OrgID cannot be empty")
	case r.Partner == "":
		return fmt.Errorf("Partner cannot be empty")
	case r.SvcName == "":
		return fmt.Errorf("SvcName cannot be empty")
	case r.FileType == "":
		return fmt.Errorf("FileType cannot be empty")
	}

	q := `UPDATE upload_service_request SET 
	status= 'ACCEPTED',
	verify_by= COALESCE(NULLIF(:verify_by, ''), verify_by),
	verified= now()
	WHERE org_id = :org_id AND partner = :partner AND service_name = :service_name AND file_type = :file_type
	RETURNING *`

	stmt, err := s.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		return fmt.Errorf("executing upload service  request accept: %w", err)
	}
	return nil
}

func (s *Storage) RejectUploadSvcRequest(ctx context.Context, r storage.UploadServiceRequest) error {
	switch {
	case r.VerifyBy == "":
		return fmt.Errorf("VerifyBy cannot be empty")
	case r.OrgID == "":
		return fmt.Errorf("OrgID cannot be empty")
	case r.Partner == "":
		return fmt.Errorf("Partner cannot be empty")
	case r.SvcName == "":
		return fmt.Errorf("SvcName cannot be empty")
	case r.FileType == "":
		return fmt.Errorf("FileType cannot be empty")
	}

	q := `UPDATE upload_service_request SET 
	status= 'REJECTED',
	verify_by= COALESCE(NULLIF(:verify_by, ''), verify_by),
	verified= now()
	WHERE org_id = :org_id AND partner = :partner AND service_name = :service_name AND file_type = :file_type
	RETURNING *`

	stmt, err := s.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		return fmt.Errorf("executing upload service request reject: %w", err)
	}
	return nil
}

func (s *Storage) ListUploadSvcRequest(ctx context.Context, f storage.UploadSvcRequestFilter) ([]storage.UploadServiceRequest, error) {
	b := NewBuilder("SELECT *, count(*) OVER() AS total FROM upload_service_request")
	b.Where("org_id", "=", f.OrgID).
		Any("status", f.Status).
		Any("service_name", f.SvcName).
		Any("partner", f.Partner)
	stmt, err := s.db.PrepareNamed(b.query)
	if err != nil {
		return nil, err
	}
	r := []storage.UploadServiceRequest{}
	if err := stmt.Select(&r, b.args); err != nil {
		return nil, fmt.Errorf("executing upload request service list history: %w", err)
	}
	return r, nil
}

func (s *Storage) RemoveUploadSvcRequest(ctx context.Context, f storage.UploadServiceRequest) error {
	switch {
	case f.OrgID == "":
		return status.Errorf(codes.InvalidArgument, "OrgID cannot be empty")
	case f.Partner == "":
		return status.Errorf(codes.InvalidArgument, "Partner cannot be empty")
	case f.SvcName == "":
		return status.Errorf(codes.InvalidArgument, "SvcName cannot be empty")
	case f.FileType == "":
		return status.Errorf(codes.InvalidArgument, "FileType cannot be empty")
	}
	if _, err := s.db.Exec("DELETE FROM upload_service_request WHERE org_id=$1 AND partner=$2 AND service_name=$3 AND file_type=$4", f.OrgID, f.Partner, f.SvcName, f.FileType); err != nil {
		if err == sql.ErrNoRows {
			return storage.NotFound
		}
		return status.Errorf(codes.Internal, "Remove Upload Svc Request failed")
	}
	return nil
}

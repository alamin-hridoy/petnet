package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"brank.as/petnet/profile/storage"
	"github.com/lib/pq"
)

func (s *Storage) CreateFileUpload(ctx context.Context, f storage.FileUpload) (*storage.FileUpload, error) {
	const createFile = `
INSERT INTO file_upload (
    org_id,
    user_id,
    upload_type,
	file_names,
    bucket_name,
	submitted,
	checked,
	file_name
) VALUES (
    :org_id,
    :user_id,
    :upload_type,
	:file_names,
	:bucket_name,
	:submitted,
	:checked,
	:file_name
)
RETURNING file_id, created`
	if f.FileNames != "" {
		f.Submitted = 1
	}
	stmt, err := s.db.PrepareNamedContext(ctx, createFile)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if err := stmt.Get(&f, f); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, err
	}
	return &f, nil
}

func (s *Storage) UpsertFileUpload(ctx context.Context, f storage.FileUpload) (*storage.FileUpload, error) {
	const createFile = `
UPDATE
	file_upload
SET
	file_names= :file_names,
	submitted= COALESCE(NULLIF(:submitted, 0), submitted),
	checked= COALESCE(:checked, checked),
	file_name= COALESCE(:file_name, file_name)
WHERE
	org_id=:org_id AND upload_type=:upload_type 
RETURNING file_id, created`
	if f.FileNames != "" {
		f.Submitted = 1
	}
	stmt, err := s.db.PrepareNamedContext(ctx, createFile)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	if err := stmt.Get(&f, f); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return &f, nil
}

func (s *Storage) GetFileUpload(ctx context.Context, id string) (*storage.FileUpload, error) {
	const getFile = `SELECT * FROM file_upload WHERE file_id = $1`
	var f storage.FileUpload
	if err := s.db.Get(&f, getFile, id); err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *Storage) ListFileUploads(ctx context.Context, org string, f storage.FileUploadFilter) ([]storage.FileUpload, error) {
	const getFiles = `SELECT * FROM file_upload WHERE org_id = $1 AND COALESCE(upload_type = ANY($2), TRUE)`
	fu := []storage.FileUpload{}
	if err := s.db.Select(&fu, getFiles, org, pq.StringArray(f.UploadTypes)); err != nil {
		return nil, err
	}
	return fu, nil
}

func (s *Storage) DeleteFileUpload(ctx context.Context, id string, name string) error {
	uploadedFIles, err := s.GetFileUpload(ctx, id)
	if err != nil {
		return err
	}
	if uploadedFIles.FileNames == "" {
		return errors.New("No File Found")
	}
	fileLists := strings.Split(uploadedFIles.FileNames, ",")
	newFileLists := []string{}
	for _, f := range fileLists {
		if f != name {
			newFileLists = append(newFileLists, f)
		}
	}
	uploadedFIles.FileNames = strings.Join(newFileLists, ",")
	if _, err := s.UpsertFileUpload(ctx, *uploadedFIles); err != nil {
		return fmt.Errorf("executing file_upload delete: %w", err)
	}
	return nil
}

package file

import (
	"context"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/profile/storage/postgres"
)

type FileStore interface {
	CreateFileUpload(ctx context.Context, f storage.FileUpload) (*storage.FileUpload, error)
	UpsertFileUpload(ctx context.Context, f storage.FileUpload) (*storage.FileUpload, error)
	ListFileUploads(ctx context.Context, org string, f storage.FileUploadFilter) ([]storage.FileUpload, error)
	DeleteFileUpload(ctx context.Context, id string, name string) error
}

type Svc struct {
	st FileStore
}

func New(st *postgres.Storage) *Svc {
	return &Svc{st: st}
}

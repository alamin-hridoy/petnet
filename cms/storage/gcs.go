package storage

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

// GCS is the common interface for onboarding provider.
type GCS interface {
	// Store stores the given file and returns a public accessible URL to file.
	Store(context.Context, *File) (string, error)
	// Get reads the given file and returns it.
	Get(context.Context, string) (*storage.ObjectHandle, error)
	// True if mocked else false.
	IsMock() bool
}

// File is a specific file of an org.
type File struct {
	Content []byte
	OrgID   string
}

// NewFile creates a new instance of file.
func NewFile(content []byte, orgID string) *File {
	return &File{Content: content, OrgID: orgID}
}

// GoogleCloudStorage is an implementation of storage.GCS that uses GoogleCloudStorage.
type GoogleCloudStorage struct {
	bucketName string
	bucket     *storage.BucketHandle
}

// NewGoogleCloudStorage creates new GoogleCloudStorage instance.
func NewGoogleCloudStorage(client *storage.Client, bucketName string) *GoogleCloudStorage {
	return &GoogleCloudStorage{
		bucketName: bucketName,
		bucket:     client.Bucket(bucketName),
	}
}

func (g *GoogleCloudStorage) writeFile(ctx context.Context, l *File) (*storage.ObjectHandle, error) {
	fileName := uuid.New().String()
	obj := g.bucket.Object(fileName)
	w := obj.NewWriter(ctx)
	defer func() { _ = w.Close() }()

	_, err := w.Write(l.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to write content, err: %w", err)
	}

	if err = w.Close(); err != nil {
		return nil, fmt.Errorf("failed to close, err: %w", err)
	}
	return obj, nil
}

// Store stores the given file to GCS.
func (g *GoogleCloudStorage) Store(ctx context.Context, l *File) (string, error) {
	obj, err := g.writeFile(ctx, l)
	if err != nil {
		return "", fmt.Errorf("failed to write file, err: %w", err)
	}
	return obj.ObjectName(), nil
}

// Get gets the given file from gcs
func (g *GoogleCloudStorage) Get(ctx context.Context, fn string) (*storage.ObjectHandle, error) {
	if err := validation.Validate(&fn, validation.Required, is.UUID); err != nil {
		return nil, err
	}
	obj := g.bucket.Object(fn)
	return obj, nil
}

// IsMock checks if storage is mocked
func (g *GoogleCloudStorage) IsMock() bool {
	return false
}

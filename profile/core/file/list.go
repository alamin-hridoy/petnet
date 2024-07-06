package file

import (
	"context"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
)

// ListFiles ...
func (s *Svc) ListFiles(ctx context.Context, oid string, f storage.FileUploadFilter) ([]storage.FileUpload, error) {
	log := logging.FromContext(ctx)

	res, err := s.st.ListFileUploads(ctx, oid, f)
	if err != nil {
		logging.WithError(err, log).Error("listing file uploads")
		return nil, err
	}
	return res, nil
}

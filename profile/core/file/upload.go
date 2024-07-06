package file

import (
	"context"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
)

// UpsertFileUpload ...
func (s *Svc) UploadFile(ctx context.Context, fu storage.FileUpload) (*storage.FileUpload, error) {
	log := logging.FromContext(ctx)

	res, err := s.st.CreateFileUpload(ctx, fu)
	if err != nil {
		if err != storage.Conflict {
			logging.WithError(err, log).Error("creating file upload")
			return nil, err
		}
		res, err = s.st.UpsertFileUpload(ctx, fu)
		if err != nil {
			logging.WithError(err, log).Error("upserting file upload")
			return nil, err
		}
	}
	return res, nil
}

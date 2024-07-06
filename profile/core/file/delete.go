package file

import (
	"context"

	fpb "brank.as/petnet/gunk/dsa/v2/file"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteFileUpload ...
func (s *Svc) DeleteFileUpload(ctx context.Context, fu *fpb.DeleteFileUploadRequest) error {
	log := logging.FromContext(ctx)
	if fu.ID == "" {
		return status.Error(codes.Internal, "File ID is missing")
	}
	if fu.FileNames == "" {
		return status.Error(codes.Internal, "File Names is missing")
	}
	err := s.st.DeleteFileUpload(ctx, fu.ID, fu.FileNames)
	if err != nil {
		logging.WithError(err, log).Error("Delete uploaded file")
		return err
	}
	return nil
}

package file

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fpb "brank.as/petnet/gunk/dsa/v2/file"
	"brank.as/petnet/profile/storage"
)

func (s *Svc) DeleteFileUpload(ctx context.Context, fu *fpb.DeleteFileUploadRequest) (*fpb.DeleteFileUploadResponse, error) {
	if fu.ID == "" {
		return nil, status.Error(codes.Internal, "File ID is missing")
	}
	if fu.FileNames == "" {
		return nil, status.Error(codes.Internal, "File Names is missing")
	}
	if fu.OrgID == "" {
		return nil, status.Error(codes.Internal, "Org Id is missing")
	}
	files, err := s.core.ListFiles(ctx, fu.GetOrgID(), storage.FileUploadFilter{})
	if err != nil {
		return nil, status.Error(codes.Internal, "You don't have permission to delete this file")
	}
	var fileExist bool
	for _, fl := range files {
		fileLists := strings.Split(fl.FileNames, ",")
		for _, fn := range fileLists {
			if fn == fu.FileNames {
				fileExist = true
			}
		}
	}
	if !fileExist {
		return nil, status.Error(codes.Internal, "You don't have permission to delete this file")
	}
	if err := s.core.DeleteFileUpload(ctx, fu); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete File")
	}
	return &fpb.DeleteFileUploadResponse{}, nil
}

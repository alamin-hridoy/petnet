package file

import (
	"context"
	"encoding/json"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/profile/storage"

	fpb "brank.as/petnet/gunk/dsa/v2/file"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) ListFiles(ctx context.Context, req *fpb.ListFilesRequest) (*fpb.ListFilesResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, validation.Required, is.UUID),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	f := storage.FileUploadFilter{}
	for _, t := range req.GetTypes() {
		f.UploadTypes = append(f.UploadTypes, t.String())
	}
	fs, err := s.core.ListFiles(ctx, req.OrgID, f)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list file uploads")
	}

	return &fpb.ListFilesResponse{
		FileUploads: storageToProto(fs),
	}, nil
}

func storageToProto(fs []storage.FileUpload) []*fpb.FileUpload {
	pfs := []*fpb.FileUpload{}
	for _, f := range fs {
		fName := map[string]string{}
		json.Unmarshal([]byte(f.FileName), &fName)
		pf := &fpb.FileUpload{
			ID:         f.FileID,
			OrgID:      f.OrgID,
			UserID:     f.UserID,
			Type:       fpb.UploadType(fpb.UploadType_value[f.UploadType]),
			FileNames:  split(f.FileNames, ","),
			BucketName: f.BucketName,
			Submitted:  ppb.Boolean(f.Submitted),
			FileName:   fName,
		}
		if f.Checked.Valid {
			pf.DateChecked = tspb.New(f.Checked.Time)
		}
		pfs = append(pfs, pf)
	}
	return pfs
}

func split(s string, sep string) []string {
	ss := strings.Split(s, sep)
	if ss[0] == "" {
		return nil
	}
	return ss
}

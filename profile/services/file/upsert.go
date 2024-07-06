package file

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/profile/storage"

	fpb "brank.as/petnet/gunk/dsa/v2/file"
)

func (s *Svc) UpsertFiles(ctx context.Context, req *fpb.UpsertFilesRequest) (*fpb.UpsertFilesResponse, error) {
	for _, f := range req.FileUploads {
		if err := validation.ValidateStruct(f,
			validation.Field(&f.OrgID, validation.Required, is.UUID),
			// validation.Field(&f.UserID, validation.Required, is.UUID),
			// validation.Field(&f.FileNames, validation.Required, validation.Each(is.UUID)),
			validation.Field(&f.Type, validation.Required),
		); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	sfs := filesToStorage(req.FileUploads)
	var fu []*fpb.FileUpload
	for _, sf := range sfs {
		ur, err := s.core.UploadFile(ctx, sf)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to create file upload")
		}
		fName := map[string]string{}
		json.Unmarshal([]byte(ur.FileName), &fName)
		fu = append(fu, &fpb.FileUpload{
			ID:        ur.FileID,
			Type:      fpb.UploadType(fpb.UploadType_value[sf.UploadType]),
			Submitted: 1,
			FileName:  fName,
		})
	}
	return &fpb.UpsertFilesResponse{
		FileUploads: fu,
	}, nil
}

func filesToStorage(fs []*fpb.FileUpload) []storage.FileUpload {
	sfs := []storage.FileUpload{}
	for _, f := range fs {
		fsss, _ := json.Marshal(f.FileName)
		nt := sql.NullTime{}
		if ts := f.GetDateChecked(); ts.IsValid() {
			nt.Time = ts.AsTime()
			nt.Valid = true
		}
		sf := storage.FileUpload{
			FileID:     f.GetID(),
			OrgID:      f.GetOrgID(),
			UserID:     f.GetUserID(),
			UploadType: f.GetType().String(),
			FileNames:  strings.Join(f.GetFileNames(), ","),
			BucketName: f.GetBucketName(),
			Submitted:  int(f.GetSubmitted()),
			Checked:    nt,
			FileName:   string(fsss),
		}
		sfs = append(sfs, sf)
	}
	return sfs
}

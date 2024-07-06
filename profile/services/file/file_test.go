package file

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus/hooks/test"
	"google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/profile/storage/postgres"

	fpb "brank.as/petnet/gunk/dsa/v2/file"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	fc "brank.as/petnet/profile/core/file"
)

func TestFiles(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	sd := map[string]string{
		"1test":  "1test",
		"1test1": "1test1",
		"test":   "test",
		"test1":  "test1",
	}
	oid := uuid.New().String()
	req := &fpb.UpsertFilesRequest{
		FileUploads: []*fpb.FileUpload{
			{
				OrgID:       oid,
				UserID:      uuid.New().String(),
				FileNames:   []string{uuid.New().String(), uuid.New().String()},
				BucketName:  "bucket1",
				Type:        fpb.UploadType_Picture,
				Submitted:   ppb.Boolean_True,
				DateChecked: timestamppb.Now(),
				FileName: map[string]string{
					"test":  "test",
					"test1": "test1",
				},
			},
			{
				OrgID:       oid,
				UserID:      uuid.New().String(),
				FileNames:   []string{uuid.New().String(), uuid.New().String()},
				BucketName:  "bucket2",
				Type:        fpb.UploadType_NBIClearance,
				Submitted:   ppb.Boolean_True,
				DateChecked: timestamppb.Now(),
				FileName: map[string]string{
					"1test":  "1test",
					"1test1": "1test1",
				},
			},
		},
	}
	test.NewNullLogger()
	st, cleanup := postgres.NewTestStorage(os.Getenv("DATABASE_CONNECTION"), filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(cleanup)

	_, err := st.CreateOrgProfile(ctx, &storage.OrgProfile{
		OrgID:  oid,
		UserID: uuid.NewString(),
	})
	if err != nil {
		t.Fatal(err)
	}

	s := New(fc.New(st))
	if _, err := s.UpsertFiles(ctx, req); err != nil {
		t.Fatal("CreateFiles: ", err)
	}

	req.FileUploads[0].FileNames = []string{uuid.New().String(), uuid.New().String()}
	req.FileUploads[1].FileNames = []string{uuid.New().String(), uuid.New().String()}
	req.FileUploads[0].FileName = sd
	req.FileUploads[1].FileName = sd
	if _, err = s.UpsertFiles(ctx, req); err != nil {
		t.Fatal("UpdateFiles: ", err)
	}

	o := cmp.Options{
		cmpopts.IgnoreFields(fpb.FileUpload{}, "ID"),
		cmpopts.IgnoreFields(timestamppb.Timestamp{}, "Nanos", "Seconds"),
		cmpopts.IgnoreUnexported(
			fpb.FileUpload{}, timestamppb.Timestamp{},
		),
	}
	res, err := s.ListFiles(ctx, &fpb.ListFilesRequest{
		OrgID: oid,
	})
	if err != nil {
		t.Error("ListFiles: ", err)
	}
	if !cmp.Equal(req.FileUploads, res.FileUploads, o) {
		t.Error("ListFiles (-want +got): ", cmp.Diff(req.FileUploads, res.FileUploads, o))
	}

	res, err = s.ListFiles(ctx, &fpb.ListFilesRequest{
		OrgID: oid,
		Types: []fpb.UploadType{fpb.UploadType_Picture},
	})
	if err != nil {
		t.Error("ListFiles: ", err)
	}
	if !cmp.Equal([]*fpb.FileUpload{req.FileUploads[0]}, res.FileUploads, o) {
		t.Error("ListFiles (-want +got): ", cmp.Diff(req.FileUploads, res.FileUploads, o))
	}
}

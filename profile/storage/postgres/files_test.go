package postgres_test

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"brank.as/petnet/profile/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestFileUpload(t *testing.T) {
	ts := newTestStorage(t)

	ctx := context.Background()
	testOrg := uuid.NewString()
	_, err := ts.CreateOrgProfile(ctx, &storage.OrgProfile{
		OrgID:  testOrg,
		UserID: uuid.NewString(),
	})
	if err != nil {
		t.Fatal(err)
	}

	upl := []storage.FileUpload{
		{
			OrgID:      testOrg,
			UserID:     uuid.NewString(),
			UploadType: "IDPhoto",
			BucketName: "bucket",
			FileNames:  "file1,file2,file3",
			Submitted:  1,
			FileName:   "test",
		},
		{
			OrgID:      testOrg,
			UserID:     uuid.NewString(),
			UploadType: "NBIClearance",
			BucketName: "bucket",
			FileNames:  "file4,file5,file6",
			Submitted:  1,
			FileName:   "test2",
		},
	}
	opt := []cmp.Option{
		cmp.FilterPath(func(p cmp.Path) bool {
			return p.Last().String() == ".FileID"
		}, cmp.Comparer(func(a, b string) bool {
			c := a
			if a == "" {
				c = b
			}
			_, err := uuid.Parse(c)
			return err == nil
		})),
		cmpopts.IgnoreFields(storage.FileUpload{}, "Created"),
	}
	for _, u := range upl {
		got, err := ts.CreateFileUpload(ctx, u)
		if err != nil {
			t.Error(err)
		}
		if !cmp.Equal(u, *got, opt...) {
			t.Error(cmp.Diff(u, *got, opt...))
		}
		got2, err := ts.GetFileUpload(ctx, got.FileID)
		if err != nil {
			t.Error(err)
		}
		if !cmp.Equal(got, got2) {
			t.Error(cmp.Diff(got, got2))
		}

		u.FileNames = "updated"
		u.Submitted = 1
		u.Checked = sql.NullTime{Time: time.Unix(1515151515, 0), Valid: true}
		got3, err := ts.UpsertFileUpload(ctx, u)
		if err != nil {
			t.Error(err)
		}
		got4, err := ts.GetFileUpload(ctx, got.FileID)
		if err != nil {
			t.Error(err)
		}
		if !cmp.Equal(got3, got4) {
			t.Error(cmp.Diff(got3, got4))
		}
	}
	upl[0].FileNames = "updated"
	upl[0].Submitted = 1
	upl[0].Checked = sql.NullTime{Time: time.Unix(1515151515, 0), Valid: true}
	upl[1].FileNames = "updated"
	upl[1].Submitted = 1
	upl[1].Checked = sql.NullTime{Time: time.Unix(1515151515, 0), Valid: true}
	list, err := ts.ListFileUploads(ctx, testOrg, storage.FileUploadFilter{
		UploadTypes: []string{upl[0].UploadType},
	})
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal([]storage.FileUpload{upl[0]}, list, opt...) {
		t.Error(cmp.Diff([]storage.FileUpload{upl[0]}, list, opt...))
	}
	list2, err := ts.ListFileUploads(ctx, testOrg, storage.FileUploadFilter{})
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(upl, list2, opt...) {
		t.Error(cmp.Diff(upl, list2, opt...))
	}
	deletedFileName := "updated"
	deletedFileID := ""
	for _, f := range list2 {
		deletedFileID = f.FileID
		err := ts.DeleteFileUpload(ctx, f.FileID, deletedFileName)
		if err != nil {
			t.Error(err)
		}
	}
	list3, err := ts.ListFileUploads(ctx, testOrg, storage.FileUploadFilter{
		UploadTypes: []string{upl[0].UploadType},
	})
	if err != nil {
		t.Error(err)
	}
	for _, f := range list3 {
		if f.FileID == deletedFileID {
			if f.FileNames != "" {
				fileLists := strings.Split(f.FileNames, ",")
				for _, fn := range fileLists {
					if fn == deletedFileName {
						t.Error(errors.New("Failed to deleted file."))
					}
				}
			}
		}
	}
}

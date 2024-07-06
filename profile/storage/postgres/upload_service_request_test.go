package postgres_test

import (
	"context"
	"testing"

	"brank.as/petnet/profile/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestUploadServiceRequest(t *testing.T) {
	st := newTestStorage(t)

	o := cmp.Options{
		cmpopts.IgnoreFields(storage.UploadServiceRequest{}, "Verified", "Created", "ID"),
	}

	oid := uuid.NewString()
	uid := uuid.NewString()
	want := storage.UploadServiceRequest{
		OrgID:    oid,
		Partner:  "ALL",
		SvcName:  "svcName",
		Status:   "",
		FileType: "membership",
		FileID:   "123456",
		Total:    "1",
		CreateBy: uid,
	}
	// create upload request
	csres, err := st.CreateUploadSvcRequest(context.TODO(), want)
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if csres.Created.IsZero() {
		t.Fatal("returned empty Created time")
	}
	// list upload request
	got1, err := st.ListUploadSvcRequest(context.TODO(), storage.UploadSvcRequestFilter{})
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if len(got1) != 1 {
		t.Fatalf("want 1 service request, got %v", len(got1))
	}
	if !cmp.Equal(want, got1[0], o) {
		t.Fatal(cmp.Diff(want, got1[0], o))
	}
	// updated upload request
	want.FileID = "98765"
	want.VerifyBy = uid
	usres, err := st.UpdateUploadSvcRequest(context.TODO(), want)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(&want, usres, o) {
		t.Fatal(cmp.Diff(want, usres, o))
	}
	// accepted upload request
	err = st.AcceptUploadSvcRequest(context.TODO(), want)
	if err != nil {
		t.Fatal(err)
	}
	want.Status = "ACCEPTED"
	got2, err := st.ListUploadSvcRequest(context.TODO(), storage.UploadSvcRequestFilter{})
	if err != nil {
		t.Fatalf("got 2 error %v, want nil", err)
	}
	if len(got2) != 1 {
		t.Fatalf("want 2 upload service request, got %v", len(got2))
	}
	if !cmp.Equal(want, got2[0], o) {
		t.Fatal(cmp.Diff(want, got2[0], o))
	}
	// rejected upload request
	err = st.RejectUploadSvcRequest(context.TODO(), want)
	if err != nil {
		t.Fatal(err)
	}
	want.Status = "REJECTED"
	got3, err := st.ListUploadSvcRequest(context.TODO(), storage.UploadSvcRequestFilter{})
	if err != nil {
		t.Fatalf("got 3 error %v, want nil", err)
	}
	if len(got3) != 1 {
		t.Fatalf("want 3 upload service request, got %v", len(got3))
	}
	if !cmp.Equal(want, got3[0], o) {
		t.Fatal(cmp.Diff(want, got3[0], o))
	}
	// removed upload request
	err = st.RemoveUploadSvcRequest(context.TODO(), want)
	if err != nil {
		t.Fatal(err)
	}
	got4, err := st.ListUploadSvcRequest(context.TODO(), storage.UploadSvcRequestFilter{})
	if err != nil {
		t.Fatalf("got 4 error %v, want nil", err)
	}
	if len(got4) != 0 {
		t.Fatalf("want 4 upload service request, got %v", len(got4))
	}
}

func TestUploadListServiceRequest(t *testing.T) {
	st := newTestStorage(t)
	oid1 := uuid.NewString()
	oid2 := uuid.NewString()
	oid3 := uuid.NewString()
	uid := uuid.NewString()
	p1 := storage.UploadServiceRequest{
		OrgID:    oid1,
		Partner:  "P1",
		SvcName:  "svcName",
		Status:   "ACCEPTED",
		FileType: "membership",
		FileID:   "123456",
		CreateBy: uid,
		VerifyBy: uid,
		Total:    "1",
	}
	p2 := storage.UploadServiceRequest{
		OrgID:    oid1,
		Partner:  "P2",
		SvcName:  "svcName1",
		Status:   "ACCEPTED",
		FileType: "membership",
		FileID:   "123456",
		Total:    "1",
		CreateBy: uid,
		VerifyBy: uid,
	}
	p3 := storage.UploadServiceRequest{
		OrgID:    oid1,
		Partner:  "P3",
		SvcName:  "svcName2",
		Status:   "ACCEPTED",
		FileType: "membership",
		FileID:   "123456",
		Total:    "1",
		CreateBy: uid,
		VerifyBy: uid,
	}
	p4 := storage.UploadServiceRequest{
		OrgID:    oid2,
		Partner:  "P4",
		SvcName:  "svcName3",
		Status:   "ACCEPTED",
		FileType: "membership",
		FileID:   "123456",
		Total:    "1",
		CreateBy: uid,
		VerifyBy: uid,
	}
	p5 := storage.UploadServiceRequest{
		OrgID:    oid3,
		Partner:  "P5",
		SvcName:  "svcName4",
		Status:   "REJECTED",
		FileType: "membership",
		FileID:   "123456",
		Total:    "1",
		CreateBy: uid,
		VerifyBy: uid,
	}
	reqs := []storage.UploadServiceRequest{p2, p1, p4, p3, p5}

	for _, r := range reqs {
		_, err := st.CreateUploadSvcRequest(context.TODO(), r)
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
		if r.Status == "ACCEPTED" {
			if err := st.AcceptUploadSvcRequest(context.TODO(), r); err != nil {
				t.Fatalf("got error %v, want nil", err)
			}
		} else {
			if err := st.RejectUploadSvcRequest(context.TODO(), r); err != nil {
				t.Fatalf("got error %v, want nil", err)
			}
		}
	}

	tests := []struct {
		desc   string
		want   []storage.UploadServiceRequest
		filter storage.UploadSvcRequestFilter
	}{
		{
			desc:   "All",
			filter: storage.UploadSvcRequestFilter{},
			want:   []storage.UploadServiceRequest{p2, p1, p4, p3, p5},
		},
		{
			desc: "Limit Org",
			filter: storage.UploadSvcRequestFilter{
				OrgID: oid3,
			},
			want: []storage.UploadServiceRequest{p5},
		},
		{
			desc: "Limit Status",
			filter: storage.UploadSvcRequestFilter{
				OrgID:  oid1,
				Status: []string{"ACCEPTED"},
			},
			want: []storage.UploadServiceRequest{p1, p2, p3},
		},
		{
			desc: "Limit SvcName",
			filter: storage.UploadSvcRequestFilter{
				OrgID:   oid1,
				SvcName: []string{"svcName", "svcName1"},
			},
			want: []storage.UploadServiceRequest{p1, p2},
		},
		{
			desc: "Limit Partner",
			filter: storage.UploadSvcRequestFilter{
				OrgID:   oid1,
				Partner: []string{"P1", "P2"},
			},
			want: []storage.UploadServiceRequest{p1, p2},
		},
	}

	o := cmp.Options{
		cmpopts.IgnoreFields(storage.UploadServiceRequest{}, "Created", "ID", "Verified", "Total"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			got, err := st.ListUploadSvcRequest(context.TODO(), test.filter)
			if err != nil {
				t.Fatalf("got error %v, want nil", err)
			}
			if !cmp.Equal(test.want, got, o) {
				t.Fatal(cmp.Diff(test.want, got, o))
			}
		})
	}
}

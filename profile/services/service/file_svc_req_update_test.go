package service

import (
	"context"
	"testing"

	spb "brank.as/petnet/gunk/dsa/v2/service"
	emlc "brank.as/petnet/profile/core/email"
	"brank.as/petnet/profile/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestUploadService(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	test.NewNullLogger()
	st := newTestStorage(t)
	mlr := newTestMailer(t)
	s := New(st, emlc.New(mlr))
	uid := uuid.NewString()
	oid1 := uuid.NewString()
	oid2 := uuid.NewString()
	oid3 := uuid.NewString()
	in := []*spb.AddUploadSvcRequestRequest{
		{
			OrgID:    oid1,
			Partner:  "P1",
			SvcName:  "SvcName1",
			FileType: "FileType1",
			FileID:   "FileID1",
			CreateBy: uid,
		},
		{
			OrgID:    oid1,
			Partner:  "P2",
			SvcName:  "SvcName2",
			FileType: "FileType2",
			FileID:   "FileID2",
			CreateBy: uid,
		},
		{
			OrgID:    oid1,
			Partner:  "P3",
			SvcName:  "SvcName3",
			FileType: "FileType3",
			FileID:   "FileID3",
			CreateBy: uid,
		},
		{
			OrgID:    oid2,
			Partner:  "P4",
			SvcName:  "SvcName4",
			FileType: "FileType4",
			FileID:   "FileID4",
			CreateBy: uid,
		},
		{
			OrgID:    oid3,
			Partner:  "P5",
			SvcName:  "SvcName5",
			FileType: "FileType5",
			FileID:   "FileID5",
			CreateBy: uid,
		},
	}
	want := []*spb.UploadSvcResponse{
		{
			OrgID:    oid1,
			Partner:  "P1",
			SvcName:  "SvcName1",
			FileType: "FileType1",
			FileID:   "FileID1",
			CreateBy: uid,
			Total:    "5",
			Status:   "ACCEPTED",
			VerifyBy: uid,
		},
		{
			OrgID:    oid1,
			Partner:  "P2",
			SvcName:  "SvcName2",
			FileType: "FileType2",
			FileID:   "FileID2",
			CreateBy: uid,
			Total:    "5",
			Status:   "ACCEPTED",
			VerifyBy: uid,
		},
		{
			OrgID:    oid1,
			Partner:  "P3",
			SvcName:  "SvcName3",
			FileType: "FileType3",
			FileID:   "FileID3",
			CreateBy: uid,
			Total:    "5",
			Status:   "ACCEPTED",
			VerifyBy: uid,
		},
		{
			OrgID:    oid2,
			Partner:  "P4",
			SvcName:  "SvcName4",
			FileType: "FileType4",
			FileID:   "FileID4",
			CreateBy: uid,
			Total:    "5",
			Status:   "ACCEPTED",
			VerifyBy: uid,
		},
		{
			OrgID:    oid3,
			Partner:  "P5",
			SvcName:  "SvcName5",
			FileType: "FileType5",
			FileID:   "FileID5",
			CreateBy: uid,
			Total:    "5",
			Status:   "REJECTED",
			VerifyBy: uid,
		},
	}
	if _, err := st.CreateOrgProfile(ctx, &storage.OrgProfile{
		OrgID:  oid1,
		UserID: uid,
	}); err != nil {
		t.Fatal(err)
	}

	for _, v := range in {
		if _, err := s.AddUploadSvcRequest(ctx, &spb.AddUploadSvcRequestRequest{
			OrgID:    v.OrgID,
			Partner:  v.Partner,
			SvcName:  v.SvcName,
			FileType: v.FileType,
			FileID:   v.FileID,
			CreateBy: v.CreateBy,
		}); err != nil {
			t.Fatal(err)
		}
	}

	got, err := s.ListUploadSvcRequest(ctx, &spb.ListUploadSvcRequestRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.GetResults()) != len(want) {
		t.Fatalf("got: %v, want: %v", len(got.GetResults()), len(want))
	}
	if got.GetResults()[0].Created.AsTime().IsZero() {
		t.Fatal("Created time not set")
	}

	o := cmp.Options{
		cmpopts.IgnoreFields(spb.UploadSvcResponse{}, "Verified", "Created", "ID", "VerifyBy", "Status"),
		cmpopts.IgnoreUnexported(spb.UploadSvcResponse{}),
	}
	if !cmp.Equal(want, got.GetResults(), o) {
		t.Fatal(cmp.Diff(want, got.GetResults(), o))
	}
	// Accept Service Request
	if _, err := s.AcceptUploadSvcRequest(ctx, &spb.AcceptUploadSvcRequestRequest{
		OrgID:    in[0].OrgID,
		Partner:  in[0].Partner,
		SvcName:  in[0].SvcName,
		VerifyBy: uid,
		FileType: in[0].FileType,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.AcceptUploadSvcRequest(ctx, &spb.AcceptUploadSvcRequestRequest{
		OrgID:    in[1].OrgID,
		Partner:  in[1].Partner,
		SvcName:  in[1].SvcName,
		VerifyBy: uid,
		FileType: in[1].FileType,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.AcceptUploadSvcRequest(ctx, &spb.AcceptUploadSvcRequestRequest{
		OrgID:    in[2].OrgID,
		Partner:  in[2].Partner,
		SvcName:  in[2].SvcName,
		VerifyBy: uid,
		FileType: in[2].FileType,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.AcceptUploadSvcRequest(ctx, &spb.AcceptUploadSvcRequestRequest{
		OrgID:    in[3].OrgID,
		Partner:  in[3].Partner,
		SvcName:  in[3].SvcName,
		VerifyBy: uid,
		FileType: in[3].FileType,
	}); err != nil {
		t.Fatal(err)
	}
	AcGot, err := s.ListUploadSvcRequest(ctx, &spb.ListUploadSvcRequestRequest{
		Status: []string{"ACCEPTED"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	if AcGot.GetResults()[0].Status != "ACCEPTED" {
		t.Fatal("Accepted Is not set")
	}
	if AcGot.GetResults()[1].Status != "ACCEPTED" {
		t.Fatal("Accepted Is not set")
	}

	// Rejected Service Request
	if _, err := s.RejectUploadSvcRequest(ctx, &spb.RejectUploadSvcRequestRequest{
		OrgID:    in[4].OrgID,
		Partner:  in[4].Partner,
		SvcName:  in[4].SvcName,
		VerifyBy: uid,
		FileType: in[4].FileType,
	}); err != nil {
		t.Fatal(err)
	}
	ReGot, err := s.ListUploadSvcRequest(ctx, &spb.ListUploadSvcRequestRequest{
		Status: []string{"REJECTED"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	if ReGot.GetResults()[0].Status != "REJECTED" {
		t.Fatal("Rejected Is not set")
	}
	// List Test
	tests := []struct {
		desc   string
		want   []*spb.UploadSvcResponse
		filter *spb.ListUploadSvcRequestRequest
	}{
		{
			desc:   "All",
			filter: &spb.ListUploadSvcRequestRequest{},
			want:   []*spb.UploadSvcResponse{want[0], want[1], want[2], want[3], want[4]},
		},
		{
			desc: "Limit Org",
			filter: &spb.ListUploadSvcRequestRequest{
				OrgID: oid3,
			},
			want: []*spb.UploadSvcResponse{want[4]},
		},
		{
			desc: "Limit Status",
			filter: &spb.ListUploadSvcRequestRequest{
				OrgID:  oid1,
				Status: []string{"ACCEPTED"},
			},
			want: []*spb.UploadSvcResponse{want[0], want[1], want[2]},
		},
		{
			desc: "Limit SvcName",
			filter: &spb.ListUploadSvcRequestRequest{
				OrgID:    oid1,
				SvcNames: []string{"SvcName1", "SvcName2", "SvcName3"},
			},
			want: []*spb.UploadSvcResponse{want[0], want[1], want[2]},
		},
		{
			desc: "Limit Partner",
			filter: &spb.ListUploadSvcRequestRequest{
				OrgID:    oid1,
				Partners: []string{"P1", "P2"},
			},
			want: []*spb.UploadSvcResponse{want[0], want[1]},
		},
	}
	to := cmp.Options{
		cmpopts.IgnoreFields(spb.UploadSvcResponse{}, "Created", "ID", "Verified", "Total"),
		cmpopts.IgnoreUnexported(spb.UploadSvcResponse{}),
	}
	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			got, err := s.ListUploadSvcRequest(context.TODO(), test.filter)
			if err != nil {
				t.Fatalf("got error %v, want nil", err)
			}
			if !cmp.Equal(test.want, got.GetResults(), to) {
				t.Fatal(cmp.Diff(test.want, got.GetResults(), to))
			}
		})
	}

	// Updated Service Request
	if _, err := s.UpdateUploadSvcRequest(ctx, &spb.UpdateUploadSvcRequestRequest{
		OrgID:    in[1].OrgID,
		Partner:  in[1].Partner,
		SvcName:  in[1].SvcName,
		Status:   "ACCEPTED",
		FileType: in[1].FileType,
		FileID:   "123456",
		CreateBy: uid,
		VerifyBy: uid,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.UpdateUploadSvcRequest(ctx, &spb.UpdateUploadSvcRequestRequest{
		OrgID:    in[3].OrgID,
		Partner:  in[3].Partner,
		SvcName:  in[3].SvcName,
		Status:   "ACCEPTED",
		FileType: in[3].FileType,
		FileID:   "987654",
		CreateBy: uid,
		VerifyBy: uid,
	}); err != nil {
		t.Fatal(err)
	}
	UpGot, err := s.ListUploadSvcRequest(ctx, &spb.ListUploadSvcRequestRequest{
		Status: []string{""},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	if UpGot.GetResults()[0].FileID != "123456" {
		t.Fatal("Rejected Is not set")
	}
	if UpGot.GetResults()[1].FileID != "987654" {
		t.Fatal("Rejected Is not set")
	}
}

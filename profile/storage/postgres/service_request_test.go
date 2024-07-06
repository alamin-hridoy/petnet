package postgres_test

import (
	"context"
	"testing"

	"brank.as/petnet/profile/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestServiceRequest(t *testing.T) {
	st := newTestStorage(t)
	ctx := context.TODO()
	o := cmp.Options{
		cmpopts.IgnoreFields(storage.ServiceRequest{}, "Updated", "Created", "ID", "Applied"),
	}

	oid := uuid.NewString()
	want := storage.ServiceRequest{
		OrgID:       oid,
		SvcName:     "svcName",
		CompanyName: "compName",
		Partner:     "WU",
		Total:       "1",
	}

	csres, err := st.CreateSvcRequest(ctx, want)
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if csres.Updated.IsZero() {
		t.Fatal("returned empty Updated time")
	}
	if csres.Created.IsZero() {
		t.Fatal("returned empty Created time")
	}
	if !csres.Applied.Time.IsZero() {
		t.Fatal("applied should not be set")
	}
	got1, err := st.ListSvcRequest(ctx, storage.SvcRequestFilter{
		Status: []string{"PENDING", "ACCEPTED", "REJECTED"},
	})
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if len(got1) != 0 {
		t.Fatalf("want o service request, got %v", len(got1))
	}
	if err = st.ApplySvcRequest(ctx, want); err != nil {
		t.Fatalf("got error %v, want nil", err)
	}

	got2, err := st.ListSvcRequest(ctx, storage.SvcRequestFilter{})
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if len(got2) != 1 {
		t.Fatalf("want 1 service request, got %v", len(got2))
	}
	if got2[0].Updated == want.Updated {
		t.Fatal("Updated field has not changed")
	}
	if got2[0].Applied.Time.IsZero() {
		t.Fatal("returned empty Applied time")
	}
	want.Status = "PENDING"
	if !cmp.Equal(want, got2[0], o) {
		t.Fatal(cmp.Diff(want, got2[0], o))
	}

	want.Remarks = "remarks"
	want.Status = "ACCEPTED"
	want.Partner = "WU"
	want.UpdatedBy = "John"
	want.Enabled = true
	if err = st.AcceptSvcRequest(ctx, want); err != nil {
		t.Fatalf("got error %v, want nil", err)
	}

	got3, err := st.ListSvcRequest(ctx, storage.SvcRequestFilter{})
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if len(got3) != 1 {
		t.Fatalf("want 1 service request, got %v", len(got3))
	}
	if got3[0].Updated == got2[0].Updated {
		t.Fatal("Updated field has not changed")
	}
	if got3[0].Applied.Time.IsZero() {
		t.Fatal("returned empty Applied time")
	}
	if !cmp.Equal(want, got3[0], o) {
		t.Fatal(cmp.Diff(want, got3[0], o))
	}

	want.Remarks = "remarks"
	want.Status = "REJECTED"
	want.Partner = "WU"
	want.UpdatedBy = "John"
	if err = st.RejectSvcRequest(ctx, want); err != nil {
		t.Fatalf("got error %v, want nil", err)
	}

	got4, err := st.ListSvcRequest(ctx, storage.SvcRequestFilter{})
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if len(got4) != 1 {
		t.Fatalf("want 1 service request, got %v", len(got4))
	}
	if got4[0].Updated == got3[0].Updated {
		t.Fatal("Updated field has not changed")
	}
	if got4[0].Applied.Time.IsZero() {
		t.Fatal("returned empty Applied time")
	}
	if !cmp.Equal(want, got4[0], o) {
		t.Fatal(cmp.Diff(want, got4[0], o))
	}

	want.UpdatedBy = "Emily"
	want.Enabled = true
	if err = st.EnableSvcRequest(ctx, want); err != nil {
		t.Fatalf("got error %v, want nil", err)
	}

	got5, err := st.ListSvcRequest(ctx, storage.SvcRequestFilter{})
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if len(got5) != 1 {
		t.Fatalf("want 1 service request, got %v", len(got5))
	}
	if got5[0].Updated == got4[0].Updated {
		t.Fatal("Updated field has not changed")
	}
	if got5[0].Applied.Time.IsZero() {
		t.Fatal("returned empty Applied time")
	}
	if !cmp.Equal(want, got5[0], o) {
		t.Fatal(cmp.Diff(want, got5[0], o))
	}

	want.UpdatedBy = "John"
	want.Enabled = false
	if err = st.DisableSvcRequest(ctx, want); err != nil {
		t.Fatalf("got error %v, want nil", err)
	}

	got6, err := st.ListSvcRequest(ctx, storage.SvcRequestFilter{})
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if len(got6) != 1 {
		t.Fatalf("want 1 service request, got %v", len(got6))
	}
	if got6[0].Updated == got5[0].Updated {
		t.Fatal("Updated field has not changed")
	}
	if got6[0].Applied.Time.IsZero() {
		t.Fatal("returned empty Applied time")
	}
	if !cmp.Equal(want, got6[0], o) {
		t.Fatal(cmp.Diff(want, got6[0], o))
	}
	want.Remarks = got6[0].Remarks
	want.Status = "ACCEPTED"
	want.Partner = got6[0].Partner
	want.SvcName = got6[0].SvcName
	want.UpdatedBy = got6[0].UpdatedBy
	want.Enabled = true
	if err = st.AcceptSvcRequest(ctx, want); err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	valRes, err := st.ValidateSvcRequest(ctx, storage.ValidateSvcRequestFilter{
		OrgID:               oid,
		Partner:             want.Partner,
		SvcName:             want.SvcName,
		IsAnyPartnerEnabled: false,
	})
	if err != nil {
		t.Fatalf("validate partner error: %v", err)
	}
	if valRes.Enabled != true {
		t.Fatal("validated is not set")
	}
	val2Res, err := st.ValidateSvcRequest(ctx, storage.ValidateSvcRequestFilter{
		OrgID:               oid,
		Partner:             "TF",
		SvcName:             want.SvcName,
		IsAnyPartnerEnabled: true,
	})
	if err != nil {
		t.Fatalf("validate partner error: %v", err)
	}
	if val2Res.Enabled != true {
		t.Fatal("any validated is not set")
	}
	val3Res, err := st.ValidateSvcRequest(ctx, storage.ValidateSvcRequestFilter{
		OrgID:               oid,
		Partner:             "TF",
		SvcName:             want.SvcName,
		IsAnyPartnerEnabled: false,
	})
	if err != nil {
		t.Fatalf("validate partner error: %v", err)
	}
	if val3Res.Enabled != false {
		t.Fatal("validation is not set")
	}
	want.Status = "PARTNERDRAFT"
	if err = st.SetStatusSvcRequest(ctx, want); err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	got7, err := st.ListSvcRequest(ctx, storage.SvcRequestFilter{})
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if got7[0].Status != want.Status {
		t.Fatal("set status field has not changed")
	}
	if got7[0].Status != "PARTNERDRAFT" {
		t.Fatal("set status field has not changed")
	}
	// add remark
	want.Remarks = "Sample Remak"
	want.OrgID = got7[0].OrgID
	want.SvcName = got7[0].SvcName
	if err = st.AddRemarkSvcRequest(ctx, want); err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	got8, err := st.ListSvcRequest(ctx, storage.SvcRequestFilter{})
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	if got8[0].Remarks != "Sample Remak" {
		t.Fatal("add remark field has not changed")
	}
	wantt := storage.UpdateServiceRequestOrgID{
		OldOrgID: csres.OrgID,
		NewOrgID: "3dc3fea6-50df-41f9-aa50-80ab4128246a",
		Status:   "REJECTED",
	}
	_, err = st.UpdateServiceRequestByOrgID(ctx, wantt)
	if err != nil {
		t.Fatalf("got error %v, want nil", err)
	}
	// reject
	lRes, _ := st.ListSvcRequest(ctx, storage.SvcRequestFilter{})
	if len(lRes) > 0 {
		for _, lr := range lRes {
			st.RejectSvcRequest(ctx, storage.ServiceRequest{
				SvcName:   lr.SvcName,
				OrgID:     lr.OrgID,
				Partner:   lr.Partner,
				UpdatedBy: lr.UpdatedBy,
			})
		}
	}
	if len(lRes) > 0 {
		for _, lr := range lRes {
			err = st.RemoveSvcRequest(ctx, storage.ServiceRequest{
				OrgID:   lr.OrgID,
				Partner: lr.Partner,
				SvcName: lr.SvcName,
			})
			if err != nil {
				if err != storage.NotFound {
					t.Fatal("Remove Svc Request failed")
				}
			}
		}
	}
	lNewRes, err := st.ListSvcRequest(ctx, storage.SvcRequestFilter{})
	if len(lNewRes) > 0 {
		t.Fatal("Remove Svc Request failed")
	}
}

func TestListServiceRequest(t *testing.T) {
	st := newTestStorage(t)

	oid1 := uuid.NewString()
	oid2 := uuid.NewString()
	oid3 := uuid.NewString()
	p2 := storage.ServiceRequest{
		OrgID:       oid1,
		SvcName:     "svcName2",
		CompanyName: "compName2",
		Partner:     "P2",
		Status:      "ACCEPTED",
		UpdatedBy:   "P2",
		Enabled:     true,
	}
	p1 := storage.ServiceRequest{
		OrgID:       oid1,
		SvcName:     "svcName1",
		CompanyName: "compName1",
		Partner:     "P1",
		Status:      "ACCEPTED",
		UpdatedBy:   "P1",
		Enabled:     true,
	}
	p4 := storage.ServiceRequest{
		OrgID:       oid2,
		SvcName:     "svcName4",
		CompanyName: "compName4",
		Partner:     "P4",
		Status:      "REJECTED",
		UpdatedBy:   "P4",
	}
	p3 := storage.ServiceRequest{
		OrgID:       oid1,
		SvcName:     "svcName3",
		CompanyName: "compName3",
		Partner:     "P3",
		Status:      "ACCEPTED",
		UpdatedBy:   "P3",
		Enabled:     true,
	}
	p5 := storage.ServiceRequest{
		OrgID:       oid3,
		SvcName:     "svcName5",
		CompanyName: "compName5",
		Partner:     "P5",
		Status:      "REJECTED",
		UpdatedBy:   "P5",
	}
	reqs := []storage.ServiceRequest{p2, p1, p4, p3, p5}

	for _, r := range reqs {
		_, err := st.CreateSvcRequest(context.TODO(), r)
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}

		if r.Status == "ACCEPTED" {
			if err := st.AcceptSvcRequest(context.TODO(), r); err != nil {
				t.Fatalf("got error %v, want nil", err)
			}
		} else {
			if err := st.RejectSvcRequest(context.TODO(), r); err != nil {
				t.Fatalf("got error %v, want nil", err)
			}
		}
	}

	tests := []struct {
		desc   string
		want   []storage.ServiceRequest
		filter storage.SvcRequestFilter
	}{
		{
			desc: "All",
			filter: storage.SvcRequestFilter{
				OrgID: []string{oid1, oid2, oid3},
			},
			want: []storage.ServiceRequest{p5, p3, p4, p1, p2},
		},
		{
			desc: "Sort CreatedCol ASC",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1, oid2, oid3},
				SortOrder:    "ASC",
				SortByColumn: "CREATED",
			},
			want: []storage.ServiceRequest{p2, p1, p4, p3, p5},
		},
		{
			desc: "Sort CreatedCol DESC",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1, oid2, oid3},
				SortOrder:    "DESC",
				SortByColumn: "CREATED",
			},
			want: []storage.ServiceRequest{p5, p3, p4, p1, p2},
		},
		{
			desc: "Sort CompanyNameCol ASC",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1, oid2, oid3},
				SortOrder:    "ASC",
				SortByColumn: "COMPANYNAME",
			},
			want: []storage.ServiceRequest{p1, p2, p3, p4, p5},
		},
		{
			desc: "Sort CompanyNameCol DESC",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1, oid2, oid3},
				SortOrder:    "DESC",
				SortByColumn: "COMPANYNAME",
			},
			want: []storage.ServiceRequest{p5, p4, p3, p2, p1},
		},
		{
			desc: "Sort ServiceNameCol ASC",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1, oid2, oid3},
				SortOrder:    "ASC",
				SortByColumn: "SERVICENAME",
			},
			want: []storage.ServiceRequest{p1, p2, p3, p4, p5},
		},
		{
			desc: "Sort ServiceNameCol DESC",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1, oid2, oid3},
				SortOrder:    "DESC",
				SortByColumn: "SERVICENAME",
			},
			want: []storage.ServiceRequest{p5, p4, p3, p2, p1},
		},
		{
			desc: "Sort StatusCol ASC",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1, oid2, oid3},
				SortOrder:    "ASC",
				SortByColumn: "STATUS",
			},
			want: []storage.ServiceRequest{p2, p1, p3, p4, p5},
		},
		{
			desc: "Sort StatusCol DESC",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1, oid2, oid3},
				SortOrder:    "DESC",
				SortByColumn: "STATUS",
			},
			want: []storage.ServiceRequest{p4, p5, p2, p1, p3},
		},
		{
			desc: "Sort PartnerCol ASC",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1, oid2, oid3},
				SortOrder:    "ASC",
				SortByColumn: "PARTNER",
			},
			want: []storage.ServiceRequest{p1, p2, p3, p4, p5},
		},
		{
			desc: "Sort PartnerCol DESC",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1, oid2, oid3},
				SortOrder:    "DESC",
				SortByColumn: "PARTNER",
			},
			want: []storage.ServiceRequest{p5, p4, p3, p2, p1},
		},
		{
			desc: "Sort LastUpdatedCol ASC",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1, oid2, oid3},
				SortOrder:    "ASC",
				SortByColumn: "LASTUPDATED",
			},
			want: []storage.ServiceRequest{p2, p1, p4, p3, p5},
		},
		{
			desc: "Sort LastUpdatedCol DESC",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1, oid2, oid3},
				SortOrder:    "DESC",
				SortByColumn: "LASTUPDATED",
			},
			want: []storage.ServiceRequest{p5, p3, p4, p1, p2},
		},
		{
			desc: "Sort LastUpdatedByCol ASC",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1, oid2, oid3},
				SortOrder:    "ASC",
				SortByColumn: "UPDATEDBY",
			},
			want: []storage.ServiceRequest{p1, p2, p3, p4, p5},
		},
		{
			desc: "Sort LastUpdatedByCol DESC",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1, oid2, oid3},
				SortOrder:    "DESC",
				SortByColumn: "UPDATEDBY",
			},
			want: []storage.ServiceRequest{p5, p4, p3, p2, p1},
		},
		{
			desc: "Limit OrgID",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1},
				SortOrder:    "ASC",
				SortByColumn: "COMPANYNAME",
			},
			want: []storage.ServiceRequest{p1, p2, p3},
		},
		{
			desc: "Limit Status",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1, oid2, oid3},
				Status:       []string{"ACCEPTED"},
				SortOrder:    "ASC",
				SortByColumn: "COMPANYNAME",
			},
			want: []storage.ServiceRequest{p1, p2, p3},
		},
		{
			desc: "Limit SvcName",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1, oid2, oid3},
				SvcName:      []string{"svcName1", "svcName2"},
				SortOrder:    "ASC",
				SortByColumn: "COMPANYNAME",
			},
			want: []storage.ServiceRequest{p1, p2},
		},
		{
			desc: "Limit Partner",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1, oid2, oid3},
				Partner:      []string{"P1", "P2"},
				SortOrder:    "ASC",
				SortByColumn: "COMPANYNAME",
			},
			want: []storage.ServiceRequest{p1, p2},
		},
		{
			desc: "Limit",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1, oid2, oid3},
				Limit:        3,
				SortOrder:    "ASC",
				SortByColumn: "COMPANYNAME",
			},
			want: []storage.ServiceRequest{p1, p2, p3},
		},
		{
			desc: "Offset",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1, oid2, oid3},
				Offset:       3,
				SortOrder:    "ASC",
				SortByColumn: "COMPANYNAME",
			},
			want: []storage.ServiceRequest{p4, p5},
		},
		{
			desc: "Limit/Offset",
			filter: storage.SvcRequestFilter{
				OrgID:        []string{oid1, oid2, oid3},
				Limit:        2,
				Offset:       2,
				SortOrder:    "ASC",
				SortByColumn: "COMPANYNAME",
			},
			want: []storage.ServiceRequest{p3, p4},
		},
	}

	o := cmp.Options{
		cmpopts.IgnoreFields(storage.ServiceRequest{}, "Updated", "Created", "ID", "Applied", "Total", "OrgID"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			got, err := st.ListSvcRequest(context.TODO(), test.filter)
			if err != nil {
				t.Fatalf("got error %v, want nil", err)
			}

			if !cmp.Equal(test.want, got, o) {
				t.Fatal(cmp.Diff(test.want, got, o))
			}
		})
	}
}

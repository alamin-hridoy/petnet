package service

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	spb "brank.as/petnet/gunk/dsa/v2/service"
	emlc "brank.as/petnet/profile/core/email"
	"brank.as/petnet/profile/integrations/email"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/profile/storage/postgres"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestService(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	test.NewNullLogger()

	st := newTestStorage(t)
	mlr := newTestMailer(t)
	s := New(st, emlc.New(mlr))

	ptnrs := []*storage.PartnerList{
		{
			Stype:       "WU",
			Name:        "Western Union",
			Status:      "ENABLED",
			ServiceName: "REMITTANCE",
		},
		{
			Stype:       "IR",
			Name:        "IRemit",
			Status:      "ENABLED",
			ServiceName: "REMITTANCE",
		},
		{
			Stype:       "TF",
			Name:        "Transfast",
			Status:      "ENABLED",
			ServiceName: "REMITTANCE",
		},
	}

	oid := uuid.NewString()
	uid := uuid.NewString()
	if _, err := st.CreateOrgProfile(ctx, &storage.OrgProfile{
		OrgID:  oid,
		UserID: uid,
	}); err != nil {
		t.Fatal(err)
	}

	st.CreatePartnerList(ctx, &storage.PartnerList{
		ID:               "",
		Stype:            "",
		Name:             "",
		Created:          time.Time{},
		Updated:          time.Time{},
		Deleted:          sql.NullTime{},
		Status:           "",
		TransactionTypes: "",
		ServiceName:      "",
		UpdatedBy:        "",
		DisableReason:    sql.NullString{},
		Platform:         "",
	})

	if len(ptnrs) > 0 {
		for _, v := range ptnrs {
			if _, err := st.CreatePartnerList(ctx, v); err != nil {
				continue
			}
		}
	}

	if _, err := s.AddServiceRequest(ctx, &spb.AddServiceRequestRequest{
		OrgID: oid,
		Type:  spb.ServiceType_REMITTANCE,
		Partners: []string{
			spb.RemittancePartner_WU.String(),
			spb.RemittancePartner_IR.String(),
			spb.RemittancePartner_TF.String(),
		},
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := s.AddServiceRequest(ctx, &spb.AddServiceRequestRequest{
		OrgID: oid,
		Type:  spb.ServiceType_REMITTANCE,
		Partners: []string{
			spb.RemittancePartner_WU.String(),
			spb.RemittancePartner_IR.String(),
			spb.RemittancePartner_TF.String(),
		},
	}); err != nil {
		t.Fatal("should be able to update same services")
	}

	want := []*spb.ServiceRequest{
		{
			OrgID:       oid,
			CompanyName: "",
			Partner:     "TF",
			Type:        spb.ServiceType_REMITTANCE,
			Status:      spb.ServiceRequestStatus_NOSTATUS,
			Enabled:     false,
			Remarks:     "",
			Applied:     nil,
			UpdatedBy:   "",
		},
		{
			OrgID:       oid,
			CompanyName: "",
			Partner:     "IR",
			Type:        spb.ServiceType_REMITTANCE,
			Status:      spb.ServiceRequestStatus_NOSTATUS,
			Enabled:     false,
			Remarks:     "",
			Applied:     nil,
			UpdatedBy:   "",
		},
		{
			OrgID:       oid,
			CompanyName: "",
			Partner:     "WU",
			Type:        spb.ServiceType_REMITTANCE,
			Status:      spb.ServiceRequestStatus_NOSTATUS,
			Enabled:     false,
			Remarks:     "",
			Applied:     nil,
			UpdatedBy:   "",
		},
	}
	got, err := s.ListServiceRequest(ctx, &spb.ListServiceRequestRequest{
		Types:  []spb.ServiceType{spb.ServiceType_REMITTANCE},
		OrgIDs: []string{oid},
		Partners: []string{
			spb.RemittancePartner_WU.String(),
			spb.RemittancePartner_IR.String(),
			spb.RemittancePartner_TF.String(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.GetServiceRequst()) != 3 {
		t.Fatalf("got: %v, want: %v", len(got.GetServiceRequst()), 3)
	}
	if got.GetServiceRequst()[0].Created.AsTime().IsZero() {
		t.Fatal("Created time not set")
	}
	if got.GetServiceRequst()[0].Updated.AsTime().IsZero() {
		t.Fatal("Updated time not set")
	}

	o := cmp.Options{
		cmpopts.IgnoreFields(spb.ServiceRequest{}, "Created", "Updated", "ID"),
		cmpopts.IgnoreUnexported(spb.ServiceRequest{}),
	}
	if !cmp.Equal(want, got.GetServiceRequst(), o) {
		t.Fatal(cmp.Diff(want, got.GetServiceRequst(), o))
	}

	if _, err := s.ApplyServiceRequest(ctx, &spb.ApplyServiceRequestRequest{
		OrgID: oid,
		Type:  spb.ServiceType_REMITTANCE,
	}); err != nil {
		t.Fatal(err)
	}

	got, err = s.ListServiceRequest(ctx, &spb.ListServiceRequestRequest{
		Types:  []spb.ServiceType{spb.ServiceType_REMITTANCE},
		OrgIDs: []string{oid},
		Partners: []string{
			spb.RemittancePartner_WU.String(),
			spb.RemittancePartner_IR.String(),
			spb.RemittancePartner_TF.String(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.GetServiceRequst()[0].Applied.AsTime().IsZero() {
		t.Fatal("Applied time not set")
	}

	o = cmp.Options{
		cmpopts.IgnoreFields(spb.ServiceRequest{}, "Created", "Updated", "Applied", "ID"),
		cmpopts.IgnoreUnexported(spb.ServiceRequest{}),
	}
	for _, r := range want {
		r.Status = spb.ServiceRequestStatus_PENDING
	}
	if !cmp.Equal(want, got.GetServiceRequst(), o) {
		t.Fatal(cmp.Diff(want, got.GetServiceRequst(), o))
	}
	// Accept Service Request
	if _, err := s.AcceptServiceRequest(ctx, &spb.ServiceStatusRequestRequest{
		OrgID:     oid,
		UpdatedBy: uid,
		Partner:   spb.RemittancePartner_TF.String(),
		SvcName:   spb.ServiceType_REMITTANCE.String(),
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := s.AcceptServiceRequest(ctx, &spb.ServiceStatusRequestRequest{
		OrgID:     oid,
		UpdatedBy: uid,
		Partner:   spb.RemittancePartner_IR.String(),
		SvcName:   spb.ServiceType_REMITTANCE.String(),
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := s.AcceptServiceRequest(ctx, &spb.ServiceStatusRequestRequest{
		OrgID:     oid,
		UpdatedBy: uid,
		Partner:   spb.RemittancePartner_WU.String(),
		SvcName:   spb.ServiceType_REMITTANCE.String(),
	}); err != nil {
		t.Fatal(err)
	}

	AcGot, err := s.ListServiceRequest(ctx, &spb.ListServiceRequestRequest{
		Types:  []spb.ServiceType{spb.ServiceType_REMITTANCE},
		OrgIDs: []string{oid},
		Partners: []string{
			spb.RemittancePartner_WU.String(),
			spb.RemittancePartner_IR.String(),
			spb.RemittancePartner_TF.String(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if AcGot.GetServiceRequst()[0].Status.String() != spb.ServiceRequestStatus_ACCEPTED.String() {
		t.Fatal("Accepted Is not set")
	}

	o = cmp.Options{
		cmpopts.IgnoreFields(spb.ServiceRequest{}, "Created", "Updated", "Applied", "UpdatedBy", "ID"),
		cmpopts.IgnoreUnexported(spb.ServiceRequest{}),
	}
	for _, r := range want {
		r.Status = spb.ServiceRequestStatus_ACCEPTED
		r.Enabled = true
	}
	if !cmp.Equal(want, AcGot.GetServiceRequst(), o) {
		t.Fatal(cmp.Diff(want, AcGot.GetServiceRequst(), o))
	}
	// Reject Service Request
	if _, err := s.RejectServiceRequest(ctx, &spb.ServiceStatusRequestRequest{
		OrgID:     oid,
		UpdatedBy: uid,
		Partner:   spb.RemittancePartner_TF.String(),
		SvcName:   spb.ServiceType_REMITTANCE.String(),
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := s.RejectServiceRequest(ctx, &spb.ServiceStatusRequestRequest{
		OrgID:     oid,
		UpdatedBy: uid,
		Partner:   spb.RemittancePartner_IR.String(),
		SvcName:   spb.ServiceType_REMITTANCE.String(),
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := s.RejectServiceRequest(ctx, &spb.ServiceStatusRequestRequest{
		OrgID:     oid,
		UpdatedBy: uid,
		Partner:   spb.RemittancePartner_WU.String(),
		SvcName:   spb.ServiceType_REMITTANCE.String(),
	}); err != nil {
		t.Fatal(err)
	}

	ReGot, err := s.ListServiceRequest(ctx, &spb.ListServiceRequestRequest{
		Types:  []spb.ServiceType{spb.ServiceType_REMITTANCE},
		OrgIDs: []string{oid},
		Partners: []string{
			spb.RemittancePartner_WU.String(),
			spb.RemittancePartner_IR.String(),
			spb.RemittancePartner_TF.String(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if ReGot.GetServiceRequst()[0].Status.String() != spb.ServiceRequestStatus_REJECTED.String() {
		t.Fatal("Rejected Is not set")
	}

	o = cmp.Options{
		cmpopts.IgnoreFields(spb.ServiceRequest{}, "Created", "Updated", "Applied", "UpdatedBy", "ID"),
		cmpopts.IgnoreUnexported(spb.ServiceRequest{}),
	}
	for _, r := range want {
		r.Status = spb.ServiceRequestStatus_REJECTED
	}
	if !cmp.Equal(want, ReGot.GetServiceRequst(), o) {
		t.Fatal(cmp.Diff(want, ReGot.GetServiceRequst(), o))
	}
	// Enable Service Request
	if _, err := s.EnableServiceRequest(ctx, &spb.ServiceStatusRequestRequest{
		OrgID:     oid,
		UpdatedBy: uid,
		Partner:   spb.RemittancePartner_TF.String(),
		SvcName:   spb.ServiceType_REMITTANCE.String(),
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := s.EnableServiceRequest(ctx, &spb.ServiceStatusRequestRequest{
		OrgID:     oid,
		UpdatedBy: uid,
		Partner:   spb.RemittancePartner_IR.String(),
		SvcName:   spb.ServiceType_REMITTANCE.String(),
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := s.EnableServiceRequest(ctx, &spb.ServiceStatusRequestRequest{
		OrgID:     oid,
		UpdatedBy: uid,
		Partner:   spb.RemittancePartner_WU.String(),
		SvcName:   spb.ServiceType_REMITTANCE.String(),
	}); err != nil {
		t.Fatal(err)
	}

	EnGot, err := s.ListServiceRequest(ctx, &spb.ListServiceRequestRequest{
		Types:  []spb.ServiceType{spb.ServiceType_REMITTANCE},
		OrgIDs: []string{oid},
		Partners: []string{
			spb.RemittancePartner_WU.String(),
			spb.RemittancePartner_IR.String(),
			spb.RemittancePartner_TF.String(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if EnGot.GetServiceRequst()[0].Enabled != true {
		t.Fatal("Enabled is not set")
	}

	o = cmp.Options{
		cmpopts.IgnoreFields(spb.ServiceRequest{}, "Created", "Updated", "Applied", "UpdatedBy", "ID"),
		cmpopts.IgnoreUnexported(spb.ServiceRequest{}),
	}
	for _, r := range want {
		r.Enabled = true
	}
	if !cmp.Equal(want, EnGot.GetServiceRequst(), o) {
		t.Fatal(cmp.Diff(want, EnGot.GetServiceRequst(), o))
	}
	// Disable Service Request
	if _, err := s.DisableServiceRequest(ctx, &spb.ServiceStatusRequestRequest{
		OrgID:     oid,
		UpdatedBy: uid,
		Partner:   spb.RemittancePartner_TF.String(),
		SvcName:   spb.ServiceType_REMITTANCE.String(),
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := s.DisableServiceRequest(ctx, &spb.ServiceStatusRequestRequest{
		OrgID:     oid,
		UpdatedBy: uid,
		Partner:   spb.RemittancePartner_IR.String(),
		SvcName:   spb.ServiceType_REMITTANCE.String(),
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := s.DisableServiceRequest(ctx, &spb.ServiceStatusRequestRequest{
		OrgID:     oid,
		UpdatedBy: uid,
		Partner:   spb.RemittancePartner_WU.String(),
		SvcName:   spb.ServiceType_REMITTANCE.String(),
	}); err != nil {
		t.Fatal(err)
	}

	DsGot, err := s.ListServiceRequest(ctx, &spb.ListServiceRequestRequest{
		Types:  []spb.ServiceType{spb.ServiceType_REMITTANCE},
		OrgIDs: []string{oid},
		Partners: []string{
			spb.RemittancePartner_WU.String(),
			spb.RemittancePartner_IR.String(),
			spb.RemittancePartner_TF.String(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if DsGot.GetServiceRequst()[0].Enabled != false {
		t.Fatal("Disable is not set")
	}

	o = cmp.Options{
		cmpopts.IgnoreFields(spb.ServiceRequest{}, "Created", "Updated", "Applied", "UpdatedBy", "ID"),
		cmpopts.IgnoreUnexported(spb.ServiceRequest{}),
	}
	for _, r := range want {
		r.Enabled = false
	}
	if !cmp.Equal(want, DsGot.GetServiceRequst(), o) {
		t.Fatal(cmp.Diff(want, DsGot.GetServiceRequst(), o))
	}
	// Validate Service Request
	if _, err := s.AcceptServiceRequest(ctx, &spb.ServiceStatusRequestRequest{
		OrgID:     oid,
		UpdatedBy: uid,
		Partner:   spb.RemittancePartner_TF.String(),
		SvcName:   spb.ServiceType_REMITTANCE.String(),
	}); err != nil {
		t.Fatal(err)
	}
	ValGot, err := s.ValidateServiceAccess(ctx, &spb.ValidateServiceAccessRequest{
		OrgID:               oid,
		Partner:             spb.RemittancePartner_TF.String(),
		SvcName:             spb.ServiceType_REMITTANCE.String(),
		IsAnyPartnerEnabled: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if ValGot.Enabled != true {
		t.Fatal("validated is not set")
	}
	Val2Got, err := s.ValidateServiceAccess(ctx, &spb.ValidateServiceAccessRequest{
		OrgID:               oid,
		Partner:             spb.RemittancePartner_IR.String(),
		SvcName:             spb.ServiceType_REMITTANCE.String(),
		IsAnyPartnerEnabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if Val2Got.Enabled != true {
		t.Fatal("any validated is not set")
	}
	// Set Status Service Request
	if _, err := s.SetStatusUploadSvcRequest(ctx, &spb.SetStatusUploadSvcRequestRequest{
		OrgID:    oid,
		Partners: []string{spb.RemittancePartner_TF.String()},
		SvcName:  spb.ServiceType_REMITTANCE.String(),
		Status:   spb.ServiceRequestStatus_PARTNERDRAFT,
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := s.SetStatusUploadSvcRequest(ctx, &spb.SetStatusUploadSvcRequestRequest{
		OrgID:    oid,
		Status:   spb.ServiceRequestStatus_PARTNERDRAFT,
		Partners: []string{spb.RemittancePartner_IR.String()},
		SvcName:  spb.ServiceType_REMITTANCE.String(),
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := s.SetStatusUploadSvcRequest(ctx, &spb.SetStatusUploadSvcRequestRequest{
		OrgID:    oid,
		Status:   spb.ServiceRequestStatus_PARTNERDRAFT,
		Partners: []string{spb.RemittancePartner_WU.String()},
		SvcName:  spb.ServiceType_REMITTANCE.String(),
	}); err != nil {
		t.Fatal(err)
	}
	sSGot, err := s.ListServiceRequest(ctx, &spb.ListServiceRequestRequest{
		Types:  []spb.ServiceType{spb.ServiceType_REMITTANCE},
		OrgIDs: []string{oid},
		Partners: []string{
			spb.RemittancePartner_WU.String(),
			spb.RemittancePartner_IR.String(),
			spb.RemittancePartner_TF.String(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if sSGot.GetServiceRequst()[0].GetStatus() != spb.ServiceRequestStatus_PARTNERDRAFT {
		t.Fatal("set status is failed")
	}
	o = cmp.Options{
		cmpopts.IgnoreFields(spb.ServiceRequest{}, "Created", "Updated", "Applied", "UpdatedBy", "Enabled", "ID"),
		cmpopts.IgnoreUnexported(spb.ServiceRequest{}),
	}
	for _, r := range want {
		r.Status = spb.ServiceRequestStatus_PARTNERDRAFT
	}
	if !cmp.Equal(want, sSGot.GetServiceRequst(), o) {
		t.Fatal(cmp.Diff(want, sSGot.GetServiceRequst(), o))
	}
	// add remark Service Request
	if len(sSGot.GetServiceRequst()) > 0 {
		for _, v := range sSGot.GetServiceRequst() {
			if _, err := s.AddRemarkSvcRequest(ctx, &spb.AddRemarkSvcRequestRequest{
				OrgID:     v.OrgID,
				Remark:    "Sample Remark",
				UpdatedBy: uid,
				SvcName:   v.Type.String(),
			}); err != nil {
				t.Fatal(err)
			}
		}
	}
	rRGot, err := s.ListServiceRequest(ctx, &spb.ListServiceRequestRequest{
		Types:  []spb.ServiceType{spb.ServiceType_REMITTANCE},
		OrgIDs: []string{oid},
		Partners: []string{
			spb.RemittancePartner_WU.String(),
			spb.RemittancePartner_IR.String(),
			spb.RemittancePartner_TF.String(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if rRGot.GetServiceRequst()[0].GetRemarks() != "Sample Remark" {
		t.Fatal("Add Remark is failed")
	}
	o = cmp.Options{
		cmpopts.IgnoreFields(spb.ServiceRequest{}, "Created", "Updated", "Applied", "UpdatedBy", "Enabled", "ID"),
		cmpopts.IgnoreUnexported(spb.ServiceRequest{}),
	}
	for _, r := range want {
		r.Remarks = "Sample Remark"
	}
	if !cmp.Equal(want, rRGot.GetServiceRequst(), o) {
		t.Fatal(cmp.Diff(want, rRGot.GetServiceRequst(), o))
	}
}

func TestListService(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	test.NewNullLogger()

	st := newTestStorage(t)
	mlr := newTestMailer(t)
	s := New(st, emlc.New(mlr))
	ptnrs := []*storage.PartnerList{
		{
			Stype:       "WU",
			Name:        "Western Union",
			Status:      "ENABLED",
			ServiceName: "REMITTANCE",
		},
		{
			Stype:       "IR",
			Name:        "IRemit",
			Status:      "ENABLED",
			ServiceName: "REMITTANCE",
		},
		{
			Stype:       "TF",
			Name:        "Transfast",
			Status:      "ENABLED",
			ServiceName: "REMITTANCE",
		},
	}

	oid1 := uuid.NewString()
	oid2 := uuid.NewString()
	oids := []string{oid1, oid2}
	cn := []string{"A", "B"}
	for i, id := range oids {
		if _, err := st.CreateOrgProfile(ctx, &storage.OrgProfile{
			OrgID:  id,
			UserID: uuid.NewString(),
			BusinessInfo: storage.BusinessInfo{
				CompanyName: cn[i],
			},
		}); err != nil {
			t.Fatal(err)
		}

		if len(ptnrs) > 0 {
			for _, v := range ptnrs {
				if _, err := st.CreatePartnerList(ctx, v); err != nil {
					continue
				}
			}
		}

		if _, err := s.AddServiceRequest(ctx, &spb.AddServiceRequestRequest{
			OrgID: id,
			Type:  spb.ServiceType_REMITTANCE,
			Partners: []string{
				spb.RemittancePartner_WU.String(),
				spb.RemittancePartner_IR.String(),
				spb.RemittancePartner_TF.String(),
			},
		}); err != nil {
			t.Fatal(err)
		}
	}

	if _, err := s.ApplyServiceRequest(ctx, &spb.ApplyServiceRequestRequest{
		OrgID: oid1,
		Type:  spb.ServiceType_REMITTANCE,
	}); err != nil {
		t.Fatal(err)
	}

	tf1 := spb.ServiceRequest{
		OrgID:       oid1,
		CompanyName: "A",
		Partner:     "TF",
		Type:        spb.ServiceType_REMITTANCE,
		Status:      spb.ServiceRequestStatus_PENDING,
		Enabled:     false,
		Remarks:     "",
		UpdatedBy:   "",
	}
	ir1 := spb.ServiceRequest{
		OrgID:       oid1,
		CompanyName: "A",
		Partner:     "IR",
		Type:        spb.ServiceType_REMITTANCE,
		Status:      spb.ServiceRequestStatus_PENDING,
		Enabled:     false,
		Remarks:     "",
		UpdatedBy:   "",
	}
	wu1 := spb.ServiceRequest{
		OrgID:       oid1,
		CompanyName: "A",
		Partner:     "WU",
		Type:        spb.ServiceType_REMITTANCE,
		Status:      spb.ServiceRequestStatus_PENDING,
		Enabled:     false,
		Remarks:     "",
		UpdatedBy:   "",
	}
	tf2 := spb.ServiceRequest{
		OrgID:       oid2,
		CompanyName: "B",
		Partner:     "TF",
		Type:        spb.ServiceType_REMITTANCE,
		Status:      spb.ServiceRequestStatus_NOSTATUS,
		Enabled:     false,
		Remarks:     "",
		UpdatedBy:   "",
	}
	ir2 := spb.ServiceRequest{
		OrgID:       oid2,
		CompanyName: "B",
		Partner:     "IR",
		Type:        spb.ServiceType_REMITTANCE,
		Status:      spb.ServiceRequestStatus_NOSTATUS,
		Enabled:     false,
		Remarks:     "",
		UpdatedBy:   "",
	}
	wu2 := spb.ServiceRequest{
		OrgID:       oid2,
		CompanyName: "B",
		Partner:     "WU",
		Type:        spb.ServiceType_REMITTANCE,
		Status:      spb.ServiceRequestStatus_NOSTATUS,
		Enabled:     false,
		Remarks:     "",
		UpdatedBy:   "",
	}

	tests := []struct {
		desc string
		want []*spb.ServiceRequest
		in   *spb.ListServiceRequestRequest
	}{
		{
			desc: "All",
			in: &spb.ListServiceRequestRequest{
				OrgIDs: []string{oid1, oid2},
			},
			want: []*spb.ServiceRequest{&tf2, &ir2, &wu2, &tf1, &ir1, &wu1},
		},
		{
			desc: "Sort CreatedCol ASC",
			in: &spb.ListServiceRequestRequest{
				OrgIDs:       []string{oid1, oid2},
				SortByColumn: spb.ServiceSort_CREATED,
				SortBy:       spb.SortBy_ASC,
			},
			want: []*spb.ServiceRequest{&wu1, &ir1, &tf1, &wu2, &ir2, &tf2},
		},
		{
			desc: "Sort CreatedCol DESC",
			in: &spb.ListServiceRequestRequest{
				OrgIDs:       []string{oid1, oid2},
				SortByColumn: spb.ServiceSort_CREATED,
				SortBy:       spb.SortBy_DESC,
			},
			want: []*spb.ServiceRequest{&tf2, &ir2, &wu2, &tf1, &ir1, &wu1},
		},
		{
			desc: "Sort CompanyNameCol ASC",
			in: &spb.ListServiceRequestRequest{
				OrgIDs:       []string{oid1, oid2},
				SortByColumn: spb.ServiceSort_COMPANYNAME,
				SortBy:       spb.SortBy_ASC,
			},
			want: []*spb.ServiceRequest{&ir1, &tf1, &wu1, &wu2, &ir2, &tf2},
		},
		{
			desc: "Sort CompanyNameCol DESC",
			in: &spb.ListServiceRequestRequest{
				OrgIDs:       []string{oid1, oid2},
				SortByColumn: spb.ServiceSort_COMPANYNAME,
				SortBy:       spb.SortBy_DESC,
			},
			want: []*spb.ServiceRequest{&wu2, &ir2, &tf2, &ir1, &tf1, &wu1},
		},
		{
			desc: "Sort ServiceNameCol ASC",
			in: &spb.ListServiceRequestRequest{
				OrgIDs:       []string{oid1, oid2},
				SortByColumn: spb.ServiceSort_SERVICENAME,
				SortBy:       spb.SortBy_ASC,
			},
			want: []*spb.ServiceRequest{&wu2, &ir2, &tf2, &ir1, &tf1, &wu1},
		},
		{
			desc: "Sort ServiceNameCol DESC",
			in: &spb.ListServiceRequestRequest{
				OrgIDs:       []string{oid1, oid2},
				SortByColumn: spb.ServiceSort_SERVICENAME,
				SortBy:       spb.SortBy_DESC,
			},
			want: []*spb.ServiceRequest{&wu2, &ir2, &tf2, &ir1, &tf1, &wu1},
		},
		{
			desc: "Sort StatusCol ASC",
			in: &spb.ListServiceRequestRequest{
				OrgIDs:       []string{oid1, oid2},
				SortByColumn: spb.ServiceSort_STATUS,
				SortBy:       spb.SortBy_ASC,
			},
			want: []*spb.ServiceRequest{&wu2, &ir2, &tf2, &ir1, &tf1, &wu1},
		},
		{
			desc: "Sort StatusCol DESC",
			in: &spb.ListServiceRequestRequest{
				OrgIDs:       []string{oid1, oid2},
				SortByColumn: spb.ServiceSort_STATUS,
				SortBy:       spb.SortBy_DESC,
			},
			want: []*spb.ServiceRequest{&ir1, &tf1, &wu1, &wu2, &ir2, &tf2},
		},
		{
			desc: "Sort PartnerCol ASC",
			in: &spb.ListServiceRequestRequest{
				OrgIDs:       []string{oid1, oid2},
				SortByColumn: spb.ServiceSort_PARTNER,
				SortBy:       spb.SortBy_ASC,
			},
			want: []*spb.ServiceRequest{&ir2, &ir1, &tf2, &tf1, &wu2, &wu1},
		},
		{
			desc: "Sort PartnerCol DESC",
			in: &spb.ListServiceRequestRequest{
				OrgIDs:       []string{oid1, oid2},
				SortByColumn: spb.ServiceSort_PARTNER,
				SortBy:       spb.SortBy_DESC,
			},
			want: []*spb.ServiceRequest{&wu2, &wu1, &tf2, &tf1, &ir2, &ir1},
		},
		{
			desc: "Sort LastUpdatedCol ASC",
			in: &spb.ListServiceRequestRequest{
				OrgIDs:       []string{oid1, oid2},
				SortByColumn: spb.ServiceSort_LASTUPDATED,
				SortBy:       spb.SortBy_ASC,
			},
			want: []*spb.ServiceRequest{&wu2, &ir2, &tf2, &ir1, &tf1, &wu1},
		},
		{
			desc: "Sort LastUpdatedCol DESC",
			in: &spb.ListServiceRequestRequest{
				OrgIDs:       []string{oid1, oid2},
				SortByColumn: spb.ServiceSort_LASTUPDATED,
				SortBy:       spb.SortBy_DESC,
			},
			want: []*spb.ServiceRequest{&ir1, &tf1, &wu1, &tf2, &ir2, &wu2},
		},
		{
			desc: "Limit OrgID",
			in: &spb.ListServiceRequestRequest{
				OrgIDs: []string{oid1},
			},
			want: []*spb.ServiceRequest{&tf1, &ir1, &wu1},
		},
		{
			desc: "Limit Status",
			in: &spb.ListServiceRequestRequest{
				OrgIDs:   []string{oid1, oid2},
				Statuses: []spb.ServiceRequestStatus{spb.ServiceRequestStatus_PENDING},
			},
			want: []*spb.ServiceRequest{&tf1, &ir1, &wu1},
		},
		{
			desc: "Limit SvcName",
			in: &spb.ListServiceRequestRequest{
				OrgIDs: []string{oid1, oid2},
				Types:  []spb.ServiceType{spb.ServiceType_REMITTANCE},
			},
			want: []*spb.ServiceRequest{&tf2, &ir2, &wu2, &tf1, &ir1, &wu1},
		},
		{
			desc: "Limit Partner",
			in: &spb.ListServiceRequestRequest{
				OrgIDs:   []string{oid1, oid2},
				Partners: []string{spb.RemittancePartner_WU.String(), spb.RemittancePartner_IR.String()},
			},
			want: []*spb.ServiceRequest{&ir2, &wu2, &ir1, &wu1},
		},
		{
			desc: "Limit",
			in: &spb.ListServiceRequestRequest{
				OrgIDs: []string{oid1, oid2},
				Limit:  3,
			},
			want: []*spb.ServiceRequest{&tf2, &ir2, &wu2},
		},
		{
			desc: "Offset",
			in: &spb.ListServiceRequestRequest{
				OrgIDs: []string{oid1, oid2},
				Offset: 3,
			},
			want: []*spb.ServiceRequest{&tf1, &ir1, &wu1},
		},
		{
			desc: "Limit/Offset",
			in: &spb.ListServiceRequestRequest{
				OrgIDs: []string{oid1, oid2},
				Limit:  2,
				Offset: 2,
			},
			want: []*spb.ServiceRequest{&wu2, &tf1},
		},
		// todo add tests for LastUpdatedBy once we can approve and rejects svc requests
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			got, err := s.ListServiceRequest(ctx, test.in)
			if err != nil {
				t.Fatal(err)
			}
			o := cmp.Options{
				cmpopts.IgnoreFields(spb.ServiceRequest{}, "Created", "Updated", "Applied", "ID"),
				cmpopts.IgnoreUnexported(spb.ServiceRequest{}),
			}
			if !cmp.Equal(test.want, got.GetServiceRequst(), o) {
				t.Fatal(cmp.Diff(test.want, got.GetServiceRequst(), o))
			}

			if _, err := s.UpdateServiceRequestByOrgID(ctx, &spb.UpdateServiceRequestByOrgIDRequest{
				OldOrgID: "10000000-0000-0000-0000-000000000000",
				NewOrgID: "20000000-0000-0000-0000-000000000000",
				Status:   "Approved",
			}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

var _testStorage *postgres.Storage

func TestMain(m *testing.M) {
	const dbConnEnv = "DATABASE_CONNECTION"
	ddlConnStr := os.Getenv(dbConnEnv)
	if ddlConnStr == "" {
		log.Printf("%s is not set, skipping", dbConnEnv)
		return
	}

	var teardown func()
	_testStorage, teardown = postgres.NewTestStorage(ddlConnStr, filepath.Join("..", "..", "migrations", "sql"))

	exitCode := m.Run()

	if teardown != nil {
		teardown()
	}
	os.Exit(exitCode)
}

func newTestStorage(tb testing.TB) *postgres.Storage {
	if testing.Short() {
		tb.Skip("skipping tests that use postgres on -short")
	}
	return _testStorage
}

func newTestMailer(tb testing.TB) *email.MailSender {
	port, _ := strconv.Atoi(os.Getenv("smtp.port"))
	mailer := email.New(
		os.Getenv("smtp.host"),
		port,
		os.Getenv("smtp.username"),
		os.Getenv("smtp.password"),
		os.Getenv("smtp.fromAddr"),
		os.Getenv("smtp.fromName"),
		os.Getenv("cms.url"),
	)

	return mailer
}

package partner

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus/hooks/test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/profile/storage/postgres"

	ppb "brank.as/petnet/gunk/dsa/v2/partner"
	sc "brank.as/petnet/profile/core/partner"
	tppb "google.golang.org/protobuf/types/known/timestamppb"
)

func TestPartner(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	test.NewNullLogger()
	st, cleanup := postgres.NewTestStorage(os.Getenv("DATABASE_CONNECTION"), filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(cleanup)

	stPend := ppb.PartnerStatusType_PENDING
	stDsbl := ppb.PartnerStatusType_DISABLED
	stEnbl := ppb.PartnerStatusType_ENABLED

	tests := []struct {
		desc    string
		want    *ppb.Partners
		svcType ppb.PartnerType
	}{
		{
			desc: "WU",
			want: &ppb.Partners{
				WesternUnionPartner: &ppb.WesternUnionPartner{
					Status: stPend,
				},
				PartnerStatuses: map[string]string{
					ppb.PartnerType_WU.String(): stPend.String(),
				},
			},
			svcType: ppb.PartnerType_WU,
		},
		{
			desc: "IR",
			want: &ppb.Partners{
				IRemitPartner: &ppb.IRemitPartner{
					Param1: "p-1",
					Param2: "p-2",
					Status: stPend,
				},
				PartnerStatuses: map[string]string{
					ppb.PartnerType_IR.String(): stPend.String(),
				},
			},
			svcType: ppb.PartnerType_IR,
		},
		{
			desc: "TF",
			want: &ppb.Partners{
				TransfastPartner: &ppb.TransfastPartner{
					Param1: "p-1",
					Param2: "p-2",
					Status: stPend,
				},
				PartnerStatuses: map[string]string{
					ppb.PartnerType_TF.String(): stPend.String(),
				},
			},
			svcType: ppb.PartnerType_TF,
		},
		{
			desc: "RM",
			want: &ppb.Partners{
				RemitlyPartner: &ppb.RemitlyPartner{
					Param1: "p-1",
					Param2: "p-2",
					Status: stPend,
				},
				PartnerStatuses: map[string]string{
					ppb.PartnerType_RM.String(): stPend.String(),
				},
			},
			svcType: ppb.PartnerType_RM,
		},
		{
			desc: "RIA",
			want: &ppb.Partners{
				RiaPartner: &ppb.RiaPartner{
					Param1: "p-1",
					Param2: "p-2",
					Status: stPend,
				},
				PartnerStatuses: map[string]string{
					ppb.PartnerType_RIA.String(): stPend.String(),
				},
			},
			svcType: ppb.PartnerType_RIA,
		},
		{
			desc: "MB",
			want: &ppb.Partners{
				MetroBankPartner: &ppb.MetroBankPartner{
					Param1: "p-1",
					Param2: "p-2",
					Status: stPend,
				},
				PartnerStatuses: map[string]string{
					ppb.PartnerType_MB.String(): stPend.String(),
				},
			},
			svcType: ppb.PartnerType_MB,
		},
		{
			desc: "BPI",
			want: &ppb.Partners{
				BPIPartner: &ppb.BPIPartner{
					Param1: "p-1",
					Param2: "p-2",
					Status: stPend,
				},
				PartnerStatuses: map[string]string{
					ppb.PartnerType_BPI.String(): stPend.String(),
				},
			},
			svcType: ppb.PartnerType_BPI,
		},
		{
			desc: "USSC",
			want: &ppb.Partners{
				USSCPartner: &ppb.USSCPartner{
					Param1: "p-1",
					Param2: "p-2",
					Status: stPend,
				},
				PartnerStatuses: map[string]string{
					ppb.PartnerType_USSC.String(): stPend.String(),
				},
			},
			svcType: ppb.PartnerType_USSC,
		},
		{
			desc: "JPR",
			want: &ppb.Partners{
				JapanRemitPartner: &ppb.JapanRemitPartner{
					Param1: "p-1",
					Param2: "p-2",
					Status: stPend,
				},
				PartnerStatuses: map[string]string{
					ppb.PartnerType_JPR.String(): stPend.String(),
				},
			},
			svcType: ppb.PartnerType_JPR,
		},
		{
			desc: "IC",
			want: &ppb.Partners{
				InstantCashPartner: &ppb.InstantCashPartner{
					Param1: "p-1",
					Param2: "p-2",
					Status: stPend,
				},
				PartnerStatuses: map[string]string{
					ppb.PartnerType_IC.String(): stPend.String(),
				},
			},
			svcType: ppb.PartnerType_IC,
		},
		{
			desc: "UNT",
			want: &ppb.Partners{
				UnitellerPartner: &ppb.UnitellerPartner{
					Param1: "p-1",
					Param2: "p-2",
					Status: stPend,
				},
				PartnerStatuses: map[string]string{
					ppb.PartnerType_UNT.String(): stPend.String(),
				},
			},
			svcType: ppb.PartnerType_UNT,
		},
		{
			desc: "CEB",
			want: &ppb.Partners{
				CebuanaPartner: &ppb.CebuanaPartner{
					Param1: "p-1",
					Param2: "p-2",
					Status: stPend,
				},
				PartnerStatuses: map[string]string{
					ppb.PartnerType_CEB.String(): stPend.String(),
				},
			},
			svcType: ppb.PartnerType_CEB,
		},
		{
			desc: "WISE",
			want: &ppb.Partners{
				TransferWisePartner: &ppb.TransferWisePartner{
					Param1: "p-1",
					Param2: "p-2",
					Status: stPend,
				},
				PartnerStatuses: map[string]string{
					ppb.PartnerType_WISE.String(): stPend.String(),
				},
			},
			svcType: ppb.PartnerType_WISE,
		},
		{
			desc: "CEBI",
			want: &ppb.Partners{
				CebuanaIntlPartner: &ppb.CebuanaIntlPartner{
					Param1: "p-1",
					Param2: "p-2",
					Status: stPend,
				},
				PartnerStatuses: map[string]string{
					ppb.PartnerType_CEBI.String(): stPend.String(),
				},
			},
			svcType: ppb.PartnerType_CEBI,
		},
		{
			desc: "AYA",
			want: &ppb.Partners{
				AyannahPartner: &ppb.AyannahPartner{
					Param1: "p-1",
					Param2: "p-2",
					Status: stPend,
				},
				PartnerStatuses: map[string]string{
					ppb.PartnerType_AYA.String(): stPend.String(),
				},
			},
			svcType: ppb.PartnerType_AYA,
		},
		{
			desc: "IE",
			want: &ppb.Partners{
				IntelExpressPartner: &ppb.IntelExpressPartner{
					Param1: "p-1",
					Param2: "p-2",
					Status: stPend,
				},
				PartnerStatuses: map[string]string{
					ppb.PartnerType_IE.String(): stPend.String(),
				},
			},
			svcType: ppb.PartnerType_IE,
		},
	}

	s := New(sc.New(st))
	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			oid := uuid.New().String()
			_, err := st.CreateOrgProfile(ctx, &storage.OrgProfile{
				OrgID:  oid,
				UserID: uuid.NewString(),
			})
			if err != nil {
				t.Fatal(err)
			}
			test.want.OrgID = oid

			got, err := s.CreatePartners(ctx, &ppb.CreatePartnersRequest{Partners: test.want})
			if err != nil {
				t.Fatal("create service: ", err)
			}
			got2, err := s.GetPartners(ctx, &ppb.GetPartnersRequest{OrgID: oid})
			if err != nil {
				t.Fatal("get services: ", err)
			}
			checkTimestamps(t, got2, test.desc)
			checkID(t, got2.GetPartners(), test.desc)
			o := cmp.Options{
				cmpopts.IgnoreFields(
					ppb.WesternUnionPartner{}, "ID", "Created", "Updated", "StartDate", "EndDate",
				),
				cmpopts.IgnoreFields(
					ppb.IRemitPartner{}, "ID", "Created", "Updated", "StartDate", "EndDate",
				),
				cmpopts.IgnoreFields(
					ppb.TransfastPartner{}, "ID", "Created", "Updated", "StartDate", "EndDate",
				),
				cmpopts.IgnoreFields(
					ppb.RemitlyPartner{}, "ID", "Created", "Updated", "StartDate", "EndDate",
				),
				cmpopts.IgnoreFields(
					ppb.RiaPartner{}, "ID", "Created", "Updated", "StartDate", "EndDate",
				),
				cmpopts.IgnoreFields(
					ppb.MetroBankPartner{}, "ID", "Created", "Updated", "StartDate", "EndDate",
				),
				cmpopts.IgnoreFields(
					ppb.BPIPartner{}, "ID", "Created", "Updated", "StartDate", "EndDate",
				),
				cmpopts.IgnoreFields(
					ppb.USSCPartner{}, "ID", "Created", "Updated", "StartDate", "EndDate",
				),
				cmpopts.IgnoreFields(
					ppb.JapanRemitPartner{}, "ID", "Created", "Updated", "StartDate", "EndDate",
				),
				cmpopts.IgnoreFields(
					ppb.InstantCashPartner{}, "ID", "Created", "Updated", "StartDate", "EndDate",
				),
				cmpopts.IgnoreFields(
					ppb.UnitellerPartner{}, "ID", "Created", "Updated", "StartDate", "EndDate",
				),
				cmpopts.IgnoreFields(
					ppb.CebuanaPartner{}, "ID", "Created", "Updated", "StartDate", "EndDate",
				),
				cmpopts.IgnoreFields(
					ppb.TransferWisePartner{}, "ID", "Created", "Updated", "StartDate", "EndDate",
				),
				cmpopts.IgnoreFields(
					ppb.CebuanaIntlPartner{}, "ID", "Created", "Updated", "StartDate", "EndDate",
				),
				cmpopts.IgnoreFields(
					ppb.AyannahPartner{}, "ID", "Created", "Updated", "StartDate", "EndDate",
				),
				cmpopts.IgnoreFields(
					ppb.IntelExpressPartner{}, "ID", "Created", "Updated", "StartDate", "EndDate",
				),
				cmpopts.IgnoreUnexported(
					ppb.Partners{}, ppb.WesternUnionPartner{},
					ppb.IRemitPartner{}, ppb.TransfastPartner{},
					ppb.RemitlyPartner{}, ppb.RiaPartner{},
					ppb.MetroBankPartner{}, ppb.BPIPartner{},
					ppb.USSCPartner{}, ppb.JapanRemitPartner{},
					ppb.InstantCashPartner{}, ppb.UnitellerPartner{},
					ppb.CebuanaPartner{}, ppb.TransferWisePartner{},
					ppb.CebuanaIntlPartner{}, ppb.AyannahPartner{},
					ppb.IntelExpressPartner{},
				),
			}
			if !cmp.Equal(test.want, got.Partners, o) {
				t.Error("(-want +got): ", cmp.Diff(test.want, got.Partners, o))
			}

			id := setUpdateFields(t, got.Partners, test.want, test.desc)
			if _, err := s.UpdatePartners(ctx,
				&ppb.UpdatePartnersRequest{Partners: test.want}); err != nil {
				t.Fatal("update service: ", err)
			}
			got2, err = s.GetPartners(ctx, &ppb.GetPartnersRequest{OrgID: oid})
			if err != nil {
				t.Fatal("get services: ", err)
			}
			checkTimestamps(t, got2, test.desc)
			if !cmp.Equal(test.want, got2.Partners, o) {
				t.Error("(-want +got): ", cmp.Diff(test.want, got2.Partners, o))
			}

			if _, err := s.EnablePartner(ctx, &ppb.EnablePartnerRequest{
				OrgID:   oid,
				Partner: test.svcType,
			}); err != nil {
				t.Fatal("get services: ", err)
			}
			got2, err = s.GetPartners(ctx, &ppb.GetPartnersRequest{OrgID: oid})
			if err != nil {
				t.Fatal("get services: ", err)
			}
			sts := getStatus(t, got2.GetPartners(), test.desc)
			if sts != stEnbl {
				t.Error("want: ENABLED, got: ", sts)
			}
			if len(got2.GetPartners().PartnerStatuses) != 1 {
				t.Error("want: 1")
			}
			ssts := got2.GetPartners().PartnerStatuses[test.svcType.String()]
			if ssts != stEnbl.String() {
				t.Error("want: ENABLED, got: ", ssts)
			}

			if _, err := s.ValidatePartnerAccess(ctx, &ppb.ValidatePartnerAccessRequest{
				OrgID:   oid,
				Partner: test.svcType,
			}); err != nil {
				t.Fatal("service should be accessable: ", err)
			}

			if _, err := s.DisablePartner(ctx, &ppb.DisablePartnerRequest{
				OrgID:   oid,
				Partner: test.svcType,
			}); err != nil {
				t.Fatal("get services: ", err)
			}
			got2, err = s.GetPartners(ctx, &ppb.GetPartnersRequest{OrgID: oid})
			if err != nil {
				t.Fatal("get services: ", err)
			}
			sts = getStatus(t, got2.GetPartners(), test.desc)
			if sts != stDsbl {
				t.Error("want: DISABLED, got: ", sts)
			}
			if len(got2.GetPartners().PartnerStatuses) != 1 {
				t.Error("want: 1")
			}
			ssts = got2.GetPartners().PartnerStatuses[test.svcType.String()]
			if ssts != stDsbl.String() {
				t.Error("want: DISABLED, got: ", ssts)
			}

			if _, err := s.DeletePartner(ctx, &ppb.DeletePartnerRequest{ID: id}); err != nil {
				t.Fatal("delete service: ", err)
			}
			got2, err = s.GetPartners(ctx, &ppb.GetPartnersRequest{OrgID: oid})
			if status.Code(err) != codes.NotFound {
				t.Fatal("get services: ", err)
			}
		})
	}
}

func TestConflictNotFound(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	test.NewNullLogger()
	st, cleanup := postgres.NewTestStorage(os.Getenv("DATABASE_CONNECTION"), filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(cleanup)

	oid := uuid.New().String()
	_, err := st.CreateOrgProfile(ctx, &storage.OrgProfile{
		OrgID:  oid,
		UserID: uuid.NewString(),
	})
	if err != nil {
		t.Fatal(err)
	}

	stPend := ppb.PartnerStatusType_PENDING
	want := []*ppb.Partners{
		{
			OrgID: oid,
			WesternUnionPartner: &ppb.WesternUnionPartner{
				Status:  stPend,
				Created: tppb.Now(),
			},
			PartnerStatuses: map[string]string{
				ppb.PartnerType_WU.String(): stPend.String(),
			},
		},
		{
			OrgID: oid,
			WesternUnionPartner: &ppb.WesternUnionPartner{
				Coy:        "coy2",
				TerminalID: "tid2",
				Status:     stPend,
				Created:    tppb.Now(),
				Updated:    tppb.Now(),
			},
			PartnerStatuses: map[string]string{
				ppb.PartnerType_WU.String(): stPend.String(),
			},
		},
	}

	s := New(sc.New(st))
	_, err = s.GetPartners(ctx, &ppb.GetPartnersRequest{OrgID: uuid.New().String()})
	if status.Code(err) != codes.NotFound {
		t.Fatal("want: notfound, got: ", err)
	}

	got, err := s.CreatePartners(ctx, &ppb.CreatePartnersRequest{Partners: want[0]})
	if err != nil {
		t.Fatal("creating service: ", err)
	}
	_, err = s.CreatePartners(ctx, &ppb.CreatePartnersRequest{Partners: want[0]})
	if status.Code(err) != codes.AlreadyExists {
		t.Fatal("want: alreadyexists, got: ", err)
	}

	want[1].WesternUnionPartner.ID = uuid.New().String()
	if _, err := s.UpdatePartners(ctx, &ppb.UpdatePartnersRequest{Partners: want[1]}); err == nil || grpc.Code(err) != codes.NotFound {
		t.Fatal("should not be found: ", err)
	}
	id := got.Partners.WesternUnionPartner.ID
	want[1].WesternUnionPartner.ID = id
	_, err = s.UpdatePartners(ctx, &ppb.UpdatePartnersRequest{Partners: want[1]})
	if err != nil {
		t.Fatal("update service: ", err)
	}
}

func checkTimestamps(t *testing.T, got *ppb.GetPartnersResponse, svc string) {
	switch svc {
	case "WU":
		if got.Partners.WesternUnionPartner.Created.AsTime().IsZero() {
			t.Fatal("service created timestamp is empty")
		}
		if got.Partners.WesternUnionPartner.Updated.AsTime().IsZero() {
			t.Fatal("service updated timestamp is empty")
		}
	case "IR":
		if got.Partners.IRemitPartner.Created.AsTime().IsZero() {
			t.Fatal("service created timestamp is empty")
		}
		if got.Partners.IRemitPartner.Updated.AsTime().IsZero() {
			t.Fatal("service updated timestamp is empty")
		}
	case "TF":
		if got.Partners.TransfastPartner.Created.AsTime().IsZero() {
			t.Fatal("service created timestamp is empty")
		}
		if got.Partners.TransfastPartner.Updated.AsTime().IsZero() {
			t.Fatal("service updated timestamp is empty")
		}
	case "RM":
		if got.Partners.RemitlyPartner.Created.AsTime().IsZero() {
			t.Fatal("service created timestamp is empty")
		}
		if got.Partners.RemitlyPartner.Updated.AsTime().IsZero() {
			t.Fatal("service updated timestamp is empty")
		}
	case "RIA":
		if got.Partners.RiaPartner.Created.AsTime().IsZero() {
			t.Fatal("service created timestamp is empty")
		}
		if got.Partners.RiaPartner.Updated.AsTime().IsZero() {
			t.Fatal("service updated timestamp is empty")
		}
	case "MB":
		if got.Partners.MetroBankPartner.Created.AsTime().IsZero() {
			t.Fatal("service created timestamp is empty")
		}
		if got.Partners.MetroBankPartner.Updated.AsTime().IsZero() {
			t.Fatal("service updated timestamp is empty")
		}
	case "BPI":
		if got.Partners.BPIPartner.Created.AsTime().IsZero() {
			t.Fatal("service created timestamp is empty")
		}
		if got.Partners.BPIPartner.Updated.AsTime().IsZero() {
			t.Fatal("service updated timestamp is empty")
		}
	case "USSC":
		if got.Partners.USSCPartner.Created.AsTime().IsZero() {
			t.Fatal("service created timestamp is empty")
		}
		if got.Partners.USSCPartner.Updated.AsTime().IsZero() {
			t.Fatal("service updated timestamp is empty")
		}
	case "JPR":
		if got.Partners.JapanRemitPartner.Created.AsTime().IsZero() {
			t.Fatal("service created timestamp is empty")
		}
		if got.Partners.JapanRemitPartner.Updated.AsTime().IsZero() {
			t.Fatal("service updated timestamp is empty")
		}
	case "IC":
		if got.Partners.InstantCashPartner.Created.AsTime().IsZero() {
			t.Fatal("service created timestamp is empty")
		}
		if got.Partners.InstantCashPartner.Updated.AsTime().IsZero() {
			t.Fatal("service updated timestamp is empty")
		}
	case "UNT":
		if got.Partners.UnitellerPartner.Created.AsTime().IsZero() {
			t.Fatal("service created timestamp is empty")
		}
		if got.Partners.UnitellerPartner.Updated.AsTime().IsZero() {
			t.Fatal("service updated timestamp is empty")
		}
	case "CEB":
		if got.Partners.CebuanaPartner.Created.AsTime().IsZero() {
			t.Fatal("service created timestamp is empty")
		}
		if got.Partners.CebuanaPartner.Updated.AsTime().IsZero() {
			t.Fatal("service updated timestamp is empty")
		}
	case "WISE":
		if got.Partners.TransferWisePartner.Created.AsTime().IsZero() {
			t.Fatal("service created timestamp is empty")
		}
		if got.Partners.TransferWisePartner.Updated.AsTime().IsZero() {
			t.Fatal("service updated timestamp is empty")
		}
	case "CEBI":
		if got.Partners.CebuanaIntlPartner.Created.AsTime().IsZero() {
			t.Fatal("service created timestamp is empty")
		}
		if got.Partners.CebuanaIntlPartner.Updated.AsTime().IsZero() {
			t.Fatal("service updated timestamp is empty")
		}
	case "AYA":
		if got.Partners.AyannahPartner.Created.AsTime().IsZero() {
			t.Fatal("service created timestamp is empty")
		}
		if got.Partners.AyannahPartner.Updated.AsTime().IsZero() {
			t.Fatal("service updated timestamp is empty")
		}
	case "IE":
		if got.Partners.IntelExpressPartner.Created.AsTime().IsZero() {
			t.Fatal("service created timestamp is empty")
		}
		if got.Partners.IntelExpressPartner.Updated.AsTime().IsZero() {
			t.Fatal("service updated timestamp is empty")
		}
	}
}

func checkID(t *testing.T, got *ppb.Partners, svc string) {
	switch svc {
	case "WU":
		if got.WesternUnionPartner.ID == "" {
			t.Error("id is empty")
		}
	case "IR":
		if got.IRemitPartner.ID == "" {
			t.Error("id is empty")
		}
	case "TF":
		if got.TransfastPartner.ID == "" {
			t.Error("id is empty")
		}
	case "RM":
		if got.RemitlyPartner.ID == "" {
			t.Error("id is empty")
		}
	case "RIA":
		if got.RiaPartner.ID == "" {
			t.Error("id is empty")
		}
	case "MB":
		if got.MetroBankPartner.ID == "" {
			t.Error("id is empty")
		}
	case "BPI":
		if got.BPIPartner.ID == "" {
			t.Error("id is empty")
		}
	case "USSC":
		if got.USSCPartner.ID == "" {
			t.Error("id is empty")
		}
	case "JPR":
		if got.JapanRemitPartner.ID == "" {
			t.Error("id is empty")
		}
	case "IC":
		if got.InstantCashPartner.ID == "" {
			t.Error("id is empty")
		}
	case "UNT":
		if got.UnitellerPartner.ID == "" {
			t.Error("id is empty")
		}
	case "CEB":
		if got.CebuanaPartner.ID == "" {
			t.Error("id is empty")
		}
	case "WISE":
		if got.TransferWisePartner.ID == "" {
			t.Error("id is empty")
		}
	case "CEBI":
		if got.CebuanaIntlPartner.ID == "" {
			t.Error("id is empty")
		}
	case "AYA":
		if got.AyannahPartner.ID == "" {
			t.Error("id is empty")
		}
	case "IE":
		if got.IntelExpressPartner.ID == "" {
			t.Error("id is empty")
		}
	}
}

func getStatus(t *testing.T, got *ppb.Partners, svc string) ppb.PartnerStatusType {
	var sts ppb.PartnerStatusType
	switch svc {
	case "WU":
		sts = got.GetWesternUnionPartner().GetStatus()
	case "IR":
		sts = got.GetIRemitPartner().GetStatus()
	case "TF":
		sts = got.GetTransfastPartner().GetStatus()
	case "RM":
		sts = got.GetRemitlyPartner().GetStatus()
	case "RIA":
		sts = got.GetRiaPartner().GetStatus()
	case "MB":
		sts = got.GetMetroBankPartner().GetStatus()
	case "BPI":
		sts = got.GetBPIPartner().GetStatus()
	case "USSC":
		sts = got.GetUSSCPartner().GetStatus()
	case "JPR":
		sts = got.GetJapanRemitPartner().GetStatus()
	case "IC":
		sts = got.GetInstantCashPartner().GetStatus()
	case "UNT":
		sts = got.GetUnitellerPartner().GetStatus()
	case "CEB":
		sts = got.GetCebuanaPartner().GetStatus()
	case "WISE":
		sts = got.GetTransferWisePartner().GetStatus()
	case "CEBI":
		sts = got.GetCebuanaIntlPartner().GetStatus()
	case "AYA":
		sts = got.GetAyannahPartner().GetStatus()
	case "IE":
		sts = got.GetIntelExpressPartner().GetStatus()
	}
	return sts
}

func setUpdateFields(t *testing.T, got *ppb.Partners, want *ppb.Partners, svc string) string {
	var id string
	switch svc {
	case "WU":
		id = got.WesternUnionPartner.ID
		want.WesternUnionPartner.ID = id
		want.WesternUnionPartner.Coy = "coy"
		want.WesternUnionPartner.TerminalID = "tid"
	case "IR":
		id = got.IRemitPartner.ID
		want.IRemitPartner.ID = id
	case "TF":
		id = got.TransfastPartner.ID
		want.TransfastPartner.ID = id
	case "RM":
		id = got.RemitlyPartner.ID
		want.RemitlyPartner.ID = id
	case "RIA":
		id = got.RiaPartner.ID
		want.RiaPartner.ID = id
	case "MB":
		id = got.MetroBankPartner.ID
		want.MetroBankPartner.ID = id
	case "BPI":
		id = got.BPIPartner.ID
		want.BPIPartner.ID = id
	case "USSC":
		id = got.USSCPartner.ID
		want.USSCPartner.ID = id
	case "JPR":
		id = got.JapanRemitPartner.ID
		want.JapanRemitPartner.ID = id
	case "IC":
		id = got.InstantCashPartner.ID
		want.InstantCashPartner.ID = id
	case "UNT":
		id = got.UnitellerPartner.ID
		want.UnitellerPartner.ID = id
	case "CEB":
		id = got.CebuanaPartner.ID
		want.CebuanaPartner.ID = id
	case "WISE":
		id = got.TransferWisePartner.ID
		want.TransferWisePartner.ID = id
	case "CEBI":
		id = got.CebuanaIntlPartner.ID
		want.CebuanaIntlPartner.ID = id
	case "AYA":
		id = got.AyannahPartner.ID
		want.AyannahPartner.ID = id
	case "IE":
		id = got.IntelExpressPartner.ID
		want.IntelExpressPartner.ID = id
	}
	return id
}

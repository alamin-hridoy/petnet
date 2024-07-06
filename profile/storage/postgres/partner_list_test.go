package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	spb "brank.as/petnet/gunk/dsa/v2/partner"
	"brank.as/petnet/profile/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestPartnerList(t *testing.T) {
	ts := newTestStorage(t)
	ctx := context.Background()
	tests := []struct {
		desc          string
		input         []*storage.PartnerList
		payload       interface{}
		payloadUpdate interface{}
	}{
		{
			desc: "partnerList",
			input: []*storage.PartnerList{
				{
					Stype:       spb.PartnerType_WU.String(),
					Name:        "Weatern Union",
					Created:     time.Now(),
					Updated:     time.Now(),
					Status:      spb.PartnerStatusType_ENABLED.String(),
					ServiceName: "REMITTANCE",
					DisableReason: sql.NullString{
						String: "Profile not completed",
						Valid:  true,
					},
					UpdatedBy:        "394d67ad-5e91-4f10-af1b-3305d583639a",
					Platform:         "Perahub",
					IsProvider:       false,
					PerahubPartnerID: "123",
					RemcoID:          "1",
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			for _, pnrList := range test.input {
				_, err := ts.CreatePartnerList(ctx, pnrList)
				if err != nil {
					t.Fatal(err)
				}
			}
			res, err := ts.GetPartnerList(ctx, &storage.PartnerList{
				ServiceName: "REMITTANCE",
			})
			if err != nil {
				t.Fatal(err)
			}
			if res[0].Created.IsZero() {
				t.Error("created is empty")
			}
			if res[0].Updated.IsZero() {
				t.Error("updated is empty")
			}
			_, err = ts.GetPartnerByStype(ctx, spb.PartnerType_WU.String())
			if err != nil {
				t.Fatal(err)
			}
			o := cmp.Options{
				cmpopts.IgnoreFields(
					storage.PartnerList{}, "ID", "Created", "Updated", "Deleted",
				),
			}
			if !cmp.Equal(res, test.input, o) {
				t.Error(cmp.Diff(res, test.input, o))
			}
			_, err = ts.GetDSAPartnerList(ctx, &storage.GetDSAPartnerListRequest{})
			if err != nil {
				t.Fatal(err)
			}

			test.input[0].Status = spb.PartnerStatusType_DISABLED.String()
			_, err = ts.UpdatePartnerList(ctx, test.input[0])
			if err != nil {
				t.Error("update failed")
			}
			Updatedres, err := ts.GetPartnerList(ctx, &storage.PartnerList{
				ServiceName: "REMITTANCE",
			})
			if err != nil {
				t.Fatal(err)
			}
			if Updatedres[0].Status != test.input[0].Status {
				t.Error("should be disabled")
			}
			styps := []string{}
			disableReason := sql.NullString{}
			updatedBy := ""
			for _, us := range Updatedres {
				disableReason = us.DisableReason
				styps = append(styps, us.Stype)
				_, err = ts.EnablePartnerList(ctx, us.Stype)
				if err != nil {
					t.Fatal(err)
				}
			}
			_, err = ts.DisableMultiplePartnerList(ctx, styps, disableReason.String, updatedBy)
			if err != nil {
				t.Fatal(err)
			}
			_, err = ts.EnableMultiplePartnerList(ctx, styps, updatedBy)
			if err != nil {
				t.Fatal(err)
			}
			Enabledres, err := ts.GetPartnerList(ctx, &storage.PartnerList{})
			if err != nil {
				t.Fatal(err)
			}
			for _, us := range Enabledres {
				if us.Status == spb.PartnerStatusType_DISABLED.String() {
					t.Error("should be enabled")
				}
			}
			for _, us := range Enabledres {
				_, err = ts.DisablePartnerList(ctx, us.Stype, us.DisableReason.String)
				if err != nil {
					t.Fatal(err)
				}
			}
			Disabledres, err := ts.GetPartnerList(ctx, &storage.PartnerList{
				ServiceName: "REMITTANCE",
			})
			if err != nil {
				t.Fatal(err)
			}
			for _, ds := range Disabledres {
				if ds.Status == spb.PartnerStatusType_ENABLED.String() {
					t.Error("should be Disabled")
				}
			}
			for _, ds := range Disabledres {
				_, err = ts.DeletePartnerList(ctx, ds.Stype)
				if err != nil {
					t.Fatal(err)
				}
			}
			deletedRes, err := ts.GetPartnerList(ctx, &storage.PartnerList{
				ServiceName: "REMITTANCE",
			})
			if err != nil && err != storage.NotFound {
				t.Fatal(err)
			}
			if deletedRes != nil {
				t.Error("should be Deleted")
			}
		})
	}
}

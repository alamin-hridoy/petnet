package postgres_test

import (
	"context"
	"testing"
	"time"

	spb "brank.as/petnet/gunk/dsa/v2/partner"
	"brank.as/petnet/profile/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCICOPartnerList(t *testing.T) {
	ts := newTestStorage(t)
	ctx := context.Background()
	tests := []struct {
		desc          string
		input         []*storage.CICOPartnerList
		payload       interface{}
		payloadUpdate interface{}
	}{
		{
			desc: "cicoPartnerList",
			input: []*storage.CICOPartnerList{
				{
					Stype:   spb.PartnerType_WU.String(),
					Name:    "Weatern Union",
					Created: time.Now(),
					Updated: time.Now(),
					Status:  spb.PartnerStatusType_ENABLED.String(),
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			for _, pnrList := range test.input {
				_, err := ts.CreateCICOPartnerList(ctx, pnrList)
				if err != nil {
					t.Fatal(err)
				}
			}
			res, err := ts.GetCICOPartnerList(ctx, &storage.CICOPartnerList{})
			if err != nil {
				t.Fatal(err)
			}
			if res[0].Created.IsZero() {
				t.Error("created is empty")
			}
			if res[0].Updated.IsZero() {
				t.Error("updated is empty")
			}
			o := cmp.Options{
				cmpopts.IgnoreFields(
					storage.PartnerList{}, "ID", "Created", "Updated", "Deleted",
				),
			}
			if !cmp.Equal(res, test.input, o) {
				t.Error(cmp.Diff(res, test.input, o))
			}

			test.input[0].Status = spb.PartnerStatusType_DISABLED.String()
			_, err = ts.UpdateCICOPartnerList(ctx, test.input[0])
			if err != nil {
				t.Error("update failed")
			}
			Updatedres, err := ts.GetCICOPartnerList(ctx, &storage.CICOPartnerList{})
			if err != nil {
				t.Fatal(err)
			}
			if Updatedres[0].Status != test.input[0].Status {
				t.Error("should be disabled")
			}
			for _, us := range Updatedres {
				_, err = ts.EnableCICOPartnerList(ctx, us.Stype)
				if err != nil {
					t.Fatal(err)
				}
			}
			Enabledres, err := ts.GetCICOPartnerList(ctx, &storage.CICOPartnerList{})
			if err != nil {
				t.Fatal(err)
			}
			for _, us := range Enabledres {
				if us.Status == spb.PartnerStatusType_DISABLED.String() {
					t.Error("should be enabled")
				}
			}
			for _, us := range Enabledres {
				_, err = ts.DisableCICOPartnerList(ctx, us.Stype)
				if err != nil {
					t.Fatal(err)
				}
			}
			Disabledres, err := ts.GetCICOPartnerList(ctx, &storage.CICOPartnerList{})
			if err != nil {
				t.Fatal(err)
			}
			for _, ds := range Disabledres {
				if ds.Status == spb.PartnerStatusType_ENABLED.String() {
					t.Error("should be Disabled")
				}
			}
			for _, ds := range Disabledres {
				_, err = ts.DeleteCICOPartnerList(ctx, ds.Stype)
				if err != nil {
					t.Fatal(err)
				}
			}
			deletedRes, err := ts.GetCICOPartnerList(ctx, &storage.CICOPartnerList{})
			if err != nil && err != storage.NotFound {
				t.Fatal(err)
			}
			if deletedRes != nil {
				t.Error("should be Deleted")
			}
		})
	}
}

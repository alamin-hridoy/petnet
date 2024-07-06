package postgres_test

import (
	"context"
	"encoding/json"
	"testing"

	spb "brank.as/petnet/gunk/dsa/v2/partner"
	"brank.as/petnet/profile/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestPartner(t *testing.T) {
	ts := newTestStorage(t)

	ctx := context.Background()
	oid := uuid.NewString()
	uid := uuid.NewString()
	if _, err := ts.CreateOrgProfile(ctx, &storage.OrgProfile{
		OrgID:  oid,
		UserID: uid,
	}); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		desc          string
		input         *storage.Partner
		payload       interface{}
		payloadUpdate interface{}
	}{
		{
			desc: "WesternUnion",
			input: &storage.Partner{
				Type:      spb.PartnerType_WU.String(),
				OrgID:     oid,
				UpdatedBy: uid,
				Status:    spb.PartnerStatusType_PENDING.String(),
			},
			payload: storage.WesternUnionPartner{
				Coy:        "coy",
				TerminalID: "terminal-id",
			},
			payloadUpdate: storage.WesternUnionPartner{
				Coy:        "coy-u",
				TerminalID: "terminal-id-u",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			b, err := json.Marshal(&test.payload)
			if err != nil {
				t.Fatal(err)
			}
			test.input.Partner = string(b)

			id, err := ts.CreatePartner(ctx, test.input)
			if err != nil {
				t.Fatal(err)
			}
			gs, err := ts.GetPartner(ctx, oid, test.input.Type)
			if err != nil {
				t.Fatal(err)
			}
			if gs.ID != id {
				t.Fatal("Get Partner")
			}
			_, err = ts.CreatePartner(ctx, test.input)
			if err != storage.Conflict {
				t.Error("want: conflict, got: ", err)
			}
			if id == "" {
				t.Error("id is empty")
			}
			res, err := ts.GetPartners(ctx, oid)
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
					storage.Partner{}, "ID", "Created", "Updated", "Partner",
				),
			}
			test.input.Status = spb.PartnerStatusType_PENDING.String()
			if !cmp.Equal(res, []storage.Partner{*test.input}, o) {
				t.Error(cmp.Diff(res, []storage.Partner{*test.input}, o))
			}

			switch test.input.Type {
			case spb.PartnerType_WU.String():
				sv := &storage.WesternUnionPartner{}
				err = json.Unmarshal([]byte(res[0].Partner), sv)
				if err != nil {
					t.Fatal(err)
				}
				if !cmp.Equal(test.payload, *sv) {
					t.Error(cmp.Diff(test.payload, *sv))
				}
			}

			test.input.ID = id
			b, err = json.Marshal(&test.payloadUpdate)
			if err != nil {
				t.Fatal(err)
			}
			test.input.Partner = string(b)
			sid := test.input.ID
			test.input.ID = uuid.New().String()
			id, err = ts.UpdatePartner(ctx, test.input)
			if err != storage.NotFound {
				t.Error("want: notfound, got: ", err)
			}
			test.input.ID = sid
			id, err = ts.UpdatePartner(ctx, test.input)
			if err != nil {
				t.Fatal(err)
			}
			if id == "" {
				t.Error("id is empty")
			}
			res2, err := ts.GetPartners(ctx, oid)
			if err != nil {
				t.Fatal(err)
			}
			if len(res2) != 1 {
				t.Error("want: 1, got: ", len(res2))
			}
			if res2[0].Created.IsZero() {
				t.Error("created is empty")
			}
			if res2[0].Updated.IsZero() {
				t.Error("updated is empty")
			}
			if res2[0].Created != res[0].Created {
				t.Error("created has been changed")
			}
			if res2[0].Updated == res[0].Updated {
				t.Error("updated not updated")
			}

			switch test.input.Type {
			case spb.PartnerType_WU.String():
				sv := &storage.WesternUnionPartner{}
				err = json.Unmarshal([]byte(res2[0].Partner), sv)
				if err != nil {
					t.Fatal(err)
				}
				if !cmp.Equal(test.payloadUpdate, *sv) {
					t.Error(cmp.Diff(test.payloadUpdate, *sv))
				}
			}

			if err := ts.ValidatePartnerAccess(ctx, oid, test.input.Type); err == nil {
				t.Fatal("partner should be disabled and invalid")
			}

			err = ts.EnablePartner(ctx, oid, test.input.Type)
			if err != nil {
				t.Fatal(err)
			}
			res, err = ts.GetPartners(ctx, oid)
			if err != nil {
				t.Fatal(err)
			}
			if res[0].Status != spb.PartnerStatusType_ENABLED.String() {
				t.Error("should be enabled")
			}

			if err := ts.ValidatePartnerAccess(ctx, oid, test.input.Type); err != nil {
				t.Fatal("partner should be enabled and valid")
			}

			err = ts.DisablePartner(ctx, oid, test.input.Type)
			if err != nil {
				t.Fatal(err)
			}
			res, err = ts.GetPartners(ctx, oid)
			if err != nil {
				t.Fatal(err)
			}
			if res[0].Status != spb.PartnerStatusType_DISABLED.String() {
				t.Error("should be disabled")
			}

			if err := ts.ValidatePartnerAccess(ctx, oid, test.input.Type); err == nil {
				t.Fatal("partner should be disabled and invalid")
			}

			id, err = ts.DeletePartner(ctx, id)
			if err != nil {
				t.Fatal(err)
			}
			if id == "" {
				t.Error("id is empty")
			}
			res, err = ts.GetPartners(ctx, oid)
			if err != nil {
				t.Fatal(err)
			}
			if len(res) != 0 {
				t.Error("should be 0 after deletion")
			}
		})
	}
}

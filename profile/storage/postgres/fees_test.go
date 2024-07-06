package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"brank.as/petnet/profile/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestCreateFees(t *testing.T) {
	ts := newTestStorage(t)

	ctx := context.Background()
	oid := uuid.NewString()
	_, err := ts.CreateOrgProfile(ctx, &storage.OrgProfile{
		OrgID:  oid,
		UserID: uuid.NewString(),
	})
	if err != nil {
		t.Fatal(err)
	}
	nt := sql.NullTime{Time: time.Unix(1414141414, 0), Valid: true}
	nt2 := sql.NullTime{Time: time.Unix(1515151515, 0), Valid: true}

	o := cmpopts.IgnoreFields(storage.FeeCommission{}, "ID", "Type", "Created", "Updated")
	fs := []storage.FeeCommission{
		{
			OrgID:            oid,
			OrgProfileID:     uuid.NewString(),
			FeeAmount:        "100",
			CommissionAmount: "100",
			StartDate:        nt,
			EndDate:          nt,
		},
		{
			OrgID:            oid,
			OrgProfileID:     uuid.NewString(),
			FeeAmount:        "200",
			CommissionAmount: "200",
			StartDate:        nt2,
			EndDate:          nt2,
		},
	}
	got, err := ts.CreateOrgFees(ctx, fs[0])
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(fs[0], *got, o) {
		t.Error("(-want +got): ", cmp.Diff(fs[0], *got, o))
	}
	got, err = ts.UpsertOrgFees(ctx, fs[1])
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(fs[1], *got, o) {
		t.Error("(-want +got): ", cmp.Diff(fs[1], *got, o))
	}
}

func TestListFees(t *testing.T) {
	ts := newTestStorage(t)

	ctx := context.Background()
	oid := uuid.NewString()
	_, err := ts.CreateOrgProfile(ctx, &storage.OrgProfile{
		OrgID:  oid,
		UserID: uuid.NewString(),
	})
	if err != nil {
		t.Fatal(err)
	}
	nt := sql.NullTime{Time: time.Unix(1414141414, 0), Valid: true}

	o := cmpopts.IgnoreFields(storage.FeeCommission{}, "ID", "Type", "Created", "Updated")
	fs := []storage.FeeCommission{
		{
			OrgID:            oid,
			OrgProfileID:     uuid.NewString(),
			FeeAmount:        "1",
			CommissionAmount: "1",
			FeeStatus:        2,
			StartDate:        nt,
			EndDate:          nt,
			Count:            5,
		},
		{
			OrgID:            oid,
			OrgProfileID:     uuid.NewString(),
			FeeAmount:        "2",
			CommissionAmount: "2",
			FeeStatus:        2,
			StartDate:        nt,
			EndDate:          nt,
			Count:            5,
		},
		{
			OrgID:            oid,
			OrgProfileID:     uuid.NewString(),
			FeeAmount:        "3",
			CommissionAmount: "3",
			FeeStatus:        2,
			StartDate:        nt,
			EndDate:          nt,
			Count:            5,
		},
		{
			OrgID:            oid,
			OrgProfileID:     uuid.NewString(),
			FeeAmount:        "4",
			CommissionAmount: "4",
			FeeStatus:        2,
			StartDate:        nt,
			EndDate:          nt,
			Count:            5,
		},
		{
			OrgID:            oid,
			OrgProfileID:     uuid.NewString(),
			FeeAmount:        "5",
			CommissionAmount: "5",
			FeeStatus:        2,
			StartDate:        nt,
			EndDate:          nt,
			Count:            5,
		},
	}
	for _, f := range fs {
		if _, err := ts.CreateOrgFees(ctx, f); err != nil {
			t.Error(err)
		}
	}

	tests := []struct {
		name string
		f    storage.LimitOffsetFilter
		want []storage.FeeCommission
	}{
		{
			name: "No Limit",
			f:    storage.LimitOffsetFilter{},
			want: []storage.FeeCommission{fs[0], fs[1], fs[2], fs[3], fs[4]},
		},
		{
			name: "First Two",
			f: storage.LimitOffsetFilter{
				Limit:  2,
				Offset: 0,
			},
			want: []storage.FeeCommission{fs[0], fs[1]},
		},
		{
			name: "Next Two",
			f: storage.LimitOffsetFilter{
				Limit:  2,
				Offset: 2,
			},
			want: []storage.FeeCommission{fs[2], fs[3]},
		},
		{
			name: "Last One",
			f: storage.LimitOffsetFilter{
				Limit:  2,
				Offset: 4,
			},
			want: []storage.FeeCommission{fs[4]},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := ts.ListOrgFees(ctx, oid, test.f)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.want, got, o) {
				t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
			}
		})
	}
}

func TestCreateFeeCommissionRate(t *testing.T) {
	ts := newTestStorage(t)

	ctx := context.Background()
	fcid := uuid.NewString()

	o := cmpopts.IgnoreFields(storage.Rate{}, "ID")
	fs := []storage.Rate{
		{
			FeeCommissionID: fcid,
			MinVolume:       "50",
			MaxVolume:       "100",
			TxnRate:         "50%",
		},
	}
	got, err := ts.CreateFeeCommissionRate(ctx, fs[0])
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(fs[0], *got, o) {
		t.Error("(-want +got): ", cmp.Diff(fs[0], *got, o))
	}
}

func TestUpsertRate(t *testing.T) {
	ts := newTestStorage(t)

	ctx := context.Background()
	fcid := uuid.NewString()

	o := cmpopts.IgnoreFields(storage.Rate{}, "ID")
	fs := []storage.Rate{
		{
			FeeCommissionID: fcid,
			MinVolume:       "50",
			MaxVolume:       "100",
			TxnRate:         "50%",
		},
	}
	got, err := ts.UpsertRate(ctx, fs[0])
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(fs[0], *got, o) {
		t.Error("(-want +got): ", cmp.Diff(fs[0], *got, o))
	}
}

func TestListFeesCommissionRate(t *testing.T) {
	ts := newTestStorage(t)
	ctx := context.Background()

	fcid := uuid.NewString()
	o := cmpopts.IgnoreFields(storage.Rate{}, "ID")

	fs := []storage.Rate{
		{
			FeeCommissionID: fcid,
			MinVolume:       "50",
			MaxVolume:       "100",
			TxnRate:         "50%",
		},
		{
			FeeCommissionID: fcid,
			MinVolume:       "70",
			MaxVolume:       "160",
			TxnRate:         "20%",
		},
		{
			FeeCommissionID: fcid,
			MinVolume:       "10",
			MaxVolume:       "90",
			TxnRate:         "10%",
		},
	}
	for _, f := range fs {
		if _, err := ts.CreateFeeCommissionRate(ctx, f); err != nil {
			t.Error(err)
		}
	}

	tests := []struct {
		name string
		want []storage.Rate
	}{
		{
			name: "Get List",
			want: []storage.Rate{fs[0], fs[1], fs[2]},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := ts.ListFeesCommissionRate(ctx, fcid)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.want, got, o) {
				t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
			}
		})
	}
}

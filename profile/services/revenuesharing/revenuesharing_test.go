package revenuesharing

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	spb "brank.as/petnet/gunk/dsa/v2/partner"
	rcm "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	rcc "brank.as/petnet/profile/core/revenuesharing"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/profile/storage/postgres"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus/hooks/test"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestRevenueSharing(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	test.NewNullLogger()
	st := newTestStorage(t)
	s := New(rcc.New(st))

	uid := uuid.NewString()
	oid := uuid.New().String()

	tests := []struct {
		desc             string
		input            rcm.CreateRevenueSharingRequest
		inputRes         *rcm.CreateRevenueSharingResponse
		inputtier        []rcm.CreateRevenueSharingTierRequest
		inputtierRes     []*rcm.RevenueSharingTier
		payload          storage.RevenueSharing
		payloadUpdate    rcm.UpdateRevenueSharingRequest
		payloadGet       rcm.GetRevenueSharingListRequest
		payloadUpdateRes *rcm.UpdateRevenueSharingResponse
		payloadGetRes    []*rcm.RevenueSharing
	}{
		{
			desc: "insert remittance revenue sharing",
			input: rcm.CreateRevenueSharingRequest{
				OrgID:           oid,
				UserID:          uid,
				Partner:         spb.PartnerType_AYA.String(),
				BoundType:       rcm.BoundType_INBOUND,
				RemitType:       rcm.RemitType_REMITTANCE,
				TransactionType: rcm.TransactionType_DIGITAL,
				TierType:        rcm.TierType_PERCENTAGE,
				Amount:          "10",
				CreatedBy:       uid,
			},
			inputRes: &rcm.CreateRevenueSharingResponse{
				OrgID:           oid,
				UserID:          uid,
				Partner:         spb.PartnerType_AYA.String(),
				BoundType:       rcm.BoundType_INBOUND,
				RemitType:       rcm.RemitType_REMITTANCE,
				TransactionType: rcm.TransactionType_DIGITAL,
				TierType:        rcm.TierType_PERCENTAGE,
				Amount:          "10",
				CreatedBy:       uid,
			},
			payloadUpdate: rcm.UpdateRevenueSharingRequest{
				OrgID:           oid,
				UserID:          uid,
				Partner:         spb.PartnerType_AYA.String(),
				BoundType:       rcm.BoundType_INBOUND,
				RemitType:       rcm.RemitType_REMITTANCE,
				TransactionType: rcm.TransactionType_DIGITAL,
				TierType:        rcm.TierType_TIERPERCENTAGE,
				Amount:          "",
				CreatedBy:       uid,
			},
			payloadUpdateRes: &rcm.UpdateRevenueSharingResponse{
				Partner:         spb.PartnerType_AYA.String(),
				BoundType:       rcm.BoundType_INBOUND,
				RemitType:       rcm.RemitType_REMITTANCE,
				TransactionType: rcm.TransactionType_DIGITAL,
				TierType:        rcm.TierType_TIERPERCENTAGE,
				Amount:          "",
				CreatedBy:       uid,
			},
			inputtier: []rcm.CreateRevenueSharingTierRequest{
				{
					MinValue: "1000",
					MaxValue: "2000",
					Amount:   "10",
				},
				{
					MinValue: "2001",
					MaxValue: "3000",
					Amount:   "20",
				},
				{
					MinValue: "3001",
					MaxValue: "4000",
					Amount:   "30",
				},
			},
			inputtierRes: []*rcm.RevenueSharingTier{
				{
					MinValue: "1000",
					MaxValue: "2000",
					Amount:   "10",
				},
				{
					MinValue: "2001",
					MaxValue: "3000",
					Amount:   "20",
				},
				{
					MinValue: "3001",
					MaxValue: "4000",
					Amount:   "30",
				},
			},
			payloadGet: rcm.GetRevenueSharingListRequest{
				OrgID:           oid,
				UserID:          uid,
				Partner:         spb.PartnerType_AYA.String(),
				BoundType:       rcm.BoundType_INBOUND,
				RemitType:       rcm.RemitType_REMITTANCE,
				TransactionType: rcm.TransactionType_DIGITAL,
				TierType:        rcm.TierType_TIERPERCENTAGE,
				Amount:          "",
				CreatedBy:       uid,
			},
			payloadGetRes: []*rcm.RevenueSharing{
				{
					OrgID:           oid,
					UserID:          uid,
					Partner:         spb.PartnerType_AYA.String(),
					BoundType:       rcm.BoundType_INBOUND,
					RemitType:       rcm.RemitType_REMITTANCE,
					TransactionType: rcm.TransactionType_DIGITAL,
					TierType:        rcm.TierType_TIERPERCENTAGE,
					Amount:          "",
					CreatedBy:       uid,
				},
			},
		},
		{
			desc: "insert bills payment Partner commission",
			input: rcm.CreateRevenueSharingRequest{
				OrgID:           oid,
				UserID:          uid,
				Partner:         spb.PartnerType_WU.String(),
				BoundType:       rcm.BoundType_INBOUND,
				RemitType:       rcm.RemitType_BILLSPAYMENT,
				TransactionType: rcm.TransactionType_OTC,
				TierType:        rcm.TierType_PERCENTAGE,
				Amount:          "20",
				CreatedBy:       uid,
			},
			inputRes: &rcm.CreateRevenueSharingResponse{
				OrgID:           oid,
				UserID:          uid,
				Partner:         spb.PartnerType_WU.String(),
				BoundType:       rcm.BoundType_INBOUND,
				RemitType:       rcm.RemitType_BILLSPAYMENT,
				TransactionType: rcm.TransactionType_OTC,
				TierType:        rcm.TierType_PERCENTAGE,
				Amount:          "20",
				CreatedBy:       uid,
			},
			payloadUpdate: rcm.UpdateRevenueSharingRequest{
				OrgID:           oid,
				UserID:          uid,
				Partner:         spb.PartnerType_WU.String(),
				BoundType:       rcm.BoundType_INBOUND,
				RemitType:       rcm.RemitType_BILLSPAYMENT,
				TransactionType: rcm.TransactionType_OTC,
				TierType:        rcm.TierType_PERCENTAGE,
				Amount:          "",
				CreatedBy:       uid,
			},
			payloadUpdateRes: &rcm.UpdateRevenueSharingResponse{
				Partner:         spb.PartnerType_WU.String(),
				BoundType:       rcm.BoundType_INBOUND,
				RemitType:       rcm.RemitType_BILLSPAYMENT,
				TransactionType: rcm.TransactionType_OTC,
				TierType:        rcm.TierType_PERCENTAGE,
				Amount:          "",
				CreatedBy:       uid,
			},
			payloadGet: rcm.GetRevenueSharingListRequest{
				OrgID:           oid,
				UserID:          uid,
				Partner:         spb.PartnerType_WU.String(),
				BoundType:       rcm.BoundType_INBOUND,
				RemitType:       rcm.RemitType_BILLSPAYMENT,
				TransactionType: rcm.TransactionType_OTC,
				TierType:        rcm.TierType_PERCENTAGE,
				Amount:          "",
				CreatedBy:       uid,
			},
			payloadGetRes: []*rcm.RevenueSharing{
				{
					OrgID:           oid,
					UserID:          uid,
					Partner:         spb.PartnerType_WU.String(),
					BoundType:       rcm.BoundType_INBOUND,
					RemitType:       rcm.RemitType_BILLSPAYMENT,
					TransactionType: rcm.TransactionType_OTC,
					TierType:        rcm.TierType_PERCENTAGE,
					Amount:          "",
					CreatedBy:       uid,
				},
			},
			inputtier: []rcm.CreateRevenueSharingTierRequest{
				{
					MinValue: "1",
					MaxValue: "10000",
					Amount:   "60",
				},
				{
					MinValue: "10001",
					MaxValue: "30000",
					Amount:   "70",
				},
				{
					MinValue: "30001",
					MaxValue: "40000",
					Amount:   "80",
				},
			},
			inputtierRes: []*rcm.RevenueSharingTier{
				{
					MinValue: "1",
					MaxValue: "10000",
					Amount:   "60",
				},
				{
					MinValue: "10001",
					MaxValue: "30000",
					Amount:   "70",
				},
				{
					MinValue: "30001",
					MaxValue: "40000",
					Amount:   "80",
				},
			},
		},
	}
	ignore := cmp.Options{cmpopts.IgnoreFields(rcm.CreateRevenueSharingResponse{}, "ID", "Created", "Updated"), cmpopts.IgnoreUnexported(rcm.CreateRevenueSharingResponse{}, timestamppb.Timestamp{})}

	ignoreU := cmp.Options{cmpopts.IgnoreFields(rcm.UpdateRevenueSharingResponse{}, "ID", "Created", "Updated"), cmpopts.IgnoreUnexported(rcm.UpdateRevenueSharingResponse{}, timestamppb.Timestamp{})}

	ignoreTier := cmp.Options{cmpopts.IgnoreFields(rcm.RevenueSharingTier{}, "ID", "RevenueSharingID"), cmpopts.IgnoreUnexported(rcm.RevenueSharingTier{})}

	ignoreGet := cmp.Options{cmpopts.IgnoreFields(rcm.RevenueSharing{}, "ID", "Created", "Updated", "Count"), cmpopts.IgnoreUnexported(rcm.RevenueSharing{})}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			res, err := s.CreateRevenueSharing(ctx, &test.input)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.inputRes, res, ignore) {
				t.Error(cmp.Diff(test.inputRes, res))
			}
			test.payloadUpdate.ID = res.ID
			uRes, err := s.UpdateRevenueSharing(ctx, &test.payloadUpdate)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.payloadUpdateRes, uRes, ignoreU) {
				t.Error(cmp.Diff(test.payloadUpdateRes, uRes))
			}
			for _, it := range test.inputtier {
				it.RevenueSharingID = res.ID
				_, err := s.CreateRevenueSharingTier(ctx, &it)
				if err != nil {
					t.Fatal(err)
				}
			}
			comTierList, err := s.GetRevenueSharingTierList(ctx, &rcm.GetRevenueSharingTierListRequest{
				RevenueSharingID: res.ID,
			})
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(test.inputtierRes, comTierList.Results, ignoreTier) {
				t.Error(cmp.Diff(test.inputtierRes, comTierList.Results))
			}
			comList, err := s.GetRevenueSharingList(ctx, &test.payloadGet)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.payloadGetRes, comList.Results, ignoreGet) {
				t.Error(cmp.Diff(test.payloadGetRes, comList.Results))
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

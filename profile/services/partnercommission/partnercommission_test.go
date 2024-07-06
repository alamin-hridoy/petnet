package partnercommission

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	spb "brank.as/petnet/gunk/dsa/v2/partner"
	rcm "brank.as/petnet/gunk/dsa/v2/partnercommission"
	rcc "brank.as/petnet/profile/core/partnercommission"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/profile/storage/postgres"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus/hooks/test"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPartnerCommission(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	test.NewNullLogger()
	st := newTestStorage(t)
	s := New(rcc.New(st))

	uid := uuid.NewString()

	tests := []struct {
		desc             string
		input            rcm.CreatePartnerCommissionRequest
		inputRes         *rcm.CreatePartnerCommissionResponse
		inputtier        []rcm.CreatePartnerCommissionTierRequest
		inputtierRes     []*rcm.PartnerCommissionTier
		payload          storage.PartnerCommission
		payloadUpdate    rcm.UpdatePartnerCommissionRequest
		payloadGet       rcm.GetPartnerCommissionsListRequest
		payloadUpdateRes *rcm.UpdatePartnerCommissionResponse
		payloadGetRes    []*rcm.PartnerCommission
	}{
		{
			desc: "insert remittance partner commission",
			input: rcm.CreatePartnerCommissionRequest{
				Partner:         spb.PartnerType_AYA.String(),
				BoundType:       rcm.BoundType_INBOUND,
				RemitType:       rcm.RemitType_REMITTANCE,
				TransactionType: rcm.TransactionType_DIGITAL,
				TierType:        rcm.TierType_FIXED,
				Amount:          "10",
				StartDate:       timestamppb.Now(),
				EndDate:         timestamppb.Now(),
				CreatedBy:       uid,
			},
			inputRes: &rcm.CreatePartnerCommissionResponse{
				Partner:         spb.PartnerType_AYA.String(),
				BoundType:       rcm.BoundType_INBOUND,
				RemitType:       rcm.RemitType_REMITTANCE,
				TransactionType: rcm.TransactionType_DIGITAL,
				TierType:        rcm.TierType_FIXED,
				Amount:          "10",
				StartDate:       timestamppb.Now(),
				EndDate:         timestamppb.Now(),
				CreatedBy:       uid,
			},
			payloadUpdate: rcm.UpdatePartnerCommissionRequest{
				Partner:         spb.PartnerType_AYA.String(),
				BoundType:       rcm.BoundType_INBOUND,
				RemitType:       rcm.RemitType_REMITTANCE,
				TransactionType: rcm.TransactionType_DIGITAL,
				TierType:        rcm.TierType_TIERPERCENTAGE,
				Amount:          "",
				CreatedBy:       uid,
			},
			payloadUpdateRes: &rcm.UpdatePartnerCommissionResponse{
				Partner:         spb.PartnerType_AYA.String(),
				BoundType:       rcm.BoundType_INBOUND,
				RemitType:       rcm.RemitType_REMITTANCE,
				TransactionType: rcm.TransactionType_DIGITAL,
				TierType:        rcm.TierType_TIERPERCENTAGE,
				Amount:          "",
				CreatedBy:       uid,
			},
			inputtier: []rcm.CreatePartnerCommissionTierRequest{
				{
					CommissionTier: []*rcm.PartnerCommissionTier{
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
				},
			},
			inputtierRes: []*rcm.PartnerCommissionTier{
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
			payloadGet: rcm.GetPartnerCommissionsListRequest{
				Partner:         spb.PartnerType_AYA.String(),
				BoundType:       rcm.BoundType_INBOUND,
				RemitType:       rcm.RemitType_REMITTANCE,
				TransactionType: rcm.TransactionType_DIGITAL,
				TierType:        rcm.TierType_TIERPERCENTAGE,
				Amount:          "",
				CreatedBy:       uid,
			},
			payloadGetRes: []*rcm.PartnerCommission{
				{
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
			input: rcm.CreatePartnerCommissionRequest{
				Partner:         spb.PartnerType_WU.String(),
				BoundType:       rcm.BoundType_INBOUND,
				RemitType:       rcm.RemitType_BILLSPAYMENT,
				TransactionType: rcm.TransactionType_OTC,
				TierType:        rcm.TierType_PERCENTAGE,
				Amount:          "20",
				StartDate:       timestamppb.Now(),
				EndDate:         timestamppb.Now(),
				CreatedBy:       uid,
			},
			inputRes: &rcm.CreatePartnerCommissionResponse{
				Partner:         spb.PartnerType_WU.String(),
				BoundType:       rcm.BoundType_INBOUND,
				RemitType:       rcm.RemitType_BILLSPAYMENT,
				TransactionType: rcm.TransactionType_OTC,
				TierType:        rcm.TierType_PERCENTAGE,
				Amount:          "20",
				StartDate:       timestamppb.Now(),
				EndDate:         timestamppb.Now(),
				CreatedBy:       uid,
			},
			payloadUpdate: rcm.UpdatePartnerCommissionRequest{
				Partner:         spb.PartnerType_WU.String(),
				BoundType:       rcm.BoundType_INBOUND,
				RemitType:       rcm.RemitType_BILLSPAYMENT,
				TransactionType: rcm.TransactionType_OTC,
				TierType:        rcm.TierType_FIXED,
				Amount:          "",
				CreatedBy:       uid,
			},
			payloadUpdateRes: &rcm.UpdatePartnerCommissionResponse{
				Partner:         spb.PartnerType_WU.String(),
				BoundType:       rcm.BoundType_INBOUND,
				RemitType:       rcm.RemitType_BILLSPAYMENT,
				TransactionType: rcm.TransactionType_OTC,
				TierType:        rcm.TierType_FIXED,
				Amount:          "",
				CreatedBy:       uid,
			},
			payloadGet: rcm.GetPartnerCommissionsListRequest{
				Partner:         spb.PartnerType_WU.String(),
				BoundType:       rcm.BoundType_INBOUND,
				RemitType:       rcm.RemitType_BILLSPAYMENT,
				TransactionType: rcm.TransactionType_OTC,
				TierType:        rcm.TierType_FIXED,
				Amount:          "",
				CreatedBy:       uid,
			},
			payloadGetRes: []*rcm.PartnerCommission{
				{
					Partner:         spb.PartnerType_WU.String(),
					BoundType:       rcm.BoundType_INBOUND,
					RemitType:       rcm.RemitType_BILLSPAYMENT,
					TransactionType: rcm.TransactionType_OTC,
					TierType:        rcm.TierType_FIXED,
					Amount:          "",
					CreatedBy:       uid,
				},
			},
			inputtier: []rcm.CreatePartnerCommissionTierRequest{
				{
					CommissionTier: []*rcm.PartnerCommissionTier{
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
			},
			inputtierRes: []*rcm.PartnerCommissionTier{
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
	ignore := cmp.Options{cmpopts.IgnoreFields(rcm.CreatePartnerCommissionResponse{}, "ID", "Created", "Updated", "StartDate", "EndDate"), cmpopts.IgnoreUnexported(rcm.CreatePartnerCommissionResponse{}, timestamppb.Timestamp{})}

	ignoreU := cmp.Options{cmpopts.IgnoreFields(rcm.UpdatePartnerCommissionResponse{}, "ID", "Created", "Updated", "StartDate", "EndDate"), cmpopts.IgnoreUnexported(rcm.UpdatePartnerCommissionResponse{}, timestamppb.Timestamp{})}

	ignoreTier := cmp.Options{cmpopts.IgnoreFields(rcm.PartnerCommissionTier{}, "ID", "PartnerCommissionID"), cmpopts.IgnoreUnexported(rcm.PartnerCommissionTier{})}

	ignoreGet := cmp.Options{cmpopts.IgnoreFields(rcm.PartnerCommission{}, "ID", "Created", "Updated", "StartDate", "EndDate", "Count"), cmpopts.IgnoreUnexported(rcm.PartnerCommission{})}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			res, err := s.CreatePartnerCommission(ctx, &test.input)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.inputRes, res, ignore) {
				t.Error(cmp.Diff(test.inputRes, res))
			}
			test.payloadUpdate.ID = res.ID
			uRes, err := s.UpdatePartnerCommission(ctx, &test.payloadUpdate)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.payloadUpdateRes, uRes, ignoreU) {
				t.Error(cmp.Diff(test.payloadUpdateRes, uRes))
			}
			for _, it := range test.inputtier {
				for i := range it.GetCommissionTier() {
					it.CommissionTier[i].PartnerCommissionID = res.ID
				}
				_, err := s.CreatePartnerCommissionTier(ctx, &it)
				if err != nil {
					t.Fatal(err)
				}
			}
			comTierList, err := s.GetPartnerCommissionsTierList(ctx, &rcm.GetPartnerCommissionsTierListRequest{
				PartnerCommissionID: res.ID,
			})
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(test.inputtierRes, comTierList.Results, ignoreTier) {
				t.Error(cmp.Diff(test.inputtierRes, comTierList.Results))
			}
			comList, err := s.GetPartnerCommissionsList(ctx, &test.payloadGet)
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

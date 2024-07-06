package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	spb "brank.as/petnet/gunk/dsa/v2/partner"
	svc "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/profile/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestPartnerCommission(t *testing.T) {
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
		input         storage.PartnerCommission
		inputtier     []storage.PartnerCommissionTier
		payload       storage.PartnerCommission
		payloadUpdate storage.PartnerCommission
	}{
		{
			desc: "insert remittance partner commission",
			input: storage.PartnerCommission{
				Partner:         spb.PartnerType_AYA.String(),
				RemitType:       svc.ServiceType_REMITTANCE.String(),
				TransactionType: string(storage.TransactionType_Digital),
				TierType:        string(storage.TierType_Fixed_Amount),
				BoundType:       string(storage.BoundType_In),
				StartDate:       time.Now(),
				EndDate:         time.Now(),
				Amount:          "10",
				CreatedBy:       uid,
				UpdatedBy:       uid,
			},
			inputtier: []storage.PartnerCommissionTier{
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
			payloadUpdate: storage.PartnerCommission{
				Partner:         spb.PartnerType_AYA.String(),
				RemitType:       svc.ServiceType_REMITTANCE.String(),
				TransactionType: string(storage.TransactionType_Digital),
				TierType:        string(storage.TierType_Fixed_Tier_Percentage),
				BoundType:       string(storage.BoundType_In),
				Amount:          "",
				CreatedBy:       uid,
				UpdatedBy:       uid,
			},
		},
		{
			desc: "insert bills payment partner commission",
			input: storage.PartnerCommission{
				Partner:         spb.PartnerType_WU.String(),
				RemitType:       svc.ServiceType_BILLSPAYMENT.String(),
				TransactionType: string(storage.TransactionType_Otc),
				TierType:        string(storage.TierType_Fixed_Percentage),
				BoundType:       string(storage.BoundType_Others),
				StartDate:       time.Now(),
				EndDate:         time.Now(),
				Amount:          "20",
				CreatedBy:       uid,
				UpdatedBy:       uid,
			},
			inputtier: []storage.PartnerCommissionTier{
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
			payloadUpdate: storage.PartnerCommission{
				Partner:         spb.PartnerType_WU.String(),
				RemitType:       svc.ServiceType_BILLSPAYMENT.String(),
				TransactionType: string(storage.TransactionType_Otc),
				TierType:        string(storage.TierType_Fixed_Tier_Amount),
				BoundType:       string(storage.BoundType_Others),
				Amount:          "",
				CreatedBy:       uid,
				UpdatedBy:       uid,
			},
		},
	}
	ignore := cmp.Options{
		cmpopts.IgnoreFields(
			storage.PartnerCommission{}, "ID", "Created", "Updated", "Count",
		),
	}
	ignoreTier := cmp.Options{
		cmpopts.IgnoreFields(
			storage.PartnerCommissionTier{}, "ID", "PartnerCommissionID",
		),
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			res, err := ts.CreatePartnerCommission(ctx, test.input)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(&test.input, res, ignore) {
				t.Error(cmp.Diff(&test.input, res))
			}
			comId := res.ID
			test.payloadUpdate.ID = comId
			test.payloadUpdate.Updated = time.Now()
			uRes, err := ts.UpdatePartnerCommission(ctx, test.payloadUpdate)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(&test.payloadUpdate, uRes, ignore) {
				t.Error(cmp.Diff(&test.payloadUpdate, uRes))
			}
			for _, it := range test.inputtier {
				it.PartnerCommissionID = comId
				_, err := ts.CreatePartnerCommissionTier(ctx, it)
				if err != nil {
					t.Fatal(err)
				}
			}
			comTierList, err := ts.GetPartnerCommissionsTierList(ctx, storage.PartnerCommissionTier{
				PartnerCommissionID: comId,
			})
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.inputtier, comTierList, ignoreTier) {
				t.Error(cmp.Diff(test.inputtier, comTierList))
			}
			comList, err := ts.GetPartnerCommissionsList(ctx, test.payloadUpdate)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.payloadUpdate, comList[0], ignore) {
				t.Error(cmp.Diff(test.payloadUpdate, comList[0]))
			}
			for _, it := range comTierList {
				it.Amount = "50"
				_, err := ts.UpdatePartnerCommissionTier(ctx, it)
				if err != nil {
					t.Fatal(err)
				}
			}
			ucomTierList, err := ts.GetPartnerCommissionsTierList(ctx, storage.PartnerCommissionTier{
				PartnerCommissionID: comId,
			})
			if err != nil {
				t.Fatal(err)
			}
			comTierList[0].Amount = "50"
			if !cmp.Equal(ucomTierList[0], comTierList[0], ignoreTier) {
				t.Error(cmp.Diff(ucomTierList[0], comTierList[0]))
			}
			if err := ts.DeletePartnerCommissionTier(ctx, storage.PartnerCommissionTier{
				PartnerCommissionID: comId,
			}); err != nil {
				t.Fatal(err)
			}
			adcomTierList, err := ts.GetPartnerCommissionsTierList(ctx, storage.PartnerCommissionTier{
				PartnerCommissionID: comId,
			})
			if err != nil {
				t.Fatal(err)
			}
			if len(adcomTierList) != 0 {
				t.Fatal(errors.New("commission tier should be 0"))
			}
			if err := ts.DeletePartnerCommission(ctx, storage.PartnerCommission{
				Partner:         test.payloadUpdate.Partner,
				RemitType:       test.payloadUpdate.RemitType,
				TransactionType: test.payloadUpdate.TransactionType,
				BoundType:       test.payloadUpdate.BoundType,
			}); err != nil {
				t.Fatal(err)
			}
			aDcomList, err := ts.GetPartnerCommissionsList(ctx, test.payloadUpdate)
			if err != nil {
				t.Fatal(err)
			}
			if len(aDcomList) != 0 {
				t.Fatal(errors.New("commission should be 0"))
			}
		})
	}
}

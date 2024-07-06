package util

import (
	"context"
	"testing"
	"time"

	"brank.as/petnet/api/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestRecordCiCo(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	o := cmp.Options{
		cmpopts.IgnoreFields(
			storage.CashInCashOutHistory{}, "ID", "OrgID", "TrxDate", "Details", "CreatedBy", "UpdatedBy", "Count", "SortByColumn", "SortOrder", "Limit", "Offset", "Total", "Created", "Updated",
		),
	}
	tests := []struct {
		name string
		in   storage.CashInCashOutHistory
		res  storage.CashInCashOutHistoryRes
		want storage.CashInCashOutHistory
	}{
		{
			name: "Success Stage",
			in: storage.CashInCashOutHistory{
				OrgID:              uuid.NewString(),
				PartnerCode:        "DSA",
				Provider:           "PRV",
				PetnetTrackingNo:   "KLMNOPQRST",
				TrxType:            "TRX",
				ProviderTrackingNo: "UVWXYZ",
				ReferenceNumber:    "ABCDEFGHIJ",
				PrincipalAmount:    1000,
				Charges:            10,
				TotalAmount:        110,
				SvcProvider:        "GCASH",
				TrxDate:            time.Now(),
				Details:            []byte{},
				TxnStatus:          "SUCCESS",
				ErrorCode:          "",
				ErrorMessage:       "",
				ErrorTime:          "",
				ErrorType:          "",
				CreatedBy:          uuid.NewString(),
				UpdatedBy:          uuid.NewString(),
				Created:            time.Now(),
				Updated:            time.Now(),
			},
			res: storage.CashInCashOutHistoryRes{
				Code:    200,
				Message: "Success",
				Result: storage.CashInCashOutHistoryDetails{
					ID:                 "id",
					PartnerCode:        "DSA",
					Provider:           "PRV",
					PetnetTrackingno:   "KLMNOPQRST",
					TrxType:            "TRX",
					ProviderTrackingNo: "UVWXYZ",
					ReferenceNumber:    "ABCDEFGHIJ",
					PrincipalAmount:    1000,
					Charges:            10,
					TotalAmount:        110,
				},
			},
			want: storage.CashInCashOutHistory{
				OrgID:              uuid.NewString(),
				PartnerCode:        "DSA",
				Provider:           "PRV",
				PetnetTrackingNo:   "KLMNOPQRST",
				TrxType:            "TRX",
				ProviderTrackingNo: "UVWXYZ",
				ReferenceNumber:    "ABCDEFGHIJ",
				PrincipalAmount:    1000,
				Charges:            10,
				TotalAmount:        110,
				SvcProvider:        "GCASH",
				TrxDate:            time.Now(),
				Details:            []byte{},
				TxnStatus:          "SUCCESS",
				ErrorCode:          "",
				ErrorMessage:       "",
				ErrorTime:          "",
				ErrorType:          "",
				CreatedBy:          uuid.NewString(),
				UpdatedBy:          uuid.NewString(),
				Created:            time.Now(),
				Updated:            time.Now(),
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			var err error
			aRes, aErr := RecordCiCo(ctx, st, test.in, test.res, err)
			if aErr != nil {
				t.Fatal(aErr)
			}
			if !cmp.Equal(&test.want, aRes, o) {
				t.Error(cmp.Diff(&test.want, aRes, o))
			}
		})
	}
}

package bills_payment

import (
	"context"
	"testing"

	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBCTransact(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      BCTransactRequest
		want    *BCTransactResponse
		wantErr bool
	}{
		{
			name: "Success",
			in: BCTransactRequest{
				UserID:                  "5500",
				CustomerID:              6925598,
				LocationID:              "371",
				LocationName:            "Information Technology Department",
				Coy:                     "usp",
				CallbackURL:             "",
				BillID:                  "2",
				BillerTag:               "ADMSN",
				BillerName:              "ADAMSON UNIVERSITY",
				TrxDate:                 "2022-10-17",
				Amount:                  "1000",
				ServiceCharge:           "0.00",
				PartnerCharge:           "0",
				TotalAmount:             1000,
				Identifier:              "Jenson Masangcay",
				AccountNumber:           "200820248",
				PaymentMethod:           "CASH",
				ClientReferenceNumber:   "PP012134103321",
				ReferenceNumber:         "0156600054001128",
				ValidationNumber:        "d3a5f72b-27e4-4fe1-aa57-2c417d7a8a2d",
				ReceiptValidationNumber: "PP01211870000010",
				TpaID:                   "PP01",
				CurrencyID:              "1",
				FormType:                "OAR",
				FormNumber:              "HOA20217610",
				OtherInfo: BCOtherInfo{
					LastName:        "Serato",
					FirstName:       "Mike Edward",
					MiddleName:      "Secret",
					PaymentType:     "B2",
					Course:          "BSCpE",
					TotalAssessment: "0.00",
					SchoolYear:      "2021-2022",
					Term:            "1",
				},
				Type: "Batch",
			},
			want: &BCTransactResponse{
				Code:    200,
				Message: "Success",
				Result: BCTransactResult{
					TransactionID:   "21187PP01024FF38I",
					ReferenceNumber: "511111",
					ClientReference: "6217d18d-5cff-4f1e-affd-6503883dflk0",
					BillerReference: "PP012118724FF38I",
					PaymentMethod:   "CASH",
					Amount:          "100.00",
					OtherCharges:    "0.00",
					Status:          "PENDING",
					Message:         "The payment was successfully created.",
					Details:         []*bp.Details{},
					CreatedAt:       "2021-07-06 16:57:29",
					Timestamp:       "2021-07-06 08:57:25",
				},
				RemcoID: 2,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.BCTransact(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("BCTransact() error = %v, wantErr %v", err, test.wantErr)
			}
			tOps := []cmp.Option{
				cmpopts.IgnoreFields(BCTransactInquireResult{}, "Details"),
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBCRetry(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      BCRetryRequest
		want    *BCRetryResponse
		wantErr bool
	}{
		{
			name: "Success",
			in: BCRetryRequest{
				Coy:      "usp",
				Type:     "Batch",
				Amount:   "100",
				TpaID:    "PP01",
				BillID:   "2",
				UserID:   "5500",
				TrxDate:  "2021-10-13",
				FormType: "OAR",
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
				BillerTag:               "ADMSN",
				Identifier:              "Jenson Masangcay",
				BillerName:              "ADAMSON UNIVERSITY",
				CurrencyID:              "1",
				CustomerID:              6925598,
				FormNumber:              "",
				LocationID:              "371",
				TotalAmount:             100,
				LocationName:            "Information Technology Department",
				AccountNumber:           "200820248",
				PartnerCharge:           "0",
				PaymentMethod:           "CASH",
				ServiceCharge:           "0.00",
				ReferenceNumber:         "0156600054001128",
				ValidationNumber:        "d3a5f72b-27e4-4fe1-aa57-2c417d7a8a2d",
				ClientReferenceNumber:   "PP012134103320",
				ReceiptValidationNumber: "PP012134103320",
				ID:                      2902,
			},
			want: &BCRetryResponse{
				Code:    200,
				Message: "Success",
				Result: BCRetryResult{
					TransactionID:   "21187PP01024FF38I",
					ReferenceNumber: "511111",
					ClientReference: "6217d18d-5cff-4f1e-affd-6503883dflk0",
					BillerReference: "PP012118724FF38I",
					PaymentMethod:   "CASH",
					Amount:          "100.00",
					OtherCharges:    "0.00",
					Status:          "PENDING",
					Message:         "The payment was successfully created.",
					Details:         []interface{}{},
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
			got, err := s.BCRetry(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("BCRetry() error = %v, wantErr %v", err, test.wantErr)
			}
			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

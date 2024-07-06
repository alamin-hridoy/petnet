package bills_payment

import (
	"context"
	"testing"

	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBCTransactInquire(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      BCTransactInquireRequest
		want    *BCTransactInquireResponse
		wantErr bool
	}{
		{
			name: "Success",
			in: BCTransactInquireRequest{
				Code:            "ASVCA",
				ClientReference: "PP01212860000002",
			},
			want: &BCTransactInquireResponse{
				Code:    200,
				Message: "Success",
				Result: BCTransactInquireResult{
					TransactionID:   "21152PP0120000006",
					ReferenceNumber: "0136402637",
					ClientReference: "a1mef46c-d994-487e-a086-oaesdc78a42f",
					BillerReference: "21152PP0120000006",
					PaymentMethod:   "CASH",
					Amount:          "3974.83",
					OtherCharges:    "0.00",
					Status:          "PENDING",
					Message:         BCTransactInquireResultMessage{Header: "Payment Receipt", Message: "Sweet! We have received your MERALCO bill payment and are currently processing it. Thank you. Have a great day ahead!", Footer: "Please note that payments made after 7PM will be posted 7AM the next day."},
					Details:         []*bp.Details{},
					CreatedAt:       "2021-06-01 16:41:41",
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
			got, err := s.BCTransactInquire(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("BCTransactInquire() error = %v, wantErr %v", err, test.wantErr)
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

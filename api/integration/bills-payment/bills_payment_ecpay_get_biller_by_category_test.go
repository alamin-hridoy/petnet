package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBillsPaymentEcpayBillerByCategory(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *BillsPaymentEcpayBillerByCategoryResponseBody
		wantErr bool
	}{
		{
			name: "Success",
			want: &BillsPaymentEcpayBillerByCategoryResponseBody{
				Code:    200,
				Message: "Success",
				Result: []BillsPaymentEcpayBillerByCategoryResult{
					{
						BillerTag:         "MANILAWATER",
						Description:       "MANILA WATER COMPANY",
						FirstField:        "8 Digit Contract Account Number",
						FirstFieldFormat:  "Numeric",
						FirstFieldWidth:   "8",
						SecondField:       "Account Name",
						SecondFieldFormat: "Alphanumeric",
						SecondFieldWidth:  "30",
						ServiceCharge:     2,
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.BillsPaymentEcpayBillerByCategory(context.Background(), 1)
			if (err != nil) != test.wantErr {
				t.Fatalf("BillsPaymentEcpayBillerByCategory() error = %v, wantErr %v", err, test.wantErr)
			}

			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

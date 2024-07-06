package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBillsPaymentEcpayBillerlist(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *BillsPaymentEcpayBillerlistResponseBody
		wantErr bool
	}{
		{
			name: "Success",
			want: &BillsPaymentEcpayBillerlistResponseBody{
				Code:    200,
				Message: "Success",
				Result: []BillsPaymentEcpayBillerlistResult{
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
				RemcoID: 1,
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.BillsPaymentEcpayBillerlist(context.Background())
			if (err != nil) != test.wantErr {
				t.Fatalf("BillsPaymentEcpayBillerlist() error = %v, wantErr %v", err, test.wantErr)
			}

			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

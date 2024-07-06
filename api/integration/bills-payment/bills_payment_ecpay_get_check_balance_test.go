package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBillsPaymentEcpayCheckBalance(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *BillsPaymentEcpayCheckBalanceResponseBody
		wantErr bool
	}{
		{
			name: "Success",
			want: &BillsPaymentEcpayCheckBalanceResponseBody{
				Code:    200,
				Message: "Success",
				Result: BillsPaymentEcpayCheckBalanceResult{
					RemBal: "10000.00",
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
			got, err := s.BillsPaymentEcpayCheckBalance(context.Background())
			if (err != nil) != test.wantErr {
				t.Fatalf("BillsPaymentEcpayCheckBalance() error = %v, wantErr %v", err, test.wantErr)
			}

			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

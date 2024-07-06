package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBillsPaymentEcpayRetry(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      BillsPaymentEcpayRetryRequest
		want    *BillsPaymentEcpayRetryResponseBody
		wantErr bool
	}{
		{
			name: "Success",
			in: BillsPaymentEcpayRetryRequest{
				ID: 3115,
			},
			want: &BillsPaymentEcpayRetryResponseBody{
				Code:    "0",
				Message: "Success",
				Result: BillsPaymentEcpayRetryResult{
					Status:          "0",
					Message:         "SUCCESS! REF #72482FD0A467",
					ServiceCharge:   10,
					Timestamp:       "2021-03-28 08:58:28",
					ReferenceNumber: "72482FD0A467",
				},
				RemcoID: 1,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.BillsPaymentEcpayRetry(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("BillsPaymentEcpayRetry() error = %v, wantErr %v", err, test.wantErr)
			}
			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

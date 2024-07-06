package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBillsPaymentMultiPayBillerInquire(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      BillsPaymentMultiPayBillerInquireRequest
		want    *BillsPaymentMultiPayBillerInquireResponseBody
		wantErr bool
	}{
		{
			name: "Success",
			in: BillsPaymentMultiPayBillerInquireRequest{
				AccountNumber: "MP1MN7JMJV",
				Amount:        51,
				ContactNumber: "09123456789",
			},
			want: &BillsPaymentMultiPayBillerInquireResponseBody{
				Status: 200,
				Reason: "OK",
				Data: BillsPaymentMultiPayBillerInquireData{
					AccountNumber: "MP1MN7JMJV",
					Amount:        51,
					Biller:        "MSYS_TEST_BILLER",
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.BillsPaymentMultiPayBillerInquire(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("BillsPaymentMultiPayBillerInquire() error = %v, wantErr %v", err, test.wantErr)
			}
			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

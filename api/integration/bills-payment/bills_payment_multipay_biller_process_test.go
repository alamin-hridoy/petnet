package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBillsPaymentMultiPayBillerProcess(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      BillsPaymentMultiPayBillerProcessRequest
		want    *BillsPaymentMultiPayBillerProcessResponseBody
		wantErr bool
	}{
		{
			name: "Success",
			in: BillsPaymentMultiPayBillerProcessRequest{
				AccountNumber: "MP1MN7JMJV",
				Amount:        51,
				ContactNumber: "09123456789",
			},
			want: &BillsPaymentMultiPayBillerProcessResponseBody{
				Status: 200,
				Reason: "You have successfully processed this transaction.",
				Data: BillsPaymentMultiPayBillerProcessData{
					Refno:  "MP1MN7JMJV",
					Txnid:  "TEST-62F0802B98B6E",
					Biller: "MSYS_TEST_BILLER",
					Meta:   []interface{}{},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.BillsPaymentMultiPayBillerProcess(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("BillsPaymentMultiPayBillerProcess() error = %v, wantErr %v", err, test.wantErr)
			}
			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBillsPaymentMultiPayTransact(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      BillsPaymentMultiPayTransactRequest
		want    *BillsPaymentMultiPayTransactResponseBody
		wantErr bool
	}{
		{
			name: "Success",
			in: BillsPaymentMultiPayTransactRequest{
				Amount: "51",
				Txnid:  "TEST-62F0802B98B6E",
			},
			want: &BillsPaymentMultiPayTransactResponseBody{
				Data: BillsPaymentMultiPayTransactData{
					URL: "https://pgi-staging.multipay.ph/MP0XQPUWDTVSRUONACH",
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.BillsPaymentMultiPayTransact(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("BillsPaymentMultiPayTransact() error = %v, wantErr %v", err, test.wantErr)
			}
			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

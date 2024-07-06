package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBillsPaymentMultiPaySearchTrx(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      BillsPaymentMultiPaySearchTrxRequest
		want    *BillsPaymentMultiPaySearchTrxResponseBody
		wantErr bool
	}{
		{
			name: "Success",
			in: BillsPaymentMultiPaySearchTrxRequest{
				ReferenceNo: "MP9UYHSEHSIKRTY4RHC",
			},
			want: &BillsPaymentMultiPaySearchTrxResponseBody{
				Code:    200,
				Message: "Good",
				Result: BillsPaymentMultiPaySearchTrxResult{
					Txnid:                "TEST-62F0802B98B6E",
					Refno:                "MP9UYHSEHSIKRTY4RHC",
					Amount:               "51.00",
					Fee:                  "0.00",
					Status:               "V",
					PaymentChannel:       "PGI",
					IsTransactionExpired: true,
					CreatedAt:            "2022-08-09 16:48:58",
					ExpiresAt:            "2022-08-12 16:48:58",
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.BillsPaymentMultiPaySearchTrx(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("BillsPaymentMultiPaySearchTrx() error = %v, wantErr %v", err, test.wantErr)
			}
			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBillsPaymentMultiPayVoidTrx(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      BillsPaymentMultiPayVoidTrxRequest
		want    *BillsPaymentMultiPayVoidTrxResponseBody
		wantErr bool
	}{
		{
			name: "Success",
			in: BillsPaymentMultiPayVoidTrxRequest{
				ReferenceNo: "MP9UYHSEHSIKRTY4RHC",
			},
			want: &BillsPaymentMultiPayVoidTrxResponseBody{
				Data: BillsPaymentMultiPayVoidTrxData{
					Txnid:                "TEST-62F0802B98B6E",
					Refno:                "MP0XQPUWDTVSRUONACH",
					Amount:               "51.00",
					Fee:                  "0.00",
					Status:               "V",
					PaymentChannel:       "PGI",
					IsTransactionExpired: true,
					CreatedAt:            "2022-09-05 07:41:36",
					ExpiresAt:            "2022-09-08 07:41:36",
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.BillsPaymentMultiPayVoidTrx(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("BillsPaymentMultiPayVoidTrx() error = %v, wantErr %v", err, test.wantErr)
			}
			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

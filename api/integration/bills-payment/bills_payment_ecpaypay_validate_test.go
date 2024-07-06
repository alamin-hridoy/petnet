package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBillsPaymentEcpayValidate(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      BillsPaymentEcpayValidateRequest
		want    *BillsPaymentEcpayValidateResponseBody
		wantErr bool
	}{
		{
			name: "Success",
			in: BillsPaymentEcpayValidateRequest{
				AccountNo:  "99999222229",
				Identifier: "Masangcay,Jenson",
				BillerTag:  "DAVAOLIGHT",
			},
			want: &BillsPaymentEcpayValidateResponseBody{
				Code:    200,
				Message: "Success",
				Result:  "DAVAOLIGHT",
				RemcoID: 1,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.BillsPaymentEcpayValidate(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("BillsPaymentEcpayValidate() error = %v, wantErr %v", err, test.wantErr)
			}
			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

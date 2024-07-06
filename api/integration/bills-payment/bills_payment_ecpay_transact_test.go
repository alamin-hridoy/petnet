package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBillsPaymentEcpayTransact(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      BillsPaymentEcpayTransactRequest
		want    *BillsPaymentEcpayTransactResponseBody
		wantErr bool
	}{
		{
			name: "Success",
			in: BillsPaymentEcpayTransactRequest{
				BillID:                1,
				BillerTag:             "DAVAOLIGHT",
				TrxDate:               "2022-04-21",
				UserID:                "5500",
				RemoteUserID:          "5500",
				CustomerID:            "6925598",
				LocationID:            "371",
				RemoteLocationID:      "371",
				LocationName:          "Information Technology",
				Coy:                   "usp",
				CurrencyID:            "1",
				FormType:              "OAR",
				FormNumber:            "IT03",
				AccountNumber:         "99999222229",
				Identifier:            "Masangcay,Jenson",
				Amount:                203,
				ServiceCharge:         10,
				TotalAmount:           213,
				ClientReferenceNumber: "PRHBfdwdASDFSADFSF",
			},
			want: &BillsPaymentEcpayTransactResponseBody{
				Code:    "0",
				Message: "Success",
				Result: BillsPaymentEcpayTransactResult{
					Status:          "0",
					Message:         "SUCCESS! REF #F6L2S84JN00M",
					ServiceCharge:   10,
					Timestamp:       "2022-05-19 06:46:05",
					ReferenceNumber: "F6L2S84JN00M",
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
			got, err := s.BillsPaymentEcpayTransact(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("BillsPaymentEcpayTransact() error = %v, wantErr %v", err, test.wantErr)
			}
			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

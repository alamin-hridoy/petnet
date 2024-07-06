package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBillsPaymentMultiPayBillerByCategory(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *BillsPaymentMultiPayBillerByCategoryResponseBody
		wantErr bool
	}{
		{
			name: "Success",
			want: &BillsPaymentMultiPayBillerByCategoryResponseBody{
				Code:    200,
				Message: "Good",
				Result: []BillsPaymentMultiPayBillerByCategoryResult{
					{
						PartnerID:   3,
						BillerTag:   "MULTIPAY-AppendPay",
						Description: "AppendPay",
						Category:    1,
						FieldList: []FieldList{
							{
								ID:    "amount",
								Type:  "numeric",
								Label: "Amount",
								Order: 1,
								Rules: []Rules{
									{
										Code:    1,
										Type:    "required",
										Value:   "",
										Format:  "",
										Message: "Please provide the amount.",
										Options: "",
									},
								},
								Description: "Amount to be paid",
								Placeholder: "Insert Amount",
							},
						},
						ServiceCharge: 0,
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.BillsPaymentMultiPayBillerByCategory(context.Background(), 1)
			if (err != nil) != test.wantErr {
				t.Fatalf("BillsPaymentMultiPayBillerByCategory() error = %v, wantErr %v", err, test.wantErr)
			}

			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

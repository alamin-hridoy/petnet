package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBillsPaymentMultiPayBillerlist(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *BillsPaymentMultiPayBillerlistResponseBody
		wantErr bool
	}{
		{
			name: "Success",
			want: &BillsPaymentMultiPayBillerlistResponseBody{
				Code:    200,
				Message: "Good",
				Result: []BillsPaymentMultiPayBillerlistResult{
					{
						PartnerID:   3,
						BillerTag:   "MULTIPAY-NBI",
						Description: "NBI",
						Category:    4,
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
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.BillsPaymentMultiPayBillerlist(context.Background())
			if (err != nil) != test.wantErr {
				t.Fatalf("BillsPaymentMultiPayBillerlist() error = %v, wantErr %v", err, test.wantErr)
			}

			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

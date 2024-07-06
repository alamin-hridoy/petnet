package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBillsPaymentMultiPayBillerCategory(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *BillsPaymentMultiPayBillerCategoryResponseBody
		wantErr bool
	}{
		{
			name: "Success",
			want: &BillsPaymentMultiPayBillerCategoryResponseBody{
				Code:    200,
				Message: "Good",
				Result: []BillsPaymentMultiPayBillerCategoryResult{
					{
						ID:           6,
						BillID:       1,
						CategoryName: "Airlines",
						CreatedAt:    "2022-04-08T04:26:01.000000Z",
						UpdatedAt:    "2022-04-08T04:26:01.000000Z",
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
			got, err := s.BillsPaymentMultiPayBillerCategory(context.Background())
			if (err != nil) != test.wantErr {
				t.Fatalf("BillsPaymentMultiPayBillerCategory() error = %v, wantErr %v", err, test.wantErr)
			}

			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

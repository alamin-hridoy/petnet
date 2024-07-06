package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBCGettingOtherCharges(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      BCGettingOtherChargesRequest
		want    *BCGettingOtherChargesResponse
		wantErr bool
	}{
		{
			name: "Success",
			in: BCGettingOtherChargesRequest{
				Code:       "PWCOR",
				Amount:     "1000",
				UserID:     "1",
				LocationID: "371",
			},
			want: &BCGettingOtherChargesResponse{
				Code:    200,
				Message: "Success",
				Result: BCGettingOtherChargesResult{
					OtherCharges: "5.00",
				},
				RemcoID: 2,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.BCGettingOtherCharges(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("BCGettingOtherCharges() error = %v, wantErr %v", err, test.wantErr)
			}
			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

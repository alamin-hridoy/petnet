package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBCGetWalletBalance(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *BCGetWalletBalanceResponse
		wantErr bool
	}{
		{
			name: "Success",
			want: &BCGetWalletBalanceResponse{
				Code:    200,
				Message: "Success",
				Result: BCGetWalletBalanceResult{
					Balance: "0.00",
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
			got, err := s.BCGetWalletBalance(context.Background())
			if (err != nil) != test.wantErr {
				t.Fatalf("BCGetWalletBalance() error = %v, wantErr %v", err, test.wantErr)
			}

			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

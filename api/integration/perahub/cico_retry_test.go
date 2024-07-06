package perahub

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCicoRetry(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      CicoRetryRequest
		want    *CicoRetryResponse
		wantErr bool
	}{
		{
			name: "Success",
			in: CicoRetryRequest{
				PartnerCode:      "DSA",
				PetnetTrackingno: "5a269417e107691f3d7c",
				TrxDate:          "2022-05-17",
			},
			want: &CicoRetryResponse{
				Code:    200,
				Message: "SUCCESS TRANSACTION.",
				Result: &CicoRetryResult{
					PartnerCode:        "DSA",
					Provider:           "GCASH",
					PetnetTrackingno:   "5a269417e107691f3d7c",
					TrxDate:            "2022-05-17",
					TrxType:            "Cash In",
					ProviderTrackingno: "09654767706",
					ReferenceNumber:    "09654767706",
					PrincipalAmount:    10,
					Charges:            0,
					TotalAmount:        10,
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.CicoRetry(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("CicoRetry() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

package perahub

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCicoExecute(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      CicoExecuteRequest
		want    *CicoExecuteResponse
		wantErr bool
	}{
		{
			name: "Success",
			in: CicoExecuteRequest{
				PartnerCode:      "DSA",
				PetnetTrackingno: "5a269417e107691f3d7c",
				TrxDate:          "2022-05-17",
			},
			want: &CicoExecuteResponse{
				Code:    200,
				Message: "Successful",
				Result: &CicoExecuteResult{
					PartnerCode:        "DSA",
					Provider:           "GCASH",
					PetnetTrackingno:   "5a269417e107691f3d7c",
					TrxDate:            "2022-05-17",
					TrxType:            "Cash In",
					ProviderTrackingno: "7000001521345",
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
			got, err := s.CicoExecute(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("CicoExecute() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

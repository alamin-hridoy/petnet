package perahub

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCicoInquire(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      CicoInquireRequest
		want    *CicoInquireResponse
		wantErr bool
	}{
		{
			name: "Success",
			in: CicoInquireRequest{
				PartnerCode:      "DSA",
				Provider:         "GCASH",
				TrxType:          "Cash In",
				ReferenceNumber:  "09654767706",
				PetnetTrackingno: "3115bc3f587d747cf8f5",
			},
			want: &CicoInquireResponse{
				Code:    200,
				Message: "Successful",
				Result: &CicoInquireResult{
					StatusMessage:    "SUCCESSFUL CASHIN",
					PetnetTrackingno: "238a8006885b57765cd8",
					TrxType:          "Cash In",
					ReferenceNumber:  "09654767706",
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.CicoInquire(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("CicoInquire() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

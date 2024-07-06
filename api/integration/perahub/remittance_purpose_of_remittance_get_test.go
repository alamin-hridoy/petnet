package perahub

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestPurposeOfRemittanceGet(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *PurposeOfRemittanceGetRes
		wantErr bool
	}{
		{
			name: "Success",
			want: &PurposeOfRemittanceGetRes{
				Code:    200,
				Message: "Good",
				Result: &PurposeOfRemittanceGetResult{
					ID:                  1,
					PurposeOfRemittance: "Donation",
					CreatedAt:           time.Now(),
					UpdatedAt:           time.Now(),
					DeletedAt:           time.Now(),
				},
			},
			wantErr: false,
		},
	}
	tOps := []cmp.Option{
		cmpopts.IgnoreFields(PurposeOfRemittanceGetResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.PurposeOfRemittanceGet(context.Background(), "1")
			if err != nil && !test.wantErr {
				t.Errorf("PurposeOfRemittanceGet() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

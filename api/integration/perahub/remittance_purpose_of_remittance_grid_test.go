package perahub

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestPurposeOfRemittanceGrid(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *PurposeOfRemittanceGridRes
		wantErr bool
	}{
		{
			name: "Success",
			want: &PurposeOfRemittanceGridRes{
				Code:    200,
				Message: "Good",
				Result: []PurposeOfRemittanceGridResult{
					{
						ID:                  "1",
						PurposeOfRemittance: "Gift",
						CreatedAt:           time.Now(),
						UpdatedAt:           time.Now(),
						DeletedAt:           time.Now(),
					},
					{
						ID:                  "2",
						PurposeOfRemittance: "Fund",
						CreatedAt:           time.Now(),
						UpdatedAt:           time.Now(),
						DeletedAt:           time.Now(),
					},
					{
						ID:                  "3",
						PurposeOfRemittance: "Allowance",
						CreatedAt:           time.Now(),
						UpdatedAt:           time.Now(),
						DeletedAt:           time.Now(),
					},
				},
			},
			wantErr: false,
		},
	}
	tOps := []cmp.Option{
		cmpopts.IgnoreFields(PurposeOfRemittanceGridResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.PurposeOfRemittanceGrid(context.Background())
			if err != nil && !test.wantErr {
				t.Errorf("PurposeOfRemittanceGrid() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

package perahub

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSourceOfFundGrid(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *SourceOfFundGridRes
		wantErr bool
	}{
		{
			name: "Success",
			want: &SourceOfFundGridRes{
				Code:    200,
				Message: "Good",
				Result: []SourceOfFundGridResult{
					{
						ID:           "1",
						SourceOfFund: "SALARY",
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
						DeletedAt:    time.Now(),
					},
					{
						ID:           "2",
						SourceOfFund: "BUSINESS",
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
						DeletedAt:    time.Now(),
					},
					{
						ID:           "3",
						SourceOfFund: "REMITTANCE",
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
						DeletedAt:    time.Now(),
					},
				},
			},
			wantErr: false,
		},
	}
	tOps := []cmp.Option{
		cmpopts.IgnoreFields(SourceOfFundGridResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.SourceOfFundGrid(context.Background())
			if err != nil && !test.wantErr {
				t.Errorf("SourceOfFundGrid() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

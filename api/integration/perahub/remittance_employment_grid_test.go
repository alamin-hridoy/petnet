package perahub

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRemittanceEmploymentGrid(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *RemittanceEmploymentGridRes
		wantErr bool
	}{
		{
			name: "Success",
			want: &RemittanceEmploymentGridRes{
				Code:    200,
				Message: "Good",
				Result: []RemittanceEmploymentGridResult{
					{
						ID:               1,
						EmploymentNature: "REGULAR",
						CreatedAt:        time.Now(),
						UpdatedAt:        time.Now(),
						DeletedAt:        time.Now(),
					},
					{
						ID:               2,
						EmploymentNature: "PROBATIONARY",
						CreatedAt:        time.Now(),
						UpdatedAt:        time.Now(),
						DeletedAt:        time.Now(),
					},
					{
						ID:               3,
						EmploymentNature: "CONTRACTUAL",
						CreatedAt:        time.Now(),
						UpdatedAt:        time.Now(),
						DeletedAt:        time.Now(),
					},
				},
			},
			wantErr: false,
		},
	}
	tOps := []cmp.Option{
		cmpopts.IgnoreFields(RemittanceEmploymentGridResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.RemittanceEmploymentGrid(context.Background())
			if err != nil && !test.wantErr {
				t.Errorf("RemittanceEmploymentGrid() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

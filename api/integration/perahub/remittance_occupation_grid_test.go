package perahub

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestOccupationGrid(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *OccupationGridRes
		wantErr bool
	}{
		{
			name: "Success",
			want: &OccupationGridRes{
				Code:    200,
				Message: "Good",
				Result: []OccupationGridResult{
					{
						ID:         1,
						Occupation: "Programmer",
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
						DeletedAt:  time.Now(),
					},
					{
						ID:         2,
						Occupation: "Engineer",
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
						DeletedAt:  time.Now(),
					},
					{
						ID:         3,
						Occupation: "Doctor",
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
						DeletedAt:  time.Now(),
					},
				},
			},
			wantErr: false,
		},
	}
	tOps := []cmp.Option{
		cmpopts.IgnoreFields(OccupationGridResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.OccupationGrid(context.Background())
			if err != nil && !test.wantErr {
				t.Errorf("OccupationGrid() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

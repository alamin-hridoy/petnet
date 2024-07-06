package perahub

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestOccupationGet(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *OccupationGetRes
		wantErr bool
	}{
		{
			name: "Success",
			want: &OccupationGetRes{
				Code:    200,
				Message: "Good",
				Result: &OccupationGetResult{
					ID:         1,
					Occupation: "Programmer",
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
					DeletedAt:  time.Now(),
				},
			},
			wantErr: false,
		},
	}
	tOps := []cmp.Option{
		cmpopts.IgnoreFields(OccupationGetResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.OccupationGet(context.Background(), "1")
			if err != nil && !test.wantErr {
				t.Errorf("OccupationGet() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

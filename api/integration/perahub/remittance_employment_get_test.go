package perahub

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRemittanceEmploymentGet(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *RemittanceEmploymentGetRes
		wantErr bool
	}{
		{
			name: "Success",
			want: &RemittanceEmploymentGetRes{
				Code:    200,
				Message: "Good",
				Result: &RemittanceEmploymentGetResult{
					ID:               1,
					EmploymentNature: "REGULAR",
					CreatedAt:        time.Now(),
					UpdatedAt:        time.Now(),
					DeletedAt:        time.Now(),
				},
			},
			wantErr: false,
		},
	}
	tOps := []cmp.Option{
		cmpopts.IgnoreFields(RemittanceEmploymentGetResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.RemittanceEmploymentGet(context.Background(), "1")
			if err != nil && !test.wantErr {
				t.Errorf("RemittanceEmploymentGet() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

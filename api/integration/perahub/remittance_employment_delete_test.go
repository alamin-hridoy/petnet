package perahub

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRemittanceEmploymentDelete(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *RemittanceEmploymentDeleteRes
		wantErr bool
	}{
		{
			name: "Success",
			want: &RemittanceEmploymentDeleteRes{
				Code:    200,
				Message: "Good",
				Result: &RemittanceEmploymentDeleteResult{
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
		cmpopts.IgnoreFields(RemittanceEmploymentDeleteResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.RemittanceEmploymentDelete(context.Background(), "1")
			if err != nil && !test.wantErr {
				t.Errorf("RemittanceEmploymentDelete() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

package perahub

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSourceOfFundDelete(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *SourceOfFundDeleteRes
		wantErr bool
	}{
		{
			name: "Success",
			want: &SourceOfFundDeleteRes{
				Code:    200,
				Message: "Good",
				Result: SourceOfFundDeleteResult{
					ID:           1,
					SourceOfFund: "SALARY",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
					DeletedAt:    time.Now(),
				},
			},
			wantErr: false,
		},
	}
	tOps := []cmp.Option{
		cmpopts.IgnoreFields(SourceOfFundDeleteResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.SourceOfFundDelete(context.Background(), "1")
			if err != nil && !test.wantErr {
				t.Errorf("SourceOfFundDelete() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

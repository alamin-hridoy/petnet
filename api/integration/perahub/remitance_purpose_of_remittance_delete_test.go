package perahub

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestPurposeOfRemittanceDelete(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *PurposeOfRemittanceDeleteRes
		wantErr bool
	}{
		{
			name: "Success",
			want: &PurposeOfRemittanceDeleteRes{
				Code:    200,
				Message: "Good",
				Result: PurposeOfRemittanceDeleteResult{
					ID:                  1,
					PurposeOfRemittance: "Gift",
					CreatedAt:           time.Now(),
					UpdatedAt:           time.Now(),
					DeletedAt:           time.Now(),
				},
			},
			wantErr: false,
		},
	}
	tOps := []cmp.Option{
		cmpopts.IgnoreFields(PurposeOfRemittanceDeleteResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.PurposeOfRemittanceDelete(context.Background(), "1")
			if err != nil && !test.wantErr {
				t.Errorf("PurposeOfRemittanceDelete() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

package perahub

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestPurposeOfRemittanceUpdate(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	req := PurposeOfRemittanceUpdateReq{
		PurposeOfRemittance: "Donation",
	}
	tests := []struct {
		name        string
		in          PurposeOfRemittanceUpdateReq
		expectedReq PurposeOfRemittanceUpdateReq
		want        *PurposeOfRemittanceUpdateRes
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &PurposeOfRemittanceUpdateRes{
				Code:    200,
				Message: "Good",
				Result: PurposeOfRemittanceUpdateResult{
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
		cmpopts.IgnoreFields(PurposeOfRemittanceUpdateResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.PurposeOfRemittanceUpdate(context.Background(), test.in, "1")
			if err != nil && !test.wantErr {
				t.Errorf("PurposeOfRemittanceUpdate() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq PurposeOfRemittanceUpdateReq
			if err := json.Unmarshal(m.GetMockRequest(), &newReq); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.expectedReq, newReq) {
				t.Error(cmp.Diff(test.expectedReq, newReq))
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

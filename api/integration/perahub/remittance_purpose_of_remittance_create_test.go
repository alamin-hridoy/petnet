package perahub

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestPurposeOfRemittanceCreate(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	req := PurposeOfRemittanceCreateReq{
		PurposeOfRemittance: "Donation",
	}
	tests := []struct {
		name        string
		in          PurposeOfRemittanceCreateReq
		expectedReq PurposeOfRemittanceCreateReq
		want        *PurposeOfRemittanceCreateRes
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &PurposeOfRemittanceCreateRes{
				Code:    200,
				Message: "Good",
				Result: PurposeOfRemittanceCreateResult{
					ID:                  1,
					PurposeOfRemittance: "USP",
					CreatedAt:           time.Now(),
					UpdatedAt:           time.Now(),
				},
			},
			wantErr: false,
		},
	}
	tOps := []cmp.Option{
		cmpopts.IgnoreFields(PurposeOfRemittanceCreateResult{}, "CreatedAt", "UpdatedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.PurposeOfRemittanceCreate(context.Background(), test.in)
			if err != nil && !test.wantErr {
				t.Errorf("PurposeOfRemittanceCreate() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq PurposeOfRemittanceCreateReq
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

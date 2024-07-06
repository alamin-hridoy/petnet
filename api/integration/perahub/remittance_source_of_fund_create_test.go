package perahub

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSourceOfFundCreateCreate(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	req := SourceOfFundCreateReq{
		SourceOfFund: "SALARY",
	}
	tests := []struct {
		name        string
		in          SourceOfFundCreateReq
		expectedReq SourceOfFundCreateReq
		want        *SourceOfFundCreateRes
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &SourceOfFundCreateRes{
				Code:    200,
				Message: "Good",
				Result: SourceOfFundCreateResult{
					ID:           1,
					SourceOfFund: "SALARY",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				},
			},
			wantErr: false,
		},
	}
	tOps := []cmp.Option{
		cmpopts.IgnoreFields(SourceOfFundCreateResult{}, "CreatedAt", "UpdatedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.SourceOfFundCreate(context.Background(), test.in)
			if err != nil && !test.wantErr {
				t.Errorf("SourceOfFundCreate() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq SourceOfFundCreateReq
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

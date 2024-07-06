package perahub

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRemittanceEmploymentUpdate(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	req := RemittanceEmploymentUpdateReq{
		Employment:       "REGULAR",
		EmploymentNature: "REGULAR",
	}
	tests := []struct {
		name        string
		in          RemittanceEmploymentUpdateReq
		expectedReq RemittanceEmploymentUpdateReq
		want        *RemittanceEmploymentUpdateRes
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &RemittanceEmploymentUpdateRes{
				Code:    200,
				Message: "Good",
				Result: RemittanceEmploymentUpdateResult{
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
		cmpopts.IgnoreFields(RemittanceEmploymentUpdateResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.RemittanceEmploymentUpdate(context.Background(), test.in, "1")
			if err != nil && !test.wantErr {
				t.Errorf("RemittanceEmploymentUpdate() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq RemittanceEmploymentUpdateReq
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

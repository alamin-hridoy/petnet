package perahub

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestOccupationUpdate(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	req := OccupationUpdateReq{
		Occupation: "Programmer",
	}
	tests := []struct {
		name        string
		in          OccupationUpdateReq
		expectedReq OccupationUpdateReq
		want        *OccupationUpdateRes
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &OccupationUpdateRes{
				Code:    200,
				Message: "Good",
				Result: OccupationUpdateResult{
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
		cmpopts.IgnoreFields(OccupationUpdateResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.OccupationUpdate(context.Background(), test.in, "1")
			if err != nil && !test.wantErr {
				t.Errorf("OccupationUpdate() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq OccupationUpdateReq
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

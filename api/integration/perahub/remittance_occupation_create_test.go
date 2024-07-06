package perahub

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestOccupationCreate(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	req := OccupationCreateReq{
		Occupation: "Programmer",
	}
	tests := []struct {
		name        string
		in          OccupationCreateReq
		expectedReq OccupationCreateReq
		want        *OccupationCreateRes
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &OccupationCreateRes{
				Code:    200,
				Message: "Good",
				Result: OccupationCreateResult{
					ID:         1,
					Occupation: "Programmer",
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				},
			},
			wantErr: false,
		},
	}
	tOps := []cmp.Option{
		cmpopts.IgnoreFields(OccupationCreateResult{}, "CreatedAt", "UpdatedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.OccupationCreate(context.Background(), test.in)
			if err != nil && !test.wantErr {
				t.Errorf("OccupationCreate() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq OccupationCreateReq
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

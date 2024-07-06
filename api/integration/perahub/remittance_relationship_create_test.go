package perahub

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRemittanceRelationshipCreate(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	req := RemittanceRelationshipCreateReq{
		Relationship: "Friend",
	}
	tests := []struct {
		name        string
		in          RemittanceRelationshipCreateReq
		expectedReq RemittanceRelationshipCreateReq
		want        *RemittanceRelationshipCreateRes
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &RemittanceRelationshipCreateRes{
				Code:    200,
				Message: "Good",
				Result: RemittanceRelationshipCreateResult{
					Relationship: "Friend",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
					ID:           1,
				},
			},
			wantErr: false,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.RemittanceRelationshipCreate(context.Background(), test.in)
			if err != nil && !test.wantErr {
				t.Errorf("RemittanceRelationshipCreate() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq RemittanceRelationshipCreateReq
			if err := json.Unmarshal(m.GetMockRequest(), &newReq); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.expectedReq, newReq) {
				t.Error(cmp.Diff(test.expectedReq, newReq))
			}
			tOps := []cmp.Option{
				cmpopts.IgnoreFields(RemittanceRelationshipCreateResult{}, "CreatedAt", "UpdatedAt"),
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

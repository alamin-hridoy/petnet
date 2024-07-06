package perahub

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRemittanceRelationshiptUpdate(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	req := RelationshipUpdateReq{
		Relationship: "Friend",
	}
	tests := []struct {
		name        string
		in          RelationshipUpdateReq
		expectedReq RelationshipUpdateReq
		want        *RelationshipUpdateRes
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &RelationshipUpdateRes{
				Code:    200,
				Message: "Good",
				Result: RelationshipUpdateResult{
					ID:           1,
					Relationship: "Friend",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
					DeletedAt:    time.Now(),
				},
			},
			wantErr: false,
		},
	}
	tOps := []cmp.Option{
		cmpopts.IgnoreFields(RelationshipUpdateResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.RelationshipUpdate(context.Background(), test.in, "1")
			if err != nil && !test.wantErr {
				t.Errorf("RelationshipUpdate() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq RelationshipUpdateReq
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

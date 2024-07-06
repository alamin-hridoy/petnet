package perahub

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRemittanceRelationshipGrid(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *RemittanceRelationshipGridRes
		wantErr bool
	}{
		{
			name: "Success",
			want: &RemittanceRelationshipGridRes{
				Code:    200,
				Message: "Good",
				Result: []RemittanceRelationshipGridResult{
					{
						ID:           1,
						Relationship: "Friend",
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
						DeletedAt:    time.Now(),
					},
					{
						ID:           2,
						Relationship: "Father",
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
						DeletedAt:    time.Now(),
					},
					{
						ID:           3,
						Relationship: "Mother",
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
						DeletedAt:    time.Now(),
					},
				},
			},
			wantErr: false,
		},
	}
	tOps := []cmp.Option{
		cmpopts.IgnoreFields(RemittanceRelationshipGridResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.RemittanceRelationshiptGrid(context.Background())
			if err != nil && !test.wantErr {
				t.Errorf("RemittanceRelationshiptGrid() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

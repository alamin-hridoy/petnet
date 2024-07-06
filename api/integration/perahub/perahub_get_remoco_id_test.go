package perahub

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPerahubGetRemcoID(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *PerahubGetRemcoIDResponse
		wantErr bool
	}{
		{
			name: "Success",
			want: &PerahubGetRemcoIDResponse{
				Code:    "200",
				Message: "Good",
				Result: []PerahubGetRemcoIDResult{
					{
						ID:   "1",
						Name: "iRemit",
					},
					{
						ID:   "2",
						Name: "BPI",
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.PerahubGetRemcoID(context.Background())
			if (err != nil) != test.wantErr {
				t.Fatalf("PerahubGetRemcoID() error = %v, wantErr %v", err, test.wantErr)
			}
			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

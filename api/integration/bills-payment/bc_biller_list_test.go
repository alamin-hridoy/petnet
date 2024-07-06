package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBCBillerList(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *BCBillerListResponse
		wantErr bool
	}{
		{
			name: "Success",
			want: &BCBillerListResponse{
				Code:    200,
				Message: "Success",
				Result: []BCBillerListResult{
					{
						Name:            "MERALCO",
						Code:            "MECOR",
						Description:     "Meralco Real-Time Posting",
						Category:        "Electricity",
						Type:            "RTP",
						Logo:            "https://stg-bc-api-images.s3-ap-southeast-1.amazonaws.com/biller-logos/250/MECOR.png",
						IsMultipleBills: 1,
						IsCde:           0,
						IsAsync:         1,
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
			got, err := s.BCBillerList(context.Background())
			if (err != nil) != test.wantErr {
				t.Fatalf("BCBillerList() error = %v, wantErr %v", err, test.wantErr)
			}

			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

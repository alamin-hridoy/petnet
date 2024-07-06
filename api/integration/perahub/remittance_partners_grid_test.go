package perahub

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRemittancePartnersGrid(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *RemittancePartnersGridRes
		wantErr bool
	}{
		{
			name: "Success",
			want: &RemittancePartnersGridRes{
				Code:    200,
				Message: "Good",
				Result: []RemittancePartnersGridResult{
					{
						ID:           1,
						PartnerCode:  "DRP",
						PartnerName:  "PERA HUB",
						ClientSecret: "26da230221d9e506b1fd823df1869875",
						Status:       1,
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
						DeletedAt:    time.Now(),
					},
					{
						ID:           2,
						PartnerCode:  "USP",
						PartnerName:  "PERA HUB",
						ClientSecret: "12358fbef0bb08d7a7bab57df956a335",
						Status:       1,
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
		cmpopts.IgnoreFields(RemittancePartnersGridResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.RemittancePartnersGrid(context.Background())
			if err != nil && !test.wantErr {
				t.Errorf("RemittancePartnersGrid() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

package perahub

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRemittancePartnersDelete(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      RemittancePartnersDeleteReq
		want    *RemittancePartnersDeleteRes
		wantErr bool
	}{
		{
			name: "Success",
			in: RemittancePartnersDeleteReq{
				ID:          "1",
				PartnerCode: "USP",
				PartnerName: "PERA HUB",
			},
			want: &RemittancePartnersDeleteRes{
				Code:    200,
				Message: "Good",
				Result: RemittancePartnersDeleteResult{
					ID:           1,
					PartnerCode:  "USP",
					PartnerName:  "PERA HUB",
					ClientSecret: "adawdawdawd",
					Status:       1,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
					DeletedAt:    time.Now(),
				},
			},
			wantErr: false,
		},
	}
	tOps := []cmp.Option{
		cmpopts.IgnoreFields(RemittancePartnersDeleteResult{}, "CreatedAt", "UpdatedAt", "DeletedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.RemittancePartnersDelete(context.Background(), test.in)
			if err != nil && !test.wantErr {
				t.Errorf("RemittancePartnersDelete() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

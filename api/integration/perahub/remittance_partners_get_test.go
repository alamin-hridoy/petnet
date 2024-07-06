package perahub

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRemittancePartnersGet(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		want    *RemittancePartnersGetRes
		wantErr bool
	}{
		{
			name: "Success",
			want: &RemittancePartnersGetRes{
				Code:    200,
				Message: "Good",
				Result: &RemittancePartnersGetResult{
					ID:           1,
					PartnerCode:  "DRP",
					PartnerName:  "BRANKAS",
					ClientSecret: "4fab1de660a6b7faef0168ca4788408a",
					Status:       1,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
					DeletedAt:    "null",
				},
			},
			wantErr: false,
		},
	}
	tOps := []cmp.Option{
		cmpopts.IgnoreFields(RemittancePartnersGetResult{}, "CreatedAt", "UpdatedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.RemittancePartnersGet(context.Background(), "1")
			if err != nil && !test.wantErr {
				t.Errorf("RemittancePartnersGet() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

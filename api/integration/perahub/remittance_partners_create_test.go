package perahub

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRemittancePartnersCreate(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	req := RemittancePartnersCreateReq{
		PartnerCode: "DRP",
		PartnerName: "BRANKAS",
	}
	tests := []struct {
		name        string
		in          RemittancePartnersCreateReq
		expectedReq RemittancePartnersCreateReq
		want        *RemittancePartnersCreateRes
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &RemittancePartnersCreateRes{
				Code:    200,
				Message: "Good",
				Result: RemittancePartnersCreateResult{
					ID:           1,
					PartnerCode:  "USP",
					PartnerName:  "PERA HUB",
					ClientSecret: "adawdawdawd",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
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
			got, err := s.RemittancePartnersCreate(context.Background(), test.in)
			if err != nil && !test.wantErr {
				t.Errorf("RemittancePartnersCreate() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq RemittancePartnersCreateReq
			if err := json.Unmarshal(m.GetMockRequest(), &newReq); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.expectedReq, newReq) {
				t.Error(cmp.Diff(test.expectedReq, newReq))
			}
			tOps := []cmp.Option{
				cmpopts.IgnoreFields(RemittancePartnersCreateResult{}, "CreatedAt", "UpdatedAt"),
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

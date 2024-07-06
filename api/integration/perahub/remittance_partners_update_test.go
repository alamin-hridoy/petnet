package perahub

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestRemittancePartnersUpdate(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	req := RemittancePartnersUpdateReq{
		PartnerCode: "USP",
		PartnerName: "PERA HUB",
	}
	tests := []struct {
		name        string
		in          RemittancePartnersUpdateReq
		expectedReq RemittancePartnersUpdateReq
		want        *RemittancePartnersUpdateRes
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &RemittancePartnersUpdateRes{
				Code:    200,
				Message: "Good",
				Result: RemittancePartnersUpdateResult{
					ID:           1,
					PartnerCode:  "USP",
					PartnerName:  "PERA HUB",
					ClientSecret: "adawdawdawd",
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
		cmpopts.IgnoreFields(RemittancePartnersUpdateResult{}, "CreatedAt", "UpdatedAt"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.RemittancePartnersUpdate(context.Background(), test.in, "1")
			if err != nil && !test.wantErr {
				t.Errorf("RemittancePartnersUpdate() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq RemittancePartnersUpdateReq
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

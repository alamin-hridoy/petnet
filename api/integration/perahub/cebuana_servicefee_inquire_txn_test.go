package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCebuanaSFInquiry(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          CebuanaSFInquiryRequest
		expectedReq CebuanaSFInquiryRequest
		want        *CebuanaSFInquiryRespBody
		wantErr     bool
	}{
		{
			name: "Success",
			in: CebuanaSFInquiryRequest{
				PrincipalAmount: "100",
				CurrencyID:      "6",
				AgentCode:       "01030063",
			},
			expectedReq: CebuanaSFInquiryRequest{
				PrincipalAmount: "100",
				CurrencyID:      "6",
				AgentCode:       "01030063",
			},
			want: &CebuanaSFInquiryRespBody{
				Code:    "0",
				Message: "Successful",
				Result: CebuanaSFInquiryResult{
					ResultStatus: "Successful",
					MessageID:    "0",
					LogID:        "0",
					ServiceFee:   "1.00",
				},
				RemcoID: "9",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.CebuanaSFInquiry(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("CebuanaSFInquiry() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq CebuanaSFInquiryRequest
			if err := json.Unmarshal(m.GetMockRequest(), &newReq); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.expectedReq, newReq) {
				t.Error(cmp.Diff(test.expectedReq, newReq))
			}

			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}

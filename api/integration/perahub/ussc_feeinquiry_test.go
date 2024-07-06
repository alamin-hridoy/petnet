package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUSSCFeeInquiry(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	req := USSCFeeInquiryRequest{
		Panalokard: "",
		Amount:     "3424.33",
		USSCPromo:  "",
		BranchCode: "branch1",
	}

	tests := []struct {
		name        string
		in          USSCFeeInquiryRequest
		expectedReq USSCFeeInquiryRequest
		want        *USSCFeeInquiryRespBody
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &USSCFeeInquiryRespBody{
				Code:    "200",
				Message: "OK",
				Result: USSCFeeInquiryResult{
					PnplAmount:    "1000.00",
					ServiceCharge: "1.00",
					Msg:           "",
					Code:          "0",
					NewScreen:     "0",
					JournalNo:     "000000202",
					ProcessDate:   "null",
					RefNo:         "1",
					TotAmount:     "1001.00",
					SendOTP:       "Y",
				},
				RemcoID: 10,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.USSCFeeInquiry(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("USSCFeeInquiry() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq USSCFeeInquiryRequest
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

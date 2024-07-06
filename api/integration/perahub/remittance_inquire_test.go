package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemitanceInquire(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	req := RemitanceInquireReq{
		Phrn: "PH1654789142",
	}
	tests := []struct {
		name        string
		in          RemitanceInquireReq
		expectedReq RemitanceInquireReq
		want        *RemitanceInquireRes
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &RemitanceInquireRes{
				Code:    200,
				Message: "PeraHUB Reference Number (PHRN) is available for Payout",
				Result: RemitanceInquireResult{
					Phrn:                  "PH1658296732",
					PrincipalAmount:       10000,
					IsoCurrency:           "PHP",
					ConversionRate:        1,
					IsoOriginatingCountry: "PHP",
					IsoDestinationCountry: "PHP",
					SenderLastName:        "HERMO",
					SenderFirstName:       "IRENE",
					SenderMiddleName:      "M",
					ReceiverLastName:      "HERMO",
					ReceiverFirstName:     "SONNY",
					ReceiverMiddleName:    "D",
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
			got, err := s.RemitanceInquire(context.Background(), test.in)
			if err != nil && !test.wantErr {
				t.Errorf("RemitanceInquire() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq RemitanceInquireReq
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

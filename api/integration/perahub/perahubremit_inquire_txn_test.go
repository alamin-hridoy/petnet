package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPerahubRemitInquire(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          PerahubRemitInquireRequest
		expectedReq PerahubRemitInquireRequest
		want        *PerahubRemitInquireResponse
		wantErr     bool
	}{
		{
			name: "Success",
			in: PerahubRemitInquireRequest{
				ControlNumber: "PH1655176065",
				LocationID:    371,
			},
			expectedReq: PerahubRemitInquireRequest{
				ControlNumber: "PH1655176065",
				LocationID:    371,
			},
			want: &PerahubRemitInquireResponse{
				Code:    200,
				Message: "PeraHUB Reference Number (PHRN) is available for Payout",
				Result: PerahubRemitInquireResult{
					PrincipalAmount:    179,
					IsoCurrency:        "PHP",
					ConversionRate:     1,
					SenderLastName:     "Sauer",
					SenderFirstName:    "Mittie",
					SenderMiddleName:   "O",
					ReceiverLastName:   "Soto",
					ReceiverFirstName:  "Blanche",
					ReceiverMiddleName: "G",
					ControlNumber:      "PH1655176065",
					OriginatingCountry: "PH",
					DestinationCountry: "PH",
					SenderName:         "Sauer, Mittie, O",
					ReceiverName:       "Soto, Blanche, G",
					PartnerCode:        "DRP",
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s, m := newTestSvc(t, st)
			got, err := s.PerahubRemitInquire(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("PerahubRemitInquire() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq PerahubRemitInquireRequest
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

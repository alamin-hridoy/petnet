package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRiaInquire(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          RiaInquireRequest
		expectedReq RiaInquireRequest
		want        *RiaInquireResponseBody
		wantErr     bool
	}{
		{
			name: "Success",
			in: RiaInquireRequest{
				DeviceID:        "desktop123",
				AgentID:         "84424911",
				AgentCode:       "ITD",
				ReferenceNumber: "TEST0987654321TEST12",
				ControlNumber:   "13020065349",
				LocationID:      "191",
				UserID:          "1893",
				LocationName:    "loc-name",
			},
			expectedReq: RiaInquireRequest{
				DeviceID:        "desktop123",
				AgentID:         "84424911",
				AgentCode:       "ITD",
				ReferenceNumber: "TEST0987654321TEST12",
				ControlNumber:   "13020065349",
				LocationID:      "191",
				UserID:          "1893",
				LocationName:    "loc-name",
			},
			want: &RiaInquireResponseBody{
				Code: "200",
				Msg:  "Order is available for payout.",
				Result: RiaResult{
					ControlNumber:      "CTRL1",
					ClientReferenceNo:  "1",
					OriginatingCountry: "TH",
					DestinationCountry: "PH",
					SenderName:         "John, Michael Doe",
					ReceiverName:       "Jane, Emily Doe",
					PrincipalAmount:    "1000.00",
					Currency:           "PHP",
					IsDomestic:         "1",
					OrderNo:            "TH1950882455",
				},
				RemcoID: "12",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s, m := newTestSvc(t, st)
			got, err := s.RiaInquire(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("RiaInquire() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq RiaInquireRequest
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

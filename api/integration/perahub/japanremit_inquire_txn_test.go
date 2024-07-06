package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestJPRInquire(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          JPRInquireRequest
		expectedReq JPRInquireRequest
		want        *JPRInquireResponseBody
		wantErr     bool
	}{
		{
			name: "Success",
			in: JPRInquireRequest{
				Branch:          "Information Technology Department",
				OutletCode:      "HOA",
				ReferenceNumber: "REF1",
				ControlNumber:   "CTRL1",
				LocationID:      "371",
				UserID:          5188,
				LocationName:    "Information Technology Department",
				DeviceId:        "a2c91f7155d18ba34c08782c86b07338b2fd1a6f06e30d1243",
				AgentId:         "84424911",
				AgentCode:       "HOA",
				BranchCode:      "00371",
				LocationCode:    "HOA",
				Currency:        "PHP",
			},
			expectedReq: JPRInquireRequest{
				Branch:          "Information Technology Department",
				OutletCode:      "HOA",
				ReferenceNumber: "REF1",
				ControlNumber:   "CTRL1",
				LocationID:      "371",
				UserID:          5188,
				LocationName:    "Information Technology Department",
				DeviceId:        "a2c91f7155d18ba34c08782c86b07338b2fd1a6f06e30d1243",
				AgentId:         "84424911",
				AgentCode:       "HOA",
				BranchCode:      "00371",
				LocationCode:    "HOA",
				Currency:        "PHP",
			},
			want: &JPRInquireResponseBody{
				Code:    "0",
				Message: "Available For Pickup",
				Result: JPRResult{
					ControlNumber:      "CTRL1",
					ReferenceNumber:    "REF1",
					OriginatingCountry: "UNITED ARAB EMIRATES",
					DestinationCountry: "PHILIPPINES",
					SenderName:         "John Michael Doe",
					ReceiverName:       "Jane Emily Doe",
					PrincipalAmount:    "1000.00",
					Currency:           "PHP",
					PayTokenId:         "729764",
				},
				RemcoID: "17",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s, m := newTestSvc(t, st)
			got, err := s.JPRInquire(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("JPRInquire() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq JPRInquireRequest
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

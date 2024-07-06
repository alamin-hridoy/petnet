package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInstaCashInquire(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          InstaCashInquireRequest
		expectedReq InstaCashInquireRequest
		want        *InstaCashInquireResponseBody
		wantErr     bool
	}{
		{
			name: "Success",
			in: InstaCashInquireRequest{
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
			},
			expectedReq: InstaCashInquireRequest{
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
			},
			want: &InstaCashInquireResponseBody{
				Code:    "1",
				Message: "Client Details",
				Result: InstaCashResult{
					ControlNumber:      "CTRL1",
					ReferenceNumber:    "REF1",
					OriginatingCountry: "UNITED ARAB EMIRATES",
					DestinationCountry: "PHILIPPINES",
					SenderName:         "John Michael Doe",
					ReceiverName:       "Jane Emily Doe",
					PrincipalAmount:    "1000.00",
					Currency:           "PHP",
					Purpose:            "FAMILY MAINTENANCE",
					Status:             "For Pick Up",
				},
				RemcoID: "16",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s, m := newTestSvc(t, st)
			got, err := s.InstaCashInquire(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("InstaCashInquire() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq InstaCashInquireRequest
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

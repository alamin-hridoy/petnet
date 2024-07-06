package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAYANNAHInquire(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          AYANNAHInquireRequest
		expectedReq AYANNAHInquireRequest
		want        *AYANNAHInquireResponseBody
		wantErr     bool
	}{
		{
			name: "Success",
			in: AYANNAHInquireRequest{
				Branch:          "Information Technology Department",
				OutletCode:      "HOA",
				ReferenceNumber: "1",
				ControlNumber:   "CTRL1",
				LocationID:      "371",
				UserID:          "5188",
				LocationName:    "Information Technology Department",
				DeviceID:        "044379868a90415aab26e8a469ed028075c15a7392aade4343",
				AgentCode:       "HOA",
				LocationCode:    "PERAH-00001",
				Currency:        "PHP",
			},
			expectedReq: AYANNAHInquireRequest{
				Branch:          "Information Technology Department",
				OutletCode:      "HOA",
				ReferenceNumber: "1",
				ControlNumber:   "CTRL1",
				LocationID:      "371",
				UserID:          "5188",
				LocationName:    "Information Technology Department",
				DeviceID:        "044379868a90415aab26e8a469ed028075c15a7392aade4343",
				AgentCode:       "HOA",
				LocationCode:    "PERAH-00001",
				Currency:        "PHP",
			},
			want: &AYANNAHInquireResponseBody{
				Code:    "200",
				Message: "Success",
				Result: AYANNAHInquireResult{
					ResponseCode:       "AVAILABLE",
					ResponseMessage:    "Transaction Available For Payout.",
					ControlNumber:      "CTRL1",
					PrincipalAmount:    "1000.00",
					Currency:           "PHP",
					CreationDate:       "2022-01-14 14:57:26 ",
					ReceiverName:       "Octaviano Luis Rafael",
					SenderName:         "John, Michael Doe",
					Address:            "IMORTALONE",
					City:               "ANTIPOLO,RIZAL",
					Country:            "PH",
					ZipCode:            "null",
					OriginatingCountry: "PH",
					DestinationCountry: "PH",
					ContactNumber:      "null",
					ReferenceNumber:    "1",
				},
				RemcoID: "22",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.AYANNAHInquire(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("AYANNAHInquire() error = %v, wantErr %v", err, test.wantErr)
			}
			var newReq AYANNAHInquireRequest
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

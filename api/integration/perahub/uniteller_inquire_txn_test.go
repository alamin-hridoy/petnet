package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUNTInquire(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          UNTInquireRequest
		expectedReq UNTInquireRequest
		want        *UNTInquireResponseBody
		wantErr     bool
	}{
		{
			name: "Success",
			in: UNTInquireRequest{
				Branch:          "Information Technology Department",
				OutletCode:      "HOA",
				ReferenceNumber: "20211022PHB016307325",
				ControlNumber:   "800118975",
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
			expectedReq: UNTInquireRequest{
				Branch:          "Information Technology Department",
				OutletCode:      "HOA",
				ReferenceNumber: "20211022PHB016307325",
				ControlNumber:   "800118975",
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
			want: &UNTInquireResponseBody{
				Code:    "00000000",
				Message: "Success",
				Result: UNTResult{
					ResponseCode:       "00000000",
					ControlNumber:      "CTRL1",
					PrincipalAmount:    "1000.00",
					Currency:           "PHP",
					CreationDate:       "2021-06-15T19:10:16.000-0400",
					ReceiverName:       "Jane Emily Doe",
					Address:            "TEST",
					City:               "MORONG",
					Country:            "PH",
					SenderName:         "John Michael Doe",
					ZipCode:            "36978",
					OriginatingCountry: "US",
					DestinationCountry: "PH",
					ContactNumber:      "1540254852",
					FmtSenderName:      "John, Michael Doe",
					FmtReceiverName:    "Jane, Emily Doe",
				},
				RemcoID: "20",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s, m := newTestSvc(t, st)
			got, err := s.UNTInquire(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("UNTInquire() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq UNTInquireRequest
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

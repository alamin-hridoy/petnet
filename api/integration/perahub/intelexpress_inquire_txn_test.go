package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIEInquire(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          IEInquireRequest
		expectedReq IEInquireRequest
		want        *IEInquireResponse
		wantErr     bool
	}{
		{
			name: "Success",
			in: IEInquireRequest{
				Branch:          "Information Technology Department",
				OutletCode:      "HOA",
				ReferenceNumber: "19564301",
				ControlNumber:   "07K84790",
				LocationID:      "371",
				UserID:          "1893",
				LocationName:    "MALOLOS",
				DeviceID:        "044379868a90415aab26e8a469ed028075c15a7392aade4343",
				AgentID:         "84424911",
				AgentCode:       "HOA",
				BranchCode:      "00371",
				LocationCode:    "PERAH-00001",
				Currency:        "PHP",
			},
			expectedReq: IEInquireRequest{
				Branch:          "Information Technology Department",
				OutletCode:      "HOA",
				ReferenceNumber: "19564301",
				ControlNumber:   "07K84790",
				LocationID:      "371",
				UserID:          "1893",
				LocationName:    "MALOLOS",
				DeviceID:        "044379868a90415aab26e8a469ed028075c15a7392aade4343",
				AgentID:         "84424911",
				AgentCode:       "HOA",
				BranchCode:      "00371",
				LocationCode:    "PERAH-00001",
				Currency:        "PHP",
			},
			want: &IEInquireResponse{
				Code:    "200",
				Message: "Success",
				Result: IEInquireResult{
					ControlNumber:      "CTRL1",
					TrxDate:            "07/03/2021",
					PrincipalAmount:    "1000.00",
					Currency:           "PHP",
					ReceiverName:       "Jane Emily Doe",
					SenderName:         "John Michael Doe",
					Address:            "TBILISI,KOKHREIDZIS - 8",
					Country:            "GEO",
					OriginatingCountry: "GEO",
					DestinationCountry: "PHL",
					ContactNumber:      "",
					ReferenceNumber:    "1",
				},
				RemcoID: "24",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.IEInquire(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("IEInquire() error = %v, wantErr %v", err, test.wantErr)
			}
			var newReq IEInquireRequest
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

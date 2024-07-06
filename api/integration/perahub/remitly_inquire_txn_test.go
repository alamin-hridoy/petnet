package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRMInquire(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          RMInquireRequest
		expectedReq RMInquireRequest
		want        *RMInquireResponseBody
		wantErr     bool
	}{
		{
			name: "Success",
			in: RMInquireRequest{
				Branch:       "branch",
				OutletCode:   "outlet-code",
				RefNo:        "refno",
				ControlNo:    "cntrl-no",
				LocationID:   "1",
				UserID:       "2",
				LocationName: "loc-name",
				DeviceID:     "dv-id",
				AgentID:      "ag-id",
				AgentCode:    "ag-code",
				BranchCode:   "b-code",
				LocationCode: "loc-code",
				CurrencyCode: "cur-code",
			},
			expectedReq: RMInquireRequest{
				Branch:       "branch",
				OutletCode:   "outlet-code",
				RefNo:        "refno",
				ControlNo:    "cntrl-no",
				LocationID:   "1",
				UserID:       "2",
				LocationName: "loc-name",
				DeviceID:     "dv-id",
				AgentID:      "ag-id",
				AgentCode:    "ag-code",
				BranchCode:   "b-code",
				LocationCode: "loc-code",
				CurrencyCode: "cur-code",
			},
			want: &RMInquireResponseBody{
				Code: "200",
				Msg:  "PAYABLE",
				Result: RMInqResult{
					ControlNo:     "CTRL1",
					PnplAmt:       "1000.00",
					SenderName:    "John, Michael Doe",
					RcvName:       "Jane, Emily Doe",
					Address:       "7118 Street",
					CurrencyCode:  "PHP",
					ContactNumber: "9162427505",
					City:          "Manila",
					Country:       "PHL",
				},
				RemcoID: "21",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.RMInquire(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("RMInquire() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			var newReq RMInquireRequest
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

package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBPInquire(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          BPInquireRequest
		expectedReq BPInquireRequest
		want        *BPInquireResponseBody
		wantErr     bool
	}{
		{
			name: "Success",
			in: BPInquireRequest{
				RefNo:        "1",
				ControlNo:    "CTRL1",
				LocationID:   "371",
				UserID:       "5188",
				TrxDate:      "2021-10-11",
				LocationName: "IT Department",
			},
			expectedReq: BPInquireRequest{
				RefNo:        "1",
				ControlNo:    "CTRL1",
				LocationID:   "371",
				UserID:       "5188",
				TrxDate:      "2021-10-11",
				LocationName: "IT Department",
			},
			want: &BPInquireResponseBody{
				Code: "200",
				Msg:  "IN PROCESS:TRANSACTION PROCESS ONGOING",
				Result: BPInquireResult{
					Status:            "T",
					Desc:              "TRANSMIT",
					ControlNo:         "CTRL1",
					RefNo:             "1",
					ClientReferenceNo: "CL1",
					PnplAmt:           "1000.00",
					SenderName:        "John, Michael Doe",
					RcvName:           "Jane, Emily Doe",
					Address:           "PLA",
					Currency:          "PHP",
					ContactNumber:     "09190000000",
					OrgnCtry:          "SINGAPORE",
					DestCtry:          "PHILIPPINES",
				},
				RemcoID: "2",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.BPInquire(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("BPInquire() error = %v, wantErr %v", err, test.wantErr)
			}
			var newReq BPInquireRequest
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

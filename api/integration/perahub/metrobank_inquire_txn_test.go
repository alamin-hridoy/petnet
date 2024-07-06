package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMBInquire(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          MBInquireRequest
		expectedReq MBInquireRequest
		want        *MBInquireResponseBody
		wantErr     bool
	}{
		{
			name: "Success",
			in: MBInquireRequest{
				RefNo:        "refno",
				ControlNo:    "cntrl-no",
				LocationID:   "1",
				UserID:       "1",
				LocationName: "loc-name",
			},
			expectedReq: MBInquireRequest{
				RefNo:        "refno",
				ControlNo:    "cntrl-no",
				LocationID:   "1",
				UserID:       "1",
				LocationName: "loc-name",
			},
			want: &MBInquireResponseBody{
				Code: "0",
				Msg:  "Available for pick-up",
				Result: MBInqResult{
					RefNo:           "REF1",
					ControlNo:       "CTRL1",
					StatusText:      "0",
					PrincipalAmount: "1000.00",
					RcvName:         "Jane, Emily Doe",
					Address:         "7118 Street",
					ContactNumber:   "9162427505",
					Currency:        "PHP",
				},
				RemcoID: "8",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.MBInquire(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("MBInquire() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq MBInquireRequest
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

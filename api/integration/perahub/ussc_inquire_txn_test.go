package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUSSCInquire(t *testing.T) {
	t.Parallel()
	req := USSCInquireRequest{
		RefNo:        "TEST0987654321TEST12",
		ControlNo:    "SP323762267861",
		LocationID:   "191",
		UserID:       "1893",
		LocationName: "MALOLOS",
		BranchCode:   "MAL",
	}
	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          USSCInquireRequest
		expectedReq USSCInquireRequest
		want        *USSCInquireResponseBody
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &USSCInquireResponseBody{
				Code: "000000",
				Msg:  "OK",
				Result: USSCInqResult{
					RcvName:            "Iglesia, Julius James",
					ControlNo:          "CTRL1",
					PrincipalAmount:    "1000.00",
					ContactNumber:      "0922261616161",
					RefNo:              "1",
					SenderName:         "John Michael Doe",
					TrxDate:            "20211214",
					SenderLastName:     "Doe",
					SenderFirstName:    "John",
					SenderMiddleName:   "Michael",
					ServiceCharge:      "1.00",
					TotalAmount:        "1001.00",
					ReceiverFirstName:  "CAMITAN",
					ReceiverMiddleName: "ALVIN",
					ReceiverLastName:   "JOMAR TE TEST",
					RelationTo:         "Family",
					PurposeTransaction: "Family Support/Living Expenses",
				},
				RemcoID: "10",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.USSCInquire(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("USSCInquire() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq USSCInquireRequest
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

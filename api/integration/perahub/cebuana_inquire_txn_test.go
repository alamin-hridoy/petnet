package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCEBInquire(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          CEBInquireRequest
		expectedReq CEBInquireRequest
		want        *CEBInquireResponseBody
		wantErr     bool
	}{
		{
			name: "Success",
			in: CEBInquireRequest{
				ControlNo:    "CTRL1",
				LocationID:   "191",
				UserID:       "1893",
				LocationName: "MALOLOS",
			},
			expectedReq: CEBInquireRequest{
				ControlNo:    "CTRL1",
				LocationID:   "191",
				UserID:       "1893",
				LocationName: "MALOLOS",
			},
			want: &CEBInquireResponseBody{
				Code:    "1",
				Message: "Successful",
				Result: CEBInquireResult{
					ResultStatus:      "Successful",
					MessageID:         "125",
					LogID:             "0",
					ClientReferenceNo: "1",
					ControlNo:         "CTRL1",
					SenderName:        "John Michael Doe",
					RcvName:           "ESTOCAPIO, SHAIRA MIKA, MADJALIS",
					PnplAmt:           "1000",
					ServiceCharge:     "1",
					BirthDate:         "1980-08-10T00:00:00",
					Currency:          "PHP",
					BeneficiaryID:     "5342",
					RemStatusID:       "1",
					RemStatusDes:      "Outstanding",
				},
				RemcoID: "9",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.CEBInquire(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("CEBInquire() error = %v, wantErr %v", err, test.wantErr)
			}
			var newReq CEBInquireRequest
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

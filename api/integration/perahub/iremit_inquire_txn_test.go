package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIRInquire(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          IRInquireRequest
		expectedReq IRInquireRequest
		want        *IRInquireResponseBody
		wantErr     bool
	}{
		{
			name: "Success",
			in: IRInquireRequest{
				Branch:       "branch",
				RefNo:        "REF1",
				ControlNo:    "CTRL1",
				LocationID:   "33",
				UserID:       "44",
				LocationName: "loc-name",
			},
			expectedReq: IRInquireRequest{
				Branch:       "branch",
				RefNo:        "REF1",
				ControlNo:    "CTRL1",
				LocationID:   "33",
				UserID:       "44",
				LocationName: "loc-name",
			},
			want: &IRInquireResponseBody{
				Code: "0",
				Msg:  "Available for Pick-up",
				Result: IRResult{
					Status:        "0",
					Desc:          "Available for Pick-up",
					ControlNo:     "CTRL1",
					RefNo:         "REF1",
					PnplAmt:       "1000.00",
					SenderName:    "John, Michael Doe",
					RcvName:       "Jane, Emily Doe",
					Address:       "PLA",
					CurrencyCode:  "PHP",
					ContactNumber: "09190000000",
					RcvLastName:   "Doe",
					RcvFirstName:  "Jane",
				},
				RemcoID: "1",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s, m := newTestSvc(t, st)
			got, err := s.IRemitInquire(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("IRInquire() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			var newReq IRInquireRequest
			if err := json.Unmarshal(m.GetMockRequest(), &newReq); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.expectedReq, newReq) {
				t.Error(cmp.Diff(test.expectedReq, newReq))
			}

			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
			if got.Result.Address != test.want.Result.Address {
				t.Errorf("Result.Address, want: %s, got: %s", test.want.Result.Address, got.Result.Address)
			}
			if got.Result.ContactNumber != test.want.Result.ContactNumber {
				t.Errorf("Result.ContactNumber, want: %s, got: %s", test.want.Result.ContactNumber, got.Result.ContactNumber)
			}
		})
	}
}

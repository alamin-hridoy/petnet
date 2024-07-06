package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTFInquire(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          TFInquireRequest
		expectedReq TFInquireRequest
		want        *TFInquireResponseBody
		wantErr     bool
	}{
		{
			name: "Success",
			in: TFInquireRequest{
				Branch:       "branch",
				RefNo:        "11",
				ControlNo:    "22",
				LocationID:   "33",
				UserID:       "44",
				LocationName: "loc-name",
			},
			expectedReq: TFInquireRequest{
				Branch:       "branch",
				RefNo:        "11",
				ControlNo:    "22",
				LocationID:   "33",
				UserID:       "44",
				LocationName: "loc-name",
			},
			want: &TFInquireResponseBody{
				Code: "200",
				Msg:  "Success!",
				Result: TFResult{
					Status:        "T",
					Desc:          "TRANSMIT",
					ControlNo:     "CTRL1",
					RefNo:         "1",
					PnplAmt:       "1000.00",
					SenderName:    "John, Michael Doe",
					RcvName:       "Jane, Emily Doe",
					Address:       "PLA",
					CurrencyCode:  "PHP",
					ContactNumber: "09190000000",
					RcvLastName:   "Doe",
					RcvFirstName:  "Jane",
					OrgnCtry:      "UNITED ARAB EMIRATES",
					DestCtry:      "PHILIPPINES",
					TxnDate:       "2021-07-15T13:26:01.493-04:00",
					IsDomestic:    "0",
					IDType:        "0",
					RcvCtryCode:   "PH",
					RcvStateID:    "PH023",
					RcvStateName:  "METRO MANILA",
					RcvCityID:     "1",
					RcvCityName:   "METRO MANILA",
					RcvIDType:     "1",
					RcvIsIndiv:    "True",
					PrpsOfRmtID:   "1",
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
			got, err := s.TFInquire(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("TFInquire() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq TFInquireRequest
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

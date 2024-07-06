package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCEBINTInquire(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          CEBINTInquireRequest
		expectedReq CEBINTInquireRequest
		want        *CEBINTInquireResponseBody
		wantErr     bool
	}{
		{
			name: "Success",
			in: CEBINTInquireRequest{
				ControlNumber:            "CTRL1",
				LocationID:               "191",
				UserID:                   "1893",
				LocationName:             "MALOLOS",
				InternationalPartnerCode: "PNG-CQ",
				DeviceID:                 "722231cee610d10f5c2a1ffb9fa54f2c691f0100829b2b0a2c",
				AgentID:                  "84424911",
				AgentCode:                "01030063",
				BranchCode:               "00292",
				LocationCode:             "MAL",
				Branch:                   "MALOLOS",
				OutletCode:               "MAL",
				ReferenceNumber:          "REF1",
			},
			expectedReq: CEBINTInquireRequest{
				ControlNumber:            "CTRL1",
				LocationID:               "191",
				UserID:                   "1893",
				LocationName:             "MALOLOS",
				InternationalPartnerCode: "PNG-CQ",
				DeviceID:                 "722231cee610d10f5c2a1ffb9fa54f2c691f0100829b2b0a2c",
				AgentID:                  "84424911",
				AgentCode:                "01030063",
				BranchCode:               "00292",
				LocationCode:             "MAL",
				Branch:                   "MALOLOS",
				OutletCode:               "MAL",
				ReferenceNumber:          "REF1",
			},
			want: &CEBINTInquireResponseBody{
				Code:    "1",
				Message: "Successful",
				Result: CEBINTInquireResult{
					IsDomestic:                  "0",
					ResultStatus:                "Successful",
					MessageID:                   "0",
					LogID:                       "0",
					ClientReferenceNo:           "1",
					ControlNumber:               "CTRL1",
					SenderName:                  "John Michael Doe",
					ReceiverName:                "ESTOCAPIO, SHAIRA MIKA, MADJALIS",
					PrincipalAmount:             "1000",
					ServiceCharge:               "1",
					BirthDate:                   "1980-08-10T00:00:00",
					Currency:                    "PHP",
					BeneficiaryID:               "5342",
					RemittanceStatusID:          "1",
					RemittanceStatusDescription: "Outstanding",
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
			got, err := s.CEBINTInquire(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("CEBINTInquire() error = %v, wantErr %v", err, test.wantErr)
			}
			var newReq CEBINTInquireRequest
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

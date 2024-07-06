package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemitanceConfirmReceiveMoney(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	req := RemitanceConfirmReceiveMoneyReq{
		PayoutValidateReferenceNumber: "4f8a09d3b293807aa50305f66d6cc73c",
	}
	tests := []struct {
		name        string
		in          RemitanceConfirmReceiveMoneyReq
		expectedReq RemitanceConfirmReceiveMoneyReq
		want        *RemitanceConfirmReceiveMoneyRes
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &RemitanceConfirmReceiveMoneyRes{
				Code:    200,
				Message: "Successful",
				Result: RemitanceConfirmReceiveMoneyResult{
					Phrn:                  "PH1654789142",
					PrincipalAmount:       10000,
					IsoOriginatingCountry: "PHP",
					IsoDestinationCountry: "PHP",
					SenderLastName:        "HERMO",
					SenderFirstName:       "IRENE",
					SenderMiddleName:      "M",
					ReceiverLastName:      "HERMO",
					ReceiverFirstName:     "SONNY",
					ReceiverMiddleName:    "D",
				},
			},
			wantErr: false,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.RemitanceConfirmReceiveMoney(context.Background(), test.in)
			if err != nil && !test.wantErr {
				t.Errorf("RemitanceConfirmReceiveMoney() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq RemitanceConfirmReceiveMoneyReq
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

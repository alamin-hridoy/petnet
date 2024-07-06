package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemitanceCancelSendMoney(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	req := RemitanceCancelSendMoneyReq{
		Phrn:        "PH1654789142",
		PartnerCode: "DRP",
		Remarks:     "Sample cancel send",
	}
	tests := []struct {
		name        string
		in          RemitanceCancelSendMoneyReq
		expectedReq RemitanceCancelSendMoneyReq
		want        *RemitanceCancelSendMoneyRes
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &RemitanceCancelSendMoneyRes{
				Code:    200,
				Message: "Cancel Send Remittance successful",
				Result: RemitanceCancelSendMoneyResult{
					Phrn:                      "PH1654789142",
					CancelSendDate:            "2022-06-14 19:35:02",
					CancelSendReferenceNumber: "6a6e74400561a300d627aba12107bb6c",
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
			got, err := s.RemitanceCancelSendMoney(context.Background(), test.in)
			if err != nil && !test.wantErr {
				t.Errorf("RemitanceCancelSendMoney() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq RemitanceCancelSendMoneyReq
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

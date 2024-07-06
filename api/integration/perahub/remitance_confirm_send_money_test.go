package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRemitanceConfirmSendMoney(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	req := RemitanceConfirmSendMoneyReq{
		SendValidateReferenceNumber: "084396e59ed3a4acc8da4bd8885bec01",
	}
	tests := []struct {
		name        string
		in          RemitanceConfirmSendMoneyReq
		expectedReq RemitanceConfirmSendMoneyReq
		want        *RemitanceConfirmSendMoneyRes
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &RemitanceConfirmSendMoneyRes{
				Code:    200,
				Message: "Successful",
				Result: RemitanceConfirmSendMoneyResult{
					Phrn: "PH1654787564",
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
			got, err := s.RemitanceConfirmSendMoney(context.Background(), test.in)
			if err != nil && !test.wantErr {
				t.Errorf("RemitanceConfirmSendMoney() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq RemitanceConfirmSendMoneyReq
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

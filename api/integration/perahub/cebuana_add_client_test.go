package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCebAddClient(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name        string
		in          CebAddClientReq
		expectedReq CebAddClientReq
		want        *CebAddClientResp
		wantErr     bool
	}{
		{
			name: "Success",
			in: CebAddClientReq{
				FirstName: "John",
				LastName:  "Doe",
				BirthDate: "1981-08-12",
				CpCtryID:  "63",
				ContactNo: "63",
				TpCtryID:  "0",
				TpArCode:  "",
				CrtyAdID:  "166",
				PAdd:      "1953 Ph 3b Block 6 Lot 9 Camarin Caloocan City",
				CAdd:      "2863",
				UserID:    "2863",
				SOFID:     "1",
				Tin:       "",
				TpNo:      "",
				AgentCode: "01030063",
			},
			expectedReq: CebAddClientReq{
				FirstName: "John",
				LastName:  "Doe",
				BirthDate: "1981-08-12",
				CpCtryID:  "63",
				ContactNo: "63",
				TpCtryID:  "0",
				TpArCode:  "",
				CrtyAdID:  "166",
				PAdd:      "1953 Ph 3b Block 6 Lot 9 Camarin Caloocan City",
				CAdd:      "2863",
				UserID:    "2863",
				SOFID:     "1",
				Tin:       "",
				TpNo:      "",
				AgentCode: "01030063",
			},
			want: &CebAddClientResp{
				Code:    0,
				Message: "Successful",
				Result: RUResult{
					ResultStatus: "Successful",
					MessageID:    0,
					LogID:        0,
					ClientID:     3673,
					ClientNo:     "EWFHM0000070155828",
				},
				RemcoID: 9,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.CebAddClient(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("CebAddClient() error = %v, wantErr %v", err, test.wantErr)
			}
			var newReq CebAddClientReq
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

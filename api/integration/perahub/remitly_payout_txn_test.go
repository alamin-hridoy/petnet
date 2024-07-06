package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRMPayout(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	req := RMPayoutRequest{
		ClientRefNo:   "cl-ref",
		RefNo:         "ref-no",
		LocationCode:  "1",
		LocationID:    "1",
		LocationName:  "loc-name",
		Gender:        "M",
		ControlNo:     "ctrl-no",
		CurrencyCode:  "4",
		PnplAmt:       "14",
		IDNumber:      "12",
		IDType:        "id-type",
		IDIssBy:       "id-iss-by",
		IDIssueDate:   "iss-date",
		IDExpDate:     "exp-date",
		ContactNumber: "54321",
		Address:       "addr",
		City:          "city",
		Province:      "prov",
		Country:       "country",
		ZipCode:       "12345",
		State:         "state",
		Natl:          "natl",
		BirthDate:     "birty",
		BirthCountry:  "bcountry",
		Occupation:    "occ",
		UserID:        "2",
		TxnDate:       "txn-date",
		CustomerID:    "3",
		CurrencyID:    "4",
		RemcoID:       "5",
		TxnType:       "6",
		IsDomestic:    "7",
		CustomerName:  "cust-name",
		ServiceCharge: "8",
		RmtLocID:      "9",
		DstAmount:     "10",
		TotalAmount:   "11",
		BBAmt:         "bb-amt",
		MCRateID:      "rateid",
		MCRate:        "rate",
		RmtIPAddr:     "127.0.0.1",
		RmtUserID:     "12",
		OrgnCtry:      "orgn-ctry",
		DestCtry:      "dest-ctry",
		PurposeTxn:    "prps",
		SourceFund:    "src",
		RelationTo:    "rel",
		BirthPlace:    "bplace",
		Barangay:      "barang",
		RiskScore:     "1",
		RiskCriteria:  "1",
		FormType:      "f-type",
		FormNumber:    "f-no",
		PayoutType:    "13",
		SenderName:    "send",
		RcvName:       "rcv",
		SenderFName:   "s-fn",
		SenderMName:   "s-mn",
		SenderLName:   "s-ln",
		RcvFName:      "r-fn",
		RcvMName:      "r-mn",
		RcvLName:      "r-ln",
		AgentID:       "a-id",
		AgentCode:     "a-code",
		IPAddr:        "ip-addr",
		RateCat:       "rate-cat",
		DsaCode:       "TEST_DSA",
		DsaTrxType:    "digital",
	}

	tests := []struct {
		name        string
		in          RMPayoutRequest
		expectedReq RMPayoutRequest
		want        *RMPayoutResponseBody
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &RMPayoutResponseBody{
				Code: "200",
				Msg:  "Success",
				Result: RMPayResult{
					RefNo:      "REF1",
					Created:    "2021-10-21",
					State:      "PAID",
					Type:       "CASH_PICKUP",
					PayerCodes: "",
				},
				RemcoID: "21",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.RMPayout(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("RMPayout() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq RMPayoutRequest
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

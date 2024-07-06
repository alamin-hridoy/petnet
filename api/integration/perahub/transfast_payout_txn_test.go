package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTFPayout(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	req := TFPayoutRequest{
		RiskScore:       "1",
		RiskCriteria:    "1",
		LocationID:      "1",
		LocationName:    "loc-name",
		UserID:          "2",
		TxnDate:         "txn-date",
		CustomerID:      "3",
		CurrencyID:      "4",
		RemcoID:         "5",
		TxnType:         "6",
		IsDomestic:      "7",
		CustomerName:    "cust-name",
		ServiceCharge:   "8",
		RmtLocID:        "9",
		DstAmount:       "10",
		TotalAmount:     "11",
		RmtUserID:       "12",
		RmtIPAddr:       "127.0.0.1",
		PurposeTxn:      "prps",
		SourceFund:      "src",
		Occupation:      "occ",
		RelationTo:      "rel",
		BirthDate:       "birty",
		BirthPlace:      "bplace",
		BirthCountry:    "bcountry",
		IDType:          "id-type",
		IDNumber:        "id-no",
		Address:         "addr",
		Barangay:        "barang",
		City:            "city",
		Province:        "prov",
		ZipCode:         "12345",
		Country:         "country",
		ContactNumber:   "54321",
		PayoutType:      "13",
		SenderName:      "send",
		RcvName:         "rcv",
		PnplAmt:         "14",
		ControlNo:       "ctrl-no",
		RefNo:           "ref-no",
		ClientRefNo:     "cl-ref",
		RcvOccupationID: "15",
		RcvStateID:      "stateid",
		RcvCityID:       "16",
		KYCVerified:     true,
		IDExpDate:       "expdate",
		Gender:          "M",
		RmtReasonID:     "17",
		ProofOfAddress:  "proof",
		DsaCode:         "TEST_DSA",
		DsaTrxType:      "digital",
	}

	tests := []struct {
		name        string
		in          TFPayoutRequest
		expectedReq TFPayoutRequest
		want        *TFPayoutResponseBody
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &TFPayoutResponseBody{
				Code:    "200",
				Msg:     "Successful.",
				RemcoID: "7",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s, m := newTestSvc(t, st)
			got, err := s.TFPayout(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("TFPayout() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq TFPayoutRequest
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

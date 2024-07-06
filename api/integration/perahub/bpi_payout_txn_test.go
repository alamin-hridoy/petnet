package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBPPayout(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	req := BPPayoutRequest{
		BBAmt:         "0",
		MCRate:        "0.00",
		RateCat:       "2",
		MCRateID:      "3",
		ControlNo:     "06PNET0915211119",
		LocationID:    "371",
		UserID:        "5188",
		TxnDate:       "2021-10-21",
		LocationName:  "work at home",
		CustomerID:    "332",
		CurrencyID:    "1",
		RemcoID:       "2",
		TxnType:       "2",
		IsDomestic:    "0",
		CustomerName:  "Juan Dela Cruz",
		ServiceCharge: "0.00",
		RmtLocID:      "1",
		DstAmount:     "0.00",
		TotalAmount:   "1000.00",
		RmtUserID:     "155",
		OrgnCtry:      "PH",
		DestCtry:      "PH",
		PurposeTxn:    "0",
		SourceFund:    "none",
		Occupation:    "none",
		RelationTo:    "none",
		BirthDate:     "1992-12-21",
		BirthPlace:    "none",
		BirthCountry:  "none",
		IDType:        "passport",
		IDNumber:      "666",
		Address:       "San ka na street",
		Barangay:      "Pio del Pilar",
		City:          "Makati",
		Province:      "NCR",
		ZipCode:       "422222",
		Country:       "PH",
		ContactNumber: "1976485",
		CurAdd:        "Unknown",
		PerAdd:        "Unknown",
		RiskScore:     "4",
		RiskCriteria:  "none",
		FormType:      "OAR",
		FormNumber:    "999",
		PayoutType:    "1",
		SenderName:    "Ako si Sender Jr",
		RcvName:       "Siya naman si Receiver",
		PnplAmt:       "1000.00",
		ClientRefNo:   "123445688",
		RefNo:         "06PNET091521111909120602",
		DsaCode:       "TEST_DSA",
		DsaTrxType:    "digital",
	}

	tests := []struct {
		name        string
		in          BPPayoutRequest
		expectedReq BPPayoutRequest
		want        *BPPayoutResponseBody
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &BPPayoutResponseBody{
				Code: "200",
				Msg:  "Success",
				Result: BPPayoutResult{
					Status:            "T",
					Desc:              "TRANSMIT",
					ControlNo:         "CTRL1",
					RefNo:             "1",
					ClientReferenceNo: "CL1",
					PnplAmt:           "1000.00",
					SenderName:        "MARLON GAVINO REYES VILLA",
					RcvName:           "DANIEL, JENI",
					Address:           "null",
					Currency:          "null",
					ContactNumber:     "null",
					RcvLastName:       "null",
					RcvFirstName:      "null",
					OrgnCtry:          "SINGAPORE",
					DestCtry:          "PHILIPPINES",
					TxnDate:           "null",
					IsDomestic:        "null",
					IDType:            "null",
					RcvCtryCode:       "null",
					RcvStateID:        "null",
					RcvStateName:      "null",
					RcvCityID:         "null",
					RcvCityName:       "null",
					RcvIDType:         "null",
					RcvIsIndiv:        "null",
					PrpsOfRmtID:       "null",
					DsaCode:           "TEST_DSA",
					DsaTrxType:        "digital",
				},
				RemcoID: "2",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s, m := newTestSvc(t, st)
			got, err := s.BPPayout(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("BPPayout() error = %v, wantErr %v", err, test.wantErr)
			}
			var newReq BPPayoutRequest
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

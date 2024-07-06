package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCEBPayout(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	req := CEBPayoutRequest{
		LocationID:    "191",
		LocationName:  "MALOLOS",
		UserID:        "5188",
		TxnDate:       "2021-10-06",
		CustomerID:    "6276847",
		CurrencyID:    "1",
		BeneficiaryID: "5342",
		IDTypeID:      "1",
		RemcoID:       "9",
		TxnType:       "2",
		IsDomestic:    "1",
		CustomerName:  "ALVIN CAMITAN",
		ServiceCharge: "90",
		RmtLocID:      "191",
		DstAmount:     "0",
		TotalAmount:   "101",
		RmtUserID:     "5188",
		RmtIpADD:      "130.211.2.203",
		OrgnCtry:      "Philippines",
		DestCtry:      "PH",
		PurposeTxn:    "Family Support/Living Expenses",
		SourceFund:    "Savings",
		Occupation:    "Unemployed",
		RelationTo:    "Family",
		BirthDate:     "1980-08-10T00:00:00",
		BirthPlace:    "ERMITA,BOHOL",
		BirthCountry:  "Philippines",
		IDType:        "445",
		IDNumber:      "ASD292929291",
		Address:       "18 SITIO PULO",
		Barangay:      "BULIHAN",
		City:          "MALOLOS",
		Province:      "BULACAN",
		ZipCode:       "3000A",
		Country:       "Philippines",
		ContactNumber: "09065595959",
		CurrentAddress: NonexAddress{
			Address1: "Block 1 Lot 12",
			Address2: "Placid Homes",
			Barangay: "BULIHAN",
			City:     "ANDA",
			Province: "BOHOL",
			ZipCode:  "1000A",
			Country:  "Philippines",
		},
		PermanentAddress: NonexAddress{
			Address1: "Block 1 Lot 12",
			Address2: "Placid Homes",
			Barangay: "BULIHAN",
			City:     "ANDA",
			Province: "BOHOL",
			ZipCode:  "1000A",
			Country:  "Philippines",
		},
		RiskScore:     "1",
		PayoutType:    "1",
		SenderName:    "REYES, MARK JEFFREY, ALOS",
		RcvName:       "ESTOCAPIO, SHAIRA MIKA, MADJALIS",
		PnplAmt:       "3000",
		ClientRefNo:   "10095",
		ControlNo:     "P002AYNRRTPT",
		RefNo:         "21181472909563000012",
		BuyBackAmount: "1",
		MCRate:        "123",
		RateCat:       "1",
		MCRateID:      "1",
		FormType:      "1",
		FormNumber:    "1",
		DsaCode:       "TEST_DSA",
		DsaTrxType:    "digital",
	}

	tests := []struct {
		name        string
		in          CEBPayoutRequest
		expectedReq CEBPayoutRequest
		want        *CEBPayoutResponseBody
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &CEBPayoutResponseBody{
				Code:    "0",
				Message: "Successful",
				Result: CEBPayoutResult{
					Message: "Successfully Payout!",
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
			got, err := s.CEBPayout(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("CEBPayout() error = %v, wantErr %v", err, test.wantErr)
			}
			var newReq CEBPayoutRequest
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

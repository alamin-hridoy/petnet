package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAYANNAHPayout(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	req := AYANNAHPayoutRequest{
		LocationID:         "191",
		LocationName:       "MALOLOS",
		UserID:             "5188",
		TrxDate:            "2021-07-13",
		CustomerID:         "7712780",
		CurrencyID:         "1",
		RemcoID:            "22",
		TrxType:            "2",
		IsDomestic:         "1",
		CustomerName:       "Cortez, Fernando, gre",
		ServiceCharge:      "0",
		RemoteLocationID:   "191",
		DstAmount:          "0",
		TotalAmount:        "3300",
		RemoteUserID:       "1893",
		RemoteIPAddress:    "130.211.2.203",
		PurposeTransaction: "Family Support/Living Expenses",
		SourceFund:         "Salary/Income",
		Occupation:         "Unemployed",
		RelationTo:         "Family",
		BirthDate:          "2000-12-15",
		BirthPlace:         "MALOLOS,BULACANqdasd",
		BirthCountry:       "Philippines",
		IDType:             "PASSPORT",
		IDNumber:           "PRND32200265569Pdasd3",
		Address:            "18 SITIO PULO",
		Barangay:           "BULIHAN",
		City:               "MALOLOS",
		Province:           "BULACAN",
		ZipCode:            "3000A",
		Country:            "Philippines",
		ContactNumber:      "09265454935",
		RiskScore:          "1",
		PayoutType:         "1",
		SenderName:         "Mercado, Marites Cueto wesf",
		ReceiverName:       "Cortez, Fernando wesfes",
		PrincipalAmount:    "3300",
		ControlNumber:      "RA5MDDXYW7RHef",
		ReferenceNumber:    "85bf86fe-7723-42c8-b511-8bd46270a699",
		McRate:             "1",
		McRateID:           "1",
		RateCategory:       "required",
		OriginatingCountry: "Philippines",
		DestinationCountry: "Philippines",
		FormType:           "0",
		FormNumber:         "0",
		ClientReferenceNo:  "ccea7bc90e207c0016e6",
		BuyBackAmount:      "1",
		IPAddress:          "::1",
		DsaCode:            "TEST_DSA",
		DsaTrxType:         "digital",
	}

	tests := []struct {
		name        string
		in          AYANNAHPayoutRequest
		expectedReq AYANNAHPayoutRequest
		want        *AYANNAHPayoutResponseBody
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &AYANNAHPayoutResponseBody{
				Code:    "200",
				Message: "Success",
				Result: AYANNAHPayoutResult{
					Message: "Successfully Payout.",
				},
				RemcoID: "22",
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s, m := newTestSvc(t, st)
			got, err := s.AYANNAHPayout(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("AYANNAHPayout() error = %v, wantErr %v", err, test.wantErr)
			}
			var newReq AYANNAHPayoutRequest
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

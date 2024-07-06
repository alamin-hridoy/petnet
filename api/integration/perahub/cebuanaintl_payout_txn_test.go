package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCEBINTPayout(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	req := CEBINTPayoutRequest{
		LocationID:         "191",
		LocationName:       "MALOLOS",
		UserID:             "5188",
		TrxDate:            "2021-10-06",
		CustomerID:         "6276847",
		CurrencyID:         "1",
		RemcoID:            "19",
		TrxType:            "2",
		IsDomestic:         "0",
		CustomerName:       "ALVIN CAMITAN",
		ServiceCharge:      "1",
		RemoteLocationID:   "191",
		DstAmount:          "0",
		TotalAmount:        "1001.00",
		RemoteUserID:       "5188",
		RemoteIPAddress:    "130.211.2.203",
		OriginatingCountry: "Malaysia",
		DestinationCountry: "PH",
		PurposeTransaction: "Family Support/Living Expenses",
		SourceFund:         "Savings",
		Occupation:         "Unemployed",
		RelationTo:         "Family",
		BirthDate:          "1995-09-06",
		BirthPlace:         "ERMITA,BOHOL",
		BirthCountry:       "Philippines",
		IDType:             "NBI Clearance",
		IDNumber:           "20337019",
		Address:            "18 SITIO PULO",
		Barangay:           "BULIHAN",
		City:               "MALOLOS",
		Province:           "BULACAN",
		ZipCode:            "3000A",
		Country:            "Philippines",
		ContactNumber:      "CTRL1",
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
		RiskScore:                "1",
		RiskCriteria:             "1",
		FormType:                 "0",
		FormNumber:               "0",
		PayoutType:               "1",
		SenderName:               "Mercado, Marites Cueto",
		ReceiverName:             "Cortez, Fernando",
		PrincipalAmount:          "1000.00",
		ClientReferenceNo:        "1",
		ControlNumber:            "QZ13VCE12",
		ReferenceNumber:          "REF1",
		IdentificationTypeID:     "11",
		BeneficiaryID:            "8540",
		AgentCode:                "01030063",
		IDIssuedBy:               "PH",
		IDIssuedState:            "null",
		IDIssuedCity:             "TANZA",
		IDDateOfIssue:            "2021-02-11T00:00:00",
		IDIssuingCountry:         "Philippines",
		IDIssuingState:           "null",
		PassportIDIssuedCountry:  "166",
		InternationalPartnerCode: "PL0005",
		McRate:                   "0",
		BuyBackAmount:            "0",
		RateCategory:             "0",
		McRateID:                 "0",
		DsaCode:                  "TEST_DSA",
		DsaTrxType:               "digital",
	}
	tests := []struct {
		name        string
		in          CEBINTPayoutRequest
		expectedReq CEBINTPayoutRequest
		want        *CEBINTPayoutResponseBody
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &CEBINTPayoutResponseBody{
				Code:    "0",
				Message: "Successful",
				Result: CEBINTPayoutResult{
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
			got, err := s.CEBINTPayout(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("CEBINTPayout() error = %v, wantErr %v", err, test.wantErr)
			}
			var newReq CEBINTPayoutRequest
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

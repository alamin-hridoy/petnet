package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCebuanaSendMoney(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	req := CebuanaSendRequest{
		LocationID:         "191",
		LocationName:       "MALOLOS",
		UserID:             5188,
		TrxDate:            "2021-10-06",
		CustomerID:         "6276847",
		CurrencyID:         "1",
		RemcoID:            "9",
		TrxType:            "1",
		IsDomestic:         "1",
		CustomerName:       "ALVIN CAMITAN",
		ServiceCharge:      "1",
		RemoteLocationID:   "191",
		DstAmount:          "0",
		TotalAmount:        "101",
		RemoteUserID:       "5188",
		RemoteIPAddress:    "130.211.2.203",
		OriginatingCountry: "Philippines",
		DestinationCountry: "PH",
		PurposeTransaction: "Family Support/Living Expenses",
		SourceFund:         "Savings",
		Occupation:         "Unemployed",
		RelationTo:         "Family",
		BirthDate:          "1995-09-06",
		BirthPlace:         "ERMITA,BOHOL",
		BirthCountry:       "Philippines",
		IDType:             "445",
		IDNumber:           "ASD292929291",
		Address:            "18 SITIO PULO",
		Barangay:           "BULIHAN",
		City:               "MALOLOS",
		Province:           "BULACAN",
		ZipCode:            "3000A",
		Country:            "Philippines",
		ContactNumber:      "09065595959",
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
		RiskScore:       "1",
		FormType:        "OAR",
		FormNumber:      "MAL0101011",
		BeneficiaryID:   "5342",
		AgentCode:       "01030063",
		McRateID:        "1",
		RateCtg:         "1",
		BbAmount:        "100",
		McRate:          "1",
		SendCurrencyID:  "6",
		PrincipalAmount: "3000",
		SenderName:      "Aneek Khan",
		ReceiverName:    "Kamuzzaman",
		ClientRefNo:     "5689664",
		PayoutType:      "1",
	}

	tests := []struct {
		name        string
		in          CebuanaSendRequest
		expectedReq CebuanaSendRequest
		want        *CebuanaSendResponseBody
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &CebuanaSendResponseBody{
				Code:    "0",
				Message: "Successful",
				Result: CebuanaSendResult{
					ResultStatus: "Successful",
					MessageID:    "0",
					LogID:        "0",
					ControlNo:    "CTRL1",
					ServiceFee:   "1.00",
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
			got, err := s.CebuanaSendMoney(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("CebuanaSendMoney() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq CebuanaSendRequest
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

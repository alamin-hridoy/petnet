package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUSSCSendMoney(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	req := USSCSendRequest{
		ControlNo:          "random2223",
		McRate:             "0",
		BbAmount:           "0",
		RateCtg:            "0",
		McRateID:           "0",
		BranchCode:         "branch1",
		LocationID:         "191",
		LocationName:       "MALOLOS",
		UserID:             "5188",
		TrxDate:            "2021-10-06",
		CustomerID:         "6276847",
		CurrencyID:         "1",
		RemcoID:            "10",
		TrxType:            "2",
		IsDomestic:         "1",
		CustomerName:       "ALVIN CAMITA",
		ServiceCharge:      "60.00",
		RemoteLocationID:   "191",
		DstAmount:          "0",
		TotalAmount:        "3484.33",
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
		RiskCrt:         "1",
		FormType:        "OR",
		FormNumber:      "MAL00001",
		PayoutType:      "1",
		SenderName:      "JOMAR TE TEST",
		ReceiverName:    "ALVIN CAMITAN",
		PrincipalAmount: "3424.33",
		ClientRefNo:     "b41c69b5b0fc5e3cf7cb3",
		ReferenceNo:     "20211221PHB065653289",
		SendFName:       "John",
		SendMName:       "Michael",
		SendLName:       "Doe",
		RecFName:        "Jane",
		RecMName:        "Emily",
		RecLName:        "Doe",
		RecConNo:        "0922261616161",
		KycVer:          true,
		Gender:          "M",
	}

	tests := []struct {
		name        string
		in          USSCSendRequest
		expectedReq USSCSendRequest
		want        *USSCSendResponseBody
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &USSCSendResponseBody{
				Code:    "000000",
				Message: "OK",
				Result: USSCSendResult{
					ControlNo:          "CTRL1",
					TrxDate:            "2021-12-22",
					SendFName:          "John",
					SendMName:          "Michael",
					SendLName:          "Doe",
					PrincipalAmount:    "1000.00",
					ServiceCharge:      "1.00",
					TotalAmount:        "1001.00",
					RecFName:           "ALVIN",
					RecMName:           "JOMAR TE TEST",
					RecLName:           "CAMITAN",
					ContactNumber:      "0922261616161",
					RelationTo:         "Family",
					PurposeTransaction: "Family Support/Living Expenses",
					ReferenceNo:        "1",
				},
				RemcoID: "10",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.USSCSendMoney(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("USSCSendMoney() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq USSCSendRequest
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
